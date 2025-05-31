package main

import (
	"fmt"
	"os"

	"github.com/landmaster135/devbox/internal/independencies/claude_code_usage/cmd"
)

func main() {
	app := cmd.NewApp()
	
	if err := app.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
