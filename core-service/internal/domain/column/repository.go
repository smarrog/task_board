package column

import (
	"context"

	"github.com/smarrog/task-board/core-service/internal/domain/board"
)

type Repository interface {
	InTx(ctx context.Context, fn func(ctx context.Context) error) error

	Save(ctx context.Context, c *Column) error
	Get(ctx context.Context, id Id) (*Column, error)
	ListByBoard(ctx context.Context, boardId board.Id) ([]*Column, error)
	ListByBoards(ctx context.Context, boardIds []board.Id) ([]*Column, error)
	Delete(ctx context.Context, id Id) error

	LockBoardColumns(ctx context.Context, boardId board.Id) error
	CountInBoard(ctx context.Context, boardId board.Id) (int, error)

	ShiftPositions(ctx context.Context, boardId board.Id, fromIncl, toIncl int, delta int) error
}
