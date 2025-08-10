package main

import (
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/github/config"
	usecases "github.com/landmaster135/devbox/internal/github/usecases"
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
	case "list-issues":
		handleListIssues(cfg)
	default:
		fmt.Fprintf(os.Stderr, "エラー: 未対応の操作タイプです: %s\n", cfg.Operation)
		config.PrintUsage()
		os.Exit(1)
	}
}

// handleListIssues はイシュー一覧取得を処理する
func handleListIssues(cfg *config.Config) {
	// GitHubIssueServiceを初期化
	service := usecases.NewGitHubIssueService(cfg.Token)

	// イシュー一覧を取得
	result, err := service.HandleToListIssues(cfg.Owner, cfg.Repo, cfg.State, cfg.Sort, cfg.Direction, cfg.PerPage, cfg.Page)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を標準出力に表示
	fmt.Print(result)
}
