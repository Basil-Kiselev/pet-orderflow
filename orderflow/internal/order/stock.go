package order

import "context"

type Stocker interface {
	GetStock(context.Context, string) (int32, error)
}
