package kafka

import (
	"context"
	"encoding/json"

	"github.com/KedroPedro/realtime-order-analytics/internal/domain/entity"
	"github.com/segmentio/kafka-go"
)

type EventPublisher struct {
	writer *kafka.Writer
}

func NewEventPublisher(writer *kafka.Writer) *EventPublisher {
	return &EventPublisher{
		writer: writer,
	}
}

func (ep *EventPublisher) PublishBatch(ctx context.Context, events []entity.OrderOutbox) error {
	msgs := make([]kafka.Message, len(events))
	bEvent := make([]byte, 0, 1<<7)
	var err error

	for _, event := range events {
		bEvent, err = json.Marshal(map[string]any{
			"id":         event.Id,
			"event_type": event.EventType,
			"created_at": event.CreatedAt,
			"payload":    event.Payload,
		})
		if err != nil {
			return err
		}

		msgs = append(
			msgs,
			kafka.Message{Value: bEvent},
		)
	}

	if err := ep.writer.WriteMessages(ctx, msgs...); err != nil {
		return err
	}

	return nil
}
