package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/smarrog/task-board/core-service/internal/domain/board"
	"github.com/smarrog/task-board/core-service/internal/domain/common"
)

type BoardsRepo struct {
	pg     *pgxpool.Pool
	log    *zerolog.Logger
	outbox *OutboxRepo
}

func NewBoardsRepo(pg *pgxpool.Pool, log *zerolog.Logger, outbox *OutboxRepo) *BoardsRepo {
	return &BoardsRepo{
		pg:     pg,
		log:    log,
		outbox: outbox,
	}
}

func (r *BoardsRepo) Save(ctx context.Context, b *board.Board) error {
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

func (r *BoardsRepo) Get(ctx context.Context, id board.Id) (*board.Board, error) {
	var ownerIdRaw uuid.UUID
	var titleRaw, descRaw string
	var createdAt, updatedAt time.Time

	err := r.pg.QueryRow(ctx, `
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
	title, err := common.NewTitle(titleRaw)
	if err != nil {
		return nil, err
	}
	desc, err := common.NewDescription(descRaw)
	if err != nil {
		return nil, err
	}

	return board.Rehydrate(id, ownerId, title, desc, createdAt, updatedAt), nil
}

func (r *BoardsRepo) ListByOwner(ctx context.Context, ownerId common.UserId) ([]*board.Board, error) {
	rows, err := r.pg.Query(ctx, `
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
		t, err := common.NewTitle(titleRaw)
		if err != nil {
			return nil, err
		}
		d, err := common.NewDescription(descRaw)
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

	ct, err := tx.Exec(ctx, `DELETE FROM boards WHERE id = $1`, id.UUID())
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return board.ErrNotFound
	}

	// TODO send event to outbox

	return nil
}
