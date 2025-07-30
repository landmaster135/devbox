package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/valkey/config"
	loggerRepo "github.com/landmaster135/devbox/internal/valkey/infrastructure/logger/repository"
	valkeyRepo "github.com/landmaster135/devbox/internal/valkey/infrastructure/valkey/repository"
	usecases "github.com/landmaster135/devbox/internal/valkey/usecases"
)

// handleGetKeys はキー取得を処理する
func handleGetKeys(cfg *config.Config, service *usecases.DataService) {
	ctx := context.Background()

	keys, err := service.GetKeys(ctx, cfg.Pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("パターン '%s' に一致するキー (%d件):\n", cfg.Pattern, len(keys))
	for _, key := range keys {
		fmt.Printf("  %s\n", key)
	}
}

// handleGetValue は値取得を処理する
func handleGetValue(cfg *config.Config, service *usecases.DataService) {
	ctx := context.Background()

	value, err := service.GetValue(ctx, cfg.Key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("キー '%s' の値: %s\n", cfg.Key, value)
}

// handleGetType は型取得を処理する
func handleGetType(cfg *config.Config, service *usecases.DataService) {
	ctx := context.Background()

	keyType, err := service.GetType(ctx, cfg.Key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("キー '%s' の型: %s\n", cfg.Key, keyType)
}

// handleSetValue は値設定を処理する
func handleSetValue(cfg *config.Config, service *usecases.DataService) {
	ctx := context.Background()

	err := service.SetValue(ctx, cfg.Key, []byte(cfg.Value))
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("キー '%s' に値を設定しました\n", cfg.Key)
}

// handleDeleteKey はキー削除を処理する
func handleDeleteKey(cfg *config.Config, service *usecases.DataService) {
	ctx := context.Background()

	deleted, err := service.DeleteKey(ctx, cfg.Key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	if deleted {
		fmt.Printf("キー '%s' を削除しました\n", cfg.Key)
	} else {
		fmt.Printf("キー '%s' は存在しませんでした\n", cfg.Key)
	}
}

// handleDeleteKeys は複数キー削除を処理する
func handleDeleteKeys(cfg *config.Config, service *usecases.DataService) {
	ctx := context.Background()

	results, err := service.DeleteKeys(ctx, cfg.Keys)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	deletedCount := 0
	for key, deleted := range results {
		if deleted {
			fmt.Printf("キー '%s' を削除しました\n", key)
			deletedCount++
		} else {
			fmt.Printf("キー '%s' は存在しませんでした\n", key)
		}
	}

	fmt.Printf("合計 %d 個のキーを削除しました\n", deletedCount)
}

// handleSelectKeys は値選択を処理する
func handleSelectKeys(cfg *config.Config, service *usecases.DataService) {
	ctx := context.Background()

	result, err := service.SelectKeys(ctx, cfg.Key, cfg.Keys, cfg.Pattern, cfg.All)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果をJSON形式で出力
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "JSON変換エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("選択された値:\n%s\n", string(jsonData))
}

// handleGetAllValues は全値取得を処理する
func handleGetAllValues(cfg *config.Config, service *usecases.DataService, logger loggerRepo.Logger) {
	ctx := context.Background()

	result, err := service.GetAllValues(ctx, cfg.Keys, logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果をJSON形式で出力
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "JSON変換エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("取得された全値:\n%s\n", string(jsonData))
}

// handleDeleteData はデータ削除を処理する
func handleDeleteData(cfg *config.Config, service *usecases.DataService, logger loggerRepo.Logger) {
	ctx := context.Background()

	result, err := service.DeleteData(ctx, cfg.Key, cfg.Keys, cfg.Pattern, cfg.DryRun, logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果をJSON形式で出力
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "JSON変換エラー: %v\n", err)
		os.Exit(1)
	}

	if cfg.DryRun {
		fmt.Printf("削除対象データ（ドライラン）:\n%s\n", string(jsonData))
	} else {
		fmt.Printf("削除結果:\n%s\n", string(jsonData))
	}
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

	// Valkey接続URLを構築
	valkeyURL := cfg.BuildValkeyURL()

	// リポジトリを初期化
	repo, err := valkeyRepo.NewDataRepository(valkeyURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "リポジトリの初期化に失敗しました: %v\n", err)
		os.Exit(1)
	}

	// サービスを初期化
	service := usecases.NewDataService(repo)

	// ロガーを初期化
	logger := loggerRepo.NewDefaultLogger()

	// 操作タイプに応じて処理を実行
	switch cfg.Operation {
	case "get-keys":
		if cfg.Pattern == "" {
			fmt.Fprintf(os.Stderr, "エラー: get-keys操作にはpatternが必要です\n")
			os.Exit(1)
		}
		handleGetKeys(cfg, service)
	case "get-value":
		if cfg.Key == "" {
			fmt.Fprintf(os.Stderr, "エラー: get-value操作にはkeyが必要です\n")
			os.Exit(1)
		}
		handleGetValue(cfg, service)
	case "get-type":
		if cfg.Key == "" {
			fmt.Fprintf(os.Stderr, "エラー: get-type操作にはkeyが必要です\n")
			os.Exit(1)
		}
		handleGetType(cfg, service)
	case "set-value":
		if cfg.Key == "" || cfg.Value == "" {
			fmt.Fprintf(os.Stderr, "エラー: set-value操作にはkeyとvalueが必要です\n")
			os.Exit(1)
		}
		handleSetValue(cfg, service)
	case "delete-key":
		if cfg.Key == "" {
			fmt.Fprintf(os.Stderr, "エラー: delete-key操作にはkeyが必要です\n")
			os.Exit(1)
		}
		handleDeleteKey(cfg, service)
	case "delete-keys":
		if len(cfg.Keys) == 0 {
			fmt.Fprintf(os.Stderr, "エラー: delete-keys操作にはkeysが必要です\n")
			os.Exit(1)
		}
		handleDeleteKeys(cfg, service)
	case "select-keys":
		handleSelectKeys(cfg, service)
	case "get-all-values":
		if len(cfg.Keys) == 0 {
			fmt.Fprintf(os.Stderr, "エラー: get-all-values操作にはkeysが必要です\n")
			os.Exit(1)
		}
		handleGetAllValues(cfg, service, logger)
	case "delete-data":
		handleDeleteData(cfg, service, logger)
	default:
		fmt.Fprintf(os.Stderr, "エラー: 未対応の操作タイプです: %s\n", cfg.Operation)
		config.PrintUsage()
		os.Exit(1)
	}
}
