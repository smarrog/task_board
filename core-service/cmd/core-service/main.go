package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/smarrog/task-board/core-service/internal/app"
	"github.com/smarrog/task-board/core-service/internal/config"
	"github.com/smarrog/task-board/shared/logger"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()
	log := logger.New(cfg.AppName, cfg.LogLevel)

	a, err := app.New(cfg, log)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing app")
	}

	if err := a.Run(ctx); err != nil {
		log.Fatal().Err(err).Msg("Error running app")
	}
}
