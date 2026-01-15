package kafka

import (
	"errors"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rs/zerolog"
	"github.com/smarrog/task-board/core-service/internal/config"
)

type Producer struct {
	p   *kafka.Producer
	log *zerolog.Logger
	cfg *config.Config
}

func NewProducer(cfg *config.Config, log *zerolog.Logger) (*Producer, error) {
	kafkaCfg := &kafka.ConfigMap{
		"bootstrap.servers": cfg.KafkaBrokers,
		"acks":              cfg.KafkaAcks,
	}

	p, err := kafka.NewProducer(kafkaCfg)
	if err != nil {
		return nil, err
	}

	pr := &Producer{p: p, log: log, cfg: cfg}
	go pr.logDeliveryEvents()
	return pr, nil
}

func (p *Producer) logDeliveryEvents() {
	for e := range p.p.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				p.log.Err(ev.TopicPartition.Error).Msg("kafka delivery failed")
			}
		}
	}
}

func (p *Producer) Produce(topic string, key, value []byte) error {
	if topic == "" {
		return errors.New("kafka topic is empty")
	}

	// Produce is async; we rely on Flush at commit boundaries for bounded latency.
	return p.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
		Value:          value,
	}, nil)
}

func (p *Producer) Flush(timeoutMs int) int {
	return p.p.Flush(timeoutMs)
}

func (p *Producer) Close() {
	p.p.Close()
}
