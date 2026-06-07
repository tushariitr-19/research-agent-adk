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

//go:embed instructions/author_profile.md
var authorProfileInstruction string

func NewAuthorProfileAgent(llm model.LLM) (agent.Agent, error) {
	logger.Log.Debug("building author_profile_agent")
	client := tools.NewOpenAlexClient()

	authorTool, err := tools.NewSearchAuthorPapersTool(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create author tool: %w", err)
	}

	return llmagent.New(llmagent.Config{
		Name:        "AuthorProfileAgent",
		Model:       llm,
		Description: "Profiles the authors of a research paper using OpenAlex",
		Instruction: authorProfileInstruction,
		Tools:       []adktool.Tool{authorTool},
		OutputKey:   "author_profiles",
	})
}
