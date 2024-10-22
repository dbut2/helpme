package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/ollama/ollama/api"
)

func main() {
	const model = "codegemma:instruct"

	args := strings.TrimSpace(strings.Join(os.Args[1:], " "))
	if args == "" {
		fmt.Println("no argument provided")
		os.Exit(1)
	}

	client, err := api.ClientFromEnvironment()
	if err != nil {
		panic(err.Error())
	}

	err = client.Heartbeat(context.Background())
	if err != nil {
		panic(err.Error())
	}

	list, err := client.List(context.Background())
	if err != nil {
		panic(err.Error())
	}

	contains := false
	for _, m := range list.Models {
		if m.Model == model {
			contains = true
			break
		}
	}

	if !contains {
		fmt.Printf("Downloading %s...\n", model)
		err = client.Pull(context.Background(), &api.PullRequest{
			Model:    model,
			Insecure: true,
		}, func(response api.ProgressResponse) error {
			return nil
		})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("Done")
	}

	prompt := "You are a bash based command line application running on %s that will return a command for the users need. You must only return the command to be ran. You must not supply an explanation. If the user requests multiple commands you should only return a single back. Don't wrap in a code block."

	err = client.Chat(context.Background(), &api.ChatRequest{
		Model: model,
		Messages: []api.Message{
			{
				Role:    "system",
				Content: fmt.Sprintf(prompt, runtime.GOOS),
			},
			{
				Role:    "user",
				Content: args,
			},
		},
		Stream: ptrTo(true),
	}, func(response api.ChatResponse) error {
		fmt.Print(response.Message.Content)
		return nil
	})
	if err != nil {
		fmt.Printf("ai error: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println()
}

func ptrTo[T any](v T) *T {
	return &v
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
