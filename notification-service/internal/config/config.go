package config

import (
	"github.com/rs/zerolog"
	"github.com/smarrog/task-board/shared/env"
	"github.com/smarrog/task-board/shared/logger"
	"time"
)

type Config struct {
	AppName  string
	LogLevel zerolog.Level

	PostgresDSN             string
	PostgresTimeout         time.Duration
	PostgresMinConns        int
	PostgresMaxConns        int
	PostgresMaxConnIdleTime time.Duration
	PostgresMaxConnLifeTime time.Duration

	KafkaGroupId    string
	KafkaBrokers    []string
	KafkaTopics     []string
	KafkaDLQEnabled bool
	KafkaDlqTopic   string
}

func Load() *Config {
	cfg := Config{
		AppName: env.GetString("APP_NAME", "notification-service"),

		PostgresDSN:             env.GetString("POSTGRES_DSN", ""),
		PostgresTimeout:         env.GetDuration("POSTGRES_TIMEOUT", time.Second*5),
		PostgresMinConns:        env.GetInt("POSTGRES_MIN_CONNS", 1),
		PostgresMaxConns:        env.GetInt("POSTGRES_MAX_CONNS", 10),
		PostgresMaxConnIdleTime: env.GetDuration("POSTGRES_MAX_CONN_IDLE_TIME", time.Minute*3),
		PostgresMaxConnLifeTime: env.GetDuration("POSTGRES_MAX_CONN_LIFETIME", time.Minute*30),

		KafkaGroupId:    env.GetString("KAFKA_GROUP_ID", ""),
		KafkaBrokers:    env.GetSplitString("KAFKA_BROKERS", []string{}),
		KafkaTopics:     env.GetSplitString("KAFKA_TOPICS", []string{}),
		KafkaDLQEnabled: env.GetBool("KAFKA_DLQ_ENABLED", true),
		KafkaDlqTopic:   env.GetString("KAFKA_DLQ_TOPIC", ""),

		LogLevel: logger.StrToLogLevel(env.GetString("LOG_LEVEL", "info")),
	}

	return &cfg
}
