package main

import (
	"context"
	"errors"
	"io"

	"github.com/sashabaranov/go-openai"
)

type OpenAIConfig struct {
	ApiKey string `env:"OPENAI_TOKEN"`
	Model  string `env:"OPENAI_MODEL"`
}

type OpenAI struct {
	client *openai.Client
}

func NewOpenAI(config *OpenAIConfig) *OpenAI {
	return &OpenAI{
		client: openai.NewClient(config.ApiKey),
	}
}

func (o *OpenAI) Generate(ctx context.Context, system, prompt string, ch chan string, errCh chan error) error {
	stream, err := o.client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4o,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: system,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Stream: true,
	})
	if err != nil {
		return err
	}

	for {
		resp, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		if len(resp.Choices) == 0 {
			return errors.New("no response")
		}

		ch <- resp.Choices[0].Delta.Content
	}
}
