package persistence

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"github.com/smarrog/task-board/core-service/internal/domain/common"
)

type OutboxRepo struct {
	txm *TxManager
	log *zerolog.Logger
}

func NewOutboxRepo(txm *TxManager) *OutboxRepo {
	return &OutboxRepo{
		txm: txm,
	}
}

func (r *OutboxRepo) SaveEvents(ctx context.Context, events []common.DomainEvent) error {
	if len(events) == 0 {
		return nil
	}

	return r.txm.InTx(ctx, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `
			// TODO
		`)
		if err != nil {
			return err
		}

		return nil
	})
}
