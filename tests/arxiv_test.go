//go:build integration

package tests

import (
	"context"
	"testing"

	"github.com/tushariitr-19/research-agent-adk/logger"
	"github.com/tushariitr-19/research-agent-adk/tools"
)

func TestArxivFetchByID(t *testing.T) {
	if err := logger.Init(false); err != nil {
		t.Fatalf("failed to init logger: %v", err)
	}

	client := tools.NewArxivClient()

	paper, err := client.FetchByID(context.Background(), "1706.03762")
	if err != nil {
		t.Fatalf("FetchByID failed: %v", err)
	}

	if paper.Title == "" {
		t.Error("expected title to be non-empty")
	}
	if paper.ID != "1706.03762" {
		t.Errorf("expected ID 1706.03762 got %s", paper.ID)
	}
	if len(paper.Authors) == 0 {
		t.Error("expected authors to be non-empty")
	}
	if paper.Abstract == "" {
		t.Error("expected abstract to be non-empty")
	}
	if paper.PDFURL == "" {
		t.Error("expected PDF URL to be non-empty")
	}

	t.Logf("title: %s", paper.Title)
	t.Logf("authors: %d", len(paper.Authors))
	t.Logf("categories: %v", paper.Categories)
	t.Logf("pdf: %s", paper.PDFURL)
}
