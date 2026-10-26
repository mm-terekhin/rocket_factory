package v1

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mm-terekhin/rocket_factory/inventory/internal/converter"
	inventoryV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/proto/inventory/v1"
)

func (s *APISuite) TestSuccessCreatePart() {
	var (
		uuid                = gofakeit.UUID()
		name                = gofakeit.Name()
		description         = gofakeit.ProductDescription()
		price               = 345.6
		stockQuantity       = gofakeit.Int64()
		category            = inventoryV1.Category_CATEGORY_ENGINE
		length              = gofakeit.Float64()
		width               = gofakeit.Float64()
		height              = gofakeit.Float64()
		weight              = gofakeit.Float64()
		manufacturerName    = gofakeit.Name()
		manufacturerCountry = gofakeit.Country()
		manufacturerWebsite = gofakeit.URL()

		info = &inventoryV1.PartInfo{
			Name:          name,
			Description:   description,
			Price:         price,
			StockQuantity: stockQuantity,
			Category:      category,
			Dimensions: &inventoryV1.Dimensions{
				Length: length,
				Width:  width,
				Height: height,
				Weight: weight,
			},
			Manufacturer: &inventoryV1.Manufacturer{
				Name:    manufacturerName,
				Country: manufacturerCountry,
				Website: manufacturerWebsite,
			},
		}
	)

	req := &inventoryV1.CreatePartRequest{
		Info: info,
	}

	expectedModelInfo := converter.PartInfoToModel(info)

	s.inventoryService.On("CreatePart", testCtx(), expectedModelInfo).Return(uuid, nil)
	res, err := s.api.CreatePart(testCtx(), req)

	s.Require().NoError(err)
	s.Require().NotEmpty(res)
	s.Require().Equal(res, &inventoryV1.CreatePartResponse{Uuid: uuid})
}

func (s *APISuite) TestFailureCreatePart() {
	var (
		name                = gofakeit.Name()
		description         = gofakeit.ProductDescription()
		price               = 345.6
		stockQuantity       = gofakeit.Int64()
		category            = inventoryV1.Category_CATEGORY_ENGINE
		length              = gofakeit.Float64()
		width               = gofakeit.Float64()
		height              = gofakeit.Float64()
		weight              = gofakeit.Float64()
		manufacturerName    = gofakeit.Name()
		manufacturerCountry = gofakeit.Country()
		manufacturerWebsite = gofakeit.URL()

		info = &inventoryV1.PartInfo{
			Name:          name,
			Description:   description,
			Price:         price,
			StockQuantity: stockQuantity,
			Category:      category,
			Dimensions: &inventoryV1.Dimensions{
				Length: length,
				Width:  width,
				Height: height,
				Weight: weight,
			},
			Manufacturer: &inventoryV1.Manufacturer{
				Name:    manufacturerName,
				Country: manufacturerCountry,
				Website: manufacturerWebsite,
			},
		}

		serviceErr = gofakeit.Error()
	)

	req := &inventoryV1.CreatePartRequest{
		Info: info,
	}

	expectedModelInfo := converter.PartInfoToModel(info)

	s.inventoryService.On("CreatePart", testCtx(), expectedModelInfo).Return("", serviceErr)
	res, err := s.api.CreatePart(testCtx(), req)

	s.Require().Error(err)
	s.Require().Nil(res)
	s.Require().ErrorIs(err, serviceErr)
}
