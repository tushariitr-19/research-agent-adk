package config

import (
	"fmt"
	"os"
)

const (
	GeminiModel = "gemini-flash-latest"
)

type Config struct {
	GoogleAPIKey string
	Debug        bool
}

func Load() (*Config, error) {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GOOGLE_API_KEY environment variable is required")
	}

	return &Config{
		GoogleAPIKey: apiKey,
		Debug:        os.Getenv("DEBUG") == "true",
	}, nil
}
