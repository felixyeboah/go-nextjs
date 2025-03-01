package email

import (
	"fmt"

	"github.com/nanayaw/fullstack/internal/config"
	"github.com/nanayaw/fullstack/internal/service"
)

// NewEmailService creates a new email service based on configuration
func NewEmailService(cfg *config.EmailConfig) (service.EmailService, error) {
	// If Upstash Workflow is configured, use it
	if cfg.UpstashWorkflowURL != "" && cfg.UpstashWorkflowToken != "" {
		return NewUpstashWorkflowService(cfg)
	}

	// Otherwise, use Resend
	if cfg.ResendAPIKey != "" {
		return NewResendService(cfg)
	}

	return nil, fmt.Errorf("no email service configured")
}
