package config

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	consts "github.com/smarrog/notification-service/internal/app"
)

type Config struct {
	LogLevel     zerolog.Level
	KafkaBrokers []string
	KafkaTopics  []string
}

func Load() Config {
	cfg := Config{
		LogLevel:     strToLogLevel(getEnv("LOG_LEVEL", "info")),
		KafkaBrokers: splitStringToSlice(getEnv("KAFKA_BROKERS", "kafka:9092")),
		KafkaTopics:  splitStringToSlice(getEnv("KAFKA_TOPICS", consts.KAFKA_TOPIC_EVENTS)),
	}

	return cfg
}

func splitStringToSlice(s string) []string {
	return strings.Split(s, ",")
}

func getEnv(key, def string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return def
}

func strToLogLevel(s string) zerolog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return zerolog.DebugLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}
