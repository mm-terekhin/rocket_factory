package v1

import (
	"github.com/mm-terekhin/rocket_factory/payment/internal/service"
	paymentV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/proto/payment/v1"
)

type paymentApi struct {
	paymentV1.UnimplementedPaymentServiceServer

	paymentService service.PaymentService
}

func NewApi(paymentService service.PaymentService) *paymentApi {
	return &paymentApi{paymentService: paymentService}
}
