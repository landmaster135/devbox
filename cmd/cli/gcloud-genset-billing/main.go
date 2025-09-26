package main

import (
	"fmt"
	"os"

	cfg "github.com/landmaster135/devbox/internal/gcloud_genset_billing/config"
	"github.com/landmaster135/devbox/internal/gcloud_genset_billing/usecases"
)

func main() {
	conf, err := cfg.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		cfg.PrintUsage()
		os.Exit(1)
	}

	if conf.Help {
		cfg.PrintUsage()
		return
	}

	service := usecases.NewService()

	command, err := service.BuildCommand(conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	service.PrintHighlightedCommand(command)
}
