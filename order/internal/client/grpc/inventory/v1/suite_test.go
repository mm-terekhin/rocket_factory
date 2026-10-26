package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	inventoryV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
)

type inventoryGRPCMock struct {
	mock.Mock
}

func (m *inventoryGRPCMock) GetPart(ctx context.Context, in *inventoryV1.GetPartRequest, _ ...grpc.CallOption) (*inventoryV1.GetPartResponse, error) {
	panic("unexpected call GetPart")
}

func (m *inventoryGRPCMock) ListParts(ctx context.Context, in *inventoryV1.ListPartsRequest, _ ...grpc.CallOption) (*inventoryV1.ListPartsResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*inventoryV1.ListPartsResponse), args.Error(1)
}

func (m *inventoryGRPCMock) CreatePart(ctx context.Context, in *inventoryV1.CreatePartRequest, _ ...grpc.CallOption) (*inventoryV1.CreatePartResponse, error) {
	panic("unexpected call CreatePart")
}

type ClientSuite struct {
	suite.Suite

	ctx    context.Context
	grpc   *inventoryGRPCMock
	client *client
}

func (s *ClientSuite) SetupTest() {
	s.ctx = context.Background()
	s.grpc = new(inventoryGRPCMock)
	s.client = NewClient(s.grpc)
}

func TestClientIntegration(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}
