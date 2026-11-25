package main

import (
	"context"
	"fmt"
	"orderflow/internal/kafka"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	_ = godotenv.Load("../.env")

	consumer := kafka.NewConsumer(kafka.Config{
		Brokers: []string{os.Getenv("KAFKA_ADDR")},
		Topic:   "order.confirmed",
		GroupID: "notifications",
	})

	defer consumer.Close()

	if err := consumer.Run(ctx); err != nil {
		fmt.Printf("consumer error: %v\n", err)
	}
}
