package main

import (
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/git_info_retriever/config"
	usecases "github.com/landmaster135/devbox/internal/git_info_retriever/usecases"
)

// handleRetrieveRepositoryInfo はリポジトリ情報取得を処理する
func handleRetrieveRepositoryInfo(cfg *config.Config) {
	// GitInfoServiceを初期化
	service := usecases.NewService()

	result, err := service.RetrieveRepositoryInfo(cfg.Service, cfg.Token, cfg.SaveFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(result)
}

// handleArchiveRepositories はリポジトリアーカイブ処理を処理する
func handleArchiveRepositories(cfg *config.Config) {
	// GitInfoServiceを初期化
	service := usecases.NewService()

	result, err := service.ArchiveRepositories(cfg.Service, cfg.Token, cfg.OutputCommandFilePath, cfg.ArchiveDir, cfg.SrcFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

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

	// 操作タイプによる分岐処理
	switch cfg.Operation {
	case "retrieve":
		handleRetrieveRepositoryInfo(cfg)
	case "archive":
		handleArchiveRepositories(cfg)
	default:
		fmt.Fprintf(os.Stderr, "エラー: 未対応の操作タイプです: %s\n", cfg.Operation)
		config.PrintUsage()
		os.Exit(1)
	}
}
