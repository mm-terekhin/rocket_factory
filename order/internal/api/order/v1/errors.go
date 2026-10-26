package v1

import (
	"errors"

	"github.com/mm-terekhin/rocket_factory/order/internal/model"
	orderV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/openapi/order/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func badRequest(message string) *orderV1.BadRequestError {
	return &orderV1.BadRequestError{Code: 400, Message: message}
}

func notFound(orderUUID string) *orderV1.NotFoundError {
	return &orderV1.NotFoundError{
		Code:    404,
		Message: "Order for uuid'" + orderUUID + "' not found",
	}
}

func conflict(message string) *orderV1.ConflictError {
	return &orderV1.ConflictError{Code: 409, Message: message}
}

func asInvalidArgument(err error) (*model.InvalidArgumentError, bool) {
	var e *model.InvalidArgumentError
	return e, errors.As(err, &e)
}

func isOrderNotFound(err error) bool {
	return errors.Is(err, model.ErrOrderNotFound)
}

func isOrderConflict(err error) bool {
	return errors.Is(err, model.ErrOrderConflict)
}

func asGRPCBadRequest(err error) (*orderV1.BadRequestError, bool) {
	st, ok := status.FromError(err)
	if !ok || st.Code() != codes.InvalidArgument {
		return nil, false
	}
	return badRequest(st.Message()), true
}

func mapCreateOrderErr(err error) (orderV1.CreateOrderRes, error) {
	if inv, ok := asInvalidArgument(err); ok {
		return badRequest(inv.Message), nil
	}
	return nil, err
}

func mapGetOrderErr(err error, orderUUID string) (orderV1.GetOrderByUUIDRes, error) {
	if isOrderNotFound(err) {
		return notFound(orderUUID), nil
	}
	if inv, ok := asInvalidArgument(err); ok {
		return badRequest(inv.Message), nil
	}
	return nil, err
}

func mapCancelErr(err error, orderUUID string) (orderV1.CancelOrderByUUIDRes, error) {
	if isOrderNotFound(err) {
		return notFound(orderUUID), nil
	}
	if isOrderConflict(err) {
		return conflict(err.Error()), nil
	}
	return nil, err
}

func mapPayErr(err error, orderUUID string) (orderV1.PayOrderByUUIDRes, error) {
	if isOrderNotFound(err) {
		return notFound(orderUUID), nil
	}
	if br, ok := asGRPCBadRequest(err); ok {
		return br, nil
	}
	if inv, ok := asInvalidArgument(err); ok {
		return badRequest(inv.Message), nil
	}
	return nil, err
}
