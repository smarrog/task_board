package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/smarrog/task-board/api-service/internal/app"
	"github.com/smarrog/task-board/api-service/internal/config"
	"github.com/smarrog/task-board/shared/logger"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()
	log := logger.New("api-service", cfg.LogLevel)

	a := app.New(cfg, log)
	if err := a.Init(); err != nil {
		log.Fatal().Err(err).Msg("init failed")
	}

	if err := a.Run(ctx); err != nil {
		log.Fatal().Err(err).Msg("run failed")
	}
}
