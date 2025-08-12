package main

import (
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/notion_sync/config"
	usecases "github.com/landmaster135/devbox/internal/notion_sync/usecases"
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

	// NotionSyncServiceを初期化
	service := usecases.NewNotionSyncService()

	// Notion同期を実行
	result, err := service.HandleNotionSync(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を出力
	fmt.Print(result)
}
