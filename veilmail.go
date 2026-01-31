// Package veilmail provides a Go client for the Veil Mail API.
//
// Create a client with your API key and use the service properties to access
// the API:
//
//	client, err := veilmail.New("veil_live_xxx")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	email, err := client.Emails.Send(context.Background(), &veilmail.SendEmailParams{
//	    From:    "hello@yourdomain.com",
//	    To:      []string{"user@example.com"},
//	    Subject: "Hello",
//	    HTML:    "<p>Hello World</p>",
//	})
package veilmail

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://api.veilmail.xyz"
	defaultTimeout = 30 * time.Second
	userAgent      = "veilmail-go/0.1.0"
)

// Client is the Veil Mail API client.
type Client struct {
	// Emails provides methods for sending and managing emails.
	Emails *EmailsService

	// Domains provides methods for domain management and verification.
	Domains *DomainsService

	// Templates provides methods for template management.
	Templates *TemplatesService

	// Audiences provides methods for audience and subscriber management.
	Audiences *AudiencesService

	// Campaigns provides methods for campaign management.
	Campaigns *CampaignsService

	// Webhooks provides methods for webhook management and signature verification.
	Webhooks *WebhooksService

	// Topics provides methods for subscription topic management.
	Topics *TopicsService

	// Properties provides methods for contact property management.
	Properties *PropertiesService

	// Sequences provides methods for automation sequence management.
	Sequences *SequencesService

	// Feeds provides methods for RSS feed management.
	Feeds *FeedsService

	// Forms provides methods for signup form management.
	Forms *FormsService

	// Analytics provides methods for geo and device analytics.
	Analytics *AnalyticsService

	http *httpClient
}

// Option configures a Client.
type Option func(*Client)

// WithBaseURL sets the API base URL.
// The default is https://api.veilmail.xyz.
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.http.baseURL = strings.TrimRight(url, "/")
	}
}

// WithHTTPClient sets the underlying HTTP client used for API requests.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.http.client = client
	}
}

// WithTimeout sets the request timeout.
// The default is 30 seconds.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.http.client.Timeout = d
	}
}

// New creates a new Veil Mail API client.
// The API key must start with "veil_live_" or "veil_test_".
func New(apiKey string, opts ...Option) (*Client, error) {
	if !strings.HasPrefix(apiKey, "veil_live_") && !strings.HasPrefix(apiKey, "veil_test_") {
		return nil, fmt.Errorf("veilmail: invalid API key format, must start with 'veil_live_' or 'veil_test_'")
	}

	h := &httpClient{
		apiKey:  apiKey,
		baseURL: defaultBaseURL,
		client: &http.Client{
			Timeout: defaultTimeout,
		},
		userAgent: userAgent,
	}

	c := &Client{http: h}

	for _, opt := range opts {
		opt(c)
	}

	c.Emails = &EmailsService{client: h}
	c.Domains = &DomainsService{client: h}
	c.Templates = &TemplatesService{client: h}
	c.Audiences = &AudiencesService{client: h}
	c.Campaigns = &CampaignsService{client: h}
	c.Webhooks = &WebhooksService{client: h}
	c.Topics = &TopicsService{client: h}
	c.Properties = &PropertiesService{client: h}
	c.Sequences = &SequencesService{client: h}
	c.Feeds = &FeedsService{client: h}
	c.Forms = &FormsService{client: h}
	c.Analytics = &AnalyticsService{client: h}

	return c, nil
}

// FormatEmailAddress formats a name and email into "Name <email>" format.
// If name is empty, the email is returned as-is.
func FormatEmailAddress(name, email string) string {
	if name == "" {
		return email
	}
	return fmt.Sprintf("%s <%s>", name, email)
}
