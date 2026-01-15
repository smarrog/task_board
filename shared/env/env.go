package env

import (
	"os"
	"strconv"
	"strings"
	"time"
)

func GetString(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}

func GetInt(key string, def int) int {
	v := GetString(key, "")
	if v == "" {
		return def
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}

	return n
}

func GetDuration(key string, def time.Duration) time.Duration {
	v := GetString(key, "")
	if v == "" {
		return def
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		return def
	}

	return d
}

func GetSplitString(key string, def []string) []string {
	v := GetString(key, "")
	if v == "" {
		return def
	}

	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func GetBool(key string, def bool) bool {
	v := strings.TrimSpace(strings.ToLower(GetString(key, "")))
	if v == "" {
		return def
	}

	switch v {
	case "1", "true", "t", "yes", "y", "on":
		return true
	case "0", "false", "f", "no", "n", "off":
		return false
	default:
		return def
	}
}
