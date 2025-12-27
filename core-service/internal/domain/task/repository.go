package task

import "context"

type Repository interface {
	Save(ctx context.Context, t *Task) error
	Get(ctx context.Context, id Id) (*Task, error)
	Delete(ctx context.Context, id Id) error
}
