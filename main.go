package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/ollama/ollama/api"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Model string `env:"HELPME_MODEL, default=codegemma:instruct"`
}

func main() {
	var config Config
	err := envconfig.Process(context.Background(), &config)
	if err != nil {
		panic(err.Error())
	}

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
			panic(err.Error())
		}
		fmt.Printf("Done")
	}

	prompt := "You are a bash based command line application running on %s that will return a command for the users need. You must only return the command to be ran. You must not supply an explanation."

	err = client.Generate(context.Background(), &api.GenerateRequest{
		Model:  config.Model,
		Prompt: args,
		System: fmt.Sprintf(prompt, runtime.GOOS),
		Stream: ptrTo(true),
	}, func(response api.GenerateResponse) error {
		fmt.Print(response.Response)
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
