package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/netf/gofiber-boilerplate/internal/api/auth"
	"github.com/netf/gofiber-boilerplate/internal/errors"
	"github.com/netf/gofiber-boilerplate/internal/models"
	"github.com/netf/gofiber-boilerplate/internal/services"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var _ services.AuthService = (*MockAuthService)(nil)

// Update MockAuthService to implement AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) AuthenticateUser(name, pass string) (*models.User, error) {
	args := m.Called(name, pass)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthService) RegisterUser(name, pass, email string) (*models.User, error) {
	args := m.Called(name, pass, email)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthService) GenerateToken(user *models.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) ValidateToken(token string) (*models.User, error) {
	args := m.Called(token)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthService) GetUserByID(id uint) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

func TestLogin(t *testing.T) {
	mockService := new(MockAuthService)
	validate := validator.New()
	handler := &AuthHandler{
		authService: mockService,
		validate:    validate,
	}

	app := fiber.New()
	app.Post("/auth/login", handler.Login)

	testCases := []struct {
		name           string
		loginRequest   models.LoginRequest
		mockUser       *models.User
		mockError      error
		expectedStatus int
	}{
		{
			name:           "Success",
			loginRequest:   models.LoginRequest{Name: "testuser", Pass: "password123"},
			mockUser:       &models.User{ID: 1, Name: "testuser"},
			mockError:      nil,
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "Invalid Credentials",
			loginRequest:   models.LoginRequest{Name: "testuser", Pass: "wrongpassword"},
			mockUser:       nil,
			mockError:      errors.New("invalid credentials"),
			expectedStatus: fiber.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService.On("AuthenticateUser", tc.loginRequest.Name, tc.loginRequest.Pass).Return(tc.mockUser, tc.mockError)

			body, _ := json.Marshal(tc.loginRequest)
			req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req)

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
			mockService.AssertExpectations(t)
		})
	}
}

func TestRegister(t *testing.T) {
	mockService := new(MockAuthService)
	validate := validator.New()
	handler := &AuthHandler{
		authService: mockService,
		validate:    validate,
	}

	app := fiber.New()
	app.Post("/auth/register", handler.Register)

	testCases := []struct {
		name            string
		registerRequest models.RegisterRequest
		mockUser        *models.User
		mockError       error
		expectedStatus  int
	}{
		{
			name:            "Success",
			registerRequest: models.RegisterRequest{Name: "newuser", Pass: "password123", Email: "newuser@example.com"},
			mockUser:        &models.User{ID: 1, Name: "newuser", Email: "newuser@example.com"},
			mockError:       nil,
			expectedStatus:  fiber.StatusCreated,
		},
		{
			name:            "User Already Exists",
			registerRequest: models.RegisterRequest{Name: "existinguser", Pass: "password123", Email: "existing@example.com"},
			mockUser:        nil,
			mockError:       errors.ErrUserAlreadyExists,
			expectedStatus:  fiber.StatusConflict,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService.On("RegisterUser", tc.registerRequest.Name, tc.registerRequest.Pass, tc.registerRequest.Email).Return(tc.mockUser, tc.mockError)

			body, _ := json.Marshal(tc.registerRequest)
			req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req)

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
			mockService.AssertExpectations(t)
		})
	}
}

func TestLogout(t *testing.T) {
	mockService := new(MockAuthService)
	validate := validator.New()
	handler := &AuthHandler{
		authService: mockService,
		validate:    validate,
	}

	app := fiber.New()
	app.Post("/auth/logout", handler.Logout)

	req := httptest.NewRequest("POST", "/auth/logout", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestRefreshToken(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	app := fiber.New()
	app.Post("/auth/refresh", func(c *fiber.Ctx) error {
		// Set a valid user in the context
		c.Locals("user", &models.User{ID: 1, Name: "testuser"})
		return handler.RefreshToken(c)
	})

	req := httptest.NewRequest("POST", "/auth/refresh", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

}

func TestLoginInvalidBody(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	app := fiber.New()
	app.Post("/auth/login", handler.Login)

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestRegisterInvalidBody(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	app := fiber.New()
	app.Post("/auth/register", handler.Register)

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestRegisterInvalidData(t *testing.T) {
	mockService := new(MockAuthService)
	validate := validator.New()
	handler := &AuthHandler{
		authService: mockService,
		validate:    validate,
	}

	app := fiber.New()
	app.Post("/auth/register", handler.Register)

	invalidRequest := models.RegisterRequest{Name: "a", Pass: "short", Email: "invalid-email"}
	body, _ := json.Marshal(invalidRequest)

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestMain(m *testing.M) {
	// Initialize auth package
	if err := auth.InitAuth(); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize auth")
	}

	// Run the tests
	os.Exit(m.Run())
}
