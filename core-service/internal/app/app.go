package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/smarrog/task-board/core-service/internal/config"
	"github.com/smarrog/task-board/core-service/internal/domain/board"
	"github.com/smarrog/task-board/core-service/internal/domain/column"
	"github.com/smarrog/task-board/core-service/internal/domain/task"
	appkafka "github.com/smarrog/task-board/core-service/internal/infrastructure/kafka"
	"github.com/smarrog/task-board/core-service/internal/infrastructure/persistence"
	"github.com/smarrog/task-board/core-service/internal/transport/grpc"
	boarduc "github.com/smarrog/task-board/core-service/internal/usecase/board"
	columnuc "github.com/smarrog/task-board/core-service/internal/usecase/column"
	taskuc "github.com/smarrog/task-board/core-service/internal/usecase/task"
	"github.com/smarrog/task-board/shared/logger"
)

type App struct {
	log           *zerolog.Logger
	cfg           *config.Config
	grpc          *grpc.Server
	pg            *pgxpool.Pool
	outboxWorker  *persistence.OutboxWorker
	kafkaProducer *appkafka.Producer
}

func (a *App) Init() error {
	cfg := config.Load()
	log := logger.New(cfg.AppName, cfg.LogLevel)

	a.cfg = cfg
	a.log = log

	pg, err := newPG(cfg)
	if err != nil {
		return err
	}

	a.pg = pg

	txm := persistence.NewTxManager(pg, log)

	outboxRepo := persistence.NewOutboxRepo(txm, log)
	boardsRepo := persistence.NewBoardsRepo(txm, log, outboxRepo)
	columnsRepo := persistence.NewColumnsRepo(txm, log, outboxRepo)
	tasksRepo := persistence.NewTasksRepo(txm, log, outboxRepo)

	boardsHandler := createBoardsHandler(log, boardsRepo, columnsRepo, tasksRepo)
	columnsHandler := createColumnsHandler(log, columnsRepo, tasksRepo)
	tasksHandler := createTasksHandler(log, tasksRepo)

	a.grpc = grpc.NewServer(log, boardsHandler, columnsHandler, tasksHandler)

	producer, err := appkafka.NewProducer(cfg, log)
	if err != nil {
		return err
	}

	a.kafkaProducer = producer
	a.outboxWorker = persistence.NewOutboxWorker(txm, outboxRepo, producer, cfg.KafkaTopic[0], cfg.OutboxBatchSize, cfg.OutboxPollInterval, log)

	return nil
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	if a.outboxWorker != nil {
		go func() {
			errCh <- a.outboxWorker.Run(ctx)
		}()
	}
	go func() {
		errCh <- a.grpc.Run(":" + a.cfg.GRPCPort)
	}()

	select {
	case <-ctx.Done():
		if a.kafkaProducer != nil {
			a.kafkaProducer.Close()
		}
		a.pg.Close()
		a.grpc.Stop()
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

func createBoardsHandler(
	log *zerolog.Logger,
	boardsRepo board.Repository,
	columnsRepo column.Repository,
	tasksRepo task.Repository,
) *grpc.BoardsHandler {
	createBoard := boarduc.NewCreateBoardUseCase(boardsRepo)
	getBoard := boarduc.NewGetBoardUseCase(boardsRepo, columnsRepo, tasksRepo)
	listBoards := boarduc.NewListBoardsUseCase(boardsRepo, columnsRepo, tasksRepo)
	updateBoard := boarduc.NewUpdateBoardUseCase(boardsRepo)
	deleteBoard := boarduc.NewDeleteBoardUseCase(boardsRepo)

	boardsHandler := grpc.NewBoardsHandler(log, createBoard, getBoard, listBoards, updateBoard, deleteBoard)
	return boardsHandler
}

func createColumnsHandler(log *zerolog.Logger, columnsRepo column.Repository, tasksRepo task.Repository) *grpc.ColumnsHandler {
	createColumn := columnuc.NewCreateColumnUseCase(columnsRepo)
	getColumn := columnuc.NewGetColumnUseCase(columnsRepo, tasksRepo)
	moveColumn := columnuc.NewMoveColumnUseCase(columnsRepo)
	deleteColumn := columnuc.NewDeleteColumnUseCase(columnsRepo)

	columnsHandler := grpc.NewColumnsHandler(log, createColumn, getColumn, moveColumn, deleteColumn)
	return columnsHandler
}

func createTasksHandler(log *zerolog.Logger, tasksRepo task.Repository) *grpc.TasksHandler {
	createTask := taskuc.NewCreateTaskUseCase(tasksRepo)
	getTask := taskuc.NewGetTaskUseCase(tasksRepo)
	updateTask := taskuc.NewUpdateTaskUseCase(tasksRepo)
	moveTask := taskuc.NewMoveTaskUseCase(tasksRepo)
	deleteTask := taskuc.NewDeleteTaskUseCase(tasksRepo)

	tasksHandler := grpc.NewTasksHandler(log, createTask, getTask, updateTask, moveTask, deleteTask)
	return tasksHandler
}
