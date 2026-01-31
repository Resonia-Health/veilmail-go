package veilmail

import (
	"context"
	"fmt"
	"net/url"
)

// TopicsService handles subscription topic operations.
type TopicsService struct {
	client *httpClient
}

// Create creates a new subscription topic.
func (s *TopicsService) Create(ctx context.Context, params *CreateTopicParams) (*Topic, error) {
	var topic Topic
	err := s.client.post(ctx, "/v1/topics", params, &topic)
	if err != nil {
		return nil, err
	}
	return &topic, nil
}

// List returns all topics, optionally filtered by active status.
func (s *TopicsService) List(ctx context.Context, active *bool) ([]Topic, error) {
	q := url.Values{}
	if active != nil {
		if *active {
			q.Set("active", "true")
		} else {
			q.Set("active", "false")
		}
	}
	var resp dataWrapper[[]Topic]
	err := s.client.get(ctx, "/v1/topics", q, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Get retrieves a single topic by ID.
func (s *TopicsService) Get(ctx context.Context, id string) (*Topic, error) {
	var topic Topic
	err := s.client.get(ctx, fmt.Sprintf("/v1/topics/%s", id), nil, &topic)
	if err != nil {
		return nil, err
	}
	return &topic, nil
}

// Update updates a topic.
func (s *TopicsService) Update(ctx context.Context, id string, params *UpdateTopicParams) (*Topic, error) {
	var topic Topic
	err := s.client.patch(ctx, fmt.Sprintf("/v1/topics/%s", id), params, &topic)
	if err != nil {
		return nil, err
	}
	return &topic, nil
}

// Delete deactivates a topic (soft delete).
func (s *TopicsService) Delete(ctx context.Context, id string) error {
	return s.client.del(ctx, fmt.Sprintf("/v1/topics/%s", id))
}

// GetPreferences retrieves a subscriber's topic preferences.
func (s *TopicsService) GetPreferences(ctx context.Context, audienceID, subscriberID string) ([]TopicPreference, error) {
	path := fmt.Sprintf("/v1/audiences/%s/subscribers/%s/topics", audienceID, subscriberID)
	var resp dataWrapper[[]TopicPreference]
	err := s.client.get(ctx, path, nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// SetPreferences sets a subscriber's topic preferences.
func (s *TopicsService) SetPreferences(ctx context.Context, audienceID, subscriberID string, params *SetTopicPreferencesParams) ([]TopicPreference, error) {
	path := fmt.Sprintf("/v1/audiences/%s/subscribers/%s/topics", audienceID, subscriberID)
	var resp dataWrapper[[]TopicPreference]
	err := s.client.put(ctx, path, params, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}
