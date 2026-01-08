package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"github.com/smarrog/task-board/core-service/internal/domain/board"
	"github.com/smarrog/task-board/core-service/internal/domain/column"
	"github.com/smarrog/task-board/core-service/internal/domain/common"
)

type ColumnsRepo struct {
	txm    *TxManager
	log    *zerolog.Logger
	outbox *OutboxRepo
}

func NewColumnsRepo(txm *TxManager, log *zerolog.Logger, outbox *OutboxRepo) *ColumnsRepo {
	return &ColumnsRepo{
		txm:    txm,
		log:    log,
		outbox: outbox,
	}
}

func (r *ColumnsRepo) InTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.txm.InTx(ctx, func(ctx context.Context, _ pgx.Tx) error {
		return fn(ctx)
	})
}

func (r *ColumnsRepo) Save(ctx context.Context, c *column.Column) error {
	return r.txm.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `
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
		return r.outbox.SaveEvents(ctx, events)
	})
}

func (r *ColumnsRepo) Get(ctx context.Context, id column.Id) (*column.Column, error) {
	db := r.txm.DB(ctx)

	var boardIdRaw string
	var positionRaw int
	var createdAt, updatedAt time.Time

	err := db.QueryRow(ctx, `
		SELECT board_id, position, created_at, updated_at
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
	return r.txm.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		ct, err := tx.Exec(ctx, `DELETE FROM columns WHERE id = $1`, id.UUID())
		if err != nil {
			return err
		}
		if ct.RowsAffected() == 0 {
			return column.ErrNotFound
		}

		events := []common.DomainEvent{column.DeletedEvent{Id: id.String(), At: time.Now().UTC()}}
		if err := r.outbox.SaveEvents(ctx, events); err != nil {
			return err
		}


		return nil
	})
}

func (r *ColumnsRepo) LockBoardColumns(ctx context.Context, boardId board.Id) error {
	return r.txm.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `SELECT id FROM columns WHERE board_id=$1 FOR UPDATE`, boardId.UUID())
		return err
	})
}

func (r *ColumnsRepo) CountInBoard(ctx context.Context, boardId board.Id) (int, error) {
	var n int
	err := r.txm.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		return tx.QueryRow(ctx, `SELECT COUNT(*) FROM columns WHERE board_id=$1`, boardId.UUID()).Scan(&n)
	})
	return n, err
}

func (r *ColumnsRepo) ShiftPositions(ctx context.Context, boardId board.Id, fromIncl, toIncl int, delta int) error {
	if fromIncl > toIncl {
		return nil
	}
	return r.txm.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `
			UPDATE columns
			SET position = position + $4
			WHERE board_id=$1 AND position BETWEEN $2 AND $3
		`, boardId.UUID(), fromIncl, toIncl, delta)
		return err
	})
}
