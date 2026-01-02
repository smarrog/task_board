package persistence

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type txCtxKey struct{}
type TxFunc func(ctx context.Context, tx pgx.Tx) error

type TxManager struct {
	pg  *pgxpool.Pool
	log *zerolog.Logger
}

func NewTxManager(pg *pgxpool.Pool, log *zerolog.Logger) *TxManager {
	return &TxManager{pg: pg, log: log}
}

// DB — минимальный интерфейс, общий для pgx.Tx и pgxpool.Pool
type DB interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(context.Context, *pgx.Batch) pgx.BatchResults
}

// DB если есть транзакция в ctx — используем её, иначе pool.
func (m *TxManager) DB(ctx context.Context) DB {
	if tx, ok := txFromCtx(ctx); ok {
		return tx
	}
	return m.pg
}

// InTx открывает транзакцию, если её ещё нет в ctx.
// Если tx уже есть — просто выполняет fn(ctx) внутри существующей tx.
func (m *TxManager) InTx(ctx context.Context, fn TxFunc) (err error) {
	if tx, ok := txFromCtx(ctx); ok {
		return fn(ctx, tx)
	}

	tx, err := m.pg.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				m.log.Error().Err(rbErr).Msg("failed to rollback tx")
			}
		}
	}()

	ctx = withTx(ctx, tx)

	if err = fn(ctx, tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func withTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txCtxKey{}, tx)
}

func txFromCtx(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txCtxKey{}).(pgx.Tx)
	return tx, ok
}
