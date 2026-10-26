package service

import (
	"context"

	"github.com/mm-terekhin/rocket_factory/payment/internal/model"
)

type PaymentService interface {
	Pay(context.Context, model.PayOrderInput) (*string, error)
}
