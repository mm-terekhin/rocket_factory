package v1

import (
	"context"

	"github.com/mm-terekhin/rocket_factory/order/internal/converter"
	"github.com/mm-terekhin/rocket_factory/order/internal/model"
	orderV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/openapi/order/v1"
)

func (a *orderApi) PayOrderByUUID(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderByUUIDParams) (orderV1.PayOrderByUUIDRes, error) {
	method := converter.PaymentMethodFromHTTP(req.GetPaymentMethod())

	resp, err := a.OrderService.Pay(ctx, model.PayOrderInput{
		OrderUUID:     params.OrderUUID,
		PaymentMethod: method,
	})
	if err != nil {
		return mapPayErr(err, params.OrderUUID)
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: *resp,
	}, nil
}
