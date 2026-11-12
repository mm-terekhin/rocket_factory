package inventory

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	def "github.com/mm-terekhin/rocket_factory/inventory/internal/repository"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	parts *mongo.Collection
}

func NewRepository(db *mongo.Database) *repository {
	collection := db.Collection("parts")

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "uuid", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "info.name", Value: 1}},
			Options: options.Index().SetUnique(false),
		},
		{
			Keys:    bson.D{{Key: "info.category", Value: 1}},
			Options: options.Index().SetUnique(false),
		},
		{
			Keys:    bson.D{{Key: "info.manufacturer.country", Value: 1}},
			Options: options.Index().SetUnique(false),
		},
		{
			Keys:    bson.D{{Key: "info.tags", Value: 1}},
			Options: options.Index().SetUnique(false),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		panic(err)
	}

	return &repository{
		parts: collection,
	}
}
