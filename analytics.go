package veilmail

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// AnalyticsService handles geo and device analytics operations.
type AnalyticsService struct {
	client *httpClient
}

// Geo returns organization-level geo analytics.
func (s *AnalyticsService) Geo(ctx context.Context, params *GeoAnalyticsParams) (*GeoAnalyticsResponse, error) {
	q := url.Values{}
	if params != nil {
		if params.Days > 0 {
			q.Set("days", strconv.Itoa(params.Days))
		}
		if params.EventType != "" {
			q.Set("eventType", params.EventType)
		}
	}
	var resp GeoAnalyticsResponse
	err := s.client.get(ctx, "/v1/analytics/geo", q, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Devices returns organization-level device analytics.
func (s *AnalyticsService) Devices(ctx context.Context, params *DeviceAnalyticsParams) (*DeviceAnalyticsResponse, error) {
	q := url.Values{}
	if params != nil {
		if params.Days > 0 {
			q.Set("days", strconv.Itoa(params.Days))
		}
		if params.EventType != "" {
			q.Set("eventType", params.EventType)
		}
	}
	var resp DeviceAnalyticsResponse
	err := s.client.get(ctx, "/v1/analytics/devices", q, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// CampaignGeo returns campaign-level geo analytics.
func (s *AnalyticsService) CampaignGeo(ctx context.Context, campaignID string, params *CampaignAnalyticsParams) (*GeoAnalyticsResponse, error) {
	q := url.Values{}
	if params != nil && params.EventType != "" {
		q.Set("eventType", params.EventType)
	}
	var resp GeoAnalyticsResponse
	err := s.client.get(ctx, fmt.Sprintf("/v1/campaigns/%s/analytics/geo", campaignID), q, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// CampaignDevices returns campaign-level device analytics.
func (s *AnalyticsService) CampaignDevices(ctx context.Context, campaignID string, params *CampaignAnalyticsParams) (*DeviceAnalyticsResponse, error) {
	q := url.Values{}
	if params != nil && params.EventType != "" {
		q.Set("eventType", params.EventType)
	}
	var resp DeviceAnalyticsResponse
	err := s.client.get(ctx, fmt.Sprintf("/v1/campaigns/%s/analytics/devices", campaignID), q, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
