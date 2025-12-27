package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/smarrog/task-board/core-service/internal/domain/column"
	"github.com/smarrog/task-board/core-service/internal/domain/common"
	"github.com/smarrog/task-board/core-service/internal/domain/task"
)

type TasksRepo struct {
	pg     *pgxpool.Pool
	log    *zerolog.Logger
	outbox *OutboxRepo
}

func NewTasksRepo(pg *pgxpool.Pool, log *zerolog.Logger, outbox *OutboxRepo) *TasksRepo {
	return &TasksRepo{
		pg:     pg,
		log:    log,
		outbox: outbox,
	}
}

func (r *TasksRepo) Save(ctx context.Context, t *task.Task) error {
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
        INSERT INTO boards (id, column_id, position, title, description, assignee_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        ON CONFLICT (id) DO UPDATE
        SET owner_id    = EXCLUDED.owner_id,
            title       = EXCLUDED.title,
            description = EXCLUDED.description,
            updated_at  = EXCLUDED.updated_at
    `,
		t.Id().UUID(),
		t.ColumnId().UUID(),
		t.Position(),
		t.Title().String(),
		t.Description().String(),
		t.AssigneeId().UUID(),
		t.CreatedAt(),
		t.UpdatedAt(),
	)
	if err != nil {
		return err
	}

	events := t.PullEvents()
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

func (r *TasksRepo) Get(ctx context.Context, id task.Id) (*task.Task, error) {
	var columnIdRaw string
	var positionRaw int
	var titleRaw string
	var descRaw string
	var assigneeIdRaw string
	var createdAt, updatedAt time.Time

	err := r.pg.QueryRow(ctx, `
		SELECT column_id, position, title, description, assignee_id, createdAt, updatedAt
		FROM tasks
		WHERE id = $1
	`, id.UUID()).Scan(&columnIdRaw, &positionRaw, &titleRaw, &descRaw, &assigneeIdRaw, &createdAt, &updatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, task.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	columnId, err := column.IdFromString(columnIdRaw)
	if err != nil {
		return nil, err
	}
	position := positionRaw
	title, err := common.NewTitle(titleRaw)
	if err != nil {
		return nil, err
	}
	desc, err := common.NewDescription(descRaw)
	if err != nil {
		return nil, err
	}
	assigneeId, err := common.UserIdFromString(assigneeIdRaw)
	if err != nil {
		return nil, err
	}

	return task.Rehydrate(id, columnId, position, title, desc, assigneeId, createdAt, updatedAt), nil
}

func (r *TasksRepo) Delete(ctx context.Context, id task.Id) error {
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

	ct, err := tx.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, id.UUID())
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return task.ErrNotFound
	}

	// TODO send event to outbox

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
