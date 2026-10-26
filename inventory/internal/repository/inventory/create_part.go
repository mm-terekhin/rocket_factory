package inventory

import (
	"context"
	"time"
	"uuid"

	"github.com/mm-terekhin/rocket_factory/inventory/internal/converter"
	"github.com/mm-terekhin/rocket_factory/inventory/internal/model"
	repoModel "github.com/mm-terekhin/rocket_factory/inventory/internal/repository/model"
	"github.com/samber/lo"
)

func (r *repository) CreatePart(_ context.Context, info model.PartInfo) (string, error) {
	id := uuid.New().String()
	now := time.Now()

	part := repoModel.Part{
		UUID:      id,
		CreatedAt: now,
		Info:      converter.PartInfoToRepo(info),
		UpdatedAt: lo.ToPtr(now),
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.parts[id] = part

	return id, nil
}
