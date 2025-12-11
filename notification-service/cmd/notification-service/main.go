package notification_service

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/smarrog/notification-service/internal/config"
	"github.com/smarrog/notification-service/internal/logger"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()
	log := logger.New(cfg.LogLevel)

	log.Info().Msg("Notification service was started")

	<-ctx.Done()

	log.Info().Msg("Notification service was stopped")
}
