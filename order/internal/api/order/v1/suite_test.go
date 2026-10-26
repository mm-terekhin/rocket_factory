package v1

import (
	"context"
	"testing"

	"github.com/mm-terekhin/rocket_factory/order/internal/service/mocks"
	"github.com/stretchr/testify/suite"
)

type APISuite struct {
	suite.Suite

	ctx context.Context

	orderService *mocks.OrderService
	api          *orderApi
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()
	s.orderService = mocks.NewOrderService(s.T())
	s.api = NewApi(s.orderService)
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
