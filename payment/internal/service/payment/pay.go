package payment

import (
	"context"
	"log"
	"uuid"

	"github.com/mm-terekhin/rocket_factory/payment/internal/model"
)

func (s *service) Pay(ctx context.Context, input model.PayOrderInput) (*string, error) {
	isValid := validate(input)

	if isValid != nil {
		return nil, isValid
	}

	transaction_uuid := uuid.New().String()

	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transaction_uuid)

	return &transaction_uuid, nil
}

func validate(input model.PayOrderInput) error {
	if input.OrderUuid == "" {
		return &model.InvalidArgumentError{Message: "order_uuid is inputuired"}
	}

	if input.UserUuid == "" {
		return &model.InvalidArgumentError{Message: "user_uuid is inputuired"}
	}

	if input.PaymentMethod == model.Unspecified {
		return &model.InvalidArgumentError{Message: "payment_method is inputuired"}
	}

	return nil
}
