package task

import (
	"context"
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/column"
	"github.com/smarrog/task-board/core-service/internal/domain/common"
	"github.com/smarrog/task-board/core-service/internal/domain/task"
	"github.com/smarrog/task-board/core-service/internal/usecase/cache"
)

type CreateTaskUseCase struct {
	repo    task.Repository
	columns column.Repository
	cache   cache.Invalidator
}

type CreateTaskInput struct {
	ColumnId    string
	Position    int
	Title       string
	Description string
	AssigneeId  string
}

type CreateTaskOutput struct {
	Task *task.Task
}

func NewCreateTaskUseCase(repo task.Repository, columns column.Repository, cache cache.Invalidator) *CreateTaskUseCase {
	return &CreateTaskUseCase{repo: repo, columns: columns, cache: cache}
}

func (uc *CreateTaskUseCase) Execute(ctx context.Context, input CreateTaskInput) (output *CreateTaskOutput, err error) {
	cid, err := column.IdFromString(input.ColumnId)
	if err != nil {
		return nil, err
	}
	position, err := task.NewPosition(input.Position)
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

	t := task.New(cid, position, title, desc, aid)

	err = uc.repo.Save(ctx, t)
	if err != nil {
		return nil, fmt.Errorf("save task: %w", err)
	}

	if uc.cache != nil {
		if c, err := uc.columns.Get(ctx, cid); err == nil {
			_ = uc.cache.InvalidateBoard(ctx, c.BoardId())
		}
	}

	output = &CreateTaskOutput{
		Task: t,
	}

	return output, nil
}
