package orderanalytic

import (
	"context"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/KedroPedro/realtime-order-analytics/internal/domain/entity"
)

type OrderRepository struct {
	conn driver.Conn
}

func NewOrderRepository(conn driver.Conn) *OrderRepository {
	return &OrderRepository{
		conn: conn,
	}
}

func (or *OrderRepository) GetDailyOrderTrends(ctx context.Context) ([]entity.DailyOrderTrend, error) {
	query :=
		`
			SELECT
				toDate(created_at) AS day,
				count() AS orders_count,
				sum(total) AS revenue,
				avg(total) AS average_order
			FROM analytics.order_events 
			WHERE created_at >= now()- INTERVAL 30 DAY
			GROUP BY day
			ORDER BY day ASC;
		`

	rows, err := or.conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dailyOrderTrends := make([]entity.DailyOrderTrend, 0)

	temp := DailyOrderTrend{}

	for rows.Next() {
		if err := rows.ScanStruct(&temp); err != nil {
			return nil, err
		}

		dailyOrderTrends = append(dailyOrderTrends, temp.ToEntity())
	}

	return dailyOrderTrends, nil
}

func (or *OrderRepository) GetTopClients(ctx context.Context) ([]entity.TopClient, error) {
	query :=
		`
			SELECT
				owner_id,
				sum(total) AS total_revenue,
				avg(total) AS avg_order
			FROM analytics.order_events 
			WHERE event_type = 'order_created'
				AND created_at >= now() - INTERVAL 30 DAY
			GROUP BY owner_id
			ORDER BY total_revenue DESC
			LIMIT 10;
		`

	rows, err := or.conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	topClients := make([]entity.TopClient, 0)

	temp := TopClient{}
	for rows.Next() {
		if err := rows.ScanStruct(&temp); err != nil {
			return nil, err
		}

		topClients = append(topClients, temp.ToEntity())
	}

	return topClients, nil
}

func (or *OrderRepository) GetStatisticsByHours(ctx context.Context) ([]entity.HourStatistic, error) {
	query :=
		`
			SELECT
				toHour(create_at) AS hour,
				count() AS total_orders,
				sum(total) AS total_revenue,
				avg(total) AS avg_order
			FROM analytics.order_events 
			WHERE created_at >= now() - INTERVAL 7 DAY
			GROUP BY hour
			ORDER BY hour ASC;
		`

	rows, err := or.conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hourStatistics := make([]entity.HourStatistic, 0)

	temp := HourStatistic{}
	for rows.Next() {
		if err := rows.ScanStruct(&temp); err != nil {
			return nil, err
		}

		hourStatistics = append(hourStatistics, temp.ToEntity())
	}

	return hourStatistics, nil

}
