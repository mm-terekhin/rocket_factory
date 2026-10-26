package v1

import (
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mm-terekhin/rocket_factory/inventory/internal/converter"
	"github.com/mm-terekhin/rocket_factory/inventory/internal/model"
	inventoryV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/proto/inventory/v1"
)

func (s *APISuite) TestSuccessGetPart() {
	var (
		uuid                = gofakeit.UUID()
		name                = gofakeit.Name()
		description         = gofakeit.ProductDescription()
		price               = 345.6
		stockQuantity       = gofakeit.Int64()
		category            = model.CategoryEngine
		length              = gofakeit.Float64()
		width               = gofakeit.Float64()
		height              = gofakeit.Float64()
		weight              = gofakeit.Float64()
		manufacturerName    = gofakeit.Name()
		manufacturerCountry = gofakeit.Country()
		manufacturerWebsite = gofakeit.URL()
		createdAt           = time.Now()

		part = model.Part{
			UUID: uuid,
			Info: model.PartInfo{
				Name:          name,
				Description:   description,
				Price:         price,
				StockQuantity: stockQuantity,
				Category:      category,
				Dimensions: model.Dimensions{
					Length: length,
					Width:  width,
					Height: height,
					Weight: weight,
				},
				Manufacturer: model.Manufacturer{
					Name:    manufacturerName,
					Country: manufacturerCountry,
					Website: manufacturerWebsite,
				},
			},
			CreatedAt: createdAt,
			UpdatedAt: &createdAt,
		}
	)

	req := &inventoryV1.GetPartRequest{Uuid: uuid}

	s.inventoryService.On("GetPart", s.ctx, uuid).Return(part, nil)

	res, err := s.api.GetPart(s.ctx, req)

	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(res, &inventoryV1.GetPartResponse{Part: converter.PartToProto(part)})
}

func (s *APISuite) TestFailureGetPart() {
	var (
		uuid = gofakeit.UUID()

		part       = model.Part{}
		serviceErr = model.ErrorPartNotFound
	)

	req := &inventoryV1.GetPartRequest{Uuid: uuid}

	serviceErr = model.ErrorPartNotFound
	s.inventoryService.On("GetPart", s.ctx, uuid).Return(part, serviceErr)

	res, err := s.api.GetPart(s.ctx, req)

	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
	s.Require().Nil(res)
}
