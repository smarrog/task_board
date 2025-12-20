package config

import (
	"github.com/rs/zerolog"
	"github.com/smarrog/task-board/shared/env"
	"github.com/smarrog/task-board/shared/logger"
)

type Config struct {
	LogLevel zerolog.Level
}

func Load() *Config {
	cfg := Config{
		LogLevel: logger.StrToLogLevel(env.GetString("LOG_LEVEL", "info")),
	}

	return &cfg
}
