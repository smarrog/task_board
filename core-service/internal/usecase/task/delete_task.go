package task

import (
	"context"
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/task"
)

type DeleteTaskUseCase struct {
	repo task.Repository
}

type DeleteTaskInput struct {
	TaskId string
}

type DeleteTaskOutput struct {
}

func NewDeleteTaskUseCase(repo task.Repository) *DeleteTaskUseCase {
	return &DeleteTaskUseCase{repo: repo}
}

func (uc *DeleteTaskUseCase) Execute(ctx context.Context, input DeleteTaskInput) (*DeleteTaskOutput, error) {
	tid, err := task.IdFromString(input.TaskId)
	if err != nil {
		return nil, err
	}

	err = uc.repo.Delete(ctx, tid)
	if err != nil {
		return nil, fmt.Errorf("delete task: %w", err)
	}

	output := &DeleteTaskOutput{}

	return output, nil
}
