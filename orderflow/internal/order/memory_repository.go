package order

import (
	"context"
	"fmt"
)

type MemoryRepo struct {
	orders map[int64]*Order
	nextID int64
}

func NewRepo() *MemoryRepo {
	return &MemoryRepo{
		orders: make(map[int64]*Order),
		nextID: 1,
	}
}

func (repo *MemoryRepo) Save(ctx context.Context, order *Order) error {
	order.ID = repo.nextID
	repo.orders[order.ID] = order
	repo.nextID++
	return nil
}

func (repo MemoryRepo) GetByID(ctx context.Context, id int64) (*Order, error) {
	if order, ok := repo.orders[id]; ok {
		return order, nil
	}

	return nil, fmt.Errorf("order not found")
}
