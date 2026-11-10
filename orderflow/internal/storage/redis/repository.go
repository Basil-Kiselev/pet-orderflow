package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"orderflow/internal/order"
	"time"

	"github.com/redis/go-redis/v9"
)

type Repository struct {
	sqlRepo order.Repository
	client  *redis.Client
	ttl     time.Duration
}

func New(sqlRepo order.Repository, client *redis.Client, ttl time.Duration) *Repository {
	return &Repository{
		sqlRepo: sqlRepo,
		client:  client,
		ttl:     ttl,
	}
}

func (r Repository) Save(ctx context.Context, order *order.Order) error {
	// пока заглушка
	return r.sqlRepo.Save(ctx, order)
}

// на стадии разработки
func (r Repository) GetByID(ctx context.Context, id int64) (*order.Order, error) {
	ord := &order.Order{}
	key := fmt.Sprintf("order:%d", id)

	res, err := r.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		ord, err = r.sqlRepo.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}

		jOrder, err := json.Marshal(ord)
		if err != nil {
			// лог
			return ord, nil
		}
		r.client.Set(ctx, key, jOrder, r.ttl)

		return ord, nil
	} else if err != nil {
		// добавить лог, что редис давн
		return r.sqlRepo.GetByID(ctx, id)
	}

	err = json.Unmarshal(res, ord)

	return ord, err
}
