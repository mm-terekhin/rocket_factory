package v1

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mm-terekhin/rocket_factory/order/internal/model"
	orderV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/openapi/order/v1"
)

func (s *APISuite) TestSuccessCancelOrderByUUID() {
	orderUUID := gofakeit.UUID()
	order := &model.Order{OrderUUID: orderUUID, Status: string(model.StatusCancelled)}

	s.orderService.On("Cancel", testCtx(), orderUUID).Return(order, nil)

	res, err := s.api.CancelOrderByUUID(testCtx(), orderV1.CancelOrderByUUIDParams{OrderUUID: orderUUID})

	s.Require().NoError(err)
	s.Require().Nil(res)
}

func (s *APISuite) TestFailureCancelOrderByUUID_NotFound() {
	orderUUID := gofakeit.UUID()

	s.orderService.On("Cancel", testCtx(), orderUUID).Return(nil, model.ErrOrderNotFound)

	res, err := s.api.CancelOrderByUUID(testCtx(), orderV1.CancelOrderByUUIDParams{OrderUUID: orderUUID})

	s.Require().NoError(err)
	nf, ok := res.(*orderV1.NotFoundError)
	s.Require().True(ok)
	s.Require().Equal(404, nf.Code)
}

func (s *APISuite) TestFailureCancelOrderByUUID_Conflict() {
	orderUUID := gofakeit.UUID()

	s.orderService.On("Cancel", testCtx(), orderUUID).Return(nil, model.ErrOrderConflict)

	res, err := s.api.CancelOrderByUUID(testCtx(), orderV1.CancelOrderByUUIDParams{OrderUUID: orderUUID})

	s.Require().NoError(err)
	conf, ok := res.(*orderV1.ConflictError)
	s.Require().True(ok)
	s.Require().Equal(409, conf.Code)
}
