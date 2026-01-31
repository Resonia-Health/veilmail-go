package veilmail

import (
	"context"
	"fmt"
	"net/url"
)

// TemplatesService handles template operations.
type TemplatesService struct {
	client *httpClient
}

// Create creates a new email template.
func (s *TemplatesService) Create(ctx context.Context, params *CreateTemplateParams) (*Template, error) {
	var resp dataWrapper[Template]
	err := s.client.post(ctx, "/v1/templates", params, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// List returns a paginated list of templates.
func (s *TemplatesService) List(ctx context.Context, params *ListParams) (*ListResponse[Template], error) {
	q := url.Values{}
	addPaginationParams(q, params)
	var resp ListResponse[Template]
	err := s.client.get(ctx, "/v1/templates", q, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves a single template by ID.
func (s *TemplatesService) Get(ctx context.Context, id string) (*Template, error) {
	var resp dataWrapper[Template]
	err := s.client.get(ctx, fmt.Sprintf("/v1/templates/%s", id), nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Update updates an existing template.
func (s *TemplatesService) Update(ctx context.Context, id string, params *UpdateTemplateParams) (*Template, error) {
	var resp dataWrapper[Template]
	err := s.client.patch(ctx, fmt.Sprintf("/v1/templates/%s", id), params, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Preview renders a template with the given variables without sending.
func (s *TemplatesService) Preview(ctx context.Context, params *PreviewTemplateParams) (*PreviewTemplateResponse, error) {
	var resp PreviewTemplateResponse
	err := s.client.post(ctx, "/v1/templates/preview", params, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Delete removes a template.
func (s *TemplatesService) Delete(ctx context.Context, id string) error {
	return s.client.del(ctx, fmt.Sprintf("/v1/templates/%s", id))
}
