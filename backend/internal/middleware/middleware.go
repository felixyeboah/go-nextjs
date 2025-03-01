package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/nanayaw/fullstack/internal/errors"
	"github.com/nanayaw/fullstack/internal/service"
)

type Middleware struct {
	authService  service.AuthService
	cacheService service.CacheService
}

// Define a custom type for context keys to avoid collisions
type contextKey string

// Define constants for context keys
const (
	requestIDKey contextKey = "request_id"
)

func NewMiddleware(authService service.AuthService, cacheService service.CacheService) *Middleware {
	return &Middleware{
		authService:  authService,
		cacheService: cacheService,
	}
}

// Authenticate verifies the JWT token and sets the user in context
func (m *Middleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get("Authorization")
		if auth == "" {
			return errors.NewAuthenticationError(errors.ErrUnauthorized)
		}

		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return errors.NewAuthenticationError(errors.ErrUnauthorized)
		}

		session, err := m.authService.ValidateSession(c.Request().Context(), parts[1])
		if err != nil {
			return errors.NewAuthenticationError(errors.ErrInvalidToken)
		}

		// Set user ID in context
		c.Set("userID", session.UserID)
		c.Set("sessionID", session.ID)

		return next(c)
	}
}

// RateLimit implements rate limiting per IP address
func (m *Middleware) RateLimit(limit int, duration time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			key := fmt.Sprintf("ratelimit:%s:%s", ip, c.Request().URL.Path)

			// Check rate limit
			allowed, err := m.cacheService.CheckRateLimit(c.Request().Context(), key, limit, int(duration.Seconds()))
			if err != nil {
				return errors.NewInternalError(err)
			}

			if !allowed {
				return errors.NewRateLimitError(errors.ErrTooManyRequests)
			}

			return next(c)
		}
	}
}

// RequireRole checks if the authenticated user has the required role
func (m *Middleware) RequireRole(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get user from context (set by Authenticate middleware)
			userID := c.Get("userID").(string)
			if userID == "" {
				return errors.NewAuthenticationError(errors.ErrUnauthorized)
			}

			// TODO: Implement role checking logic
			// This would typically involve checking the user's role in the database
			// or in the JWT claims

			return next(c)
		}
	}
}

// RequestLogger logs request details
func RequestLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		err := next(c)

		stop := time.Now()
		latency := stop.Sub(start)

		req := c.Request()
		res := c.Response()

		fields := map[string]interface{}{
			"remote_ip":  c.RealIP(),
			"latency":    latency.String(),
			"method":     req.Method,
			"uri":        req.RequestURI,
			"status":     res.Status,
			"size":       res.Size,
			"user_agent": req.UserAgent(),
		}

		if err != nil {
			fields["error"] = err.Error()
		}

		// TODO: Use proper logger
		fmt.Printf("%+v\n", fields)

		return err
	}
}

// Recover recovers from panics and logs the error
func Recover(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}

				// Log the error and stack trace
				// TODO: Use proper logger
				fmt.Printf("PANIC: %v\n", err)

				// Return a 500 error to the client
				c.Error(errors.NewInternalError(err))
			}
		}()

		return next(c)
	}
}

// RequestID middleware adds a unique request ID to each request
func (m *Middleware) RequestID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		res := c.Response()

		rid := req.Header.Get(echo.HeaderXRequestID)
		if rid == "" {
			rid = uuid.New().String()
			req.Header.Set(echo.HeaderXRequestID, rid)
		}
		res.Header().Set(echo.HeaderXRequestID, rid)

		// Use the custom type for context key
		ctx := context.WithValue(req.Context(), requestIDKey, rid)
		c.SetRequest(req.WithContext(ctx))

		return next(c)
	}
}

// CORS handles Cross-Origin Resource Sharing
func CORS() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			c.Response().Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Accept,Authorization,Content-Type,X-Request-ID")
			c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
			c.Response().Header().Set("Access-Control-Max-Age", "86400")

			if c.Request().Method == http.MethodOptions {
				return c.NoContent(http.StatusNoContent)
			}

			return next(c)
		}
	}
}

// Timeout adds a timeout to the request context
func Timeout(timeout time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, cancel := context.WithTimeout(c.Request().Context(), timeout)
			defer cancel()

			c.SetRequest(c.Request().WithContext(ctx))

			done := make(chan error, 1)
			go func() {
				done <- next(c)
			}()

			select {
			case <-ctx.Done():
				return errors.NewInternalError(ctx.Err())
			case err := <-done:
				return err
			}
		}
	}
}
