package main

import (
	"context"
	"fmt"
)

type App struct {
	config    Config
	generator Generator
}

type Generator interface {
	Generate(ctx context.Context, system, prompt string, ch chan string, errCh chan error) error
}

func NewApp(ctx context.Context) (*App, error) {
	config, err := loadConfig()
	if err != nil {
		panic(err.Error())
	}

	var generator Generator
	switch config.Tool {
	case "OLLAMA":
		generator, err = NewOllamaClient(ctx, config.Ollama)
		if err != nil {
			return nil, err
		}
	case "CHATGPT":
		generator = NewOpenAI(config.OpenAI)
	default:
		fallthrough
	case "CLAUDE":
		generator, err = NewClaude(config.Claude)
		if err != nil {
			return nil, err
		}
	}

	return &App{
		config:    config,
		generator: generator,
	}, nil
}

func (a *App) Run(ctx context.Context, prompt string) error {
	ch := make(chan string)
	errCh := make(chan error)

	defer close(ch)
	defer close(errCh)

	go func() {
		errCh <- a.generator.Generate(ctx, a.config.SystemPrompt, prompt, ch, errCh)
	}()

	for {
		select {
		case response := <-ch:
			fmt.Print(response)
		case err := <-errCh:
			if err == nil {
				fmt.Println()
			}
			return err
		}
	}
}
