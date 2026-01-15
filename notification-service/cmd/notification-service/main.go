package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/smarrog/task-board/notification-service/internal/app"
	"github.com/smarrog/task-board/notification-service/internal/config"
	"github.com/smarrog/task-board/notification-service/internal/infrastructure/kafka"
	"github.com/smarrog/task-board/notification-service/internal/infrastructure/notifier"
	"github.com/smarrog/task-board/shared/logger"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()
	log := logger.New(cfg.AppName, cfg.LogLevel)

	logNotifier := notifier.NewLoggerNotifier(log)
	a := app.NewHandler(logNotifier)

	var dlqWriter *kafka.DlqWriter
	if cfg.KafkaDLQEnabled {
		w, err := kafka.NewDlqWriter(log, cfg.KafkaBrokers, cfg.KafkaDlqTopic)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create DLQ producer")
		}
		dlqWriter = w
		defer dlqWriter.Close()
	}

	msgHandler := kafka.NewOutboxHandler(log, a, dlqWriter)

	consumer := kafka.NewConsumer(cfg, log, msgHandler.HandleKafkaMessage)

	if err := consumer.Start(ctx); err != nil {
		log.Fatal().Err(err).Msg("Failed to start")
	}

	log.Info().Msg("Started")

	<-ctx.Done()

	log.Info().Msg("Stopped")
}
