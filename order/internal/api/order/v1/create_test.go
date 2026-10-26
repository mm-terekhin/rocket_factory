package v1

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mm-terekhin/rocket_factory/order/internal/model"
	orderV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/openapi/order/v1"
)

func (s *APISuite) TestSuccessCreateOrder() {
	userUUID := gofakeit.UUID()
	partUUID := gofakeit.UUID()
	order := &model.Order{
		OrderUUID:  gofakeit.UUID(),
		UserUUID:   userUUID,
		PartUuids:  []string{partUUID},
		TotalPrice: 200,
		Status:     string(model.StatusPendingPayment),
	}

	s.orderService.On("Create", testCtx(), model.CreateOrderData{
		UserUUID:  userUUID,
		PartUuids: []string{partUUID},
	}).Return(order, nil)

	res, err := s.api.CreateOrder(testCtx(), &orderV1.CreateOrderRequest{
		UserUUID:  userUUID,
		PartUuids: []string{partUUID},
	})

	s.Require().NoError(err)
	created, ok := res.(*orderV1.CreatedOrderResponse)
	s.Require().True(ok)
	s.Require().Equal(order.OrderUUID, created.GetOrderUUID().Or(""))
	s.Require().Equal(order.TotalPrice, created.GetTotalPrice().Or(0))
}

func (s *APISuite) TestFailureCreateOrder_BadRequest() {
	s.orderService.On("Create", testCtx(), model.CreateOrderData{}).
		Return(nil, &model.InvalidArgumentError{Message: "user_uuid is required"})

	res, err := s.api.CreateOrder(testCtx(), &orderV1.CreateOrderRequest{})

	s.Require().NoError(err)
	br, ok := res.(*orderV1.BadRequestError)
	s.Require().True(ok)
	s.Require().Equal(400, br.Code)
	s.Require().Equal("user_uuid is required", br.Message)
}
