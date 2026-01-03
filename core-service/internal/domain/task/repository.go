package task

import (
	"context"

	"github.com/smarrog/task-board/core-service/internal/domain/column"
)

type Repository interface {
	InTx(ctx context.Context, fn func(ctx context.Context) error) error

	Save(ctx context.Context, t *Task) error
	Get(ctx context.Context, id Id) (*Task, error)
	Delete(ctx context.Context, id Id) error

	LockColumnTasks(ctx context.Context, columnId column.Id) error
	CountInColumn(ctx context.Context, columnId column.Id) (int, error)

	ShiftPositions(ctx context.Context, columnId column.Id, fromIncl, toIncl int, delta int) error
	ShiftAfterRemove(ctx context.Context, columnId column.Id, fromPos int) error
	ShiftForInsert(ctx context.Context, columnId column.Id, toPos int) error
}
