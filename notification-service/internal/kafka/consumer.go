package kafka

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rs/zerolog"
	"github.com/smarrog/notification-service/internal/config"
)

type HandlerFunc func(ctx context.Context, msg *kafka.Message) error

type Consumer interface {
	Start(ctx context.Context)
}

type consumer struct {
	reader  *kafka.Consumer
	logger  *zerolog.Logger
	cfg     *config.Config
	handler HandlerFunc
}

func NewConsumer(cfg *config.Config, logger *zerolog.Logger, handlerFunc HandlerFunc) Consumer {
	kafkaCfg := &kafka.ConfigMap{
		"bootstrap.servers":       cfg.KafkaBrokers,
		"group.id":                cfg.KafkaGroupId,
		"auto.offset.reset":       "latest",
		"enable.auto.commit":      true,
		"auto.commit.interval.ms": 5000,
	}

	reader, err := kafka.NewConsumer(kafkaCfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create kafka consumer")
	}

	return &consumer{
		reader:  reader,
		logger:  logger,
		cfg:     cfg,
		handler: handlerFunc,
	}
}

func (c *consumer) Start(ctx context.Context) {
	if err := c.reader.SubscribeTopics(c.cfg.KafkaTopics, nil); err != nil {
		c.logger.Fatal().Err(err).Msg("Failed to subscribe to topics")
	}

	defer func(reader *kafka.Consumer) {
		err := reader.Close()
		if err != nil {
			c.logger.Err(err).Msg("Failed to close kafka consumer")
		}
	}(c.reader)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		m, err := c.reader.ReadMessage(-1)
		if err != nil {
			c.logger.Err(err).Msg("Failed to read message from kafka consumer")
			continue
		}

		if err := c.handler(ctx, m); err != nil {
			c.logger.Err(err).Msgf("Failed to handle message: %v", m)
		}
	}
}
