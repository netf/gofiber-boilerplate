package errors

// Custom error types
var (
	ErrUserAlreadyExists  = New("user already exists")
	ErrInvalidCredentials = New("invalid credentials")
	ErrResourceNotFound   = New("resource not found")
	ErrInvalidInput       = New("invalid input")
	ErrUnauthorized       = New("unauthorized")
	ErrInternalServer     = New("internal server error")
	ErrDatabaseOperation  = New("database operation error")
	ErrUserNotFound       = New("user not found")
)
