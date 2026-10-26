package v1

import (
	"context"

	"github.com/mm-terekhin/rocket_factory/order/internal/converter"
	orderV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/openapi/order/v1"
)

func (a *orderApi) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	resp, err := a.OrderService.Create(ctx, converter.CreateOrderDataToModel(req))
	if err != nil {
		return mapCreateOrderErr(err)
	}

	return &orderV1.CreatedOrderResponse{
		OrderUUID:  orderV1.NewOptString(resp.OrderUUID),
		TotalPrice: orderV1.NewOptFloat64(resp.TotalPrice),
	}, nil

}
