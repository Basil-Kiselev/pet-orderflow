package httpapi

import "context"

type Stocker interface {
	GetStock(context.Context, string) (int32, error)
}

type FakeStock struct {
	stock map[string]int32
}

func NewFakeStock() *FakeStock {
	return &FakeStock{
		stock: map[string]int32{
			"Keyboard": 52,
			"Mouse":    24,
		},
	}
}

func (f FakeStock) GetStock(ctx context.Context, productName string) (int32, error) {
	return f.stock[productName], nil
}
