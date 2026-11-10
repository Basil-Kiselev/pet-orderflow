package order

import (
	"context"
	"testing"
)

func TestRepo(t *testing.T) {
	repo := NewRepo()
	ctx := context.Background()

	items := []OrderItem{
		{
			Product: Product{
				Name:  "Test1",
				Price: 100,
			},
			Quantity: 0,
		},
		{
			Product: Product{
				Name:  "Test2",
				Price: 50,
			},
			Quantity: 2,
		},
	}

	testOrder := &Order{
		Status: StatusNew,
		Items:  items,
	}

	secondOrder := &Order{
		Status: StatusNew,
		Items:  items,
	}

	err := repo.Save(ctx, testOrder)
	if err != nil {
		t.Errorf("fail to save order")
	}

	err = repo.Save(ctx, secondOrder)
	if err != nil {
		t.Errorf("fail to save order")
	}

	if ord, err := repo.GetByID(ctx, 1); err != nil || ord == nil {
		t.Errorf("fail to get order")
	}

	if secondOrder.ID != 2 {
		t.Errorf("wrong increment id")
	}
}

func TestRepoGetWrongOrder(t *testing.T) {
	repo := NewRepo()
	ctx := context.Background()

	ord, err := repo.GetByID(ctx, 1)

	if err == nil {
		t.Errorf("fail to get order")
	}

	if ord != nil {
		t.Errorf("fail to get order")
	}
}
