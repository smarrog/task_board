package config

import (
	"strings"

	"github.com/rs/zerolog"
	"github.com/smarrog/task-board/shared/env"
	"github.com/smarrog/task-board/shared/logger"
)

type Config struct {
	LogLevel     zerolog.Level
	KafkaGroupId string
	KafkaBrokers string
	KafkaTopics  []string
}

func Load() *Config {
	cfg := Config{
		LogLevel:     logger.StrToLogLevel(env.GetEnv("LOG_LEVEL", "info")),
		KafkaGroupId: env.GetEnv("KAFKA_GROUP_ID", "notification-service"),
		KafkaBrokers: env.GetEnv("KAFKA_BROKERS", "kafka:9092"),
		KafkaTopics:  strings.Split(env.GetEnv("KAFKA_TOPICS", "board-events"), ","),
	}

	return &cfg
}
