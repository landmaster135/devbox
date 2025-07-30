package usecases

import (
	"context"
	"encoding/json"
	"fmt"

	valkeyRepo "github.com/landmaster135/devbox/internal/memory/infrastructure/valkey/repository"
)

// ValkeyRepository はValkeyベースの知識グラフリポジトリ
type ValkeyRepository struct {
	dataRepo valkeyRepo.DataRepository
	key      string
}

// インターフェースを実装していることを確認
var _ KnowledgeGraphRepository = (*ValkeyRepository)(nil)

// NewValkeyRepository は新しいValkeyRepositoryを作成する
func NewValkeyRepository(valkeyURL string) (*ValkeyRepository, error) {
	dataRepo, err := valkeyRepo.NewDataRepository(valkeyURL)
	if err != nil {
		return nil, fmt.Errorf("valkeyリポジトリの初期化に失敗しました: %v", err)
	}

	return &ValkeyRepository{
		dataRepo: dataRepo,
		key:      "knowledge_graph:main",
	}, nil
}

// NewValkeyRepositoryWithKey はカスタムキーを使用してValkeyRepositoryを作成する
func NewValkeyRepositoryWithKey(valkeyURL string, key string) (*ValkeyRepository, error) {
	dataRepo, err := valkeyRepo.NewDataRepository(valkeyURL)
	if err != nil {
		return nil, fmt.Errorf("valkeyリポジトリの初期化に失敗しました: %v", err)
	}

	return &ValkeyRepository{
		dataRepo: dataRepo,
		key:      key,
	}, nil
}

// LoadGraph は知識グラフをValkeyから読み込む
func (r *ValkeyRepository) LoadGraph() (*KnowledgeGraph, error) {
	ctx := context.Background()

	// キーが存在するかチェック
	keyType, err := r.dataRepo.GetType(ctx, r.key)
	if err != nil {
		return nil, fmt.Errorf("キーの型取得エラー: %v", err)
	}

	// キーが存在しない場合は空のグラフを返す
	if keyType == "none" {
		return &KnowledgeGraph{
			Entities:  []Entity{},
			Relations: []Relation{},
		}, nil
	}

	// データを取得
	data, err := r.dataRepo.GetValueAsByte(ctx, r.key)
	if err != nil {
		return nil, fmt.Errorf("データ取得エラー: %v", err)
	}

	// JSONをパース
	var graph KnowledgeGraph
	if err := json.Unmarshal(data, &graph); err != nil {
		return nil, fmt.Errorf("JSON解析エラー: %v", err)
	}

	return &graph, nil
}

// SaveGraph は知識グラフをValkeyに保存する
func (r *ValkeyRepository) SaveGraph(graph *KnowledgeGraph) error {
	ctx := context.Background()

	// JSONに変換
	data, err := json.Marshal(graph)
	if err != nil {
		return fmt.Errorf("JSON変換エラー: %v", err)
	}

	// Valkeyに保存
	if err := r.dataRepo.SetValue(ctx, r.key, data); err != nil {
		return fmt.Errorf("データ保存エラー: %v", err)
	}

	return nil
}

// GetKey は使用しているキーを返す
func (r *ValkeyRepository) GetKey() string {
	return r.key
}

// StartServer はValkeyサーバーを起動する
func (r *ValkeyRepository) StartServer() error {
	return r.dataRepo.StartServer()
}

// StopServer はValkeyサーバーを停止する
func (r *ValkeyRepository) StopServer() error {
	return r.dataRepo.StopServer()
}

// IsServerRunning はValkeyサーバーが起動しているかチェックする
func (r *ValkeyRepository) IsServerRunning() bool {
	return r.dataRepo.IsServerRunning()
}
