package app

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/smarrog/task-board/notification-service/internal/config"
	"github.com/smarrog/task-board/notification-service/internal/handler"
	"github.com/smarrog/task-board/notification-service/internal/infrastructure/kafka"
	"github.com/smarrog/task-board/notification-service/internal/infrastructure/notifier"
	"github.com/smarrog/task-board/notification-service/internal/infrastructure/persistence"
	"github.com/smarrog/task-board/shared/logger"
)

var ErrInvalidConfig = errors.New("invalid config")

type App struct {
	log      *zerolog.Logger
	cfg      *config.Config
	pg       *pgxpool.Pool
	dlq      *kafka.DlqWriter
	consumer *kafka.Consumer
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

	nRepo := persistence.NewNotificationsRepo(pg, log)
	logNotifier := notifier.NewLoggerNotifier(log)
	h := handler.NewHandler(logNotifier, nRepo)

	var dlqWriter *kafka.DlqWriter
	if cfg.KafkaDLQEnabled {
		w, err := kafka.NewDlqWriter(log, cfg.KafkaBrokers, cfg.KafkaDlqTopic)
		if err != nil {
			return err
		}
		dlqWriter = w
	}
	a.dlq = dlqWriter

	msgHandler := kafka.NewOutboxHandler(log, h, dlqWriter)
	consumer := kafka.NewConsumer(cfg, log, msgHandler.HandleKafkaMessage)
	a.consumer = consumer

	return nil
}

func (a *App) Run(ctx context.Context) error {
	if a.consumer == nil {
		return ErrInvalidConfig
	}
	defer a.closeAll()

	a.log.Info().Msg("Started")
	err := a.consumer.Start(ctx)
	a.log.Info().Msg("Stopped")
	return err
}

func (a *App) Log() *zerolog.Logger { return a.log }

func (a *App) closeAll() {
	if a.dlq != nil {
		a.dlq.Close()
	}
	if a.pg != nil {
		a.pg.Close()
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
