package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	config "github.com/landmaster135/devbox/internal/memory/config"
	valkeyRepo "github.com/landmaster135/devbox/internal/memory/infrastructure/valkey/repository"
)

// #==============================================================#
// ##          Interface                                         ##
// #==============================================================#
// KnowledgeGraphRepository は知識グラフの永続化を抽象化するインターフェース
type KnowledgeGraphRepository interface {
	// LoadGraph は知識グラフを読み込む
	LoadGraph() (*KnowledgeGraph, error)

	// SaveGraph は知識グラフを保存する
	SaveGraph(graph *KnowledgeGraph) error
}

// #==============================================================#
// ##          FileRepository                                    ##
// #==============================================================#
// FileRepository はファイルベースの知識グラフリポジトリ
type FileRepository struct {
	fileReader config.FileReader
	fileWriter config.FileWriter
	memoryFile string
}

// インターフェースを実装していることを確認
var _ KnowledgeGraphRepository = (*FileRepository)(nil)

// NewFileRepository は新しいFileRepositoryを作成する
func NewFileRepository(memoryFile string) *FileRepository {
	return &FileRepository{
		fileReader: &config.StandardFileReader{},
		fileWriter: &config.StandardFileWriter{},
		memoryFile: memoryFile,
	}
}

// NewFileRepositoryWithDependencies は依存性注入版のFileRepositoryを作成する
func NewFileRepositoryWithDependencies(fileReader config.FileReader, fileWriter config.FileWriter, memoryFile string) *FileRepository {
	return &FileRepository{
		fileReader: fileReader,
		fileWriter: fileWriter,
		memoryFile: memoryFile,
	}
}

// LoadGraph は知識グラフをファイルから読み込む
func (r *FileRepository) LoadGraph() (*KnowledgeGraph, error) {
	// ファイルが存在しない場合は空のグラフを返す
	if _, err := os.Stat(r.memoryFile); os.IsNotExist(err) {
		return &KnowledgeGraph{
			Entities:  []Entity{},
			Relations: []Relation{},
		}, nil
	}

	data, err := r.fileReader.ReadFile(r.memoryFile)
	if err != nil {
		return nil, fmt.Errorf("ファイル読み込みエラー: %v", err)
	}

	var graph KnowledgeGraph
	if err := json.Unmarshal(data, &graph); err != nil {
		return nil, fmt.Errorf("JSON解析エラー: %v", err)
	}

	return &graph, nil
}

// SaveGraph は知識グラフをファイルに保存する
func (r *FileRepository) SaveGraph(graph *KnowledgeGraph) error {
	// ディレクトリが存在しない場合は作成
	dir := filepath.Dir(r.memoryFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("ディレクトリ作成エラー: %v", err)
	}

	data, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON変換エラー: %v", err)
	}

	if err := r.fileWriter.WriteFile(r.memoryFile, data, 0644); err != nil {
		return fmt.Errorf("ファイル書き込みエラー: %v", err)
	}

	return nil
}

// #==============================================================#
// ##          ValkeyRepository                                  ##
// #==============================================================#
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
