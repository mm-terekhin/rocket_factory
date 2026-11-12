package inventory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/mm-terekhin/rocket_factory/inventory/internal/repository/mocks"
)

func testCtx() context.Context {
	return context.Background()
}

type ServiceSuite struct {
	suite.Suite

	inventoryRepository *mocks.InventoryRepository

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.inventoryRepository = mocks.NewInventoryRepository(s.T())

	s.service = NewService(s.inventoryRepository)
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
