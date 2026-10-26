package v1

import (
	"github.com/mm-terekhin/rocket_factory/order/internal/service"
	orderV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/openapi/order/v1"
)

var _ orderV1.Handler = (*orderApi)(nil)

type orderApi struct {
	OrderService service.OrderService
}

func NewApi(
	service service.OrderService,
) *orderApi {
	return &orderApi{
		OrderService: service,
	}
}
