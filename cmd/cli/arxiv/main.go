package main

import (
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/arxiv/config"
	usecases "github.com/landmaster135/devbox/internal/arxiv/usecases"
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
	case "search":
		handleSearch(cfg)
	case "get_by_id":
		handleGetById(cfg)
	default:
		fmt.Fprintf(os.Stderr, "エラー: 未対応の操作タイプです: %s\n", cfg.Operation)
		config.PrintUsage()
		os.Exit(1)
	}
}

// handleSearch は論文検索を処理する
func handleSearch(cfg *config.Config) {
	// ArxivServiceを初期化
	service := usecases.NewArxivService()

	// 検索を実行
	result, err := service.HandleSearch(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を出力
	fmt.Print(result)
}

// handleGetById はID指定での論文取得を処理する
func handleGetById(cfg *config.Config) {
	// ArxivServiceを初期化
	service := usecases.NewArxivService()

	// ID指定取得を実行
	result, err := service.HandleSearch(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を出力
	fmt.Print(result)
}
