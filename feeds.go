package veilmail

import (
	"context"
	"fmt"
	"net/url"
)

// FeedsService handles RSS feed operations.
type FeedsService struct {
	client *httpClient
}

// Create creates a new RSS feed.
func (s *FeedsService) Create(ctx context.Context, params *CreateRssFeedParams) (*RssFeed, error) {
	var resp RssFeed
	err := s.client.post(ctx, "/v1/feeds", params, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// List returns all RSS feeds.
func (s *FeedsService) List(ctx context.Context) (*RssFeedListResponse, error) {
	var resp RssFeedListResponse
	err := s.client.get(ctx, "/v1/feeds", nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves a single feed by ID (includes recent items).
func (s *FeedsService) Get(ctx context.Context, id string) (*RssFeedDetail, error) {
	var resp RssFeedDetail
	err := s.client.get(ctx, fmt.Sprintf("/v1/feeds/%s", id), nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Update updates a feed.
func (s *FeedsService) Update(ctx context.Context, id string, params *UpdateRssFeedParams) (*RssFeed, error) {
	var resp RssFeed
	err := s.client.put(ctx, fmt.Sprintf("/v1/feeds/%s", id), params, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Delete deletes a feed and all its items.
func (s *FeedsService) Delete(ctx context.Context, id string) error {
	return s.client.del(ctx, fmt.Sprintf("/v1/feeds/%s", id))
}

// Poll manually triggers a feed poll.
func (s *FeedsService) Poll(ctx context.Context, id string) (*FeedPollResult, error) {
	var resp FeedPollResult
	err := s.client.post(ctx, fmt.Sprintf("/v1/feeds/%s/poll", id), nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Pause pauses an active feed.
func (s *FeedsService) Pause(ctx context.Context, id string) (*RssFeed, error) {
	var resp RssFeed
	err := s.client.post(ctx, fmt.Sprintf("/v1/feeds/%s/pause", id), nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Resume resumes a paused or errored feed.
func (s *FeedsService) Resume(ctx context.Context, id string) (*RssFeed, error) {
	var resp RssFeed
	err := s.client.post(ctx, fmt.Sprintf("/v1/feeds/%s/resume", id), nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListItems returns a paginated list of feed items.
func (s *FeedsService) ListItems(ctx context.Context, feedID string, params *ListFeedItemsParams) (*ListResponse[RssFeedItem], error) {
	q := url.Values{}
	if params != nil {
		addPaginationParams(q, &params.ListParams)
		if params.Processed != nil {
			if *params.Processed {
				q.Set("processed", "true")
			} else {
				q.Set("processed", "false")
			}
		}
	}
	var resp ListResponse[RssFeedItem]
	err := s.client.get(ctx, fmt.Sprintf("/v1/feeds/%s/items", feedID), q, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
