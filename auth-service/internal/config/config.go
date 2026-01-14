package config

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/smarrog/task-board/shared/env"
	"github.com/smarrog/task-board/shared/logger"
)

type Config struct {
	LogLevel zerolog.Level

	AppName string
	GRPCPort string

	PostgresDSN string
	PostgresMinConns int
	PostgresMaxConns int
	PostgresTimeout time.Duration
	PostgresMaxConnIdleTime time.Duration
	PostgresMaxConnLifeTime time.Duration

	JWTSecret string
	AccessTokenTTL time.Duration
}

func Load() *Config {
	cfg := Config{
		AppName:  env.GetString("APP_NAME", "auth-service"),
		LogLevel: logger.StrToLogLevel(env.GetString("LOG_LEVEL", "info")),
		GRPCPort: env.GetString("GRPC_PORT", "50052"),

		PostgresDSN: env.GetString("POSTGRES_DSN", "postgres://postgres:postgres@auth-db:5432/auth?sslmode=disable"),
		PostgresMinConns: env.GetInt("POSTGRES_MIN_CONNS", 1),
		PostgresMaxConns: env.GetInt("POSTGRES_MAX_CONNS", 10),
		PostgresTimeout: env.GetDuration("POSTGRES_TIMEOUT", 5*time.Second),
		PostgresMaxConnIdleTime: env.GetDuration("POSTGRES_MAX_CONN_IDLE_TIME", 30*time.Second),
		PostgresMaxConnLifeTime: env.GetDuration("POSTGRES_MAX_CONN_LIFE_TIME", 5*time.Minute),

		JWTSecret: env.GetString("JWT_SECRET", "dev-secret"),
		AccessTokenTTL: env.GetDuration("ACCESS_TOKEN_TTL", 24*time.Hour),
	}

	return &cfg
}
