package order

import (
	"context"

	"github.com/mm-terekhin/rocket_factory/order/internal/model"
)

func (s *service) Pay(ctx context.Context, input model.PayOrderInput) (*string, error) {
	order, err := s.loadOrder(ctx, input.OrderUUID)
	if err != nil {
		return nil, err
	}

	resp, err := s.PaymentClient.PayOrder(ctx, model.PayOrderData{
		OrderUUID:     input.OrderUUID,
		UserUuid:      order.UserUUID,
		PaymentMethod: input.PaymentMethod,
	})
	if err != nil {
		return nil, err
	}

	order.TransactionUUID = &resp
	order.Status = string(model.StatusPaid)

	if err := s.OrderRepository.UpdateOrder(ctx, *order); err != nil {
		return nil, err
	}

	return order.TransactionUUID, nil

}
