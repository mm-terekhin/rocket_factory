package service

import (
	"context"

	"github.com/mm-terekhin/rocket_factory/inventory/internal/model"
)

type InventoryService interface {
	GetPart(context.Context, string) (model.Part, error)
	ListParts(context.Context, model.PartsFilter) ([]model.Part, error)
	CreatePart(context.Context, model.PartInfo) (string, error)
}
