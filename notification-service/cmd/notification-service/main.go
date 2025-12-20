package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/smarrog/task-board/notification-service/internal/config"
	"github.com/smarrog/task-board/notification-service/internal/kafka"
	"github.com/smarrog/task-board/notification-service/internal/processor"
	"github.com/smarrog/task-board/shared/logger"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()
	log := logger.New(cfg.AppName, cfg.LogLevel)
	proc := processor.NewProcessor(log)
	consumer := kafka.NewConsumer(cfg, log, proc.Handle)

	if err := consumer.Start(ctx); err != nil {
		log.Fatal().Err(err).Msg("Failed to start")
	}

	log.Info().Msg("Started")

	<-ctx.Done()

	log.Info().Msg("Stopped")
}
