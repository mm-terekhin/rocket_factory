package order

import (
	"context"

	"github.com/mm-terekhin/rocket_factory/order/internal/model"
	"github.com/mm-terekhin/rocket_factory/order/internal/repository/converter"
)

func (r *repository) UpdateOrder(ctx context.Context, order model.Order) error {
	repoOrder := converter.OrderToRepo(order)

	tag, err := r.con.Exec(ctx, `
		UPDATE orders
		SET status = $2,
		transaction_uuid = $3,
		payment_method = $4
		WHERE order_uuid = $1
	`,
		repoOrder.OrderUUID,
		repoOrder.Status,
		repoOrder.TransactionUUID,
		repoOrder.PaymentMethod,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return model.ErrOrderNotFound
	}

	return nil
}
