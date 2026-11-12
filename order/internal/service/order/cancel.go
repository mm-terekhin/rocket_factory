package order

import (
	"context"

	"github.com/mm-terekhin/rocket_factory/order/internal/model"
)

func (s *service) Cancel(ctx context.Context, uuid string) (*model.Order, error) {
	order, err := s.OrderRepository.GetOrder(ctx, uuid)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, model.ErrOrderNotFound
	}

	if order.Status == string(model.StatusPaid) {
		return nil, model.ErrOrderConflict
	}

	if order.Status != string(model.StatusCancelled) {
		order.Status = string(model.StatusCancelled)

		if err := s.OrderRepository.UpdateOrder(ctx, *order); err != nil {
			return nil, err
		}
	}

	return order, nil
}
