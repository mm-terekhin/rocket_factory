package v1

import (
	"context"

	"github.com/mm-terekhin/rocket_factory/inventory/internal/converter"
	inventoryV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/proto/inventory/v1"
)

func (a *inventoryApi) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	parts, err := a.inventoryService.ListParts(ctx, converter.FilterToModel(req.GetFilter()))

	if err != nil {
		return nil, err
	}

	partsProto := make([]*inventoryV1.Part, 0, len(parts))

	for _, part := range parts {

		partsProto = append(partsProto, converter.PartToProto(part))
	}

	return &inventoryV1.ListPartsResponse{Parts: partsProto}, nil
}
