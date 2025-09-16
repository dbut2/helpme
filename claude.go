package main

import (
	"context"
	"errors"

	"github.com/liushuangls/go-anthropic/v2"
)

type ClaudeConfig struct {
	ApiKey string `env:"ANTHROPIC_API_TOKEN"`
}

type Claude struct {
	config *ClaudeConfig
	client *anthropic.Client
}

func NewClaude(config *ClaudeConfig) (*Claude, error) {
	if config.ApiKey == "" {
		return nil, errors.New("please set ANTHROPIC_API_TOKEN")
	}
	return &Claude{
		client: anthropic.NewClient(config.ApiKey),
	}, nil
}

func (c *Claude) Generate(ctx context.Context, system, prompt string, ch chan string, errCh chan error) error {
	_, err := c.client.CreateMessagesStream(ctx, anthropic.MessagesStreamRequest{
		MessagesRequest: anthropic.MessagesRequest{
			Model: "claude-sonnet-4-20250514",
			Messages: []anthropic.Message{
				{
					Role: anthropic.RoleUser,
					Content: []anthropic.MessageContent{
						{
							Type: "text",
							Text: &prompt,
						},
					},
				},
			},
			MaxTokens: 8192,
			System:    system,
		},
		OnError: func(response anthropic.ErrorResponse) {
			errCh <- response.Error
		},
		OnContentBlockDelta: func(data anthropic.MessagesEventContentBlockDeltaData) {
			ch <- *data.Delta.Text
		},
	})
	return err
}
