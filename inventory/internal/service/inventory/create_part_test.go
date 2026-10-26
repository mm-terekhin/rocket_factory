package inventory

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mm-terekhin/rocket_factory/inventory/internal/model"
)

func (s *ServiceSuite) TestSuccessCreatePart() {
	info := model.PartInfo{
		Name:          gofakeit.Name(),
		Description:   gofakeit.ProductDescription(),
		Price:         345.6,
		StockQuantity: gofakeit.Int64(),
		Category:      model.CategoryEngine,
		Dimensions: model.Dimensions{
			Length: gofakeit.Float64(),
			Width:  gofakeit.Float64(),
			Height: gofakeit.Float64(),
			Weight: gofakeit.Float64(),
		},
		Manufacturer: model.Manufacturer{
			Name:    gofakeit.Name(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		},
	}

	expectedUUID := gofakeit.UUID()

	s.inventoryRepository.On("CreatePart", testCtx(), info).Return(expectedUUID, nil)

	res, err := s.service.CreatePart(testCtx(), info)

	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedUUID, res)
}

func (s *ServiceSuite) TestFailureCreatePart() {
	info := model.PartInfo{
		Name:          gofakeit.Name(),
		Description:   gofakeit.ProductDescription(),
		Price:         345.6,
		StockQuantity: gofakeit.Int64(),
		Category:      model.CategoryEngine,
		Dimensions: model.Dimensions{
			Length: gofakeit.Float64(),
			Width:  gofakeit.Float64(),
			Height: gofakeit.Float64(),
			Weight: gofakeit.Float64(),
		},
		Manufacturer: model.Manufacturer{
			Name:    gofakeit.Name(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		},
	}

	expectedUUID := gofakeit.UUID()

	serviceErr := gofakeit.Error()

	s.inventoryRepository.On("CreatePart", testCtx(), info).Return(expectedUUID, serviceErr)

	res, err := s.service.CreatePart(testCtx(), info)

	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
	s.Require().NotNil(res)
}
