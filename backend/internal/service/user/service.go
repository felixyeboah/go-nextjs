package user

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/nanayaw/fullstack/internal/models"
)

// Repository defines the interface for user repository
type Repository interface {
	GetUser(ctx context.Context, id string) (*models.User, error)
	UpdateUser(ctx context.Context, id string, user *models.UpdateUserRequest) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
	GetUserActivity(ctx context.Context, userID string, page, pageSize int) ([]*models.AuditLog, int, error)
}

// Service handles user-related business logic
type Service struct {
	repo Repository
}

// NewService creates a new user service
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// GetUser retrieves a user by ID
func (s *Service) GetUser(c echo.Context, id string) (*models.User, error) {
	if s.repo == nil {
		// Return mock data for now
		return &models.User{
			ID:            id,
			Email:         "user@example.com",
			FullName:      "John Doe",
			AvatarURL:     "https://example.com/avatar.jpg",
			EmailVerified: true,
		}, nil
	}

	return s.repo.GetUser(c.Request().Context(), id)
}

// UpdateUser updates a user
func (s *Service) UpdateUser(c echo.Context, id string, req *models.UpdateUserRequest) (*models.User, error) {
	if s.repo == nil {
		// Return mock data for now
		return &models.User{
			ID:            id,
			Email:         "user@example.com",
			FullName:      *req.FullName,
			AvatarURL:     "https://example.com/avatar.jpg",
			EmailVerified: true,
		}, nil
	}

	return s.repo.UpdateUser(c.Request().Context(), id, req)
}

// DeleteUser deletes a user
func (s *Service) DeleteUser(c echo.Context, id string) error {
	if s.repo == nil {
		// Return success for now
		return nil
	}

	return s.repo.DeleteUser(c.Request().Context(), id)
}

// GetUserActivity retrieves user activity
func (s *Service) GetUserActivity(c echo.Context, userID string, page, pageSize int) ([]*models.AuditLog, int, error) {
	if s.repo == nil {
		// Return mock data for now
		activities := []*models.AuditLog{
			{
				ID:         "1",
				UserID:     userID,
				Action:     "login",
				EntityType: "session",
				EntityID:   "session-1",
				Metadata:   "Login from Chrome on Windows",
			},
		}
		return activities, 1, nil
	}

	return s.repo.GetUserActivity(c.Request().Context(), userID, page, pageSize)
}

// UpdateProfile updates a user's profile
func (s *Service) UpdateProfile(c echo.Context, userID string, req *models.UpdateUserRequest) (*models.User, error) {
	return s.UpdateUser(c, userID, req)
}
