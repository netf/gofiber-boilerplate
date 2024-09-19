package db

import (
	"github.com/netf/gofiber-boilerplate/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(databaseURL string) (*gorm.DB, error) {
	// Initialize GORM
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(models.ModelsToMigrate...)
}
