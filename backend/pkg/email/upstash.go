package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// UpstashWorkflowClient is a client for the Upstash Workflow API
type UpstashWorkflowClient struct {
	baseURL   string
	token     string
	client    *http.Client
	fromEmail string
	fromName  string
}

// WorkflowRequest represents a request to the Upstash Workflow API
type WorkflowRequest struct {
	EmailType   string                 `json:"emailType"`
	To          string                 `json:"to"`
	Subject     string                 `json:"subject"`
	Data        map[string]interface{} `json:"data"`
	FromEmail   string                 `json:"fromEmail"`
	FromName    string                 `json:"fromName"`
	ScheduledAt *time.Time             `json:"scheduledAt,omitempty"`
}

// WorkflowResponse represents a response from the Upstash Workflow API
type WorkflowResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// NewUpstashWorkflowClient creates a new Upstash Workflow client
func NewUpstashWorkflowClient(baseURL, token, fromEmail, fromName string) *UpstashWorkflowClient {
	return &UpstashWorkflowClient{
		baseURL:   baseURL,
		token:     token,
		client:    &http.Client{Timeout: 10 * time.Second},
		fromEmail: fromEmail,
		fromName:  fromName,
	}
}

// SendEmail sends an email using the Upstash Workflow API
func (c *UpstashWorkflowClient) SendEmail(ctx context.Context, emailType, to, subject string, data map[string]interface{}) error {
	req := WorkflowRequest{
		EmailType: emailType,
		To:        to,
		Subject:   subject,
		Data:      data,
		FromEmail: c.fromEmail,
		FromName:  c.fromName,
	}

	return c.sendWorkflowRequest(ctx, req)
}

// ScheduleEmail schedules an email to be sent at a specific time
func (c *UpstashWorkflowClient) ScheduleEmail(ctx context.Context, emailType, to, subject string, data map[string]interface{}, scheduledAt time.Time) error {
	req := WorkflowRequest{
		EmailType:   emailType,
		To:          to,
		Subject:     subject,
		Data:        data,
		FromEmail:   c.fromEmail,
		FromName:    c.fromName,
		ScheduledAt: &scheduledAt,
	}

	return c.sendWorkflowRequest(ctx, req)
}

// sendWorkflowRequest sends a request to the Upstash Workflow API
func (c *UpstashWorkflowClient) sendWorkflowRequest(ctx context.Context, req WorkflowRequest) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal workflow request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create workflow request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send workflow request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("workflow request failed with status code: %d", resp.StatusCode)
	}

	var workflowResp WorkflowResponse
	if err := json.NewDecoder(resp.Body).Decode(&workflowResp); err != nil {
		return fmt.Errorf("failed to decode workflow response: %w", err)
	}

	if workflowResp.Status != "queued" && workflowResp.Status != "scheduled" {
		return fmt.Errorf("workflow request failed with status: %s", workflowResp.Status)
	}

	return nil
}
