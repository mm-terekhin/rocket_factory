package inventory

import (
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mm-terekhin/rocket_factory/inventory/internal/model"
)

func (s *ServiceSuite) TestSuccessGetPart() {
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
		}
	)

	s.inventoryRepository.On("GetPart", testCtx(), uuid).Return(part, nil)

	res, err := s.service.GetPart(testCtx(), uuid)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(part, res)
}

func (s *ServiceSuite) TestFailureGetPart() {
	var (
		uuid        = gofakeit.UUID()
		notFoundErr = model.ErrPartNotFound
		part        = model.Part{}
	)

	s.inventoryRepository.On("GetPart", testCtx(), uuid).Return(part, notFoundErr)

	res, err := s.service.GetPart(testCtx(), uuid)
	s.Require().Error(err)
	s.Require().ErrorIs(err, notFoundErr)
	s.Require().Equal(res, part)
}
