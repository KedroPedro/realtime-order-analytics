package kafka

import (
	"fmt"
	"net"
	"os"
	"strconv"
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

func (k *Kafka) SetupOrderOutboxTopic() error {
	conn, err := kafka.Dial("tcp", k.brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	br, err := conn.Controller()
	if err != nil {
		return err
	}

	controllerAddr := net.JoinHostPort(
		br.Host,
		strconv.Itoa(br.Port),
	)

	controller, err := kafka.Dial("tcp", controllerAddr)
	if err != nil {
		return err
	}
	defer controller.Close()

	controller.CreateTopics(
		kafka.TopicConfig{
			Topic:             orderOutboxTopic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	)

	return nil
}

func (k *Kafka) newWriter(topic string) *kafka.Writer {
	return kafka.NewWriter(
		kafka.WriterConfig{
			Brokers:      k.brokers,
			Topic:        topic,
			BatchSize:    100,
			MaxAttempts:  10,
			Balancer:     &kafka.RoundRobin{},
			RequiredAcks: int(kafka.RequireOne),
		},
	)
}

func getBrokers() ([]string, error) {
	brokersStr := os.Getenv(envKafkaBrokers)
	brokersStr = strings.TrimSpace(brokersStr)
	if brokersStr == "" {
		return nil, fmt.Errorf("environment variable %q is not set", envKafkaBrokers)
	}

	return strings.Split(brokersStr, ";"), nil
}
