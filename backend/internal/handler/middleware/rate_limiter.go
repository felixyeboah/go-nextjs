package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nanayaw/fullstack/internal/handler/response"
)

// RateLimiterConfig contains configuration for the rate limiter
type RateLimiterConfig struct {
	// Max number of requests allowed in the time window
	Limit int
	// Time window for rate limiting
	Window time.Duration
	// Message to return when rate limit is exceeded
	Message string
	// Optional function to get a custom key for rate limiting (e.g., user ID, API key)
	// If nil, client IP will be used
	KeyFunc func(echo.Context) string
	// Routes to exclude from rate limiting
	ExcludedRoutes []string
}

// DefaultRateLimiterConfig returns a default configuration for the rate limiter
func DefaultRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		Limit:   100,
		Window:  time.Minute,
		Message: "Rate limit exceeded. Please try again later.",
		KeyFunc: func(c echo.Context) string {
			// Get client IP by default
			ip, _, err := net.SplitHostPort(c.Request().RemoteAddr)
			if err != nil {
				return c.Request().RemoteAddr
			}
			return ip
		},
		ExcludedRoutes: []string{},
	}
}

// AuthRateLimiterConfig returns a configuration for authentication endpoints
func AuthRateLimiterConfig() RateLimiterConfig {
	config := DefaultRateLimiterConfig()
	config.Limit = 5
	config.Window = time.Minute * 15
	config.Message = "Too many login attempts. Please try again later."
	return config
}

// client represents a client with request count and reset time
type client struct {
	count      int
	lastAccess time.Time
	resetTime  time.Time
}

// RateLimiter implements a simple in-memory rate limiter
type RateLimiter struct {
	mu      sync.RWMutex
	clients map[string]*client
	config  RateLimiterConfig
}

// NewRateLimiter creates a new rate limiter with the given configuration
func NewRateLimiter(config RateLimiterConfig) *RateLimiter {
	limiter := &RateLimiter{
		clients: make(map[string]*client),
		config:  config,
	}

	// Start a goroutine to clean up expired clients
	go limiter.cleanupLoop()

	return limiter
}

// cleanupLoop periodically removes expired clients
func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.config.Window)
	defer ticker.Stop()

	for range ticker.C {
		rl.cleanup()
	}
}

// cleanup removes expired clients
func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for key, c := range rl.clients {
		if now.After(c.resetTime) {
			delete(rl.clients, key)
		}
	}
}

// isExcluded checks if the current route is excluded from rate limiting
func (rl *RateLimiter) isExcluded(path string) bool {
	for _, route := range rl.config.ExcludedRoutes {
		if route == path {
			return true
		}
	}
	return false
}

// Middleware returns an Echo middleware function for rate limiting
func (rl *RateLimiter) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip rate limiting for excluded routes
			if rl.isExcluded(c.Path()) {
				return next(c)
			}

			// Get client key
			key := rl.config.KeyFunc(c)

			rl.mu.Lock()
			now := time.Now()
			cl, exists := rl.clients[key]

			// If client doesn't exist or window has expired, create a new entry
			if !exists || now.After(cl.resetTime) {
				rl.clients[key] = &client{
					count:      1,
					lastAccess: now,
					resetTime:  now.Add(rl.config.Window),
				}
				rl.mu.Unlock()
				return next(c)
			}

			// Update last access time
			cl.lastAccess = now

			// Check if rate limit is exceeded
			if cl.count >= rl.config.Limit {
				rl.mu.Unlock()

				// Set rate limit headers
				c.Response().Header().Set("X-RateLimit-Limit", string(rune(rl.config.Limit)))
				c.Response().Header().Set("X-RateLimit-Remaining", "0")
				c.Response().Header().Set("X-RateLimit-Reset", string(rune(cl.resetTime.Unix())))
				c.Response().Header().Set("Retry-After", string(rune(int(cl.resetTime.Sub(now).Seconds()))))

				return c.JSON(http.StatusTooManyRequests, response.ErrorResponse{
					Error: rl.config.Message,
				})
			}

			// Increment request count
			cl.count++
			remaining := rl.config.Limit - cl.count
			rl.mu.Unlock()

			// Set rate limit headers
			c.Response().Header().Set("X-RateLimit-Limit", string(rune(rl.config.Limit)))
			c.Response().Header().Set("X-RateLimit-Remaining", string(rune(remaining)))
			c.Response().Header().Set("X-RateLimit-Reset", string(rune(cl.resetTime.Unix())))

			return next(c)
		}
	}
}

// RateLimit returns a middleware function with default configuration
func RateLimit() echo.MiddlewareFunc {
	return NewRateLimiter(DefaultRateLimiterConfig()).Middleware()
}

// AuthRateLimit returns a middleware function for authentication endpoints
func AuthRateLimit() echo.MiddlewareFunc {
	return NewRateLimiter(AuthRateLimiterConfig()).Middleware()
}
