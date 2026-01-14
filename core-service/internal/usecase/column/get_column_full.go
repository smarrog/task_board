package column

import (
    "context"
    "fmt"

    columndo "github.com/smarrog/task-board/core-service/internal/domain/column"
    taskdo "github.com/smarrog/task-board/core-service/internal/domain/task"
)

type GetColumnFullUseCase struct {
    columns columndo.Repository
    tasks   taskdo.Repository
}

type GetColumnFullInput struct {
    ColumnId string
}

type GetColumnFullOutput struct {
    Column *columndo.Column
    Tasks  []*taskdo.Task
}

func NewGetColumnFullUseCase(columns columndo.Repository, tasks taskdo.Repository) *GetColumnFullUseCase {
    return &GetColumnFullUseCase{columns: columns, tasks: tasks}
}

func (uc *GetColumnFullUseCase) Execute(ctx context.Context, input GetColumnFullInput) (*GetColumnFullOutput, error) {
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

    return &GetColumnFullOutput{Column: c, Tasks: ts}, nil
}
