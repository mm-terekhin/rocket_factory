package order

import (
	"github.com/jackc/pgx/v5"
	def "github.com/mm-terekhin/rocket_factory/order/internal/repository"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	con *pgx.Conn
}

func NewRepository(con *pgx.Conn) *repository {
	return &repository{
		con: con,
	}
}
