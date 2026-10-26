package converter

import (
	"github.com/mm-terekhin/rocket_factory/order/internal/model"
	orderV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/openapi/order/v1"
)

func CreateOrderDataToModel(data *orderV1.CreateOrderRequest) model.CreateOrderData {
	return model.CreateOrderData{
		UserUUID:  data.UserUUID,
		PartUuids: data.PartUuids,
	}
}

func OrderToResponse(order *model.Order) *orderV1.Order {
	return &orderV1.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: optStringFromPtr(order.TransactionUUID),
		PaymentMethod:   optStringFromPtr(order.PaymentMethod),
		Status:          order.Status,
	}
}

func optStringFromPtr(p *string) orderV1.OptString {
	if p == nil {
		return orderV1.OptString{} // Set == false
	}
	return orderV1.NewOptString(*p)
}

func PaymentMethodFromHTTP(s string) model.PaymentMethod {
	switch s {
	case "CARD", "PAYMENT_METHOD_CARD":
		return model.PaymentMethodCard
	case "SBP", "PAYMENT_METHOD_SBP":
		return model.PaymentMethodSBP
	case "CREDIT_CARD", "PAYMENT_METHOD_CREDIT_CARD":
		return model.PaymentMethodCreditCard
	case "INVESTOR_MONEY", "PAYMENT_METHOD_INVESTOR_MONEY":
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnspecified
	}
}
