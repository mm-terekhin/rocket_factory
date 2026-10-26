package inventory

import (
	"github.com/mm-terekhin/rocket_factory/inventory/internal/repository"
	def "github.com/mm-terekhin/rocket_factory/inventory/internal/service"
)

var _ def.InventoryService = (*service)(nil)

type service struct {
	InventoryRepository repository.InventoryRepository
}

func NewService(repository repository.InventoryRepository) *service {
	return &service{
		InventoryRepository: repository,
	}
}
