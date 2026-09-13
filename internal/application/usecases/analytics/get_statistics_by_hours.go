package analyticsuc

import (
	"context"

	"github.com/KedroPedro/realtime-order-analytics/internal/domain/entity"
	"github.com/KedroPedro/realtime-order-analytics/internal/domain/interfaces"
)

type GetStatisticsByHoursUsecase struct {
	analytic interfaces.OrderAnalytic
}

func NewGetStatisticsByHoursUsecase(
	analytic interfaces.OrderAnalytic,
) *GetStatisticsByHoursUsecase {
	return &GetStatisticsByHoursUsecase{
		analytic: analytic,
	}
}

func (uc *GetStatisticsByHoursUsecase) Execute(ctx context.Context) ([]entity.HourStatistic, error) {
	return uc.analytic.GetStatisticsByHours(ctx)
}
