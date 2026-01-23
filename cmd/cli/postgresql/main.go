package main

import (
	"context"
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/postgresql/config"
	usecases "github.com/landmaster135/devbox/internal/postgresql/usecases"
)

func handleDump(cfg *config.Config) {
	// ダンプを実行
	result, err := usecases.HandleToDumpTable(context.Background(), cfg.DatabaseURL, cfg.TableName, cfg.OutputPath, cfg.Format, cfg.Limit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: テーブルダンプの実行に失敗しました: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(result)
}

func handleDumpAllTables(cfg *config.Config) {
	// 全テーブルダンプを実行
	result, err := usecases.HandleToDumpAllTables(context.Background(), cfg.DatabaseURL, cfg.OutputPath, cfg.Format, cfg.Limit, cfg.Concurrency)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: 全テーブルダンプに失敗しました: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(result)
}

func handleListTablesMinimum(cfg *config.Config) {
	// PostgreSQLサービスを初期化
	service, err := usecases.NewPostgreSQLService(cfg.DatabaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: PostgreSQLサービスの初期化に失敗しました: %v\n", err)
		os.Exit(1)
	}
	defer service.Close()

	// テーブル一覧（最小限）を取得
	result, err := service.HandleToListTablesMinimum(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: テーブル一覧の取得に失敗しました: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(result)
}

func handleListTables(cfg *config.Config) {
	// PostgreSQLサービスを初期化
	service, err := usecases.NewPostgreSQLService(cfg.DatabaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: PostgreSQLサービスの初期化に失敗しました: %v\n", err)
		os.Exit(1)
	}
	defer service.Close()

	// テーブル一覧（詳細）を取得
	if cfg.Format == "text" {
		// テキスト形式で取得
		result, err := service.HandleToListTables(context.Background())
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: テーブル一覧の取得に失敗しました: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(result)
	} else {
		// JSON形式で取得
		result, err := service.HandleGetAllTableSummaries(context.Background())
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: テーブル一覧の取得に失敗しました: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(result)
	}
}

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
	case "dump-all-tables":
		handleDumpAllTables(cfg)
	case "list-tables-minimum":
		handleListTablesMinimum(cfg)
	case "list-tables":
		handleListTables(cfg)
	default:
		fmt.Fprintf(os.Stderr, "エラー: 未対応の操作タイプです: %s\n", cfg.Operation)
		config.PrintUsage()
		os.Exit(1)
	}
}
