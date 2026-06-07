package agents

import (
	_ "embed"
	"fmt"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/agent/workflowagents/sequentialagent"
	"google.golang.org/adk/model"

	"github.com/tushariitr-19/research-agent-adk/logger"
	"github.com/tushariitr-19/research-agent-adk/tools"
)

//go:embed instructions/synthesize.md
var synthesizeInstruction string

// NewDeepResearchAgent creates a sequential agent that:
// 1. Fetches the paper
// 2. In parallel: finds related papers + profiles authors
// 3. Synthesizes everything into a comprehensive report
func NewDeepResearchAgent(llm model.LLM) (agent.Agent, error) {
	logger.Log.Debug("building deep_research_agent")
	arxivClient := tools.NewArxivClient()

	fetchAgent, err := NewFetchAgent(llm, arxivClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create fetch agent: %w", err)
	}

	parallelResearchAgent, err := NewParallelResearchAgent(llm)
	if err != nil {
		return nil, fmt.Errorf("failed to create parallel research agent: %w", err)
	}

	synthesizeAgent, err := llmagent.New(llmagent.Config{
		Name:        "SynthesizeAgent",
		Model:       llm,
		Description: "Synthesizes paper details, related papers and author profiles into a comprehensive report",
		Instruction: synthesizeInstruction,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create synthesize agent: %w", err)
	}

	return sequentialagent.New(sequentialagent.Config{
		AgentConfig: agent.Config{
			Name:        "DeepResearchAgent",
			Description: "Performs deep research on a paper — fetches, runs parallel analysis, and synthesizes a comprehensive report",
			SubAgents:   []agent.Agent{fetchAgent, parallelResearchAgent, synthesizeAgent},
		},
	})
}
