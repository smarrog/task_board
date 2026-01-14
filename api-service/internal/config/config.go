package config

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/smarrog/task-board/shared/env"
	"github.com/smarrog/task-board/shared/logger"
)

type Config struct {
	LogLevel zerolog.Level

	HTTPAddr        string
	CoreGRPCAddr    string
	AuthGRPCAddr    string
	JWTSecret       string
	RequestTimeout  time.Duration
	ShutdownTimeout time.Duration

	FiberAppName      string
	FiberIdleTimeout  time.Duration
	FiberReadTimeout  time.Duration
	FiberWriteTimeout time.Duration
}

func Load() *Config {
	cfg := &Config{
		LogLevel: logger.StrToLogLevel(env.GetString("LOG_LEVEL", "info")),

		HTTPAddr:        env.GetString("HTTP_ADDR", ":8080"),
		CoreGRPCAddr:    env.GetString("CORE_GRPC_ADDR", "core-service:50051"),
		AuthGRPCAddr:    env.GetString("AUTH_GRPC_ADDR", "auth-service:50052"),
		JWTSecret:       env.GetString("JWT_SECRET", "dev-secret"),
		RequestTimeout:  env.GetDuration("REQUEST_TIMEOUT", 5*time.Second),
		ShutdownTimeout: env.GetDuration("SHUTDOWN_TIMEOUT", 10*time.Second),

		FiberAppName:      env.GetString("FIBER_APP_NAME", "task-board-api-gateway"),
		FiberIdleTimeout:  env.GetDuration("FIBER_IDLE_TIMEOUT", 30*time.Second),
		FiberReadTimeout:  env.GetDuration("FIBER_READ_TIMEOUT", 30*time.Second),
		FiberWriteTimeout: env.GetDuration("FIBER_WRITE_TIMEOUT", 30*time.Second),
	}

	return cfg
}
