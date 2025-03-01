package router

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	authHandler "github.com/nanayaw/fullstack/internal/handler/auth"
	appMiddleware "github.com/nanayaw/fullstack/internal/handler/middleware"
	userHandler "github.com/nanayaw/fullstack/internal/handler/user"
	"github.com/nanayaw/fullstack/internal/service/auth"
)

// Router handles all the routes for the application
type Router struct {
	Echo        *echo.Echo
	AuthHandler *authHandler.Handler
	UserHandler *userHandler.Handler
	AuthService auth.Service
}

// NewRouter creates a new router
func NewRouter(e *echo.Echo, authHandler *authHandler.Handler, userHandler *userHandler.Handler, authService auth.Service) *Router {
	return &Router{
		Echo:        e,
		AuthHandler: authHandler,
		UserHandler: userHandler,
		AuthService: authService,
	}
}

// SetupRoutes sets up all the routes for the application
func (r *Router) SetupRoutes() {
	// Middleware
	r.Echo.Use(middleware.Logger())
	r.Echo.Use(middleware.Recover())
	r.Echo.Use(middleware.CORS())
	r.Echo.Use(middleware.RequestID())

	// API v1 group
	v1 := r.Echo.Group("/api/v1")

	// Health check endpoint
	v1.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "ok",
			"version": "1.0.0",
		})
	})

	// Auth routes
	auth := v1.Group("/auth")
	r.AuthHandler.RegisterRoutes(auth)

	// User routes
	users := v1.Group("/users")
	users.Use(appMiddleware.AuthMiddleware(r.AuthService))
	r.UserHandler.RegisterRoutes(users)

	// Swagger documentation
	r.Echo.GET("/swagger/*", echoSwagger.WrapHandler)
}

// SetupTimeoutMiddleware sets up the timeout middleware
func (r *Router) SetupTimeoutMiddleware(timeoutSeconds int) {
	r.Echo.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: time.Duration(timeoutSeconds) * time.Second,
	}))
}
