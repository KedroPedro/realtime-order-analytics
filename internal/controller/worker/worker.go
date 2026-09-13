package worker

import (
	"context"
	"errors"
	"fmt"
	"time"

	workeruc "github.com/KedroPedro/realtime-order-analytics/internal/application/usecases/worker"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type OutboxWorker struct {
	uc     *workeruc.ReadAndPublishOrderUsecase
	ticker time.Ticker
	stopCh chan struct{}
	name   string
}

const (
	tickerDuration = time.Millisecond * 500
	ctxTimeout     = time.Second * 10
)

func NewOutboxWorker(
	uc *workeruc.ReadAndPublishOrderUsecase,
	name string,
) *OutboxWorker {
	return &OutboxWorker{
		uc:     uc,
		stopCh: make(chan struct{}),
		ticker: *time.NewTicker(tickerDuration),
		name:   name,
	}
}

func (ow *OutboxWorker) Start() {
	defer ow.ticker.Stop()
	for {
		select {
		case <-ow.stopCh:
			return
		case <-ow.ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)

			if err := ow.uc.Execute(ctx, ow.name); err != nil {
				if !errors.Is(pgx.ErrNoRows, err) {
					log.Err(err).Msg(fmt.Sprintf("worker name: %q", ow.name))
				}
			}

			cancel()
		}

	}
}

func (ow *OutboxWorker) Stop() {
	ow.stopCh <- struct{}{}
	close(ow.stopCh)
}
