package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	paymentV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const grpcPort = 50052

type paymentService struct {
	paymentV1.UnimplementedPaymentServiceServer
}

func (c *paymentService) PayOrder(_ context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	isValid := validate(req)

	if isValid != nil {
		return nil, isValid
	}

	transaction_uuid := uuid.New().String()

	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transaction_uuid)

	return &paymentV1.PayOrderResponse{TransactionUuid: transaction_uuid}, nil

}

func validate(req *paymentV1.PayOrderRequest) error {
	if req.GetOrderUuid() == "" {
		return status.Error(codes.InvalidArgument, "order_uuid is required")
	}

	if req.GetUserUuid() == "" {
		return status.Error(codes.InvalidArgument, "user_uuid is required")
	}

	if req.GetPaymentMethod() == paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED {
		return status.Error(codes.InvalidArgument, "payment_method is required")
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("Failed to close listener %v\n", cerr)
		}
	}()

	s := grpc.NewServer()
	service := &paymentService{}
	paymentV1.RegisterPaymentServiceServer(s, service)

	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err := s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
