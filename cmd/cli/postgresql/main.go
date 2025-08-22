package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/postgresql/config"
	usecases "github.com/landmaster135/devbox/internal/postgresql/usecases"
)

func main() {
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		config.PrintUsage()
		os.Exit(1)
	}

	if cfg.Help {
		config.PrintUsage()
		return
	}

	switch cfg.Operation {
	case "dump":
		handleDump(cfg)
	default:
		fmt.Fprintf(os.Stderr, "エラー: 未対応の操作タイプです: %s\n", cfg.Operation)
		config.PrintUsage()
		os.Exit(1)
	}
}

func handleDump(cfg *config.Config) {
	// PostgreSQLサービスを初期化
	service, err := usecases.NewPostgreSQLService(cfg.DatabaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: PostgreSQLサービスの初期化に失敗しました: %v\n", err)
		os.Exit(1)
	}
	defer service.Close()

	// ダンプを実行
	result, err := service.HandleToDumpTable(context.Background(), cfg.TableName, cfg.OutputPath, cfg.Format, cfg.Limit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: テーブルダンプの実行に失敗しました: %v\n", err)
		os.Exit(1)
	}

	// 結果をJSON形式で標準出力に表示
	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: 結果のJSON変換に失敗しました: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(string(jsonResult))
}
