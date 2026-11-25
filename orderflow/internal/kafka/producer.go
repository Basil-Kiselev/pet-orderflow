package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	order "orderflow/internal/order/events"
	"strconv"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	w *kafka.Writer
}

func NewProducer(cfg Config) *Producer {
	return &Producer{
		w: &kafka.Writer{
			Addr:         kafka.TCP(cfg.Brokers...),
			Topic:        cfg.Topic,
			RequiredAcks: kafka.RequireAll,
		},
	}
}

func (p *Producer) Close() error {
	return p.w.Close()
}

func (p *Producer) PublishOrderConfirmed(ctx context.Context, event order.OrderConfirmedEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("fail to marshaling data: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(strconv.FormatInt(event.OrderID, 10)),
		Value: data,
	}

	err = p.w.WriteMessages(ctx, msg)
	if err != nil {
		return fmt.Errorf("fail to write msg: %w", err)
	}

	return nil
}
