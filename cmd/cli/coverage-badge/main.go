package main

import (
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/coverage_badge/config"
	usecases "github.com/landmaster135/devbox/internal/coverage_badge/usecases"
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

	service := usecases.NewCoverageBadgeService()

	switch cfg.Operation {
	case config.OperationCreateBadge:
		handleCreateBadgeOperation(service, cfg)
	case config.OperationPatchBadge:
		handlePatchBadgeOperation(service, cfg)
	default:
		fmt.Fprintf(os.Stderr, "エラー: 未対応のoperationです: %s\n", cfg.Operation)
		config.PrintUsage()
		os.Exit(1)
	}
}

func handleCreateBadgeOperation(service *usecases.CoverageBadgeService, cfg *config.Config) {
	badge, err := service.CreateBadge(usecases.CreateBadgeInput{
		BadgeTitle:      cfg.BadgeTitle,
		CoverageFile:    cfg.CoverageFile,
		GreenThreshold:  cfg.GreenThreshold,
		YellowThreshold: cfg.YellowThreshold,
		ForceColor:      cfg.ForceColor,
		BadgeLink:       cfg.BadgeLink,
		BadgeValue:      cfg.BadgeValue,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(badge)
}

func handlePatchBadgeOperation(service *usecases.CoverageBadgeService, cfg *config.Config) {
	result, err := service.PatchBadge(usecases.PatchBadgeInput{
		CreateBadgeInput: usecases.CreateBadgeInput{
			BadgeTitle:      cfg.BadgeTitle,
			CoverageFile:    cfg.CoverageFile,
			GreenThreshold:  cfg.GreenThreshold,
			YellowThreshold: cfg.YellowThreshold,
			ForceColor:      cfg.ForceColor,
			BadgeLink:       cfg.BadgeLink,
			BadgeValue:      cfg.BadgeValue,
		},
		TargetFile: cfg.TargetFile,
		DryRun:     cfg.DryRun,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	if cfg.DryRun {
		fmt.Print(result.PatchedContent)
		return
	}

	if result.ContentModified {
		fmt.Printf("カバレッジバッジを更新しました: %s\n", result.TargetFile)
		return
	}

	fmt.Printf("変更はありませんでした: %s\n", result.TargetFile)
}
