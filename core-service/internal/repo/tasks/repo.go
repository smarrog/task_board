package tasks

import (
	"context"

	"github.com/google/uuid"
	"github.com/smarrog/task-board/core-service/internal/model"
)

type Repo interface {
	Create(ctx context.Context, t *model.Task) error
	Get(ctx context.Context, id uuid.UUID) (*model.Task, error)
	List(ctx context.Context, columnId uuid.UUID) ([]model.Task, error)
	Update(ctx context.Context, t *model.Task) error
	Delete(ctx context.Context, id uuid.UUID) error
}
