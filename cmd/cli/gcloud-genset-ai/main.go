package main

import (
	"fmt"
	"os"

	cfg "github.com/landmaster135/devbox/internal/gcloud-genset-ai/config"
	usecases "github.com/landmaster135/devbox/internal/gcloud-genset-ai/usecases"
)

func main() {
	config, err := cfg.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		cfg.PrintUsage()
		os.Exit(1)
	}

	if config.Help {
		cfg.PrintUsage()
		return
	}

	service := usecases.NewService()

	switch config.Operation {
	case cfg.OperationUndeployProcessorVersion:
		command, err := service.BuildUndeployProcessorVersionCommand(usecases.UndeployProcessorVersionParams{
			Region:        config.Region,
			ProjectNumber: config.ProjectNumber,
			ProcessorID:   config.ProcessorID,
			VersionID:     config.VersionID,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			os.Exit(1)
		}
		service.PrintHighlightedCommand(command)
	default:
		fmt.Fprintf(os.Stderr, "エラー: 未対応の操作です: %s\n", config.Operation)
		cfg.PrintUsage()
		os.Exit(1)
	}
}
