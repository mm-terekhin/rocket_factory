package v1

import (
	"context"
	"fmt"

	"github.com/mm-terekhin/rocket_factory/order/internal/model"
	generatedPaymentV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/proto/payment/v1"
)

func paymentMethodToProto(method model.PaymentMethod) (generatedPaymentV1.PaymentMethod, error) {
	switch method {
	case model.PaymentMethodUnspecified:
		return generatedPaymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED, nil
	case model.PaymentMethodCard:
		return generatedPaymentV1.PaymentMethod_PAYMENT_METHOD_CARD, nil
	case model.PaymentMethodSBP:
		return generatedPaymentV1.PaymentMethod_PAYMENT_METHOD_SBP, nil
	case model.PaymentMethodCreditCard:
		return generatedPaymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD, nil
	case model.PaymentMethodInvestorMoney:
		return generatedPaymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY, nil
	default:
		return generatedPaymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED, fmt.Errorf("unknown payment method: %d", method)
	}
}

func (c *client) PayOrder(ctx context.Context, data model.PayOrderData) (string, error) {
	paymentMethod, err := paymentMethodToProto(data.PaymentMethod)
	if err != nil {
		return "", err
	}

	resp, err := c.generatedClient.PayOrder(ctx, &generatedPaymentV1.PayOrderRequest{
		OrderUuid:     data.OrderUUID,
		UserUuid:      data.UserUuid,
		PaymentMethod: paymentMethod,
	})
	if err != nil {
		return "", err
	}

	return resp.GetTransactionUuid(), nil
}
