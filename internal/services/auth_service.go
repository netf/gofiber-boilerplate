package services

import (
	"github.com/netf/gofiber-boilerplate/internal/errors"
	"github.com/netf/gofiber-boilerplate/internal/models"
	"github.com/netf/gofiber-boilerplate/internal/repositories"
	"github.com/netf/gofiber-boilerplate/internal/utils"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// AuthService defines the interface for authentication-related operations
type AuthService interface {
	AuthenticateUser(name, password string) (*models.User, error)
	RegisterUser(name, password, email string) (*models.User, error)
}

// authService implements the AuthService interface
type authService struct {
	authRepo repositories.AuthRepository
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(authRepo repositories.AuthRepository) AuthService {
	return &authService{authRepo: authRepo}
}

// AuthenticateUser authenticates a user with the given credentials
func (s *authService) AuthenticateUser(name, password string) (*models.User, error) {
	if name == "" || password == "" {
		return nil, errors.New("username and password are required")
	}

	hashedPassword, err := utils.PwHash([]byte(password))
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		return nil, err
	}

	user, err := s.authRepo.FindUserByCredentials(name, hashedPassword)
	if err != nil {
		log.Error().Err(err).Str("name", name).Msg("Failed to find user")
		return nil, err
	}

	return user, nil
}

// RegisterUser creates a new user account
func (s *authService) RegisterUser(name, password, email string) (*models.User, error) {
	// Check if user already exists
	existingUser, err := s.authRepo.FindUserByName(name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.ErrUserAlreadyExists
	}

	// Hash the password
	hashedPassword, err := utils.PwHash([]byte(password))
	if err != nil {
		return nil, err
	}

	// Create new user
	newUser := &models.User{
		Name:  name,
		Pass:  hashedPassword,
		Email: email,
	}

	if err := s.authRepo.CreateUser(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}
