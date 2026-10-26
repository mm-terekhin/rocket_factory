package v1

import (
	"github.com/mm-terekhin/rocket_factory/inventory/internal/service"
	inventoryV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/proto/inventory/v1"
)

type inventoryApi struct {
	inventoryV1.UnimplementedInventoryServiceServer

	inventoryService service.InventoryService
}

func NewApi(inventoryService service.InventoryService) *inventoryApi {
	return &inventoryApi{inventoryService: inventoryService}
}
