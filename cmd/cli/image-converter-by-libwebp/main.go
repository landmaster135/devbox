package main

import (
	"context"
	"fmt"
	"io"
	"os"

	config "github.com/landmaster135/devbox/internal/image_converter_by_libwebp/config"
	usecases "github.com/landmaster135/devbox/internal/image_converter_by_libwebp/usecases"
)

type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeError
)

func run(args []string, stdout io.Writer, stderr io.Writer) exitCode {
	cfg, err := config.ParseFlagsWithArgs(args)
	if err != nil {
		fmt.Fprintf(stderr, "エラー: %v\n", err)
		config.PrintUsage()
		return exitCodeError
	}
	if cfg.Help {
		config.PrintUsage()
		return exitCodeOK
	}

	service := usecases.NewService()
	result, err := service.Convert(context.Background(), cfg)
	if err != nil {
		fmt.Fprintf(stderr, "エラー: %v\n", err)
		if result == nil {
			return exitCodeError
		}
	}

	if result != nil {
		fmt.Fprintln(stdout, "画像変換が完了しました")
		fmt.Fprintf(stdout, "  成功: %d ファイル\n", result.SuccessCount)
		fmt.Fprintf(stdout, "  失敗: %d ファイル\n", result.ErrorCount)
		fmt.Fprintf(stdout, "  出力先: %s\n", result.OutputDir)
	}

	if err != nil {
		return exitCodeError
	}
	return exitCodeOK
}

func main() {
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}
