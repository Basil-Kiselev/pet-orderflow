package order

import (
	"context"
	"testing"
)

type FakeStock struct {
	stock map[string]int32
}

func NewFakeStock() *FakeStock {
	return &FakeStock{
		stock: map[string]int32{
			"Test1": 52,
			"Test2": 24,
		},
	}
}

func TestService(t *testing.T) {
	repo := NewRepo()

	service := NewService(repo, NewFakeStock())
	ctx := context.Background()

	items := []OrderItem{
		{
			Product: Product{
				Name:  "Test1",
				Price: 100,
			},
			Quantity: 1,
		},
		{
			Product: Product{
				Name:  "Test2",
				Price: 50,
			},
			Quantity: 2,
		},
	}

	newOrder, err := service.Create(ctx, items)
	if err != nil {
		t.Errorf("fail to create order")
	}

	if newOrder == nil {
		t.Errorf("fail to create order")
	}

	gettingOrder, err := service.GetByID(ctx, 1)
	if err != nil {
		t.Fatalf("fail to get order")
	}

	if gettingOrder == nil {
		t.Fatalf("fail to get order")
	}

	if gettingOrder.ID != 1 {
		t.Errorf("wrong order id")
	}
}

func TestServiceConfirm(t *testing.T) {
	repo := NewRepo()
	service := NewService(repo, NewFakeStock())
	ctx := context.Background()

	items := []OrderItem{
		{
			Product: Product{
				Name:  "Test1",
				Price: 100,
			},
			Quantity: 1,
		},
		{
			Product: Product{
				Name:  "Test2",
				Price: 50,
			},
			Quantity: 2,
		},
	}

	newOrder, err := service.Create(ctx, items)
	if err != nil {
		t.Fatalf("fail to create order")
	}

	updOrd, err := service.Confirm(ctx, newOrder.ID)
	if err != nil {
		t.Fatalf("fail to confirm order: %v", err)
	}

	if updOrd.Status != StatusConfirmed {
		t.Fatal("wrong status after update")
	}
}

func (f FakeStock) GetStock(ctx context.Context, productName string) (int32, error) {
	return f.stock[productName], nil
}
