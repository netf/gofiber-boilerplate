package models

import (
	"time"

	"gorm.io/gorm"
)

// Todo represents a task to be done.
type Todo struct {
	// The ID of the todo item.
	// example: 1
	ID uint `gorm:"primaryKey" json:"id"`
	// The title of the todo item.
	// example: Buy groceries
	Title string `json:"title" validate:"required,min=3,max=255"`
	// The status of the todo item.
	// example: false
	Completed bool           `json:"completed"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
