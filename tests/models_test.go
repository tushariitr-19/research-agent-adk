package tests

import (
	"testing"

	"github.com/tushariitr-19/research-agent-adk/models"
)

func TestPaperModel(t *testing.T) {
	paper := models.Paper{
		ID:              "1706.03762",
		Title:           "Attention Is All You Need",
		Abstract:        "We propose a new simple network architecture...",
		Authors:         []models.Author{{Name: "Ashish Vaswani"}},
		Categories:      []string{"cs.CL", "cs.LG"},
		PrimaryCategory: "cs.CL",
		Published:       "2017-06-12",
		Source:          "arxiv",
	}

	if paper.ID == "" {
		t.Error("expected ID to be non-empty")
	}
	if paper.Title == "" {
		t.Error("expected Title to be non-empty")
	}
	if len(paper.Authors) == 0 {
		t.Error("expected Authors to be non-empty")
	}
	if len(paper.Categories) == 0 {
		t.Error("expected Categories to be non-empty")
	}
	if paper.Abstract == "" {
		t.Error("expected Abstract to be non-empty")
	}
	if paper.PrimaryCategory == "" {
		t.Error("expected PrimaryCategory to be non-empty")
	}
	if paper.Published == "" {
		t.Error("expected Published to be non-empty")
	}
	if paper.Source == "" {
		t.Error("expected Source to be non-empty")
	}

	t.Logf("paper: %s by %s", paper.Title, paper.Authors[0].Name)
}

func TestAuthorModel(t *testing.T) {
	author := models.Author{
		Name:        "Ashish Vaswani",
		Affiliation: "Google Brain",
	}

	if author.Name == "" {
		t.Error("expected Name to be non-empty")
	}

	t.Logf("author: %s (%s)", author.Name, author.Affiliation)
}

func TestSearchResultModel(t *testing.T) {
	result := models.SearchResult{
		Query:      "transformer architecture",
		TotalFound: 100,
		Papers:     []models.Paper{{Title: "Attention Is All You Need"}},
		Source:     "arxiv",
	}

	if result.Query == "" {
		t.Error("expected Query to be non-empty")
	}
	if len(result.Papers) == 0 {
		t.Error("expected Papers to be non-empty")
	}
	if result.TotalFound == 0 {
		t.Error("expected TotalFound to be non-zero")
	}
	if result.Source == "" {
		t.Error("expected Source to be non-empty")
	}

	t.Logf("search result: %d papers found for '%s'", result.TotalFound, result.Query)
}
