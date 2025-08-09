package main

import (
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/ops_for_golang/config"
	usecases "github.com/landmaster135/devbox/internal/ops_for_golang/usecases"
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
	case "test-coverage":
		handleTestCoverage(cfg)
	case "test-coverage-project":
		handleTestCoverageProject(cfg)
	case "coverage-func":
		handleCoverageFunc(cfg)
	case "run":
		handleGoRun(cfg)
	default:
		fmt.Fprintf(os.Stderr, "エラー: 未対応の操作タイプです: %s\n", cfg.Operation)
		config.PrintUsage()
		os.Exit(1)
	}
}

// handleTestCoverage はテストカバレッジを処理する
func handleTestCoverage(cfg *config.Config) {
	// GolangOpsServiceを初期化
	service := usecases.NewGolangOpsService()

	// テストカバレッジを実行
	result, err := service.HandleTestCoverage(cfg.Directory, cfg.GrepPattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を標準出力に表示
	fmt.Print(result)
}

// handleTestCoverageProject はプロジェクト全体のテストカバレッジを処理する
func handleTestCoverageProject(cfg *config.Config) {
	// GolangOpsServiceを初期化
	service := usecases.NewGolangOpsService()

	// プロジェクト全体のテストカバレッジを実行
	result, err := service.HandleTestCoverageProject(cfg.Directory)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を標準出力に表示
	fmt.Print(result)
}

// handleCoverageFunc はカバレッジ関数情報を処理する
func handleCoverageFunc(cfg *config.Config) {
	// GolangOpsServiceを初期化
	service := usecases.NewGolangOpsService()

	// カバレッジ関数情報を実行
	result, err := service.HandleCoverageFunc(cfg.CoverageFile, cfg.GrepPattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を標準出力に表示
	fmt.Print(result)
}

// handleGoRun はgo runを処理する
func handleGoRun(cfg *config.Config) {
	// GolangOpsServiceを初期化
	service := usecases.NewGolangOpsService()

	// go runを実行
	result, err := service.HandleGoRun(cfg.ExecutionFile, cfg.Parameters)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を標準出力に表示
	fmt.Print(result)
}
