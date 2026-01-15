package config

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/smarrog/task-board/shared/env"
	"github.com/smarrog/task-board/shared/logger"
)

type Config struct {
	AppName string

	GRPCPort string

	PostgresDSN             string
	PostgresMinConns        int
	PostgresMaxConns        int
	PostgresMaxConnIdleTime time.Duration
	PostgresMaxConnLifeTime time.Duration
	PostgresTimeout         time.Duration

	RedisAddr     string
	RedisPassword string
	RedisDB       int
	RedisCacheTtl time.Duration

	KafkaGroupId string
	KafkaBrokers string
	KafkaAcks    int
	KafkaTopic   []string

	OutboxPollInterval time.Duration
	OutboxBatchSize    int

	LogLevel zerolog.Level
}

func Load() *Config {
	cfg := &Config{
		AppName: env.GetString("APP_NAME", "core-service"),

		GRPCPort: env.GetString("GRPC_PORT", "50052"),

		PostgresDSN:             env.GetString("POSTGRES_DSN", ""),
		PostgresMinConns:        env.GetInt("POSTGRES_MIN_CONNS", 1),
		PostgresMaxConns:        env.GetInt("POSTGRES_MAX_CONNS", 10),
		PostgresMaxConnIdleTime: env.GetDuration("POSTGRES_MAX_CONN_IDLE_TIME", 5*time.Minute),
		PostgresMaxConnLifeTime: env.GetDuration("POSTGRES_MAX_CONN_LIFE_TIME", 30*time.Minute),
		PostgresTimeout:         env.GetDuration("POSTGRES_TIMEOUT", 2*time.Second),

		RedisAddr:     env.GetString("REDIS_ADDR", ""),
		RedisPassword: env.GetString("REDIS_PASSWORD", ""),
		RedisDB:       env.GetInt("REDIS_DB", 0),
		RedisCacheTtl: env.GetDuration("REDIS_CACHE_TTL", 30*time.Second),

		KafkaGroupId: env.GetString("KAFKA_GROUP_ID", ""),
		KafkaBrokers: env.GetString("KAFKA_BROKERS", ""),
		KafkaAcks:    env.GetInt("KAFKA_ACKS", -1),
		KafkaTopic:   env.GetSplitString("KAFKA_TOPICS", []string{}),

		OutboxPollInterval: env.GetDuration("OUTBOX_POLL_INTERVAL", 5000*time.Millisecond),
		OutboxBatchSize:    env.GetInt("OUTBOX_BATCH_SIZE", 50),

		LogLevel: logger.StrToLogLevel(env.GetString("LOG_LEVEL", "info")),
	}

	return cfg
}
