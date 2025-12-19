package main

import (
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/data_converter/config"
	usecases "github.com/landmaster135/devbox/internal/data_converter/usecases"
)

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

	// データ変換サービスを初期化
	service := usecases.NewDataConverterService()

	// データ変換を実行
	result, err := service.ConvertData(cfg.InputFormat, cfg.OutputFormat, cfg.Input, cfg.InputFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を出力
	fmt.Print(result)
}
