package v1

import (
	"context"

	"github.com/mm-terekhin/rocket_factory/order/internal/converter"
	orderV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/openapi/order/v1"
)

func (a *orderApi) GetOrderByUUID(ctx context.Context, params orderV1.GetOrderByUUIDParams) (orderV1.GetOrderByUUIDRes, error) {
	resp, err := a.OrderService.Get(ctx, params.OrderUUID)
	if err != nil {
		return mapGetOrderErr(err, params.OrderUUID)
	}

	return converter.OrderToResponse(resp), nil
}
