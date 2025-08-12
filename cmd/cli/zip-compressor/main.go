package main

import (
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/zip_compressor/config"
	usecases "github.com/landmaster135/devbox/internal/zip_compressor/usecases"
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

	// 操作タイプに応じて処理を実行
	switch cfg.Operation {
	case "compress":
		handleCompress(cfg)
	case "decompress":
		handleDecompress(cfg)
	default:
		fmt.Fprintf(os.Stderr, "エラー: 未対応の操作タイプです: %s\n", cfg.Operation)
		config.PrintUsage()
		os.Exit(1)
	}
}

// handleCompress はファイル/ディレクトリの圧縮を処理する
func handleCompress(cfg *config.Config) {
	// ZipCompressorServiceを初期化
	service := usecases.NewZipCompressorService()

	// 圧縮を実行
	result, err := service.HandleCompress(cfg.Path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を出力
	fmt.Print(result)
}

// handleDecompress はZipファイルの展開を処理する
func handleDecompress(cfg *config.Config) {
	// ZipCompressorServiceを初期化
	service := usecases.NewZipCompressorService()

	// 展開を実行
	result, err := service.HandleDecompress(cfg.Path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を出力
	fmt.Print(result)
}
