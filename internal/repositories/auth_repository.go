package repositories

import (
	"github.com/netf/gofiber-boilerplate/internal/errors"
	"github.com/netf/gofiber-boilerplate/internal/models"
	"gorm.io/gorm"
)

// AuthRepository handles database operations related to authentication
type AuthRepository struct {
	db *gorm.DB
}

// NewAuthRepository creates a new AuthRepository instance
func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// FindUserByCredentials retrieves a user by their name and hashed password
func (r *AuthRepository) FindUserByCredentials(name string, hashedPassword []byte) (*models.User, error) {
	var user models.User
	err := r.db.Select("id, name, email").
		Where("name = ? AND pass = ?", name, hashedPassword).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.ErrDatabaseOperation
	}

	return &user, nil
}

// FindUserByName retrieves a user by their username
func (r *AuthRepository) FindUserByName(name string) (*models.User, error) {
	var user models.User
	err := r.db.Where("name = ?", name).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.ErrDatabaseOperation
	}

	return &user, nil
}

// CreateUser creates a new user in the database
func (r *AuthRepository) CreateUser(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return errors.ErrDatabaseOperation
	}

	return nil
}
