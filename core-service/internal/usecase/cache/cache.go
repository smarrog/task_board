package cache

import (
	"context"
	"time"

	"github.com/smarrog/task-board/core-service/internal/domain/board"
	"github.com/smarrog/task-board/core-service/internal/domain/column"
	"github.com/smarrog/task-board/core-service/internal/domain/task"
)

type BoardData struct {
	Board   *board.Board
	Columns []*column.Column
	Tasks   []*task.Task
}

type Cacher interface {
	Invalidator

	GetBoard(ctx context.Context, id board.Id) (out *BoardData, hit bool, err error)
	SetBoard(ctx context.Context, id board.Id, out *BoardData, ttl time.Duration) error
}

type Invalidator interface {
	InvalidateBoard(ctx context.Context, id board.Id) error
}
