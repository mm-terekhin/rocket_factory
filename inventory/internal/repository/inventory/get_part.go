package inventory

import (
	"context"

	"github.com/mm-terekhin/rocket_factory/inventory/internal/model"
	"github.com/mm-terekhin/rocket_factory/inventory/internal/repository/converter"
)

func (r *repository) GetPart(_ context.Context, uuid string) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	part, ok := r.parts[uuid]

	if !ok {
		return model.Part{}, model.ErrPartNotFound
	}

	return converter.PartToModel(part), nil
}
