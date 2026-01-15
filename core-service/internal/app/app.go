package app

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/smarrog/task-board/core-service/internal/config"
	boarddo "github.com/smarrog/task-board/core-service/internal/domain/board"
	columndo "github.com/smarrog/task-board/core-service/internal/domain/column"
	taskdo "github.com/smarrog/task-board/core-service/internal/domain/task"
	appcache "github.com/smarrog/task-board/core-service/internal/infrastructure/cache"
	appkafka "github.com/smarrog/task-board/core-service/internal/infrastructure/kafka"
	"github.com/smarrog/task-board/core-service/internal/infrastructure/persistence"
	"github.com/smarrog/task-board/core-service/internal/transport/grpc"
	boarduc "github.com/smarrog/task-board/core-service/internal/usecase/board"
	commonuc "github.com/smarrog/task-board/core-service/internal/usecase/cache"
	columnuc "github.com/smarrog/task-board/core-service/internal/usecase/column"
	taskuc "github.com/smarrog/task-board/core-service/internal/usecase/task"
	"github.com/smarrog/task-board/shared/logger"
)

var ErrDisabled = errors.New("redis disabled")

type App struct {
	log           *zerolog.Logger
	cfg           *config.Config
	grpc          *grpc.Server
	pg            *pgxpool.Pool
	redis         *redis.Client
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

	rdb, err := newRedis(cfg)
	if err != nil {
		// если редис не поднялся, да и хрен с ним
		log.Err(err).Msg("Redis will not be started")
	}

	a.redis = rdb

	var cache commonuc.Cacher
	if rdb != nil {
		cache = appcache.NewRedisCache(rdb)
	}

	txm := persistence.NewTxManager(pg, log)

	outboxRepo := persistence.NewOutboxRepo(txm, log)
	boardsRepo := persistence.NewBoardsRepo(txm, log, outboxRepo)
	columnsRepo := persistence.NewColumnsRepo(txm, log, outboxRepo)
	tasksRepo := persistence.NewTasksRepo(txm, log, outboxRepo)

	boardsHandler := createBoardsHandler(log, boardsRepo, columnsRepo, tasksRepo, cache, cfg.RedisCacheTtl)
	columnsHandler := createColumnsHandler(log, columnsRepo, tasksRepo, cache)
	tasksHandler := createTasksHandler(log, tasksRepo, columnsRepo, cache)

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
		a.closeAll()
		return nil
	case err := <-errCh:
		return err
	}
}

func (a *App) Log() *zerolog.Logger { return a.log }

func (a *App) closeAll() {
	if a.redis != nil {
		_ = a.redis.Close()
	}
	a.kafkaProducer.Close()
	a.pg.Close()
	a.grpc.Stop()
}

func newRedis(cfg *config.Config) (*redis.Client, error) {
	if cfg.RedisAddr == "" {
		// Redis is optional; caching will be effectively disabled.
		return nil, ErrDisabled
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		_ = rdb.Close()
		return nil, err
	}

	return rdb, nil
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

func createBoardsHandler(
	log *zerolog.Logger,
	boardsRepo boarddo.Repository,
	columnsRepo columndo.Repository,
	tasksRepo taskdo.Repository,
	cache commonuc.Cacher,
	redisCacheTTL time.Duration,
) *grpc.BoardsHandler {
	createBoard := boarduc.NewCreateBoardUseCase(boardsRepo)
	getBoard := boarduc.NewGetBoardUseCase(boardsRepo, columnsRepo, tasksRepo, cache, redisCacheTTL)
	listBoards := boarduc.NewListBoardsUseCase(boardsRepo, columnsRepo, tasksRepo, cache, redisCacheTTL)
	updateBoard := boarduc.NewUpdateBoardUseCase(boardsRepo, cache)
	deleteBoard := boarduc.NewDeleteBoardUseCase(boardsRepo, cache)

	boardsHandler := grpc.NewBoardsHandler(log, createBoard, getBoard, listBoards, updateBoard, deleteBoard)
	return boardsHandler
}

func createColumnsHandler(
	log *zerolog.Logger,
	columnsRepo columndo.Repository,
	tasksRepo taskdo.Repository,
	cache commonuc.Cacher,
) *grpc.ColumnsHandler {
	createColumn := columnuc.NewCreateColumnUseCase(columnsRepo, cache)
	getColumn := columnuc.NewGetColumnUseCase(columnsRepo, tasksRepo)
	moveColumn := columnuc.NewMoveColumnUseCase(columnsRepo, cache)
	deleteColumn := columnuc.NewDeleteColumnUseCase(columnsRepo, cache)

	columnsHandler := grpc.NewColumnsHandler(log, createColumn, getColumn, moveColumn, deleteColumn)
	return columnsHandler
}

func createTasksHandler(
	log *zerolog.Logger,
	tasksRepo taskdo.Repository,
	columnsRepo columndo.Repository,
	cache commonuc.Cacher,
) *grpc.TasksHandler {
	createTask := taskuc.NewCreateTaskUseCase(tasksRepo, columnsRepo, cache)
	getTask := taskuc.NewGetTaskUseCase(tasksRepo)
	updateTask := taskuc.NewUpdateTaskUseCase(tasksRepo, columnsRepo, cache)
	moveTask := taskuc.NewMoveTaskUseCase(tasksRepo, columnsRepo, cache)
	deleteTask := taskuc.NewDeleteTaskUseCase(tasksRepo, columnsRepo, cache)

	tasksHandler := grpc.NewTasksHandler(log, createTask, getTask, updateTask, moveTask, deleteTask)
	return tasksHandler
}
