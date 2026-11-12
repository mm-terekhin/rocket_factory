package inventory

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/mm-terekhin/rocket_factory/inventory/internal/converter"
	"github.com/mm-terekhin/rocket_factory/inventory/internal/model"
)

func (r *repository) CreatePart(ctx context.Context, part model.Part) (primitive.ObjectID, error) {
	repoPart := converter.PartToRepo(part)

	res, err := r.parts.InsertOne(ctx, repoPart)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}
