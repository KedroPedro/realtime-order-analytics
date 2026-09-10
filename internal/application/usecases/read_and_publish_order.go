package usecases

import (
	"context"

	"github.com/KedroPedro/realtime-order-analytics/internal/domain/interfaces"
)

type ReadAndPublishOrderUsecase struct {
	outbox    interfaces.OrdersOutbox
	publisher interfaces.EventPublisher
}

func NewReadAndPublishOrderUsecase(
	outbox interfaces.OrdersOutbox,
	publisher interfaces.EventPublisher,
) *ReadAndPublishOrderUsecase {
	return &ReadAndPublishOrderUsecase{
		outbox:    outbox,
		publisher: publisher,
	}
}

func (uc *ReadAndPublishOrderUsecase) Execute(ctx context.Context) error {
	events, err := uc.outbox.GetEventsBatch(ctx)
	if err != nil {
		return err
	}

	if err := uc.publisher.PublishBatch(ctx, events); err != nil {
		return err
	}

	if err := uc.outbox.PublishEvents(ctx, events); err != nil {
		return err
	}

	return nil
}
