package payment

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

func testCtx() context.Context {
	return context.Background()
}

type ServiceSuite struct {
	suite.Suite

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.service = NewService()
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
