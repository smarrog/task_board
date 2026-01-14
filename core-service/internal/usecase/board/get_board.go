package board

import (
	"context"
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/board"
	"github.com/smarrog/task-board/core-service/internal/domain/column"
	"github.com/smarrog/task-board/core-service/internal/domain/task"
)

type GetBoardUseCase struct {
	boards  board.Repository
	columns column.Repository
	tasks   task.Repository
}

type GetBoardInput struct {
	BoardId string
}

type GetBoardOutput struct {
	Board   *board.Board
	Columns []*column.Column
	Tasks   []*task.Task
}

func NewGetBoardUseCase(boards board.Repository, columns column.Repository, tasks task.Repository) *GetBoardUseCase {
	return &GetBoardUseCase{boards: boards, columns: columns, tasks: tasks}
}

func (uc *GetBoardUseCase) Execute(ctx context.Context, input GetBoardInput) (*GetBoardOutput, error) {
	id, err := board.IdFromString(input.BoardId)
	if err != nil {
		return nil, err
	}

	b, err := uc.boards.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get board: %w", err)
	}

	cols, err := uc.columns.ListByBoard(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("list columns: %w", err)
	}

	columnIds := make([]column.Id, 0, len(cols))
	for _, c := range cols {
		columnIds = append(columnIds, c.Id())
	}

	tasksOut, err := uc.tasks.ListByColumns(ctx, columnIds)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}

	return &GetBoardOutput{Board: b, Columns: cols, Tasks: tasksOut}, nil
}
