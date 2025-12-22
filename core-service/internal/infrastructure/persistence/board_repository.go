package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/smarrog/task-board/core-service/internal/domain/board"
)

type BoardRepository struct {
	pg *pgxpool.Pool
}

func NewBoardRepository(pg *pgxpool.Pool) *BoardRepository {
	return &BoardRepository{pg: pg}
}

func (r *BoardRepository) Save(ctx context.Context, b *board.Board) error {
	_, err := r.pg.Exec(ctx, `
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

	return nil
}

func (r *BoardRepository) Get(ctx context.Context, id board.BoardId) (*board.Board, error) {
	var (
		ownerID   uuid.UUID
		titleRaw  string
		descRaw   string
		createdAt time.Time
		updatedAt time.Time
	)

	err := r.pg.QueryRow(ctx, `
        SELECT owner_id, title, description, created_at, updated_at
        FROM boards
        WHERE id = $1
    `, id.UUID()).Scan(&ownerID, &titleRaw, &descRaw, &createdAt, &updatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, board.ErrBoardNotFound
	}
	if err != nil {
		return nil, err
	}

	o, err := board.UserIdFromUUID(ownerID)
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

	return board.RehydrateBoard(id, o, t, d, createdAt, updatedAt), nil
}

func (r *BoardRepository) ListByOwner(ctx context.Context, ownerID board.UserId) ([]*board.Board, error) {
	rows, err := r.pg.Query(ctx, `
        SELECT id, title, description, created_at, updated_at
        FROM boards
        WHERE owner_id = $1
        ORDER BY created_at DESC, id
    `, ownerID.UUID())
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
		id, err := board.BoardIdFromUUID(idRaw)
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
		out = append(out, board.RehydrateBoard(id, ownerID, t, d, createdAt, updatedAt))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *BoardRepository) Delete(ctx context.Context, id board.BoardId) error {
	ct, err := r.pg.Exec(ctx, `DELETE FROM boards WHERE id = $1`, id.UUID())
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return board.ErrBoardNotFound
	}
	return nil
}
