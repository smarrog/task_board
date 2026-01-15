package kafka

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
	"github.com/smarrog/task-board/notification-service/internal/config"
)

type HandlerFunc func(ctx context.Context, msg *kafka.Message) error

type Consumer struct {
	readers []*kafka.Reader
	logger  *zerolog.Logger
	cfg     *config.Config
	handler HandlerFunc
}

func NewConsumer(cfg *config.Config, logger *zerolog.Logger, handlerFunc HandlerFunc) *Consumer {
	return &Consumer{logger: logger, cfg: cfg, handler: handlerFunc}
}

func (c *Consumer) Start(ctx context.Context) error {
	c.readers = make([]*kafka.Reader, 0, len(c.cfg.KafkaBrokers))
	for _, topic := range c.cfg.KafkaTopics {
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers: c.cfg.KafkaBrokers,
			Topic:   topic,
			GroupID: c.cfg.KafkaGroupId,
		})
		c.readers = append(c.readers, r)
	}

	errCh := make(chan error, len(c.readers))
	for _, r := range c.readers {
		go func() {
			defer func(r *kafka.Reader) {
				_ = r.Close()
			}(r)
			for {
				m, err := r.FetchMessage(ctx)
				if err != nil {
					if errors.Is(err, context.Canceled) {
						errCh <- nil
						return
					}
					c.logger.Err(err).Msg("failed to fetch kafka message")
					continue
				}

				hErr := c.handler(ctx, &m)
				if err := r.CommitMessages(ctx, m); err != nil {
					c.logger.Err(err).Msg("failed to commit kafka message")
				}
				if hErr != nil {
					c.logger.Err(hErr).Msg("handler returned error")
				}
			}
		}()
	}

	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		return err
	}
}
