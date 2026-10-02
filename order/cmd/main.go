package main

import (
	// "sync"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"uuid"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	orderV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const (
	inventoryAddress  = "localhost:50051"
	paymentAddress    = "localhost:50052"
	httpPort          = "8080"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type Status string

const (
	pending  Status = "PENDING_PAYMENT"
	paid     Status = "PAID"
	canceled Status = "CANCELLED"
)

type OrdersStorage struct {
	mu     sync.RWMutex
	orders map[string]*orderV1.Order
}

func NewOrdersStorage() *OrdersStorage {
	return &OrdersStorage{
		orders: make(map[string]*orderV1.Order),
	}
}

func (s *OrdersStorage) GetOrder(uuid string) *orderV1.Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.orders[uuid]

	if !ok {
		return nil
	}

	return val
}

func (s *OrdersStorage) UpdateOrder(order *orderV1.Order) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[order.GetOrderUUID()] = order
}

type OrdersHandler struct {
	storage   *OrdersStorage
	inventory inventoryV1.InventoryServiceClient
	payment   paymentV1.PaymentServiceClient
}

func (h *OrdersHandler) CancelOrderByUUID(_ context.Context, params orderV1.CancelOrderByUUIDParams) (orderV1.CancelOrderByUUIDRes, error) {
	order := h.storage.GetOrder(params.OrderUUID)

	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order for uuid'" + params.OrderUUID + "' not found",
		}, nil
	}

	if order.Status == string(paid) {
		return &orderV1.ConflictError{
			Code:    409,
			Message: "The order has already been paid for and cannot be cancelled.",
		}, nil
	}

	if order.Status != string(canceled) {
		order.Status = string(canceled)

		h.storage.UpdateOrder(order)
	}

	return nil, nil

}

func badRequest(message string) *orderV1.BadRequestError {
	return &orderV1.BadRequestError{
		Code:    400,
		Message: message,
	}
}

func (h *OrdersHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	if req.GetUserUUID() == "" {
		return badRequest("user_uuid is required"), nil
	}

	partUUIDs := req.GetPartUuids()
	if len(partUUIDs) == 0 {
		return badRequest("order has no parts"), nil
	}

	unique := make(map[string]struct{}, len(partUUIDs))
	for _, id := range partUUIDs {
		if id == "" {
			return badRequest("part_uuid must not be empty"), nil
		}
		unique[id] = struct{}{}
	}

	uniqueIDs := make([]string, 0, len(unique))
	for id := range unique {
		uniqueIDs = append(uniqueIDs, id)
	}

	resp, err := h.inventory.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{Uuids: uniqueIDs},
	})
	if err != nil {
		return nil, err
	}

	byUUID := make(map[string]*inventoryV1.Part, len(resp.GetParts()))
	for _, p := range resp.GetParts() {
		byUUID[p.GetUuid()] = p
	}

	var totalPrice float64
	for _, id := range partUUIDs {
		part, ok := byUUID[id]
		if !ok {
			return badRequest(fmt.Sprintf("part not found: %s", id)), nil
		}
		totalPrice += part.GetInfo().GetPrice()
	}

	orderUUID := uuid.New().String()
	h.storage.UpdateOrder(&orderV1.Order{
		OrderUUID:  orderUUID,
		UserUUID:   req.GetUserUUID(),
		PartUuids:  partUUIDs, // как пришло от клиента
		TotalPrice: totalPrice,
		Status:     string(pending),
	})
	return &orderV1.CreatedOrderResponse{
		OrderUUID:  orderV1.NewOptString(orderUUID),
		TotalPrice: orderV1.NewOptFloat64(totalPrice),
	}, nil
}

func (h *OrdersHandler) GetOrderByUUID(ctx context.Context, params orderV1.GetOrderByUUIDParams) (orderV1.GetOrderByUUIDRes, error) {
	order := h.storage.GetOrder(params.OrderUUID)

	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order for uuid'" + params.OrderUUID + "' not found",
		}, nil
	}

	return order, nil
}

func (h *OrdersHandler) PayOrderByUUID(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderByUUIDParams) (orderV1.PayOrderByUUIDRes, error) {
	order := h.storage.GetOrder(params.OrderUUID)

	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order for uuid'" + params.OrderUUID + "' not found",
		}, nil
	}

	method, ok := paymentMethodFromHTTP(req.GetPaymentMethod())
	if !ok {
		method = paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
	resp, err := h.payment.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     order.OrderUUID,
		UserUuid:      order.UserUUID,
		PaymentMethod: method,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.InvalidArgument {
			return badRequest(st.Message()), nil
		}
		return nil, err
	}

	order.TransactionUUID = orderV1.NewOptString(resp.GetTransactionUuid())
	order.Status = string(paid)
	h.storage.orders[order.OrderUUID] = order

	return &orderV1.PayOrderResponse{
		TransactionUUID: resp.GetTransactionUuid(),
	}, nil
}

func paymentMethodFromHTTP(s string) (paymentV1.PaymentMethod, bool) {
	switch s {
	case "CARD", "PAYMENT_METHOD_CARD":
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD, true
	case "SBP", "PAYMENT_METHOD_SBP":
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SBP, true
	case "CREDIT_CARD", "PAYMENT_METHOD_CREDIT_CARD":
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD, true
	case "INVESTOR_MONEY", "PAYMENT_METHOD_INVESTOR_MONEY":
		return paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY, true
	default:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED, false
	}
}

func NewOrdersHandler(storage *OrdersStorage, inventory inventoryV1.InventoryServiceClient, payment paymentV1.PaymentServiceClient) *OrdersHandler {
	return &OrdersHandler{
		storage:   storage,
		inventory: inventory,
		payment:   payment,
	}
}

func main() {
	storage := NewOrdersStorage()

	invConn, err := grpc.NewClient(
		inventoryAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect: %v\n", err)
	}
	defer func() {
		if cerr := invConn.Close(); cerr != nil {
			log.Printf("Failed to close listener %v\n", cerr)
		}
	}()

	inventory := inventoryV1.NewInventoryServiceClient(invConn)

	payConn, err := grpc.NewClient(
		paymentAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	payment := paymentV1.NewPaymentServiceClient(payConn)

	ordersHandler := NewOrdersHandler(storage, inventory, payment)

	ordersServer, err := orderV1.NewServer(ordersHandler)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Mount("/", ordersServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err := server.ListenAndServe()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")

}
