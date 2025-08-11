package usecases

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// テスト用のMarkdownコンテンツ定数
const (
	validMarkdownWithMarker = `# テストブログ記事

これは前書きの部分です。この部分は抽出されません。

## 概要

この記事の概要です。この部分も抽出されません。

# Content

## はじまり

ここからが実際のブログコンテンツです。この部分が抽出されます。

### セクション1

これは重要なコンテンツの一部です。

- リスト項目1
- リスト項目2
- リスト項目3

### セクション2

` + "```javascript\nconsole.log(\"Hello, World!\");\n```" + `

このコードも含まれます。

## まとめ

これがまとめの部分です。すべて抽出対象に含まれます。`

	markdownWithoutMarker = `# テストブログ記事

これは通常のMarkdownファイルです。

## セクション1

コンテンツマーカーが含まれていません。

## セクション2

このファイルは抽出対象になりません。`

	incompleteMarker = `# テストブログ記事

前書き部分

# Content
## はじまり

空行がないため、マーカーとして認識されません。`

	expectedExtractedContent = `# Content

## はじまり

ここからが実際のブログコンテンツです。この部分が抽出されます。

### セクション1

これは重要なコンテンツの一部です。

- リスト項目1
- リスト項目2
- リスト項目3

### セクション2

` + "```javascript\nconsole.log(\"Hello, World!\");\n```" + `

このコードも含まれます。

## まとめ

これがまとめの部分です。すべて抽出対象に含まれます。`
)

// MockFileOperator はFileOperatorのモック実装
type MockFileOperator struct {
	ReadFileFunc  func(filename string) ([]byte, error)
	WriteFileFunc func(filename string, data []byte, perm os.FileMode) error
	MkdirAllFunc  func(path string, perm os.FileMode) error
	WalkDirFunc   func(root string, fn fs.WalkDirFunc) error
	StatFunc      func(name string) (os.FileInfo, error)

	// テスト用のファイルシステム状態
	files       map[string][]byte
	directories map[string]bool
	errors      map[string]error
}

// NewMockFileOperator は新しいMockFileOperatorを作成する
func NewMockFileOperator() *MockFileOperator {
	return &MockFileOperator{
		files:       make(map[string][]byte),
		directories: make(map[string]bool),
		errors:      make(map[string]error),
	}
}

func (m *MockFileOperator) ReadFile(filename string) ([]byte, error) {
	if m.ReadFileFunc != nil {
		return m.ReadFileFunc(filename)
	}
	if err, exists := m.errors[filename]; exists {
		return nil, err
	}
	if content, exists := m.files[filename]; exists {
		return content, nil
	}
	return nil, fmt.Errorf("file not found: %s", filename)
}

func (m *MockFileOperator) WriteFile(filename string, data []byte, perm os.FileMode) error {
	if m.WriteFileFunc != nil {
		return m.WriteFileFunc(filename, data, perm)
	}
	if err, exists := m.errors[filename]; exists {
		return err
	}
	m.files[filename] = data
	return nil
}

func (m *MockFileOperator) MkdirAll(path string, perm os.FileMode) error {
	if m.MkdirAllFunc != nil {
		return m.MkdirAllFunc(path, perm)
	}
	if err, exists := m.errors[path]; exists {
		return err
	}
	m.directories[path] = true
	return nil
}

func (m *MockFileOperator) WalkDir(root string, fn fs.WalkDirFunc) error {
	if m.WalkDirFunc != nil {
		return m.WalkDirFunc(root, fn)
	}
	if err, exists := m.errors[root]; exists {
		return err
	}

	// モックファイルシステムを走査
	for filePath := range m.files {
		if strings.HasPrefix(filePath, root) {
			info := &MockFileInfo{name: filepath.Base(filePath), isDir: false}
			entry := &MockDirEntry{info: info}
			if err := fn(filePath, entry, nil); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *MockFileOperator) Stat(name string) (os.FileInfo, error) {
	if m.StatFunc != nil {
		return m.StatFunc(name)
	}
	if err, exists := m.errors[name]; exists {
		return nil, err
	}
	if _, exists := m.files[name]; exists {
		return &MockFileInfo{name: filepath.Base(name), isDir: false}, nil
	}
	if _, exists := m.directories[name]; exists {
		return &MockFileInfo{name: filepath.Base(name), isDir: true}, nil
	}
	return nil, os.ErrNotExist
}

// SetFile はモックファイルシステムにファイルを追加する
func (m *MockFileOperator) SetFile(path string, content []byte) {
	m.files[path] = content
}

// SetDirectory はモックファイルシステムにディレクトリを追加する
func (m *MockFileOperator) SetDirectory(path string) {
	m.directories[path] = true
}

// SetError は特定のパスでエラーを発生させる
func (m *MockFileOperator) SetError(path string, err error) {
	m.errors[path] = err
}

// MockFileInfo はos.FileInfoのモック実装
type MockFileInfo struct {
	name  string
	isDir bool
}

func (m *MockFileInfo) Name() string       { return m.name }
func (m *MockFileInfo) Size() int64        { return 0 }
func (m *MockFileInfo) Mode() os.FileMode  { return 0644 }
func (m *MockFileInfo) ModTime() time.Time { return time.Now() }
func (m *MockFileInfo) IsDir() bool        { return m.isDir }
func (m *MockFileInfo) Sys() interface{}   { return nil }

// MockDirEntry はfs.DirEntryのモック実装
type MockDirEntry struct {
	info *MockFileInfo
}

func (m *MockDirEntry) Name() string               { return m.info.Name() }
func (m *MockDirEntry) IsDir() bool                { return m.info.IsDir() }
func (m *MockDirEntry) Type() fs.FileMode          { return m.info.Mode().Type() }
func (m *MockDirEntry) Info() (fs.FileInfo, error) { return m.info, nil }

func TestNewService_Normal(t *testing.T) {
	service := NewService()

	if service == nil {
		t.Fatal("Expected service to be non-nil")
	}

	if service.fileOperator == nil {
		t.Error("Expected fileOperator to be set")
	}

	// DefaultFileOperatorが設定されていることを確認
	if _, ok := service.fileOperator.(*DefaultFileOperator); !ok {
		t.Error("Expected DefaultFileOperator to be set")
	}
}

func TestNewServiceWithFileOperator_Normal(t *testing.T) {
	mockFileOperator := NewMockFileOperator()
	service := NewServiceWithFileOperator(mockFileOperator)

	if service == nil {
		t.Fatal("Expected service to be non-nil")
	}

	if service.fileOperator != mockFileOperator {
		t.Error("Expected fileOperator to be the provided mock")
	}
}

func TestExtractContentAfterMarker_Normal(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		expectedResult string
		expectError    bool
		errorMessage   string
	}{
		{
			name:           "ValidMarker_Normal",
			content:        validMarkdownWithMarker,
			expectedResult: expectedExtractedContent,
			expectError:    false,
		},
		{
			name:         "NoMarker_Error",
			content:      markdownWithoutMarker,
			expectError:  true,
			errorMessage: "指定されたマーカーが見つかりません",
		},
		{
			name:         "IncompleteMarker_Error",
			content:      incompleteMarker,
			expectError:  true,
			errorMessage: "指定されたマーカーが見つかりません",
		},
		{
			name:         "EmptyContent_Error",
			content:      "",
			expectError:  true,
			errorMessage: "指定されたマーカーが見つかりません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService()
			result, err := service.extractContentAfterMarker(tt.content)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if tt.errorMessage != "" && !strings.Contains(err.Error(), tt.errorMessage) {
					t.Errorf("Expected error message to contain '%s', got '%s'", tt.errorMessage, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if result != tt.expectedResult {
					t.Errorf("Expected result:\n%s\n\nGot:\n%s", tt.expectedResult, result)
				}
			}
		})
	}
}

func TestValidateDirectory_Normal(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func(*MockFileOperator)
		dirPath      string
		description  string
		expectError  bool
		errorMessage string
	}{
		{
			name: "ValidDirectory_Normal",
			setupMock: func(mock *MockFileOperator) {
				mock.SetDirectory("/valid/dir")
			},
			dirPath:     "/valid/dir",
			description: "テストディレクトリ",
			expectError: false,
		},
		{
			name: "DirectoryNotExists_Error",
			setupMock: func(mock *MockFileOperator) {
				// ディレクトリを設定しない
			},
			dirPath:      "/nonexistent/dir",
			description:  "テストディレクトリ",
			expectError:  true,
			errorMessage: "テストディレクトリが存在しません",
		},
		{
			name: "PathIsFile_Error",
			setupMock: func(mock *MockFileOperator) {
				mock.SetFile("/path/to/file.txt", []byte("content"))
			},
			dirPath:      "/path/to/file.txt",
			description:  "テストディレクトリ",
			expectError:  true,
			errorMessage: "指定されたパスはディレクトリではありません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFileOperator := NewMockFileOperator()
			tt.setupMock(mockFileOperator)

			service := NewServiceWithFileOperator(mockFileOperator)
			err := service.validateDirectory(tt.dirPath, tt.description)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if tt.errorMessage != "" && !strings.Contains(err.Error(), tt.errorMessage) {
					t.Errorf("Expected error message to contain '%s', got '%s'", tt.errorMessage, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}

func TestCreateDestinationDirectory_Normal(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func(*MockFileOperator)
		destDir      string
		expectError  bool
		errorMessage string
	}{
		{
			name: "Success_Normal",
			setupMock: func(mock *MockFileOperator) {
				// 正常なケース
			},
			destDir:     "/output/dir",
			expectError: false,
		},
		{
			name: "MkdirAllFailed_Error",
			setupMock: func(mock *MockFileOperator) {
				mock.SetError("/output/dir", fmt.Errorf("permission denied"))
			},
			destDir:      "/output/dir",
			expectError:  true,
			errorMessage: "permission denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFileOperator := NewMockFileOperator()
			tt.setupMock(mockFileOperator)

			service := NewServiceWithFileOperator(mockFileOperator)
			err := service.createDestinationDirectory(tt.destDir)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if tt.errorMessage != "" && !strings.Contains(err.Error(), tt.errorMessage) {
					t.Errorf("Expected error message to contain '%s', got '%s'", tt.errorMessage, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}

func TestFindContentFiles_Normal(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(*MockFileOperator)
		srcDir        string
		expectedFiles []string
		expectError   bool
		errorMessage  string
	}{
		{
			name: "FilesWithMarker_Normal",
			setupMock: func(mock *MockFileOperator) {
				mock.SetDirectory("/src")
				mock.SetFile("/src/file1.md", []byte(validMarkdownWithMarker))
				mock.SetFile("/src/file2.md", []byte(markdownWithoutMarker))
				mock.SetFile("/src/file3.md", []byte(validMarkdownWithMarker))
				mock.SetFile("/src/file4.txt", []byte("not markdown"))
			},
			srcDir:        "/src",
			expectedFiles: []string{"/src/file1.md", "/src/file3.md"},
			expectError:   false,
		},
		{
			name: "NoMarkerFiles_Normal",
			setupMock: func(mock *MockFileOperator) {
				mock.SetDirectory("/src")
				mock.SetFile("/src/file1.md", []byte(markdownWithoutMarker))
				mock.SetFile("/src/file2.md", []byte(incompleteMarker))
			},
			srcDir:        "/src",
			expectedFiles: []string{},
			expectError:   false,
		},
		{
			name: "WalkDirFailed_Error",
			setupMock: func(mock *MockFileOperator) {
				mock.SetError("/src", fmt.Errorf("walk error"))
			},
			srcDir:       "/src",
			expectError:  true,
			errorMessage: "ディレクトリの走査に失敗しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFileOperator := NewMockFileOperator()
			tt.setupMock(mockFileOperator)

			service := NewServiceWithFileOperator(mockFileOperator)
			files, err := service.findContentFiles(tt.srcDir)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if tt.errorMessage != "" && !strings.Contains(err.Error(), tt.errorMessage) {
					t.Errorf("Expected error message to contain '%s', got '%s'", tt.errorMessage, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if len(files) != len(tt.expectedFiles) {
					t.Errorf("Expected %d files, got %d", len(tt.expectedFiles), len(files))
				}
				// 順序に依存しない検証：期待されるファイルがすべて含まれているかチェック
				expectedMap := make(map[string]bool)
				for _, expectedFile := range tt.expectedFiles {
					expectedMap[expectedFile] = true
				}
				for _, file := range files {
					if !expectedMap[file] {
						t.Errorf("Unexpected file found: %s", file)
					}
					delete(expectedMap, file)
				}
				for missingFile := range expectedMap {
					t.Errorf("Expected file not found: %s", missingFile)
				}
			}
		})
	}
}

func TestExtractBlogContent_Normal(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockFileOperator)
		srcDir         string
		destDir        string
		expectedResult string
		expectError    bool
		errorMessage   string
	}{
		{
			name: "Success_Normal",
			setupMock: func(mock *MockFileOperator) {
				mock.SetDirectory("/src")
				mock.SetDirectory("/dest")
				mock.SetFile("/src/file1.md", []byte(validMarkdownWithMarker))
				mock.SetFile("/src/file2.md", []byte(validMarkdownWithMarker))
			},
			srcDir:         "/src",
			destDir:        "/dest",
			expectedResult: "処理完了: 2件のファイルからコンテンツを抽出しました。",
			expectError:    false,
		},
		{
			name: "NoMarkerFiles_Normal",
			setupMock: func(mock *MockFileOperator) {
				mock.SetDirectory("/src")
				mock.SetDirectory("/dest")
				mock.SetFile("/src/file1.md", []byte(markdownWithoutMarker))
			},
			srcDir:         "/src",
			destDir:        "/dest",
			expectedResult: "指定されたディレクトリにコンテンツマーカーを含むMarkdownファイルが見つかりませんでした。",
			expectError:    false,
		},
		{
			name: "SourceDirNotExists_Error",
			setupMock: func(mock *MockFileOperator) {
				// ソースディレクトリを設定しない
			},
			srcDir:       "/nonexistent",
			destDir:      "/dest",
			expectError:  true,
			errorMessage: "ソースディレクトリが存在しません",
		},
		{
			name: "DestDirCreationFailed_Error",
			setupMock: func(mock *MockFileOperator) {
				mock.SetDirectory("/src")
				mock.SetError("/dest", fmt.Errorf("permission denied"))
			},
			srcDir:       "/src",
			destDir:      "/dest",
			expectError:  true,
			errorMessage: "出力ディレクトリの作成に失敗しました",
		},
		{
			name: "PartialSuccess_WithErrors",
			setupMock: func(mock *MockFileOperator) {
				mock.SetDirectory("/src")
				mock.SetDirectory("/dest")
				mock.SetFile("/src/file1.md", []byte(validMarkdownWithMarker))
				mock.SetFile("/src/file2.md", []byte(validMarkdownWithMarker))
				// file2.mdの書き込みでエラーを発生させる
				mock.SetError("/dest/file2.md", fmt.Errorf("write error"))
			},
			srcDir:         "/src",
			destDir:        "/dest",
			expectedResult: "処理完了: 1件のファイルからコンテンツを抽出しました。\n\nエラーが発生したファイル:\nファイル file2.md: ファイルの保存に失敗しました: write error",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFileOperator := NewMockFileOperator()
			tt.setupMock(mockFileOperator)

			service := NewServiceWithFileOperator(mockFileOperator)
			result, err := service.ExtractBlogContent(tt.srcDir, tt.destDir)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if tt.errorMessage != "" && !strings.Contains(err.Error(), tt.errorMessage) {
					t.Errorf("Expected error message to contain '%s', got '%s'", tt.errorMessage, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if result != tt.expectedResult {
					t.Errorf("Expected result:\n%s\n\nGot:\n%s", tt.expectedResult, result)
				}
			}
		})
	}
}

func TestExtractContentFromFile_Normal(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func(*MockFileOperator)
		filePath     string
		destDir      string
		expectError  bool
		errorMessage string
	}{
		{
			name: "Success_Normal",
			setupMock: func(mock *MockFileOperator) {
				mock.SetFile("/src/test.md", []byte(validMarkdownWithMarker))
			},
			filePath:    "/src/test.md",
			destDir:     "/dest",
			expectError: false,
		},
		{
			name: "FileReadError_Error",
			setupMock: func(mock *MockFileOperator) {
				mock.SetError("/src/test.md", fmt.Errorf("read error"))
			},
			filePath:     "/src/test.md",
			destDir:      "/dest",
			expectError:  true,
			errorMessage: "ファイルの読み込みに失敗しました",
		},
		{
			name: "NoMarker_Error",
			setupMock: func(mock *MockFileOperator) {
				mock.SetFile("/src/test.md", []byte(markdownWithoutMarker))
			},
			filePath:     "/src/test.md",
			destDir:      "/dest",
			expectError:  true,
			errorMessage: "コンテンツの抽出に失敗しました",
		},
		{
			name: "FileWriteError_Error",
			setupMock: func(mock *MockFileOperator) {
				mock.SetFile("/src/test.md", []byte(validMarkdownWithMarker))
				mock.SetError("/dest/test.md", fmt.Errorf("write error"))
			},
			filePath:     "/src/test.md",
			destDir:      "/dest",
			expectError:  true,
			errorMessage: "ファイルの保存に失敗しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFileOperator := NewMockFileOperator()
			tt.setupMock(mockFileOperator)

			service := NewServiceWithFileOperator(mockFileOperator)
			err := service.extractContentFromFile(tt.filePath, tt.destDir)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if tt.errorMessage != "" && !strings.Contains(err.Error(), tt.errorMessage) {
					t.Errorf("Expected error message to contain '%s', got '%s'", tt.errorMessage, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				// 正常な場合は、ファイルが正しく書き込まれていることを確認
				expectedPath := filepath.Join(tt.destDir, filepath.Base(tt.filePath))
				if content, exists := mockFileOperator.files[expectedPath]; !exists {
					t.Error("Expected file to be written")
				} else if string(content) != expectedExtractedContent {
					t.Errorf("Expected content:\n%s\n\nGot:\n%s", expectedExtractedContent, string(content))
				}
			}
		})
	}
}
