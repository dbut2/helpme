package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

func main() {
	args := strings.TrimSpace(strings.Join(os.Args[1:], " "))
	if args == "" {
		fmt.Println("no argument provided")
		os.Exit(1)
	}

	openaiToken := os.Getenv("OPENAI_TOKEN")
	if openaiToken == "" {
		fmt.Println("env var not set")
		os.Exit(1)
	}

	client := openai.NewClient(openaiToken)

	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT4,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a bash based command line application running on macOS that will return a command for the users need. You must only return the command to be ran. You must not supply an explanation. If the user requests multiple commands you should only return a single back.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: args,
			},
		},
		Temperature: 0.1,
	})

	if err != nil {
		fmt.Printf("ai error: %s\n", err.Error())
		os.Exit(1)
	}

	if len(resp.Choices) == 0 {
		fmt.Println("no response")
		os.Exit(1)
	}

	fmt.Println(resp.Choices[0].Message.Content)
}
