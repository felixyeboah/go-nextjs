package auth

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"testing"
	"time"

	"github.com/nanayaw/fullstack/internal/config"
	"github.com/nanayaw/fullstack/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Define UserActivity type if it's not defined in models
type UserActivity struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Action    string    `json:"action"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

// Define UpdateProfileRequest type if it's not defined in models
type UpdateProfileRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
	Avatar   string `json:"avatar"`
}

// Mock services
type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockUserService) GetUser(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// Add the missing DeleteUser method
func (m *mockUserService) DeleteUser(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Update the GetUserActivity method signature
func (m *mockUserService) GetUserActivity(ctx context.Context, userID string, page, pageSize int) ([]*models.AuditLog, int, error) {
	args := m.Called(ctx, userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*models.AuditLog), args.Int(1), args.Error(2)
}

// Add the missing ListUsers method
func (m *mockUserService) ListUsers(ctx context.Context, page, pageSize int) ([]*models.User, int, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*models.User), args.Int(1), args.Error(2)
}

// Add the missing RemoveAvatar method
func (m *mockUserService) RemoveAvatar(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// Add the missing UpdateProfile method
func (m *mockUserService) UpdateProfile(ctx context.Context, userID string, req *models.UpdateUserRequest) (*models.User, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// Add the missing UpdateUser method
func (m *mockUserService) UpdateUser(ctx context.Context, userID string, req *models.UpdateUserRequest) (*models.User, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// Add the missing UploadAvatar method
func (m *mockUserService) UploadAvatar(ctx context.Context, userID string, fileData []byte, fileType string) (string, error) {
	args := m.Called(ctx, userID, fileData, fileType)
	return args.String(0), args.Error(1)
}

// Implement other UserService methods...

type mockEmailService struct {
	mock.Mock
}

func (m *mockEmailService) SendVerificationEmail(ctx context.Context, to, token string) error {
	args := m.Called(ctx, to, token)
	return args.Error(0)
}

func (m *mockEmailService) SendWelcomeEmail(ctx context.Context, to, userName string) error {
	args := m.Called(ctx, to, userName)
	return args.Error(0)
}

// Add the missing ParseTemplate method
func (m *mockEmailService) ParseTemplate(templateName string, data interface{}) (string, error) {
	args := m.Called(templateName, data)
	return args.String(0), args.Error(1)
}

// Add the missing SendLoginNotificationEmail method
func (m *mockEmailService) SendLoginNotificationEmail(ctx context.Context, to, location, device string) error {
	args := m.Called(ctx, to, location, device)
	return args.Error(0)
}

// Add the missing SendPasswordChangedEmail method
func (m *mockEmailService) SendPasswordChangedEmail(ctx context.Context, to string) error {
	args := m.Called(ctx, to)
	return args.Error(0)
}

// Add the missing SendPasswordResetEmail method
func (m *mockEmailService) SendPasswordResetEmail(ctx context.Context, to, token string) error {
	args := m.Called(ctx, to, token)
	return args.Error(0)
}

// Add the missing ValidateEmailAddress method
func (m *mockEmailService) ValidateEmailAddress(email string) bool {
	args := m.Called(email)
	return args.Bool(0)
}

// Implement other EmailService methods...

type mockCacheService struct {
	mock.Mock
}

func (m *mockCacheService) CheckRateLimit(ctx context.Context, key string, limit, duration int) (bool, error) {
	args := m.Called(ctx, key, limit, duration)
	return args.Bool(0), args.Error(1)
}

func (m *mockCacheService) StoreSession(ctx context.Context, sessionID, userID string, expiration time.Duration) error {
	args := m.Called(ctx, sessionID, userID, expiration)
	return args.Error(0)
}

// Fix the CacheData method to use int for expiration
func (m *mockCacheService) CacheData(ctx context.Context, key string, value interface{}, expiration int) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

// Add the missing GetCachedData method
func (m *mockCacheService) GetCachedData(ctx context.Context, key string, dest interface{}) error {
	args := m.Called(ctx, key, dest)
	return args.Error(0)
}

// Add the missing InvalidateCache method
func (m *mockCacheService) InvalidateCache(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

// Add the missing InvalidateCachePattern method
func (m *mockCacheService) InvalidateCachePattern(ctx context.Context, pattern string) error {
	args := m.Called(ctx, pattern)
	return args.Error(0)
}

// Add the missing ResetRateLimit method
func (m *mockCacheService) ResetRateLimit(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

// Add missing GetSession and InvalidateSession methods referenced in tests
func (m *mockCacheService) GetSession(ctx context.Context, sessionID string) (string, error) {
	args := m.Called(ctx, sessionID)
	return args.String(0), args.Error(1)
}

func (m *mockCacheService) InvalidateSession(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

// Implement other CacheService methods...

// Helper function to create a test config with keys
func createTestConfig() *config.AuthConfig {
	// Generate a test key pair
	publicKey, privateKey, _ := ed25519.GenerateKey(nil)

	return &config.AuthConfig{
		VerificationTTL:  time.Hour * 24,
		PasswordResetTTL: time.Hour * 24,
		AccessTokenTTL:   time.Minute * 15,
		RefreshTokenTTL:  time.Hour * 24 * 7,
		MaxLoginAttempts: 5,
		PrivateKey:       hex.EncodeToString(privateKey),
		PublicKey:        hex.EncodeToString(publicKey),
	}
}

func TestPasetoService_Register(t *testing.T) {
	// Setup
	userSvc := new(mockUserService)
	emailSvc := new(mockEmailService)
	cacheSvc := new(mockCacheService)

	cfg := createTestConfig()

	service, err := NewPasetoService(cfg, userSvc, emailSvc, cacheSvc)
	assert.NoError(t, err)

	// Test data
	req := &models.CreateUserRequest{
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
	}

	user := &models.User{
		ID:       "user123",
		Email:    req.Email,
		FullName: req.FullName,
	}

	// Setup expectations
	userSvc.On("CreateUser", mock.Anything, req).Return(user, nil)
	userSvc.On("GetUser", mock.Anything, user.ID).Return(user, nil)
	emailSvc.On("SendVerificationEmail", mock.Anything, user.Email, mock.AnythingOfType("string")).Return(nil)
	emailSvc.On("SendWelcomeEmail", mock.Anything, user.Email, user.FullName).Return(nil)

	// Execute
	result, err := service.Register(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, user, result)
	userSvc.AssertExpectations(t)
	emailSvc.AssertExpectations(t)
}

func TestPasetoService_Login(t *testing.T) {
	// Setup
	userSvc := new(mockUserService)
	emailSvc := new(mockEmailService)
	cacheSvc := new(mockCacheService)

	cfg := createTestConfig()

	service, err := NewPasetoService(cfg, userSvc, emailSvc, cacheSvc)
	assert.NoError(t, err)

	// Test data
	req := &models.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	user := &models.User{
		ID:       "user123",
		Email:    req.Email,
		FullName: "Test User",
	}

	// Setup expectations
	cacheSvc.On("CheckRateLimit", mock.Anything, mock.AnythingOfType("string"), cfg.MaxLoginAttempts, int(cfg.LockoutDuration.Seconds())).Return(true, nil)
	userSvc.On("GetUser", mock.Anything, req.Email).Return(user, nil)
	cacheSvc.On("StoreSession", mock.Anything, mock.AnythingOfType("string"), user.ID, cfg.RefreshTokenTTL).Return(nil)

	// Execute
	result, err := service.Login(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user, &result.User)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
	userSvc.AssertExpectations(t)
	cacheSvc.AssertExpectations(t)
}

func TestPasetoService_RefreshToken(t *testing.T) {
	// Setup
	userSvc := new(mockUserService)
	emailSvc := new(mockEmailService)
	cacheSvc := new(mockCacheService)

	cfg := createTestConfig()

	service, err := NewPasetoService(cfg, userSvc, emailSvc, cacheSvc)
	assert.NoError(t, err)

	// Test data
	userID := "user123"
	oldRefreshToken, err := service.generateToken(userID, "refresh", cfg.RefreshTokenTTL)
	assert.NoError(t, err)

	req := &models.RefreshTokenRequest{
		RefreshToken: oldRefreshToken,
	}

	// Setup expectations
	cacheSvc.On("GetSession", mock.Anything, oldRefreshToken).Return(userID, nil)
	cacheSvc.On("InvalidateSession", mock.Anything, oldRefreshToken).Return(nil)
	cacheSvc.On("StoreSession", mock.Anything, mock.AnythingOfType("string"), userID, cfg.RefreshTokenTTL).Return(nil)

	// Execute
	result, err := service.RefreshToken(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
	assert.NotEqual(t, oldRefreshToken, result.RefreshToken)
	cacheSvc.AssertExpectations(t)
}

func TestPasetoService_ValidateSession(t *testing.T) {
	// Setup
	userSvc := new(mockUserService)
	emailSvc := new(mockEmailService)
	cacheSvc := new(mockCacheService)

	cfg := createTestConfig()

	service, err := NewPasetoService(cfg, userSvc, emailSvc, cacheSvc)
	assert.NoError(t, err)

	// Test data
	userID := "user123"
	token, err := service.generateToken(userID, "access", cfg.AccessTokenTTL)
	assert.NoError(t, err)

	// Execute
	session, err := service.ValidateSession(context.Background(), token)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, userID, session.UserID)
}
