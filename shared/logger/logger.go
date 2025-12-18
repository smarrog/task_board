package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func New(service string, globalLevel zerolog.Level) *zerolog.Logger {
	zerolog.SetGlobalLevel(globalLevel)

	l := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Caller().
		Str("service", service).
		Logger()

	return &l
}
