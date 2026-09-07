package interfaces

import (
	"context"

	"github.com/KedroPedro/realtime-order-analytics/internal/domain/entity"
)

type Orders interface {
	CreateOrder(cxt context.Context, order *entity.Order) error
}
