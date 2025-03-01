package security

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/nanayaw/fullstack/internal/config"
	"github.com/nanayaw/fullstack/internal/model"
	"github.com/nanayaw/fullstack/internal/service"
	"github.com/nanayaw/fullstack/pkg/email/templates"
	"github.com/nanayaw/fullstack/pkg/logger"
)

// Repository defines the interface for security-related database operations
type Repository interface {
	// RecordLoginAttempt records a login attempt for a user
	RecordLoginAttempt(ctx context.Context, attempt *model.LoginAttempt) error

	// GetRecentLoginAttempts gets recent login attempts for a user
	GetRecentLoginAttempts(ctx context.Context, userID string, limit int) ([]*model.LoginAttempt, error)

	// LockAccount locks a user account
	LockAccount(ctx context.Context, userID string, until time.Time, reason string) error

	// UnlockAccount unlocks a user account
	UnlockAccount(ctx context.Context, userID string) error

	// IsAccountLocked checks if a user account is locked
	IsAccountLocked(ctx context.Context, userID string) (bool, time.Time, string, error)

	// RecordSecurityEvent records a security event
	RecordSecurityEvent(ctx context.Context, event *model.SecurityEvent) error

	// GetUserSecurityEvents gets security events for a user
	GetUserSecurityEvents(ctx context.Context, userID string, limit int) ([]*model.SecurityEvent, error)
}

// Service provides security-related functionality
type Service struct {
	repo        Repository
	emailSvc    service.EmailService
	config      *config.Config
	logger      logger.Logger
	geoIPLookup GeoIPLookup
}

// GeoIPLookup defines the interface for IP geolocation
type GeoIPLookup interface {
	// GetLocation gets the location information for an IP address
	GetLocation(ip string) (*Location, error)
}

// Location represents geolocation information
type Location struct {
	Country     string
	Region      string
	City        string
	Coordinates struct {
		Latitude  float64
		Longitude float64
	}
}

// NewService creates a new security service
func NewService(repo Repository, emailSvc service.EmailService, config *config.Config, logger logger.Logger) *Service {
	return &Service{
		repo:     repo,
		emailSvc: emailSvc,
		config:   config,
		logger:   logger,
		// Use a simple IP lookup implementation by default
		geoIPLookup: &SimpleGeoIPLookup{},
	}
}

// SetGeoIPLookup sets the GeoIPLookup implementation
func (s *Service) SetGeoIPLookup(lookup GeoIPLookup) {
	s.geoIPLookup = lookup
}

// RecordLoginAttempt records a login attempt and handles security measures
func (s *Service) RecordLoginAttempt(ctx context.Context, userID, email, ipAddress, userAgent string, successful bool) error {
	// Create login attempt record
	attempt := &model.LoginAttempt{
		UserID:      userID,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		Successful:  successful,
		AttemptedAt: time.Now(),
	}

	// Get location information
	location, err := s.getLocationString(ipAddress)
	if err != nil {
		s.logger.Warn("Failed to get location for IP", "ip", ipAddress, "error", err)
		location = "Unknown location"
	}
	attempt.Location = location

	// Record the login attempt
	if err := s.repo.RecordLoginAttempt(ctx, attempt); err != nil {
		return fmt.Errorf("failed to record login attempt: %w", err)
	}

	// If login was successful, check if it's from a new location/device
	if successful {
		return s.handleSuccessfulLogin(ctx, userID, email, attempt)
	}

	// If login failed, check for brute force attempts
	return s.handleFailedLogin(ctx, userID, email, attempt)
}

// handleSuccessfulLogin handles a successful login attempt
func (s *Service) handleSuccessfulLogin(ctx context.Context, userID, email string, attempt *model.LoginAttempt) error {
	// Get recent successful logins
	recentLogins, err := s.repo.GetRecentLoginAttempts(ctx, userID, 5)
	if err != nil {
		return fmt.Errorf("failed to get recent login attempts: %w", err)
	}

	// Check if this is a login from a new location or device
	isNewLocation := true
	isNewDevice := true

	for _, login := range recentLogins {
		// Skip the current login
		if login.ID == attempt.ID {
			continue
		}

		// Skip failed logins
		if !login.Successful {
			continue
		}

		// Check if location matches
		if login.Location == attempt.Location {
			isNewLocation = false
		}

		// Check if user agent matches (simple device fingerprinting)
		if login.UserAgent == attempt.UserAgent {
			isNewDevice = false
		}
	}

	// If this is a login from a new location or device, send a notification
	if isNewLocation || isNewDevice {
		// Send login notification email
		if err := s.sendLoginNotification(ctx, email, attempt); err != nil {
			s.logger.Error("Failed to send login notification", "error", err)
			// Don't return the error, as the login was successful
		}

		// Record security event
		eventType := "new_device_login"
		if isNewLocation {
			eventType = "new_location_login"
		}

		event := &model.SecurityEvent{
			UserID:      userID,
			EventType:   eventType,
			IPAddress:   attempt.IPAddress,
			UserAgent:   attempt.UserAgent,
			Location:    attempt.Location,
			Description: fmt.Sprintf("Login from %s", attempt.Location),
			CreatedAt:   time.Now(),
		}

		if err := s.repo.RecordSecurityEvent(ctx, event); err != nil {
			s.logger.Error("Failed to record security event", "error", err)
		}
	}

	return nil
}

// handleFailedLogin handles a failed login attempt
func (s *Service) handleFailedLogin(ctx context.Context, userID, email string, attempt *model.LoginAttempt) error {
	// If userID is empty, we can't do much (user not found)
	if userID == "" {
		return nil
	}

	// Get recent failed login attempts
	recentAttempts, err := s.repo.GetRecentLoginAttempts(ctx, userID, 10)
	if err != nil {
		return fmt.Errorf("failed to get recent login attempts: %w", err)
	}

	// Count failed attempts in the last hour
	failedCount := 0
	for _, a := range recentAttempts {
		// Skip successful logins
		if a.Successful {
			continue
		}

		// Count failed attempts in the last hour
		if time.Since(a.AttemptedAt) < time.Hour {
			failedCount++
		}
	}

	// If too many failed attempts, lock the account
	if failedCount >= s.config.Security.MaxLoginAttempts {
		lockDuration := s.config.Security.AccountLockDuration
		if lockDuration == 0 {
			lockDuration = 30 * time.Minute
		}

		unlockTime := time.Now().Add(lockDuration)
		reason := fmt.Sprintf("Too many failed login attempts (%d)", failedCount)

		if err := s.repo.LockAccount(ctx, userID, unlockTime, reason); err != nil {
			return fmt.Errorf("failed to lock account: %w", err)
		}

		// Record security event
		event := &model.SecurityEvent{
			UserID:      userID,
			EventType:   "account_locked",
			IPAddress:   attempt.IPAddress,
			UserAgent:   attempt.UserAgent,
			Location:    attempt.Location,
			Description: reason,
			CreatedAt:   time.Now(),
		}

		if err := s.repo.RecordSecurityEvent(ctx, event); err != nil {
			s.logger.Error("Failed to record security event", "error", err)
		}

		// Send account locked email
		if err := s.sendAccountLockedEmail(ctx, email, unlockTime, failedCount); err != nil {
			s.logger.Error("Failed to send account locked email", "error", err)
		}
	}

	return nil
}

// IsAccountLocked checks if a user account is locked
func (s *Service) IsAccountLocked(ctx context.Context, userID string) (bool, time.Time, string, error) {
	return s.repo.IsAccountLocked(ctx, userID)
}

// UnlockAccount unlocks a user account
func (s *Service) UnlockAccount(ctx context.Context, userID string) error {
	return s.repo.UnlockAccount(ctx, userID)
}

// DetectSuspiciousActivity detects suspicious activity for a user
func (s *Service) DetectSuspiciousActivity(ctx context.Context, userID, email, ipAddress, userAgent, activityType string) error {
	// Get location information
	location, err := s.getLocationString(ipAddress)
	if err != nil {
		s.logger.Warn("Failed to get location for IP", "ip", ipAddress, "error", err)
		location = "Unknown location"
	}

	// Get recent security events
	events, err := s.repo.GetUserSecurityEvents(ctx, userID, 10)
	if err != nil {
		return fmt.Errorf("failed to get user security events: %w", err)
	}

	// Check for suspicious patterns
	isSuspicious := false
	suspiciousReason := ""

	// Check for activity from multiple locations in a short time
	locations := make(map[string]bool)
	for _, event := range events {
		if time.Since(event.CreatedAt) < 24*time.Hour {
			locations[event.Location] = true
		}
	}

	// If activity from more than 3 distinct locations in 24 hours, flag as suspicious
	if len(locations) >= 3 && !locations[location] {
		isSuspicious = true
		suspiciousReason = "Activity from multiple locations in a short time period"
	}

	// Record the activity
	event := &model.SecurityEvent{
		UserID:      userID,
		EventType:   activityType,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		Location:    location,
		Description: fmt.Sprintf("%s from %s", activityType, location),
		CreatedAt:   time.Now(),
	}

	if err := s.repo.RecordSecurityEvent(ctx, event); err != nil {
		return fmt.Errorf("failed to record security event: %w", err)
	}

	// If suspicious, send notification
	if isSuspicious {
		// Update event type
		event.EventType = "suspicious_activity"
		event.Description = suspiciousReason

		// Record suspicious activity event
		if err := s.repo.RecordSecurityEvent(ctx, event); err != nil {
			s.logger.Error("Failed to record suspicious activity event", "error", err)
		}

		// Send suspicious activity email
		if err := s.sendSuspiciousActivityEmail(ctx, email, event); err != nil {
			s.logger.Error("Failed to send suspicious activity email", "error", err)
		}
	}

	return nil
}

// NotifyPasswordChanged sends a notification when a password is changed
func (s *Service) NotifyPasswordChanged(ctx context.Context, userID, email, ipAddress, userAgent string) error {
	// Get location information
	location, err := s.getLocationString(ipAddress)
	if err != nil {
		s.logger.Warn("Failed to get location for IP", "ip", ipAddress, "error", err)
		location = "Unknown location"
	}

	// Record security event
	event := &model.SecurityEvent{
		UserID:      userID,
		EventType:   "password_changed",
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		Location:    location,
		Description: "Password changed",
		CreatedAt:   time.Now(),
	}

	if err := s.repo.RecordSecurityEvent(ctx, event); err != nil {
		s.logger.Error("Failed to record security event", "error", err)
	}

	// Send password changed email
	return s.emailSvc.SendPasswordChangedEmail(ctx, email)
}

// sendLoginNotification sends a login notification email
func (s *Service) sendLoginNotification(ctx context.Context, email string, attempt *model.LoginAttempt) error {
	deviceInfo := s.getDeviceInfo(attempt.UserAgent)
	location := attempt.Location

	// Use the EmailService interface method
	return s.emailSvc.SendLoginNotificationEmail(ctx, email, deviceInfo, location)
}

// sendAccountLockedEmail sends an account locked email
func (s *Service) sendAccountLockedEmail(ctx context.Context, email string, unlockTime time.Time, failedAttempts int) error {
	// Currently there's no specific method for account locked emails in the interface
	// We could add one or use a more generic method
	// For now, we'll just log it
	s.logger.Info("Account locked", "email", email, "unlockTime", unlockTime, "failedAttempts", failedAttempts)
	return nil
}

// sendSuspiciousActivityEmail sends a suspicious activity email
func (s *Service) sendSuspiciousActivityEmail(ctx context.Context, email string, event *model.SecurityEvent) error {
	// Currently there's no specific method for suspicious activity emails in the interface
	// We could add one or use a more generic method
	// For now, we'll just log it
	s.logger.Info("Suspicious activity detected",
		"email", email,
		"eventType", event.EventType,
		"location", event.Location,
		"ipAddress", event.IPAddress)
	return nil
}

// getTemplateData returns the common template data
func (s *Service) getTemplateData() templates.TemplateData {
	return templates.TemplateData{
		AppName:      s.config.App.Name,
		SupportEmail: "support@example.com", // Use a default or get from config
		BaseURL:      s.config.App.URL,
		Year:         time.Now().Year(),
	}
}

// getDeviceInfo extracts device information from user agent
func (s *Service) getDeviceInfo(userAgent string) string {
	// This is a very simple implementation
	// In a production environment, you might want to use a more sophisticated
	// user agent parser library

	userAgent = strings.ToLower(userAgent)

	var device, os, browser string

	// Detect device
	switch {
	case strings.Contains(userAgent, "iphone"):
		device = "iPhone"
	case strings.Contains(userAgent, "ipad"):
		device = "iPad"
	case strings.Contains(userAgent, "android") && strings.Contains(userAgent, "mobile"):
		device = "Android Phone"
	case strings.Contains(userAgent, "android"):
		device = "Android Tablet"
	case strings.Contains(userAgent, "macintosh"):
		device = "Mac"
	case strings.Contains(userAgent, "windows"):
		device = "Windows PC"
	case strings.Contains(userAgent, "linux"):
		device = "Linux PC"
	default:
		device = "Unknown Device"
	}

	// Detect OS
	switch {
	case strings.Contains(userAgent, "windows nt 10"):
		os = "Windows 10"
	case strings.Contains(userAgent, "windows nt 6.3"):
		os = "Windows 8.1"
	case strings.Contains(userAgent, "windows nt 6.2"):
		os = "Windows 8"
	case strings.Contains(userAgent, "windows nt 6.1"):
		os = "Windows 7"
	case strings.Contains(userAgent, "mac os x"):
		os = "macOS"
	case strings.Contains(userAgent, "android"):
		os = "Android"
	case strings.Contains(userAgent, "ios"):
		os = "iOS"
	case strings.Contains(userAgent, "linux"):
		os = "Linux"
	default:
		os = "Unknown OS"
	}

	// Detect browser
	switch {
	case strings.Contains(userAgent, "chrome") && !strings.Contains(userAgent, "edg"):
		browser = "Chrome"
	case strings.Contains(userAgent, "firefox"):
		browser = "Firefox"
	case strings.Contains(userAgent, "safari") && !strings.Contains(userAgent, "chrome"):
		browser = "Safari"
	case strings.Contains(userAgent, "edg"):
		browser = "Edge"
	case strings.Contains(userAgent, "opera"):
		browser = "Opera"
	default:
		browser = "Unknown Browser"
	}

	return fmt.Sprintf("%s (%s, %s)", device, os, browser)
}

// getLocationString gets a formatted location string from an IP address
func (s *Service) getLocationString(ipAddress string) (string, error) {
	// Skip for localhost or private IPs
	if ipAddress == "127.0.0.1" || ipAddress == "::1" || isPrivateIP(ipAddress) {
		return "Local Network", nil
	}

	location, err := s.geoIPLookup.GetLocation(ipAddress)
	if err != nil {
		return "", err
	}

	if location.City != "" && location.Country != "" {
		return fmt.Sprintf("%s, %s", location.City, location.Country), nil
	} else if location.Country != "" {
		return location.Country, nil
	}

	return "Unknown Location", nil
}

// isPrivateIP checks if an IP address is private
func isPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	// Check if IPv4 private address
	if ip4 := ip.To4(); ip4 != nil {
		// Check private IPv4 ranges
		// 10.0.0.0/8
		if ip4[0] == 10 {
			return true
		}
		// 172.16.0.0/12
		if ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31 {
			return true
		}
		// 192.168.0.0/16
		if ip4[0] == 192 && ip4[1] == 168 {
			return true
		}
	}

	// Check if IPv6 local address
	return ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast()
}

// SimpleGeoIPLookup is a simple implementation of GeoIPLookup
// In a production environment, you would use a proper GeoIP database
// like MaxMind GeoIP or IP-API
type SimpleGeoIPLookup struct{}

// GetLocation returns a dummy location for demonstration purposes
func (s *SimpleGeoIPLookup) GetLocation(ip string) (*Location, error) {
	// In a real implementation, you would look up the IP in a GeoIP database
	// For now, return a dummy location
	return &Location{
		Country: "Unknown Country",
		Region:  "Unknown Region",
		City:    "Unknown City",
		Coordinates: struct {
			Latitude  float64
			Longitude float64
		}{
			Latitude:  0,
			Longitude: 0,
		},
	}, nil
}
