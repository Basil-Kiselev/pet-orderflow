package order

import (
	"context"
	"fmt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
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
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fail to get order: %w", err)
	}

	return order, nil
}
