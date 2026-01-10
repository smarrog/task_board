package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func AccessLog(log *zerolog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		evt := log.Info()
		if err != nil {
			evt = log.Error().Err(err)
		}
		evt.
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", c.Response().StatusCode()).
			Dur("duration", time.Since(start)).
			Str("request_id", c.Get(RequestIDHeader)).
			Msg("http")
		return err
	}
}
