package main

import (
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/ocr_executor/config"
	usecases "github.com/landmaster135/devbox/internal/ocr_executor/usecases"
)

func main() {
	// コマンドライン引数を解析
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		config.PrintUsage()
		os.Exit(1)
	}

	// サービスを作成
	service := usecases.NewOCRExecutorService(
		cfg.Path,
		cfg.Recursive,
		cfg.GetTesseractLanguages(),
		cfg.OutputDir,
		cfg.OutputFormat,
	)

	// OCR処理を実行
	result, err := service.ExecuteFromPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を標準出力に表示
	var output string
	if cfg.OutputFormat == "json" {
		output, err = result.FormatAsJSON()
		if err != nil {
			fmt.Fprintf(os.Stderr, "JSON出力エラー: %v\n", err)
			os.Exit(1)
		}
	} else {
		output = result.FormatAsText()
	}

	fmt.Print(output)
}
