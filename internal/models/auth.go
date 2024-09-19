package models

type LoginRequest struct {
	Name string `json:"name" validate:"required"`
	Pass string `json:"pass" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegisterRequest struct {
	Name  string `json:"name" validate:"required,min=3,max=50"`
	Pass  string `json:"pass" validate:"required,min=8"`
	Email string `json:"email" validate:"required,email"`
}

type RegisterResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}

type RefreshTokenResponse struct {
	Token string `json:"token"`
}
