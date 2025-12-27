package persistence

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/smarrog/task-board/core-service/internal/domain/common"
)

type OutboxRepo struct {
	pg *pgxpool.Pool
}

func NewOutboxRepo(pg *pgxpool.Pool) *OutboxRepo {
	return &OutboxRepo{
		pg: pg,
	}
}

func (r *OutboxRepo) SaveEvents(ctx context.Context, tx pgx.Tx, events []common.DomainEvent) error {
	// TODO
	return nil
}
