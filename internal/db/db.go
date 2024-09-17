package db

import (
	"github.com/netf/gofiber-boilerplate/internal/models"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(databaseURL string) (*gorm.DB, error) {
	// Initialize GORM
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to database")
		return nil, err
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	// Auto-migrate models
	err := db.AutoMigrate(&models.Todo{})
	if err != nil {
		log.Error().Err(err).Msg("Database migration failed")
		return err
	}

	log.Info().Msg("Database migration completed")
	return nil
}
