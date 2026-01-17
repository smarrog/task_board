package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/smarrog/task-board/shared/domain/shared"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"github.com/smarrog/task-board/core-service/internal/domain/board"
	"github.com/smarrog/task-board/core-service/internal/domain/column"
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

func (r *ColumnsRepo) ListByBoard(ctx context.Context, boardId board.Id) ([]*column.Column, error) {
	db := r.txm.DB(ctx)

	rows, err := db.Query(ctx, `
		SELECT id, position, created_at, updated_at
		FROM columns
		WHERE board_id = $1
		ORDER BY position ASC
	`, boardId.UUID())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]*column.Column, 0)
	for rows.Next() {
		var idRaw string
		var positionRaw int
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&idRaw, &positionRaw, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		id, err := column.IdFromString(idRaw)
		if err != nil {
			return nil, err
		}
		pos, err := column.NewPosition(positionRaw)
		if err != nil {
			return nil, err
		}
		out = append(out, column.Rehydrate(id, boardId, pos, createdAt, updatedAt))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (r *ColumnsRepo) ListByBoards(ctx context.Context, boardIds []board.Id) ([]*column.Column, error) {
	if len(boardIds) == 0 {
		return []*column.Column{}, nil
	}
	db := r.txm.DB(ctx)

	ids := make([]uuid.UUID, 0, len(boardIds))
	for _, bid := range boardIds {
		ids = append(ids, bid.UUID())
	}

	rows, err := db.Query(ctx, `
		SELECT id, board_id, position, created_at, updated_at
		FROM columns
		WHERE board_id = ANY($1)
		ORDER BY board_id ASC, position ASC
	`, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]*column.Column, 0)
	for rows.Next() {
		var (
			idRaw       string
			boardIdRaw  string
			positionRaw int
			createdAt   time.Time
			updatedAt   time.Time
		)
		if err := rows.Scan(&idRaw, &boardIdRaw, &positionRaw, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		id, err := column.IdFromString(idRaw)
		if err != nil {
			return nil, err
		}
		bid, err := board.IdFromString(boardIdRaw)
		if err != nil {
			return nil, err
		}
		pos, err := column.NewPosition(positionRaw)
		if err != nil {
			return nil, err
		}
		out = append(out, column.Rehydrate(id, bid, pos, createdAt, updatedAt))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
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

		events := []shared.DomainEvent{column.DeletedEvent{Id: id.String(), At: time.Now().UTC()}}
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
