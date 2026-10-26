package order

import (
	"context"

	"github.com/mm-terekhin/rocket_factory/order/internal/model"
	"github.com/mm-terekhin/rocket_factory/order/internal/repository/converter"
)

func (r *repository) UpdateOrder(ctx context.Context, order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.orders[order.OrderUUID]

	if !ok {
		return model.ErrorOrderNotFound
	}

	r.orders[order.OrderUUID] = converter.OrderToRepo(order)

	return nil
}
