package model

import "errors"

var ErrorOrderNotFound = errors.New("Order not found")
var ErrorOrderConflict = errors.New("The order has already been paid for and cannot be cancelled.")

type InvalidArgumentError struct{ Message string }

func (e *InvalidArgumentError) Error() string { return e.Message }
