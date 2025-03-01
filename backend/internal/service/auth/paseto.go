package auth

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/nanayaw/fullstack/internal/config"
	"github.com/nanayaw/fullstack/internal/errors"
	"github.com/nanayaw/fullstack/internal/models"
	"github.com/nanayaw/fullstack/internal/service"
	"github.com/o1egl/paseto/v2"
)

type TokenClaims struct {
	ID        string    `json:"jti"`
	Subject   string    `json:"sub"`
	IssuedAt  time.Time `json:"iat"`
	ExpiresAt time.Time `json:"exp"`
	Type      string    `json:"type"`
}

type PasetoService struct {
	publicKey  ed25519.PublicKey
	privateKey ed25519.PrivateKey
	config     *config.AuthConfig
	userSvc    service.UserService
	emailSvc   service.EmailService
	cacheSvc   service.CacheService
}

func NewPasetoService(
	cfg *config.AuthConfig,
	userSvc service.UserService,
	emailSvc service.EmailService,
	cacheSvc service.CacheService,
) (*PasetoService, error) {
	// Decode keys from hex
	privateKeyBytes, err := hex.DecodeString(cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %w", err)
	}

	publicKeyBytes, err := hex.DecodeString(cfg.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}

	return &PasetoService{
		publicKey:  publicKeyBytes,
		privateKey: privateKeyBytes,
		config:     cfg,
		userSvc:    userSvc,
		emailSvc:   emailSvc,
		cacheSvc:   cacheSvc,
	}, nil
}

func (s *PasetoService) Register(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	// Create user
	user, err := s.userSvc.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	// Send verification email
	if err := s.SendVerificationEmail(ctx, user.ID); err != nil {
		// Log error but don't fail registration
		fmt.Printf("failed to send verification email: %v\n", err)
	}

	// Send welcome email
	if err := s.emailSvc.SendWelcomeEmail(ctx, user.Email, user.FullName); err != nil {
		// Log error but don't fail registration
		fmt.Printf("failed to send welcome email: %v\n", err)
	}

	return user, nil
}

func (s *PasetoService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	// Check rate limiting
	key := fmt.Sprintf("login_attempts:%s", req.Email)
	allowed, err := s.cacheSvc.CheckRateLimit(ctx, key, s.config.MaxLoginAttempts, int(s.config.LockoutDuration.Seconds()))
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, errors.NewRateLimitError("too many login attempts")
	}

	// Get user
	user, err := s.userSvc.GetUser(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	// TODO: Verify password hash

	// Generate tokens
	accessToken, err := s.generateToken(user.ID, "access", s.config.AccessTokenTTL)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateToken(user.ID, "refresh", s.config.RefreshTokenTTL)
	if err != nil {
		return nil, err
	}

	// Store refresh token
	if err := s.cacheSvc.StoreSession(ctx, refreshToken, user.ID, s.config.RefreshTokenTTL); err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *PasetoService) RefreshToken(ctx context.Context, req *models.RefreshTokenRequest) (*models.RefreshTokenResponse, error) {
	// Validate refresh token
	claims, err := s.validateToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	if claims.Type != "refresh" {
		return nil, errors.NewAuthenticationError("invalid token type")
	}

	// Check if token is still valid in cache
	userID, err := s.cacheSvc.GetSession(ctx, req.RefreshToken)
	if err != nil || userID == "" {
		return nil, errors.NewAuthenticationError("invalid or expired token")
	}

	// Generate new tokens
	accessToken, err := s.generateToken(claims.Subject, "access", s.config.AccessTokenTTL)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateToken(claims.Subject, "refresh", s.config.RefreshTokenTTL)
	if err != nil {
		return nil, err
	}

	// Invalidate old refresh token
	if err := s.cacheSvc.InvalidateSession(ctx, req.RefreshToken); err != nil {
		return nil, err
	}

	// Store new refresh token
	if err := s.cacheSvc.StoreSession(ctx, refreshToken, claims.Subject, s.config.RefreshTokenTTL); err != nil {
		return nil, err
	}

	return &models.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *PasetoService) Logout(ctx context.Context, sessionID string) error {
	return s.cacheSvc.InvalidateSession(ctx, sessionID)
}

func (s *PasetoService) ValidateSession(ctx context.Context, token string) (*models.Session, error) {
	claims, err := s.validateToken(token)
	if err != nil {
		return nil, err
	}

	if claims.Type != "access" {
		return nil, errors.NewAuthenticationError("invalid token type")
	}

	return &models.Session{
		ID:     claims.ID,
		UserID: claims.Subject,
	}, nil
}

func (s *PasetoService) InvalidateAllSessions(ctx context.Context, userID string) error {
	// Invalidate all sessions for user
	pattern := fmt.Sprintf("session:*:%s", userID)
	return s.cacheSvc.InvalidateCachePattern(ctx, pattern)
}

func (s *PasetoService) SendVerificationEmail(ctx context.Context, userID string) error {
	// Get user
	user, err := s.userSvc.GetUser(ctx, userID)
	if err != nil {
		return err
	}

	// Generate verification token
	token, err := s.generateToken(userID, "verification", s.config.VerificationTTL)
	if err != nil {
		return err
	}

	// Send email
	if err := s.emailSvc.SendVerificationEmail(ctx, user.Email, token); err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}

// Helper functions
func (s *PasetoService) generateToken(subject, tokenType string, expiration time.Duration) (string, error) {
	now := time.Now()
	claims := TokenClaims{
		ID:        generateUUID(),
		Subject:   subject,
		IssuedAt:  now,
		ExpiresAt: now.Add(expiration),
		Type:      tokenType,
	}

	token, err := paseto.NewV2().Sign(s.privateKey, claims, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func (s *PasetoService) validateToken(token string) (*TokenClaims, error) {
	var claims TokenClaims
	err := paseto.NewV2().Verify(token, s.publicKey, &claims, nil)
	if err != nil {
		return nil, errors.NewAuthenticationError("invalid token")
	}

	if time.Now().After(claims.ExpiresAt) {
		return nil, errors.NewAuthenticationError("token expired")
	}

	return &claims, nil
}

func generateUUID() string {
	// TODO: Implement proper UUID generation
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// Implement missing methods required by the auth.Service interface

func (s *PasetoService) VerifyEmail(ctx context.Context, req *models.VerifyEmailRequest) error {
	// Validate verification token
	claims, err := s.validateToken(req.Token)
	if err != nil {
		return err
	}

	if claims.Type != "verification" {
		return errors.NewAuthenticationError("invalid token type")
	}

	// Update user's email verification status
	updateReq := &models.UpdateUserRequest{
		// Set email as verified in the database
	}
	_, err = s.userSvc.UpdateUser(ctx, claims.Subject, updateReq)
	if err != nil {
		return err
	}

	return nil
}

func (s *PasetoService) SendPasswordResetEmail(ctx context.Context, email string) error {
	// Get user
	user, err := s.userSvc.GetUser(ctx, email)
	if err != nil {
		// Don't reveal if email exists or not
		return nil
	}

	// Generate reset token
	token, err := s.generateToken(user.ID, "password_reset", s.config.PasswordResetTTL)
	if err != nil {
		return err
	}

	// Send email
	if err := s.emailSvc.SendPasswordResetEmail(ctx, user.Email, token); err != nil {
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	return nil
}

func (s *PasetoService) ResetPassword(ctx context.Context, req *models.ResetPasswordRequest) error {
	// Validate reset token
	claims, err := s.validateToken(req.Token)
	if err != nil {
		return err
	}

	if claims.Type != "password_reset" {
		return errors.NewAuthenticationError("invalid token type")
	}

	// Update user's password
	updateReq := &models.UpdateUserRequest{
		Password: &req.Password,
	}
	_, err = s.userSvc.UpdateUser(ctx, claims.Subject, updateReq)
	if err != nil {
		return err
	}

	// Invalidate all sessions for user
	if err := s.InvalidateAllSessions(ctx, claims.Subject); err != nil {
		// Log error but don't fail password reset
		fmt.Printf("failed to invalidate sessions: %v\n", err)
	}

	// Send password changed notification
	user, err := s.userSvc.GetUser(ctx, claims.Subject)
	if err == nil {
		if err := s.emailSvc.SendPasswordChangedEmail(ctx, user.Email); err != nil {
			// Log error but don't fail password reset
			fmt.Printf("failed to send password changed notification: %v\n", err)
		}
	}

	return nil
}

func (s *PasetoService) ChangePassword(ctx context.Context, userID string, oldPassword, newPassword string) error {
	// Get user to verify old password
	user, err := s.userSvc.GetUser(ctx, userID)
	if err != nil {
		return err
	}

	// TODO: Verify old password hash against stored hash
	// This would require a password verification utility

	// Update password
	updateReq := &models.UpdateUserRequest{
		Password: &newPassword,
	}
	_, err = s.userSvc.UpdateUser(ctx, userID, updateReq)
	if err != nil {
		return err
	}

	// Invalidate all sessions for user
	if err := s.InvalidateAllSessions(ctx, userID); err != nil {
		// Log error but don't fail password change
		fmt.Printf("failed to invalidate sessions: %v\n", err)
	}

	// Send password changed notification
	if err := s.emailSvc.SendPasswordChangedEmail(ctx, user.Email); err != nil {
		// Log error but don't fail password change
		fmt.Printf("failed to send password changed notification: %v\n", err)
	}

	return nil
}

func (s *PasetoService) HandleOAuthLogin(ctx context.Context, req *models.OAuthLoginRequest) (*models.LoginResponse, error) {
	// TODO: Implement OAuth login
	// This would require fetching user profile from OAuth provider using the code

	// For now, create a mock user
	email := fmt.Sprintf("%s_user@example.com", req.Provider)
	fullName := fmt.Sprintf("%s User", req.Provider)

	// Check if user exists with this email
	user, err := s.userSvc.GetUser(ctx, email)
	if err != nil {
		// If user doesn't exist, create a new one
		createReq := &models.CreateUserRequest{
			Email:     email,
			FullName:  fullName,
			Password:  generateUUID(), // Generate random password
			AvatarURL: "",
		}
		user, err = s.Register(ctx, createReq)
		if err != nil {
			return nil, err
		}
	}

	// Generate tokens
	accessToken, err := s.generateToken(user.ID, "access", s.config.AccessTokenTTL)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateToken(user.ID, "refresh", s.config.RefreshTokenTTL)
	if err != nil {
		return nil, err
	}

	// Store refresh token
	if err := s.cacheSvc.StoreSession(ctx, refreshToken, user.ID, s.config.RefreshTokenTTL); err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *PasetoService) LinkOAuthAccount(ctx context.Context, userID string, req *models.OAuthLoginRequest) error {
	// TODO: Implement linking OAuth account
	// This would require storing OAuth provider and ID in the database
	return nil
}

func (s *PasetoService) UnlinkOAuthAccount(ctx context.Context, userID string, provider string) error {
	// TODO: Implement unlinking OAuth account
	// This would require removing OAuth provider and ID from the database
	return nil
}
