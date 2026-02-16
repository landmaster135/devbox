package filesystem

import (
	"fmt"

	"github.com/landmaster135/devbox/internal/file_line_deduper/domain/models"
)

// MockRepository は Repository インターフェースのモック実装です。
type MockRepository struct {
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

func (m *MockRepository) ReadFile(path string) (*models.FileContent, error) {
	if m.ReadFileFunc == nil {
		return nil, fmt.Errorf("ReadFileFunc が未設定です")
	}
	return m.ReadFileFunc(path)
}

func (m *MockRepository) WriteFile(path string, content *models.FileContent) error {
	if m.WriteFileFunc == nil {
		return fmt.Errorf("WriteFileFunc が未設定です")
	}
	return m.WriteFileFunc(path, content)
}

func (m *MockRepository) FileExists(path string) bool {
	if m.FileExistsFunc == nil {
		return false
	}
	return m.FileExistsFunc(path)
}

func (m *MockRepository) FindFilesByExt(dirPath, ext string) ([]string, error) {
	if m.FindFilesByExtFunc == nil {
		return nil, fmt.Errorf("FindFilesByExtFunc が未設定です")
	}
	return m.FindFilesByExtFunc(dirPath, ext)
}

func (m *MockRepository) HasFilesWithExt(dirPath, ext string) (bool, error) {
	if m.HasFilesWithExtFunc == nil {
		return false, fmt.Errorf("HasFilesWithExtFunc が未設定です")
	}
	return m.HasFilesWithExtFunc(dirPath, ext)
}

func (m *MockRepository) ReadJSONFile(path string) (interface{}, error) {
	if m.ReadJSONFileFunc == nil {
		return nil, fmt.Errorf("ReadJSONFileFunc が未設定です")
	}
	return m.ReadJSONFileFunc(path)
}

func (m *MockRepository) GetDirectoryPath(path string) string {
	if m.GetDirectoryPathFunc == nil {
		return ""
	}
	return m.GetDirectoryPathFunc(path)
}

func (m *MockRepository) CreateDirectory(dirPath string) error {
	if m.CreateDirectoryFunc == nil {
		return fmt.Errorf("CreateDirectoryFunc が未設定です")
	}
	return m.CreateDirectoryFunc(dirPath)
}

func (m *MockRepository) ReadDir(dirPath string) ([]*models.DirEntry, error) {
	if m.ReadDirFunc == nil {
		return nil, fmt.Errorf("ReadDirFunc が未設定です")
	}
	return m.ReadDirFunc(dirPath)
}
