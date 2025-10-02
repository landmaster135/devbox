package main

import (
	"fmt"
	"os"

	"github.com/landmaster135/devbox/internal/git_diff_recorder/config"
	"github.com/landmaster135/devbox/internal/git_diff_recorder/usecases"
)

// runOutputMode は出力モードを実行する
func runOutputMode(cfg *config.Config) error {
	// 生成サービスを作成
	service := usecases.NewGitDiffGeneratorService(cfg.GitDir, cfg)

	// 詳細差分を取得
	detailedDiff, err := service.GetCurrentDetailedDiff()
	if err != nil {
		return err
	}

	// 詳細差分のみを出力
	fmt.Print(detailedDiff)
	return nil
}

// runReadMode は読み取りモードを実行する
func runReadMode(cfg *config.Config) error {
	// 読み取りサービスを作成
	readerService := usecases.NewDiffReaderService(cfg)

	// 詳細差分を読み取り表示
	return readerService.ReadAndDisplayDetailedDiff()
}

// runRecordMode は記録モードを実行する
func runRecordMode(cfg *config.Config) error {
	var workingDir string
	var err error

	if cfg.GitDir != "" {
		// 指定されたディレクトリを使用
		workingDir = cfg.GitDir
	} else {
		// 現在のディレクトリを使用（後方互換性）
		workingDir, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("現在のディレクトリを取得できませんでした: %w", err)
		}
	}

	// 記録サービスを作成
	service := usecases.NewGitDiffRecorderService(workingDir, cfg)

	// 差分を記録
	return service.RecordDiff()
}

func main() {
	// コマンドライン引数を解析
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		config.PrintUsage()
		os.Exit(1)
	}

	if cfg.GenMode {
		// 出力モード
		if err := runOutputMode(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			os.Exit(1)
		}
	} else if cfg.ReadMode {
		// 読み取りモード
		if err := runReadMode(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			os.Exit(1)
		}
	} else {
		// 記録モード（既存機能）
		if err := runRecordMode(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			os.Exit(1)
		}
	}
}
