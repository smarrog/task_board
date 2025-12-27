package task

import (
	"context"
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/column"
	"github.com/smarrog/task-board/core-service/internal/domain/common"
	"github.com/smarrog/task-board/core-service/internal/domain/task"
)

type UpdateTaskUseCase struct {
	repo task.Repository
}

type UpdateTaskInput struct {
	TaskId      string
	ColumnId    string
	Position    int
	Title       string
	Description string
	AssigneeId  string
}

type UpdateTaskOutput struct {
	Task *task.Task
}

func NewUpdateTaskUseCase(repo task.Repository) *UpdateTaskUseCase {
	return &UpdateTaskUseCase{repo: repo}
}

func (uc *UpdateTaskUseCase) Execute(ctx context.Context, input UpdateTaskInput) (output *UpdateTaskOutput, err error) {
	tid, err := task.IdFromString(input.TaskId)
	if err != nil {
		return nil, err
	}
	cid, err := column.IdFromString(input.ColumnId)
	if err != nil {
		return nil, err
	}
	position := input.Position
	title, err := common.NewTitle(input.Title)
	if err != nil {
		return nil, fmt.Errorf("column: %w", err)
	}
	desc, err := common.NewDescription(input.Description)
	if err != nil {
		return nil, fmt.Errorf("column: %w", err)
	}
	aid, err := common.UserIdFromString(input.AssigneeId)
	if err != nil {
		return nil, fmt.Errorf("column assignee_id: %w", err)
	}

	t, err := uc.repo.Get(ctx, tid)
	if err != nil {
		return nil, err
	}

	t.Update(cid, position, title, desc, aid)

	err = uc.repo.Save(ctx, t)
	if err != nil {
		return nil, fmt.Errorf("save column: %w", err)
	}

	output = &UpdateTaskOutput{
		Task: t,
	}

	return output, nil
}
