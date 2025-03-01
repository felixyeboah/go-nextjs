package user

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/nanayaw/fullstack/internal/handler/response"
	"github.com/nanayaw/fullstack/internal/models"
	"github.com/nanayaw/fullstack/internal/service/auth"
)

// UserService defines the interface for user service
type UserService interface {
	GetUser(ctx echo.Context, id string) (*models.User, error)
	UpdateUser(ctx echo.Context, id string, req *models.UpdateUserRequest) (*models.User, error)
	DeleteUser(ctx echo.Context, id string) error
	GetUserActivity(ctx echo.Context, userID string, page, pageSize int) ([]*models.AuditLog, int, error)
	UpdateProfile(ctx echo.Context, userID string, req *models.UpdateUserRequest) (*models.User, error)
}

// Handler handles user-related requests
type Handler struct {
	userService UserService
	authService auth.Service
}

// NewHandler creates a new user handler
func NewHandler(userService UserService, authService auth.Service) *Handler {
	return &Handler{
		userService: userService,
		authService: authService,
	}
}

// GetUser handles getting user information
func (h *Handler) GetUser(c echo.Context) error {
	// Extract user ID from context (set by auth middleware)
	userID := c.Get("user_id").(string)

	// Call service
	user, err := h.userService.GetUser(c, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to get user information"))
	}

	return c.JSON(http.StatusOK, user)
}

// UpdateUser handles updating user information
func (h *Handler) UpdateUser(c echo.Context) error {
	// Extract user ID from context (set by auth middleware)
	userID := c.Get("user_id").(string)

	var req models.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request format"))
	}

	// Validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
	}

	// Call service
	user, err := h.userService.UpdateUser(c, userID, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update user information"))
	}

	return c.JSON(http.StatusOK, user)
}

// DeleteUser handles deleting a user
func (h *Handler) DeleteUser(c echo.Context) error {
	// Extract user ID from context (set by auth middleware)
	userID := c.Get("user_id").(string)

	// Call service
	err := h.userService.DeleteUser(c, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete user"))
	}

	return c.JSON(http.StatusOK, response.NewSuccessResponse("User deleted successfully", nil))
}

// GetUserActivity handles getting user activity
func (h *Handler) GetUserActivity(c echo.Context) error {
	// Extract user ID from context (set by auth middleware)
	userID := c.Get("user_id").(string)

	// Get pagination parameters
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Call service
	activities, total, err := h.userService.GetUserActivity(c, userID, page, pageSize)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to get user activity"))
	}

	// Convert to response model
	activityItems := make([]ActivityItem, len(activities))
	for i, activity := range activities {
		activityItems[i] = ActivityItem{
			ID:        activity.ID,
			Action:    activity.Action,
			Timestamp: activity.CreatedAt.Format(http.TimeFormat),
			Details:   activity.Metadata,
		}
	}

	// Calculate total pages
	totalPages := total / pageSize
	if total%pageSize > 0 {
		totalPages++
	}

	// Create response
	resp := UserActivityResponse{
		Activities: activityItems,
		Pagination: PaginationInfo{
			CurrentPage: page,
			PageSize:    pageSize,
			TotalItems:  total,
			TotalPages:  totalPages,
		},
	}

	return c.JSON(http.StatusOK, resp)
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get the current user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} UserProfileResponse "User profile"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /api/v1/users/profile [get]
func (h *Handler) GetProfile(c echo.Context) error {
	// Extract user ID from context (set by auth middleware)
	userID := c.Get("user_id").(string)

	// Call service
	user, err := h.userService.GetUser(c, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to get user profile"))
	}

	// Split full name into first and last name
	firstName, lastName := splitFullName(user.FullName)

	// Create response
	resp := UserProfileResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: firstName,
		LastName:  lastName,
		CreatedAt: user.CreatedAt.Format(http.TimeFormat),
		UpdatedAt: user.UpdatedAt.Format(http.TimeFormat),
	}

	return c.JSON(http.StatusOK, resp)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the current user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateProfileRequest true "Profile update details"
// @Success 200 {object} UserProfileResponse "Updated profile"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /api/v1/users/profile [put]
func (h *Handler) UpdateProfile(c echo.Context) error {
	// Extract user ID from context (set by auth middleware)
	userID := c.Get("user_id").(string)

	var req UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request format"))
	}

	// Validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
	}

	// Convert to service model
	fullName := req.FirstName + " " + req.LastName
	updateReq := &models.UpdateUserRequest{
		FullName: &fullName,
	}

	// Call service
	user, err := h.userService.UpdateProfile(c, userID, updateReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update profile"))
	}

	// Split full name into first and last name
	firstName, lastName := splitFullName(user.FullName)

	// Create response
	resp := UserProfileResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: firstName,
		LastName:  lastName,
		CreatedAt: user.CreatedAt.Format(http.TimeFormat),
		UpdatedAt: user.UpdatedAt.Format(http.TimeFormat),
	}

	return c.JSON(http.StatusOK, resp)
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change the current user's password
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangePasswordRequest true "Password change details"
// @Success 200 {object} ChangePasswordResponse "Password changed successfully"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /api/v1/users/change-password [post]
func (h *Handler) ChangePassword(c echo.Context) error {
	// Extract user ID from context (set by auth middleware)
	userID := c.Get("user_id").(string)

	var req ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request format"))
	}

	// Validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
	}

	// Call service
	err := h.authService.ChangePassword(c.Request().Context(), userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Failed to change password. Current password may be incorrect."))
	}

	// Create response
	resp := ChangePasswordResponse{
		Message: "Password changed successfully",
	}

	return c.JSON(http.StatusOK, resp)
}

// DeleteAccount godoc
// @Summary Delete user account
// @Description Delete the current user's account
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} DeleteAccountResponse "Account deleted successfully"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /api/v1/users/account [delete]
func (h *Handler) DeleteAccount(c echo.Context) error {
	// Extract user ID from context (set by auth middleware)
	userID := c.Get("user_id").(string)

	// Call service
	err := h.userService.DeleteUser(c, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete account"))
	}

	// Invalidate all sessions
	err = h.authService.InvalidateAllSessions(c.Request().Context(), userID)
	if err != nil {
		// Log error but continue
		// The user account is already deleted, so we should still return success
	}

	// Create response
	resp := DeleteAccountResponse{
		Message: "Account deleted successfully",
	}

	return c.JSON(http.StatusOK, resp)
}

// RegisterRoutes registers all user routes
func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.GET("/me", h.GetUser)
	g.PUT("/me", h.UpdateUser)
	g.DELETE("/me", h.DeleteUser)
	g.GET("/me/activity", h.GetUserActivity)
	g.GET("/profile", h.GetProfile)
	g.PUT("/profile", h.UpdateProfile)
	g.POST("/change-password", h.ChangePassword)
	g.DELETE("/account", h.DeleteAccount)
}

// Helper function to split full name into first and last name
func splitFullName(fullName string) (string, string) {
	// Split the full name by space
	parts := []rune(fullName)

	// Find the last space
	lastSpaceIndex := -1
	for i, r := range parts {
		if r == ' ' {
			lastSpaceIndex = i
		}
	}

	// If no space found, return the full name as first name and empty last name
	if lastSpaceIndex == -1 {
		return fullName, ""
	}

	// Split the name
	firstName := string(parts[:lastSpaceIndex])
	lastName := string(parts[lastSpaceIndex+1:])

	return firstName, lastName
}
