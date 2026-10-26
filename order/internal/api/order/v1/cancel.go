package v1

import (
	"context"

	orderV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/openapi/order/v1"
)

func (a *orderApi) CancelOrderByUUID(ctx context.Context, params orderV1.CancelOrderByUUIDParams) (orderV1.CancelOrderByUUIDRes, error) {
	_, err := a.OrderService.Cancel(ctx, params.OrderUUID)
	if err != nil {
		return mapCancelErr(err, params.OrderUUID)
	}

	return nil, nil
}
