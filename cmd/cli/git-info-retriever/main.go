package main

import (
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/git_info_retriever/config"
	usecases "github.com/landmaster135/devbox/internal/git_info_retriever/usecases"
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

	// Git情報取得処理を実行
	handleGetRepositoryInfo(cfg)
}

// handleGetRepositoryInfo はリポジトリ情報取得を処理する
func handleGetRepositoryInfo(cfg *config.Config) {
	// GitInfoServiceを初期化
	service := usecases.NewService()

	// ファイル出力が指定されている場合
	if cfg.SaveFile != "" {
		_, err := service.GetAndSaveRepositoryInfo(cfg.Service, cfg.Token, cfg.SaveFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("結果をファイルに保存しました: %s\n", cfg.SaveFile)
		return
	}

	// 標準出力に表示
	result, err := service.GetRepositoryInfo(cfg.Service, cfg.Token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(result)
}
