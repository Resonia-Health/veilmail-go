package veilmail

import (
	"context"
	"fmt"
	"net/url"
)

// FormsService handles signup form operations.
type FormsService struct {
	client *httpClient
}

// Create creates a new signup form.
func (s *FormsService) Create(ctx context.Context, params *CreateFormParams) (*Form, error) {
	var resp Form
	err := s.client.post(ctx, "/v1/forms", params, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// List returns a paginated list of forms.
func (s *FormsService) List(ctx context.Context, params *ListParams) (*ListResponse[Form], error) {
	q := url.Values{}
	addPaginationParams(q, params)
	var resp ListResponse[Form]
	err := s.client.get(ctx, "/v1/forms", q, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves a single form by ID.
func (s *FormsService) Get(ctx context.Context, id string) (*Form, error) {
	var resp Form
	err := s.client.get(ctx, fmt.Sprintf("/v1/forms/%s", id), nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Update updates a form.
func (s *FormsService) Update(ctx context.Context, id string, params *UpdateFormParams) (*Form, error) {
	var resp Form
	err := s.client.put(ctx, fmt.Sprintf("/v1/forms/%s", id), params, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Delete removes a form.
func (s *FormsService) Delete(ctx context.Context, id string) error {
	return s.client.del(ctx, fmt.Sprintf("/v1/forms/%s", id))
}
