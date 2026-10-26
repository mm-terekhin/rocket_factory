package model

import "errors"

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrOrderConflict = errors.New("the order has already been paid for and cannot be cancelled")
)

type InvalidArgumentError struct{ Message string }

func (e *InvalidArgumentError) Error() string { return e.Message }
