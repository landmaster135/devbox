package main

import (
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/git_commit_history_retriever/config"
	usecases "github.com/landmaster135/devbox/internal/git_commit_history_retriever/usecases"
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
	service := usecases.NewGitCommitHistoryService(cfg.GitDir, cfg)

	// コミット履歴と詳細を取得
	result, err := service.GetCommitHistoryWithDetails()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を出力
	fmt.Print(result)
}
