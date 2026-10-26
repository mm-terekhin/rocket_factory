package v1

import (
	"context"

	clientConverter "github.com/mm-terekhin/rocket_factory/order/internal/client/converter"
	"github.com/mm-terekhin/rocket_factory/order/internal/model"
	generatedInventoryV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	resp, err := c.generatedClient.ListParts(ctx, &generatedInventoryV1.ListPartsRequest{
		Filter: clientConverter.PartsFilterToProto(filter),
	})

	if err != nil {
		return nil, err
	}

	return clientConverter.PartListToModel(resp.Parts), nil
}
