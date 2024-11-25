package main

import (
	"context"
	"fmt"
	"os"
	"strings"
)

func main() {
	ctx := context.Background()

	args := strings.TrimSpace(strings.Join(os.Args[1:], " "))
	if args == "" {
		fmt.Println("no argument provided")
		os.Exit(1)
	}

	app, err := NewApp(ctx)
	if err != nil {
		panic(err.Error())
	}

	err = app.Run(ctx, args)
	if err != nil {
		panic(err.Error())
	}
}
