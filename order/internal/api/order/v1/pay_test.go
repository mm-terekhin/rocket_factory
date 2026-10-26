package v1

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mm-terekhin/rocket_factory/order/internal/model"
	orderV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/openapi/order/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *APISuite) TestSuccessPayOrderByUUID() {
	orderUUID := gofakeit.UUID()
	txUUID := gofakeit.UUID()

	s.orderService.On("Pay", s.ctx, model.PayOrderInput{
		OrderUUID:     orderUUID,
		PaymentMethod: model.PaymentMethodCard,
	}).Return(&txUUID, nil)

	res, err := s.api.PayOrderByUUID(s.ctx, &orderV1.PayOrderRequest{
		PaymentMethod: "CARD",
	}, orderV1.PayOrderByUUIDParams{OrderUUID: orderUUID})

	s.Require().NoError(err)
	paid, ok := res.(*orderV1.PayOrderResponse)
	s.Require().True(ok)
	s.Require().Equal(txUUID, paid.TransactionUUID)
}

func (s *APISuite) TestFailurePayOrderByUUID_NotFound() {
	orderUUID := gofakeit.UUID()

	s.orderService.On("Pay", s.ctx, model.PayOrderInput{
		OrderUUID:     orderUUID,
		PaymentMethod: model.PaymentMethodUnspecified,
	}).Return(nil, model.ErrorOrderNotFound)

	res, err := s.api.PayOrderByUUID(s.ctx, &orderV1.PayOrderRequest{}, orderV1.PayOrderByUUIDParams{OrderUUID: orderUUID})

	s.Require().NoError(err)
	nf, ok := res.(*orderV1.NotFoundError)
	s.Require().True(ok)
	s.Require().Equal(404, nf.Code)
}

func (s *APISuite) TestFailurePayOrderByUUID_BadRequest() {
	orderUUID := gofakeit.UUID()
	grpcErr := status.Error(codes.InvalidArgument, "payment method is required")

	s.orderService.On("Pay", s.ctx, model.PayOrderInput{
		OrderUUID:     orderUUID,
		PaymentMethod: model.PaymentMethodUnspecified,
	}).Return(nil, grpcErr)

	res, err := s.api.PayOrderByUUID(s.ctx, &orderV1.PayOrderRequest{}, orderV1.PayOrderByUUIDParams{OrderUUID: orderUUID})

	s.Require().NoError(err)
	br, ok := res.(*orderV1.BadRequestError)
	s.Require().True(ok)
	s.Require().Equal(400, br.Code)
	s.Require().Equal("payment method is required", br.Message)
}
