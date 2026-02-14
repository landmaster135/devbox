package main

import (
	"fmt"
	"io"
	"os"

	"github.com/landmaster135/devbox/internal/data_converter/config"
	"github.com/landmaster135/devbox/internal/data_converter/usecases"
)

type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeError
)

func run(args []string, stdout, stderr io.Writer) exitCode {
	cfg, err := config.ParseFlags(args)
	if err != nil {
		fmt.Fprintf(stderr, "エラー: %v\n", err)
		config.PrintUsage(stderr)
		return exitCodeError
	}

	if cfg.Help {
		config.PrintUsage(stdout)
		return exitCodeOK
	}

	service := usecases.NewService()
	message, err := service.ConvertFile(
		cfg.InputFilePath,
		cfg.OutputFilePath,
		cfg.InputFormat,
		cfg.OutputFormat,
	)
	if err != nil {
		fmt.Fprintf(stderr, "エラー: %v\n", err)
		return exitCodeError
	}

	fmt.Fprintln(stdout, message)
	return exitCodeOK
}

func main() {
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}
