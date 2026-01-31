package veilmail

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	t.Run("valid live key", func(t *testing.T) {
		client, err := New("veil_live_test123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if client == nil {
			t.Fatal("expected client, got nil")
		}
		if client.Emails == nil {
			t.Error("expected Emails service")
		}
		if client.Domains == nil {
			t.Error("expected Domains service")
		}
		if client.Templates == nil {
			t.Error("expected Templates service")
		}
		if client.Audiences == nil {
			t.Error("expected Audiences service")
		}
		if client.Campaigns == nil {
			t.Error("expected Campaigns service")
		}
		if client.Webhooks == nil {
			t.Error("expected Webhooks service")
		}
		if client.Topics == nil {
			t.Error("expected Topics service")
		}
		if client.Properties == nil {
			t.Error("expected Properties service")
		}
	})

	t.Run("valid test key", func(t *testing.T) {
		client, err := New("veil_test_test123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if client == nil {
			t.Fatal("expected client, got nil")
		}
	})

	t.Run("invalid key format", func(t *testing.T) {
		_, err := New("invalid_key")
		if err == nil {
			t.Fatal("expected error for invalid key")
		}
	})

	t.Run("empty key", func(t *testing.T) {
		_, err := New("")
		if err == nil {
			t.Fatal("expected error for empty key")
		}
	})
}

func TestClientOptions(t *testing.T) {
	t.Run("custom base URL", func(t *testing.T) {
		client, err := New("veil_test_xxx", WithBaseURL("https://custom.api.com/"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if client.http.baseURL != "https://custom.api.com" {
			t.Errorf("expected base URL without trailing slash, got %s", client.http.baseURL)
		}
	})

	t.Run("custom timeout", func(t *testing.T) {
		client, err := New("veil_test_xxx", WithTimeout(5*time.Second))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if client.http.client.Timeout != 5*time.Second {
			t.Errorf("expected 5s timeout, got %v", client.http.client.Timeout)
		}
	})

	t.Run("custom HTTP client", func(t *testing.T) {
		custom := &http.Client{Timeout: 10 * time.Second}
		client, err := New("veil_test_xxx", WithHTTPClient(custom))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if client.http.client != custom {
			t.Error("expected custom HTTP client")
		}
	})
}

func TestFormatEmailAddress(t *testing.T) {
	tests := []struct {
		testName   string
		inputName  string
		inputEmail string
		want       string
	}{
		{"email only", "", "hello@example.com", "hello@example.com"},
		{"with name", "Alice", "hello@example.com", "Alice <hello@example.com>"},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			got := FormatEmailAddress(tt.inputName, tt.inputEmail)
			if got != tt.want {
				t.Errorf("FormatEmailAddress(%q, %q) = %q, want %q", tt.inputName, tt.inputEmail, got, tt.want)
			}
		})
	}
}

func TestErrorParsing(t *testing.T) {
	t.Run("401 authentication error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(401)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "invalid_api_key",
					"message": "Invalid API key",
				},
			})
		}))
		defer server.Close()

		client, _ := New("veil_test_xxx", WithBaseURL(server.URL))
		_, err := client.Emails.Get(context.Background(), "test")
		if err == nil {
			t.Fatal("expected error")
		}
		if !IsAuthenticationError(err) {
			t.Errorf("expected authentication error, got %T", err)
		}
	})

	t.Run("404 not found error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "not_found",
					"message": "Email not found",
				},
			})
		}))
		defer server.Close()

		client, _ := New("veil_test_xxx", WithBaseURL(server.URL))
		_, err := client.Emails.Get(context.Background(), "nonexistent")
		if err == nil {
			t.Fatal("expected error")
		}
		if !IsNotFoundError(err) {
			t.Errorf("expected not found error, got %T", err)
		}
	})

	t.Run("422 PII detected error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(422)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "pii_detected",
					"message": "PII detected in email content",
					"details": map[string]interface{}{
						"piiTypes": []string{"PHONE_NUMBER", "EMAIL_ADDRESS"},
					},
				},
			})
		}))
		defer server.Close()

		client, _ := New("veil_test_xxx", WithBaseURL(server.URL))
		_, err := client.Emails.Send(context.Background(), &SendEmailParams{
			From:    "test@example.com",
			To:      []string{"user@example.com"},
			Subject: "Test",
			HTML:    "<p>Test</p>",
		})
		if err == nil {
			t.Fatal("expected error")
		}
		if !IsPiiDetectedError(err) {
			t.Errorf("expected PII detected error, got %T", err)
		}
		var apiErr *Error
		if ok := errorAs(err, &apiErr); ok {
			if len(apiErr.PiiTypes) != 2 {
				t.Errorf("expected 2 PII types, got %d", len(apiErr.PiiTypes))
			}
		}
	})

	t.Run("429 rate limit error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Retry-After", "30")
			w.WriteHeader(429)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "rate_limited",
					"message": "Rate limit exceeded",
				},
			})
		}))
		defer server.Close()

		client, _ := New("veil_test_xxx", WithBaseURL(server.URL))
		_, err := client.Emails.List(context.Background(), nil)
		if err == nil {
			t.Fatal("expected error")
		}
		if !IsRateLimitError(err) {
			t.Errorf("expected rate limit error, got %T", err)
		}
		var apiErr *Error
		if ok := errorAs(err, &apiErr); ok {
			if apiErr.RetryAfter != 30 {
				t.Errorf("expected RetryAfter=30, got %d", apiErr.RetryAfter)
			}
		}
	})

	t.Run("500 server error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "internal_error",
					"message": "Internal server error",
				},
			})
		}))
		defer server.Close()

		client, _ := New("veil_test_xxx", WithBaseURL(server.URL))
		_, err := client.Emails.Get(context.Background(), "test")
		if err == nil {
			t.Fatal("expected error")
		}
		if !IsServerError(err) {
			t.Errorf("expected server error, got %T", err)
		}
	})
}

func TestEmailSend(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/emails" {
			t.Errorf("expected /v1/emails, got %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer veil_test_xxx" {
			t.Errorf("expected bearer auth, got %s", r.Header.Get("Authorization"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected JSON content type, got %s", r.Header.Get("Content-Type"))
		}
		if r.Header.Get("User-Agent") != userAgent {
			t.Errorf("expected user agent %s, got %s", userAgent, r.Header.Get("User-Agent"))
		}

		var body SendEmailParams
		json.NewDecoder(r.Body).Decode(&body)
		if body.From != "hello@example.com" {
			t.Errorf("expected from=hello@example.com, got %s", body.From)
		}

		w.WriteHeader(201)
		json.NewEncoder(w).Encode(Email{
			ID:      "email_123",
			From:    body.From,
			To:      body.To,
			Subject: body.Subject,
			Status:  EmailStatusQueued,
		})
	}))
	defer server.Close()

	client, _ := New("veil_test_xxx", WithBaseURL(server.URL))
	email, err := client.Emails.Send(context.Background(), &SendEmailParams{
		From:    "hello@example.com",
		To:      []string{"user@example.com"},
		Subject: "Test",
		HTML:    "<p>Hello</p>",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if email.ID != "email_123" {
		t.Errorf("expected id=email_123, got %s", email.ID)
	}
	if email.Status != EmailStatusQueued {
		t.Errorf("expected status=queued, got %s", email.Status)
	}
}

func TestDomainCreate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": Domain{
				ID:     "dom_123",
				Domain: "mail.example.com",
				Status: DomainStatusPending,
			},
		})
	}))
	defer server.Close()

	client, _ := New("veil_test_xxx", WithBaseURL(server.URL))
	domain, err := client.Domains.Create(context.Background(), &CreateDomainParams{
		Domain: "mail.example.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if domain.ID != "dom_123" {
		t.Errorf("expected id=dom_123, got %s", domain.ID)
	}
}

func TestContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		w.WriteHeader(200)
	}))
	defer server.Close()

	client, _ := New("veil_test_xxx", WithBaseURL(server.URL))
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	_, err := client.Emails.Get(ctx, "test")
	if err == nil {
		t.Fatal("expected error from cancelled context")
	}
}

// errorAs is a helper wrapping errors.As for test readability.
func errorAs(err error, target interface{}) bool {
	switch t := target.(type) {
	case **Error:
		e, ok := err.(*Error)
		if ok {
			*t = e
			return true
		}
	}
	return false
}
