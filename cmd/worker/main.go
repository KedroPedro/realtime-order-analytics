package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	workeruc "github.com/KedroPedro/realtime-order-analytics/internal/application/usecases/worker"
	"github.com/KedroPedro/realtime-order-analytics/internal/controller/worker"
	"github.com/KedroPedro/realtime-order-analytics/internal/controller/worker/cleaner"
	"github.com/KedroPedro/realtime-order-analytics/internal/infrastructure/kafka"
	"github.com/KedroPedro/realtime-order-analytics/internal/infrastructure/postgresql"
	"github.com/rs/zerolog/log"
)

const (
	numWorkers = 5
)

func main() {
	psql, err := postgresql.NewPostgres()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	kfk, err := kafka.NewKafka()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	if err := kfk.SetupOrderOutboxTopic(); err != nil {
		log.Fatal().Err(err).Send()
	}

	workUsecase := workeruc.NewReadAndPublishOrderUsecase(
		psql.NewOrdersOutboxRepo(),
		kfk.NewEventPublisher(),
	)

	workers := make([]*worker.OutboxWorker, 0, numWorkers)

	for i := range numWorkers {
		workers = append(workers, worker.NewOutboxWorker(workUsecase, fmt.Sprintf("%d", i)))
		go func(i int) {
			log.Debug().Msg(fmt.Sprintf("worker %d started", i))
			workers[i].Start()
		}(i)
	}

	cleanUsecase := workeruc.NewCleanOutboxUsecasee(
		psql.NewCleanOutboxRepo(),
	)

	cleaner := cleaner.NewOutboxCleaner(cleanUsecase.Execute)

	go func() {
		cleaner.Start()
	}()

	stopCh := make(chan os.Signal, 2)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)

	<-stopCh
	log.Debug().Msg("shutting down")
	for _, w := range workers {
		w.Stop()
	}
	cleaner.Stop()

}
