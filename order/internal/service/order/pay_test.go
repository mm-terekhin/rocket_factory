package order

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mm-terekhin/rocket_factory/order/internal/model"
	"github.com/stretchr/testify/mock"
)

func (s *ServiceSuite) TestSuccessPay() {
	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	txUUID := gofakeit.UUID()
	order := &model.Order{
		OrderUUID: orderUUID,
		UserUUID:  userUUID,
		Status:    string(model.StatusPendingPayment),
	}

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(order)
	s.paymentClient.On("PayOrder", s.ctx, model.PayOrderData{
		OrderUUID:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: model.PaymentMethodCard,
	}).Return(txUUID, nil)
	s.orderRepository.On("UpdateOrder", s.ctx, mock.MatchedBy(func(o model.Order) bool {
		return o.OrderUUID == orderUUID &&
			o.Status == string(model.StatusPaid) &&
			o.TransactionUUID != nil &&
			*o.TransactionUUID == txUUID
	})).Return(nil)

	res, err := s.service.Pay(s.ctx, model.PayOrderInput{
		OrderUUID:     orderUUID,
		PaymentMethod: model.PaymentMethodCard,
	})

	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(txUUID, *res)
}

func (s *ServiceSuite) TestFailurePay_NotFound() {
	orderUUID := gofakeit.UUID()

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(nil)

	res, err := s.service.Pay(s.ctx, model.PayOrderInput{OrderUUID: orderUUID})

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrorOrderNotFound)
	s.Require().Nil(res)
}

func (s *ServiceSuite) TestFailurePay_PaymentError() {
	orderUUID := gofakeit.UUID()
	order := &model.Order{OrderUUID: orderUUID, UserUUID: gofakeit.UUID()}
	serviceErr := gofakeit.Error()

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(order)
	s.paymentClient.On("PayOrder", s.ctx, mock.Anything).Return("", serviceErr)

	res, err := s.service.Pay(s.ctx, model.PayOrderInput{OrderUUID: orderUUID})

	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
	s.Require().Nil(res)
}

func (s *ServiceSuite) TestFailurePay_UpdateError() {
	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	txUUID := gofakeit.UUID()
	order := &model.Order{OrderUUID: orderUUID, UserUUID: userUUID}
	serviceErr := gofakeit.Error()

	s.orderRepository.On("GetOrder", s.ctx, orderUUID).Return(order)
	s.paymentClient.On("PayOrder", s.ctx, mock.Anything).Return(txUUID, nil)
	s.orderRepository.On("UpdateOrder", s.ctx, mock.Anything).Return(serviceErr)

	res, err := s.service.Pay(s.ctx, model.PayOrderInput{OrderUUID: orderUUID})

	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
	s.Require().Nil(res)
}
