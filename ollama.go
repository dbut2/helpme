package main

import (
	"context"
	"fmt"

	"github.com/ollama/ollama/api"
)

type OllamaConfig struct {
	Model string `env:"HELPME_MODEL"`
}

const ollamaDefaultModel = "codegemma:instruct"

type Ollama struct {
	config *OllamaConfig
	client *api.Client
}

func NewOllamaClient(ctx context.Context, config *OllamaConfig, opts ...OllamaOpt) (*Ollama, error) {
	for _, opt := range opts {
		opt(config)
	}

	if config.Model == "" {
		config.Model = ollamaDefaultModel
	}

	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, err
	}

	err = client.Heartbeat(ctx)
	if err != nil {
		return nil, err
	}

	list, err := client.List(ctx)
	if err != nil {
		panic(err.Error())
	}

	contains := false
	for _, m := range list.Models {
		if m.Model == config.Model {
			contains = true
			break
		}
	}

	if !contains {
		fmt.Printf("Downloading %s...\n", config.Model)
		err = client.Pull(context.Background(), &api.PullRequest{
			Model:    config.Model,
			Insecure: true,
		}, func(response api.ProgressResponse) error {
			return nil
		})
		if err != nil {
			return nil, err
		}
		fmt.Printf("Done")
	}

	return &Ollama{
		config: config,
		client: client,
	}, nil
}

type OllamaOpt func(*OllamaConfig)

func WithModel(model string) OllamaOpt {
	return func(c *OllamaConfig) {
		c.Model = model
	}
}

func (o *Ollama) Generate(ctx context.Context, system, prompt string, ch chan string, errCh chan error) error {
	go func() {
		errCh <- o.client.Generate(ctx, &api.GenerateRequest{
			Model:  o.config.Model,
			Prompt: prompt,
			System: system,
		}, func(response api.GenerateResponse) error {
			ch <- response.Response
			return nil
		})
	}()
	return nil
}
