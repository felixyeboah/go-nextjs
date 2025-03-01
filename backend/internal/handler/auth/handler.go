package auth

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nanayaw/fullstack/internal/handler/response"
	"github.com/nanayaw/fullstack/internal/models"
	"github.com/nanayaw/fullstack/internal/service/auth"
)

// Handler handles authentication-related requests
type Handler struct {
	authService auth.Service
}

// NewHandler creates a new auth handler
func NewHandler(authService auth.Service) *Handler {
	return &Handler{
		authService: authService,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} RegisterResponse "User registered successfully"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /api/v1/auth/register [post]
func (h *Handler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request format"))
	}

	// Validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
	}

	// Convert to service model
	createUserReq := &models.CreateUserRequest{
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FirstName + " " + req.LastName,
	}

	// Call service
	user, err := h.authService.Register(c.Request().Context(), createUserReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err.Error()))
	}

	// Create response
	resp := RegisterResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	return c.JSON(http.StatusCreated, resp)
}

// Login godoc
// @Summary Login a user
// @Description Authenticate a user and return tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse "Login successful"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /api/v1/auth/login [post]
func (h *Handler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request format"))
	}

	// Validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
	}

	// Convert to service model
	loginReq := &models.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	// Call service
	result, err := h.authService.Login(c.Request().Context(), loginReq)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, response.NewErrorResponse("Invalid credentials"))
	}

	// Create response
	resp := LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    3600, // 1 hour in seconds
	}

	return c.JSON(http.StatusOK, resp)
}

// RefreshToken godoc
// @Summary Refresh authentication tokens
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} RefreshTokenResponse "Token refreshed successfully"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 401 {object} ErrorResponse "Invalid token"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /api/v1/auth/refresh [post]
func (h *Handler) RefreshToken(c echo.Context) error {
	var req RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request format"))
	}

	// Validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
	}

	// Convert to service model
	refreshReq := &models.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	}

	// Call service
	result, err := h.authService.RefreshToken(c.Request().Context(), refreshReq)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, response.NewErrorResponse("Invalid refresh token"))
	}

	// Create response
	resp := RefreshTokenResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    3600, // 1 hour in seconds
	}

	return c.JSON(http.StatusOK, resp)
}

// Logout godoc
// @Summary Logout a user
// @Description Invalidate the user's tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LogoutRequest true "Logout request"
// @Success 200 {object} LogoutResponse "Logout successful"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /api/v1/auth/logout [post]
func (h *Handler) Logout(c echo.Context) error {
	var req LogoutRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request format"))
	}

	// Validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
	}

	// Call service
	err := h.authService.Logout(c.Request().Context(), req.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to logout"))
	}

	// Create response
	resp := LogoutResponse{
		Message: "Logout successful",
	}

	return c.JSON(http.StatusOK, resp)
}

// VerifyEmail godoc
// @Summary Verify email
// @Description Verify a user's email address using a token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body VerifyEmailRequest true "Email verification token"
// @Success 200 {object} VerifyEmailResponse
// @Failure 400 {object} ErrorResponse "Invalid token"
// @Router /api/v1/auth/verify-email [post]
func (h *Handler) VerifyEmail(c echo.Context) error {
	var req VerifyEmailRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request"))
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
	}

	// Convert to service model
	verifyReq := &models.VerifyEmailRequest{
		Token: req.Token,
	}

	if err := h.authService.VerifyEmail(c.Request().Context(), verifyReq); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid verification token"))
	}

	return c.JSON(http.StatusOK, VerifyEmailResponse{
		Message: "Email verified successfully",
	})
}

// ForgotPassword godoc
// @Summary Request password reset
// @Description Send a password reset email to the user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ForgotPasswordRequest true "Email address"
// @Success 200 {object} ForgotPasswordResponse "Password reset email sent"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /api/v1/auth/forgot-password [post]
func (h *Handler) ForgotPassword(c echo.Context) error {
	var req ForgotPasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request format"))
	}

	// Validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
	}

	// Call service
	err := h.authService.SendPasswordResetEmail(c.Request().Context(), req.Email)
	if err != nil {
		// Don't expose whether the email exists or not for security reasons
		// Just return success even if the email doesn't exist
		log.Printf("Error sending password reset email: %v", err)
	}

	// Create response
	resp := ForgotPasswordResponse{
		Message: "If your email is registered, you will receive a password reset link",
	}

	return c.JSON(http.StatusOK, resp)
}

// ResetPassword godoc
// @Summary Reset user password
// @Description Reset a user's password using a token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "Password reset details"
// @Success 200 {object} ResetPasswordResponse "Password reset successful"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 401 {object} ErrorResponse "Invalid token"
// @Failure 500 {object} ErrorResponse "Server error"
// @Router /api/v1/auth/reset-password [post]
func (h *Handler) ResetPassword(c echo.Context) error {
	var req ResetPasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request format"))
	}

	// Validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse(err.Error()))
	}

	// Convert to service model
	resetReq := &models.ResetPasswordRequest{
		Token:    req.Token,
		Password: req.NewPassword,
	}

	// Call service
	err := h.authService.ResetPassword(c.Request().Context(), resetReq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid reset token"))
	}

	// Create response
	resp := ResetPasswordResponse{
		Message: "Password reset successful",
	}

	return c.JSON(http.StatusOK, resp)
}

// RegisterRoutes registers all auth routes
func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.POST("/register", h.Register)
	g.POST("/login", h.Login)
	g.POST("/refresh", h.RefreshToken)
	g.POST("/logout", h.Logout)
	g.POST("/verify-email", h.VerifyEmail)
	g.POST("/forgot-password", h.ForgotPassword)
	g.POST("/reset-password", h.ResetPassword)
}
