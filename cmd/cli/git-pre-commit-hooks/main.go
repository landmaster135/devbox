package main

import (
	"fmt"
	"os"

	"github.com/landmaster135/devbox/internal/git_pre_commit_hooks/config"
	"github.com/landmaster135/devbox/internal/git_pre_commit_hooks/usecases"
)

func main() {
	// コマンドライン引数を解析
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		config.PrintUsage()
		os.Exit(1)
	}

	// バージョンが要求された場合
	if cfg.Version {
		fmt.Println("Secret Detector v1.0.0")
		return
	}

	// ヘルプが要求された場合
	if cfg.Help {
		config.PrintUsage()
		return
	}

	// 依存性を初期化
	commandExecutor := usecases.NewCommandExecutor()
	outputWriter := usecases.NewOutputWriter()

	// シークレット検知サービスを初期化
	service := usecases.NewSecretDetectorService(cfg.Verbose, cfg.DryRun, cfg.ConfigFile, commandExecutor, outputWriter)

	// シークレット検知を実行
	secretExitCode, err := service.ExecuteSecretDetection()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// ホームパス検知を実行
	homePathExitCode, err := service.ExecuteHomePathDetection()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// どちらかで問題が検出された場合はコミットを拒否
	if secretExitCode != 0 || homePathExitCode != 0 {
		fmt.Println()
		fmt.Printf("%s🚫 FINAL RESULT: COMMIT BLOCKED%s\n", "\033[31m", "\033[0m")
		fmt.Println("One or more issues were detected. Please resolve them before committing.")
		os.Exit(1)
	}

	fmt.Println()
	fmt.Printf("%s✅ FINAL RESULT: COMMIT ALLOWED%s\n", "\033[32m", "\033[0m")
	fmt.Println("All checks passed successfully.")
	os.Exit(0)
}
