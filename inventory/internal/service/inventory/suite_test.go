package inventory

import (
	"context"
	"testing"

	"github.com/mm-terekhin/rocket_factory/inventory/internal/repository/mocks"
	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	inventoryRepository *mocks.InventoryRepository

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.inventoryRepository = mocks.NewInventoryRepository(s.T())

	s.service = NewService(s.inventoryRepository)
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
