package repositories

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// #==============================================================#
// ##       Mocks for FileReader                                 ##
// #==============================================================#
// MockFileReader はFileReaderRepositoryのモック実装
type MockFileReader struct {
	ReadFileFunc     func(filePath string) ([]byte, error)
	LoadJSONFileFunc func(filePath string) (map[string]interface{}, error)
}

func (m *MockFileReader) ReadFile(filePath string) ([]byte, error) {
	if m.ReadFileFunc != nil {
		return m.ReadFileFunc(filePath)
	}
	return nil, fmt.Errorf("ReadFileFunc not implemented")
}

func (m *MockFileReader) LoadJSONFile(filePath string) (map[string]interface{}, error) {
	if m.LoadJSONFileFunc != nil {
		return m.LoadJSONFileFunc(filePath)
	}
	return nil, fmt.Errorf("LoadJSONFileFunc not implemented")
}

// #==============================================================#
// ##       Interfaces for FileReader                            ##
// #==============================================================#
// FileReaderRepository はファイル読み込み操作のインターフェースです
type FileReaderRepository interface {
	ReadFile(filePath string) ([]byte, error)
	LoadJSONFile(filePath string) (map[string]interface{}, error)
}

// #==============================================================#
// ##       Implementations for FileReader                       ##
// #==============================================================#
// OSFileReader はOSファイルシステムを使用したFileReaderRepositoryの実装です
type OSFileReader struct{}

// NewOSFileReader は新しいOSFileReaderインスタンスを作成します
func NewOSFileReader() FileReaderRepository {
	return &OSFileReader{}
}

// ReadFile はファイルを読み込んでバイト配列として返します
func (r *OSFileReader) ReadFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("ファイルのオープンに失敗しました: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("ファイルの読み込みに失敗しました: %w", err)
	}

	return data, nil
}

// LoadJSONFile はJSONファイルを読み込んでmapとして返します
func (r *OSFileReader) LoadJSONFile(filePath string) (map[string]interface{}, error) {
	data, err := r.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("JSONのパースに失敗しました: %w", err)
	}

	return result, nil
}
