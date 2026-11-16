package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"orderflow/internal/order"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type MockRepo struct {
	client *redis.Client
	ttl    time.Duration
}

func TestGetByIDMiss(t *testing.T) {
	ctx := context.Background()
	ordRepo := order.NewRepo()
	redisRepo := newTestRepo(ctx, t, ordRepo)

	ord, err := order.New([]order.OrderItem{
		{
			Product: order.Product{
				Name:  "Test item",
				Price: 200,
			},
			Quantity: 4,
		},
	})
	if err != nil {
		t.Fatalf("create order: %v", err)
	}

	// id старт с 1
	err = ordRepo.Save(ctx, ord)
	if err != nil {
		t.Fatalf("fail to save order: %v", err)
	}

	// очищаем ддля уверенности что ключа нет
	err = redisRepo.client.Del(ctx, "order:1").Err()
	if err != nil {
		t.Fatalf("fail to clear cache: %v", err)
	}

	ordFromRedis, err := redisRepo.GetByID(ctx, 1)
	if err != nil {
		t.Fatalf("fail to get order from redisRepo: %v", err)
	}

	if ordFromRedis.ID != 1 {
		t.Error("wrong order id after redis")
	}

	res, err := redisRepo.client.Get(ctx, "order:1").Bytes()
	if err != nil {
		t.Fatalf("fail to check redis order: %v", err)
	}

	cachedOrd := &order.Order{}
	err = json.Unmarshal(res, cachedOrd)
	if err != nil {
		t.Fatal("fail to unmarshaling cached order")
	}

	if cachedOrd.ID != 1 {
		t.Error("wrong order id in cache")
	}
}

func TestGetByIDHit(t *testing.T) {
	ctx := context.Background()
	ordRepo := order.NewRepo()
	redisRepo := newTestRepo(ctx, t, ordRepo)

	ord, err := order.New([]order.OrderItem{
		{
			Product: order.Product{
				Name:  "Test item",
				Price: 200,
			},
			Quantity: 4,
		},
	})
	if err != nil {
		t.Fatalf("create order: %v", err)
	}

	// Костыль на совпадение с тестом miss, для очистки
	ord.ID = 1

	jOrder, err := json.Marshal(ord)
	if err != nil {
		t.Fatalf("fail to data marshaling: %v", err)
	}

	key := fmt.Sprintf("order:%d", ord.ID)
	err = redisRepo.client.Set(ctx, key, jOrder, redisRepo.ttl).Err()
	if err != nil {
		t.Fatalf("fail to set cache: %v", err)
	}

	resOrd, err := redisRepo.GetByID(ctx, ord.ID)
	if err != nil {
		t.Fatalf("fail take order from cache: %v", err)
	}

	if resOrd.ID != ord.ID {
		t.Error("wrong ID resOrd after cache")
	}

}

func newTestRepo(ctx context.Context, t *testing.T, ordRepo order.Repository) *Repository {
	t.Helper()

	_ = godotenv.Load("../../../../.env")

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		t.Skip("REDIS_ADDR is not set")
	}

	redisCli := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   15,
	})

	t.Cleanup(func() {
		redisCli.Del(ctx, "order:1")
		redisCli.Close()

	})

	return &Repository{
		ordRepo: ordRepo,
		ttl:     30 * time.Second,
		client:  redisCli,
	}
}
