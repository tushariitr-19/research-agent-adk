//go:build integration

package tests

import (
	"context"
	"testing"

	"github.com/tushariitr-19/research-agent-adk/logger"
	"github.com/tushariitr-19/research-agent-adk/tools"
)

func TestOpenAlexSearchRelatedPapers(t *testing.T) {
	if err := logger.Init(false); err != nil {
		t.Fatalf("failed to init logger: %v", err)
	}

	client := tools.NewOpenAlexClient()

	papers, err := client.SearchRelatedPapers(context.Background(), "transformer attention mechanism", 3)
	if err != nil {
		t.Fatalf("SearchRelatedPapers failed: %v", err)
	}

	if len(papers) == 0 {
		t.Error("expected at least one paper")
	}

	for _, p := range papers {
		if p.Title == "" {
			t.Error("expected title to be non-empty")
		}
		if p.Source != "openalex" {
			t.Errorf("expected source openalex got %s", p.Source)
		}
	}

	t.Logf("found %d related papers", len(papers))
	t.Logf("first paper: %s (%s)", papers[0].Title, papers[0].Published)
}

func TestOpenAlexSearchAuthorPapers(t *testing.T) {
	if err := logger.Init(false); err != nil {
		t.Fatalf("failed to init logger: %v", err)
	}

	client := tools.NewOpenAlexClient()

	papers, err := client.SearchAuthorPapers(context.Background(), "Ashish Vaswani", 3)
	if err != nil {
		t.Fatalf("SearchAuthorPapers failed: %v", err)
	}

	if len(papers) == 0 {
		t.Error("expected at least one paper")
	}

	for _, p := range papers {
		if p.Title == "" {
			t.Error("expected title to be non-empty")
		}
	}

	t.Logf("found %d papers for Ashish Vaswani", len(papers))
	t.Logf("first paper: %s (%s)", papers[0].Title, papers[0].Published)
}

func TestOpenAlexConcurrentRequests(t *testing.T) {
	if err := logger.Init(false); err != nil {
		t.Fatalf("failed to init logger: %v", err)
	}

	client := tools.NewOpenAlexClient()

	// Test that concurrent requests work without rate limiting issues
	done := make(chan error, 2)

	go func() {
		_, err := client.SearchAuthorPapers(context.Background(), "Ashish Vaswani", 2)
		done <- err
	}()

	go func() {
		_, err := client.SearchAuthorPapers(context.Background(), "Noam Shazeer", 2)
		done <- err
	}()

	for i := 0; i < 2; i++ {
		if err := <-done; err != nil {
			t.Errorf("concurrent request failed: %v", err)
		}
	}

	t.Log("concurrent requests completed successfully")
}
