package inventory

import (
	"sync"

	def "github.com/mm-terekhin/rocket_factory/inventory/internal/repository"
	repoModel "github.com/mm-terekhin/rocket_factory/inventory/internal/repository/model"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	mu    sync.RWMutex
	parts map[string]repoModel.Part
}

func NewRepository() *repository {
	return &repository{
		parts: make(map[string]repoModel.Part),
	}
}
