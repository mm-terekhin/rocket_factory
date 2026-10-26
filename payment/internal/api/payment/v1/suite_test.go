package v1

import (
	"context"
	"testing"

	"github.com/mm-terekhin/rocket_factory/payment/internal/service/mocks"
	"github.com/stretchr/testify/suite"
)

func testCtx() context.Context {
	return context.Background()
}

type APISuite struct {
	suite.Suite

	paymentService *mocks.PaymentService
	api            *paymentApi
}

func (s *APISuite) SetupTest() {
	s.paymentService = mocks.NewPaymentService(s.T())
	s.api = NewApi(s.paymentService)
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
