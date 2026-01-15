package persistence

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	appkafka "github.com/smarrog/task-board/core-service/internal/infrastructure/kafka"
	msg "github.com/smarrog/task-board/shared/messaging"
)

type OutboxWorker struct {
	txm          *TxManager
	repo         *OutboxRepo
	producer     *appkafka.Producer
	topic        string
	batchSize    int
	pollInterval time.Duration
	log          *zerolog.Logger
}

func NewOutboxWorker(
	txm *TxManager,
	repo *OutboxRepo,
	producer *appkafka.Producer,
	topic string,
	batchSize int,
	pollInterval time.Duration,
	log *zerolog.Logger,
) *OutboxWorker {
	return &OutboxWorker{
		txm:          txm,
		repo:         repo,
		producer:     producer,
		topic:        topic,
		batchSize:    batchSize,
		pollInterval: pollInterval,
		log:          log,
	}
}

func (w *OutboxWorker) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}

		if err := w.processOnce(ctx); err != nil {
			w.log.Err(err).Msg("outbox worker iteration failed")
		}
	}
}

func (w *OutboxWorker) processOnce(ctx context.Context) error {
	return w.txm.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		rows, err := w.repo.FetchUnpublishedForUpdate(ctx, tx, w.batchSize)
		if err != nil {
			return err
		}

		if len(rows) == 0 {
			return nil
		}

		ids := make([]uuid.UUID, 0, len(rows))
		for _, r := range rows {
			oMsg := msg.OutboxMessage{
				Id:            r.ID.String(),
				EventType:     r.EventType,
				AggregateType: r.AggregateType,
				AggregateId:   r.AggregateID.String(),
				CreatedAt:     r.CreatedAt.Time.UTC(),
				Payload:       json.RawMessage(r.Payload),
				Version:       1,
			}

			b, err := json.Marshal(oMsg)
			if err != nil {
				return err
			}

			if err := w.producer.Produce(w.topic, []byte(r.ID.String()), b); err != nil {
				return err
			}

			ids = append(ids, r.ID)
		}

		return w.repo.MarkPublished(ctx, tx, ids)
	})
}
