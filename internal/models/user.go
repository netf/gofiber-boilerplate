package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system.
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"uniqueIndex" json:"name" validate:"required,min=3,max=50"`
	Pass      []byte         `json:"-" validate:"required"`
	Email     string         `gorm:"uniqueIndex" json:"email" validate:"required,email"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
