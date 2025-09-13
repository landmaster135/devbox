package main

import (
	"fmt"
	"os"

	"github.com/landmaster135/devbox/internal/secret_detector/config"
	"github.com/landmaster135/devbox/internal/secret_detector/usecases"
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

	// メイン処理を実行
	exitCode, err := service.ExecuteSecretDetection()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	os.Exit(exitCode)
}
