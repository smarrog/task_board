package kafka

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
	"github.com/smarrog/task-board/core-service/internal/config"
)

type Producer struct {
	w   *kafka.Writer
	log *zerolog.Logger
	cfg *config.Config
}

func NewProducer(cfg *config.Config, log *zerolog.Logger) (*Producer, error) {
	w := &kafka.Writer{
		Addr:         kafka.TCP(cfg.KafkaBrokers...),
		Balancer:     &kafka.Hash{},
		RequiredAcks: kafka.RequiredAcks(cfg.KafkaAcks),
	}

	return &Producer{w: w, log: log, cfg: cfg}, nil
}

func (p *Producer) Produce(topic string, key, value []byte) error {
	return p.w.WriteMessages(context.Background(), kafka.Message{
		Topic: topic,
		Key:   key,
		Value: value,
	})
}

func (p *Producer) Close() {
	_ = p.w.Close()
}
