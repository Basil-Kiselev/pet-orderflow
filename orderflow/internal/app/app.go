package app

import (
	"context"
	"fmt"
	"net/http"
	"orderflow/internal/order"
	"orderflow/internal/storage/postgres"
	"orderflow/internal/storage/redis"
	"orderflow/internal/transport/httpapi"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	goredis "github.com/redis/go-redis/v9"
)

const RedisTTL = 1 * time.Minute

func Run() error {
	ctx := context.Background()
	_ = godotenv.Load("../.env")

	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("fail to create db pool: %w", err)
	}

	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		fmt.Printf("fail to ping db: %v", err)
		return err
	}

	redisCli := goredis.NewClient(&goredis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	defer redisCli.Close()

	redisRepo := redis.New(postgres.NewRepo(pool), redisCli, RedisTTL)

	httpHandler := httpapi.NewHandler(order.NewService(redisRepo))

	fmt.Println("OrderFlow server started on :8080")

	err = http.ListenAndServe(":8080", httpHandler.Routes())
	if err != nil {
		fmt.Printf("failed to start server: %v", err)
		return err
	}

	return nil
}
