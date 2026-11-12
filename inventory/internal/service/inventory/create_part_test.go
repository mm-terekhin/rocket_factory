package inventory

import (
	"reflect"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/mm-terekhin/rocket_factory/inventory/internal/model"
)

func partInfoMatcher(info model.PartInfo) any {
	return mock.MatchedBy(func(p model.Part) bool {
		return reflect.DeepEqual(p.Info, info) && p.UUID != "" && !p.CreatedAt.IsZero()
	})
}

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

	s.inventoryRepository.On("CreatePart", testCtx(), partInfoMatcher(info)).Return(primitive.NewObjectID(), nil)

	res, err := s.service.CreatePart(testCtx(), info)

	s.Require().NoError(err)
	s.Require().NotEmpty(res)
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

	serviceErr := gofakeit.Error()

	s.inventoryRepository.On("CreatePart", testCtx(), partInfoMatcher(info)).Return(primitive.NilObjectID, serviceErr)

	res, err := s.service.CreatePart(testCtx(), info)

	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
	s.Require().Empty(res)
}
