package analyticsuc

import (
	"context"

	"github.com/KedroPedro/realtime-order-analytics/internal/domain/entity"
	"github.com/KedroPedro/realtime-order-analytics/internal/domain/interfaces"
)

type GetTopClientsUsecase struct {
	analytic interfaces.OrderAnalytic
}

func NewGetTopClientsUsecase(
	analytic interfaces.OrderAnalytic,
) *GetTopClientsUsecase {
	return &GetTopClientsUsecase{
		analytic: analytic,
	}
}

func (uc *GetTopClientsUsecase) Execute(ctx context.Context) ([]entity.TopClient, error) {
	return uc.analytic.GetTopClients(ctx)
}
