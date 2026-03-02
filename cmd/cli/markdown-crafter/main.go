package main

import (
	"fmt"
	"os"

	"github.com/landmaster135/devbox/internal/markdown_crafter/config"
	"github.com/landmaster135/devbox/internal/markdown_crafter/usecases"
)

func main() {
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		config.PrintUsage()
		os.Exit(1)
	}

	if cfg.Help {
		config.PrintUsage()
		return
	}

	service := usecases.NewService(nil)
	result, err := service.ExecuteByConfig(cfg)

	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		config.PrintUsage()
		os.Exit(1)
	}

	fmt.Print(result)
}
