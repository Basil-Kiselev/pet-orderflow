package order

import "fmt"

const (
	StatusNew       Status = "new"
	StatusConfirmed Status = "confirmed"
)

type Product struct {
	Name  string
	Price int
}

type OrderItem struct {
	Product  Product
	Quantity int
}

type Order struct {
	ID     int64
	Status Status
	Items  []OrderItem
}

type Status string

func (o Order) Total() int {
	total := 0
	for _, item := range o.Items {
		total += item.Quantity * item.Product.Price
	}
	return total
}

func (o *Order) Confirm() error {
	if err := o.Validate(); err != nil {
		return fmt.Errorf("bad order: %v", err)
	}
	o.Status = StatusConfirmed
	return nil
}

func (o Order) Validate() error {
	if len(o.Items) == 0 {
		return fmt.Errorf("Empty order")
	}

	for _, item := range o.Items {
		if item.Quantity < 1 {
			return fmt.Errorf("wrong quantity")
		}

		if item.Product.Name == "" {
			return fmt.Errorf("empty product name")
		}

		if item.Product.Price < 1 {
			return fmt.Errorf("bad price in product")
		}
	}

	return nil
}

func New(items []OrderItem) (*Order, error) {
	newOrder := &Order{
		Status: StatusNew,
		Items:  items,
	}

	err := newOrder.Validate()
	if err != nil {
		return nil, err
	}

	return newOrder, nil
}
