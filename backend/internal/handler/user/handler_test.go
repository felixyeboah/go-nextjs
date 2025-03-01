package user

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

// MockUserService is a mock implementation of the user service
type MockUserService struct {
	mock.Mock
}

// GetUser mocks the GetUser method
func (m *MockUserService) GetUser(c echo.Context, id string) (*models.User, error) {
	args := m.Called(c, id)
	return args.Get(0).(*models.User), args.Error(1)
}

// UpdateUser mocks the UpdateUser method
func (m *MockUserService) UpdateUser(c echo.Context, id string, req *models.UpdateUserRequest) (*models.User, error) {
	args := m.Called(c, id, req)
	return args.Get(0).(*models.User), args.Error(1)
}

// DeleteUser mocks the DeleteUser method
func (m *MockUserService) DeleteUser(c echo.Context, id string) error {
	args := m.Called(c, id)
	return args.Error(0)
}

// GetUserActivity mocks the GetUserActivity method
func (m *MockUserService) GetUserActivity(c echo.Context, userID string, page, pageSize int) ([]*models.AuditLog, int, error) {
	args := m.Called(c, userID, page, pageSize)
	return args.Get(0).([]*models.AuditLog), args.Int(1), args.Error(2)
}

// UpdateProfile mocks the UpdateProfile method
func (m *MockUserService) UpdateProfile(c echo.Context, userID string, req *models.UpdateUserRequest) (*models.User, error) {
	args := m.Called(c, userID, req)
	return args.Get(0).(*models.User), args.Error(1)
}

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

// TestGetUser tests the GetUser handler
func TestGetUser(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create mock services
	mockUserService := new(MockUserService)
	mockAuthService := new(MockAuthService)

	// Create a new user handler with the mock services
	handler := NewHandler(mockUserService, mockAuthService)

	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set user ID in context
	c.Set("user_id", "123")

	// Set up expectations
	mockUser := &models.User{
		ID:        "123",
		Email:     "test@example.com",
		FullName:  "John Doe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockUserService.On("GetUser", mock.Anything, "123").Return(mockUser, nil)

	// Call the handler
	if assert.NoError(t, handler.GetUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse the response
		var resp models.User
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)

		// Check the response
		assert.Equal(t, "123", resp.ID)
		assert.Equal(t, "test@example.com", resp.Email)
		assert.Equal(t, "John Doe", resp.FullName)
	}

	// Verify expectations
	mockUserService.AssertExpectations(t)
}

// TestUpdateUser tests the UpdateUser handler
func TestUpdateUser(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create mock services
	mockUserService := new(MockUserService)
	mockAuthService := new(MockAuthService)

	// Create a new user handler with the mock services
	handler := NewHandler(mockUserService, mockAuthService)

	// Create a request body
	fullName := "Jane Doe"
	reqBody := models.UpdateUserRequest{
		FullName: &fullName,
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/me", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set user ID in context
	c.Set("user_id", "123")

	// Set up expectations
	mockUser := &models.User{
		ID:        "123",
		Email:     "test@example.com",
		FullName:  "Jane Doe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockUserService.On("UpdateUser", mock.Anything, "123", mock.MatchedBy(func(req *models.UpdateUserRequest) bool {
		return *req.FullName == "Jane Doe"
	})).Return(mockUser, nil)

	// Mock the validator
	e.Validator = &MockValidator{}

	// Call the handler
	if assert.NoError(t, handler.UpdateUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse the response
		var resp models.User
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)

		// Check the response
		assert.Equal(t, "123", resp.ID)
		assert.Equal(t, "test@example.com", resp.Email)
		assert.Equal(t, "Jane Doe", resp.FullName)
	}

	// Verify expectations
	mockUserService.AssertExpectations(t)
}

// MockValidator is a mock implementation of the validator
type MockValidator struct{}

// Validate mocks the Validate method
func (m *MockValidator) Validate(i interface{}) error {
	return nil
}
