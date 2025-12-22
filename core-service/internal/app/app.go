package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/smarrog/task-board/core-service/internal/config"
	"github.com/smarrog/task-board/core-service/internal/domain/board"
	"github.com/smarrog/task-board/core-service/internal/infrastructure/persistence"
	"github.com/smarrog/task-board/core-service/internal/transport/grpc"
	boarduc "github.com/smarrog/task-board/core-service/internal/usecase/board"
)

type App struct {
	cfg  *config.Config
	grpc *grpc.Server
	pg   *pgxpool.Pool
	// usecases (kept for future extension)
	createBoard *boarduc.CreateBoardUseCase
}

func New(cfg *config.Config, log *zerolog.Logger) (*App, error) {
	pg, err := newPG(cfg)
	if err != nil {
		return nil, err
	}

	boardsRepo := persistence.NewBoardRepository(pg)

	boardsHandler := createBoardsHandler(log, boardsRepo)
	grpcServer := grpc.NewServer(log, boardsHandler)

	return &App{
		cfg:  cfg,
		grpc: grpcServer,
		pg:   pg,
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

func createBoardsHandler(log *zerolog.Logger, boardsRepo board.Repository) *grpc.BoardsHandler {
	createBoard := boarduc.NewCreateBoardUseCase(boardsRepo)
	getBoard := boarduc.NewGetBoardUseCase(boardsRepo)
	listBoards := boarduc.NewListBoardsUseCase(boardsRepo)
	updateBoard := boarduc.NewUpdateBoardUseCase(boardsRepo)
	deleteBoard := boarduc.NewDeleteBoardUseCase(boardsRepo)

	boardsHandler := grpc.NewBoardsHandler(log, createBoard, getBoard, listBoards, updateBoard, deleteBoard)
	return boardsHandler
}
