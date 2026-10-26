package v1

import (
	"context"

	"github.com/mm-terekhin/rocket_factory/inventory/internal/converter"
	inventoryV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/proto/inventory/v1"
)

func (a *inventoryApi) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	part, err := a.inventoryService.GetPart(ctx, req.GetUuid())

	if err != nil {
		return nil, err
	}

	return &inventoryV1.GetPartResponse{Part: converter.PartToProto(part)}, err
}
