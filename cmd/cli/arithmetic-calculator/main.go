package main

import (
	"fmt"
	"io"
	"os"

	"github.com/landmaster135/devbox/internal/arithmetic_calculator/config"
	"github.com/landmaster135/devbox/internal/arithmetic_calculator/usecases"
)

type operationExecutor interface {
	ExecuteByConfig(cfg *config.Config) (string, error)
}

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

	if err := execute(cfg, usecases.NewService(), os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
}

func execute(cfg *config.Config, executor operationExecutor, stdout io.Writer) error {
	result, err := executor.ExecuteByConfig(cfg)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(stdout, result)
	return err
}
