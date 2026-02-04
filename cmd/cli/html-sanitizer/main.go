package main

import (
	"fmt"
	"os"

	htmlconfig "github.com/landmaster135/devbox/internal/html_sanitizer/config"
	htmlusecases "github.com/landmaster135/devbox/internal/html_sanitizer/usecases"
)

func main() {
	cfg, fs, err := htmlconfig.ParseFlags(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "フラグの解析に失敗しました: %v\n", err)
		if fs != nil {
			htmlconfig.PrintUsage(fs)
		}
		os.Exit(1)
	}

	if cfg.ShowHelp {
		htmlconfig.PrintUsage(fs)
		return
	}

	svc := htmlusecases.NewSanitizerService()
	if err := run(svc, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
}

func run(svc *htmlusecases.SanitizerService, cfg *htmlconfig.Config) error {
	_, err := svc.SanitizeFile(cfg.InputPath, cfg.OutputPath, cfg.OmitsFullBody)
	return err
}
