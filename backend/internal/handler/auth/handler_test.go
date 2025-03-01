package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nanayaw/fullstack/internal/models"
)

// MockAuthService is a mock implementation of the auth service
type MockAuthService struct {
	mock.Mock
}

// Register mocks the Register method
func (m *MockAuthService) Register(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*models.User), args.Error(1)
}

// Login mocks the Login method
func (m *MockAuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*models.LoginResponse), args.Error(1)
}

// RefreshToken mocks the RefreshToken method
func (m *MockAuthService) RefreshToken(ctx context.Context, req *models.RefreshTokenRequest) (*models.RefreshTokenResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*models.RefreshTokenResponse), args.Error(1)
}

// Logout mocks the Logout method
func (m *MockAuthService) Logout(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

// SendVerificationEmail mocks the SendVerificationEmail method
func (m *MockAuthService) SendVerificationEmail(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// VerifyEmail mocks the VerifyEmail method
func (m *MockAuthService) VerifyEmail(ctx context.Context, req *models.VerifyEmailRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

// SendPasswordResetEmail mocks the SendPasswordResetEmail method
func (m *MockAuthService) SendPasswordResetEmail(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

// ResetPassword mocks the ResetPassword method
func (m *MockAuthService) ResetPassword(ctx context.Context, req *models.ResetPasswordRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

// ChangePassword mocks the ChangePassword method
func (m *MockAuthService) ChangePassword(ctx context.Context, userID string, oldPassword, newPassword string) error {
	args := m.Called(ctx, userID, oldPassword, newPassword)
	return args.Error(0)
}

// HandleOAuthLogin mocks the HandleOAuthLogin method
func (m *MockAuthService) HandleOAuthLogin(ctx context.Context, req *models.OAuthLoginRequest) (*models.LoginResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*models.LoginResponse), args.Error(1)
}

// LinkOAuthAccount mocks the LinkOAuthAccount method
func (m *MockAuthService) LinkOAuthAccount(ctx context.Context, userID string, req *models.OAuthLoginRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

// UnlinkOAuthAccount mocks the UnlinkOAuthAccount method
func (m *MockAuthService) UnlinkOAuthAccount(ctx context.Context, userID string, provider string) error {
	args := m.Called(ctx, userID, provider)
	return args.Error(0)
}

// ValidateSession mocks the ValidateSession method
func (m *MockAuthService) ValidateSession(ctx context.Context, sessionID string) (*models.Session, error) {
	args := m.Called(ctx, sessionID)
	return args.Get(0).(*models.Session), args.Error(1)
}

// InvalidateAllSessions mocks the InvalidateAllSessions method
func (m *MockAuthService) InvalidateAllSessions(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// TestRegister tests the Register handler
func TestRegister(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a mock auth service
	mockAuthService := new(MockAuthService)

	// Create a new auth handler with the mock service
	handler := NewHandler(mockAuthService)

	// Create a request body
	reqBody := RegisterRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set up expectations
	mockUser := &models.User{
		ID:        "123",
		Email:     "test@example.com",
		FullName:  "John Doe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockAuthService.On("Register", mock.Anything, mock.MatchedBy(func(req *models.CreateUserRequest) bool {
		return req.Email == "test@example.com" && req.Password == "password123" && req.FullName == "John Doe"
	})).Return(mockUser, nil)

	// Mock the validator
	e.Validator = &MockValidator{}

	// Call the handler
	if assert.NoError(t, handler.Register(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		// Parse the response
		var resp RegisterResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)

		// Check the response
		assert.Equal(t, "123", resp.ID)
		assert.Equal(t, "test@example.com", resp.Email)
		assert.Equal(t, "John", resp.FirstName)
		assert.Equal(t, "Doe", resp.LastName)
	}

	// Verify expectations
	mockAuthService.AssertExpectations(t)
}

// TestLogin tests the Login handler
func TestLogin(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a mock auth service
	mockAuthService := new(MockAuthService)

	// Create a new auth handler with the mock service
	handler := NewHandler(mockAuthService)

	// Create a request body
	reqBody := LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set up expectations
	mockResponse := &models.LoginResponse{
		User:         models.User{ID: "123", Email: "test@example.com"},
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
	}
	mockAuthService.On("Login", mock.Anything, mock.MatchedBy(func(req *models.LoginRequest) bool {
		return req.Email == "test@example.com" && req.Password == "password123"
	})).Return(mockResponse, nil)

	// Mock the validator
	e.Validator = &MockValidator{}

	// Call the handler
	if assert.NoError(t, handler.Login(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse the response
		var resp LoginResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)

		// Check the response
		assert.Equal(t, "access-token", resp.AccessToken)
		assert.Equal(t, "refresh-token", resp.RefreshToken)
		assert.Equal(t, 3600, resp.ExpiresIn)
	}

	// Verify expectations
	mockAuthService.AssertExpectations(t)
}

// MockValidator is a mock implementation of the validator
type MockValidator struct{}

// Validate mocks the Validate method
func (m *MockValidator) Validate(i interface{}) error {
	return nil
}
