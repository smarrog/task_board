package task

import (
	"context"
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/common"
	"github.com/smarrog/task-board/core-service/internal/domain/task"
)

type UpdateTaskUseCase struct {
	repo task.Repository
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

func NewUpdateTaskUseCase(repo task.Repository) *UpdateTaskUseCase {
	return &UpdateTaskUseCase{repo: repo}
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

	output = &UpdateTaskOutput{
		Task: t,
	}

	return output, nil
}
