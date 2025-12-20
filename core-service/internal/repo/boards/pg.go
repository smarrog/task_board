package boards

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/smarrog/task-board/core-service/internal/model"
)

var ErrNotFound = errors.New("board not found")

type pgRepo struct {
	pg *pgxpool.Pool
}

func NewRepo(pg *pgxpool.Pool) Repo {
	return &pgRepo{pg: pg}
}

func (r *pgRepo) Create(ctx context.Context, b *model.Board) error {
	return r.pg.QueryRow(ctx, `
		INSERT INTO boards (owner_id, title, description)
		VALUES ($1, $2, $3)
		RETURNING id
	`, b.OwnerId, b.Title, b.Description).Scan(&b.Id)
}

func (r *pgRepo) Get(ctx context.Context, id uuid.UUID) (*model.Board, error) {
	var b model.Board

	err := r.pg.QueryRow(ctx, `
		SELECT id, owner_id, title, description
		FROM boards
		WHERE id = $1
	`, id).Scan(&b.Id, &b.OwnerId, &b.Title, &b.Description)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &b, nil
}

func (r *pgRepo) List(ctx context.Context, ownerID uuid.UUID) ([]model.Board, error) {
	rows, err := r.pg.Query(ctx, `
		SELECT id, owner_id, title, description
		FROM boards
		WHERE owner_id = $1
		ORDER BY id
	`, ownerID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	out := make([]model.Board, 0)
	for rows.Next() {
		var b model.Board
		if err := rows.Scan(&b.Id, &b.OwnerId, &b.Title, &b.Description); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

func (r *pgRepo) Update(ctx context.Context, b *model.Board) error {
	ct, err := r.pg.Exec(ctx, `
		UPDATE boards
		SET title = $1,
		    description = $2,
			owner_id = $3
		WHERE id = $4
	`, b.Title, b.Description, b.OwnerId, b.Id)

	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *pgRepo) Delete(ctx context.Context, id uuid.UUID) error {
	ct, err := r.pg.Exec(ctx, `
		DELETE FROM boards
		WHERE id = $1
	`, id)

	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
