package notifier

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/smarrog/task-board/notification-service/internal/domain/notification"
)

type LoggerNotifier struct {
	log *zerolog.Logger
}

func NewLoggerNotifier(log *zerolog.Logger) *LoggerNotifier {
	return &LoggerNotifier{log: log}
}

func (n *LoggerNotifier) Notify(ctx context.Context, notif notification.Notification) error {
	n.log.Info().Msg(notif.Text)
	return nil
}
