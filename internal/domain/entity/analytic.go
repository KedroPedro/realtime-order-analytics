package entity

import (
	"time"
	"uuid"

	"github.com/shopspring/decimal"
)

type DailyOrderTrend struct {
	Day     time.Time
	Orders  uint64
	Revenue decimal.Decimal
	Average decimal.Decimal
}

type TopClient struct {
	ClientId     uuid.UUID
	TotalRevenue decimal.Decimal
	Average      decimal.Decimal
}

type HourStatistic struct {
	Hour         time.Time
	TotalOrders  uint64
	TotalRevenue decimal.Decimal
	Average      decimal.Decimal
}
