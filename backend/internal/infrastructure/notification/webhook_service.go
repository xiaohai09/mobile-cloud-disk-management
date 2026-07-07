package infrastructure

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"caiyun/internal/domain/entity"
)

// WebhookService handles webhook notifications
type WebhookService struct {
	httpClient *http.Client
}

// NewWebhookService creates a new webhook service
func NewWebhookService() *WebhookService {
	return &WebhookService{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// WebhookPayload represents the payload sent to webhook endpoints
type WebhookPayload struct {
	Event     domain.WebhookEventType `json:"event"`
	Timestamp time.Time               `json:"timestamp"`
	Data      map[string]interface{}  `json:"data"`
	Signature string                  `json:"signature"`
}

// SendWebhook sends a webhook notification to the specified endpoint
func (s *WebhookService) SendWebhook(endpoint *domain.WebhookEndpoint, event domain.WebhookEventType, data map[string]interface{}) error {
	if !endpoint.IsActive {
		return nil
	}

	payload := WebhookPayload{
		Event:     event,
		Timestamp: time.Now(),
		Data:      data,
	}

	// Generate signature if secret is provided
	if endpoint.Secret != "" {
		payload.Signature = s.generateSignature(endpoint.Secret, payload)
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal webhook payload failed: %w", err)
	}

	req, err := http.NewRequest("POST", endpoint.URL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("create webhook request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "caiyun-webhook/1.0")

	// Add custom headers
	for key, value := range endpoint.Headers {
		req.Header.Set(key, value)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send webhook failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned non-success status: %d", resp.StatusCode)
	}

	return nil
}

// generateSignature generates HMAC-SHA256 signature for webhook payload
func (s *WebhookService) generateSignature(secret string, payload WebhookPayload) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(fmt.Sprintf("%s:%d", payload.Event, payload.Timestamp.Unix())))
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

// VerifySignature verifies the webhook signature
func (s *WebhookService) VerifySignature(secret, signature string, payload WebhookPayload) bool {
	expected := s.generateSignature(secret, payload)
	return hmac.Equal([]byte(signature), []byte(expected))
}
