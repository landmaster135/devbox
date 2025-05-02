package services

import "github.com/landmaster135/devbox/internal/domain/models"

// ///////////////////////
// MockAPIRepository はAPIRepositoryのモック実装です
// ///////////////////////
type MockAPIRepository struct {
	LoadJSONFileFunc func(filePath string) ([]byte, error)
	SendRequestFunc  func(request *models.APIRequest) (*models.APIResponse, error)
}

// LoadJSONFile はJSONファイルを読み込むモックメソッドです
func (m *MockAPIRepository) LoadJSONFile(filePath string) ([]byte, error) {
	return m.LoadJSONFileFunc(filePath)
}

// SendRequest はAPIリクエストを送信するモックメソッドです
func (m *MockAPIRepository) SendRequest(request *models.APIRequest) (*models.APIResponse, error) {
	return m.SendRequestFunc(request)
}

// ///////////////////////
// MockFileRepository はFileRepositoryインターフェースのモック実装です
// ///////////////////////
type MockFileRepository struct {
	ReadFileFunc         func(path string) (*models.FileContent, error)
	WriteFileFunc        func(path string, content *models.FileContent) error
	FileExistsFunc       func(path string) bool
	FindFilesByExtFunc   func(dirPath, ext string) ([]string, error)
	HasFilesWithExtFunc  func(dirPath, ext string) (bool, error)
	ReadJSONFileFunc     func(path string) (interface{}, error)
	GetDirectoryPathFunc func(path string) string
	CreateDirectoryFunc  func(dirPath string) error
	ReadDirFunc          func(dirPath string) ([]*models.DirEntry, error)
}

// ReadFile はモックの実装です
func (m *MockFileRepository) ReadFile(path string) (*models.FileContent, error) {
	return m.ReadFileFunc(path)
}

// WriteFile はモックの実装です
func (m *MockFileRepository) WriteFile(path string, content *models.FileContent) error {
	return m.WriteFileFunc(path, content)
}

// FileExists はモックの実装です
func (m *MockFileRepository) FileExists(path string) bool {
	return m.FileExistsFunc(path)
}

// FindFilesByExt はモックの実装です
func (m *MockFileRepository) FindFilesByExt(dirPath, ext string) ([]string, error) {
	return m.FindFilesByExtFunc(dirPath, ext)
}

// HasFilesWithExt はモックの実装です
func (m *MockFileRepository) HasFilesWithExt(dirPath, ext string) (bool, error) {
	return m.HasFilesWithExtFunc(dirPath, ext)
}

// ReadJSONFile はモックの実装です
func (m *MockFileRepository) ReadJSONFile(path string) (interface{}, error) {
	return m.ReadJSONFileFunc(path)
}

// GetDirectoryPath はモックの実装です
func (m *MockFileRepository) GetDirectoryPath(path string) string {
	return m.GetDirectoryPathFunc(path)
}

// CreateDirectory はモックの実装です
func (m *MockFileRepository) CreateDirectory(dirPath string) error {
	return m.CreateDirectoryFunc(dirPath)
}

// ReadDir はモックの実装です
func (m *MockFileRepository) ReadDir(dirPath string) ([]*models.DirEntry, error) {
	return m.ReadDirFunc(dirPath)
}

// ///////////////////////
// モックEnvRepository //
// //////////////////////
type MockEnvRepository struct {
	LoadEnvFromYamlFunc func(path string) (*models.EnvConfig, error)
	SetEnvFunc          func(key, value string) error
	GetEnvFunc          func(key string) string
}

func (m *MockEnvRepository) LoadEnvFromYaml(path string) (*models.EnvConfig, error) {
	return m.LoadEnvFromYamlFunc(path)
}

func (m *MockEnvRepository) SetEnv(key, value string) error {
	return m.SetEnvFunc(key, value)
}

func (m *MockEnvRepository) GetEnv(key string) string {
	return m.GetEnvFunc(key)
}

// ///////////////////////
// MockJSONRepository はJSONRepositoryインターフェースのモック実装です
// ///////////////////////
type MockJSONRepository struct {
	ConvertFileFunc     func(filePath, key string, dryRun bool) (bool, error)
	ProcessJSONDataFunc func(data interface{}, targetKey string) (interface{}, bool)
	FindJSONFilesFunc   func(dirPath string, recursive bool) ([]string, error)
}

// ConvertFile は単一のJSONファイルを処理します
func (m *MockJSONRepository) ConvertFile(filePath, key string, dryRun bool) (bool, error) {
	if m.ConvertFileFunc != nil {
		return m.ConvertFileFunc(filePath, key, dryRun)
	}
	return false, nil
}

// ProcessJSONData はJSONデータを再帰的に処理します
func (m *MockJSONRepository) ProcessJSONData(data interface{}, targetKey string) (interface{}, bool) {
	if m.ProcessJSONDataFunc != nil {
		return m.ProcessJSONDataFunc(data, targetKey)
	}
	return data, false
}

// FindJSONFiles はディレクトリ内のJSONファイルを検索します
func (m *MockJSONRepository) FindJSONFiles(dirPath string, recursive bool) ([]string, error) {
	if m.FindJSONFilesFunc != nil {
		return m.FindJSONFilesFunc(dirPath, recursive)
	}
	return []string{}, nil
}

// ///////////////////////
// MockISO8601Repository はISO8601Repositoryインターフェースのモック実装です
// ///////////////////////
type MockISO8601Repository struct {
	ParseISO8601Func func(dateStr string) (int64, error)
}

// ParseISO8601 はISO8601形式の日時文字列をUNIXタイムスタンプに変換します
func (m *MockISO8601Repository) ParseISO8601(dateStr string) (int64, error) {
	if m.ParseISO8601Func != nil {
		return m.ParseISO8601Func(dateStr)
	}
	return 0, nil
}
