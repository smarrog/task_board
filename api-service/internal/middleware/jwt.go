package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const (
	LocalUserID = "user_id"
)

func JWT(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" || !strings.HasPrefix(strings.ToLower(auth), "bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing_or_invalid_authorization_header",
			})
		}

		raw := strings.TrimSpace(auth[len("Bearer "):])

		tok, err := jwt.Parse(raw, func(t *jwt.Token) (any, error) {
			if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "unexpected_jwt_alg")
			}
			return []byte(secret), nil
		})
		if err != nil || tok == nil || !tok.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid_token",
			})
		}

		claims, ok := tok.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid_token_claims",
			})
		}

		sub, _ := claims["sub"].(string)
		if sub == "" {
			sub, _ = claims["user_id"].(string)
		}
		if sub == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing_subject_claim",
			})
		}

		c.Locals(LocalUserID, sub)
		return c.Next()
	}
}
