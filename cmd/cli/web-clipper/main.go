package main

import (
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/web_clipper/config"
	usecases "github.com/landmaster135/devbox/internal/web_clipper/usecases"
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
	case config.OperationPatchMarkdown:
		result, err = service.PatchMarkdown(
			cfg.TargetTitle,
			cfg.TargetURL,
			cfg.SrcMarkdownContent,
			cfg.SrcMarkdownFile,
			cfg.OutFilePath,
			cfg.TopHeadingLevel,
		)
	case config.OperationRenameAttachments:
		result, err = service.RenameAttachments(usecases.RenameAttachmentsOptions{
			SrcDir:     cfg.SrcDir,
			Slug:       cfg.Slug,
			Start:      cfg.Start,
			Digits:     cfg.Digits,
			SortByTime: cfg.SortByTime,
			SortByName: cfg.SortByName,
			JSON:       cfg.JSON,
			Verbose:    cfg.Verbose,
		})
	default:
		err = fmt.Errorf("未対応のoperationです: %s", cfg.Operation)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(result)
}
