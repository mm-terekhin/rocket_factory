package client

import (
	"context"

	"github.com/mm-terekhin/rocket_factory/order/internal/model"
)

type InventoryClient interface {
	ListParts(context.Context, model.PartsFilter) ([]model.Part, error)
}

type PaymentClient interface {
	PayOrder(context.Context, model.PayOrderData) (string, error)
}
