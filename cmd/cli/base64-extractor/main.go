package main

import (
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/base64_extractor/config"
	usecases "github.com/landmaster135/devbox/internal/base64_extractor/usecases"
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
	service := usecases.NewBase64ExtractorService(cfg.Path, cfg.Recursive)

	// 画像ファイルを抽出してbase64に変換
	result, err := service.ExtractFromPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を出力
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
