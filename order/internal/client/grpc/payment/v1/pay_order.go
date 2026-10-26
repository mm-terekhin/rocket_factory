package v1

import (
	"context"
	"github.com/mm-terekhin/rocket_factory/order/internal/model"
	generatedPaymentV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, data model.PayOrderData) (string, error) {
	resp, err := c.generatedClient.PayOrder(ctx, &generatedPaymentV1.PayOrderRequest{
		OrderUuid:     data.OrderUUID,
		UserUuid:      data.UserUuid,
		PaymentMethod: generatedPaymentV1.PaymentMethod(data.PaymentMethod),
	})
	if err != nil {
		return "", err
	}

	return resp.GetTransactionUuid(), nil
}
