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

	service := usecases.NewService(usecases.ServiceOptions{})
	result, err := service.ExecuteByConfig(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(result)
}
