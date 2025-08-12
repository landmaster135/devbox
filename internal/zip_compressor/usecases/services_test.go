package usecases

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// MockFileSystemOperator はテスト用のFileSystemOperatorモック
type MockFileSystemOperator struct {
	statFunc     func(name string) (os.FileInfo, error)
	createFunc   func(name string) (*os.File, error)
	openFunc     func(name string) (*os.File, error)
	mkdirAllFunc func(path string, perm os.FileMode) error
	walkFunc     func(root string, walkFn filepath.WalkFunc) error
}

func (m *MockFileSystemOperator) Stat(name string) (os.FileInfo, error) {
	if m.statFunc != nil {
		return m.statFunc(name)
	}
	return nil, errors.New("not implemented")
}

func (m *MockFileSystemOperator) Create(name string) (*os.File, error) {
	if m.createFunc != nil {
		return m.createFunc(name)
	}
	return nil, errors.New("not implemented")
}

func (m *MockFileSystemOperator) Open(name string) (*os.File, error) {
	if m.openFunc != nil {
		return m.openFunc(name)
	}
	return nil, errors.New("not implemented")
}

func (m *MockFileSystemOperator) MkdirAll(path string, perm os.FileMode) error {
	if m.mkdirAllFunc != nil {
		return m.mkdirAllFunc(path, perm)
	}
	return errors.New("not implemented")
}

func (m *MockFileSystemOperator) Walk(root string, walkFn filepath.WalkFunc) error {
	if m.walkFunc != nil {
		return m.walkFunc(root, walkFn)
	}
	return errors.New("not implemented")
}

// MockFileInfo はテスト用のos.FileInfoモック
type MockFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	isDir   bool
	modTime string
}

func (m *MockFileInfo) Name() string      { return m.name }
func (m *MockFileInfo) Size() int64       { return m.size }
func (m *MockFileInfo) Mode() os.FileMode { return m.mode }
func (m *MockFileInfo) IsDir() bool       { return m.isDir }
func (m *MockFileInfo) ModTime() time.Time  { return time.Time{} }
func (m *MockFileInfo) Sys() interface{}  { return nil }

// MockFile はテスト用のos.Fileモック
type MockFile struct {
	content  string
	position int
	fileInfo os.FileInfo
	closed   bool
}

func NewMockFile(content string, fileInfo os.FileInfo) *MockFile {
	return &MockFile{
		content:  content,
		position: 0,
		fileInfo: fileInfo,
		closed:   false,
	}
}

func (m *MockFile) Read(p []byte) (n int, err error) {
	if m.closed {
		return 0, errors.New("file is closed")
	}
	if m.position >= len(m.content) {
		return 0, io.EOF
	}
	n = copy(p, m.content[m.position:])
	m.position += n
	return n, nil
}

func (m *MockFile) Write(p []byte) (n int, err error) {
	if m.closed {
		return 0, errors.New("file is closed")
	}
	return len(p), nil
}

func (m *MockFile) Close() error {
	m.closed = true
	return nil
}

func (m *MockFile) Stat() (os.FileInfo, error) {
	return m.fileInfo, nil
}

func TestNewZipCompressorService_Normal(t *testing.T) {
	service := NewZipCompressorService()
	if service == nil {
		t.Error("NewZipCompressorService() returned nil")
		return
	}
	if service.fileSystem == nil {
		t.Error("NewZipCompressorService() fileSystem is nil")
	}
}

func TestNewZipCompressorServiceWithFileSystem_Normal(t *testing.T) {
	mockFS := &MockFileSystemOperator{}
	service := NewZipCompressorServiceWithFileSystem(mockFS)
	if service == nil {
		t.Error("NewZipCompressorServiceWithFileSystem() returned nil")
		return
	}
	if service.fileSystem != mockFS {
		t.Error("NewZipCompressorServiceWithFileSystem() fileSystem not set correctly")
	}
}

func TestHandleCompress_FileNotExists(t *testing.T) {
	mockFS := &MockFileSystemOperator{
		statFunc: func(name string) (os.FileInfo, error) {
			return nil, errors.New("file not found")
		},
	}

	service := NewZipCompressorServiceWithFileSystem(mockFS)
	result, err := service.HandleCompress("/nonexistent/file")

	if err == nil {
		t.Error("HandleCompress() expected error but got nil")
	}
	if result != "" {
		t.Errorf("HandleCompress() result = %v, want empty string", result)
	}
	if !strings.Contains(err.Error(), "指定されたパスが存在しません") {
		t.Errorf("HandleCompress() error = %v, want contains '指定されたパスが存在しません'", err.Error())
	}
}

func TestHandleCompress_CreateZipFileFailed(t *testing.T) {
	mockFileInfo := &MockFileInfo{
		name:  "test.txt",
		size:  100,
		mode:  0644,
		isDir: false,
	}

	mockFS := &MockFileSystemOperator{
		statFunc: func(name string) (os.FileInfo, error) {
			return mockFileInfo, nil
		},
		createFunc: func(name string) (*os.File, error) {
			return nil, errors.New("create failed")
		},
	}

	service := NewZipCompressorServiceWithFileSystem(mockFS)
	result, err := service.HandleCompress("/path/to/test.txt")

	if err == nil {
		t.Error("HandleCompress() expected error but got nil")
	}
	if result != "" {
		t.Errorf("HandleCompress() result = %v, want empty string", result)
	}
	if !strings.Contains(err.Error(), "zipファイルの作成に失敗しました") {
		t.Errorf("HandleCompress() error = %v, want contains 'zipファイルの作成に失敗しました'", err.Error())
	}
}

func TestHandleDecompress_FileNotExists(t *testing.T) {
	mockFS := &MockFileSystemOperator{
		statFunc: func(name string) (os.FileInfo, error) {
			return nil, errors.New("file not found")
		},
	}

	service := NewZipCompressorServiceWithFileSystem(mockFS)
	result, err := service.HandleDecompress("/nonexistent/file.zip")

	if err == nil {
		t.Error("HandleDecompress() expected error but got nil")
	}
	if result != "" {
		t.Errorf("HandleDecompress() result = %v, want empty string", result)
	}
	if !strings.Contains(err.Error(), "指定されたZipファイルが存在しません") {
		t.Errorf("HandleDecompress() error = %v, want contains '指定されたZipファイルが存在しません'", err.Error())
	}
}

func TestHandleDecompress_PathIsDirectory(t *testing.T) {
	mockFileInfo := &MockFileInfo{
		name:  "directory",
		size:  0,
		mode:  0755,
		isDir: true,
	}

	mockFS := &MockFileSystemOperator{
		statFunc: func(name string) (os.FileInfo, error) {
			return mockFileInfo, nil
		},
	}

	service := NewZipCompressorServiceWithFileSystem(mockFS)
	result, err := service.HandleDecompress("/path/to/directory")

	if err == nil {
		t.Error("HandleDecompress() expected error but got nil")
	}
	if result != "" {
		t.Errorf("HandleDecompress() result = %v, want empty string", result)
	}
	if !strings.Contains(err.Error(), "指定されたパスはディレクトリです") {
		t.Errorf("HandleDecompress() error = %v, want contains '指定されたパスはディレクトリです'", err.Error())
	}
}

func TestHandleDecompress_NotZipFile(t *testing.T) {
	mockFileInfo := &MockFileInfo{
		name:  "test.txt",
		size:  100,
		mode:  0644,
		isDir: false,
	}

	mockFS := &MockFileSystemOperator{
		statFunc: func(name string) (os.FileInfo, error) {
			return mockFileInfo, nil
		},
	}

	service := NewZipCompressorServiceWithFileSystem(mockFS)
	result, err := service.HandleDecompress("/path/to/test.txt")

	if err == nil {
		t.Error("HandleDecompress() expected error but got nil")
	}
	if result != "" {
		t.Errorf("HandleDecompress() result = %v, want empty string", result)
	}
	if !strings.Contains(err.Error(), "指定されたファイルはZipファイルではありません") {
		t.Errorf("HandleDecompress() error = %v, want contains '指定されたファイルはZipファイルではありません'", err.Error())
	}
}

func TestHandleDecompress_MkdirAllFailed(t *testing.T) {
	mockFileInfo := &MockFileInfo{
		name:  "test.zip",
		size:  100,
		mode:  0644,
		isDir: false,
	}

	mockFS := &MockFileSystemOperator{
		statFunc: func(name string) (os.FileInfo, error) {
			return mockFileInfo, nil
		},
		mkdirAllFunc: func(path string, perm os.FileMode) error {
			return errors.New("mkdir failed")
		},
	}

	service := NewZipCompressorServiceWithFileSystem(mockFS)
	result, err := service.HandleDecompress("/path/to/test.zip")

	if err == nil {
		t.Error("HandleDecompress() expected error but got nil")
	}
	if result != "" {
		t.Errorf("HandleDecompress() result = %v, want empty string", result)
	}
	if !strings.Contains(err.Error(), "出力ディレクトリの作成に失敗しました") {
		t.Errorf("HandleDecompress() error = %v, want contains '出力ディレクトリの作成に失敗しました'", err.Error())
	}
}

func TestExtractFile_PathTraversalAttack(t *testing.T) {
	service := NewZipCompressorService()

	// zip.Fileを模擬（実際のテストでは複雑になるため、エラーケースのみテスト）
	mockZipFile := &zip.File{
		FileHeader: zip.FileHeader{
			Name: "../../../etc/passwd",
		},
	}

	// FileInfoメソッドをオーバーライド
	mockZipFile.FileHeader.SetMode(0644)

	err := service.extractFile(mockZipFile, "/safe/output/dir")

	if err == nil {
		t.Error("extractFile() expected error for path traversal attack but got nil")
	}
	if !strings.Contains(err.Error(), "不正なパスが検出されました") {
		t.Errorf("extractFile() error = %v, want contains '不正なパスが検出されました'", err.Error())
	}
}

func TestCompressFile_OpenFileFailed(t *testing.T) {
	mockFS := &MockFileSystemOperator{
		openFunc: func(name string) (*os.File, error) {
			return nil, errors.New("open failed")
		},
	}

	service := NewZipCompressorServiceWithFileSystem(mockFS)

	// 実際のzip.Writerを作成（テスト用の簡易実装）
	// この部分は実際のテストでは、より詳細なモックが必要
	err := service.compressFile("/path/to/file", nil)

	if err == nil {
		t.Error("compressFile() expected error but got nil")
	}
	if !strings.Contains(err.Error(), "ファイルを開けませんでした") {
		t.Errorf("compressFile() error = %v, want contains 'ファイルを開けませんでした'", err.Error())
	}
}
