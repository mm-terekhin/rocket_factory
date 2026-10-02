package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net"
	"os"
	"os/signal"
	"slices"
	"strings"
	"sync"
	"syscall"

	inventoryV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const grpcPort = 50051

type inventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer

	mu    sync.RWMutex
	parts map[string]*inventoryV1.Part
}

func (s *inventoryService) CreatePart(_ context.Context, req *inventoryV1.CreatePartRequest) (*inventoryV1.CreatePartResponse, error) {
	if req.GetInfo() == nil {
		return nil, status.Error(codes.InvalidArgument, "info is required")
	}

	id := uuid.New()
	now := timestamppb.Now()

	part := &inventoryV1.Part{
		Uuid: id.String(),
		Info: &inventoryV1.PartInfo{
			Name:          req.GetInfo().GetName(),
			Description:   req.GetInfo().GetDescription(),
			Price:         req.GetInfo().GetPrice(),
			StockQuantity: req.GetInfo().GetStockQuantity(),
			Category:      req.GetInfo().GetCategory(),
			Dimensions:    req.GetInfo().GetDimensions(),
			Manufacturer:  req.GetInfo().GetManufacturer(),
			Tags:          req.GetInfo().GetTags(),
			Metadata:      req.GetInfo().GetMetadata(),
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.parts[part.Uuid] = part

	return &inventoryV1.CreatePartResponse{Uuid: part.Uuid}, nil

}

func (s *inventoryService) GetPart(_ context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	part, ok := s.parts[req.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.GetUuid())
	}

	return &inventoryV1.GetPartResponse{Part: part}, nil
}

func (s *inventoryService) ListParts(_ context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	all := s.collectAllParts()

	filter := req.GetFilter()
	if filter == nil {
		return &inventoryV1.ListPartsResponse{Parts: sortedParts(all)}, nil

	}

	pf := prepareFilter(filter)

	out := make([]*inventoryV1.Part, 0, len(all))
	for _, part := range all {
		if matches(part, pf) {
			out = append(out, part)
		}
	}

	return &inventoryV1.ListPartsResponse{Parts: sortedParts(out)}, nil
}

func sortedParts(parts []*inventoryV1.Part) []*inventoryV1.Part {
	out := append([]*inventoryV1.Part(nil), parts...)
	slices.SortFunc(out, func(a, b *inventoryV1.Part) int {
		return strings.Compare(a.GetUuid(), b.GetUuid())
	})
	return out
}

type preparedFilter struct {
	uuids      map[string]struct{}
	names      map[string]struct{}
	categories map[inventoryV1.Category]struct{}
	countries  map[string]struct{}
	tags       map[string]struct{}
}

func prepareFilter(f *inventoryV1.PartsFilter) preparedFilter {
	pf := preparedFilter{}

	if len(f.GetUuids()) > 0 {
		pf.uuids = toStringSet(f.GetUuids())
	}

	if len(f.GetNames()) > 0 {
		pf.names = toStringSet(f.GetNames())
	}

	if len(f.GetCategories()) > 0 {
		pf.categories = toCategorySet(f.GetCategories())
	}

	if len(f.GetManufacturerCountries()) > 0 {
		pf.countries = toStringSet(f.GetManufacturerCountries())
	}

	if len(f.GetTags()) > 0 {
		pf.tags = toStringSet(f.GetTags())
	}

	return pf
}

func toStringSet(items []string) map[string]struct{} {
	m := make(map[string]struct{}, len(items))

	for _, s := range items {
		m[s] = struct{}{}
	}

	return m
}

func toCategorySet(items []inventoryV1.Category) map[inventoryV1.Category]struct{} {
	m := make(map[inventoryV1.Category]struct{}, len(items))

	for _, s := range items {
		m[s] = struct{}{}
	}

	return m
}

func matches(part *inventoryV1.Part, pf preparedFilter) bool {
	if len(pf.uuids) > 0 {
		if _, ok := pf.uuids[part.GetUuid()]; !ok {
			return false
		}
	}

	if len(pf.names) > 0 {
		if _, ok := pf.names[part.GetInfo().GetName()]; !ok {
			return false
		}
	}

	if len(pf.categories) > 0 {
		if _, ok := pf.categories[part.GetInfo().GetCategory()]; !ok {
			return false
		}
	}

	if len(pf.countries) > 0 {
		if _, ok := pf.countries[part.GetInfo().GetManufacturer().GetCountry()]; !ok {
			return false
		}
	}

	if len(pf.tags) > 0 {
		matched := false
		for _, t := range part.GetInfo().GetTags() {
			if _, ok := pf.tags[t]; ok {
				matched = true
				break
			}
		}

		if !matched {
			return false
		}
	}

	return true
}

func (s *inventoryService) collectAllParts() []*inventoryV1.Part {
	parts := make([]*inventoryV1.Part, 0, len(s.parts))

	for _, p := range s.parts {
		parts = append(parts, p)
	}

	return parts
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

	service := &inventoryService{
		parts: make(map[string]*inventoryV1.Part),
	}

	inventoryV1.RegisterInventoryServiceServer(s, service)

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
