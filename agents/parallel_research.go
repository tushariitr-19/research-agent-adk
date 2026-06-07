package agents

import (
	"fmt"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/workflowagents/parallelagent"
	"google.golang.org/adk/model"

	"github.com/tushariitr-19/research-agent-adk/logger"
)

// NewParallelResearchAgent creates a parallel agent that simultaneously
// fetches related papers and profiles authors
func NewParallelResearchAgent(llm model.LLM) (agent.Agent, error) {
	logger.Log.Debug("building parallel_research_agent")
	relatedPapersAgent, err := NewRelatedPapersAgent(llm)
	if err != nil {
		return nil, fmt.Errorf("failed to create related papers agent: %w", err)
	}

	authorProfileAgent, err := NewAuthorProfileAgent(llm)
	if err != nil {
		return nil, fmt.Errorf("failed to create author profile agent: %w", err)
	}

	return parallelagent.New(parallelagent.Config{
		AgentConfig: agent.Config{
			Name:        "ParallelResearchAgent",
			Description: "Simultaneously fetches related papers and profiles authors",
			SubAgents:   []agent.Agent{relatedPapersAgent, authorProfileAgent},
		},
	})
}
