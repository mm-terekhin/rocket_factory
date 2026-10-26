package v1

import (
	"context"

	"github.com/mm-terekhin/rocket_factory/inventory/internal/converter"
	inventoryV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/proto/inventory/v1"
)

func (a *inventoryApi) CreatePart(ctx context.Context, req *inventoryV1.CreatePartRequest) (*inventoryV1.CreatePartResponse, error) {
	resp, err := a.inventoryService.CreatePart(ctx, converter.PartInfoToModel(req.GetInfo()))
	if err != nil {
		return nil, err
	}

	return &inventoryV1.CreatePartResponse{Uuid: resp}, nil
}
