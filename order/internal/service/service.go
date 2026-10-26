package service

import (
	"context"

	"github.com/mm-terekhin/rocket_factory/order/internal/model"
)

type OrderService interface {
	Create(context.Context, model.CreateOrderData) (*model.Order, error)
	Get(context.Context, string) (*model.Order, error)
	Pay(context.Context, model.PayOrderInput) (*string, error)
	Cancel(context.Context, string) (*model.Order, error)
}
