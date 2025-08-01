package main

import (
	"fmt"
	"os"

	"github.com/landmaster135/devbox/internal/mcp_remote/config"
	"github.com/landmaster135/devbox/internal/mcp_remote/usecases"
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

	// プロキシサービスを初期化
	proxyService := usecases.NewProxyService()

	// プロキシを実行
	if err := proxyService.RunProxy(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "プロキシの実行中にエラーが発生しました: %v\n", err)
		os.Exit(1)
	}
}
