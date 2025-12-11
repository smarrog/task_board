package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func New(globalLevel zerolog.Level) zerolog.Logger {
	zerolog.SetGlobalLevel(globalLevel)

	return zerolog.
		New(os.Stdout).
		With().
		Timestamp().
		Caller().
		Logger()
}
