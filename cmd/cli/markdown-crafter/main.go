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

	var result string
	switch cfg.Operation {
	case config.OperationSplitHeadings:
		result, err = service.SplitHeadings(cfg.FilePath, cfg.HeadingLevel, cfg.OutputDir)
	case config.OperationAddFrontMatter:
		result, err = service.AddFrontMatter(cfg.FilePath, cfg.KVPairs)
	case config.OperationAddTags:
		result, err = service.AddTags(cfg.FilePath, cfg.Tags)
	default:
		err = fmt.Errorf("未サポートのoperationです: %s", cfg.Operation)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		config.PrintUsage()
		os.Exit(1)
	}

	fmt.Print(result)
}
