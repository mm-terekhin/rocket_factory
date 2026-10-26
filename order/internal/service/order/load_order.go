package order

import (
	"context"

	"github.com/mm-terekhin/rocket_factory/order/internal/model"
)

func (s *service) loadOrder(ctx context.Context, uuid string) (*model.Order, error) {
	order := s.OrderRepository.GetOrder(ctx, uuid)

	if order == nil {
		return nil, model.ErrOrderNotFound
	}

	return order, nil
}
