package board

import (
	"context"

	"github.com/smarrog/task-board/core-service/internal/domain/common"
)

type Repository interface {
	InTx(ctx context.Context, fn func(ctx context.Context) error) error

	Save(ctx context.Context, b *Board) error
	Get(ctx context.Context, id Id) (*Board, error)
	ListByOwner(ctx context.Context, ownerId common.UserId) ([]*Board, error)
	Delete(ctx context.Context, id Id) error
}
