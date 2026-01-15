package domain

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, u *User) error

	GetById(ctx context.Context, id UserId) (*User, error)
	GetByEmail(ctx context.Context, email Email) (*User, error)
}
