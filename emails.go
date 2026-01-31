package veilmail

import (
	"context"
	"fmt"
	"net/url"
)

// EmailsService handles email operations.
type EmailsService struct {
	client *httpClient
}

// Send sends a single email.
func (s *EmailsService) Send(ctx context.Context, params *SendEmailParams) (*Email, error) {
	var email Email
	err := s.client.post(ctx, "/v1/emails", params, &email)
	if err != nil {
		return nil, err
	}
	return &email, nil
}

// SendBatch sends a batch of up to 100 emails.
// Returns a result with per-email status including any errors.
func (s *EmailsService) SendBatch(ctx context.Context, emails []SendEmailParams) (*BatchSendResult, error) {
	body := map[string]interface{}{
		"emails": emails,
	}
	var result BatchSendResult
	err := s.client.post(ctx, "/v1/emails/batch", body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// List returns a paginated list of emails.
func (s *EmailsService) List(ctx context.Context, params *ListEmailsParams) (*ListResponse[Email], error) {
	q := url.Values{}
	if params != nil {
		addPaginationParams(q, &params.ListParams)
		if params.Status != "" {
			q.Set("status", string(params.Status))
		}
		if params.Tag != "" {
			q.Set("tag", params.Tag)
		}
		if params.After != "" {
			q.Set("after", params.After)
		}
		if params.Before != "" {
			q.Set("before", params.Before)
		}
	}
	var resp ListResponse[Email]
	err := s.client.get(ctx, "/v1/emails", q, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves a single email by ID.
func (s *EmailsService) Get(ctx context.Context, id string) (*Email, error) {
	var email Email
	err := s.client.get(ctx, fmt.Sprintf("/v1/emails/%s", id), nil, &email)
	if err != nil {
		return nil, err
	}
	return &email, nil
}

// Cancel cancels a scheduled email.
func (s *EmailsService) Cancel(ctx context.Context, id string) (*CancelEmailResult, error) {
	var result CancelEmailResult
	err := s.client.post(ctx, fmt.Sprintf("/v1/emails/%s/cancel", id), nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Update reschedules a scheduled email.
func (s *EmailsService) Update(ctx context.Context, id string, params *UpdateEmailParams) (*Email, error) {
	var email Email
	err := s.client.patch(ctx, fmt.Sprintf("/v1/emails/%s", id), params, &email)
	if err != nil {
		return nil, err
	}
	return &email, nil
}

// Links returns tracked link analytics for a specific email.
func (s *EmailsService) Links(ctx context.Context, id string, params *ListLinksParams) (*LinkStats, error) {
	q := url.Values{}
	if params != nil {
		if params.Limit > 0 {
			q.Set("limit", fmt.Sprintf("%d", params.Limit))
		}
		if params.Sort != "" {
			q.Set("sort", params.Sort)
		}
		if params.Order != "" {
			q.Set("order", params.Order)
		}
	}
	var resp LinkStats
	err := s.client.get(ctx, fmt.Sprintf("/v1/emails/%s/links", id), q, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
