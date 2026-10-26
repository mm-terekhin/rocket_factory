package v1

import (
	"errors"

	"github.com/mm-terekhin/rocket_factory/payment/internal/model"
	paymentV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func badRequest(message string) error {
	return status.Error(codes.InvalidArgument, message)
}

func asInvalidArgument(err error) (*model.InvalidArgumentError, bool) {
	var e *model.InvalidArgumentError
	return e, errors.As(err, &e)
}

func mapPayOrderErr(err error) (*paymentV1.PayOrderResponse, error) {
	if inv, ok := asInvalidArgument(err); ok {
		return nil, badRequest(inv.Message)
	}
	return nil, err
}
