package v1

import (
	"context"
	"testing"

	"github.com/mm-terekhin/rocket_factory/inventory/internal/service/mocks"
	"github.com/stretchr/testify/suite"
)

type APISuite struct {
	suite.Suite

	ctx context.Context

	inventoryService *mocks.InventoryService

	api *inventoryApi
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()

	s.inventoryService = mocks.NewInventoryService(s.T())

	s.api = NewApi(s.inventoryService)
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
