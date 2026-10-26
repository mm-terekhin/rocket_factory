package model

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
