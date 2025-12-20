package columns

import (
	"context"

	"github.com/google/uuid"
	"github.com/smarrog/task-board/core-service/internal/model"
)

type Repo interface {
	Create(ctx context.Context, c *model.Column) error
	Get(ctx context.Context, id uuid.UUID) (*model.Column, error)
	List(ctx context.Context, boardId uuid.UUID) ([]model.Column, error)
	Update(ctx context.Context, c *model.Column) error
	Delete(ctx context.Context, id uuid.UUID) error
}
