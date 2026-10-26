package converter

import (
	"log"
	"time"

	"github.com/mm-terekhin/rocket_factory/order/internal/model"
	generatedInventoryV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/proto/inventory/v1"
)

func PartListToModel(inventoryParts []*generatedInventoryV1.Part) []model.Part {
	parts := make([]model.Part, 0, len(inventoryParts))
	for _, part := range inventoryParts {
		parts = append(parts, PartToModel(part))
	}

	return parts
}

func PartToModel(inventoryPart *generatedInventoryV1.Part) model.Part {
	var updatedAt *time.Time
	if inventoryPart.UpdatedAt != nil {
		tmp := inventoryPart.UpdatedAt.AsTime()
		updatedAt = &tmp
	}

	return model.Part{
		UUID: inventoryPart.Uuid,
		Info: model.PartInfo{
			Name:          inventoryPart.Info.Name,
			Description:   inventoryPart.Info.Description,
			Price:         inventoryPart.Info.Price,
			StockQuantity: inventoryPart.Info.StockQuantity,
			Category:      model.Category(inventoryPart.Info.Category),
			Dimensions:    DimensionsToModel(inventoryPart.Info.Dimensions),
			Manufacturer:  ManufacturerToModel(inventoryPart.Info.Manufacturer),
			Tags:          inventoryPart.Info.Tags,
			Metadata:      MetadataToModel(inventoryPart.Info.Metadata),
		},
		CreatedAt: inventoryPart.CreatedAt.AsTime(),
		UpdatedAt: updatedAt,
	}
}

func DimensionsToModel(dimension *generatedInventoryV1.Dimensions) *model.Dimensions {
	if dimension == nil {
		return nil
	}

	return &model.Dimensions{
		Length: dimension.Length,
		Width:  dimension.Width,
		Height: dimension.Height,
		Weight: dimension.Weight,
	}
}

func ManufacturerToModel(manufacturer *generatedInventoryV1.Manufacturer) *model.Manufacturer {
	if manufacturer == nil {
		return nil
	}

	return &model.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func MetadataToModel(metadata map[string]*generatedInventoryV1.Value) map[string]any {
	res := make(map[string]any, len(metadata))

	for key, value := range metadata {
		if value == nil || value.Kind == nil {
			continue
		}

		switch v := value.Kind.(type) {
		case *generatedInventoryV1.Value_StringValue:
			res[key] = v.StringValue
		case *generatedInventoryV1.Value_BoolValue:
			res[key] = v.BoolValue
		case *generatedInventoryV1.Value_DoubleValue:
			res[key] = v.DoubleValue
		case *generatedInventoryV1.Value_Int64Value:
			res[key] = v.Int64Value
		default:
			log.Printf("unknown metadata kind for key %q: %T", key, value)
		}
	}

	return res
}

func PartsFilterToProto(filter model.PartsFilter) *generatedInventoryV1.PartsFilter {
	categories := make([]generatedInventoryV1.Category, 0, len(filter.Categories))
	for _, category := range filter.Categories {
		categories = append(categories, generatedInventoryV1.Category(category))
	}

	return &generatedInventoryV1.PartsFilter{
		Uuids:                 filter.Uuids,
		Names:                 filter.Names,
		Categories:            categories,
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}
