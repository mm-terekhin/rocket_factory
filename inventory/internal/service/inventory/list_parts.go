package inventory

import (
	"context"

	"github.com/mm-terekhin/rocket_factory/inventory/internal/model"
)

func (s *service) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	parts, err := s.InventoryRepository.ListParts(ctx, filter)
	if err != nil {
		return make([]model.Part, 0), err
	}

	return parts, nil
}
