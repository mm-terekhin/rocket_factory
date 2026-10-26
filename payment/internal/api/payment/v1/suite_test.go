package v1

import (
	"context"
	"testing"

	"github.com/mm-terekhin/rocket_factory/payment/internal/service/mocks"
	"github.com/stretchr/testify/suite"
)

type APISuite struct {
	suite.Suite

	ctx context.Context

	paymentService *mocks.PaymentService
	api            *paymentApi
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()
	s.paymentService = mocks.NewPaymentService(s.T())
	s.api = NewApi(s.paymentService)
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
