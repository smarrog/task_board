package task

import (
	"context"
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/column"
	"github.com/smarrog/task-board/core-service/internal/domain/common"
	"github.com/smarrog/task-board/core-service/internal/domain/task"
)

type CreateTaskUseCase struct {
	repo task.Repository
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

func NewCreateTaskUseCase(repo task.Repository) *CreateTaskUseCase {
	return &CreateTaskUseCase{repo: repo}
}

func (uc *CreateTaskUseCase) Execute(ctx context.Context, input CreateTaskInput) (output *CreateTaskOutput, err error) {
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

	t := task.New(cid, position, title, desc, aid)

	err = uc.repo.Save(ctx, t)
	if err != nil {
		return nil, fmt.Errorf("save column: %w", err)
	}

	output = &CreateTaskOutput{
		Task: t,
	}

	return output, nil
}
