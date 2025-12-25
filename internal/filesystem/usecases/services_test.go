package usecases

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// #==============================================================#
// ##          Mock Implementations                              ##
// #==============================================================#

// MockFileOpener はテスト用のFileOpenerモックです
type MockFileOpener struct {
	OpenFunc func(name string) (*os.File, error)
}

func (m *MockFileOpener) Open(name string) (*os.File, error) {
	if m.OpenFunc != nil {
		return m.OpenFunc(name)
	}
	return nil, errors.New("mock error")
}

// MockFileWriter はテスト用のFileWriterモックです
type MockFileWriter struct {
	WriteFileFunc func(filename string, data []byte, perm os.FileMode) error
}

func (m *MockFileWriter) WriteFile(filename string, data []byte, perm os.FileMode) error {
	if m.WriteFileFunc != nil {
		return m.WriteFileFunc(filename, data, perm)
	}
	return errors.New("mock error")
}

// MockDirectoryReader はテスト用のDirectoryReaderモックです
type MockDirectoryReader struct {
	ReadDirFunc func(name string) ([]os.DirEntry, error)
}

func (m *MockDirectoryReader) ReadDir(name string) ([]os.DirEntry, error) {
	if m.ReadDirFunc != nil {
		return m.ReadDirFunc(name)
	}
	return nil, errors.New("mock error")
}

// MockFileStat はテスト用のFileStatモックです
type MockFileStat struct {
	StatFunc func(name string) (os.FileInfo, error)
}

func (m *MockFileStat) Stat(name string) (os.FileInfo, error) {
	if m.StatFunc != nil {
		return m.StatFunc(name)
	}
	return nil, errors.New("mock error")
}

// MockJSONMarshaler はテスト用のJSONMarshalerモックです
type MockJSONMarshaler struct {
	MarshalIndentFunc func(v interface{}, prefix, indent string) ([]byte, error)
}

func (m *MockJSONMarshaler) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	if m.MarshalIndentFunc != nil {
		return m.MarshalIndentFunc(v, prefix, indent)
	}
	return nil, errors.New("mock error")
}

// MockYAMLMarshaler はテスト用のYAMLMarshalerモックです
type MockYAMLMarshaler struct {
	MarshalFunc func(v interface{}) ([]byte, error)
}

func (m *MockYAMLMarshaler) Marshal(v interface{}) ([]byte, error) {
	if m.MarshalFunc != nil {
		return m.MarshalFunc(v)
	}
	return nil, errors.New("mock error")
}

// MockFileInfo はテスト用のos.FileInfoモックです
type MockFileInfo struct {
	NameFunc    func() string
	SizeFunc    func() int64
	ModeFunc    func() os.FileMode
	ModTimeFunc func() time.Time
	IsDirFunc   func() bool
	SysFunc     func() interface{}
}

func (m *MockFileInfo) Name() string {
	if m.NameFunc != nil {
		return m.NameFunc()
	}
	return "test.txt"
}

func (m *MockFileInfo) Size() int64 {
	if m.SizeFunc != nil {
		return m.SizeFunc()
	}
	return 100
}

func (m *MockFileInfo) Mode() os.FileMode {
	if m.ModeFunc != nil {
		return m.ModeFunc()
	}
	return 0644
}

func (m *MockFileInfo) ModTime() time.Time {
	if m.ModTimeFunc != nil {
		return m.ModTimeFunc()
	}
	return time.Now()
}

func (m *MockFileInfo) IsDir() bool {
	if m.IsDirFunc != nil {
		return m.IsDirFunc()
	}
	return false
}

func (m *MockFileInfo) Sys() interface{} {
	if m.SysFunc != nil {
		return m.SysFunc()
	}
	return nil
}

// MockDirEntry はテスト用のos.DirEntryモックです
type MockDirEntry struct {
	NameFunc  func() string
	IsDirFunc func() bool
	TypeFunc  func() os.FileMode
	InfoFunc  func() (os.FileInfo, error)
}

func (m *MockDirEntry) Name() string {
	if m.NameFunc != nil {
		return m.NameFunc()
	}
	return "test.txt"
}

func (m *MockDirEntry) IsDir() bool {
	if m.IsDirFunc != nil {
		return m.IsDirFunc()
	}
	return false
}

func (m *MockDirEntry) Type() os.FileMode {
	if m.TypeFunc != nil {
		return m.TypeFunc()
	}
	return 0644
}

func (m *MockDirEntry) Info() (os.FileInfo, error) {
	if m.InfoFunc != nil {
		return m.InfoFunc()
	}
	return &MockFileInfo{}, nil
}

// #==============================================================#
// ##          Test Cases                                        ##
// #==============================================================#

// TestFileSystemService_NewFileSystemService_Normal は正常系のテストです
func TestFileSystemService_NewFileSystemService_Normal(t *testing.T) {
	// Arrange
	allowedDirs := []string{"/tmp", "/home/user"}

	// Act
	service := NewFileSystemService(allowedDirs)

	// Assert
	if service == nil {
		t.Fatal("サービスがnilです")
	}
	if len(service.allowedDirectories) != 2 {
		t.Errorf("許可されたディレクトリの数が期待値と異なります。期待値: 2, 実際: %d", len(service.allowedDirectories))
	}
}

// TestFileSystemService_ValidatePath_Normal は正常系のテストです
func TestFileSystemService_ValidatePath_Normal(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	allowedDirs := []string{tempDir}
	service := NewFileSystemService(allowedDirs)
	testPath := filepath.Join(tempDir, "test.txt")

	// Act
	validPath, err := service.ValidatePath(testPath)

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
	if !strings.Contains(validPath, "test.txt") {
		t.Errorf("パスが期待値と異なります。実際: %s", validPath)
	}
}

// TestFileSystemService_ValidatePath_AccessDenied はアクセス拒否のテストです
func TestFileSystemService_ValidatePath_AccessDenied(t *testing.T) {
	// Arrange
	allowedDirs := []string{"/tmp"}
	service := NewFileSystemService(allowedDirs)
	testPath := "/etc/passwd"

	// Act
	_, err := service.ValidatePath(testPath)

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、発生しませんでした")
	}
	if !strings.Contains(err.Error(), "アクセス拒否") {
		t.Errorf("期待されたエラーメッセージが含まれていません。実際: %s", err.Error())
	}
}

// TestFileSystemService_ReadFile_Normal は正常系のテストです
func TestFileSystemService_ReadFile_Normal(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "test content"

	// テストファイルを作成
	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}

	allowedDirs := []string{tempDir}
	service := NewFileSystemService(allowedDirs)

	// Act
	content, err := service.ReadFile(testFile, 1, 2000)

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
	if content != "L1: "+testContent {
		t.Errorf("内容が期待値と異なります。期待値: %s, 実際: %s", "L1: "+testContent, content)
	}
}

func TestFileSystemService_ReadFile_AddsLineNumbersAndHandlesCRLF(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	content := "foo\r\nbar\r\n\r\n"

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}

	service := NewFileSystemService([]string{tempDir})
	got, err := service.ReadFile(testFile, 1, 2000)
	if err != nil {
		t.Fatalf("read_fileでエラー: %v", err)
	}

	expected := "L1: foo\nL2: bar\nL3: "
	if got != expected {
		t.Errorf("内容が期待値と異なります。期待値: %s, 実際: %s", expected, got)
	}
}

func TestFileSystemService_ReadFile_TruncatesLongLines(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	line := strings.Repeat("あ", 600)

	if err := os.WriteFile(testFile, []byte(line), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}

	service := NewFileSystemService([]string{tempDir})
	got, err := service.ReadFile(testFile, 1, 2000)
	if err != nil {
		t.Fatalf("read_fileでエラー: %v", err)
	}

	expected := "L1: " + strings.Repeat("あ", readFileLineMaxLength)
	if got != expected {
		t.Errorf("500文字制限が期待通りに働いていません。期待値: %d文字, 実際: %d文字", len(expected), len(got))
	}
}

func TestFileSystemService_ReadFile_RespectsOffsetAndLimit(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	content := "line1\nline2\nline3\n"

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}

	service := NewFileSystemService([]string{tempDir})
	got, err := service.ReadFile(testFile, 2, 1)
	if err != nil {
		t.Fatalf("read_fileでエラー: %v", err)
	}
	if got != "L2: line2" {
		t.Errorf("offsetの結果が期待値と異なります。実際: %s", got)
	}

	got, err = service.ReadFile(testFile, 3, 5)
	if err != nil {
		t.Fatalf("read_fileでエラー: %v", err)
	}
	if got != "L3: line3" {
		t.Errorf("limitが末尾で正しく丸められていません。実際: %s", got)
	}
}

func TestFileSystemService_ReadFile_InvalidOffsetOrLimit(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("only-one-line"), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}

	service := NewFileSystemService([]string{tempDir})

	if _, err := service.ReadFile(testFile, 0, 10); err == nil {
		t.Fatal("offset=0でエラーが発生しませんでした")
	}
	if _, err := service.ReadFile(testFile, 1, 0); err == nil {
		t.Fatal("limit=0でエラーが発生しませんでした")
	}
	if _, err := service.ReadFile(testFile, 5, 10); err == nil {
		t.Fatal("存在しない行を指定してもエラーになりませんでした")
	}
}

// TestFileSystemService_WriteFile_Normal は正常系のテストです
func TestFileSystemService_WriteFile_Normal(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "test content"
	allowedDirs := []string{tempDir}
	service := NewFileSystemService(allowedDirs)

	// Act
	err := service.WriteFile(testFile, testContent)

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}

	// ファイルが作成されたことを確認
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Errorf("ファイルの読み取りに失敗しました: %v", err)
	}
	if string(content) != testContent {
		t.Errorf("内容が期待値と異なります。期待値: %s, 実際: %s", testContent, string(content))
	}
}

// TestFileSystemService_CreateDirectory_Normal は正常系のテストです
func TestFileSystemService_CreateDirectory_Normal(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	testDir := filepath.Join(tempDir, "testdir")
	allowedDirs := []string{tempDir}
	service := NewFileSystemService(allowedDirs)

	// Act
	err := service.CreateDirectory(testDir)

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}

	// ディレクトリが作成されたことを確認
	info, err := os.Stat(testDir)
	if err != nil {
		t.Errorf("ディレクトリの確認に失敗しました: %v", err)
	}
	if !info.IsDir() {
		t.Error("ディレクトリが作成されていません")
	}
}

// TestFileSystemService_ListDirectory_Normal は正常系のテストです
func TestFileSystemService_ListDirectory_Normal(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()

	// テストファイルとディレクトリを作成
	testFile := filepath.Join(tempDir, "test.txt")
	testSubDir := filepath.Join(tempDir, "subdir")

	err := os.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}

	err = os.Mkdir(testSubDir, 0755)
	if err != nil {
		t.Fatalf("テストディレクトリの作成に失敗しました: %v", err)
	}

	allowedDirs := []string{tempDir}
	service := NewFileSystemService(allowedDirs)

	// Act
	entries, err := service.ListDirectory(tempDir)

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("エントリ数が期待値と異なります。期待値: 2, 実際: %d", len(entries))
	}

	// ファイルとディレクトリが正しく識別されているか確認
	hasFile := false
	hasDir := false
	for _, entry := range entries {
		if strings.Contains(entry, "[FILE] test.txt") {
			hasFile = true
		}
		if strings.Contains(entry, "[DIR] subdir") {
			hasDir = true
		}
	}
	if !hasFile {
		t.Error("ファイルが正しく識別されていません")
	}
	if !hasDir {
		t.Error("ディレクトリが正しく識別されていません")
	}
}

// TestFileSystemService_GetFileInfo_Normal は正常系のテストです
func TestFileSystemService_GetFileInfo_Normal(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "test content"

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}

	allowedDirs := []string{tempDir}
	service := NewFileSystemService(allowedDirs)

	// Act
	info, err := service.GetFileInfo(testFile)

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
	if info == nil {
		t.Fatal("ファイル情報がnilです")
	}
	if info.Size != int64(len(testContent)) {
		t.Errorf("ファイルサイズが期待値と異なります。期待値: %d, 実際: %d", len(testContent), info.Size)
	}
	if info.IsDirectory {
		t.Error("ファイルがディレクトリとして識別されています")
	}
	if !info.IsFile {
		t.Error("ファイルがファイルとして識別されていません")
	}
}

// TestFileSystemService_SearchFiles_Normal は正常系のテストです
func TestFileSystemService_SearchFiles_Normal(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()

	// テストファイルを作成
	testFile1 := filepath.Join(tempDir, "test1.txt")
	testFile2 := filepath.Join(tempDir, "test2.go")
	testFile3 := filepath.Join(tempDir, "other.txt")

	for _, file := range []string{testFile1, testFile2, testFile3} {
		err := os.WriteFile(file, []byte("test"), 0644)
		if err != nil {
			t.Fatalf("テストファイルの作成に失敗しました: %v", err)
		}
	}

	allowedDirs := []string{tempDir}
	service := NewFileSystemService(allowedDirs)

	// Act
	results, err := service.SearchFiles(tempDir, "test", []string{})

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("検索結果数が期待値と異なります。期待値: 2, 実際: %d", len(results))
	}
}

// TestFileSystemService_GetAllowedDirectories_Normal は正常系のテストです
func TestFileSystemService_GetAllowedDirectories_Normal(t *testing.T) {
	// Arrange
	allowedDirs := []string{"/tmp", "/home/user"}
	service := NewFileSystemService(allowedDirs)

	// Act
	dirs := service.GetAllowedDirectories()

	// Assert
	if len(dirs) != 2 {
		t.Errorf("許可されたディレクトリ数が期待値と異なります。期待値: 2, 実際: %d", len(dirs))
	}
}

// TestFileSystemService_GetAllowedDirectoriesAsText_Normal は正常系のテストです
func TestFileSystemService_GetAllowedDirectoriesAsText_Normal(t *testing.T) {
	// Arrange
	allowedDirs := []string{"/tmp", "/home/user"}
	service := NewFileSystemService(allowedDirs)

	// Act
	text := service.GetAllowedDirectoriesAsText()

	// Assert
	if !strings.Contains(text, "許可されたディレクトリ:") {
		t.Error("期待されたテキストが含まれていません")
	}
}

// TestExpandHome_Normal は正常系のテストです
func TestExpandHome_Normal(t *testing.T) {
	// Arrange
	testCases := []struct {
		input    string
		expected string
	}{
		{"/absolute/path", "/absolute/path"},
		{"relative/path", "relative/path"},
	}

	// Act & Assert
	for _, tc := range testCases {
		result := expandHome(tc.input)
		if result != tc.expected && !strings.HasPrefix(tc.input, "~") {
			t.Errorf("expandHome(%s) = %s, 期待値: %s", tc.input, result, tc.expected)
		}
	}
}

// TestFileSystemService_WithMockDependencies_Normal はモック依存性を使用したテストです
func TestFileSystemService_WithMockDependencies_Normal(t *testing.T) {
	// Arrange
	allowedDirs := []string{"/tmp"}

	mockFileOpener := &MockFileOpener{
		OpenFunc: func(name string) (*os.File, error) {
			return nil, nil // 実際のファイルは開かない
		},
	}

	mockFileWriter := &MockFileWriter{
		WriteFileFunc: func(filename string, data []byte, perm os.FileMode) error {
			return nil // 成功を返す
		},
	}

	mockDirectoryReader := &MockDirectoryReader{
		ReadDirFunc: func(name string) ([]os.DirEntry, error) {
			return []os.DirEntry{
				&MockDirEntry{
					NameFunc:  func() string { return "test.txt" },
					IsDirFunc: func() bool { return false },
				},
			}, nil
		},
	}

	mockFileStat := &MockFileStat{
		StatFunc: func(name string) (os.FileInfo, error) {
			return &MockFileInfo{
				SizeFunc:  func() int64 { return 100 },
				IsDirFunc: func() bool { return false },
			}, nil
		},
	}

	mockJSONMarshaler := &MockJSONMarshaler{
		MarshalIndentFunc: func(v interface{}, prefix, indent string) ([]byte, error) {
			return []byte(`[{"name":"test.txt","type":"file"}]`), nil
		},
	}

	mockYAMLMarshaler := &MockYAMLMarshaler{
		MarshalFunc: func(v interface{}) ([]byte, error) {
			return []byte("- name: test.txt\n  type: file\n"), nil
		},
	}

	service := NewFileSystemServiceWithDependencies(
		allowedDirs,
		mockFileOpener,
		mockFileWriter,
		mockDirectoryReader,
		mockFileStat,
		mockJSONMarshaler,
		mockYAMLMarshaler,
	)

	// Act & Assert
	if service == nil {
		t.Fatal("サービスがnilです")
	}

	// WriteFileのテスト
	err := service.WriteFile("/tmp/test.txt", "test content")
	if err != nil {
		t.Errorf("WriteFileでエラーが発生しました: %v", err)
	}

	// ListDirectoryのテスト
	entries, err := service.ListDirectory("/tmp")
	if err != nil {
		t.Errorf("ListDirectoryでエラーが発生しました: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("エントリ数が期待値と異なります。期待値: 1, 実際: %d", len(entries))
	}

	// GetFileInfoのテスト
	info, err := service.GetFileInfo("/tmp/test.txt")
	if err != nil {
		t.Errorf("GetFileInfoでエラーが発生しました: %v", err)
	}
	if info.Size != 100 {
		t.Errorf("ファイルサイズが期待値と異なります。期待値: 100, 実際: %d", info.Size)
	}

	// GetDirectoryTreeAsJSONのテスト
	jsonStr, err := service.GetDirectoryTreeAsJSON("/tmp")
	if err != nil {
		t.Errorf("GetDirectoryTreeAsJSONでエラーが発生しました: %v", err)
	}
	if !strings.Contains(jsonStr, "test.txt") {
		t.Error("JSONにtest.txtが含まれていません")
	}

	// GetDirectoryTreeAsYAMLのテスト
	yamlStr, err := service.GetDirectoryTreeAsYAML("/tmp")
	if err != nil {
		t.Errorf("GetDirectoryTreeAsYAMLでエラーが発生しました: %v", err)
	}
	if !strings.Contains(yamlStr, "test.txt") {
		t.Error("YAMLにtest.txtが含まれていません")
	}
}
func TestGetDirectoryTreeWithOptionsDepthAndPagination(t *testing.T) {
	tempDir := t.TempDir()

	alpha := filepath.Join(tempDir, "alpha", "nested")
	if err := os.MkdirAll(alpha, 0o755); err != nil {
		t.Fatalf("nestedディレクトリの作成に失敗しました: %v", err)
	}
	if err := os.WriteFile(filepath.Join(filepath.Dir(alpha), "child.txt"), []byte("child"), 0o644); err != nil {
		t.Fatalf("child.txtの作成に失敗しました: %v", err)
	}

	beta := filepath.Join(tempDir, "beta")
	if err := os.MkdirAll(beta, 0o755); err != nil {
		t.Fatalf("betaディレクトリの作成に失敗しました: %v", err)
	}
	if err := os.WriteFile(filepath.Join(beta, "beta.txt"), []byte("beta"), 0o644); err != nil {
		t.Fatalf("beta.txtの作成に失敗しました: %v", err)
	}

	rootFile := filepath.Join(tempDir, "root.txt")
	if err := os.WriteFile(rootFile, []byte("root"), 0o644); err != nil {
		t.Fatalf("root.txtの作成に失敗しました: %v", err)
	}

	service := NewFileSystemService([]string{tempDir})

	shallow, err := service.GetDirectoryTreeWithOptions(tempDir, DirectoryTreeOptions{Offset: 1, Limit: 0, Depth: 1})
	if err != nil {
		t.Fatalf("深さ1でのツリー取得に失敗しました: %v", err)
	}
	if len(shallow) != 3 {
		t.Fatalf("期待するエントリ数3に対して%v件でした", len(shallow))
	}
	for _, entry := range shallow {
		if entry.Name == "alpha" {
			if len(entry.Children) != 0 {
				t.Fatalf("depth=1では子要素が含まれないはずです")
			}
			if entry.TruncatedChildCount != 2 {
				t.Fatalf("alphaには2件の子があるはずですが%v件でした", entry.TruncatedChildCount)
			}
		}
		if entry.Name == "beta" && entry.TruncatedChildCount != 1 {
			t.Fatalf("betaには1件の子があるはずですが%v件でした", entry.TruncatedChildCount)
		}
	}

	deeper, err := service.GetDirectoryTreeWithOptions(tempDir, DirectoryTreeOptions{Offset: 1, Limit: 0, Depth: 2})
	if err != nil {
		t.Fatalf("深さ2でのツリー取得に失敗しました: %v", err)
	}
	var alphaHasChildren bool
	for _, entry := range deeper {
		if entry.Name == "alpha" {
			if len(entry.Children) > 0 {
				alphaHasChildren = true
			}
			if entry.TruncatedChildCount != 0 {
				t.Fatalf("depth=2でalphaを展開した場合はTruncatedChildCountは0のはずです")
			}
		}
	}
	if !alphaHasChildren {
		t.Fatalf("depth=2ではalphaの子要素を含む必要があります")
	}

	paged, err := service.GetDirectoryTreeWithOptions(tempDir, DirectoryTreeOptions{Offset: 2, Limit: 1, Depth: 0})
	if err != nil {
		t.Fatalf("ページングされたツリー取得に失敗しました: %v", err)
	}
	if len(paged) != 1 {
		t.Fatalf("limit=1なので1件のみ返るはずですが%v件返りました", len(paged))
	}
	if paged[0].Name != "beta" {
		t.Fatalf("offset=2の場合はbetaが返るはずですが%qでした", paged[0].Name)
	}
}

func TestGetDirectoryTreeWithOptionsOffsetError(t *testing.T) {
	tempDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tempDir, "only.txt"), []byte("data"), 0o644); err != nil {
		t.Fatalf("only.txtの作成に失敗しました: %v", err)
	}

	service := NewFileSystemService([]string{tempDir})
	_, err := service.GetDirectoryTreeWithOptions(tempDir, DirectoryTreeOptions{Offset: 5, Limit: 1, Depth: 1})
	if err == nil {
		t.Fatalf("存在しないoffsetに対してエラーが発生しませんでした")
	}
	if !strings.Contains(err.Error(), "offset exceeds directory entry count") {
		t.Fatalf("予期しないエラーメッセージ: %v", err)
	}
}
