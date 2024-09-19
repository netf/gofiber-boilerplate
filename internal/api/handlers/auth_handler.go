package handlers

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/netf/gofiber-boilerplate/internal/api/auth"
	apiUtils "github.com/netf/gofiber-boilerplate/internal/api/utils"
	"github.com/netf/gofiber-boilerplate/internal/errors"
	"github.com/netf/gofiber-boilerplate/internal/models"
	"github.com/netf/gofiber-boilerplate/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
	validate    *validator.Validate
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService, validate: validator.New()}
}

// Login handles user authentication and returns a JWT token
// @Summary User login
// @Description Authenticate a user and return a JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param login body models.LoginRequest true "Login credentials"
// @Success 200 {object} apiUtils.Response[models.LoginResponse]
// @Failure 400 {object} apiUtils.ErrorResponse
// @Failure 401 {object} apiUtils.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var login models.LoginRequest
	if err := c.BodyParser(&login); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(apiUtils.CreateErrorResponse("Invalid request body", fiber.StatusBadRequest))
	}

	user, err := h.authService.AuthenticateUser(login.Name, login.Pass)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(apiUtils.CreateErrorResponse("Invalid credentials", fiber.StatusUnauthorized))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES512, auth.Claims{
		StandardClaims: jwt.StandardClaims{
			Subject:   user.Name,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	})

	ts, err := token.SignedString(auth.PrivateKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(apiUtils.CreateErrorResponse("Could not generate token", fiber.StatusInternalServerError))
	}

	response := apiUtils.CreateResponse[models.LoginResponse](models.LoginResponse{Token: ts})
	return c.Status(fiber.StatusOK).JSON(response)
}

// Register handles user registration
// @Summary User registration
// @Description Register a new user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body models.RegisterRequest true "User registration details"
// @Success 201 {object} apiUtils.Response[models.RegisterResponse]
// @Failure 400 {object} apiUtils.ErrorResponse
// @Failure 409 {object} apiUtils.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var register models.RegisterRequest
	if err := c.BodyParser(&register); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(apiUtils.CreateErrorResponse("Invalid request body", fiber.StatusBadRequest))
	}

	// Add validation
	if err := h.validate.Struct(register); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(apiUtils.CreateErrorResponse(err.Error(), fiber.StatusBadRequest))
	}

	user, err := h.authService.RegisterUser(register.Name, register.Pass, register.Email)
	if err != nil {
		if errors.Is(err, errors.ErrUserAlreadyExists) {
			return c.Status(fiber.StatusConflict).JSON(apiUtils.CreateErrorResponse("User already exists", fiber.StatusConflict))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(apiUtils.CreateErrorResponse("Could not register user", fiber.StatusInternalServerError))
	}

	response := apiUtils.CreateResponse[models.RegisterResponse](models.RegisterResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
	return c.Status(fiber.StatusCreated).JSON(response)
}

// Logout handles user logout
// @Summary User logout
// @Description Logout a user (client-side only in this implementation)
// @Tags Authentication
// @Produce json
// @Success 200 {object} apiUtils.Response[models.LogoutResponse]
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	response := apiUtils.CreateResponse[models.LogoutResponse](models.LogoutResponse{
		Message: "Successfully logged out",
	})
	return c.Status(fiber.StatusOK).JSON(response)
}

// RefreshToken handles token refresh
// @Summary Refresh JWT token
// @Description Refresh the JWT token for a logged-in user
// @Tags Authentication
// @Produce json
// @Success 200 {object} apiUtils.Response[models.RefreshTokenResponse]
// @Failure 401 {object} apiUtils.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	token := jwt.NewWithClaims(jwt.SigningMethodES512, auth.Claims{
		StandardClaims: jwt.StandardClaims{
			Subject:   user.Name,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	})

	ts, err := token.SignedString(auth.PrivateKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(apiUtils.CreateErrorResponse("Could not refresh token", fiber.StatusInternalServerError))
	}

	response := apiUtils.CreateResponse[models.RefreshTokenResponse](models.RefreshTokenResponse{Token: ts})
	return c.Status(fiber.StatusOK).JSON(response)
}
