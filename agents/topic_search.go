package agents

import (
	"fmt"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/workflowagents/sequentialagent"
	"google.golang.org/adk/model"

	"github.com/tushariitr-19/research-agent-adk/logger"
)

// NewTopicSearchAgent creates a sequential agent that searches then summarizes results
func NewTopicSearchAgent(llm model.LLM) (agent.Agent, error) {
	logger.Log.Debug("building topic_search_agent")
	searchAgent, err := NewSearchAgent(llm)
	if err != nil {
		return nil, fmt.Errorf("failed to create search agent: %w", err)
	}

	summarizeAgent, err := NewSummarizeAgent(llm)
	if err != nil {
		return nil, fmt.Errorf("failed to create summarize agent: %w", err)
	}

	return sequentialagent.New(sequentialagent.Config{
		AgentConfig: agent.Config{
			Name:        "TopicSearchAgent",
			Description: "Searches ArXiv by topic then summarizes the results into a readable report",
			SubAgents:   []agent.Agent{searchAgent, summarizeAgent},
		},
	})
}
