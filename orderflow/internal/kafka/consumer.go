package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	order "orderflow/internal/order/events"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	r *kafka.Reader
}

func NewConsumer(cfg Config) *Consumer {
	return &Consumer{
		r: kafka.NewReader(
			kafka.ReaderConfig{
				Brokers: cfg.Brokers,
				GroupID: cfg.GroupID,
				Topic:   cfg.Topic,
			}),
	}
}

func (c *Consumer) Close() error {
	return c.r.Close()
}

func (c *Consumer) Run(ctx context.Context) error {
	for {
		var msg order.OrderConfirmedEvent
		rawMsg, err := c.r.FetchMessage(ctx)
		if err != nil {
			return fmt.Errorf("fail to read kafka message: %w", err)
		}

		err = json.Unmarshal(rawMsg.Value, &msg)
		if err != nil {
			return fmt.Errorf("fail to unmarshaling msg: %w", err)
		}

		fmt.Printf("order - %d, total - %d confimed\n", msg.OrderID, msg.Total)

		err = c.r.CommitMessages(ctx, rawMsg)
		if err != nil {
			return fmt.Errorf("fail to commit msg: %w", err)
		}
	}
}
