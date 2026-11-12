package inventory

import (
	"context"
	"log"

	"github.com/mm-terekhin/rocket_factory/inventory/internal/model"
	"github.com/mm-terekhin/rocket_factory/inventory/internal/repository/converter"
	repoModel "github.com/mm-terekhin/rocket_factory/inventory/internal/repository/model"
	"go.mongodb.org/mongo-driver/bson"
)

type preparedFilter struct {
	uuids      map[string]struct{}
	names      map[string]struct{}
	categories map[model.Category]struct{}
	countries  map[string]struct{}
	tags       map[string]struct{}
}

func (r *repository) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	cursor, err := r.parts.Find(ctx, partsFilterToBSON(filter))
	if err != nil {
		return make([]model.Part, 0), err
	}
	defer func() {
		cerr := cursor.Close(ctx)
		if cerr != nil {
			log.Printf("failed to close cursor: %v\n", cerr)
		}
	}()

	var parts []repoModel.Part
	err = cursor.All(ctx, &parts)
	if err != nil {
		return nil, err
	}

	out := make([]model.Part, 0, len(parts))
	for _, part := range parts {
		partModel := converter.PartToModel(part)
		out = append(out, partModel)
	}

	return out, nil
}

func partsFilterToBSON(f model.PartsFilter) bson.M {
	query := bson.M{}

	if len(f.UUIDs) > 0 {
		query["uuid"] = bson.M{"$in": f.UUIDs}
	}
	if len(f.Names) > 0 {
		query["info.name"] = bson.M{"$in": f.Names}
	}
	if len(f.Categories) > 0 {
		cats := make([]int, 0, len(f.Categories))
		for _, c := range f.Categories {
			cats = append(cats, int(c))
		}
		query["info.category"] = bson.M{"$in": cats}
	}
	if len(f.ManufacturerCountries) > 0 {
		query["info.manufacturer.country"] = bson.M{"$in": f.ManufacturerCountries}
	}
	if len(f.Tags) > 0 {
		query["info.tags"] = bson.M{"$in": f.Tags}
	}

	return query
}
