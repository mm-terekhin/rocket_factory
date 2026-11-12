package order

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/mm-terekhin/rocket_factory/order/internal/model"
	"github.com/mm-terekhin/rocket_factory/order/internal/repository/converter"
	repoModel "github.com/mm-terekhin/rocket_factory/order/internal/repository/model"
)

func (r *repository) GetOrder(ctx context.Context, uuid string) (*model.Order, error) {
	repoOrder, err := getOrders(ctx, r.con, uuid)
	if err != nil {
		return nil, err
	}

	parts, err := getOrderParts(ctx, r.con, uuid)
	if err != nil {
		return nil, err
	}

	repoOrder.PartUuids = parts

	order := converter.OrderToModel(*repoOrder)

	return &order, nil
}

func getOrders(ctx context.Context, con *pgx.Conn, uuid string) (*repoModel.Order, error) {
	var repoOrder repoModel.Order

	err := con.QueryRow(ctx, `
        SELECT order_uuid, user_uuid, total_price, transaction_uuid, payment_method, status
        FROM orders
        WHERE order_uuid = $1
    `, uuid).Scan(
		&repoOrder.OrderUUID,
		&repoOrder.UserUUID,
		&repoOrder.TotalPrice,
		&repoOrder.TransactionUUID,
		&repoOrder.PaymentMethod,
		&repoOrder.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrOrderNotFound
		}
		return nil, err
	}

	return &repoOrder, nil
}

func getOrderParts(ctx context.Context, con *pgx.Conn, uuid string) ([]string, error) {
	rows, err := con.Query(ctx, `
    SELECT part_uuid FROM order_parts WHERE order_uuid = $1
`, uuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	parts := make([]string, 0)
	for rows.Next() {
		var partUUID string
		if err := rows.Scan(&partUUID); err != nil {
			return nil, err
		}
		parts = append(parts, partUUID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return parts, nil
}
