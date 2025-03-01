package errors

import (
	"fmt"
	"net/http"
)

// AppError represents an application-specific error
type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
	Err        error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// Error constructors
func NewValidationError(message string) *AppError {
	return &AppError{
		Code:       "VALIDATION_ERROR",
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

func NewAuthenticationError(message string) *AppError {
	return &AppError{
		Code:       "AUTHENTICATION_ERROR",
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

func NewAuthorizationError(message string) *AppError {
	return &AppError{
		Code:       "AUTHORIZATION_ERROR",
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{
		Code:       "NOT_FOUND",
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

func NewConflictError(message string) *AppError {
	return &AppError{
		Code:       "CONFLICT",
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

func NewInternalError(err error) *AppError {
	return &AppError{
		Code:       "INTERNAL_ERROR",
		Message:    "An internal error occurred",
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}

func NewRateLimitError(message string) *AppError {
	return &AppError{
		Code:       "RATE_LIMIT_EXCEEDED",
		Message:    message,
		StatusCode: http.StatusTooManyRequests,
	}
}

func NewBadRequestError(message string) *AppError {
	return &AppError{
		Code:       "BAD_REQUEST",
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

// Common error messages
const (
	ErrInvalidCredentials = "Invalid email or password"
	ErrEmailAlreadyExists = "Email already exists"
	ErrUserNotFound       = "User not found"
	ErrSessionExpired     = "Session has expired"
	// #nosec G101 - This is an error message, not a hardcoded credential
	ErrInvalidToken       = "Invalid or expired token"
	ErrPasswordTooWeak    = "Password does not meet security requirements"
	ErrInvalidEmailFormat = "Invalid email format"
	ErrUnauthorized       = "Unauthorized access"
	ErrForbidden          = "Forbidden access"
	ErrInvalidRequest     = "Invalid request"
	ErrTooManyRequests    = "Too many requests"
	ErrInternalServer     = "Internal server error"
	ErrServiceUnavailable = "Service temporarily unavailable"
	ErrDatabaseConnection = "Database connection error"
	ErrCacheConnection    = "Cache connection error"
	ErrEmailSending       = "Error sending email"
	ErrOAuthProvider      = "Error with OAuth provider"
)
