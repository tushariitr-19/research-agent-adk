package util

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
)

// OpenAlexHTTPClient is a simple HTTP client for OpenAlex
// No rate limiting needed — OpenAlex supports concurrent requests
type OpenAlexHTTPClient struct {
	client *http.Client
}

var (
	defaultOpenAlexClient *OpenAlexHTTPClient
	openAlexOnce          sync.Once
)

// GetOpenAlexHTTPClient returns the singleton OpenAlex HTTP client
func GetOpenAlexHTTPClient() *OpenAlexHTTPClient {
	openAlexOnce.Do(func() {
		defaultOpenAlexClient = &OpenAlexHTTPClient{
			client: &http.Client{Timeout: DefaultHTTPTimeout},
		}
	})
	return defaultOpenAlexClient
}

// Get performs a simple HTTP GET request
func (c *OpenAlexHTTPClient) Get(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", UserAgent)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return body, nil
}
