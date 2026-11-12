package repository

import (
	"context"

	"github.com/mm-terekhin/rocket_factory/inventory/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InventoryRepository interface {
	GetPart(context.Context, string) (model.Part, error)
	ListParts(context.Context, model.PartsFilter) ([]model.Part, error)
	CreatePart(context.Context, model.Part) (primitive.ObjectID, error)
}
