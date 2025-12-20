package boards

import (
	"context"

	"github.com/google/uuid"
	"github.com/smarrog/task-board/core-service/internal/model"
)

type Repo interface {
	Create(ctx context.Context, b *model.Board) error
	Get(ctx context.Context, id uuid.UUID) (*model.Board, error)
	List(ctx context.Context, ownerID uuid.UUID) ([]model.Board, error)
	Update(ctx context.Context, b *model.Board) error
	Delete(ctx context.Context, id uuid.UUID) error
}
