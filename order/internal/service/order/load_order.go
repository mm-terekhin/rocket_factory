package order

import (
	"context"

	"github.com/mm-terekhin/rocket_factory/order/internal/model"
)

func (s *service) loadOrder(ctx context.Context, uuid string) (*model.Order, error) {
	order, err := s.OrderRepository.GetOrder(ctx, uuid)
	if err != nil {
		return nil, model.ErrOrderNotFound
	}

	return order, nil
}
