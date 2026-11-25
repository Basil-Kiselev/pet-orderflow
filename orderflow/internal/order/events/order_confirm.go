package events

import (
	"context"
	"time"
)

type EventPublisher interface {
	PublishOrderConfirmed(context.Context, OrderConfirmedEvent) error
}

type OrderConfirmedEvent struct {
	OrderID    int64     `json:"order_id"`
	Total      int       `json:"total"`
	OccurredAt time.Time `json:"occurred_at"`
}
