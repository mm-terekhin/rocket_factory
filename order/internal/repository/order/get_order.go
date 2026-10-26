package order

import (
	"context"

	"github.com/mm-terekhin/rocket_factory/order/internal/model"
	"github.com/mm-terekhin/rocket_factory/order/internal/repository/converter"
)

func (r *repository) GetOrder(ctx context.Context, uuid string) *model.Order {
	r.mu.RLock()
	defer r.mu.RUnlock()

	val, ok := r.orders[uuid]

	if !ok {
		return nil
	}

	order := converter.OrderToModel(val)

	return &order
}
