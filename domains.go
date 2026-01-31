package veilmail

import (
	"context"
	"fmt"
	"net/url"
)

// DomainsService handles domain operations.
type DomainsService struct {
	client *httpClient
}

// Create adds a new sending domain and returns DNS records for verification.
func (s *DomainsService) Create(ctx context.Context, params *CreateDomainParams) (*Domain, error) {
	var resp dataWrapper[Domain]
	err := s.client.post(ctx, "/v1/domains", params, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// List returns a paginated list of domains.
func (s *DomainsService) List(ctx context.Context, params *ListParams) (*ListResponse[Domain], error) {
	q := url.Values{}
	addPaginationParams(q, params)
	var resp ListResponse[Domain]
	err := s.client.get(ctx, "/v1/domains", q, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves a single domain by ID.
func (s *DomainsService) Get(ctx context.Context, id string) (*Domain, error) {
	var resp dataWrapper[Domain]
	err := s.client.get(ctx, fmt.Sprintf("/v1/domains/%s", id), nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Update updates domain settings such as tracking configuration.
func (s *DomainsService) Update(ctx context.Context, id string, params *UpdateDomainParams) (*Domain, error) {
	var resp dataWrapper[Domain]
	err := s.client.patch(ctx, fmt.Sprintf("/v1/domains/%s", id), params, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Verify triggers DNS verification for a domain.
func (s *DomainsService) Verify(ctx context.Context, id string) (*Domain, error) {
	var resp dataWrapper[Domain]
	err := s.client.post(ctx, fmt.Sprintf("/v1/domains/%s/verify", id), nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Delete removes a domain.
func (s *DomainsService) Delete(ctx context.Context, id string) error {
	return s.client.del(ctx, fmt.Sprintf("/v1/domains/%s", id))
}
