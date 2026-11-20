package order

import (
	"context"
	"fmt"
	"time"
)

type Service struct {
	repo  Repository
	stock Stocker
}

func NewService(repo Repository, stock Stocker) *Service {
	return &Service{
		repo:  repo,
		stock: stock,
	}
}

func (s *Service) Create(ctx context.Context, items []OrderItem) (*Order, error) {
	newOrder, err := New(items)
	if err != nil {
		return nil, fmt.Errorf("fail to create order: %w", err)
	}

	err = s.repo.Save(ctx, newOrder)
	if err != nil {
		return nil, fmt.Errorf("fail to save order: %w", err)
	}

	return newOrder, nil
}

func (s Service) GetByID(ctx context.Context, id int64) (*Order, error) {
	ord, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fail to get order: %w", err)
	}

	return ord, nil
}

func (s *Service) Confirm(ctx context.Context, id int64) (*Order, error) {
	ord, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fail to get order: %w", err)
	}

	reqCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	check, err := s.checkAvailableItems(reqCtx, ord)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrStockUnavailable, err)
	} else if !check {
		return nil, ErrInsufficientStock
	}

	err = ord.Confirm()
	if err != nil {
		return nil, fmt.Errorf("fail to confirm order: %w", err)
	}

	err = s.repo.UpdateStatus(ctx, id, ord.Status)
	if err != nil {
		return nil, fmt.Errorf("fail to update status: %w", err)
	}

	return ord, nil
}

func (s *Service) checkAvailableItems(ctx context.Context, ord *Order) (bool, error) {
	for _, item := range ord.Items {
		quant, err := s.stock.GetStock(ctx, item.Product.Name)
		if err != nil {
			return false, fmt.Errorf("fail to check stock: %w", err)
		}

		if quant < int32(item.Quantity) {
			return false, nil
		}
	}

	return true, nil
}
