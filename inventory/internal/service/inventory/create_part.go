package inventory

import (
	"context"

	"github.com/mm-terekhin/rocket_factory/inventory/internal/model"
)

func (s *service) CreatePart(ctx context.Context, info model.PartInfo) (string, error) {
	uuid, err := s.InventoryRepository.CreatePart(ctx, info)

	if err != nil {
		return "", err
	}

	return uuid, nil
}
