package usecases

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	config "github.com/landmaster135/devbox/internal/memory/config"
)

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
