package board

import "context"

type Repository interface {
	Save(ctx context.Context, b *Board) error
	Get(ctx context.Context, id BoardId) (*Board, error)
	ListByOwner(ctx context.Context, ownerID UserId) ([]*Board, error)
	Delete(ctx context.Context, id BoardId) error
}
