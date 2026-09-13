package analytic

import (
	"net/http"

	analyticsuc "github.com/KedroPedro/realtime-order-analytics/internal/application/usecases/analytics"
)

func SetupAnalyticsOrders(
	mux *http.ServeMux,
	getDailyTrends *analyticsuc.GetDailyOrderTrendsUsecase,
	getStatisticsByHour *analyticsuc.GetStatisticsByHoursUsecase,
	getTopClients *analyticsuc.GetTopClientsUsecase,
) {
	handler := NewAnalyticsHandler(
		getDailyTrends,
		getStatisticsByHour,
		getTopClients,
	)

	mux.HandleFunc("GET /analytic/get-daily-trends", handler.DailyTrendsHandler)
	mux.HandleFunc("GET /analytic/get-stat-by-hour", handler.StatisticsByHourHandler)
	mux.HandleFunc("GET /analytic/get-top-clients", handler.TopClientsHandler)
}
