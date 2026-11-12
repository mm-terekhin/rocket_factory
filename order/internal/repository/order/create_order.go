package order

import (
	"context"
	"log"

	"github.com/mm-terekhin/rocket_factory/order/internal/model"
	"github.com/mm-terekhin/rocket_factory/order/internal/repository/converter"
)

func (r *repository) CreateOrder(ctx context.Context, order model.Order) error {
	tx, err := r.con.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		cerr := tx.Rollback(ctx)
		if cerr != nil {
			log.Printf("failed to rollback transaction: %v\n", cerr)
		}
	}()

	repoOrder := converter.OrderToRepo(order)

	if _, err := tx.Exec(
		ctx,
		"INSERT INTO orders (order_uuid, user_uuid, total_price, transaction_uuid, payment_method, status) VALUES ($1, $2, $3, $4, $5, $6)",
		repoOrder.OrderUUID, repoOrder.UserUUID, repoOrder.TotalPrice, repoOrder.TransactionUUID, repoOrder.PaymentMethod, repoOrder.Status,
	); err != nil {
		return err
	}

	for _, part := range repoOrder.PartUuids {
		if _, err := r.con.Exec(
			ctx,
			"INSERT INTO order_parts (part_uuid, order_uuid) VALUES ($1, $2)",
			part, repoOrder.OrderUUID,
		); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
