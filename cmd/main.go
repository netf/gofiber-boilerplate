package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/netf/gofiber-boilerplate/config"
	_ "github.com/netf/gofiber-boilerplate/docs"
	"github.com/netf/gofiber-boilerplate/internal/db"
	"github.com/netf/gofiber-boilerplate/internal/handlers"
	"github.com/netf/gofiber-boilerplate/internal/middleware"
	"github.com/netf/gofiber-boilerplate/internal/monitoring"

	"github.com/gofiber/fiber/v2"
	swagger "github.com/gofiber/swagger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// @title           Golang Fiber Boilerplate API
// @version         1.0
// @description     This is a sample server Todo server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @host      localhost:8080
// @BasePath  /api/v1

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Setup logging
	setupLogging(cfg.LogLevel)

	// Initialize Sentry (optional)
	if cfg.SentryDSN != "" {
		monitoring.InitSentry(cfg.SentryDSN)
		defer monitoring.FlushSentry()
	}

	// Initialize database
	database, err := db.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not initialize database")
	}

	// Run migrations
	if err := db.Migrate(database); err != nil {
		log.Fatal().Err(err).Msg("Could not run migrations")
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// Setup middlewares
	app.Use(middleware.Recover())
	app.Use(middleware.Logger())
	app.Use(middleware.SecureHeaders())
	app.Use(middleware.CORSMiddleware())
	app.Use(middleware.Compress())
	app.Use(middleware.RateLimiter())
	app.Use(middleware.RequestID())

	// API versioning
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Register routes with dependency injection
	handlers.RegisterRoutes(v1, database, cfg)

	// Swagger documentation route
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:         "/swagger/doc.json",
		DeepLinking: false,
	}))

	// Health check endpoint
	// app.Get("/health", handlers.HealthCheck)

	// Graceful shutdown
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit
		log.Info().Msg("Gracefully shutting down...")
		if err := app.Shutdown(); err != nil {
			log.Error().Err(err).Msg("Server forced to shutdown")
		}
		log.Info().Msg("Server exiting")
	}()

	// Start server
	log.Info().Msgf("Starting server on %s", cfg.ServerAddress)
	if err := app.Listen(cfg.ServerAddress); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}

	log.Info().Msg("Running cleanup tasks...")
}

func setupLogging(level string) {
	// Set global log level
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		log.Warn().Msgf("Invalid log level '%s', defaulting to 'info'", level)
		logLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(logLevel)

	// Set logging output to console
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
