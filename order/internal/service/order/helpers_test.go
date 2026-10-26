package order

import (
	"errors"

	"github.com/mm-terekhin/rocket_factory/order/internal/model"
)

func (s *ServiceSuite) requireInvalidArgument(err error, message string) {
	s.Require().Error(err)

	var inv *model.InvalidArgumentError
	s.Require().True(errors.As(err, &inv))
	s.Require().Equal(message, inv.Message)
}
