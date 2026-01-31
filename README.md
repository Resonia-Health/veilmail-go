# Veil Mail Go SDK

Official Go SDK for the [Veil Mail](https://veilmail.xyz) API. Send emails with built-in PII protection using idiomatic Go.

## Requirements

- Go 1.21 or later
- Veil Mail API key ([get one here](https://veilmail.xyz))

## Installation

```bash
go get github.com/Resonia-Health/veilmail-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    veilmail "github.com/Resonia-Health/veilmail-go"
)

func main() {
    client, err := veilmail.New("veil_live_your_api_key")
    if err != nil {
        log.Fatal(err)
    }

    email, err := client.Emails.Send(context.Background(), &veilmail.SendEmailParams{
        From:    "hello@yourdomain.com",
        To:      []string{"user@example.com"},
        Subject: "Hello from Veil Mail",
        HTML:    "<p>Welcome to Veil Mail!</p>",
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Email sent: %s (status: %s)\n", email.ID, email.Status)
}
```

## Configuration

```go
// Simple: just an API key
client, err := veilmail.New("veil_live_xxx")

// With options
client, err := veilmail.New("veil_live_xxx",
    veilmail.WithBaseURL("https://custom-api.example.com"),
    veilmail.WithTimeout(10 * time.Second),
    veilmail.WithHTTPClient(&http.Client{
        Transport: &http.Transport{
            MaxIdleConns: 10,
        },
    }),
)
```

## Emails

```go
ctx := context.Background()

// Send an email
email, err := client.Emails.Send(ctx, &veilmail.SendEmailParams{
    From:    veilmail.FormatEmailAddress("Alice", "alice@yourdomain.com"),
    To:      []string{"bob@example.com"},
    Subject: "Hello",
    HTML:    "<p>Hello Bob!</p>",
    Tags:    []string{"welcome"},
})

// Send with a template
email, err := client.Emails.Send(ctx, &veilmail.SendEmailParams{
    From:       "hello@yourdomain.com",
    To:         []string{"user@example.com"},
    TemplateID: "tmpl_xxx",
    TemplateData: map[string]interface{}{
        "name": "Alice",
    },
})

// Schedule for later
email, err := client.Emails.Send(ctx, &veilmail.SendEmailParams{
    From:         "hello@yourdomain.com",
    To:           []string{"user@example.com"},
    Subject:      "Scheduled Email",
    HTML:         "<p>This was scheduled</p>",
    ScheduledFor: "2025-06-01T09:00:00Z",
})

// Send with attachments
email, err := client.Emails.Send(ctx, &veilmail.SendEmailParams{
    From:    "hello@yourdomain.com",
    To:      []string{"user@example.com"},
    Subject: "Invoice Attached",
    HTML:    "<p>Please find your invoice attached.</p>",
    Attachments: []veilmail.EmailAttachment{
        {
            Filename:    "invoice.pdf",
            Content:     base64EncodedContent,
            ContentType: "application/pdf",
        },
    },
})

// Batch send (up to 100)
result, err := client.Emails.SendBatch(ctx, []veilmail.SendEmailParams{
    {From: "hello@yourdomain.com", To: []string{"user1@example.com"}, Subject: "Hello", HTML: "<p>Hi!</p>"},
    {From: "hello@yourdomain.com", To: []string{"user2@example.com"}, Subject: "Hello", HTML: "<p>Hi!</p>"},
})
fmt.Printf("Sent: %d, Failed: %d\n", result.Successful, result.Failed)

// List emails
emails, err := client.Emails.List(ctx, &veilmail.ListEmailsParams{
    ListParams: veilmail.ListParams{Limit: 10},
    Status:     veilmail.EmailStatusDelivered,
})

// Get email details
email, err := client.Emails.Get(ctx, "email_xxx")

// Cancel a scheduled email
result, err := client.Emails.Cancel(ctx, "email_xxx")

// Reschedule a scheduled email
email, err := client.Emails.Update(ctx, "email_xxx", &veilmail.UpdateEmailParams{
    ScheduledFor: "2025-07-01T09:00:00Z",
})
```

## Domains

```go
// Add a domain
domain, err := client.Domains.Create(ctx, &veilmail.CreateDomainParams{
    Domain: "mail.example.com",
})
// domain.DnsRecords contains the DNS records to configure

// Verify DNS records
domain, err = client.Domains.Verify(ctx, domain.ID)

// List domains
domains, err := client.Domains.List(ctx, nil)

// Update tracking settings
trueVal := true
domain, err = client.Domains.Update(ctx, domain.ID, &veilmail.UpdateDomainParams{
    TrackOpens:  &trueVal,
    TrackClicks: &trueVal,
})

// Delete a domain
err = client.Domains.Delete(ctx, domain.ID)
```

## Templates

```go
// Create a template
tmpl, err := client.Templates.Create(ctx, &veilmail.CreateTemplateParams{
    Name:    "Welcome Email",
    Subject: "Welcome, {{name}}!",
    HTML:    "<h1>Hello {{name}}</h1><p>Welcome aboard!</p>",
    Variables: []veilmail.TemplateVariable{
        {Name: "name", Type: "string", Required: true},
    },
})

// Preview with variables
preview, err := client.Templates.Preview(ctx, &veilmail.PreviewTemplateParams{
    HTML: "<h1>Hello {{name}}</h1>",
    Variables: map[string]interface{}{
        "name": "Alice",
    },
})

// List, get, update, delete
templates, err := client.Templates.List(ctx, nil)
tmpl, err = client.Templates.Get(ctx, "tmpl_xxx")
tmpl, err = client.Templates.Update(ctx, "tmpl_xxx", &veilmail.UpdateTemplateParams{...})
err = client.Templates.Delete(ctx, "tmpl_xxx")
```

## Audiences & Subscribers

```go
// Create an audience
audience, err := client.Audiences.Create(ctx, &veilmail.CreateAudienceParams{
    Name:        "Newsletter",
    Description: "Weekly newsletter subscribers",
})

// Add a subscriber
subs := client.Audiences.Subscribers(audience.ID)
subscriber, err := subs.Add(ctx, &veilmail.AddSubscriberParams{
    Email:     "user@example.com",
    FirstName: "Alice",
    LastName:  "Smith",
})

// Add with double opt-in
subscriber, err = subs.Add(ctx, &veilmail.AddSubscriberParams{
    Email:       "user@example.com",
    DoubleOptIn: true,
})
// Later: confirm the subscription
subscriber, err = subs.Confirm(ctx, subscriber.ID)

// List subscribers
subscribers, err := subs.List(ctx, &veilmail.ListSubscribersParams{
    ListParams: veilmail.ListParams{Limit: 50},
    Status:     veilmail.SubscriberStatusActive,
})

// Bulk import
result, err := subs.Import(ctx, &veilmail.ImportSubscribersParams{
    Subscribers: []veilmail.ImportSubscriberEntry{
        {Email: "user1@example.com", FirstName: "Alice"},
        {Email: "user2@example.com", FirstName: "Bob"},
    },
})

// Export as CSV
csv, err := subs.Export(ctx, nil)

// Activity timeline
activity, err := subs.Activity(ctx, subscriber.ID, &veilmail.ListActivityParams{
    ListParams: veilmail.ListParams{Limit: 20},
})
```

## Campaigns

```go
// Create a campaign
campaign, err := client.Campaigns.Create(ctx, &veilmail.CreateCampaignParams{
    Name:       "Summer Sale",
    Subject:    "Summer Sale - 50% Off!",
    From:       veilmail.FormatEmailAddress("Store", "deals@yourdomain.com"),
    AudienceID: "aud_xxx",
    HTML:       "<h1>Summer Sale!</h1><p>50% off everything</p>",
})

// Schedule for later
campaign, err = client.Campaigns.Schedule(ctx, campaign.ID, &veilmail.ScheduleCampaignParams{
    ScheduledAt: "2025-06-15T10:00:00Z",
})

// Or send immediately
campaign, err = client.Campaigns.Send(ctx, campaign.ID)

// Pause/resume/cancel
campaign, err = client.Campaigns.Pause(ctx, campaign.ID)
campaign, err = client.Campaigns.Resume(ctx, campaign.ID)
campaign, err = client.Campaigns.Cancel(ctx, campaign.ID)
```

## Webhooks

```go
// Create a webhook
webhook, err := client.Webhooks.Create(ctx, &veilmail.CreateWebhookParams{
    URL: "https://yourdomain.com/webhooks/veilmail",
    Events: []veilmail.WebhookEvent{
        veilmail.WebhookEventEmailDelivered,
        veilmail.WebhookEventEmailBounced,
        veilmail.WebhookEventEmailComplained,
    },
})

// Test the endpoint
testResult, err := client.Webhooks.Test(ctx, webhook.ID)

// Rotate the signing secret
webhook, err = client.Webhooks.RotateSecret(ctx, webhook.ID)

// Verify webhook signatures in your handler
func webhookHandler(w http.ResponseWriter, r *http.Request) {
    body, _ := io.ReadAll(r.Body)
    signature := r.Header.Get("X-Signature-Hash")

    if !veilmail.VerifySignature(body, signature, webhookSecret) {
        http.Error(w, "Invalid signature", http.StatusUnauthorized)
        return
    }

    // Process the webhook event
    var event map[string]interface{}
    json.Unmarshal(body, &event)
    // ...
}
```

## Topics

```go
// Create a subscription topic
topic, err := client.Topics.Create(ctx, &veilmail.CreateTopicParams{
    Name:        "Product Updates",
    Description: "New feature announcements",
    IsDefault:   true,
})

// List topics
active := true
topics, err := client.Topics.List(ctx, &active)

// Get subscriber preferences
prefs, err := client.Topics.GetPreferences(ctx, "aud_xxx", "sub_xxx")

// Set preferences
prefs, err = client.Topics.SetPreferences(ctx, "aud_xxx", "sub_xxx", &veilmail.SetTopicPreferencesParams{
    Topics: []veilmail.TopicPreferenceSetting{
        {TopicID: "topic_xxx", Subscribed: true},
        {TopicID: "topic_yyy", Subscribed: false},
    },
})
```

## Contact Properties

```go
// Define a property
prop, err := client.Properties.Create(ctx, &veilmail.CreateContactPropertyParams{
    Key:  "company",
    Name: "Company Name",
    Type: veilmail.ContactPropertyTypeText,
})

// Define an enum property
prop, err = client.Properties.Create(ctx, &veilmail.CreateContactPropertyParams{
    Key:         "plan",
    Name:        "Plan",
    Type:        veilmail.ContactPropertyTypeEnum,
    EnumOptions: []string{"free", "pro", "enterprise"},
})

// List properties
active := true
props, err := client.Properties.List(ctx, &active)

// Set values for a subscriber
err = client.Properties.SetValues(ctx, "aud_xxx", "sub_xxx", map[string]interface{}{
    "company": "Acme Corp",
    "plan":    "pro",
})

// Get values
values, err := client.Properties.GetValues(ctx, "aud_xxx", "sub_xxx")
```

## Error Handling

```go
email, err := client.Emails.Send(ctx, params)
if err != nil {
    if veilmail.IsAuthenticationError(err) {
        log.Fatal("Invalid API key")
    }
    if veilmail.IsValidationError(err) {
        var apiErr *veilmail.Error
        if errors.As(err, &apiErr) {
            log.Printf("Validation error: %s (details: %v)", apiErr.Message, apiErr.Details)
        }
    }
    if veilmail.IsPiiDetectedError(err) {
        var apiErr *veilmail.Error
        if errors.As(err, &apiErr) {
            log.Printf("PII detected: %v", apiErr.PiiTypes)
        }
    }
    if veilmail.IsRateLimitError(err) {
        var apiErr *veilmail.Error
        if errors.As(err, &apiErr) {
            log.Printf("Rate limited, retry after %d seconds", apiErr.RetryAfter)
        }
    }
    if veilmail.IsNotFoundError(err) {
        log.Println("Resource not found")
    }
    if veilmail.IsServerError(err) {
        log.Println("Server error, try again later")
    }
    log.Fatal(err)
}
```

## License

MIT
