package v1

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mm-terekhin/rocket_factory/order/internal/converter"
	"github.com/mm-terekhin/rocket_factory/order/internal/model"
	orderV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/openapi/order/v1"
)

func (s *APISuite) TestSuccessGetOrderByUUID() {
	orderUUID := gofakeit.UUID()
	order := &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   gofakeit.UUID(),
		PartUuids:  []string{gofakeit.UUID()},
		TotalPrice: 100,
		Status:     string(model.StatusPendingPayment),
	}

	s.orderService.On("Get", testCtx(), orderUUID).Return(order, nil)

	res, err := s.api.GetOrderByUUID(testCtx(), orderV1.GetOrderByUUIDParams{OrderUUID: orderUUID})

	s.Require().NoError(err)
	s.Require().Equal(converter.OrderToResponse(order), res)
}

func (s *APISuite) TestFailureGetOrderByUUID_NotFound() {
	orderUUID := gofakeit.UUID()

	s.orderService.On("Get", testCtx(), orderUUID).Return(nil, model.ErrOrderNotFound)

	res, err := s.api.GetOrderByUUID(testCtx(), orderV1.GetOrderByUUIDParams{OrderUUID: orderUUID})

	s.Require().NoError(err)
	nf, ok := res.(*orderV1.NotFoundError)
	s.Require().True(ok)
	s.Require().Equal(404, nf.Code)
}
