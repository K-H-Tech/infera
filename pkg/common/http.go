package common

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// RequestOptions represents optional configuration for HTTP requests
type RequestOptions func(*http.Request)

// WithHeader adds a custom header to the request
func WithHeader(key, value string) RequestOptions {
	return func(req *http.Request) {
		req.Header.Set(key, value)
	}
}

// WithQueryParam adds a query parameter to the URL (for GET requests)
func WithQueryParam(key, value string) RequestOptions {
	return func(req *http.Request) {
		// Parse existing query parameters
		parsedURL, err := url.Parse(req.URL.String())
		if err != nil {
			return // Skip if URL parsing fails
		}

		params := parsedURL.Query()
		params.Add(key, value)
		parsedURL.RawQuery = params.Encode()
		req.URL = parsedURL
	}
}

// WithQueryParams adds multiple query parameters to the URL (for GET requests)
func WithQueryParams(params map[string]string) RequestOptions {
	return func(req *http.Request) {
		parsedURL, err := url.Parse(req.URL.String())
		if err != nil {
			return
		}

		query := parsedURL.Query()
		for key, value := range params {
			query.Set(key, value)
		}
		parsedURL.RawQuery = query.Encode()
		req.URL = parsedURL
	}
}

// HTTPClient wraps the HTTP module functionality
type HTTPClient struct {
	client *http.Client
}

// NewHTTPClient creates a new HTTP module
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: &http.Client{},
	}
}

// NewHTTPClientWithTimeout creates a new HTTP module with a default timeout
func NewHTTPClientWithTimeout(timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Post makes an HTTP POST request
// url and payload are mandatory, options are optional
func (c *HTTPClient) Post(url string, payload interface{}, options ...RequestOptions) (*http.Response, error) {
	if url == "" {
		return nil, fmt.Errorf("URL is required")
	}

	var body io.Reader
	var err error

	switch p := payload.(type) {
	case []byte:
		body = bytes.NewBuffer(p)
	case string:
		body = bytes.NewBufferString(p)
	case io.Reader:
		body = p
	default:
		payloadBytes, marshalErr := json.Marshal(payload)
		if marshalErr != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", marshalErr)
		}
		body = bytes.NewBuffer(payloadBytes)
	}

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %w", err)
	}

	// Set default Content-Type if not already set
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Apply options
	for _, option := range options {
		option(req)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("POST request failed: %w", err)
	}

	return resp, nil
}

// Get makes an HTTP GET request
// url is mandatory, options are optional (for headers, query params, etc.)
func (c *HTTPClient) Get(url string, options ...RequestOptions) (*http.Response, error) {
	if url == "" {
		return nil, fmt.Errorf("URL is required")
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %w", err)
	}

	// Apply options (headers, query params, etc.)
	for _, option := range options {
		option(req)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET request failed: %w", err)
	}

	return resp, nil
}

// PostWithContext makes an HTTP POST request with context
func (c *HTTPClient) PostWithContext(ctx context.Context, url string, payload interface{}, options ...RequestOptions) (*http.Response, error) {
	if url == "" {
		return nil, fmt.Errorf("URL is required")
	}

	var body io.Reader
	var err error

	switch p := payload.(type) {
	case []byte:
		body = bytes.NewBuffer(p)
	case string:
		body = bytes.NewBufferString(p)
	case io.Reader:
		body = p
	default:
		payloadBytes, marshalErr := json.Marshal(payload)
		if marshalErr != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", marshalErr)
		}
		body = bytes.NewBuffer(payloadBytes)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request with context: %w", err)
	}

	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	for _, option := range options {
		option(req)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("POST request with context failed: %w", err)
	}

	return resp, nil
}

// GetWithContext makes an HTTP GET request with context
func (c *HTTPClient) GetWithContext(ctx context.Context, url string, options ...RequestOptions) (*http.Response, error) {
	if url == "" {
		return nil, fmt.Errorf("URL is required")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request with context: %w", err)
	}

	for _, option := range options {
		option(req)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET request with context failed: %w", err)
	}

	return resp, nil
}

// Standalone functions for convenience (using default module)

// Post makes a POST request using a default module
func Post(url string, payload interface{}, options ...RequestOptions) (*http.Response, error) {
	client := NewHTTPClient()
	return client.Post(url, payload, options...)
}

// Get makes a GET request using a default module
func Get(url string, options ...RequestOptions) (*http.Response, error) {
	client := NewHTTPClient()
	return client.Get(url, options...)
}

// PostWithContext makes a POST request with context using a default module
func PostWithContext(ctx context.Context, url string, payload interface{}, options ...RequestOptions) (*http.Response, error) {
	client := NewHTTPClient()
	return client.PostWithContext(ctx, url, payload, options...)
}

// GetWithContext makes a GET request with context using a default module
func GetWithContext(ctx context.Context, url string, options ...RequestOptions) (*http.Response, error) {
	client := NewHTTPClient()
	return client.GetWithContext(ctx, url, options...)
}

// Helper function to read response body
func ReadResponseBody(resp *http.Response) ([]byte, error) {
	if resp == nil {
		return nil, fmt.Errorf("response is nil")
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	return io.ReadAll(resp.Body)
}
