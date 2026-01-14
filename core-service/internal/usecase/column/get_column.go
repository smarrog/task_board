package column

import (
	"context"
	"fmt"

	columndo "github.com/smarrog/task-board/core-service/internal/domain/column"
	taskdo "github.com/smarrog/task-board/core-service/internal/domain/task"
)

type GetColumnUseCase struct {
	columns columndo.Repository
	tasks   taskdo.Repository
}

type GetColumnInput struct {
	ColumnId string
}

type GetColumnOutput struct {
	Column *columndo.Column
	Tasks  []*taskdo.Task
}

func NewGetColumnUseCase(columns columndo.Repository, tasks taskdo.Repository) *GetColumnUseCase {
	return &GetColumnUseCase{columns: columns, tasks: tasks}
}

func (uc *GetColumnUseCase) Execute(ctx context.Context, input GetColumnInput) (*GetColumnOutput, error) {
	id, err := columndo.IdFromString(input.ColumnId)
	if err != nil {
		return nil, err
	}

	c, err := uc.columns.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get column: %w", err)
	}

	ts, err := uc.tasks.ListByColumn(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("list tasks by column: %w", err)
	}

	return &GetColumnOutput{Column: c, Tasks: ts}, nil
}
