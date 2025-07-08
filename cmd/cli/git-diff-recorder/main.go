package main

import (
	"fmt"
	"os"

	"github.com/landmaster135/devbox/internal/git_diff_recorder/config"
	"github.com/landmaster135/devbox/internal/git_diff_recorder/usecases"
)

func main() {
	// コマンドライン引数を解析
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		config.PrintUsage()
		os.Exit(1)
	}

	// 現在のディレクトリを取得
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: 現在のディレクトリを取得できませんでした: %v\n", err)
		os.Exit(1)
	}

	// サービスを作成
	service := usecases.NewGitDiffRecorderService(workingDir, cfg)

	// 差分を記録
	if err := service.RecordDiff(); err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
}
