package interfaces

import "github.com/KedroPedro/realtime-order-analytics/internal/domain/entity"

type Orders interface {
	CreateOrder(order entity.Order) error
}
