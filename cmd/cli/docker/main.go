package main

import (
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/docker/config"
	"github.com/landmaster135/devbox/internal/docker/usecases"
)

func main() {
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		config.PrintUsage()
		os.Exit(1)
	}

	switch cfg.Operation {
	case config.OperationEnvIntoCompose:
		handleEnvIntoCompose(cfg)
	default:
		fmt.Fprintf(os.Stderr, "エラー: 未対応のoperationです: %s\n", cfg.Operation)
		config.PrintUsage()
		os.Exit(1)
	}
}

func handleEnvIntoCompose(cfg *config.Config) {
	service := usecases.NewEnvSyncService()
	count, err := service.SyncEnvIntoCompose(cfg.EnvYAMLPath, cfg.ComposePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%d件の環境変数を %s に反映しました\n", count, cfg.ComposePath)
}
