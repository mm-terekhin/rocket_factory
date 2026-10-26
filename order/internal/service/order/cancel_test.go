package order

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mm-terekhin/rocket_factory/order/internal/model"
	"github.com/stretchr/testify/mock"
)

func (s *ServiceSuite) TestSuccessCancel() {
	orderUUID := gofakeit.UUID()
	order := &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   gofakeit.UUID(),
		PartUuids:  []string{gofakeit.UUID()},
		TotalPrice: 345.6,
		Status:     string(model.StatusPendingPayment),
	}

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(order)
	s.orderRepository.On("UpdateOrder", s.ctx, mock.MatchedBy(func(o model.Order) bool {
		return o.OrderUUID == orderUUID && o.Status == string(model.StatusCancelled)
	})).Return(nil)

	res, err := s.service.Cancel(s.ctx, orderUUID)

	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(string(model.StatusCancelled), res.Status)
}

func (s *ServiceSuite) TestFailureCancel_NotFound() {
	orderUUID := gofakeit.UUID()

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(nil)

	res, err := s.service.Cancel(s.ctx, orderUUID)

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrorOrderNotFound)
	s.Require().Nil(res)
}

func (s *ServiceSuite) TestFailureCancel_Conflict() {
	orderUUID := gofakeit.UUID()
	order := &model.Order{
		OrderUUID: orderUUID,
		Status:    string(model.StatusPaid),
	}

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(order)

	res, err := s.service.Cancel(s.ctx, orderUUID)

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrorOrderConflict)
	s.Require().Nil(res)
}

func (s *ServiceSuite) TestFailureCancel_UpdateError() {
	orderUUID := gofakeit.UUID()
	order := &model.Order{
		OrderUUID: orderUUID,
		Status:    string(model.StatusPendingPayment),
	}
	serviceErr := gofakeit.Error()

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(order)
	s.orderRepository.On("UpdateOrder", s.ctx, mock.Anything).Return(serviceErr)

	res, err := s.service.Cancel(s.ctx, orderUUID)

	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
	s.Require().Nil(res)
}
