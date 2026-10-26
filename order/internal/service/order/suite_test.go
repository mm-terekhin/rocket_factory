package order

import (
	"context"
	"testing"

	mocksClient "github.com/mm-terekhin/rocket_factory/order/internal/client/mocks"
	"github.com/mm-terekhin/rocket_factory/order/internal/repository/mocks"
	"github.com/stretchr/testify/suite"
)

func testCtx() context.Context {
	return context.Background()
}

type ServiceSuite struct {
	suite.Suite

	orderRepository *mocks.OrderRepository
	inventoryClient *mocksClient.InventoryClient
	paymentClient   *mocksClient.PaymentClient
	service         *service
}

func (s *ServiceSuite) SetupTest() {
	s.orderRepository = mocks.NewOrderRepository(s.T())
	s.inventoryClient = mocksClient.NewInventoryClient(s.T())
	s.paymentClient = mocksClient.NewPaymentClient(s.T())
	s.service = NewService(s.orderRepository, s.inventoryClient, s.paymentClient)
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
