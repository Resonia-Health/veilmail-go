package veilmail

import (
	"context"
	"fmt"
	"net/url"
)

// PropertiesService handles contact property operations.
type PropertiesService struct {
	client *httpClient
}

// Create creates a new contact property definition.
func (s *PropertiesService) Create(ctx context.Context, params *CreateContactPropertyParams) (*ContactProperty, error) {
	var prop ContactProperty
	err := s.client.post(ctx, "/v1/properties", params, &prop)
	if err != nil {
		return nil, err
	}
	return &prop, nil
}

// List returns all contact properties, optionally filtered by active status.
func (s *PropertiesService) List(ctx context.Context, active *bool) ([]ContactProperty, error) {
	q := url.Values{}
	if active != nil {
		if *active {
			q.Set("active", "true")
		} else {
			q.Set("active", "false")
		}
	}
	var resp dataWrapper[[]ContactProperty]
	err := s.client.get(ctx, "/v1/properties", q, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Get retrieves a single contact property by ID.
func (s *PropertiesService) Get(ctx context.Context, id string) (*ContactProperty, error) {
	var prop ContactProperty
	err := s.client.get(ctx, fmt.Sprintf("/v1/properties/%s", id), nil, &prop)
	if err != nil {
		return nil, err
	}
	return &prop, nil
}

// Update updates a contact property definition.
func (s *PropertiesService) Update(ctx context.Context, id string, params *UpdateContactPropertyParams) (*ContactProperty, error) {
	var prop ContactProperty
	err := s.client.patch(ctx, fmt.Sprintf("/v1/properties/%s", id), params, &prop)
	if err != nil {
		return nil, err
	}
	return &prop, nil
}

// Delete deactivates a contact property (soft delete).
func (s *PropertiesService) Delete(ctx context.Context, id string) error {
	return s.client.del(ctx, fmt.Sprintf("/v1/properties/%s", id))
}

// GetValues retrieves a subscriber's property values.
func (s *PropertiesService) GetValues(ctx context.Context, audienceID, subscriberID string) (SubscriberPropertyValues, error) {
	path := fmt.Sprintf("/v1/audiences/%s/subscribers/%s/properties", audienceID, subscriberID)
	var resp dataWrapper[SubscriberPropertyValues]
	err := s.client.get(ctx, path, nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// SetValues sets property values for a subscriber.
// Pass nil for a value to delete it.
func (s *PropertiesService) SetValues(ctx context.Context, audienceID, subscriberID string, values map[string]interface{}) error {
	path := fmt.Sprintf("/v1/audiences/%s/subscribers/%s/properties", audienceID, subscriberID)
	body := map[string]interface{}{
		"values": values,
	}
	var resp struct {
		Success bool `json:"success"`
	}
	return s.client.put(ctx, path, body, &resp)
}
