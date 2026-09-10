package cleaner

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

type Cleaner struct {
	ticker *time.Ticker
	clean  func(context.Context) error
	stopCh chan struct{}
}

const (
	cleanDuration = time.Second * 60
)

func NewOutboxCleaner(
	cleanFunc func(context.Context) error,
) *Cleaner {
	return &Cleaner{
		ticker: time.NewTicker(cleanDuration),
		stopCh: make(chan struct{}),
		clean:  cleanFunc,
	}
}

func (c *Cleaner) Start() {
	for {
		select {
		case <-c.ticker.C:
			cleanCtx, cancel := context.WithTimeout(context.Background(), time.Second*30)

			if err := c.clean(cleanCtx); err != nil {
				log.Err(err).Send()
			}

			cancel()
		case <-c.stopCh:
			return
		}
	}
}

func (c *Cleaner) Stop() {
	c.stopCh <- struct{}{}
}
