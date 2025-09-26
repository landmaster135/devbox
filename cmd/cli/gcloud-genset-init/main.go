package main

import (
	"fmt"
	"os"

	cfg "github.com/landmaster135/devbox/internal/gcloud_genset_init/config"
	"github.com/landmaster135/devbox/internal/gcloud_genset_init/usecases"
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
	case cfg.OperationAuthLogin:
		command, err := service.BuildAuthLoginCommand(usecases.AuthLoginParams{
			ProjectID:      config.ProjectID,
			AdditionalArgs: config.AdditionalArgs,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			os.Exit(1)
		}
		service.PrintHighlightedCommand(command)
	case cfg.OperationSetProjectConfig:
		command, err := service.BuildSetProjectConfigCommand(usecases.SetProjectConfigParams{
			ProjectID:      config.ProjectID,
			AdditionalArgs: config.AdditionalArgs,
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
