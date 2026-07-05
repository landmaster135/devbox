package main

import (
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/disk_health/config"
	usecases "github.com/landmaster135/devbox/internal/disk_health/usecases"
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
	case config.OperationAssessSmart:
		handleAssessSmart(cfg)
	default:
		fmt.Fprintf(os.Stderr, "エラー: 未対応のoperationです: %s\n", cfg.Operation)
		config.PrintUsage()
		os.Exit(1)
	}
}

func handleAssessSmart(cfg *config.Config) {
	service := usecases.NewService(usecases.ServiceOptions{})
	result, err := service.AssessSmart(cfg.SrcFile, cfg.JSON, cfg.Verbose)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(result)
}
