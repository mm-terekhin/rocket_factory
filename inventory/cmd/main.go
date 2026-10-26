package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	inventoryApi "github.com/mm-terekhin/rocket_factory/inventory/internal/api/inventory/v1"
	inventoryRepository "github.com/mm-terekhin/rocket_factory/inventory/internal/repository/inventory"
	inventoryService "github.com/mm-terekhin/rocket_factory/inventory/internal/service/inventory"
	inventoryV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50051

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

	repo := inventoryRepository.NewRepository()
	service := inventoryService.NewService(repo)
	api := inventoryApi.NewApi(service)

	inventoryV1.RegisterInventoryServiceServer(s, api)

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
