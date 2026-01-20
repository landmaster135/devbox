package main

import (
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/taskfile/config"
	usecases "github.com/landmaster135/devbox/internal/taskfile/usecases"
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

	switch cfg.Operation {
	case config.OperationInspect:
		handleInspect(cfg)
	case config.OperationFill:
		handleFill(cfg)
	default:
		fmt.Fprintf(os.Stderr, "エラー: 未サポートのoperationです: %s\n", cfg.Operation)
		os.Exit(1)
	}
}

func handleInspect(cfg *config.Config) {
	service := usecases.NewService()
	result, err := service.Inspect(cfg.TaskType, cfg.TaskfilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	if result.HasMissingFields() {
		fmt.Fprintf(os.Stderr, "不足しているフィールドが %d 個見つかりました:\n", len(result.MissingFields))
		for _, field := range result.MissingFields {
			fmt.Fprintf(os.Stderr, "  - %s\n", field)
		}
		os.Exit(1)
	}

	fmt.Println("Taskfileには参照Taskfileのすべてのフィールドが含まれています。")
}

func handleFill(cfg *config.Config) {
	service := usecases.NewService()
	updated, err := service.Fill(cfg.TaskType, cfg.TaskfilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	if updated {
		fmt.Println("Taskfileの空欄フィールドをテンプレートの値で補完しました。")
		return
	}

	fmt.Println("補完対象の空欄フィールドは見つかりませんでした。")
}
