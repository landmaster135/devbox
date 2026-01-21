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
	case config.OperationPortsIntoCompose:
		handlePortsIntoCompose(cfg)
	case config.OperationVolumesIntoCompose:
		handleVolumesIntoCompose(cfg)
	case config.OperationUserIntoCompose:
		handleUserIntoCompose(cfg)
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

func handlePortsIntoCompose(cfg *config.Config) {
	service := usecases.NewEnvSyncService()
	if err := service.SyncPortsIntoCompose(cfg.EnvYAMLPath, cfg.ComposePath, cfg.PortKey, cfg.Service); err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s のポートを %s に更新しました\n", cfg.Service, cfg.PortKey)
}

func handleVolumesIntoCompose(cfg *config.Config) {
	service := usecases.NewEnvSyncService()
	if err := service.SyncVolumesIntoCompose(cfg.EnvYAMLPath, cfg.ComposePath, cfg.VolumeKey, cfg.Service); err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s の volumes を %s の値で更新しました\n", cfg.Service, cfg.VolumeKey)
}

func handleUserIntoCompose(cfg *config.Config) {
	service := usecases.NewEnvSyncService()
	if err := service.SyncUserIntoCompose(cfg.EnvYAMLPath, cfg.ComposePath, cfg.UserKey, cfg.Service); err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s の user を %s の値で更新しました\n", cfg.Service, cfg.UserKey)
}
