package main

import (
	"errors"
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/interactive_input/config"
	usecases "github.com/landmaster135/devbox/internal/interactive_input/usecases"
)

const (
	exitCodeOK    = 0
	exitCodeRetry = 1
	exitCodeFatal = 2
)

func main() {
	cfg, err := config.ParseFlags(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		config.PrintUsage()
		os.Exit(exitCodeFatal)
	}

	if cfg.Help {
		config.PrintUsage()
		return
	}

	service := usecases.NewService(os.Stdin, os.Stderr)
	result, err := service.Run(usecases.Config{
		Prompt:          cfg.Prompt,
		InputType:       cfg.InputType,
		Key:             cfg.Key,
		DefaultValue:    cfg.DefaultValue,
		DefaultProvided: cfg.DefaultProvided,
		ChoiceOptions:   cfg.ChoiceOptions,
		MaxAttempts:     cfg.MaxAttempts,
	})
	if err != nil {
		handleRunError(err)
	}

	fmt.Fprint(os.Stdout, result)
}

func handleRunError(err error) {
	switch {
	case errors.Is(err, usecases.ErrExceededAttempts):
		fmt.Fprintln(os.Stderr, "入力が規定回数内に完了しませんでした。")
		os.Exit(exitCodeRetry)
	case errors.Is(err, usecases.ErrUserCancelled):
		fmt.Fprintln(os.Stderr, "入力がキャンセルされました。")
		os.Exit(exitCodeRetry)
	default:
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(exitCodeFatal)
	}
}
