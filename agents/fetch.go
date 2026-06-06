package agents

import (
	_ "embed"
	"fmt"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	adktool "google.golang.org/adk/tool"

	"github.com/tushariitr-19/research-agent-adk/logger"
	"github.com/tushariitr-19/research-agent-adk/tools"
)

//go:embed instructions/fetch.md
var fetchInstruction string

// NewFetchAgent creates a sub-agent that fetches a research paper by its ArXiv ID
func NewFetchAgent(llm model.LLM, arxivClient *tools.ArxivClient) (agent.Agent, error) {
	logger.Log.Debug("building fetch_agent")
	fetchTool, err := tools.NewFetchPaperTool(arxivClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create fetch tool: %w", err)
	}

	fetchAgent, err := llmagent.New(llmagent.Config{
		Name:        "fetch_agent",
		Model:       llm,
		Description: "Fetches a research paper by its ArXiv ID",
		Instruction: fetchInstruction,
		Tools:       []adktool.Tool{fetchTool},
		OutputKey:   "paper_details",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create fetch agent: %w", err)
	}

	return fetchAgent, nil
}
