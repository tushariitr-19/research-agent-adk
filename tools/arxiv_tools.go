package tools

import (
	"context"
	"encoding/json"
	"fmt"

	adktool "google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"

	"github.com/tushariitr-19/research-agent-adk/logger"
	"github.com/tushariitr-19/research-agent-adk/models"
	"github.com/tushariitr-19/research-agent-adk/util"
	"go.uber.org/zap"
)

// PaperSource is the interface all paper sources must implement
type PaperSource interface {
	FetchByID(ctx context.Context, id string) (*models.Paper, error)
	Search(ctx context.Context, query string, limit int) (*models.SearchResult, error)
}

// ArxivClient implements PaperSource using the shared HTTP client
type ArxivClient struct {
	http *util.HTTPClient
}

// NewArxivClient creates a new ArXiv client using the singleton HTTP client
func NewArxivClient() *ArxivClient {
	return &ArxivClient{http: util.GetHTTPClient()}
}

// FetchByID fetches a single paper by its ArXiv ID
func (c *ArxivClient) FetchByID(ctx context.Context, id string) (*models.Paper, error) {
	logger.Log.Info("fetching paper by ID", zap.String("id", id))

	reqURL := util.BuildArxivURL(map[string]string{
		"id_list":     id,
		"max_results": fmt.Sprintf("%d", arxivFetchOneLimit),
	})

	data, err := c.http.Get(ctx, reqURL)
	if err != nil {
		return nil, fmt.Errorf("fetch by ID failed: %w", err)
	}

	feed, err := util.ParseArxivFeed(data)
	if err != nil {
		return nil, err
	}

	if len(feed.Entries) == 0 {
		return nil, fmt.Errorf("paper not found: %s", id)
	}

	paper := util.EntryToPaper(feed.Entries[0])
	logger.Log.Info("paper fetched", zap.String("title", paper.Title))
	return paper, nil
}

// Search searches for papers by keyword
func (c *ArxivClient) Search(ctx context.Context, query string, limit int) (*models.SearchResult, error) {
	logger.Log.Info("searching papers", zap.String("query", query), zap.Int("limit", limit))

	if limit <= 0 {
		limit = arxivDefaultLimit
	}
	if limit > arxivMaxLimit {
		limit = arxivMaxLimit
	}

	reqURL := util.BuildArxivURL(map[string]string{
		"search_query": query,
		"max_results":  fmt.Sprintf("%d", limit),
		"sortBy":       arxivSortBy,
		"sortOrder":    arxivSortOrder,
	})

	data, err := c.http.Get(ctx, reqURL)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	feed, err := util.ParseArxivFeed(data)
	if err != nil {
		return nil, err
	}

	papers := make([]models.Paper, 0, len(feed.Entries))
	for _, entry := range feed.Entries {
		papers = append(papers, *util.EntryToPaper(entry))
	}

	logger.Log.Info("search complete", zap.Int("results", len(papers)))
	return &models.SearchResult{
		Query:      query,
		TotalFound: feed.TotalResults,
		Papers:     papers,
		Source:     openAlexSource,
	}, nil
}

// ADK Tool definitions

type fetchPaperArgs struct {
	ID string `json:"id" jsonschema:"The ArXiv paper ID e.g. 1706.03762"`
}

type searchPapersArgs struct {
	Query string `json:"query" jsonschema:"The search query e.g. multi-agent systems"`
	Limit int    `json:"limit" jsonschema:"Maximum number of results to return, default 5"`
}

// NewFetchPaperTool creates an ADK tool that fetches a paper by ID
func NewFetchPaperTool(client *ArxivClient) (adktool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "fetch_paper",
			Description: "Fetch a research paper by its ArXiv ID",
		},
		func(ctx adktool.Context, args fetchPaperArgs) (map[string]any, error) {
			logger.Log.Info("tool called: fetch_paper", zap.String("id", args.ID))
			paper, err := client.FetchByID(context.Background(), args.ID)
			if err != nil {
				logger.Log.Error("fetch_paper failed", zap.String("id", args.ID), zap.Error(err))
				return nil, err
			}
			data, err := json.Marshal(paper)
			if err != nil {
				logger.Log.Error("fetch_paper: marshal failed", zap.Error(err))
				return nil, err
			}
			return map[string]any{"result": string(data)}, nil
		},
	)
}

// NewSearchPapersTool creates an ADK tool that searches papers by query
func NewSearchPapersTool(client *ArxivClient) (adktool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "search_papers",
			Description: "Search for research papers on ArXiv by keyword or topic",
		},
		func(ctx adktool.Context, args searchPapersArgs) (map[string]any, error) {
			logger.Log.Info("tool called: search_papers", zap.String("query", args.Query), zap.Int("limit", args.Limit))
			if args.Limit == 0 {
				args.Limit = 5
			}
			result, err := client.Search(context.Background(), args.Query, args.Limit)
			if err != nil {
				logger.Log.Error("search_papers failed", zap.String("query", args.Query), zap.Error(err))
				return nil, err
			}
			data, err := json.Marshal(result)
			if err != nil {
				logger.Log.Error("search_papers: marshal failed", zap.Error(err))
				return nil, err
			}
			return map[string]any{"result": string(data)}, nil
		},
	)
}
