package v1

import (
	"context"
	"testing"

	"github.com/mm-terekhin/rocket_factory/order/internal/service/mocks"
	"github.com/stretchr/testify/suite"
)

func testCtx() context.Context {
	return context.Background()
}

type APISuite struct {
	suite.Suite

	orderService *mocks.OrderService
	api          *orderApi
}

func (s *APISuite) SetupTest() {
	s.orderService = mocks.NewOrderService(s.T())
	s.api = NewApi(s.orderService)
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
