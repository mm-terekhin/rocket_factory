package v1

import (
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mm-terekhin/rocket_factory/inventory/internal/converter"
	"github.com/mm-terekhin/rocket_factory/inventory/internal/model"
	inventoryV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/proto/inventory/v1"
)

func (s *APISuite) TestSuccessListParts() {
	filter := model.PartsFilter{
		UUIDs:                 make([]string, 0),
		Names:                 make([]string, 0),
		Categories:            make([]model.Category, 0),
		ManufacturerCountries: make([]string, 0),
		Tags:                  make([]string, 0),
	}

	parts := fakeParts(3)

	req := &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Uuids:                 make([]string, 0),
			Names:                 make([]string, 0),
			Categories:            make([]inventoryV1.Category, 0),
			ManufacturerCountries: make([]string, 0),
			Tags:                  make([]string, 0),
		},
	}

	s.inventoryService.On("ListParts", s.ctx, filter).Return(parts, nil)
	expected := make([]*inventoryV1.Part, 0, len(parts))
	for _, p := range parts {
		expected = append(expected, converter.PartToProto(p))
	}
	res, err := s.api.ListParts(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotEmpty(res)

	s.Require().Equal(res, &inventoryV1.ListPartsResponse{Parts: expected})
}

func (s *APISuite) TestFailureListParts() {
	filter := model.PartsFilter{
		UUIDs:                 make([]string, 0),
		Names:                 make([]string, 0),
		Categories:            make([]model.Category, 0),
		ManufacturerCountries: make([]string, 0),
		Tags:                  make([]string, 0),
	}

	serviceErr := gofakeit.Error()

	req := &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Uuids:                 make([]string, 0),
			Names:                 make([]string, 0),
			Categories:            make([]inventoryV1.Category, 0),
			ManufacturerCountries: make([]string, 0),
			Tags:                  make([]string, 0),
		},
	}

	s.inventoryService.On("ListParts", s.ctx, filter).Return(make([]model.Part, 0), serviceErr)

	res, err := s.api.ListParts(s.ctx, req)
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
		UpdatedAt: &now,
	}
}

func fakeParts(n int) []model.Part {
	parts := make([]model.Part, 0, n)
	for i := 0; i < n; i++ {
		parts = append(parts, fakePart())
	}
	return parts
}
