package order

import (
	"testing"
)

func TestCalculateTotal(t *testing.T) {
	items := []OrderItem{
		{
			Product: Product{
				Name:  "Test1",
				Price: 100,
			},
			Quantity: 2,
		},
		{
			Product: Product{
				Name:  "Test2",
				Price: 50,
			},
			Quantity: 4,
		},
	}

	testOrder := Order{
		ID:     1,
		Status: "test-status",
		Items:  items,
	}

	if testOrder.Total() != 400 {
		t.Errorf("invalid func logic")
	}
}

func TestConfirm(t *testing.T) {
	items := []OrderItem{
		{
			Product: Product{
				Name:  "Test1",
				Price: 100,
			},
			Quantity: 2,
		},
		{
			Product: Product{
				Name:  "Test2",
				Price: 50,
			},
			Quantity: 4,
		},
	}

	testOrder := Order{
		ID:     1,
		Status: StatusNew,
		Items:  items,
	}

	if err := testOrder.Confirm(); err != nil {
		t.Errorf("fail to change status")
	}

	if testOrder.Status != StatusConfirmed {
		t.Errorf("invalid order status")
	}
}

func TestEmptyOrder(t *testing.T) {
	emptyOrder := Order{
		ID:     1,
		Status: StatusNew,
		Items:  []OrderItem{},
	}

	if err := emptyOrder.Confirm(); err == nil {
		t.Errorf("fail to empty order status change")
	}

	if emptyOrder.Status != StatusNew {
		t.Errorf("unexpected change status")
	}
}

func TestOrderValidate(t *testing.T) {
	items := []OrderItem{
		{
			Product: Product{
				Name:  "Test1",
				Price: 100,
			},
			Quantity: 2,
		},
		{
			Product: Product{
				Name:  "Test2",
				Price: 50,
			},
			Quantity: 4,
		},
	}

	testOrder := Order{
		ID:     1,
		Status: StatusNew,
		Items:  items,
	}

	err := testOrder.Validate()
	if err != nil {
		t.Errorf("invalid validate logic: %v", err)
	}
}

func TestOrderValidateQuantity(t *testing.T) {
	items := []OrderItem{
		{
			Product: Product{
				Name:  "Test1",
				Price: 100,
			},
			Quantity: 2,
		},
		{
			Product: Product{
				Name:  "Test2",
				Price: 50,
			},
			Quantity: 0,
		},
	}

	testOrder := Order{
		ID:     1,
		Status: StatusNew,
		Items:  items,
	}

	err := testOrder.Validate()
	if err == nil {
		t.Errorf("invalid validate quantity logic")
	}
}

func TestNewOrder(t *testing.T) {
	items := []OrderItem{
		{
			Product: Product{
				Name:  "Test1",
				Price: 100,
			},
			Quantity: 2,
		},
		{
			Product: Product{
				Name:  "Test2",
				Price: 50,
			},
			Quantity: 2,
		},
	}

	newOrder, err := New(items)
	if err != nil {
		t.Fatalf("invalid logic")
	}

	if newOrder.ID != 0 {
		t.Errorf("invalid id order")
	}

	if newOrder.Status != StatusNew {
		t.Errorf("invalid status")
	}
}

func TestNewOrderInvalidData(t *testing.T) {
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

	newOrder, err := New(items)
	if err == nil {
		t.Errorf("invalid logic")
	}

	if newOrder != nil {
		t.Errorf("invalid logic")
	}
}
