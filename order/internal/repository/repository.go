package repository

import (
	"context"

	"github.com/mm-terekhin/rocket_factory/order/internal/model"
)

type OrderRepository interface {
	GetOrder(context.Context, string) *model.Order
	UpdateOrder(context.Context, model.Order) error
	CreateOrder(context.Context, model.Order) error
}
