package payment

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mm-terekhin/rocket_factory/payment/internal/model"
)

func (s *ServiceSuite) TestSuccessPay() {
	input := model.PayOrderInput{
		OrderUuid:     gofakeit.UUID(),
		UserUuid:      gofakeit.UUID(),
		PaymentMethod: model.Card,
	}

	res, err := s.service.Pay(testCtx(), input)

	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().NotEmpty(*res)
}

func (s *ServiceSuite) TestFailurePay_EmptyOrderUUID() {
	_, err := s.service.Pay(testCtx(), model.PayOrderInput{
		UserUuid:      gofakeit.UUID(),
		PaymentMethod: model.Card,
	})
	s.requireInvalidArgument(err, "order_uuid is inputuired")
}

func (s *ServiceSuite) TestFailurePay_EmptyUserUUID() {
	_, err := s.service.Pay(testCtx(), model.PayOrderInput{
		OrderUuid:     gofakeit.UUID(),
		PaymentMethod: model.Card,
	})
	s.requireInvalidArgument(err, "user_uuid is inputuired")
}

func (s *ServiceSuite) TestFailurePay_UnspecifiedPaymentMethod() {
	_, err := s.service.Pay(testCtx(), model.PayOrderInput{
		OrderUuid:     gofakeit.UUID(),
		UserUuid:      gofakeit.UUID(),
		PaymentMethod: model.Unspecified,
	})
	s.requireInvalidArgument(err, "payment_method is inputuired")
}
