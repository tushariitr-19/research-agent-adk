package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	adktool "google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"

	"github.com/tushariitr-19/research-agent-adk/logger"
	"github.com/tushariitr-19/research-agent-adk/models"
	"github.com/tushariitr-19/research-agent-adk/util"
	"go.uber.org/zap"
)

// OpenAlexClient implements author search using OpenAlex API
type OpenAlexClient struct {
	http *util.OpenAlexHTTPClient
}

// NewOpenAlexClient creates a new OpenAlex client
func NewOpenAlexClient() *OpenAlexClient {
	return &OpenAlexClient{
		http: util.GetOpenAlexHTTPClient(),
	}
}

// SearchAuthorPapers searches for papers by author name
func (c *OpenAlexClient) SearchAuthorPapers(ctx context.Context, authorName string, limit int) ([]models.Paper, error) {
	logger.Log.Info("searching author papers", zap.String("author", authorName), zap.Int("limit", limit))

	if limit <= 0 {
		limit = openAlexDefaultAuthorLimit
	}

	params := url.Values{}
	params.Set("filter", fmt.Sprintf("raw_author_name.search:%s", authorName))
	params.Set("per-page", fmt.Sprintf("%d", limit))
	params.Set("select", "title,publication_year,doi,authorships,abstract_inverted_index")

	reqURL := openAlexBaseURL + "/works?" + params.Encode()

	data, err := c.http.Get(ctx, reqURL)
	if err != nil {
		return nil, fmt.Errorf("OpenAlex request failed: %w", err)
	}

	var result struct {
		Results []struct {
			Title           string `json:"title"`
			PublicationYear int    `json:"publication_year"`
			DOI             string `json:"doi"`
			Authorships     []struct {
				Author struct {
					DisplayName string `json:"display_name"`
				} `json:"author"`
			} `json:"authorships"`
		} `json:"results"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAlex response: %w", err)
	}

	papers := make([]models.Paper, 0, len(result.Results))
	for _, r := range result.Results {
		authors := make([]models.Author, 0, len(r.Authorships))
		for _, a := range r.Authorships {
			authors = append(authors, models.Author{Name: a.Author.DisplayName})
		}

		papers = append(papers, models.Paper{
			Title:     r.Title,
			Published: fmt.Sprintf("%d", r.PublicationYear),
			DOI:       r.DOI,
			Authors:   authors,
			Source:    openAlexSource,
		})
	}

	logger.Log.Info("author papers found", zap.Int("count", len(papers)))
	return papers, nil
}

// ADK Tool

type searchAuthorPapersArgs struct {
	AuthorName string `json:"author_name" jsonschema:"The full name of the author to search for"`
	Limit      int    `json:"limit,omitempty" jsonschema:"Maximum number of papers to return, default 5"`
}

// NewSearchAuthorPapersTool creates an ADK tool for searching author papers
func NewSearchAuthorPapersTool(client *OpenAlexClient) (adktool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "search_author_papers",
			Description: "Search for papers by a specific author using OpenAlex",
		},
		func(ctx adktool.Context, args searchAuthorPapersArgs) (map[string]any, error) {
			logger.Log.Info("tool called: search_author_papers", zap.String("author", args.AuthorName))
			papers, err := client.SearchAuthorPapers(context.Background(), args.AuthorName, args.Limit)
			if err != nil {
				logger.Log.Error("search_author_papers failed", zap.String("author", args.AuthorName), zap.Error(err))
				return nil, err
			}
			data, err := json.Marshal(papers)
			if err != nil {
				logger.Log.Error("search_author_papers: marshal failed", zap.Error(err))
				return nil, err
			}
			return map[string]any{"result": string(data)}, nil
		},
	)
}

type searchRelatedPapersArgs struct {
	Query string `json:"query" jsonschema:"Keywords to search for related papers"`
	Limit int    `json:"limit,omitempty" jsonschema:"Maximum number of papers, default 3"`
}

func NewSearchRelatedPapersTool(client *OpenAlexClient) (adktool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "search_related_papers",
			Description: "Search for related research papers using OpenAlex",
		},
		func(ctx adktool.Context, args searchRelatedPapersArgs) (map[string]any, error) {
			logger.Log.Info("tool called: search_related_papers", zap.String("query", args.Query))
			if args.Limit == 0 {
				args.Limit = openAlexDefaultSearchLimit
			}
			papers, err := client.SearchRelatedPapers(context.Background(), args.Query, args.Limit)
			if err != nil {
				logger.Log.Error("search_related_papers failed", zap.String("query", args.Query), zap.Error(err))
				return nil, err
			}
			data, err := json.Marshal(papers)
			if err != nil {
				logger.Log.Error("search_related_papers: marshal failed", zap.Error(err))
				return nil, err
			}
			return map[string]any{"result": string(data)}, nil
		},
	)
}

func (c *OpenAlexClient) SearchRelatedPapers(ctx context.Context, query string, limit int) ([]models.Paper, error) {
	logger.Log.Info("searching related papers", zap.String("query", query), zap.Int("limit", limit))

	params := url.Values{}
	params.Set("filter", fmt.Sprintf("title.search:%s", query))
	params.Set("per-page", fmt.Sprintf("%d", limit))
	params.Set("select", "title,publication_year,doi,authorships")

	reqURL := openAlexBaseURL + "/works?" + params.Encode()

	data, err := c.http.Get(ctx, reqURL)
	if err != nil {
		return nil, fmt.Errorf("OpenAlex request failed: %w", err)
	}

	var result struct {
		Results []struct {
			Title           string `json:"title"`
			PublicationYear int    `json:"publication_year"`
			DOI             string `json:"doi"`
			Authorships     []struct {
				Author struct {
					DisplayName string `json:"display_name"`
				} `json:"author"`
			} `json:"authorships"`
		} `json:"results"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	papers := make([]models.Paper, 0, len(result.Results))
	for _, r := range result.Results {
		authors := make([]models.Author, 0, len(r.Authorships))
		for _, a := range r.Authorships {
			authors = append(authors, models.Author{Name: a.Author.DisplayName})
		}
		papers = append(papers, models.Paper{
			Title:     r.Title,
			Published: fmt.Sprintf("%d", r.PublicationYear),
			DOI:       r.DOI,
			Authors:   authors,
			Source:    openAlexSource,
		})
	}

	logger.Log.Info("related papers found", zap.Int("count", len(papers)))
	return papers, nil
}
