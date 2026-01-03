package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"github.com/smarrog/task-board/core-service/internal/domain/column"
	"github.com/smarrog/task-board/core-service/internal/domain/common"
	"github.com/smarrog/task-board/core-service/internal/domain/task"
)

type TasksRepo struct {
	txm    *TxManager
	log    *zerolog.Logger
	outbox *OutboxRepo
}

func NewTasksRepo(txm *TxManager, log *zerolog.Logger, outbox *OutboxRepo) *TasksRepo {
	return &TasksRepo{
		txm:    txm,
		log:    log,
		outbox: outbox,
	}
}

func (r *TasksRepo) InTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.txm.InTx(ctx, func(ctx context.Context, _ pgx.Tx) error {
		return fn(ctx)
	})
}

func (r *TasksRepo) Save(ctx context.Context, t *task.Task) error {
	return r.txm.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `
			INSERT INTO tasks (id, column_id, position, title, description, assignee_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT (id) DO UPDATE
			SET column_id   = EXCLUDED.column_id,
				position    = EXCLUDED.position,
				title       = EXCLUDED.title,
				description = EXCLUDED.description,
				assignee_id = EXCLUDED.assignee_id,
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
		return r.outbox.SaveEvents(ctx, events)
	})
}

func (r *TasksRepo) Get(ctx context.Context, id task.Id) (*task.Task, error) {
	db := r.txm.DB(ctx)

	var columnIdRaw string
	var positionRaw int
	var titleRaw string
	var descRaw string
	var assigneeIdRaw string
	var createdAt, updatedAt time.Time

	err := db.QueryRow(ctx, `
		SELECT column_id, position, title, description, assignee_id, created_at, updated_at
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
	position, err := task.NewPosition(positionRaw)
	if err != nil {
		return nil, err
	}
	title, err := task.NewTitle(titleRaw)
	if err != nil {
		return nil, err
	}
	desc, err := task.NewDescription(descRaw)
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
	return r.txm.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		ct, err := tx.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, id.UUID())
		if err != nil {
			return err
		}
		if ct.RowsAffected() == 0 {
			return task.ErrNotFound
		}

		// TODO send event to outbox

		return nil
	})
}

func (r *TasksRepo) LockColumnTasks(ctx context.Context, columnId column.Id) error {
	return r.txm.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `SELECT id FROM tasks WHERE column_id=$1 FOR UPDATE`, columnId.UUID())
		return err
	})
}

func (r *TasksRepo) CountInColumn(ctx context.Context, columnId column.Id) (int, error) {
	var n int
	err := r.txm.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		return tx.QueryRow(ctx, `SELECT COUNT(*) FROM tasks WHERE column_id=$1`, columnId.UUID()).Scan(&n)
	})
	return n, err
}

func (r *TasksRepo) ShiftPositions(ctx context.Context, columnId column.Id, fromIncl, toIncl int, delta int) error {
	if fromIncl > toIncl || delta == 0 {
		return nil
	}
	return r.txm.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `
			UPDATE tasks
			SET position = position + $4
			WHERE column_id=$1 AND position BETWEEN $2 AND $3
		`, columnId.UUID(), fromIncl, toIncl, delta)
		return err
	})
}

func (r *TasksRepo) ShiftAfterRemove(ctx context.Context, columnId column.Id, fromPos int) error {
	return r.txm.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `
			UPDATE tasks
			SET position = position - 1
			WHERE column_id=$1 AND position > $2
		`, columnId.UUID(), fromPos)
		return err
	})
}

func (r *TasksRepo) ShiftForInsert(ctx context.Context, columnId column.Id, toPos int) error {
	return r.txm.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `
			UPDATE tasks
			SET position = position + 1
			WHERE column_id=$1 AND position >= $2
		`, columnId.UUID(), toPos)
		return err
	})
}
