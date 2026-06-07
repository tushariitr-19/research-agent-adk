package agents

import (
	_ "embed"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"

	"github.com/tushariitr-19/research-agent-adk/logger"
)

//go:embed instructions/summarize.md
var summarizeInstruction string

// NewSummarizeAgent creates an agent that summarizes search results
func NewSummarizeAgent(llm model.LLM) (agent.Agent, error) {
	logger.Log.Debug("building summarize_agent")
	return llmagent.New(llmagent.Config{
		Name:        "SummarizeAgent",
		Model:       llm,
		Description: "Summarizes a list of research papers into a concise readable report",
		Instruction: summarizeInstruction,
		OutputKey:   "search_summary",
	})
}
