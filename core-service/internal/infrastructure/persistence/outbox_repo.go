package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
	"github.com/smarrog/task-board/core-service/internal/domain/board"
	"github.com/smarrog/task-board/core-service/internal/domain/column"
	"github.com/smarrog/task-board/core-service/internal/domain/common"
	"github.com/smarrog/task-board/core-service/internal/domain/task"
)

type outboxEventRow struct {
	ID            uuid.UUID
	EventType     string
	AggregateType string
	AggregateID   uuid.UUID
	Payload       []byte
	CreatedAt     pgtype.Timestamptz
}

type OutboxRepo struct {
	txm *TxManager
	log *zerolog.Logger
}

func NewOutboxRepo(txm *TxManager, log *zerolog.Logger) *OutboxRepo {
	return &OutboxRepo{txm: txm, log: log}
}

func (r *OutboxRepo) SaveEvents(ctx context.Context, events []common.DomainEvent) error {
	if len(events) == 0 {
		return nil
	}

	return r.txm.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		batch := &pgx.Batch{}

		for _, ev := range events {
			aggType, aggID, err := r.aggregateInfoFromEvent(ev)
			if err != nil {
				return err
			}

			payload, err := json.Marshal(ev)
			if err != nil {
				return fmt.Errorf("marshal domain event: %w", err)
			}

			outboxID := uuid.New()
			batch.Queue(
				`INSERT INTO outbox_events (id, event_type, aggregate_type, aggregate_id, payload) VALUES ($1, $2, $3, $4, $5)`,
				outboxID,
				ev.Name(),
				aggType,
				aggID,
				payload,
			)
		}

		br := tx.SendBatch(ctx, batch)
		defer func(br pgx.BatchResults) {
			err := br.Close()
			if err != nil {
				r.log.Error().Err(err).Msg("error closing batch")
			}
		}(br)

		for range events {
			if _, err := br.Exec(); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *OutboxRepo) FetchUnpublishedForUpdate(ctx context.Context, tx pgx.Tx, limit int) ([]outboxEventRow, error) {
	if limit <= 0 {
		limit = 50
	}

	rows, err := tx.Query(ctx, `
        SELECT id, event_type, aggregate_type, aggregate_id, payload, created_at
        FROM outbox_events
        WHERE published_at IS NULL
        ORDER BY created_at
        LIMIT $1
        FOR UPDATE SKIP LOCKED
    `, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]outboxEventRow, 0, limit)
	for rows.Next() {
		var rrow outboxEventRow
		if err := rows.Scan(&rrow.ID, &rrow.EventType, &rrow.AggregateType, &rrow.AggregateID, &rrow.Payload, &rrow.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, rrow)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *OutboxRepo) MarkPublished(ctx context.Context, tx pgx.Tx, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}

	_, err := tx.Exec(ctx, `UPDATE outbox_events SET published_at = now() WHERE id = ANY($1)`, ids)
	return err
}

func (r *OutboxRepo) aggregateInfoFromEvent(ev common.DomainEvent) (string, uuid.UUID, error) {
	switch e := ev.(type) {
	case task.CreatedEvent:
		id, err := uuid.Parse(e.Id)
		return "task", id, r.wrapAggregateIdErr("task", e.Id, err)
	case task.UpdatedEvent:
		id, err := uuid.Parse(e.Id)
		return "task", id, r.wrapAggregateIdErr("task", e.Id, err)
	case task.MoveEvent:
		id, err := uuid.Parse(e.Id)
		return "task", id, r.wrapAggregateIdErr("task", e.Id, err)
	case task.DeletedEvent:
		id, err := uuid.Parse(e.Id)
		return "task", id, r.wrapAggregateIdErr("task", e.Id, err)

	case column.CreatedEvent:
		id, err := uuid.Parse(e.Id)
		return "column", id, r.wrapAggregateIdErr("column", e.Id, err)
	case column.MoveEvent:
		id, err := uuid.Parse(e.Id)
		return "column", id, r.wrapAggregateIdErr("column", e.Id, err)
	case column.DeletedEvent:
		id, err := uuid.Parse(e.Id)
		return "column", id, r.wrapAggregateIdErr("column", e.Id, err)

	case board.CreatedEvent:
		id, err := uuid.Parse(e.Id)
		return "board", id, r.wrapAggregateIdErr("board", e.Id, err)
	case board.UpdatedEvent:
		id, err := uuid.Parse(e.Id)
		return "board", id, r.wrapAggregateIdErr("board", e.Id, err)
	case board.DeletedEvent:
		id, err := uuid.Parse(e.Id)
		return "board", id, r.wrapAggregateIdErr("board", e.Id, err)
	default:
		return "", uuid.Nil, errors.New("unknown domain event type")
	}
}

func (r *OutboxRepo) wrapAggregateIdErr(aggregateType, raw string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("invalid aggregate_id for %s: %q: %w", aggregateType, raw, err)
}
