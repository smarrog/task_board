package logger

import (
	"github.com/rs/zerolog"
)

var defaultLogLevel = zerolog.InfoLevel
var strToLogLevelMap = map[string]zerolog.Level{
	"debug": zerolog.DebugLevel,
	"info":  zerolog.InfoLevel,
	"warn":  zerolog.WarnLevel,
	"error": zerolog.ErrorLevel,
	"fatal": zerolog.FatalLevel,
}

var logLevelToStringMap = map[zerolog.Level]string{
	zerolog.DebugLevel: "debug",
	zerolog.InfoLevel:  "info",
	zerolog.WarnLevel:  "warn",
	zerolog.ErrorLevel: "error",
	zerolog.FatalLevel: "fatal",
}

func StrToLogLevel(s string) zerolog.Level {
	if value, ok := strToLogLevelMap[s]; ok {
		return value
	}
	return defaultLogLevel
}

func LogLevelToString(logLevel zerolog.Level) string {
	if value, ok := logLevelToStringMap[logLevel]; ok {
		return value
	}
	return LogLevelToString(zerolog.InfoLevel)
}
