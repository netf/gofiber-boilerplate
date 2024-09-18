package auth

import (
	"github.com/gofiber/fiber/v2"
	jwt "github.com/gofiber/jwt/v3"
)

func JWTProtected(secret string) fiber.Handler {
	return jwt.New(jwt.Config{
		SigningKey:   []byte(secret),
		ContextKey:   "user",
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Unauthorized",
	})
}
