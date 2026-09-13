package analytic

import (
	"time"
	"uuid"

	"github.com/KedroPedro/realtime-order-analytics/internal/domain/entity"
	"github.com/shopspring/decimal"
)

type DailyOrderTrend struct {
	Day     time.Time       `json:"day"`
	Orders  uint64          `json:"orders"`
	Revenue decimal.Decimal `json:"revenue"`
	Average decimal.Decimal `json:"average"`
}

func (m DailyOrderTrend) FromEntity(value entity.DailyOrderTrend) DailyOrderTrend {
	return DailyOrderTrend{
		Day:     value.Day,
		Orders:  value.Orders,
		Revenue: value.Revenue,
		Average: value.Average,
	}
}

type TopClient struct {
	ClientID     uuid.UUID       `json:"client_id"`
	TotalRevenue decimal.Decimal `json:"total_revenue"`
	Average      decimal.Decimal `json:"average"`
}

func (m TopClient) FromEntity(value entity.TopClient) TopClient {
	return TopClient{
		ClientID:     value.ClientId,
		TotalRevenue: value.TotalRevenue,
		Average:      value.Average,
	}
}

type HourStatistic struct {
	Hour         time.Time       `json:"hour"`
	TotalOrders  uint64          `json:"total_orders"`
	TotalRevenue decimal.Decimal `json:"total_revenue"`
	Average      decimal.Decimal `json:"average"`
}

func (m HourStatistic) FromEntity(value entity.HourStatistic) HourStatistic {
	return HourStatistic{
		Hour:         value.Hour,
		TotalOrders:  value.TotalOrders,
		TotalRevenue: value.TotalRevenue,
		Average:      value.Average,
	}
}
