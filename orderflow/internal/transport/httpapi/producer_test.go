package httpapi

import (
	"context"
	"orderflow/internal/order/events"
)

type FakeProducer struct{}

func NewFakeProducer() *FakeProducer {
	return &FakeProducer{}
}

func (p *FakeProducer) PublishOrderConfirmed(context.Context, events.OrderConfirmedEvent) error {
	return nil
}
