package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/smarrog/task-board/core-service/internal/config"
	boardsrepo "github.com/smarrog/task-board/core-service/internal/repo/boards"
	boardsvc "github.com/smarrog/task-board/core-service/internal/service/boards"
	"github.com/smarrog/task-board/core-service/internal/transport/grpc"
)

type App struct {
	cfg  *config.Config
	grpc *grpc.Server
	pg   *pgxpool.Pool
	// services
	boards *boardsvc.Service
}

func New(cfg *config.Config, log *zerolog.Logger) (*App, error) {
	pg, err := newPG(cfg)
	if err != nil {
		return nil, err
	}

	boardsRepo := boardsrepo.NewRepo(pg)
	boardsSvc := boardsvc.NewSvc(boardsRepo)

	grpcServer := grpc.NewServer(log, boardsSvc)

	return &App{
		cfg:    cfg,
		grpc:   grpcServer,
		pg:     pg,
		boards: boardsSvc,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- a.grpc.Run(":" + a.cfg.GRPCPort)
	}()

	select {
	case <-ctx.Done():
		a.pg.Close()
		a.grpc.Stop()
		return nil
	case err := <-errCh:
		return err
	}
}

func newPG(cfg *config.Config) (*pgxpool.Pool, error) {
	pgCfg, err := pgxpool.ParseConfig(cfg.PostgresDSN)
	if err != nil {
		return nil, err
	}

	pgCfg.MinConns = int32(cfg.PostgresMinConns)
	pgCfg.MaxConns = int32(cfg.PostgresMaxConns)
	pgCfg.MaxConnIdleTime = cfg.PostgresMaxConnIdleTime
	pgCfg.MaxConnLifetime = cfg.PostgresMaxConnLifeTime

	pool, err := pgxpool.NewWithConfig(context.Background(), pgCfg)
	if err != nil {
		return nil, err
	}

	pingCtx, cancel := context.WithTimeout(context.Background(), cfg.PostgresTimeout)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
