package veilmail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type httpClient struct {
	apiKey    string
	baseURL   string
	client    *http.Client
	userAgent string
}

type requestOptions struct {
	body    interface{}
	query   url.Values
	headers http.Header
}

func (h *httpClient) do(ctx context.Context, method, path string, opts *requestOptions, result interface{}) error {
	u := h.baseURL + path

	if opts != nil && opts.query != nil {
		encoded := opts.query.Encode()
		if encoded != "" {
			u += "?" + encoded
		}
	}

	var bodyReader io.Reader
	if opts != nil && opts.body != nil {
		b, err := json.Marshal(opts.body)
		if err != nil {
			return fmt.Errorf("veilmail: failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, bodyReader)
	if err != nil {
		return fmt.Errorf("veilmail: failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+h.apiKey)
	req.Header.Set("User-Agent", h.userAgent)
	if bodyReader != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if opts != nil && opts.headers != nil {
		for key, values := range opts.headers {
			for _, v := range values {
				req.Header.Set(key, v)
			}
		}
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("veilmail: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("veilmail: failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return parseAPIError(resp.StatusCode, resp.Header, body)
	}

	if result != nil && len(body) > 0 {
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("veilmail: failed to parse response: %w", err)
		}
	}

	return nil
}

// doRaw performs an HTTP request and returns the raw response body.
// Used for endpoints that return non-JSON responses (e.g., CSV export).
func (h *httpClient) doRaw(ctx context.Context, method, path string, opts *requestOptions) ([]byte, error) {
	u := h.baseURL + path

	if opts != nil && opts.query != nil {
		encoded := opts.query.Encode()
		if encoded != "" {
			u += "?" + encoded
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u, nil)
	if err != nil {
		return nil, fmt.Errorf("veilmail: failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+h.apiKey)
	req.Header.Set("User-Agent", h.userAgent)

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("veilmail: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("veilmail: failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, parseAPIError(resp.StatusCode, resp.Header, body)
	}

	return body, nil
}

func (h *httpClient) get(ctx context.Context, path string, query url.Values, result interface{}) error {
	return h.do(ctx, http.MethodGet, path, &requestOptions{query: query}, result)
}

func (h *httpClient) post(ctx context.Context, path string, body interface{}, result interface{}) error {
	return h.do(ctx, http.MethodPost, path, &requestOptions{body: body}, result)
}

func (h *httpClient) patch(ctx context.Context, path string, body interface{}, result interface{}) error {
	return h.do(ctx, http.MethodPatch, path, &requestOptions{body: body}, result)
}

func (h *httpClient) put(ctx context.Context, path string, body interface{}, result interface{}) error {
	return h.do(ctx, http.MethodPut, path, &requestOptions{body: body}, result)
}

func (h *httpClient) del(ctx context.Context, path string) error {
	return h.do(ctx, http.MethodDelete, path, nil, nil)
}

// addPaginationParams adds standard pagination query parameters.
func addPaginationParams(q url.Values, params *ListParams) {
	if params == nil {
		return
	}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Cursor != "" {
		q.Set("cursor", params.Cursor)
	}
}
