package order

import (
	"context"

	"github.com/mm-terekhin/rocket_factory/order/internal/model"
)

func (s *service) Get(ctx context.Context, uuid string) (*model.Order, error) {
	order, err := s.loadOrder(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return order, nil
}
