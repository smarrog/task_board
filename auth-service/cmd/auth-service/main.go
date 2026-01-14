package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/smarrog/task-board/auth-service/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	a := app.App{}
	if err := a.Init(); err != nil {
		a.Log().Fatal().Err(err).Msg("Error initializing app")
	}

	if err := a.Run(ctx); err != nil {
		a.Log().Fatal().Err(err).Msg("Error running app")
	}
}
