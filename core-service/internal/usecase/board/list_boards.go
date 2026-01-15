package board

import (
	"context"
	"fmt"
	"time"

	"github.com/smarrog/task-board/core-service/internal/domain/board"
	"github.com/smarrog/task-board/core-service/internal/domain/column"
	"github.com/smarrog/task-board/core-service/internal/domain/common"
	"github.com/smarrog/task-board/core-service/internal/domain/task"
	"github.com/smarrog/task-board/core-service/internal/usecase/cache"
)

type ListBoardsUseCase struct {
	boards  board.Repository
	columns column.Repository
	tasks   task.Repository
	cache   cache.Cacher
	ttl     time.Duration
}

type ListBoardsInput struct {
	OwnerId string
}

type ListBoardsOutput struct {
	Items []*GetBoardOutput
}

func NewListBoardsUseCase(
	boards board.Repository,
	columns column.Repository,
	tasks task.Repository,
	cache cache.Cacher,
	ttl time.Duration,
) *ListBoardsUseCase {
	return &ListBoardsUseCase{boards: boards, columns: columns, tasks: tasks, cache: cache, ttl: ttl}
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

	hitByBoard := make(map[string]*cache.BoardData, len(boardsList))
	missIDs := make([]board.Id, 0)
	if uc.cache != nil {
		for _, b := range boardsList {
			bid := b.Id()
			cached, hit, err := uc.cache.GetBoard(ctx, bid)
			if err == nil && hit && cached != nil {
				hitByBoard[bid.String()] = cached
				continue
			}
			missIDs = append(missIDs, bid)
		}
	} else {
		for _, b := range boardsList {
			missIDs = append(missIDs, b.Id())
		}
	}

	missOut := make(map[string]*GetBoardOutput, len(missIDs))
	if len(missIDs) > 0 {
		boardByID := make(map[string]*board.Board, len(boardsList))
		for _, b := range boardsList {
			boardByID[b.Id().String()] = b
		}

		cols, err := uc.columns.ListByBoards(ctx, missIDs)
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

		var ts []*task.Task
		if len(columnIDs) > 0 {
			ts, err = uc.tasks.ListByColumns(ctx, columnIDs)
			if err != nil {
				return nil, fmt.Errorf("list tasks: %w", err)
			}
		}

		tasksByColumn := make(map[string][]*task.Task)
		for _, t := range ts {
			cid := t.ColumnId().String()
			tasksByColumn[cid] = append(tasksByColumn[cid], t)
		}

		for _, bid := range missIDs {
			b := boardByID[bid.String()]
			bcols := colsByBoard[bid.String()]
			flatTasks := make([]*task.Task, 0)
			for _, c := range bcols {
				flatTasks = append(flatTasks, tasksByColumn[c.Id().String()]...)
			}

			out := &GetBoardOutput{Board: b, Columns: bcols, Tasks: flatTasks}
			missOut[bid.String()] = out

			if uc.cache != nil {
				_ = uc.cache.SetBoard(ctx, bid, &cache.BoardData{Board: out.Board, Columns: out.Columns, Tasks: out.Tasks}, uc.ttl)
			}
		}
	}

	items := make([]*GetBoardOutput, 0, len(boardsList))
	for _, b := range boardsList {
		bid := b.Id().String()
		if cached, ok := hitByBoard[bid]; ok {
			items = append(items, &GetBoardOutput{Board: cached.Board, Columns: cached.Columns, Tasks: cached.Tasks})
			continue
		}
		if out, ok := missOut[bid]; ok {
			items = append(items, out)
			continue
		}
	}

	return &ListBoardsOutput{Items: items}, nil
}
