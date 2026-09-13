package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	orderuc "github.com/KedroPedro/realtime-order-analytics/internal/application/usecases/order"
	"github.com/KedroPedro/realtime-order-analytics/internal/controller/http/order"
	"github.com/KedroPedro/realtime-order-analytics/internal/controller/http/server"
	"github.com/KedroPedro/realtime-order-analytics/internal/infrastructure/postgresql"
	"github.com/rs/zerolog/log"
)

func main() {
	notifySig := make(chan os.Signal, 2)
	signal.Notify(notifySig, os.Interrupt, syscall.SIGTERM)

	pg, err := postgresql.NewPostgres()
	if err != nil {
		panic(err)
	}

	createOrderUC := orderuc.NewCreateOrderUsecase(pg.NewOrdersRepo())

	mux := http.NewServeMux()

	order.SetupOrdersRoute(mux, createOrderUC)

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
