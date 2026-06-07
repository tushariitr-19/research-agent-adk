package agents

import (
	"context"
	_ "embed"
	"fmt"
	"os"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/memory"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/session"
	adktool "google.golang.org/adk/tool"
	"google.golang.org/adk/tool/agenttool"
	"google.golang.org/genai"

	"github.com/tushariitr-19/research-agent-adk/config"
	"github.com/tushariitr-19/research-agent-adk/logger"
	"go.uber.org/zap"
)

//go:embed instructions/root.md
var rootInstruction string

func New(cfg *config.Config) (agent.Agent, error) {
	ctx := context.Background()

	logger.Log.Debug("creating gemini model")
	llm, err := gemini.NewModel(ctx, config.GeminiModel, &genai.ClientConfig{
		APIKey: cfg.GoogleAPIKey,
	})
	if err != nil {
		logger.Log.Error("failed to create gemini model", zap.Error(err))
		return nil, fmt.Errorf("failed to create model: %w", err)
	}
	logger.Log.Debug("gemini model created")

	paperAnalysisAgent, err := NewPaperAnalysisAgent(llm)
	if err != nil {
		return nil, fmt.Errorf("failed to create paper analysis agent: %w", err)
	}

	topicSearchAgent, err := NewTopicSearchAgent(llm)
	if err != nil {
		return nil, fmt.Errorf("failed to create topic search agent: %w", err)
	}

	deepResearchAgent, err := NewDeepResearchAgent(llm)
	if err != nil {
		return nil, fmt.Errorf("failed to create deep research agent: %w", err)
	}

	paperAnalysisTool := agenttool.New(paperAnalysisAgent, &agenttool.Config{})

	topicSearchTool := agenttool.New(topicSearchAgent, &agenttool.Config{})

	deepResearchTool := agenttool.New(deepResearchAgent, &agenttool.Config{})

	rootAgent, err := llmagent.New(llmagent.Config{
		Name:        "research_agent",
		Model:       llm,
		Description: "A research paper assistant powered by ArXiv",
		Instruction: rootInstruction,
		Tools:       []adktool.Tool{paperAnalysisTool, topicSearchTool, deepResearchTool},
	})
	if err != nil {
		logger.Log.Error("failed to create root agent", zap.Error(err))
		return nil, fmt.Errorf("failed to create root agent: %w", err)
	}
	logger.Log.Info("root agent ready")
	return rootAgent, nil
}

func Run(cfg *config.Config) error {
	logger.Log.Info("starting research-agent-adk")

	rootAgent, err := New(cfg)
	if err != nil {
		logger.Log.Error("failed to create root agent", zap.Error(err))
		return fmt.Errorf("failed to create root agent: %w", err)
	}

	launcherCfg := &launcher.Config{
		AgentLoader:    agent.NewSingleLoader(rootAgent),
		SessionService: session.InMemoryService(),
		MemoryService:  memory.InMemoryService(),
	}

	l := full.NewLauncher()
	logger.Log.Info("launching agent server")
	if err = l.Execute(context.Background(), launcherCfg, os.Args[1:]); err != nil {
		logger.Log.Error("launcher failed", zap.Error(err))
		return fmt.Errorf("run failed: %w", err)
	}

	return nil
}
