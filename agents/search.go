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

//go:embed instructions/search.md
var searchInstruction string

func NewSearchAgent(llm model.LLM) (agent.Agent, error) {
	logger.Log.Debug("building search_agent")
	client := tools.NewOpenAlexClient()

	searchTool, err := tools.NewSearchRelatedPapersTool(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create search tool: %w", err)
	}

	return llmagent.New(llmagent.Config{
		Name:        "SearchAgent",
		Model:       llm,
		Description: "Searches for research papers by topic or keyword using OpenAlex",
		Instruction: searchInstruction,
		Tools:       []adktool.Tool{searchTool},
		OutputKey:   "search_results",
	})
}
