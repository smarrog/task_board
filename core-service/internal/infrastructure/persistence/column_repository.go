package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/smarrog/task-board/core-service/internal/domain/board"
	"github.com/smarrog/task-board/core-service/internal/domain/column"
)

type ColumnsRepo struct {
	pg     *pgxpool.Pool
	log    *zerolog.Logger
	outbox *OutboxRepo
}

func NewColumnsRepo(pg *pgxpool.Pool, log *zerolog.Logger, outbox *OutboxRepo) *ColumnsRepo {
	return &ColumnsRepo{
		pg:     pg,
		log:    log,
		outbox: outbox,
	}
}

func (r *ColumnsRepo) Save(ctx context.Context, c *column.Column) error {
	tx, txErr := r.pg.Begin(ctx)
	if txErr != nil {
		return txErr
	}

	var err error
	defer func() {
		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				r.log.Error().Err(errRollback).Msg("failed to rollback transaction")
			}
		}
	}()

	_, err = tx.Exec(ctx, `
        INSERT INTO columns (id, board_id, position, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (id) DO UPDATE
        SET board_id    = EXCLUDED.board_id,
            position    = EXCLUDED.position,
            updated_at  = EXCLUDED.updated_at
    `,
		c.Id().UUID(),
		c.BoardId().UUID(),
		c.Position(),
		c.CreatedAt(),
		c.UpdatedAt(),
	)
	if err != nil {
		return err
	}

	events := c.PullEvents()
	if len(events) != 0 {
		err = r.outbox.SaveEvents(ctx, tx, events)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *ColumnsRepo) Get(ctx context.Context, id column.Id) (*column.Column, error) {
	var boardIdRaw string
	var positionRaw int
	var createdAt, updatedAt time.Time

	err := r.pg.QueryRow(ctx, `
		SELECT board_id, position, createdAt, updatedAt
		FROM columns
		WHERE id = $1
	`, id.UUID()).Scan(&boardIdRaw, &positionRaw, &createdAt, &updatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, column.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	boardId, err := board.IdFromString(boardIdRaw)
	if err != nil {
		return nil, err
	}
	position, err := column.NewPosition(positionRaw)
	if err != nil {
		return nil, err
	}

	return column.Rehydrate(id, boardId, position, createdAt, updatedAt), nil
}

func (r *ColumnsRepo) Delete(ctx context.Context, id column.Id) error {
	tx, txErr := r.pg.Begin(ctx)
	if txErr != nil {
		return txErr
	}

	var err error
	defer func() {
		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				r.log.Error().Err(errRollback).Msg("failed to rollback transaction")
			}
		}
	}()

	ct, err := tx.Exec(ctx, `DELETE FROM columns WHERE id = $1`, id.UUID())
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return column.ErrNotFound
	}

	// TODO send event to outbox

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
