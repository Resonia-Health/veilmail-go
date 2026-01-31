package veilmail

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
)

// WebhooksService handles webhook operations.
type WebhooksService struct {
	client *httpClient
}

// Create creates a new webhook endpoint.
func (s *WebhooksService) Create(ctx context.Context, params *CreateWebhookParams) (*Webhook, error) {
	var resp dataWrapper[Webhook]
	err := s.client.post(ctx, "/v1/webhooks", params, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// List returns a paginated list of webhooks.
func (s *WebhooksService) List(ctx context.Context, params *ListParams) (*ListResponse[Webhook], error) {
	q := url.Values{}
	addPaginationParams(q, params)
	var resp ListResponse[Webhook]
	err := s.client.get(ctx, "/v1/webhooks", q, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves a single webhook by ID.
func (s *WebhooksService) Get(ctx context.Context, id string) (*Webhook, error) {
	var resp dataWrapper[Webhook]
	err := s.client.get(ctx, fmt.Sprintf("/v1/webhooks/%s", id), nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Update updates a webhook.
func (s *WebhooksService) Update(ctx context.Context, id string, params *UpdateWebhookParams) (*Webhook, error) {
	var resp dataWrapper[Webhook]
	err := s.client.patch(ctx, fmt.Sprintf("/v1/webhooks/%s", id), params, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Delete removes a webhook.
func (s *WebhooksService) Delete(ctx context.Context, id string) error {
	return s.client.del(ctx, fmt.Sprintf("/v1/webhooks/%s", id))
}

// Test sends a test event to the webhook endpoint and returns the result.
func (s *WebhooksService) Test(ctx context.Context, id string) (*WebhookTestResult, error) {
	var resp dataWrapper[WebhookTestResult]
	err := s.client.post(ctx, fmt.Sprintf("/v1/webhooks/%s/test", id), nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// RotateSecret generates a new signing secret for the webhook.
func (s *WebhooksService) RotateSecret(ctx context.Context, id string) (*Webhook, error) {
	var resp dataWrapper[Webhook]
	err := s.client.post(ctx, fmt.Sprintf("/v1/webhooks/%s/rotate-secret", id), nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// VerifySignature verifies a webhook request signature using constant-time
// HMAC-SHA256 comparison.
//
// Use this in your webhook handler to verify that the request came from Veil Mail:
//
//	body, _ := io.ReadAll(r.Body)
//	signature := r.Header.Get("X-Signature-Hash")
//	if !veilmail.VerifySignature(body, signature, webhook.Secret) {
//	    http.Error(w, "Invalid signature", http.StatusUnauthorized)
//	    return
//	}
func VerifySignature(body []byte, signature string, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
