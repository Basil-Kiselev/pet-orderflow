package order

import "errors"

var ErrStockUnavailable = errors.New("stock unavailable")
var ErrInsufficientStock = errors.New("insufficient stock")
var ErrorderNotFound = errors.New("order not found")
