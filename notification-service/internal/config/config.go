package config

import (
	"github.com/rs/zerolog"
	"github.com/smarrog/task-board/shared/env"
	"github.com/smarrog/task-board/shared/logger"
)

type Config struct {
	AppName      string
	LogLevel     zerolog.Level
	KafkaGroupId string
	KafkaBrokers string
	KafkaTopics  []string
}

func Load() *Config {
	cfg := Config{
		AppName: env.GetString("APP_NAME", "notification-service"),

		KafkaGroupId: env.GetString("KAFKA_GROUP_ID", "notification-service"),
		KafkaBrokers: env.GetString("KAFKA_BROKERS", "kafka:9092"),
		KafkaTopics:  env.GetSplitString("KAFKA_TOPICS", []string{"board-events"}),

		LogLevel: logger.StrToLogLevel(env.GetString("LOG_LEVEL", "info")),
	}

	return &cfg
}
