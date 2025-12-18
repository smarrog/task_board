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
	log := logger.New("notification-service", cfg.LogLevel)
	proc := processor.NewProcessor(log)
	consumer := kafka.NewConsumer(cfg, log, proc.Handle)

	consumer.Start(ctx)

	log.Info().Msg("Notification service was started")

	<-ctx.Done()

	log.Info().Msg("Notification service was stopped")
}
