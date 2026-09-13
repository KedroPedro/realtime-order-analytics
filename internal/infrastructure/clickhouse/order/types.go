package orderanalytic

import (
	"time"
	"uuid"

	"github.com/KedroPedro/realtime-order-analytics/internal/domain/entity"
	"github.com/shopspring/decimal"
)

type DailyOrderTrend struct {
	Day     time.Time       `ch:"day"`
	Orders  uint64          `ch:"orders_count"`
	Revenue decimal.Decimal `ch:"revenue"`
	Average decimal.Decimal `ch:"average_order"`
}

func (r DailyOrderTrend) ToEntity() entity.DailyOrderTrend {
	return entity.DailyOrderTrend{
		Day:     r.Day,
		Orders:  r.Orders,
		Revenue: r.Revenue,
		Average: r.Average,
	}
}

type TopClient struct {
	ClientId     uuid.UUID       `ch:"owner_id"`
	TotalRevenue decimal.Decimal `ch:"total_revenue"`
	Average      decimal.Decimal `ch:"avg_order"`
}

func (r TopClient) ToEntity() entity.TopClient {
	return entity.TopClient{
		ClientId:     r.ClientId,
		TotalRevenue: r.TotalRevenue,
		Average:      r.Average,
	}
}

type HourStatistic struct {
	Hour         time.Time       `ch:"hour"`
	TotalOrders  uint64          `ch:"total_orders"`
	TotalRevenue decimal.Decimal `ch:"total_revenue"`
	Average      decimal.Decimal `ch:"avg_order"`
}

func (r HourStatistic) ToEntity() entity.HourStatistic {
	return entity.HourStatistic{
		Hour:         r.Hour,
		TotalOrders:  r.TotalOrders,
		TotalRevenue: r.TotalRevenue,
		Average:      r.Average,
	}
}
