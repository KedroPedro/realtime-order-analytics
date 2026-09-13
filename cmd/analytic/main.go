package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	analyticsuc "github.com/KedroPedro/realtime-order-analytics/internal/application/usecases/analytics"
	"github.com/KedroPedro/realtime-order-analytics/internal/controller/http/analytic"
	"github.com/KedroPedro/realtime-order-analytics/internal/controller/http/server"
	"github.com/KedroPedro/realtime-order-analytics/internal/infrastructure/clickhouse"
	"github.com/rs/zerolog/log"
)

func main() {
	notifySig := make(chan os.Signal, 2)
	signal.Notify(notifySig, os.Interrupt, syscall.SIGTERM)

	ch, err := clickhouse.NewClickhouse()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	orderAnalytic := ch.NewOrderAnalytic()

	mux := http.NewServeMux()

	analytic.SetupAnalyticsOrders(
		mux,
		analyticsuc.NewGetDailyOrderTrendsUsecase(orderAnalytic),
		analyticsuc.NewGetStatisticsByHoursUsecase(orderAnalytic),
		analyticsuc.NewGetTopClientsUsecase(orderAnalytic),
	)

	srv := server.New(mux)

	go func() {
		<-notifySig
		if err := srv.Stop(); err != nil {
			log.Fatal().Err(err).Send()
		}
	}()

	log.Debug().Msg("server started")
	if err := srv.Start(); err != nil {
		log.Fatal().Err(err).Send()
	}
}
