package column

import "context"

type Repository interface {
	Save(ctx context.Context, c *Column) error
	Get(ctx context.Context, id Id) (*Column, error)
	Delete(ctx context.Context, id Id) error
}
