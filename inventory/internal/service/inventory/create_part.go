package inventory

import (
	"context"
	"time"
	"uuid"

	"github.com/mm-terekhin/rocket_factory/inventory/internal/model"
)

func (s *service) CreatePart(ctx context.Context, info model.PartInfo) (string, error) {
	partUUID := uuid.New().String()
	now := time.Now()

	part := model.Part{
		UUID:      partUUID,
		Info:      info,
		CreatedAt: now,
	}

	_, err := s.InventoryRepository.CreatePart(ctx, part)
	if err != nil {
		return "", err
	}

	return partUUID, nil
}
