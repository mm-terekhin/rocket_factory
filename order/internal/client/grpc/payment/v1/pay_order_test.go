package v1

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mm-terekhin/rocket_factory/order/internal/model"
	paymentV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/proto/payment/v1"
	"github.com/stretchr/testify/mock"
)

func (s *ClientSuite) TestSuccessPayOrder() {
	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	txUUID := gofakeit.UUID()

	s.grpc.On("PayOrder", testCtx(), mock.MatchedBy(func(req *paymentV1.PayOrderRequest) bool {
		return req.GetOrderUuid() == orderUUID &&
			req.GetUserUuid() == userUUID &&
			req.GetPaymentMethod() == paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	})).Return(&paymentV1.PayOrderResponse{TransactionUuid: txUUID}, nil)

	res, err := s.client.PayOrder(testCtx(), model.PayOrderData{
		OrderUUID:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: model.PaymentMethodCard,
	})

	s.Require().NoError(err)
	s.Require().Equal(txUUID, res)
}

func (s *ClientSuite) TestFailurePayOrder() {
	serviceErr := gofakeit.Error()

	s.grpc.On("PayOrder", testCtx(), mock.Anything).Return(nil, serviceErr)

	_, err := s.client.PayOrder(testCtx(), model.PayOrderData{
		OrderUUID: gofakeit.UUID(),
		UserUuid:  gofakeit.UUID(),
	})

	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
}
