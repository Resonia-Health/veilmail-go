package veilmail

import (
	"context"
	"fmt"
	"net/url"
)

// AudiencesService handles audience operations.
type AudiencesService struct {
	client *httpClient
}

// Create creates a new audience.
func (s *AudiencesService) Create(ctx context.Context, params *CreateAudienceParams) (*Audience, error) {
	var resp dataWrapper[Audience]
	err := s.client.post(ctx, "/v1/audiences", params, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// List returns a paginated list of audiences.
func (s *AudiencesService) List(ctx context.Context, params *ListParams) (*ListResponse[Audience], error) {
	q := url.Values{}
	addPaginationParams(q, params)
	var resp ListResponse[Audience]
	err := s.client.get(ctx, "/v1/audiences", q, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves a single audience by ID.
func (s *AudiencesService) Get(ctx context.Context, id string) (*Audience, error) {
	var resp dataWrapper[Audience]
	err := s.client.get(ctx, fmt.Sprintf("/v1/audiences/%s", id), nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Update updates an existing audience.
func (s *AudiencesService) Update(ctx context.Context, id string, params *UpdateAudienceParams) (*Audience, error) {
	var resp dataWrapper[Audience]
	err := s.client.put(ctx, fmt.Sprintf("/v1/audiences/%s", id), params, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Delete removes an audience.
func (s *AudiencesService) Delete(ctx context.Context, id string) error {
	return s.client.del(ctx, fmt.Sprintf("/v1/audiences/%s", id))
}

// Subscribers returns a SubscribersService scoped to the given audience.
func (s *AudiencesService) Subscribers(audienceID string) *SubscribersService {
	return &SubscribersService{client: s.client, audienceID: audienceID}
}

// RecalculateEngagement recalculates engagement scores for all subscribers.
func (s *AudiencesService) RecalculateEngagement(ctx context.Context, audienceID string) (*RecalculateEngagementResult, error) {
	var resp RecalculateEngagementResult
	err := s.client.post(ctx, fmt.Sprintf("/v1/audiences/%s/recalculate-engagement", audienceID), map[string]interface{}{}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetEngagementStats returns engagement statistics for an audience.
func (s *AudiencesService) GetEngagementStats(ctx context.Context, audienceID string) (*EngagementStats, error) {
	var resp EngagementStats
	err := s.client.get(ctx, fmt.Sprintf("/v1/audiences/%s/engagement-stats", audienceID), nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// SubscribersService handles subscriber operations within an audience.
type SubscribersService struct {
	client     *httpClient
	audienceID string
}

func (s *SubscribersService) basePath() string {
	return fmt.Sprintf("/v1/audiences/%s/subscribers", s.audienceID)
}

// List returns a paginated list of subscribers.
func (s *SubscribersService) List(ctx context.Context, params *ListSubscribersParams) (*ListResponse[Subscriber], error) {
	q := url.Values{}
	if params != nil {
		addPaginationParams(q, &params.ListParams)
		if params.Status != "" {
			q.Set("status", string(params.Status))
		}
		if params.Email != "" {
			q.Set("email", params.Email)
		}
	}
	var resp ListResponse[Subscriber]
	err := s.client.get(ctx, s.basePath(), q, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Add adds a subscriber to the audience.
func (s *SubscribersService) Add(ctx context.Context, params *AddSubscriberParams) (*Subscriber, error) {
	var resp dataWrapper[Subscriber]
	err := s.client.post(ctx, s.basePath(), params, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Get retrieves a single subscriber by ID.
func (s *SubscribersService) Get(ctx context.Context, subscriberID string) (*Subscriber, error) {
	var resp dataWrapper[Subscriber]
	err := s.client.get(ctx, fmt.Sprintf("%s/%s", s.basePath(), subscriberID), nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Update updates a subscriber.
func (s *SubscribersService) Update(ctx context.Context, subscriberID string, params *UpdateSubscriberParams) (*Subscriber, error) {
	var resp dataWrapper[Subscriber]
	err := s.client.put(ctx, fmt.Sprintf("%s/%s", s.basePath(), subscriberID), params, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Remove removes a subscriber from the audience.
func (s *SubscribersService) Remove(ctx context.Context, subscriberID string) error {
	return s.client.del(ctx, fmt.Sprintf("%s/%s", s.basePath(), subscriberID))
}

// Confirm completes the double opt-in process for a subscriber.
func (s *SubscribersService) Confirm(ctx context.Context, subscriberID string) (*Subscriber, error) {
	var resp dataWrapper[Subscriber]
	err := s.client.post(ctx, fmt.Sprintf("%s/%s/confirm", s.basePath(), subscriberID), nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Import performs a bulk import of subscribers from JSON or CSV data.
func (s *SubscribersService) Import(ctx context.Context, params *ImportSubscribersParams) (*ImportSubscribersResult, error) {
	var result ImportSubscribersResult
	err := s.client.post(ctx, fmt.Sprintf("%s/import", s.basePath()), params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Export exports subscribers as CSV data.
func (s *SubscribersService) Export(ctx context.Context, params *ExportSubscribersParams) (string, error) {
	q := url.Values{}
	if params != nil && params.Status != "" {
		q.Set("status", string(params.Status))
	}
	body, err := s.client.doRaw(ctx, "GET", fmt.Sprintf("%s/export", s.basePath()), &requestOptions{query: q})
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// Activity returns a paginated timeline of events for a subscriber.
func (s *SubscribersService) Activity(ctx context.Context, subscriberID string, params *ListActivityParams) (*ListResponse[ActivityEvent], error) {
	q := url.Values{}
	if params != nil {
		addPaginationParams(q, &params.ListParams)
		if params.Type != "" {
			q.Set("type", string(params.Type))
		}
	}
	var resp ListResponse[ActivityEvent]
	err := s.client.get(ctx, fmt.Sprintf("%s/%s/activity", s.basePath(), subscriberID), q, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
