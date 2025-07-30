package main

import (
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/memory/config"
	usecases "github.com/landmaster135/devbox/internal/memory/usecases"
)

// createMemoryService は設定に基づいてMemoryServiceを作成する
func createMemoryService(cfg *config.Config) (*usecases.MemoryService, error) {
	t := cfg.StorageType
	switch t {
	case "file":
		return usecases.NewMemoryServiceWithFile(cfg.MemoryFile), nil
	case "valkey":
		valkeyURL := cfg.BuildValkeyURL()
		return usecases.NewMemoryServiceWithValkey(valkeyURL, cfg.ValkeyKey)
	default:
		return nil, fmt.Errorf("ストレージタイプが設定されていません: %v", t)
	}
}

// handleCreateEntities はエンティティ作成を処理する
func handleCreateEntities(cfg *config.Config) {
	service, err := createMemoryService(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "サービスの初期化エラー: %v\n", err)
		os.Exit(1)
	}

	result, err := service.HandleCreateEntities(cfg.Entities)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("作成されたエンティティ:\n%s\n", result)
}

// handleCreateRelations はリレーション作成を処理する
func handleCreateRelations(cfg *config.Config) {
	service, err := createMemoryService(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "サービスの初期化エラー: %v\n", err)
		os.Exit(1)
	}

	result, err := service.HandleCreateRelations(cfg.Relations)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("作成されたリレーション:\n%s\n", result)
}

// handleAddObservations は観察事項追加を処理する
func handleAddObservations(cfg *config.Config) {
	service, err := createMemoryService(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "サービスの初期化エラー: %v\n", err)
		os.Exit(1)
	}

	result, err := service.HandleAddObservations(cfg.Observations)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("追加された観察事項:\n%s\n", result)
}

// handleDeleteEntities はエンティティ削除を処理する
func handleDeleteEntities(cfg *config.Config) {
	service, err := createMemoryService(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "サービスの初期化エラー: %v\n", err)
		os.Exit(1)
	}

	err = service.HandleDeleteEntities(cfg.EntityNames)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("エンティティを削除しました: %s\n", cfg.EntityNames)
}

// handleDeleteObservations は観察事項削除を処理する
func handleDeleteObservations(cfg *config.Config) {
	service, err := createMemoryService(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "サービスの初期化エラー: %v\n", err)
		os.Exit(1)
	}

	err = service.HandleDeleteObservations(cfg.Deletions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("観察事項を削除しました\n")
}

// handleDeleteRelations はリレーション削除を処理する
func handleDeleteRelations(cfg *config.Config) {
	service, err := createMemoryService(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "サービスの初期化エラー: %v\n", err)
		os.Exit(1)
	}

	err = service.HandleDeleteRelations(cfg.Relations)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("リレーションを削除しました\n")
}

// handleReadGraph は知識グラフ全体の読み取りを処理する
func handleReadGraph(cfg *config.Config) {
	service, err := createMemoryService(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "サービスの初期化エラー: %v\n", err)
		os.Exit(1)
	}

	result, err := service.HandleReadGraph()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("知識グラフ:\n%s\n", result)
}

// handleSearchNodes はノード検索を処理する
func handleSearchNodes(cfg *config.Config) {
	service, err := createMemoryService(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "サービスの初期化エラー: %v\n", err)
		os.Exit(1)
	}

	result, err := service.HandleSearchNodes(cfg.Query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("検索結果 (クエリ: %s):\n%s\n", cfg.Query, result)
}

// handleOpenNodes は特定ノード取得を処理する
func handleOpenNodes(cfg *config.Config) {
	service, err := createMemoryService(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "サービスの初期化エラー: %v\n", err)
		os.Exit(1)
	}

	result, err := service.HandleOpenNodes(cfg.Names)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("指定されたノード (%s):\n%s\n", cfg.Names, result)
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

	// 操作タイプに応じて処理を実行
	switch cfg.Operation {
	case "create-entities":
		handleCreateEntities(cfg)
	case "create-relations":
		handleCreateRelations(cfg)
	case "add-observations":
		handleAddObservations(cfg)
	case "delete-entities":
		handleDeleteEntities(cfg)
	case "delete-observations":
		handleDeleteObservations(cfg)
	case "delete-relations":
		handleDeleteRelations(cfg)
	case "read-graph":
		handleReadGraph(cfg)
	case "search-nodes":
		handleSearchNodes(cfg)
	case "open-nodes":
		handleOpenNodes(cfg)
	default:
		fmt.Fprintf(os.Stderr, "エラー: 未対応の操作タイプです: %s\n", cfg.Operation)
		config.PrintUsage()
		os.Exit(1)
	}
}
