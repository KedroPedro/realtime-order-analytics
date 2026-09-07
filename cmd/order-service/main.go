package main

import (
	"net/http"
	"os"
	"os/signal"

	"github.com/KedroPedro/realtime-order-analytics/internal/application/usecases"
	"github.com/KedroPedro/realtime-order-analytics/internal/controller/http/order"
	"github.com/KedroPedro/realtime-order-analytics/internal/infrastructure/postgresql"
	"github.com/rs/zerolog/log"
)

func main() {
	notifySig := make(chan os.Signal)
	signal.Notify(notifySig, os.Interrupt, os.Kill)

	pg, err := postgresql.NewPostgres()
	if err != nil {
		panic(err)
	}

	createOrderUC := usecases.NewCreateOrderUsecase(pg.NewOrdersRepo())

	mux := http.NewServeMux()

	order.SetupOrdersRoute(mux, createOrderUC)

	srv := NewServer(mux)

	go func() {
		<-notifySig
		if err := srv.Stop(); err != nil {
			log.Fatal().Err(err).Send()
		}
	}()

	if err := srv.Start(); err != nil {
		log.Fatal().Err(err).Send()
	}
}
