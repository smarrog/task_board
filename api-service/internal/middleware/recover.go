package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func Recover(log *zerolog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Error().Interface("panic", r).Str("path", c.Path()).Msg("panic recovered")
				_ = c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "internal_server_error",
				})
				err = fmt.Errorf("panic: %v", r)
			}
		}()
		return c.Next()
	}
}
