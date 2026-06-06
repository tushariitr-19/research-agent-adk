package util

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/tushariitr-19/research-agent-adk/logger"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
)

// HTTPClient is a production-grade HTTP client with rate limiting,
// retry with exponential backoff, and User-Agent header injection.
type HTTPClient struct {
	client      *http.Client
	mu          sync.Mutex
	lastRequest time.Time
	minGap      time.Duration
	userAgent   string
	sem         *semaphore.Weighted // limits concurrent requests to 1
}

var (
	defaultClient *HTTPClient
	once          sync.Once
)

// GetHTTPClient returns the singleton HTTPClient instance.
func GetHTTPClient() *HTTPClient {
	once.Do(func() {
		defaultClient = &HTTPClient{
			client:    &http.Client{Timeout: DefaultHTTPTimeout},
			minGap:    DefaultMinGap,
			userAgent: UserAgent,
			sem:       semaphore.NewWeighted(1), // only 1 concurrent request
		}
	})
	return defaultClient
}

// Get performs a rate-limited, retried HTTP GET request.
func (c *HTTPClient) Get(ctx context.Context, url string) ([]byte, error) {
	var (
		body []byte
		err  error
	)

	for attempt := 0; attempt <= MaxRetries; attempt++ {
		if attempt > 0 {
			delay := RetryBaseDelay * time.Duration(attempt)
			logger.Log.Warn("retrying request",
				zap.String("url", url),
				zap.Int("attempt", attempt),
				zap.Duration("delay", delay),
			)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		body, err = c.doGet(ctx, url)
		if err == nil {
			return body, nil
		}

		// Do not retry on context cancellation
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		logger.Log.Error("request failed",
			zap.String("url", url),
			zap.Int("attempt", attempt),
			zap.Error(err),
		)
	}

	return nil, fmt.Errorf("request failed after %d attempts: %w", MaxRetries, err)
}

// doGet performs a single rate-limited HTTP GET request.
func (c *HTTPClient) doGet(ctx context.Context, url string) ([]byte, error) {
	// Acquire semaphore — only 1 request at a time globally
	if err := c.sem.Acquire(ctx, 1); err != nil {
		return nil, fmt.Errorf("context cancelled: %w", err)
	}
	defer c.sem.Release(1)

	// Enforce rate limiting
	c.mu.Lock()
	elapsed := time.Since(c.lastRequest)
	if elapsed < c.minGap {
		waitDuration := c.minGap - elapsed
		c.mu.Unlock()
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(waitDuration):
		}
		c.mu.Lock()
	}
	c.lastRequest = time.Now()
	c.mu.Unlock()

	// Build request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/xml, text/xml, */*")

	// Execute request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Handle non-200 status codes
	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("rate limited by server (429) — will retry")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	logger.Log.Debug("request successful",
		zap.String("url", url),
		zap.Int("bytes", len(body)),
	)

	return body, nil
}
