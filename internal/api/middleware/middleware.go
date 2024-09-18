package middleware

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/rs/zerolog/log"
)

// Add a new function to set up all middlewares
func SetupMiddlewares(app *fiber.App) {
	app.Use(Recover())
	app.Use(Logger())
	app.Use(SecureHeaders())
	app.Use(CORSMiddleware())
	app.Use(Compress())
	app.Use(RateLimiter())
	app.Use(RequestID())
}

func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		stop := time.Now()

		requestID := c.GetRespHeader("X-Request-ID")
		if requestID == "" {
			requestID = "unknown"
		}

		log.Info().
			Str("request_id", requestID).
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
		c.Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com; script-src 'self' 'unsafe-inline' 'unsafe-eval'; img-src 'self' data:")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		return c.Next()
	}
}

func CORSMiddleware() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept",
	})
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

func RateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
	})
}

func RequestID() fiber.Handler {
	return requestid.New()
}
