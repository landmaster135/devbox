package services

import (
	"errors"
	"fmt"
	"testing"

	"github.com/landmaster135/devbox/internal/json_file_merger/domain/models"
)

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

func TestFileService_CreateRequestBodyFromDir(t *testing.T) {
	tests := []struct {
		name                string
		dirPath             string
		keyName             string
		outputPath          string
		findFilesResult     []string
		findFilesErr        error
		readJSONResult      interface{}
		readJSONErr         error
		writeFileErr        error
		getDirectoryPathErr error
		createDirectoryErr  error
		wantErr             bool
	}{
		{
			name:            "正常系",
			dirPath:         "/test/dir",
			keyName:         "pc_stats",
			outputPath:      "/test/output.json",
			findFilesResult: []string{"/test/dir/file1.json", "/test/dir/file2.json"},
			findFilesErr:    nil,
			readJSONResult:  map[string]interface{}{"id": 1, "name": "test"},
			readJSONErr:     nil,
			writeFileErr:    nil,
			wantErr:         false,
		},
		{
			name:            "ディレクトリが存在しない",
			dirPath:         "/non_existent",
			keyName:         "pc_stats",
			outputPath:      "/test/output.json",
			findFilesResult: nil,
			findFilesErr:    fmt.Errorf("指定されたディレクトリが存在しません"),
			readJSONResult:  nil,
			readJSONErr:     nil,
			writeFileErr:    nil,
			wantErr:         true,
		},
		{
			name:            "JSONファイルが存在しない",
			dirPath:         "/empty_dir",
			keyName:         "pc_stats",
			outputPath:      "/test/output.json",
			findFilesResult: []string{},
			findFilesErr:    nil,
			readJSONResult:  nil,
			readJSONErr:     nil,
			writeFileErr:    nil,
			wantErr:         true,
		},
		{
			name:            "JSONファイルの読み込みに失敗",
			dirPath:         "/test/dir",
			keyName:         "pc_stats",
			outputPath:      "/test/output.json",
			findFilesResult: []string{"/test/dir/file1.json"},
			findFilesErr:    nil,
			readJSONResult:  nil,
			readJSONErr:     fmt.Errorf("JSONファイルの読み込みに失敗しました"),
			writeFileErr:    nil,
			wantErr:         true,
		},
		{
			name:                "出力ファイルの書き込みに失敗",
			dirPath:             "/test/dir",
			keyName:             "pc_stats",
			outputPath:          "/test/output.json",
			findFilesResult:     []string{"/test/dir/file1.json"},
			findFilesErr:        nil,
			readJSONResult:      map[string]interface{}{"id": 1, "name": "test"},
			readJSONErr:         nil,
			writeFileErr:        fmt.Errorf("ファイルの書き込みに失敗しました"),
			getDirectoryPathErr: nil,
			createDirectoryErr:  nil,
			wantErr:             true,
		},
		{
			name:                "出力パスが指定されていない",
			dirPath:             "/test/dir",
			keyName:             "pc_stats",
			outputPath:          "", // 空文字列を指定
			findFilesResult:     []string{"/test/dir/file1.json"},
			findFilesErr:        nil,
			readJSONResult:      map[string]interface{}{"id": 1, "name": "test"},
			readJSONErr:         nil,
			writeFileErr:        nil,
			getDirectoryPathErr: nil,
			createDirectoryErr:  nil,
			wantErr:             false,
		},
		{
			name:                "出力ディレクトリの作成に失敗",
			dirPath:             "/test/dir",
			keyName:             "pc_stats",
			outputPath:          "/test/output.json",
			findFilesResult:     []string{"/test/dir/file1.json"},
			findFilesErr:        nil,
			readJSONResult:      map[string]interface{}{"id": 1, "name": "test"},
			readJSONErr:         nil,
			writeFileErr:        nil,
			getDirectoryPathErr: nil,
			createDirectoryErr:  fmt.Errorf("ディレクトリの作成に失敗しました"),
			wantErr:             true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックリポジトリの設定
			mockRepo := &MockFileRepository{
				FindFilesByExtFunc: func(dirPath, ext string) ([]string, error) {
					return tt.findFilesResult, tt.findFilesErr
				},
				ReadJSONFileFunc: func(path string) (interface{}, error) {
					return tt.readJSONResult, tt.readJSONErr
				},
				WriteFileFunc: func(path string, content *models.FileContent) error {
					return tt.writeFileErr
				},
				GetDirectoryPathFunc: func(path string) string {
					return "/test"
				},
				CreateDirectoryFunc: func(dirPath string) error {
					return tt.createDirectoryErr
				},
			}

			// テスト対象のサービスを作成
			service := NewFileService(mockRepo)

			// テスト実行
			_, err := service.CreateRequestBodyFromDir(tt.dirPath, tt.keyName, tt.outputPath)

			// エラーの検証
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateRequestBodyFromDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFileService_RemoveMatchingLines(t *testing.T) {
	tests := []struct {
		name         string
		filePath     string
		startPos     int
		endPos       int
		fileExists   bool
		readFileErr  error
		writeFileErr error
		fileContent  *models.FileContent
		wantCount    int
		wantErr      bool
	}{
		{
			name:         "正常系",
			filePath:     "test.txt",
			startPos:     5,
			endPos:       200,
			fileExists:   true,
			readFileErr:  nil,
			writeFileErr: nil,
			fileContent: models.NewFileContent([]string{
				"Line 1",
				"Line 2",
				"Line 3",
				"Line 2", // 重複行
			}),
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:         "ファイルが存在しない",
			filePath:     "non_existent.txt",
			startPos:     5,
			endPos:       200,
			fileExists:   false,
			readFileErr:  nil,
			writeFileErr: nil,
			fileContent:  nil,
			wantCount:    0,
			wantErr:      true,
		},
		{
			name:         "ファイル読み込みエラー",
			filePath:     "error.txt",
			startPos:     5,
			endPos:       200,
			fileExists:   true,
			readFileErr:  errors.New("読み込みエラー"),
			writeFileErr: nil,
			fileContent:  nil,
			wantCount:    0,
			wantErr:      true,
		},
		{
			name:         "ファイル書き込みエラー",
			filePath:     "write_error.txt",
			startPos:     5,
			endPos:       200,
			fileExists:   true,
			readFileErr:  nil,
			writeFileErr: errors.New("書き込みエラー"),
			fileContent: models.NewFileContent([]string{
				"Line 1",
				"Line 2",
			}),
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックリポジトリの設定
			mockRepo := &MockFileRepository{
				FileExistsFunc: func(path string) bool {
					return tt.fileExists
				},
				ReadFileFunc: func(path string) (*models.FileContent, error) {
					if tt.readFileErr != nil {
						return nil, tt.readFileErr
					}
					return tt.fileContent, nil
				},
				WriteFileFunc: func(path string, content *models.FileContent) error {
					return tt.writeFileErr
				},
			}

			// テスト対象のサービスを作成
			service := NewFileService(mockRepo)

			// テスト実行
			gotCount, err := service.RemoveMatchingLines(tt.filePath, tt.startPos, tt.endPos)

			// エラーの検証
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveMatchingLines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 正常系の場合は結果を検証
			if !tt.wantErr && gotCount != tt.wantCount {
				t.Errorf("RemoveMatchingLines() = %v, want %v", gotCount, tt.wantCount)
			}
		})
	}
}
