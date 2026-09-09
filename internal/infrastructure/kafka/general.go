package kafka

import (
	"fmt"
	"os"
	"strings"

	"github.com/KedroPedro/realtime-order-analytics/internal/domain/interfaces"
	"github.com/segmentio/kafka-go"
)

type Kafka struct {
	brokers []string
}

const (
	envKafkaBrokers = "KAFKA_BROKERS"
)

func NewKafka() (*Kafka, error) {
	brokers, err := getBrokers()
	if err != nil {
		return nil, err
	}

	return &Kafka{
		brokers: brokers,
	}, nil
}

const (
	orderOutboxTopic = "order_outbox"
)

func (k *Kafka) NewEventPublisher() interfaces.EventPublisher {
	return NewEventPublisher(k.newWriter(orderOutboxTopic))
}

func (k *Kafka) newWriter(topic string) *kafka.Writer {
	return kafka.NewWriter(
		kafka.WriterConfig{
			Brokers:     k.brokers,
			Topic:       topic,
			BatchSize:   100,
			MaxAttempts: 10,
			Balancer:    &kafka.Hash{},
		},
	)
}

func getBrokers() ([]string, error) {
	brokersStr := os.Getenv(envKafkaBrokers)
	if brokersStr == "" {
		return nil, fmt.Errorf("environment variable %q is not set", envKafkaBrokers)
	}
	brokersStr = strings.TrimSpace(brokersStr)

	return strings.Split(brokersStr, ";"), nil
}
