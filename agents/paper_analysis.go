package agents

import (
	"fmt"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/workflowagents/sequentialagent"
	"google.golang.org/adk/model"

	"github.com/tushariitr-19/research-agent-adk/logger"
	"github.com/tushariitr-19/research-agent-adk/tools"
)

// NewPaperAnalysisAgent creates a sequential agent that fetches then analyzes a paper
func NewPaperAnalysisAgent(llm model.LLM) (agent.Agent, error) {
	logger.Log.Debug("building paper_analysis_agent")
	arxivClient := tools.NewArxivClient()

	fetchAgent, err := NewFetchAgent(llm, arxivClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create fetch agent: %w", err)
	}

	analyzeAgent, err := NewAnalyzeAgent(llm)
	if err != nil {
		return nil, fmt.Errorf("failed to create analyze agent: %w", err)
	}

	return sequentialagent.New(sequentialagent.Config{
		AgentConfig: agent.Config{
			Name:        "PaperAnalysisAgent",
			Description: "Fetches a paper by ArXiv ID then provides a comprehensive analysis",
			SubAgents:   []agent.Agent{fetchAgent, analyzeAgent},
		},
	})
}
