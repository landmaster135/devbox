package main

import (
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/git_info_retriever/config"
	usecases "github.com/landmaster135/devbox/internal/git_info_retriever/usecases"
)

// handleGetRepositoryInfo はリポジトリ情報取得を処理する
func handleGetRepositoryInfo(cfg *config.Config) {
	// GitInfoServiceを初期化
	service := usecases.NewService()

	result, err := service.RetrieveRepositoryInfo(cfg.Service, cfg.Token, cfg.SaveFilePath)
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

	// Git情報取得処理を実行
	handleGetRepositoryInfo(cfg)
}
