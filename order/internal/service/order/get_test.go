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

	s.orderRepository.On("GetOrder", testCtx(), orderUUID).Return(expected)

	res, err := s.service.Get(testCtx(), orderUUID)

	s.Require().NoError(err)
	s.Require().Equal(expected, res)
}

func (s *ServiceSuite) TestFailureGet_NotFound() {
	orderUUID := gofakeit.UUID()

	s.orderRepository.On("GetOrder", testCtx(), orderUUID).Return(nil)

	res, err := s.service.Get(testCtx(), orderUUID)

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderNotFound)
	s.Require().Nil(res)
}
