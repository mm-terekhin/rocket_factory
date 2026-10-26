package converter

import (
	"github.com/mm-terekhin/rocket_factory/inventory/internal/model"
	repoModel "github.com/mm-terekhin/rocket_factory/inventory/internal/repository/model"
)

func PartToModel(repoPart repoModel.Part) model.Part {
	return model.Part{
		UUID:      repoPart.UUID,
		Info:      PartInfoToModel(repoPart.Info),
		CreatedAt: repoPart.CreatedAt,
		UpdatedAt: repoPart.UpdatedAt,
	}
}

func PartInfoToModel(repoInfo repoModel.PartInfo) model.PartInfo {
	return model.PartInfo{
		Name:          repoInfo.Name,
		Description:   repoInfo.Description,
		Price:         repoInfo.Price,
		StockQuantity: repoInfo.StockQuantity,
		Category:      CategoryToModel(repoInfo.Category),
		Dimensions:    DimensionsToModel(repoInfo.Dimensions),
		Manufacturer:  ManufacturerToModel(repoInfo.Manufacturer),
		Tags:          repoInfo.Tags,
		Metadata:      MetadataToModel(repoInfo.Metadata),
	}
}

func DimensionsToModel(dimensions repoModel.Dimensions) model.Dimensions {
	return model.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func ManufacturerToModel(manufacturer repoModel.Manufacturer) model.Manufacturer {
	return model.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func MetadataToModel(metadate map[string]*repoModel.Value) map[string]*model.Value {
	convertedMetadate := make(map[string]*model.Value, len(metadate))

	for key, value := range metadate {
		convertedMetadate[key] = ValueToModel(value)
	}

	return convertedMetadate
}

func ValueToModel(value *repoModel.Value) *model.Value {
	return &model.Value{
		StringValue: value.StringValue,
		Int64Value:  value.Int64Value,
		DoubleValue: value.DoubleValue,
		BoolValue:   value.BoolValue,
	}
}

func CategoryToModel(c repoModel.Category) model.Category {
	return model.Category(c)
}

func CategoryToRepo(c model.Category) repoModel.Category {
	return repoModel.Category(c)
}
