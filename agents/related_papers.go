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

//go:embed instructions/related_papers.md
var relatedPapersInstruction string

func NewRelatedPapersAgent(llm model.LLM) (agent.Agent, error) {
	logger.Log.Debug("building related_papers_agent")
	client := tools.NewOpenAlexClient()

	searchTool, err := tools.NewSearchRelatedPapersTool(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create search tool: %w", err)
	}

	return llmagent.New(llmagent.Config{
		Name:        "RelatedPapersAgent",
		Model:       llm,
		Description: "Finds papers related to a given research paper using OpenAlex",
		Instruction: relatedPapersInstruction,
		Tools:       []adktool.Tool{searchTool},
		OutputKey:   "related_papers",
	})
}
