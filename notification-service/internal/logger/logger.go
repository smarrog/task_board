package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func New(globalLevel zerolog.Level) *zerolog.Logger {
	zerolog.SetGlobalLevel(globalLevel)

	logger := zerolog.
		New(os.Stdout).
		With().
		Timestamp().
		Caller().
		Str("service", "notification-service").
		Logger()

	return &logger
}
