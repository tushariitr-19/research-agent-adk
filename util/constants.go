package util

import "time"

const (
	// HTTP client
	DefaultHTTPTimeout = 30 * time.Second
	DefaultMinGap      = 8 * time.Second
	UserAgent          = "research-agent-adk/1.0 (github.com/tushariitr-19/research-agent-adk)"
	MaxRetries         = 3
	RetryBaseDelay     = 10 * time.Second

	// ArXiv API
	ArxivBaseURL = "https://export.arxiv.org/api/query"
	ArxivSource  = "arxiv"
)
