package inventory

import (
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mm-terekhin/rocket_factory/inventory/internal/model"
)

func (s *ServiceSuite) TestSuccessListParts() {
	filter := model.PartsFilter{}

	parts := fakeParts(3)

	s.inventoryRepository.On("ListParts", s.ctx, filter).Return(parts, nil)

	res, err := s.service.ListParts(s.ctx, filter)
	s.Require().NoError(err)
	s.Require().NotEmpty(res)
	s.Require().Equal(res, parts)
}

func (s *ServiceSuite) TestFailureListParts() {
	filter := model.PartsFilter{}
	parts := make([]model.Part, 0)
	serviceErr := gofakeit.Error()

	s.inventoryRepository.On("ListParts", s.ctx, filter).Return(parts, serviceErr)

	res, err := s.service.ListParts(s.ctx, filter)
	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
	s.Require().Empty(res)
}

func fakePart() model.Part {
	now := time.Now()
	return model.Part{
		UUID: gofakeit.UUID(),
		Info: model.PartInfo{
			Name:          gofakeit.Name(),
			Description:   gofakeit.ProductDescription(),
			Price:         gofakeit.Float64Range(1, 10000),
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
		},
		CreatedAt: now,
	}
}

func fakeParts(n int) []model.Part {
	parts := make([]model.Part, 0, n)
	for i := 0; i < n; i++ {
		parts = append(parts, fakePart())
	}
	return parts
}
