package config

import (
	"github.com/rs/zerolog"
	"github.com/smarrog/task-board/shared/env"
	"github.com/smarrog/task-board/shared/logger"
)

type Config struct {
	AppName         string
	LogLevel        zerolog.Level
	KafkaGroupId    string
	KafkaBrokers    string
	KafkaTopics     []string
	KafkaDLQEnabled bool
	KafkaDLQTopic   string
}

func Load() *Config {
	cfg := Config{
		AppName: env.GetString("APP_NAME", "notification-service"),

		KafkaGroupId:    env.GetString("KAFKA_GROUP_ID", ""),
		KafkaBrokers:    env.GetString("KAFKA_BROKERS", ""),
		KafkaTopics:     env.GetSplitString("KAFKA_TOPICS", []string{}),
		KafkaDLQEnabled: env.GetBool("KAFKA_DLQ_ENABLED", true),
		KafkaDLQTopic:   env.GetString("KAFKA_DLQ_TOPIC", ""),

		LogLevel: logger.StrToLogLevel(env.GetString("LOG_LEVEL", "info")),
	}

	return &cfg
}
