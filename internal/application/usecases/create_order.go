package usecases

import (
	"context"

	"github.com/KedroPedro/realtime-order-analytics/internal/domain/entity"
	"github.com/KedroPedro/realtime-order-analytics/internal/domain/interfaces"
)

type CreateOrderUsecase struct {
	orders interfaces.Orders
}

func NewCreateOrderUsecase(orders interfaces.Orders) *CreateOrderUsecase {
	return &CreateOrderUsecase{
		orders: orders,
	}
}

func (uc *CreateOrderUsecase) Execute(ctx context.Context, order *entity.Order) error {
	return uc.orders.CreateOrder(ctx, order)
}
