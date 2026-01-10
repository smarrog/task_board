package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func ErrorHandler(log *zerolog.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		msg := "internal_error"

		var fe *fiber.Error
		if errors.As(err, &fe) {
			code = fe.Code
			msg = fe.Message
		}

		if code >= 500 {
			log.Error().Err(err).
				Str("method", c.Method()).
				Str("path", c.Path()).
				Int("status", code).
				Msg("http_error")
		}

		return c.Status(code).JSON(fiber.Map{"error": msg})
	}
}
