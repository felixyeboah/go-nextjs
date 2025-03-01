package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"

	// Import the docs package to register Swagger docs
	// _ "github.com/nanayaw/fullstack/docs"

	"github.com/nanayaw/fullstack/internal/config"
	authHandler "github.com/nanayaw/fullstack/internal/handler/auth"
	userHandler "github.com/nanayaw/fullstack/internal/handler/user"
	"github.com/nanayaw/fullstack/internal/router"
	"github.com/nanayaw/fullstack/internal/service/auth"
	"github.com/nanayaw/fullstack/internal/service/cache"
	"github.com/nanayaw/fullstack/internal/service/email"
	"github.com/nanayaw/fullstack/internal/service/user"
)

// @title           Fullstack API
// @version         1.0
// @description     A modern fullstack application with Go and Next.js

// @host      localhost:8080
// @BasePath  /api/v1

func main() {
	// Parse command line flags
	port := flag.Int("port", 0, "Port to run the server on (overrides config)")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Override port if specified via flag
	if *port != 0 {
		cfg.Server.Port = *port
	}

	// Check for SERVER_PORT environment variable
	if envPort := os.Getenv("SERVER_PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			cfg.Server.Port = p
			log.Printf("Using SERVER_PORT from environment: %d", p)
		}
	}

	// Initialize Echo
	e := echo.New()

	// Initialize services
	cacheService, err := cache.NewRedisService(&cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to initialize cache service: %v", err)
	}

	emailService, err := email.NewEmailService(&cfg.Email)
	if err != nil {
		log.Fatalf("Failed to initialize email service: %v", err)
	}

	// Initialize user service
	userService := user.NewService(nil) // Replace with actual repository

	authService, err := auth.NewPasetoService(&cfg.Auth, nil, emailService, cacheService)
	if err != nil {
		log.Fatalf("Failed to initialize auth service: %v", err)
	}

	// Initialize handlers
	authHandler := authHandler.NewHandler(authService)
	userHandler := userHandler.NewHandler(userService, authService)

	// Initialize router
	r := router.NewRouter(e, authHandler, userHandler, authService)
	r.SetupRoutes()
	r.SetupTimeoutMiddleware(int(cfg.Server.ReadTimeout.Seconds()))

	// Start server
	go func() {
		addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
		log.Printf("Starting server on %s", addr)
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
