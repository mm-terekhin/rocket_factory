package order

import (
	"github.com/mm-terekhin/rocket_factory/order/internal/client"
	"github.com/mm-terekhin/rocket_factory/order/internal/repository"
	def "github.com/mm-terekhin/rocket_factory/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	OrderRepository repository.OrderRepository
	InventoryClient client.InventoryClient
	PaymentClient   client.PaymentClient
}

func NewService(
	repository repository.OrderRepository,
	inventoryClient client.InventoryClient,
	paymentClient client.PaymentClient,
) *service {
	return &service{
		OrderRepository: repository,
		InventoryClient: inventoryClient,
		PaymentClient:   paymentClient,
	}
}
