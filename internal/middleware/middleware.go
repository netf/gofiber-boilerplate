package middleware

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/rs/zerolog/log"
)

func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		stop := time.Now()

		log.Info().
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", c.Response().StatusCode()).
			Dur("latency", stop.Sub(start)).
			Msg("Handled request")

		return err
	}
}

func Recover() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("%v", r)
				log.Error().Err(err).Bytes("stack", debug.Stack()).Msg("Panic recovered")

				// Capture exception in Sentry (if initialized)
				sentry.CurrentHub().Recover(r)
				sentry.Flush(2 * time.Second)

				c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Internal Server Error",
				})
			}
		}()
		return c.Next()
	}
}

func SecureHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Security-Policy", "default-src 'self'")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		return c.Next()
	}
}

func Compress() fiber.Handler {
	return compress.New()
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	log.Error().Err(err).Msg("Unhandled error")

	// Capture exception in Sentry (if initialized)
	sentry.CaptureException(err)
	sentry.Flush(2 * time.Second)

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}

func InitSentry(dsn string) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: dsn,
	})
	if err != nil {
		log.Error().Err(err).Msg("Sentry initialization failed")
	} else {
		log.Info().Msg("Sentry initialized")
	}
}

func FlushSentry() {
	sentry.Flush(2 * time.Second)
}
