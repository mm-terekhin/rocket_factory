package model

type PayOrderInput struct {
	OrderUuid     string
	UserUuid      string
	PaymentMethod PaymentMethod
}

type PaymentMethod int32

const (
	Unspecified PaymentMethod = iota
	Card
	Sbp
	CreditCard
	InvestorMoney
)
