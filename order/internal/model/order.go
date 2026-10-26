package model

import "time"

type Order struct {
	// Уникальный идентификатор заказа.
	OrderUUID string
	// UUID пользователя.
	UserUUID string
	// Список UUID деталей.
	PartUuids []string
	// Итоговая стоимость.
	TotalPrice float64
	// UUID транзакции (если оплачен).
	TransactionUUID *string
	// Способ оплаты (если оплачен).
	PaymentMethod *string
	// Статус.
	Status string
}

type PartsFilter struct {
	Uuids                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}

type PaymentMethod int

const (
	PaymentMethodUnspecified PaymentMethod = iota
	PaymentMethodCard
	PaymentMethodSBP
	PaymentMethodCreditCard
	PaymentMethodInvestorMoney
)

type PayOrderInput struct {
	OrderUUID     string
	PaymentMethod PaymentMethod
}

type PayOrderData struct {
	OrderUUID     string
	UserUuid      string
	PaymentMethod PaymentMethod
}

type CreateOrderData struct {
	UserUUID  string
	PartUuids []string
}

type Part struct {
	UUID      string
	Info      PartInfo
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type PartInfo struct {
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      Category
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
	Tags          []string
	Metadata      map[string]any
}

type Status string

const (
	StatusPendingPayment Status = "PENDING_PAYMENT"
	StatusPaid           Status = "PAID"
	StatusCancelled      Status = "CANCELLED"
)

type Category int

const (
	CategoryEngine Category = iota + 1
	CategoryFuel
	CategoryPorthole
	CategoryWing
)

type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

type Value struct {
	StringValue *string
	Int64Value  *int64
	DoubleValue *float64
	BoolValue   *bool
}
