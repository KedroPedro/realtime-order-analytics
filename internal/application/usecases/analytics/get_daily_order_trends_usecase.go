package analyticsuc

import (
	"context"

	"github.com/KedroPedro/realtime-order-analytics/internal/domain/entity"
	"github.com/KedroPedro/realtime-order-analytics/internal/domain/interfaces"
)

type GetDailyOrderTrendsUsecase struct {
	analytic interfaces.OrderAnalytic
}

func NewGetDailyOrderTrendsUsecase(
	analytic interfaces.OrderAnalytic,
) *GetDailyOrderTrendsUsecase {
	return &GetDailyOrderTrendsUsecase{
		analytic: analytic,
	}
}

func (uc *GetDailyOrderTrendsUsecase) Execute(ctx context.Context) ([]entity.DailyOrderTrend, error) {
	return uc.analytic.GetDailyOrderTrends(ctx)
}
