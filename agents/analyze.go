package agents

import (
	_ "embed"
	"fmt"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"

	"github.com/tushariitr-19/research-agent-adk/logger"
)

//go:embed instructions/analyze.md
var analyzeInstruction string

// NewAnalyzeAgent creates a sub-agent that analyzes and summarizes a research paper
func NewAnalyzeAgent(llm model.LLM) (agent.Agent, error) {
	logger.Log.Debug("building analyze_agent")
	analyzeAgent, err := llmagent.New(llmagent.Config{
		Name:        "analyze_agent",
		Model:       llm,
		Description: "Analyzes and summarizes a research paper",
		Instruction: analyzeInstruction,
		OutputKey:   "paper_analysis",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create analyze agent: %w", err)
	}

	return analyzeAgent, nil
}
