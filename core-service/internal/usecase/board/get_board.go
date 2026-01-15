package board

import (
	"context"
	"fmt"
	"time"

	"github.com/smarrog/task-board/core-service/internal/domain/board"
	"github.com/smarrog/task-board/core-service/internal/domain/column"
	"github.com/smarrog/task-board/core-service/internal/domain/task"
	"github.com/smarrog/task-board/core-service/internal/usecase/cache"
)

type GetBoardUseCase struct {
	boards  board.Repository
	columns column.Repository
	tasks   task.Repository
	cache   cache.Cacher
	ttl     time.Duration
}

type GetBoardInput struct {
	BoardId string
}

type GetBoardOutput struct {
	Board   *board.Board
	Columns []*column.Column
	Tasks   []*task.Task
}

func NewGetBoardUseCase(
	boards board.Repository,
	columns column.Repository,
	tasks task.Repository,
	cache cache.Cacher,
	ttl time.Duration,
) *GetBoardUseCase {
	return &GetBoardUseCase{boards: boards, columns: columns, tasks: tasks, cache: cache, ttl: ttl}
}

func (uc *GetBoardUseCase) Execute(ctx context.Context, input GetBoardInput) (*GetBoardOutput, error) {
	id, err := board.IdFromString(input.BoardId)
	if err != nil {
		return nil, err
	}

	if uc.cache != nil {
		if cached, hit, err := uc.cache.GetBoard(ctx, id); err == nil && hit && cached != nil {
			return &GetBoardOutput{Board: cached.Board, Columns: cached.Columns, Tasks: cached.Tasks}, nil
		} else if err != nil {
			// cache errors must not break the request path
		}
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

	out := &GetBoardOutput{Board: b, Columns: cols, Tasks: tasksOut}

	if uc.cache != nil {
		_ = uc.cache.SetBoard(ctx, id, &cache.BoardData{Board: out.Board, Columns: out.Columns, Tasks: out.Tasks}, uc.ttl)
	}

	return out, nil
}
