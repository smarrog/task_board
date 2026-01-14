package board

import (
	"context"
	"fmt"

	"github.com/smarrog/task-board/core-service/internal/domain/board"
	"github.com/smarrog/task-board/core-service/internal/domain/column"
	"github.com/smarrog/task-board/core-service/internal/domain/common"
	"github.com/smarrog/task-board/core-service/internal/domain/task"
)

type ListBoardsUseCase struct {
	boards  board.Repository
	columns column.Repository
	tasks   task.Repository
}

type ListBoardsInput struct {
	OwnerId string
}

type ListBoardsOutput struct {
	Items []*GetBoardOutput
}

func NewListBoardsUseCase(boards board.Repository, columns column.Repository, tasks task.Repository) *ListBoardsUseCase {
	return &ListBoardsUseCase{boards: boards, columns: columns, tasks: tasks}
}

func (uc *ListBoardsUseCase) Execute(ctx context.Context, input ListBoardsInput) (*ListBoardsOutput, error) {
	oid, err := common.UserIdFromString(input.OwnerId)
	if err != nil {
		return nil, fmt.Errorf("owner_id: %w", err)
	}

	boardsList, err := uc.boards.ListByOwner(ctx, oid)
	if err != nil {
		return nil, err
	}
	if len(boardsList) == 0 {
		return &ListBoardsOutput{Items: []*GetBoardOutput{}}, nil
	}

	boardIDs := make([]board.Id, 0, len(boardsList))
	boardByID := make(map[string]*board.Board, len(boardsList))
	for _, b := range boardsList {
		boardIDs = append(boardIDs, b.Id())
		boardByID[b.Id().String()] = b
	}

	cols, err := uc.columns.ListByBoards(ctx, boardIDs)
	if err != nil {
		return nil, fmt.Errorf("list columns: %w", err)
	}

	colsByBoard := make(map[string][]*column.Column)
	columnIDs := make([]column.Id, 0, len(cols))
	for _, c := range cols {
		bid := c.BoardId().String()
		colsByBoard[bid] = append(colsByBoard[bid], c)
		columnIDs = append(columnIDs, c.Id())
	}

	ts, err := uc.tasks.ListByColumns(ctx, columnIDs)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}

	tasksByColumn := make(map[string][]*task.Task)
	for _, t := range ts {
		cid := t.ColumnId().String()
		tasksByColumn[cid] = append(tasksByColumn[cid], t)
	}

	items := make([]*GetBoardOutput, 0, len(boardsList))
	for _, b := range boardsList {
		bid := b.Id().String()
		bcols := colsByBoard[bid]
		flatTasks := make([]*task.Task, 0)
		for _, c := range bcols {
			flatTasks = append(flatTasks, tasksByColumn[c.Id().String()]...)
		}
		items = append(items, &GetBoardOutput{
			Board:   b,
			Columns: bcols,
			Tasks:   flatTasks,
		})
	}

	return &ListBoardsOutput{Items: items}, nil
}
