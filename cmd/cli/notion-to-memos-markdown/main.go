package main

import (
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/config"
	progress "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/infrastructures/progress"
	usecases "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/usecases"
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

	service := usecases.NewServiceWithReporter(nil, progress.NewWriterReporter(os.Stderr))

	var result string
	switch cfg.Operation {
	case config.OperationDistributeFiles:
		result, err = service.DistributeFiles(string(cfg.PageType), cfg.SrcJSONFile, cfg.SrcBodyDir, cfg.OutDir)
	case config.OperationCraftMarkdown:
		result, err = service.CraftMarkdown(
			string(cfg.PageType),
			cfg.Category,
			cfg.SkipsNoSrcBody,
			cfg.ConNumberStart,
			cfg.ConNumberEnd,
			cfg.SrcJSONFile,
			cfg.SrcBodyDir,
			cfg.SrcResourceDir,
			cfg.OutDir,
			cfg.OutResourceDir,
		)
	case config.OperationCheckBodyLength:
		result, err = service.CheckBodyLength(cfg.SrcBodyDir, cfg.Threshold)
	case config.OperationGrepStr:
		result, err = service.GrepStr(cfg.SrcBodyDir, cfg.TargetStr)
	case config.OperationRenameBodiesByCategoryID:
		result, err = service.RenameBodiesByCategoryID(
			string(cfg.PageType),
			cfg.ConNumberStart,
			cfg.ConNumberEnd,
			cfg.SrcJSONFile,
			cfg.SrcResourceDir,
		)
	case config.OperationMigrateToMemos:
		result, err = service.MigrateToMemos(
			string(cfg.PageType),
			cfg.BaseURL,
			cfg.APIToken,
			cfg.SrcBodyDir,
			cfg.SrcResourceDir,
		)
	default:
		err = fmt.Errorf("未対応のoperationです: %s", cfg.Operation)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(result)
}
