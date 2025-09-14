package main

import (
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/color_code_converter/config"
	usecases "github.com/landmaster135/devbox/internal/color_code_converter/usecases"
)

// handleColorConversion はカラーコード変換を処理する
func handleColorConversion(cfg *config.Config) {
	// ColorConverterServiceを初期化
	service := usecases.NewColorConverterService()

	// カラーコード変換を実行
	result, err := service.ConvertColorWithValidation(cfg.SrcFormat, cfg.DestFormat, cfg.Value)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を標準出力に表示
	fmt.Print(result)
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

	// カラーコード変換を実行
	handleColorConversion(cfg)
}
