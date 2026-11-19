package order

import "context"

type Repository interface {
	Save(context.Context, *Order) error
	GetByID(context.Context, int64) (*Order, error)
	UpdateStatus(context.Context, int64, Status) error
}
