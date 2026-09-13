package analytic

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	analyticsuc "github.com/KedroPedro/realtime-order-analytics/internal/application/usecases/analytics"
	"github.com/rs/zerolog/log"
)

type AnalyticsHandler struct {
	getDailyTrends      *analyticsuc.GetDailyOrderTrendsUsecase
	getStatisticsByHour *analyticsuc.GetStatisticsByHoursUsecase
	getTopClients       *analyticsuc.GetTopClientsUsecase
}

func NewAnalyticsHandler(
	getDailyTrends *analyticsuc.GetDailyOrderTrendsUsecase,
	getStatisticsByHour *analyticsuc.GetStatisticsByHoursUsecase,
	getTopClients *analyticsuc.GetTopClientsUsecase,
) *AnalyticsHandler {
	return &AnalyticsHandler{
		getDailyTrends:      getDailyTrends,
		getStatisticsByHour: getStatisticsByHour,
		getTopClients:       getTopClients,
	}
}

func (ah *AnalyticsHandler) DailyTrendsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	trends, err := ah.getDailyTrends.Execute(ctx)
	if err != nil {
		http.Error(w, "get statistic error", http.StatusInternalServerError)
		log.Err(err).Send()
		return
	}

	resp := make([]DailyOrderTrend, len(trends))
	conv := DailyOrderTrend{}
	for i := range trends {
		resp[i] = conv.FromEntity(trends[i])
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error().Err(err).Send()
	}
}

func (ah *AnalyticsHandler) StatisticsByHourHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	trends, err := ah.getStatisticsByHour.Execute(ctx)
	if err != nil {
		http.Error(w, "get statistic error", http.StatusInternalServerError)
		log.Err(err).Send()
		return
	}

	resp := make([]HourStatistic, len(trends))
	conv := HourStatistic{}
	for i := range trends {
		resp[i] = conv.FromEntity(trends[i])
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error().Err(err).Send()
	}
}

func (ah *AnalyticsHandler) TopClientsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	trends, err := ah.getTopClients.Execute(ctx)
	if err != nil {
		http.Error(w, "get statistic error", http.StatusInternalServerError)
		log.Err(err).Send()
		return
	}

	resp := make([]TopClient, len(trends))
	conv := TopClient{}
	for i := range trends {
		resp[i] = conv.FromEntity(trends[i])
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error().Err(err).Send()
	}
}
