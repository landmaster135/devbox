package main

import (
	"fmt"
	"os"

	"github.com/landmaster135/devbox/internal/notion_to_memos_markdown/config"
	"github.com/landmaster135/devbox/internal/notion_to_memos_markdown/usecases"
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
	case config.OperationDistributeFiles:
		result, err = service.DistributeFiles(cfg.PageType, cfg.SrcJSONPath, cfg.SrcBodyDir, cfg.OutDir)
	default:
		err = fmt.Errorf("未対応のoperationです: %s", cfg.Operation)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(result)
}
