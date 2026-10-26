package order

import (
	"context"
	"fmt"
	"uuid"

	"github.com/mm-terekhin/rocket_factory/order/internal/model"
)

func (s *service) Create(ctx context.Context, data model.CreateOrderData) (*model.Order, error) {
	if data.UserUUID == "" {
		return nil, &model.InvalidArgumentError{Message: "user_uuid is required"}
	}

	if len(data.PartUuids) == 0 {
		return nil, &model.InvalidArgumentError{Message: "order has no parts"}
	}

	unique := make(map[string]struct{}, len(data.PartUuids))
	uniqueIDs := make([]string, 0, len(data.PartUuids))
	for _, id := range data.PartUuids {
		if id == "" {
			return nil, &model.InvalidArgumentError{Message: "part_uuid must not be empty"}
		}
		if _, ok := unique[id]; ok {
			continue
		}

		unique[id] = struct{}{}
		uniqueIDs = append(uniqueIDs, id)
	}

	filter := model.PartsFilter{
		Uuids: uniqueIDs,
	}

	listParts, err := s.InventoryClient.ListParts(ctx, filter)
	if err != nil {
		return nil, err
	}

	byUUID := make(map[string]model.Part, len(listParts))
	for _, p := range listParts {
		byUUID[p.UUID] = p
	}

	var totalPrice float64
	for _, id := range data.PartUuids {
		part, ok := byUUID[id]
		if !ok {
			return nil, &model.InvalidArgumentError{Message: fmt.Sprintf("part not found: %s", id)}
		}
		totalPrice += part.Info.Price
	}

	orderUUID := uuid.New().String()

	order := model.Order{
		UserUUID:   data.UserUUID,
		PartUuids:  data.PartUuids,
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
		Status:     string(model.StatusPendingPayment),
	}

	if err := s.OrderRepository.CreateOrder(ctx, order); err != nil {
		return nil, err
	}

	return &order, nil

}
