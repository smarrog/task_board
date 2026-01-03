package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"

	"github.com/smarrog/task-board/core-service/internal/domain/board"
	"github.com/smarrog/task-board/core-service/internal/domain/common"
)

type BoardsRepo struct {
	txm    *TxManager
	log    *zerolog.Logger
	outbox *OutboxRepo
}

func NewBoardsRepo(txm *TxManager, log *zerolog.Logger, outbox *OutboxRepo) *BoardsRepo {
	return &BoardsRepo{
		txm:    txm,
		log:    log,
		outbox: outbox,
	}
}

func (r *BoardsRepo) InTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.txm.InTx(ctx, func(ctx context.Context, _ pgx.Tx) error {
		return fn(ctx)
	})
}

func (r *BoardsRepo) Save(ctx context.Context, b *board.Board) error {
	return r.txm.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `
			INSERT INTO boards (id, owner_id, title, description, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (id) DO UPDATE
			SET owner_id    = EXCLUDED.owner_id,
				title       = EXCLUDED.title,
				description = EXCLUDED.description,
				updated_at  = EXCLUDED.updated_at
		`,
			b.Id().UUID(),
			b.OwnerId().UUID(),
			b.Title().String(),
			b.Description().String(),
			b.CreatedAt(),
			b.UpdatedAt(),
		)
		if err != nil {
			return err
		}

		events := b.PullEvents()
		return r.outbox.SaveEvents(ctx, events)
	})
}

func (r *BoardsRepo) Get(ctx context.Context, id board.Id) (*board.Board, error) {
	db := r.txm.DB(ctx)

	var ownerIdRaw uuid.UUID
	var titleRaw, descRaw string
	var createdAt, updatedAt time.Time

	err := db.QueryRow(ctx, `
        SELECT owner_id, title, description, created_at, updated_at
        FROM boards
        WHERE id = $1
    `, id.UUID()).Scan(&ownerIdRaw, &titleRaw, &descRaw, &createdAt, &updatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, board.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	ownerId, err := common.UserIdFromUUID(ownerIdRaw)
	if err != nil {
		return nil, err
	}
	title, err := board.NewTitle(titleRaw)
	if err != nil {
		return nil, err
	}
	desc, err := board.NewDescription(descRaw)
	if err != nil {
		return nil, err
	}

	return board.Rehydrate(id, ownerId, title, desc, createdAt, updatedAt), nil
}

func (r *BoardsRepo) ListByOwner(ctx context.Context, ownerId common.UserId) ([]*board.Board, error) {
	db := r.txm.DB(ctx)

	rows, err := db.Query(ctx, `
        SELECT id, title, description, created_at, updated_at
        FROM boards
        WHERE owner_id = $1
        ORDER BY created_at DESC, id
    `, ownerId.UUID())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]*board.Board, 0)
	for rows.Next() {
		var (
			idRaw     uuid.UUID
			titleRaw  string
			descRaw   string
			createdAt time.Time
			updatedAt time.Time
		)
		if err := rows.Scan(&idRaw, &titleRaw, &descRaw, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		id, err := board.IdFromUUID(idRaw)
		if err != nil {
			return nil, err
		}
		t, err := board.NewTitle(titleRaw)
		if err != nil {
			return nil, err
		}
		d, err := board.NewDescription(descRaw)
		if err != nil {
			return nil, err
		}
		out = append(out, board.Rehydrate(id, ownerId, t, d, createdAt, updatedAt))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *BoardsRepo) Delete(ctx context.Context, id board.Id) error {
	return r.txm.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		ct, err := tx.Exec(ctx, `DELETE FROM boards WHERE id = $1`, id.UUID())
		if err != nil {
			return err
		}
		if ct.RowsAffected() == 0 {
			return board.ErrNotFound
		}

		// TODO send event to outbox

		return nil
	})
}
