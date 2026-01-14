package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/smarrog/task-board/auth-service/internal/config"
	ps "github.com/smarrog/task-board/auth-service/internal/infrastructure/persistence"
	"github.com/smarrog/task-board/auth-service/internal/transport/grpc"
	uc "github.com/smarrog/task-board/auth-service/internal/usecase"
	"github.com/smarrog/task-board/shared/logger"
)

type App struct {
	log  *zerolog.Logger
	cfg  *config.Config
	pg   *pgxpool.Pool
	grpc *grpc.Server
}

func (a *App) Init() error {
	cfg := config.Load()
	log := logger.New(cfg.AppName, cfg.LogLevel)
	a.cfg = cfg
	a.log = log

	pool, err := newPG(cfg)
	if err != nil {
		return err
	}
	a.pg = pool

	repo := ps.NewUsersRepo(pool, log)

	h := createAuthHandler(log, cfg, repo)

	a.grpc = grpc.NewServer(log, h)

	return nil
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- a.grpc.Run(":" + a.cfg.GRPCPort)
	}()

	select {
	case <-ctx.Done():
		a.grpc.Stop()
		a.pg.Close()
		return nil
	case err := <-errCh:
		return err
	}
}

func (a *App) Log() *zerolog.Logger { return a.log }

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

func createAuthHandler(log *zerolog.Logger, cfg *config.Config, repo *ps.UsersRepo) *grpc.AuthHandler {
	register := uc.NewRegisterUseCase(repo, cfg)
	login := uc.NewLoginUseCase(repo, cfg)

	handler := grpc.NewAuthHandler(log, register, login)
	return handler
}
