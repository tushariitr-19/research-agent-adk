package main

import (
	"log"
	"os"

	"github.com/tushariitr-19/research-agent-adk/agents"
	"github.com/tushariitr-19/research-agent-adk/config"
	"github.com/tushariitr-19/research-agent-adk/logger"
)

func main() {
	debug := os.Getenv("DEBUG") == "true"
	if err := logger.Init(debug); err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := agents.Run(cfg); err != nil {
		log.Fatalf("agent failed: %v", err)
	}
}
