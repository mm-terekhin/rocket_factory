package order

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mm-terekhin/rocket_factory/order/internal/model"
)

func (s *ServiceSuite) TestSuccessGet() {
	orderUUID := gofakeit.UUID()
	expected := &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   gofakeit.UUID(),
		PartUuids:  []string{gofakeit.UUID()},
		TotalPrice: 100,
		Status:     string(model.StatusPendingPayment),
	}

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(expected)

	res, err := s.service.Get(s.ctx, orderUUID)

	s.Require().NoError(err)
	s.Require().Equal(expected, res)
}

func (s *ServiceSuite) TestFailureGet_NotFound() {
	orderUUID := gofakeit.UUID()

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(nil)

	res, err := s.service.Get(s.ctx, orderUUID)

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrorOrderNotFound)
	s.Require().Nil(res)
}
