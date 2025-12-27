package task

import (
	"context"
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/task"
)

type GetTaskUseCase struct {
	repo task.Repository
}

type GetTaskInput struct {
	TaskId string
}

type GetTaskOutput struct {
	Task *task.Task
}

func NewGetTaskUseCase(repo task.Repository) *GetTaskUseCase {
	return &GetTaskUseCase{repo: repo}
}

func (uc *GetTaskUseCase) Execute(ctx context.Context, input GetTaskInput) (*GetTaskOutput, error) {
	id, err := task.IdFromString(input.TaskId)
	if err != nil {
		return nil, err
	}

	t, err := uc.repo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get task: %w", err)
	}

	output := &GetTaskOutput{
		Task: t,
	}

	return output, nil
}
