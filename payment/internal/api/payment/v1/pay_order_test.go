package v1

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mm-terekhin/rocket_factory/payment/internal/model"
	paymentV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/proto/payment/v1"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *APISuite) TestSuccessPayOrder() {
	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	txUUID := gofakeit.UUID()

	s.paymentService.On("Pay", testCtx(), mock.MatchedBy(func(in model.PayOrderInput) bool {
		return in.OrderUuid == orderUUID &&
			in.UserUuid == userUUID &&
			in.PaymentMethod == model.Card
	})).Return(&txUUID, nil)

	resp, err := s.api.PayOrder(testCtx(), &paymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
	})

	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Equal(txUUID, resp.GetTransactionUuid())
}

func (s *APISuite) TestFailurePayOrder_InvalidArgument() {
	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()

	s.paymentService.On("Pay", testCtx(), mock.MatchedBy(func(in model.PayOrderInput) bool {
		return in.OrderUuid == orderUUID && in.UserUuid == userUUID
	})).Return(nil, &model.InvalidArgumentError{Message: "payment_method is inputuired"})

	_, err := s.api.PayOrder(testCtx(), &paymentV1.PayOrderRequest{
		OrderUuid: orderUUID,
		UserUuid:  userUUID,
	})

	s.Require().Error(err)
	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.InvalidArgument, st.Code())
	s.Require().Equal("payment_method is inputuired", st.Message())
}
