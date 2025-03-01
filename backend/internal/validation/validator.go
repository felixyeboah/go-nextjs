package validation

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/nanayaw/fullstack/internal/errors"
)

var (
	validate *validator.Validate
	// Password requirements
	minPasswordLength = 8
	maxPasswordLength = 72 // bcrypt max length
	// Email regex pattern
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

func init() {
	validate = validator.New()
	// Register custom validation functions
	if err := validate.RegisterValidation("password", validatePassword); err != nil {
		panic(fmt.Sprintf("failed to register password validation: %v", err))
	}
	if err := validate.RegisterValidation("email", validateEmail); err != nil {
		panic(fmt.Sprintf("failed to register email validation: %v", err))
	}
}

// ValidateStruct validates a struct using validator tags
func ValidateStruct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			messages := make([]string, 0)
			for _, e := range validationErrors {
				messages = append(messages, formatValidationError(e))
			}
			return errors.NewValidationError(strings.Join(messages, "; "))
		}
		return errors.NewValidationError("Invalid input")
	}
	return nil
}

// Custom validation functions
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < minPasswordLength || len(password) > maxPasswordLength {
		return false
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

func validateEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	if !emailRegex.MatchString(email) {
		return false
	}

	// Additional RFC 5322 validation
	if _, err := mail.ParseAddress(email); err != nil {
		return false
	}

	return true
}

// Helper functions
func formatValidationError(e validator.FieldError) string {
	field := strings.ToLower(e.Field())
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "password":
		return fmt.Sprintf("%s must be at least %d characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character", field, minPasswordLength)
	case "min":
		return fmt.Sprintf("%s must be at least %s", field, e.Param())
	case "max":
		return fmt.Sprintf("%s must not exceed %s", field, e.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, e.Param())
	default:
		return fmt.Sprintf("%s failed validation: %s", field, e.Tag())
	}
}

// Convenience validation functions
func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.NewValidationError(errors.ErrInvalidEmailFormat)
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.NewValidationError(errors.ErrInvalidEmailFormat)
	}
	return nil
}

func ValidatePassword(password string) error {
	if !validatePassword(validator.FieldLevel(nil)) {
		return errors.NewValidationError(errors.ErrPasswordTooWeak)
	}
	return nil
}

func SanitizeString(s string) string {
	return strings.TrimSpace(s)
}

func SanitizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
