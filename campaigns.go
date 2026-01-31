package veilmail

import (
	"context"
	"fmt"
	"net/url"
)

// CampaignsService handles campaign operations.
type CampaignsService struct {
	client *httpClient
}

// Create creates a new campaign in draft status.
func (s *CampaignsService) Create(ctx context.Context, params *CreateCampaignParams) (*Campaign, error) {
	var resp dataWrapper[Campaign]
	err := s.client.post(ctx, "/v1/campaigns", params, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// List returns a paginated list of campaigns.
func (s *CampaignsService) List(ctx context.Context, params *ListParams) (*ListResponse[Campaign], error) {
	q := url.Values{}
	addPaginationParams(q, params)
	var resp ListResponse[Campaign]
	err := s.client.get(ctx, "/v1/campaigns", q, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves a single campaign by ID.
func (s *CampaignsService) Get(ctx context.Context, id string) (*Campaign, error) {
	var resp dataWrapper[Campaign]
	err := s.client.get(ctx, fmt.Sprintf("/v1/campaigns/%s", id), nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Update updates a draft campaign.
func (s *CampaignsService) Update(ctx context.Context, id string, params *UpdateCampaignParams) (*Campaign, error) {
	var resp dataWrapper[Campaign]
	err := s.client.patch(ctx, fmt.Sprintf("/v1/campaigns/%s", id), params, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Delete removes a campaign.
func (s *CampaignsService) Delete(ctx context.Context, id string) error {
	return s.client.del(ctx, fmt.Sprintf("/v1/campaigns/%s", id))
}

// Schedule schedules a campaign for future delivery.
func (s *CampaignsService) Schedule(ctx context.Context, id string, params *ScheduleCampaignParams) (*Campaign, error) {
	var resp dataWrapper[Campaign]
	err := s.client.post(ctx, fmt.Sprintf("/v1/campaigns/%s/schedule", id), params, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Send sends a campaign immediately.
func (s *CampaignsService) Send(ctx context.Context, id string) (*Campaign, error) {
	var resp dataWrapper[Campaign]
	err := s.client.post(ctx, fmt.Sprintf("/v1/campaigns/%s/send", id), nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Pause pauses an in-progress campaign.
func (s *CampaignsService) Pause(ctx context.Context, id string) (*Campaign, error) {
	var resp dataWrapper[Campaign]
	err := s.client.post(ctx, fmt.Sprintf("/v1/campaigns/%s/pause", id), nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Resume resumes a paused campaign.
func (s *CampaignsService) Resume(ctx context.Context, id string) (*Campaign, error) {
	var resp dataWrapper[Campaign]
	err := s.client.post(ctx, fmt.Sprintf("/v1/campaigns/%s/resume", id), nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Cancel cancels a scheduled or in-progress campaign.
func (s *CampaignsService) Cancel(ctx context.Context, id string) (*Campaign, error) {
	var resp dataWrapper[Campaign]
	err := s.client.post(ctx, fmt.Sprintf("/v1/campaigns/%s/cancel", id), nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// SendTest sends a test/preview of a campaign to 1-5 email addresses.
func (s *CampaignsService) SendTest(ctx context.Context, id string, params *SendTestParams) (*SendTestResult, error) {
	var resp SendTestResult
	err := s.client.post(ctx, fmt.Sprintf("/v1/campaigns/%s/test", id), params, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Clone clones a campaign as a new draft.
func (s *CampaignsService) Clone(ctx context.Context, id string, opts *CloneCampaignOptions) (*Campaign, error) {
	var body interface{}
	if opts != nil {
		body = opts
	} else {
		body = map[string]interface{}{}
	}
	var resp Campaign
	err := s.client.post(ctx, fmt.Sprintf("/v1/campaigns/%s/clone", id), body, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Links returns tracked link analytics for a campaign.
func (s *CampaignsService) Links(ctx context.Context, id string, params *ListLinksParams) (*LinkStats, error) {
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
	err := s.client.get(ctx, fmt.Sprintf("/v1/campaigns/%s/links", id), q, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
