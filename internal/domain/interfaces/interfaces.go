package interfaces

import (
	"context"

	"github.com/KedroPedro/realtime-order-analytics/internal/domain/entity"
)

type Orders interface {
	CreateOrder(cxt context.Context, order *entity.Order) error
}

type OrdersOutbox interface {
	GetEventsBatch(ctx context.Context, locker string) ([]entity.OrderOutbox, error)
	PublishEvents(ctx context.Context, events []entity.OrderOutbox) error
}

type EventPublisher interface {
	PublishBatch(ctx context.Context, events []entity.OrderOutbox) error
}
