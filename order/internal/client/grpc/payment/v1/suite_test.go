package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	paymentV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
)

type paymentGRPCMock struct {
	mock.Mock
}

func (m *paymentGRPCMock) PayOrder(ctx context.Context, in *paymentV1.PayOrderRequest, _ ...grpc.CallOption) (*paymentV1.PayOrderResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*paymentV1.PayOrderResponse), args.Error(1)
}

type ClientSuite struct {
	suite.Suite

	ctx    context.Context
	grpc   *paymentGRPCMock
	client *client
}

func (s *ClientSuite) SetupTest() {
	s.ctx = context.Background()
	s.grpc = new(paymentGRPCMock)
	s.client = NewClient(s.grpc)
}

func TestClientIntegration(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}
