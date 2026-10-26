package v1

import (
	"context"

	"github.com/mm-terekhin/rocket_factory/payment/internal/model"
	paymentV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/proto/payment/v1"
)

func (a *paymentApi) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	resp, err := a.paymentService.Pay(ctx, model.PayOrderInput{
		OrderUuid:     req.GetOrderUuid(),
		UserUuid:      req.UserUuid,
		PaymentMethod: model.PaymentMethod(req.GetPaymentMethod()),
	})
	if err != nil {
		return mapPayOrderErr(err)
	}

	return &paymentV1.PayOrderResponse{TransactionUuid: *resp}, nil
}
