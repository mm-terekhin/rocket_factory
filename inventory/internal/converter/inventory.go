package converter

import (
	"github.com/mm-terekhin/rocket_factory/inventory/internal/model"
	repoModel "github.com/mm-terekhin/rocket_factory/inventory/internal/repository/model"
	inventoryV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func FilterToModel(filter *inventoryV1.PartsFilter) model.PartsFilter {
	categories := make([]model.Category, 0, len(filter.GetCategories()))

	for _, category := range filter.GetCategories() {
		categories = append(categories, partCategoryToModel(category))
	}

	return model.PartsFilter{
		UUIDs:                 filter.GetUuids(),
		Names:                 filter.GetNames(),
		Categories:            categories,
		ManufacturerCountries: filter.GetManufacturerCountries(),
		Tags:                  filter.GetTags(),
	}
}

func partCategoryToModel(category inventoryV1.Category) model.Category {
	return model.Category(category)
}

func PartToProto(part model.Part) *inventoryV1.Part {
	return &inventoryV1.Part{
		Uuid:      part.UUID,
		Info:      partInfoToProto(part.Info),
		CreatedAt: timestamppb.New(part.CreatedAt),
		UpdatedAt: timestamppb.New(*part.UpdatedAt),
	}
}

func partInfoToProto(info model.PartInfo) *inventoryV1.PartInfo {
	return &inventoryV1.PartInfo{
		Name:          info.Name,
		Description:   info.Description,
		Price:         info.Price,
		StockQuantity: info.StockQuantity,
		Category:      partCategoryToProto(info.Category),
		Dimensions:    partDimensionsToProto(info.Dimensions),
		Manufacturer:  partManufacturerToProto(info.Manufacturer),
		Tags:          info.Tags,
		Metadata:      partMetadataToProto(info.Metadata),
	}
}

func partCategoryToProto(category model.Category) inventoryV1.Category {
	switch category {
	case model.CategoryEngine:
		return inventoryV1.Category_CATEGORY_ENGINE
	case model.CategoryFuel:
		return inventoryV1.Category_CATEGORY_FUEL
	case model.CategoryPorthole:
		return inventoryV1.Category_CATEGORY_PORTHOLE
	case model.CategoryWing:
		return inventoryV1.Category_CATEGORY_WING
	default:
		return inventoryV1.Category_CATEGORY_UNSPECIFIED
	}
}

func partDimensionsToProto(dimensions model.Dimensions) *inventoryV1.Dimensions {
	return &inventoryV1.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func partManufacturerToProto(manufacturer model.Manufacturer) *inventoryV1.Manufacturer {
	return &inventoryV1.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func partMetadataToProto(metadata map[string]*model.Value) map[string]*inventoryV1.Value {
	proto := make(map[string]*inventoryV1.Value, len(metadata))

	for key, value := range metadata {
		proto[key] = partValueToProto(value)
	}

	return proto
}

func partValueToProto(value *model.Value) *inventoryV1.Value {
	if value == nil {
		return nil
	}

	switch {
	case value.StringValue != nil:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_StringValue{
				StringValue: *value.StringValue,
			},
		}
	case value.Int64Value != nil:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_Int64Value{
				Int64Value: *value.Int64Value,
			},
		}
	case value.DoubleValue != nil:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_DoubleValue{
				DoubleValue: *value.DoubleValue,
			},
		}
	case value.BoolValue != nil:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_BoolValue{
				BoolValue: *value.BoolValue,
			},
		}
	default:
		return nil
	}
}

func PartInfoToModel(info *inventoryV1.PartInfo) model.PartInfo {
	modelInfo := model.PartInfo{
		Name:          info.GetName(),
		Description:   info.GetDescription(),
		Price:         info.GetPrice(),
		StockQuantity: info.GetStockQuantity(),
		Category:      InfoCategoryToModel(info.GetCategory()),
		Dimensions:    InfoDimensionsToModel(info.GetDimensions()),
		Manufacturer:  InfoManufacturerToModel(info.GetManufacturer()),
	}

	return modelInfo
}

func InfoCategoryToModel(category inventoryV1.Category) model.Category {
	return partCategoryToModel(category)
}

func InfoDimensionsToModel(dimensions *inventoryV1.Dimensions) model.Dimensions {
	return model.Dimensions{
		Length: dimensions.GetLength(),
		Width:  dimensions.GetWidth(),
		Height: dimensions.GetHeight(),
		Weight: dimensions.GetWeight(),
	}
}

func InfoManufacturerToModel(manufacturer *inventoryV1.Manufacturer) model.Manufacturer {
	return model.Manufacturer{
		Name:    manufacturer.GetName(),
		Country: manufacturer.GetCountry(),
		Website: manufacturer.GetWebsite(),
	}
}

func PartInfoToRepo(repoInfo model.PartInfo) repoModel.PartInfo {
	return repoModel.PartInfo{
		Name:          repoInfo.Name,
		Description:   repoInfo.Description,
		Price:         repoInfo.Price,
		StockQuantity: repoInfo.StockQuantity,
		Category:      CategoryToRepo(repoInfo.Category),
		Dimensions:    DimensionsToRepo(repoInfo.Dimensions),
		Manufacturer:  ManufacturerToRepo(repoInfo.Manufacturer),
		Tags:          repoInfo.Tags,
		Metadata:      MetadataToRepo(repoInfo.Metadata),
	}
}

func DimensionsToRepo(dimensions model.Dimensions) repoModel.Dimensions {
	return repoModel.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func ManufacturerToRepo(manufacturer model.Manufacturer) repoModel.Manufacturer {
	return repoModel.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func MetadataToRepo(metadate map[string]*model.Value) map[string]*repoModel.Value {
	convertedMetadate := make(map[string]*repoModel.Value, len(metadate))

	for key, value := range metadate {
		convertedMetadate[key] = ValueToRepo(value)
	}

	return convertedMetadate
}

func ValueToRepo(value *model.Value) *repoModel.Value {
	return &repoModel.Value{
		StringValue: value.StringValue,
		Int64Value:  value.Int64Value,
		DoubleValue: value.DoubleValue,
		BoolValue:   value.BoolValue,
	}
}

func CategoryToRepo(c model.Category) repoModel.Category {
	return repoModel.Category(c)
}
