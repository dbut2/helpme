package main

import (
	"context"
	"fmt"
	"runtime"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Tool         string `envconfig:"HELPME_TOOL"`
	Ollama       *OllamaConfig
	Claude       *ClaudeConfig
	OpenAI       *OpenAIConfig
	SystemPrompt string `env:"HELPME_SYSTEM_PROMPT"`
	Debug        bool   `env:"HELPME_DEBUG, default=false"`
}

var defaultSystemPrompt = fmt.Sprintf("You are a bash based command line application running on %s that will return a command for the users need. You must only return the command to be ran. You must not supply an explanation.", runtime.GOOS)

func loadConfig() (Config, error) {
	var c Config
	err := envconfig.Process(context.Background(), &c)
	if err != nil {
		return Config{}, err
	}
	if c.SystemPrompt == "" {
		c.SystemPrompt = defaultSystemPrompt
	}
	return c, nil
}
