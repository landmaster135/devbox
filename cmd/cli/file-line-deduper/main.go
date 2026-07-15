package main

import (
	"fmt"
	"io"
	"os"

	config "github.com/landmaster135/devbox/internal/file_line_deduper/config"
	services "github.com/landmaster135/devbox/internal/file_line_deduper/usecases/services"
)

// exitCode はプログラムの終了コードを表します
type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeError
)

// run はファイル処理の主要なロジックを実行します
func run(args []string, stdout, stderr io.Writer) exitCode {
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

	service := services.NewCLIService()

	result, err := service.HandleRemoveMatchingLines(cfg.FilePath, cfg.StartPos, cfg.EndPos)
	if err != nil {
		fmt.Fprintf(stderr, "エラー: %v\n", err)
		return exitCodeError
	}

	fmt.Fprint(stdout, result)
	return exitCodeOK
}

func main() {
	// run関数を呼び出し、結果に応じて終了コードを設定
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}
