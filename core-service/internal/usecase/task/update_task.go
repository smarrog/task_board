package task

import (
	"context"
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/column"
	"github.com/smarrog/task-board/core-service/internal/domain/common"
	"github.com/smarrog/task-board/core-service/internal/domain/task"
	"github.com/smarrog/task-board/core-service/internal/usecase/cache"
)

type UpdateTaskUseCase struct {
	repo    task.Repository
	columns column.Repository
	cache   cache.Invalidator
}

type UpdateTaskInput struct {
	TaskId      string
	Title       string
	Description string
	AssigneeId  string
}

type UpdateTaskOutput struct {
	Task *task.Task
}

func NewUpdateTaskUseCase(repo task.Repository, columns column.Repository, cache cache.Invalidator) *UpdateTaskUseCase {
	return &UpdateTaskUseCase{repo: repo, columns: columns, cache: cache}
}

func (uc *UpdateTaskUseCase) Execute(ctx context.Context, input UpdateTaskInput) (output *UpdateTaskOutput, err error) {
	tid, err := task.IdFromString(input.TaskId)
	if err != nil {
		return nil, err
	}
	title, err := task.NewTitle(input.Title)
	if err != nil {
		return nil, err
	}
	desc, err := task.NewDescription(input.Description)
	if err != nil {
		return nil, err
	}
	aid, err := common.UserIdFromString(input.AssigneeId)
	if err != nil {
		return nil, fmt.Errorf("task assignee_id: %w", err)
	}

	t, err := uc.repo.Get(ctx, tid)
	if err != nil {
		return nil, err
	}

	t.Update(title, desc, aid)

	err = uc.repo.Save(ctx, t)
	if err != nil {
		return nil, fmt.Errorf("save task: %w", err)
	}

	if uc.cache != nil {
		if c, err := uc.columns.Get(ctx, t.ColumnId()); err == nil {
			_ = uc.cache.InvalidateBoard(ctx, c.BoardId())
		}
	}

	output = &UpdateTaskOutput{
		Task: t,
	}

	return output, nil
}
