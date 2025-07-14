package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/melt-inc/autoupgrade"
)

func main() {
	ctx := context.Background()

	ur := autoupgrade.UpgradeBackground(ctx, "dbut.dev/helpme")
	defer func() {
		<-ur
	}()

	args := strings.TrimSpace(strings.Join(os.Args[1:], " "))
	if args == "" {
		fmt.Printf("Error: no argument provided\n")
		return
	}

	app, err := NewApp(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return
	}

	err = app.Run(ctx, args)
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return
	}
}
