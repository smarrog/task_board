package task

import (
	"context"
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/column"
	"github.com/smarrog/task-board/core-service/internal/domain/task"
	"github.com/smarrog/task-board/core-service/internal/usecase/cache"
)

type DeleteTaskUseCase struct {
	repo    task.Repository
	columns column.Repository
	cache   cache.Invalidator
}

type DeleteTaskInput struct {
	TaskId string
}

type DeleteTaskOutput struct {
}

func NewDeleteTaskUseCase(repo task.Repository, columns column.Repository, cache cache.Invalidator) *DeleteTaskUseCase {
	return &DeleteTaskUseCase{repo: repo, columns: columns, cache: cache}
}

func (uc *DeleteTaskUseCase) Execute(ctx context.Context, input DeleteTaskInput) (*DeleteTaskOutput, error) {
	tid, err := task.IdFromString(input.TaskId)
	if err != nil {
		return nil, err
	}

	t, err := uc.repo.Get(ctx, tid)
	if err != nil {
		return nil, err
	}

	err = uc.repo.Delete(ctx, tid)
	if err != nil {
		return nil, fmt.Errorf("delete task: %w", err)
	}

	if uc.cache != nil {
		if c, err := uc.columns.Get(ctx, t.ColumnId()); err == nil {
			_ = uc.cache.InvalidateBoard(ctx, c.BoardId())
		}
	}

	output := &DeleteTaskOutput{}

	return output, nil
}
