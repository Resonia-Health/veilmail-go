package veilmail

// ListParams contains pagination parameters for list operations.
type ListParams struct {
	// Limit is the maximum number of results to return (default 20, max 100).
	Limit int
	// Cursor is the pagination cursor for fetching the next page.
	Cursor string
}

// ListResponse is a paginated list response.
type ListResponse[T any] struct {
	Data       []T    `json:"data"`
	HasMore    bool   `json:"hasMore"`
	NextCursor string `json:"nextCursor,omitempty"`
}

// dataWrapper is used internally for endpoints that return { data: T }.
type dataWrapper[T any] struct {
	Data T `json:"data"`
}

// --- Email Types ---

// EmailStatus represents the status of an email.
type EmailStatus string

const (
	EmailStatusPending    EmailStatus = "pending"
	EmailStatusQueued     EmailStatus = "queued"
	EmailStatusSending    EmailStatus = "sending"
	EmailStatusSent       EmailStatus = "sent"
	EmailStatusDelivered  EmailStatus = "delivered"
	EmailStatusBounced    EmailStatus = "bounced"
	EmailStatusFailed     EmailStatus = "failed"
	EmailStatusCancelled  EmailStatus = "cancelled"
	EmailStatusScheduled  EmailStatus = "scheduled"
	EmailStatusComplained EmailStatus = "complained"
)

// SendEmailParams contains the parameters for sending an email.
type SendEmailParams struct {
	// From is the sender address. Use "Name <email>" format or just "email".
	From string `json:"from"`
	// To is the list of recipient addresses.
	To []string `json:"to"`
	// Subject is the email subject line.
	Subject string `json:"subject"`
	// HTML is the HTML body content.
	HTML string `json:"html,omitempty"`
	// Text is the plain text body content.
	Text string `json:"text,omitempty"`
	// CC is the list of carbon copy recipients.
	CC []string `json:"cc,omitempty"`
	// BCC is the list of blind carbon copy recipients.
	BCC []string `json:"bcc,omitempty"`
	// ReplyTo is the reply-to address.
	ReplyTo string `json:"replyTo,omitempty"`
	// Headers are custom email headers (X-* headers only).
	Headers map[string]string `json:"headers,omitempty"`
	// TemplateID is the ID of a template to use for content.
	TemplateID string `json:"templateId,omitempty"`
	// TemplateData is the variable data for template rendering.
	TemplateData map[string]interface{} `json:"templateData,omitempty"`
	// ScheduledFor is the ISO 8601 time to schedule the email for future delivery.
	ScheduledFor string `json:"scheduledFor,omitempty"`
	// Tags are categorization tags for the email.
	Tags []string `json:"tags,omitempty"`
	// Metadata is custom metadata attached to the email.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	// IdempotencyKey prevents duplicate sends when retrying.
	IdempotencyKey string `json:"idempotencyKey,omitempty"`
	// UnsubscribeURL is the unsubscribe link for marketing emails.
	UnsubscribeURL string `json:"unsubscribeUrl,omitempty"`
	// Type is the email classification: "transactional" or "marketing".
	Type string `json:"type,omitempty"`
	// Attachments are file attachments (max 10, 10MB total).
	Attachments []EmailAttachment `json:"attachments,omitempty"`
	// TopicID is the subscription topic ID for topic-based preferences.
	TopicID string `json:"topicId,omitempty"`
}

// EmailAttachment represents a file attachment on an email.
type EmailAttachment struct {
	// Filename is the attachment file name.
	Filename string `json:"filename"`
	// Content is the base64-encoded file content.
	Content string `json:"content,omitempty"`
	// URL is an HTTPS URL to fetch the attachment from.
	URL string `json:"url,omitempty"`
	// ContentType is the MIME type of the attachment.
	ContentType string `json:"contentType"`
	// ContentID is the Content-ID for inline attachments.
	ContentID string `json:"contentId,omitempty"`
}

// EmailSecurity contains PII scanning results.
type EmailSecurity struct {
	Scanned    bool          `json:"scanned"`
	Clean      bool          `json:"clean"`
	Detections []interface{} `json:"detections"`
}

// EmailEvent represents an event in an email's lifecycle.
type EmailEvent struct {
	Type      string      `json:"type"`
	Timestamp string      `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// Email represents an email message.
type Email struct {
	ID          string         `json:"id"`
	From        string         `json:"from"`
	To          []string       `json:"to"`
	Subject     string         `json:"subject"`
	Status      EmailStatus    `json:"status"`
	CreatedAt   string         `json:"createdAt"`
	Security    *EmailSecurity `json:"security,omitempty"`
	SentAt      *string        `json:"sentAt,omitempty"`
	DeliveredAt *string        `json:"deliveredAt,omitempty"`
	OpenedAt    *string        `json:"openedAt,omitempty"`
	OpenCount   int            `json:"openCount,omitempty"`
	ClickCount  int            `json:"clickCount,omitempty"`
	Events      []EmailEvent   `json:"events,omitempty"`
}

// ListEmailsParams contains filters for listing emails.
type ListEmailsParams struct {
	ListParams
	// Status filters by email status.
	Status EmailStatus
	// Tag filters by tag.
	Tag string
	// After filters for emails created after this ISO 8601 date.
	After string
	// Before filters for emails created before this ISO 8601 date.
	Before string
}

// BatchSendResult contains the results of a batch send operation.
type BatchSendResult struct {
	Total      int                `json:"total"`
	Successful int                `json:"successful"`
	Failed     int                `json:"failed"`
	Results    []BatchEmailResult `json:"results"`
}

// BatchEmailResult represents the result for a single email in a batch.
type BatchEmailResult struct {
	Index  int              `json:"index"`
	ID     string           `json:"id,omitempty"`
	Status string           `json:"status"`
	Error  *BatchEmailError `json:"error,omitempty"`
}

// BatchEmailError contains error details for a failed email in a batch.
type BatchEmailError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// CancelEmailResult is returned when cancelling a scheduled email.
type CancelEmailResult struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	CancelledAt string `json:"cancelledAt"`
}

// UpdateEmailParams contains the parameters for rescheduling a scheduled email.
type UpdateEmailParams struct {
	// ScheduledFor is the new ISO 8601 scheduled time.
	ScheduledFor string `json:"scheduledFor"`
}

// --- Domain Types ---

// DomainStatus represents the verification status of a domain.
type DomainStatus string

const (
	DomainStatusPending  DomainStatus = "pending"
	DomainStatusVerified DomainStatus = "verified"
	DomainStatusFailed   DomainStatus = "failed"
)

// CreateDomainParams contains the parameters for adding a domain.
type CreateDomainParams struct {
	// Domain is the domain name (e.g., "mail.example.com").
	Domain string `json:"domain"`
}

// UpdateDomainParams contains the parameters for updating domain settings.
type UpdateDomainParams struct {
	// TrackOpens enables or disables open tracking for this domain.
	TrackOpens *bool `json:"trackOpens,omitempty"`
	// TrackClicks enables or disables click tracking for this domain.
	TrackClicks *bool `json:"trackClicks,omitempty"`
}

// DnsRecord represents a DNS record required for domain verification.
type DnsRecord struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Value    string `json:"value"`
	Priority int    `json:"priority,omitempty"`
	Status   string `json:"status"`
}

// DomainTracking contains tracking configuration for a domain.
type DomainTracking struct {
	Opens  bool `json:"opens"`
	Clicks bool `json:"clicks"`
}

// DomainDnsHealth contains DNS health information for a domain.
type DomainDnsHealth struct {
	Healthy       bool     `json:"healthy"`
	Issues        []string `json:"issues"`
	LastCheckedAt *string  `json:"lastCheckedAt"`
}

// Domain represents a sending domain.
type Domain struct {
	ID         string          `json:"id"`
	Domain     string          `json:"domain"`
	Status     DomainStatus    `json:"status"`
	DnsRecords []DnsRecord     `json:"dnsRecords"`
	Tracking   DomainTracking  `json:"tracking"`
	DnsHealth  DomainDnsHealth `json:"dnsHealth"`
	VerifiedAt *string         `json:"verifiedAt"`
	CreatedAt  string          `json:"createdAt"`
	UpdatedAt  string          `json:"updatedAt"`
}

// --- Template Types ---

// TemplateType classifies a template.
type TemplateType string

const (
	TemplateTypeTransactional TemplateType = "transactional"
	TemplateTypeMarketing     TemplateType = "marketing"
)

// TemplateVariable defines a variable used in a template.
type TemplateVariable struct {
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	Required     bool        `json:"required,omitempty"`
	DefaultValue interface{} `json:"defaultValue,omitempty"`
	Description  string      `json:"description,omitempty"`
}

// CreateTemplateParams contains the parameters for creating a template.
type CreateTemplateParams struct {
	Name        string             `json:"name"`
	Subject     string             `json:"subject"`
	HTML        string             `json:"html"`
	Text        string             `json:"text,omitempty"`
	Type        TemplateType       `json:"type,omitempty"`
	Variables   []TemplateVariable `json:"variables,omitempty"`
	Description string             `json:"description,omitempty"`
}

// UpdateTemplateParams contains the parameters for updating a template.
type UpdateTemplateParams struct {
	Name        *string            `json:"name,omitempty"`
	Subject     *string            `json:"subject,omitempty"`
	HTML        *string            `json:"html,omitempty"`
	Text        *string            `json:"text,omitempty"`
	Type        TemplateType       `json:"type,omitempty"`
	Variables   []TemplateVariable `json:"variables,omitempty"`
	Description *string            `json:"description,omitempty"`
}

// Template represents an email template.
type Template struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Subject     string             `json:"subject"`
	HTML        string             `json:"html"`
	Text        *string            `json:"text"`
	Type        TemplateType       `json:"type"`
	Variables   []TemplateVariable `json:"variables"`
	Description *string            `json:"description"`
	CreatedAt   string             `json:"createdAt"`
	UpdatedAt   string             `json:"updatedAt"`
}

// PreviewTemplateParams contains the parameters for previewing a template.
type PreviewTemplateParams struct {
	HTML      string                 `json:"html"`
	Subject   string                 `json:"subject,omitempty"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// PreviewTemplateResponse is the result of a template preview.
type PreviewTemplateResponse struct {
	HTML    string `json:"html"`
	Subject string `json:"subject,omitempty"`
}

// --- Audience Types ---

// CreateAudienceParams contains the parameters for creating an audience.
type CreateAudienceParams struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// UpdateAudienceParams contains the parameters for updating an audience.
type UpdateAudienceParams struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// Audience represents a subscriber audience.
type Audience struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Description     *string `json:"description"`
	SubscriberCount int     `json:"subscriberCount"`
	CreatedAt       string  `json:"createdAt"`
	UpdatedAt       string  `json:"updatedAt"`
}

// --- Subscriber Types ---

// SubscriberStatus represents the status of a subscriber.
type SubscriberStatus string

const (
	SubscriberStatusPending      SubscriberStatus = "pending"
	SubscriberStatusActive       SubscriberStatus = "active"
	SubscriberStatusUnsubscribed SubscriberStatus = "unsubscribed"
	SubscriberStatusBounced      SubscriberStatus = "bounced"
	SubscriberStatusComplained   SubscriberStatus = "complained"
)

// ConsentType represents the type of consent given by a subscriber.
type ConsentType string

const (
	ConsentTypeExpress ConsentType = "express"
	ConsentTypeImplied ConsentType = "implied"
	ConsentTypeNotSet  ConsentType = "not_set"
)

// AddSubscriberParams contains the parameters for adding a subscriber.
type AddSubscriberParams struct {
	Email            string                 `json:"email"`
	FirstName        string                 `json:"firstName,omitempty"`
	LastName         string                 `json:"lastName,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
	DoubleOptIn      bool                   `json:"doubleOptIn,omitempty"`
	Status           string                 `json:"status,omitempty"`
	ConsentType      ConsentType            `json:"consentType,omitempty"`
	ConsentSource    string                 `json:"consentSource,omitempty"`
	ConsentDate      string                 `json:"consentDate,omitempty"`
	ConsentExpiresAt string                 `json:"consentExpiresAt,omitempty"`
	ConsentProof     string                 `json:"consentProof,omitempty"`
}

// UpdateSubscriberParams contains the parameters for updating a subscriber.
type UpdateSubscriberParams struct {
	FirstName *string                `json:"firstName,omitempty"`
	LastName  *string                `json:"lastName,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Status    string                 `json:"status,omitempty"`
}

// Subscriber represents an audience subscriber.
type Subscriber struct {
	ID               string           `json:"id"`
	Email            string           `json:"email"`
	FirstName        *string          `json:"firstName"`
	LastName         *string          `json:"lastName"`
	Status           SubscriberStatus `json:"status"`
	Metadata         interface{}      `json:"metadata"`
	Source           *string          `json:"source"`
	SubscribedAt     *string          `json:"subscribedAt"`
	ConfirmedAt      *string          `json:"confirmedAt"`
	UnsubscribedAt   *string          `json:"unsubscribedAt"`
	ConsentType      *string          `json:"consentType"`
	ConsentSource    *string          `json:"consentSource"`
	ConsentDate      *string          `json:"consentDate"`
	ConsentExpiresAt *string          `json:"consentExpiresAt"`
	CreatedAt        string           `json:"createdAt"`
	UpdatedAt        string           `json:"updatedAt,omitempty"`
}

// ListSubscribersParams contains filters for listing subscribers.
type ListSubscribersParams struct {
	ListParams
	// Status filters by subscriber status.
	Status SubscriberStatus
	// Email filters by email address (search).
	Email string
}

// ImportSubscriberEntry represents a single subscriber in a bulk import.
type ImportSubscriberEntry struct {
	Email     string                 `json:"email"`
	FirstName string                 `json:"firstName,omitempty"`
	LastName  string                 `json:"lastName,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// ImportSubscribersParams contains the parameters for bulk importing subscribers.
type ImportSubscribersParams struct {
	// Subscribers is a list of subscriber entries (use this or CSVData).
	Subscribers []ImportSubscriberEntry `json:"subscribers,omitempty"`
	// CSVData is raw CSV content for import (use this or Subscribers).
	CSVData string `json:"csvData,omitempty"`
}

// ImportError represents an error for a single row during import.
type ImportError struct {
	Row   int    `json:"row"`
	Error string `json:"error"`
}

// ImportSubscribersResult contains the results of a bulk import.
type ImportSubscribersResult struct {
	Total   int           `json:"total"`
	Created int           `json:"created"`
	Updated int           `json:"updated"`
	Skipped int           `json:"skipped"`
	Errors  []ImportError `json:"errors"`
}

// ExportSubscribersParams contains filters for exporting subscribers.
type ExportSubscribersParams struct {
	// Status filters by subscriber status.
	Status SubscriberStatus
}

// --- Activity Types ---

// ActivityEventType represents the type of an activity event.
type ActivityEventType string

const (
	ActivityEventQueued       ActivityEventType = "queued"
	ActivityEventProcessing   ActivityEventType = "processing"
	ActivityEventSent         ActivityEventType = "sent"
	ActivityEventDelivered    ActivityEventType = "delivered"
	ActivityEventOpened       ActivityEventType = "opened"
	ActivityEventClicked      ActivityEventType = "clicked"
	ActivityEventBounced      ActivityEventType = "bounced"
	ActivityEventComplained   ActivityEventType = "complained"
	ActivityEventUnsubscribed ActivityEventType = "unsubscribed"
	ActivityEventFailed       ActivityEventType = "failed"
	ActivityEventCancelled    ActivityEventType = "cancelled"
)

// ActivityEmailRef contains reference information about the email associated with an activity event.
type ActivityEmailRef struct {
	ID         string  `json:"id"`
	Subject    string  `json:"subject"`
	From       string  `json:"from"`
	CampaignID *string `json:"campaignId"`
}

// ActivityEvent represents an event in a subscriber's activity timeline.
type ActivityEvent struct {
	ID        string           `json:"id"`
	Type      string           `json:"type"`
	Timestamp string           `json:"timestamp"`
	Data      interface{}      `json:"data"`
	Email     ActivityEmailRef `json:"email"`
}

// ListActivityParams contains filters for listing activity events.
type ListActivityParams struct {
	ListParams
	// Type filters by event type.
	Type ActivityEventType
}

// --- Campaign Types ---

// CampaignStatus represents the status of a campaign.
type CampaignStatus string

const (
	CampaignStatusDraft     CampaignStatus = "draft"
	CampaignStatusScheduled CampaignStatus = "scheduled"
	CampaignStatusSending   CampaignStatus = "sending"
	CampaignStatusSent      CampaignStatus = "sent"
	CampaignStatusPaused    CampaignStatus = "paused"
	CampaignStatusCancelled CampaignStatus = "cancelled"
)

// EmailAddress represents a named email address.
type EmailAddress struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

// CreateCampaignParams contains the parameters for creating a campaign.
type CreateCampaignParams struct {
	Name string `json:"name"`
	// Subject is the email subject line.
	Subject string `json:"subject"`
	// From is the sender address. Use FormatEmailAddress() for "Name <email>" format.
	From string `json:"from"`
	// ReplyTo is the reply-to address.
	ReplyTo string `json:"replyTo,omitempty"`
	// AudienceID is the target audience for the campaign.
	AudienceID string `json:"audienceId"`
	// TemplateID is the template to use (or provide HTML/Text directly).
	TemplateID string `json:"templateId,omitempty"`
	// HTML is the HTML content (alternative to TemplateID).
	HTML string `json:"html,omitempty"`
	// Text is the plain text content.
	Text string `json:"text,omitempty"`
	// Variables are template variable values.
	Variables map[string]interface{} `json:"variables,omitempty"`
	// PreviewText is the email preview text.
	PreviewText string `json:"previewText,omitempty"`
	// Tags are categorization tags.
	Tags []string `json:"tags,omitempty"`
}

// UpdateCampaignParams contains the parameters for updating a campaign.
type UpdateCampaignParams struct {
	Name        *string                `json:"name,omitempty"`
	Subject     *string                `json:"subject,omitempty"`
	From        *string                `json:"from,omitempty"`
	ReplyTo     *string                `json:"replyTo,omitempty"`
	AudienceID  *string                `json:"audienceId,omitempty"`
	TemplateID  *string                `json:"templateId,omitempty"`
	HTML        *string                `json:"html,omitempty"`
	Text        *string                `json:"text,omitempty"`
	Variables   map[string]interface{} `json:"variables,omitempty"`
	PreviewText *string                `json:"previewText,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
}

// ScheduleCampaignParams contains the parameters for scheduling a campaign.
type ScheduleCampaignParams struct {
	// ScheduledAt is the ISO 8601 time to send the campaign.
	ScheduledAt string `json:"scheduledAt"`
}

// CampaignStats contains delivery statistics for a campaign.
type CampaignStats struct {
	Total        int `json:"total"`
	Sent         int `json:"sent"`
	Delivered    int `json:"delivered"`
	Opened       int `json:"opened"`
	Clicked      int `json:"clicked"`
	Bounced      int `json:"bounced"`
	Complained   int `json:"complained"`
	Unsubscribed int `json:"unsubscribed"`
}

// Campaign represents an email campaign.
type Campaign struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Subject     string         `json:"subject"`
	From        EmailAddress   `json:"from"`
	ReplyTo     *EmailAddress  `json:"replyTo"`
	AudienceID  string         `json:"audienceId"`
	TemplateID  *string        `json:"templateId"`
	HTML        *string        `json:"html"`
	Text        *string        `json:"text"`
	PreviewText *string        `json:"previewText"`
	Status      CampaignStatus `json:"status"`
	Tags        []string       `json:"tags"`
	Stats       CampaignStats  `json:"stats"`
	ScheduledAt *string        `json:"scheduledAt"`
	SentAt      *string        `json:"sentAt"`
	CreatedAt   string         `json:"createdAt"`
	UpdatedAt   string         `json:"updatedAt"`
}

// --- Webhook Types ---

// WebhookEvent represents a webhook event type.
type WebhookEvent string

const (
	WebhookEventEmailSent              WebhookEvent = "email.sent"
	WebhookEventEmailDelivered         WebhookEvent = "email.delivered"
	WebhookEventEmailOpened            WebhookEvent = "email.opened"
	WebhookEventEmailClicked           WebhookEvent = "email.clicked"
	WebhookEventEmailBounced           WebhookEvent = "email.bounced"
	WebhookEventEmailComplained        WebhookEvent = "email.complained"
	WebhookEventEmailFailed            WebhookEvent = "email.failed"
	WebhookEventSubscriberAdded        WebhookEvent = "subscriber.added"
	WebhookEventSubscriberRemoved      WebhookEvent = "subscriber.removed"
	WebhookEventSubscriberUnsubscribed WebhookEvent = "subscriber.unsubscribed"
	WebhookEventCampaignSent           WebhookEvent = "campaign.sent"
	WebhookEventCampaignCompleted      WebhookEvent = "campaign.completed"
)

// CreateWebhookParams contains the parameters for creating a webhook.
type CreateWebhookParams struct {
	URL         string         `json:"url"`
	Events      []WebhookEvent `json:"events"`
	Description string         `json:"description,omitempty"`
	Enabled     *bool          `json:"enabled,omitempty"`
}

// UpdateWebhookParams contains the parameters for updating a webhook.
type UpdateWebhookParams struct {
	URL         *string        `json:"url,omitempty"`
	Events      []WebhookEvent `json:"events,omitempty"`
	Description *string        `json:"description,omitempty"`
	Enabled     *bool          `json:"enabled,omitempty"`
}

// Webhook represents a webhook configuration.
type Webhook struct {
	ID              string         `json:"id"`
	URL             string         `json:"url"`
	Events          []WebhookEvent `json:"events"`
	Description     *string        `json:"description"`
	Enabled         bool           `json:"enabled"`
	Secret          string         `json:"secret"`
	LastTriggeredAt *string        `json:"lastTriggeredAt"`
	CreatedAt       string         `json:"createdAt"`
	UpdatedAt       string         `json:"updatedAt"`
}

// WebhookTestResult contains the result of testing a webhook endpoint.
type WebhookTestResult struct {
	Success      bool   `json:"success"`
	StatusCode   int    `json:"statusCode"`
	ResponseTime int    `json:"responseTime"`
	Error        string `json:"error,omitempty"`
}

// --- Topic Types ---

// CreateTopicParams contains the parameters for creating a subscription topic.
type CreateTopicParams struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IsDefault   bool   `json:"isDefault,omitempty"`
	SortOrder   int    `json:"sortOrder,omitempty"`
}

// UpdateTopicParams contains the parameters for updating a subscription topic.
type UpdateTopicParams struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	IsDefault   *bool   `json:"isDefault,omitempty"`
	SortOrder   *int    `json:"sortOrder,omitempty"`
	Active      *bool   `json:"active,omitempty"`
}

// Topic represents a subscription topic.
type Topic struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Slug            string  `json:"slug"`
	Description     *string `json:"description"`
	IsDefault       bool    `json:"isDefault"`
	SortOrder       int     `json:"sortOrder"`
	Active          bool    `json:"active"`
	SubscriberCount int     `json:"subscriberCount"`
	EmailCount      int     `json:"emailCount"`
	CreatedAt       string  `json:"createdAt"`
	UpdatedAt       string  `json:"updatedAt"`
}

// TopicPreference represents a subscriber's preference for a topic.
type TopicPreference struct {
	TopicID          string  `json:"topicId"`
	TopicName        string  `json:"topicName"`
	TopicSlug        string  `json:"topicSlug"`
	TopicDescription *string `json:"topicDescription"`
	Subscribed       bool    `json:"subscribed"`
}

// SetTopicPreferencesParams contains the parameters for setting topic preferences.
type SetTopicPreferencesParams struct {
	Topics []TopicPreferenceSetting `json:"topics"`
}

// TopicPreferenceSetting represents a single topic subscription setting.
type TopicPreferenceSetting struct {
	TopicID    string `json:"topicId"`
	Subscribed bool   `json:"subscribed"`
}

// --- Contact Property Types ---

// ContactPropertyType represents the data type of a contact property.
type ContactPropertyType string

const (
	ContactPropertyTypeText    ContactPropertyType = "text"
	ContactPropertyTypeNumber  ContactPropertyType = "number"
	ContactPropertyTypeDate    ContactPropertyType = "date"
	ContactPropertyTypeBoolean ContactPropertyType = "boolean"
	ContactPropertyTypeEnum    ContactPropertyType = "enum"
)

// CreateContactPropertyParams contains the parameters for creating a contact property.
type CreateContactPropertyParams struct {
	// Key is the unique identifier (alphanumeric + underscore).
	Key string `json:"key"`
	// Name is the display name.
	Name string `json:"name"`
	// Type is the data type (default: "text").
	Type ContactPropertyType `json:"type,omitempty"`
	// Description is an optional description.
	Description string `json:"description,omitempty"`
	// Required indicates whether this property is required for subscribers.
	Required bool `json:"required,omitempty"`
	// EnumOptions is the list of allowed values for ENUM type properties.
	EnumOptions []string `json:"enumOptions,omitempty"`
	// SortOrder controls the display order.
	SortOrder int `json:"sortOrder,omitempty"`
}

// UpdateContactPropertyParams contains the parameters for updating a contact property.
type UpdateContactPropertyParams struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Required    *bool    `json:"required,omitempty"`
	EnumOptions []string `json:"enumOptions,omitempty"`
	SortOrder   *int     `json:"sortOrder,omitempty"`
	Active      *bool    `json:"active,omitempty"`
}

// ContactProperty represents a contact property definition.
type ContactProperty struct {
	ID          string              `json:"id"`
	Key         string              `json:"key"`
	Name        string              `json:"name"`
	Type        ContactPropertyType `json:"type"`
	Description *string             `json:"description"`
	Required    bool                `json:"required"`
	EnumOptions []string            `json:"enumOptions"`
	SortOrder   int                 `json:"sortOrder"`
	Active      bool                `json:"active"`
	CreatedAt   string              `json:"createdAt"`
	UpdatedAt   string              `json:"updatedAt"`
}

// SubscriberPropertyValue represents a single property value for a subscriber.
type SubscriberPropertyValue struct {
	Value string `json:"value"`
	Type  string `json:"type"`
	Name  string `json:"name"`
}

// SubscriberPropertyValues is a map of property keys to their values for a subscriber.
type SubscriberPropertyValues map[string]SubscriberPropertyValue

// --- Sequence Types ---

// CreateSequenceParams contains the parameters for creating an automation sequence.
type CreateSequenceParams struct {
	Name          string                 `json:"name"`
	AudienceID    string                 `json:"audienceId"`
	TriggerType   string                 `json:"triggerType"`
	Description   string                 `json:"description,omitempty"`
	TriggerConfig map[string]interface{} `json:"triggerConfig,omitempty"`
}

// UpdateSequenceParams contains the parameters for updating a sequence.
type UpdateSequenceParams struct {
	Name          *string                `json:"name,omitempty"`
	Description   *string                `json:"description,omitempty"`
	TriggerType   *string                `json:"triggerType,omitempty"`
	TriggerConfig map[string]interface{} `json:"triggerConfig,omitempty"`
}

// Sequence represents an automation sequence.
type Sequence struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Description   *string                `json:"description"`
	AudienceID    string                 `json:"audienceId"`
	TriggerType   string                 `json:"triggerType"`
	TriggerConfig map[string]interface{} `json:"triggerConfig"`
	Status        string                 `json:"status"`
	Steps         []SequenceStep         `json:"steps"`
	EnrollCount   int                    `json:"enrollCount"`
	CreatedAt     string                 `json:"createdAt"`
	UpdatedAt     string                 `json:"updatedAt"`
}

// SequenceStep represents a step in a sequence.
type SequenceStep struct {
	ID              string                 `json:"id"`
	SequenceID      string                 `json:"sequenceId"`
	Position        int                    `json:"position"`
	Type            string                 `json:"type"`
	Subject         *string                `json:"subject"`
	HTML            *string                `json:"html"`
	Text            *string                `json:"text"`
	TemplateID      *string                `json:"templateId"`
	DelayAmount     *int                   `json:"delayAmount"`
	DelayUnit       *string                `json:"delayUnit"`
	ConditionType   *string                `json:"conditionType"`
	ConditionConfig map[string]interface{} `json:"conditionConfig"`
	CreatedAt       string                 `json:"createdAt"`
	UpdatedAt       string                 `json:"updatedAt"`
}

// CreateSequenceStepParams contains the parameters for adding a sequence step.
type CreateSequenceStepParams struct {
	Position        int                    `json:"position"`
	Type            string                 `json:"type"`
	Subject         string                 `json:"subject,omitempty"`
	HTML            string                 `json:"html,omitempty"`
	Text            string                 `json:"text,omitempty"`
	TemplateID      string                 `json:"templateId,omitempty"`
	DelayAmount     int                    `json:"delayAmount,omitempty"`
	DelayUnit       string                 `json:"delayUnit,omitempty"`
	ConditionType   string                 `json:"conditionType,omitempty"`
	ConditionConfig map[string]interface{} `json:"conditionConfig,omitempty"`
}

// UpdateSequenceStepParams contains the parameters for updating a sequence step.
type UpdateSequenceStepParams struct {
	Subject         *string                `json:"subject,omitempty"`
	HTML            *string                `json:"html,omitempty"`
	Text            *string                `json:"text,omitempty"`
	TemplateID      *string                `json:"templateId,omitempty"`
	DelayAmount     *int                   `json:"delayAmount,omitempty"`
	DelayUnit       *string                `json:"delayUnit,omitempty"`
	ConditionType   *string                `json:"conditionType,omitempty"`
	ConditionConfig map[string]interface{} `json:"conditionConfig,omitempty"`
}

// StepPosition represents a step ID and its target position for reordering.
type StepPosition struct {
	ID       string `json:"id"`
	Position int    `json:"position"`
}

// EnrollResult contains the result of enrolling subscribers.
type EnrollResult struct {
	Enrolled int `json:"enrolled"`
}

// SequenceEnrollment represents a subscriber enrollment in a sequence.
type SequenceEnrollment struct {
	ID           string  `json:"id"`
	SequenceID   string  `json:"sequenceId"`
	SubscriberID string  `json:"subscriberId"`
	Status       string  `json:"status"`
	CurrentStep  int     `json:"currentStep"`
	EnrolledAt   string  `json:"enrolledAt"`
	CompletedAt  *string `json:"completedAt"`
}

// --- RSS Feed Types ---

// CreateRssFeedParams contains the parameters for creating an RSS feed.
type CreateRssFeedParams struct {
	Name            string `json:"name"`
	URL             string `json:"url"`
	AudienceID      string `json:"audienceId"`
	PollInterval    string `json:"pollInterval,omitempty"`
	Mode            string `json:"mode,omitempty"`
	SubjectTemplate string `json:"subjectTemplate,omitempty"`
	HTMLTemplate    string `json:"htmlTemplate,omitempty"`
}

// UpdateRssFeedParams contains the parameters for updating an RSS feed.
type UpdateRssFeedParams struct {
	Name            *string `json:"name,omitempty"`
	URL             *string `json:"url,omitempty"`
	PollInterval    *string `json:"pollInterval,omitempty"`
	Mode            *string `json:"mode,omitempty"`
	SubjectTemplate *string `json:"subjectTemplate,omitempty"`
	HTMLTemplate    *string `json:"htmlTemplate,omitempty"`
}

// RssFeed represents an RSS feed configuration.
type RssFeed struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	URL             string  `json:"url"`
	AudienceID      string  `json:"audienceId"`
	PollInterval    string  `json:"pollInterval"`
	Mode            string  `json:"mode"`
	SubjectTemplate *string `json:"subjectTemplate"`
	HTMLTemplate    *string `json:"htmlTemplate"`
	Status          string  `json:"status"`
	LastPolledAt    *string `json:"lastPolledAt"`
	CreatedAt       string  `json:"createdAt"`
	UpdatedAt       string  `json:"updatedAt"`
}

// RssFeedListResponse is the response for listing feeds.
type RssFeedListResponse struct {
	Data []RssFeed `json:"data"`
}

// RssFeedDetail includes recent items.
type RssFeedDetail struct {
	RssFeed
	RecentItems []RssFeedItem `json:"recentItems"`
}

// RssFeedItem represents a single RSS feed item.
type RssFeedItem struct {
	ID          string  `json:"id"`
	FeedID      string  `json:"feedId"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Link        string  `json:"link"`
	Author      *string `json:"author"`
	PublishedAt string  `json:"publishedAt"`
	Processed   bool    `json:"processed"`
	CreatedAt   string  `json:"createdAt"`
}

// FeedPollResult contains the result of a manual poll.
type FeedPollResult struct {
	Success  bool `json:"success"`
	NewItems int  `json:"newItems"`
}

// ListFeedItemsParams contains parameters for listing feed items.
type ListFeedItemsParams struct {
	ListParams
	Processed *bool
}

// --- Form Types ---

// CreateFormParams contains the parameters for creating a signup form.
type CreateFormParams struct {
	Name        string                   `json:"name"`
	AudienceID  string                   `json:"audienceId"`
	Fields      []map[string]interface{} `json:"fields,omitempty"`
	DoubleOptIn *bool                    `json:"doubleOptIn,omitempty"`
	RedirectURL string                   `json:"redirectUrl,omitempty"`
	Honeypot    *bool                    `json:"honeypot,omitempty"`
	CASLConsent *bool                    `json:"caslConsent,omitempty"`
}

// UpdateFormParams contains the parameters for updating a form.
type UpdateFormParams struct {
	Name        *string                  `json:"name,omitempty"`
	Fields      []map[string]interface{} `json:"fields,omitempty"`
	DoubleOptIn *bool                    `json:"doubleOptIn,omitempty"`
	RedirectURL *string                  `json:"redirectUrl,omitempty"`
	Honeypot    *bool                    `json:"honeypot,omitempty"`
	CASLConsent *bool                    `json:"caslConsent,omitempty"`
}

// Form represents a signup form.
type Form struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	AudienceID  string                   `json:"audienceId"`
	Fields      []map[string]interface{} `json:"fields"`
	DoubleOptIn bool                     `json:"doubleOptIn"`
	RedirectURL *string                  `json:"redirectUrl"`
	Honeypot    bool                     `json:"honeypot"`
	CASLConsent bool                     `json:"caslConsent"`
	EmbedCode   string                   `json:"embedCode"`
	CreatedAt   string                   `json:"createdAt"`
	UpdatedAt   string                   `json:"updatedAt"`
}

// --- Analytics Types ---

// GeoAnalyticsParams contains parameters for geo analytics queries.
type GeoAnalyticsParams struct {
	Days      int    `json:"days,omitempty"`
	EventType string `json:"eventType,omitempty"`
}

// DeviceAnalyticsParams contains parameters for device analytics queries.
type DeviceAnalyticsParams struct {
	Days      int    `json:"days,omitempty"`
	EventType string `json:"eventType,omitempty"`
}

// CampaignAnalyticsParams contains parameters for campaign-level analytics.
type CampaignAnalyticsParams struct {
	EventType string `json:"eventType,omitempty"`
}

// GeoAnalyticsEntry represents a single geo analytics entry.
type GeoAnalyticsEntry struct {
	CountryCode string  `json:"countryCode"`
	CountryName string  `json:"countryName"`
	City        *string `json:"city"`
	Region      *string `json:"region"`
	Count       int     `json:"count"`
	Percentage  float64 `json:"percentage"`
}

// GeoAnalyticsResponse is the response for geo analytics.
type GeoAnalyticsResponse struct {
	Data       []GeoAnalyticsEntry `json:"data"`
	Total      int                 `json:"total"`
	EventType  string              `json:"eventType"`
	DateRange  DateRange           `json:"dateRange"`
}

// DateRange represents a start/end date range.
type DateRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// DeviceAnalyticsResponse is the response for device analytics.
type DeviceAnalyticsResponse struct {
	DeviceTypes      []AnalyticsBreakdown `json:"deviceTypes"`
	OperatingSystems []AnalyticsBreakdown `json:"operatingSystems"`
	Browsers         []AnalyticsBreakdown `json:"browsers"`
	EmailClients     []AnalyticsBreakdown `json:"emailClients"`
	BotCount         int                  `json:"botCount"`
	MachineOpenCount int                  `json:"machineOpenCount"`
	Total            int                  `json:"total"`
	EventType        string               `json:"eventType"`
	DateRange        DateRange            `json:"dateRange"`
}

// AnalyticsBreakdown represents a name/count/percentage analytics entry.
type AnalyticsBreakdown struct {
	Name       string  `json:"name"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

// --- Campaign Extension Types ---

// SendTestParams contains the parameters for sending a campaign test.
type SendTestParams struct {
	To []string `json:"to"`
}

// SendTestResult contains the result of a campaign test send.
type SendTestResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// CloneCampaignOptions contains options for cloning a campaign.
type CloneCampaignOptions struct {
	IncludeABTest bool `json:"includeABTest,omitempty"`
}

// LinkStats represents tracked link analytics for a campaign or email.
type LinkStats struct {
	Data  []LinkEntry `json:"data"`
	Total int         `json:"total"`
}

// LinkEntry represents a single tracked link.
type LinkEntry struct {
	URL          string  `json:"url"`
	DisplayURL   string  `json:"displayUrl"`
	TotalClicks  int     `json:"totalClicks"`
	UniqueClicks int     `json:"uniqueClicks"`
	ClickRate    float64 `json:"clickRate"`
}

// ListLinksParams contains parameters for listing tracked links.
type ListLinksParams struct {
	Limit int    `json:"limit,omitempty"`
	Sort  string `json:"sort,omitempty"`
	Order string `json:"order,omitempty"`
}

// RecalculateEngagementResult contains the result of engagement recalculation.
type RecalculateEngagementResult struct {
	Processed int `json:"processed"`
}

// EngagementStats contains engagement statistics for an audience.
type EngagementStats struct {
	AverageScore float64                `json:"averageScore"`
	Distribution map[string]int         `json:"distribution"`
	Total        int                    `json:"total"`
}
