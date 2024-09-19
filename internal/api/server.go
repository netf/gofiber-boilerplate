package api

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/netf/gofiber-boilerplate/config"
	_ "github.com/netf/gofiber-boilerplate/docs"
	"github.com/netf/gofiber-boilerplate/internal/api/auth"
	"github.com/netf/gofiber-boilerplate/internal/api/middleware"
	"github.com/netf/gofiber-boilerplate/internal/api/routes"
	"github.com/netf/gofiber-boilerplate/internal/db"
	"github.com/netf/gofiber-boilerplate/internal/monitoring"

	"github.com/gofiber/fiber/v2"
	swagger "github.com/gofiber/swagger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func Start(cfg *config.Config, database *gorm.DB) {
	setupLogging(cfg.LogLevel)

	if cfg.SentryDSN != "" {
		monitoring.InitSentry(cfg.SentryDSN)
		defer monitoring.FlushSentry()
	}

	// Add this block to run migrations
	if err := db.AutoMigrate(database); err != nil {
		log.Fatal().Err(err).Msg("Could not run auto-migrations")
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	middleware.SetupMiddlewares(app)

	if err := auth.InitAuth(); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize auth configuration")
	}

	api := app.Group("/api")
	v1 := api.Group("/v1")

	routes.RegisterRoutes(v1, database, cfg)

	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:         "/swagger/doc.json",
		DeepLinking: false,
	}))

	go gracefulShutdown(app)

	log.Info().Msgf("Starting server on %s", cfg.ServerAddress)
	if err := app.Listen(cfg.ServerAddress); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}

	log.Info().Msg("Running cleanup tasks...")
}

func gracefulShutdown(app *fiber.App) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Gracefully shutting down...")
	if err := app.Shutdown(); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	}
	log.Info().Msg("Server exiting")
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
