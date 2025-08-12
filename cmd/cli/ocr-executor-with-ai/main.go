package main

import (
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/ocr_executor_with_ai/config"
	usecases "github.com/landmaster135/devbox/internal/ocr_executor_with_ai/usecases"
)

// handleOcrExecution はOCR実行を処理する
func handleOcrExecution(cfg *config.Config) {
	// OcrExecutorServiceを初期化
	service, err := usecases.NewOcrExecutorService(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
	defer service.Close()

	// OCRを実行
	result, err := service.ProcessPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を出力
	fmt.Print(result.FormatAsText())
}

func main() {
	// コマンドライン引数を解析
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		config.PrintUsage()
		os.Exit(1)
	}

	// ヘルプが要求された場合
	if cfg.Help {
		config.PrintUsage()
		return
	}

	// OCR実行を処理
	handleOcrExecution(cfg)
}
