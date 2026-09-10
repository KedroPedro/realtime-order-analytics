package worker

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/KedroPedro/realtime-order-analytics/internal/application/usecases"
	"github.com/KedroPedro/realtime-order-analytics/internal/controller/worker"
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

	uc := usecases.NewReadAndPublishOrderUsecase(
		psql.NewOrdersOutboxRepo(),
		kfk.NewEventPublisher(),
	)

	workers := make([]*worker.OutboxWorker, numWorkers)

	for i := range numWorkers {
		workers = append(workers, worker.NewOutboxWorker(uc))
		go func(i int) {
			log.Debug().Msg(fmt.Sprintf("worker %d started", i))
			workers[i].Start()
		}(i)
	}

	stopCh := make(chan os.Signal, 2)
	signal.Notify(stopCh, os.Interrupt, os.Kill)

	<-stopCh
	log.Debug().Msg("shutting down")
	for _, w := range workers {
		w.Stop()
	}

}
