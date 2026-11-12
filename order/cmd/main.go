package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderApi "github.com/mm-terekhin/rocket_factory/order/internal/api/order/v1"
	invClient "github.com/mm-terekhin/rocket_factory/order/internal/client/grpc/inventory/v1"
	payClient "github.com/mm-terekhin/rocket_factory/order/internal/client/grpc/payment/v1"
	"github.com/mm-terekhin/rocket_factory/order/internal/migrator"
	orderRepository "github.com/mm-terekhin/rocket_factory/order/internal/repository/order"
	orderService "github.com/mm-terekhin/rocket_factory/order/internal/service/order"
	orderV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/proto/payment/v1"
)

const (
	inventoryAddress  = "localhost:50051"
	paymentAddress    = "localhost:50052"
	httpPort          = "8080"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		return err
	}

	dbURI := os.Getenv("DB_URI")

	invConn, err := grpc.NewClient(
		inventoryAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("inventory grpc connect: %w", err)
	}
	defer func() {
		if cerr := invConn.Close(); cerr != nil {
			log.Printf("failed to close inventory connection: %v\n", cerr)
		}
	}()

	inventory := invClient.NewClient(inventoryV1.NewInventoryServiceClient(invConn))

	payConn, err := grpc.NewClient(
		paymentAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("payment grpc connect: %w", err)
	}
	defer func() {
		if cerr := payConn.Close(); cerr != nil {
			log.Printf("failed to close payment connection: %v\n", cerr)
		}
	}()

	payment := payClient.NewClient(paymentV1.NewPaymentServiceClient(payConn))

	con, err := pgx.Connect(ctx, dbURI)
	if err != nil {
		return err
	}
	defer func() {
		cerr := con.Close(ctx)
		if cerr != nil {
			log.Printf("failed to close connection: %v\n", cerr)
		}
	}()

	err = con.Ping(ctx)
	if err != nil {
		return err
	}

	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	migratorRunner := migrator.NewMigrator(stdlib.OpenDB(*con.Config().Copy()), migrationsDir)

	err = migratorRunner.Up()
	if err != nil {
		return err
	}

	repo := orderRepository.NewRepository(con)
	service := orderService.NewService(repo, inventory, payment)
	api := orderApi.NewApi(service)

	ordersServer, err := orderV1.NewServer(api)
	if err != nil {
		return fmt.Errorf("openapi server: %w", err)
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

	return nil
}
