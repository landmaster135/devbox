package repositories

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/landmaster135/devbox/internal/file_character_replacer/domain"
)

// MockEncodingConverter はテスト用のEncodingConverterモック
type MockEncodingConverter struct {
	convertToUTF8Func   func(content []byte, encoding domain.EncodingType) (string, error)
	convertFromUTF8Func func(content string, encoding domain.EncodingType) ([]byte, error)
	detectEncodingFunc  func(content []byte) (domain.EncodingType, error)
}

// ConvertToUTF8 は指定されたエンコーディングのバイト列をUTF-8文字列に変換します
func (m *MockEncodingConverter) ConvertToUTF8(content []byte, encoding domain.EncodingType) (string, error) {
	if m.convertToUTF8Func != nil {
		return m.convertToUTF8Func(content, encoding)
	}
	return string(content), nil
}

// ConvertFromUTF8 はUTF-8文字列を指定されたエンコーディングのバイト列に変換します
func (m *MockEncodingConverter) ConvertFromUTF8(content string, encoding domain.EncodingType) ([]byte, error) {
	if m.convertFromUTF8Func != nil {
		return m.convertFromUTF8Func(content, encoding)
	}
	return []byte(content), nil
}

// DetectEncoding はバイト列から文字エンコーディングを推測します
func (m *MockEncodingConverter) DetectEncoding(content []byte) (domain.EncodingType, error) {
	if m.detectEncodingFunc != nil {
		return m.detectEncodingFunc(content)
	}
	return domain.EncodingUTF8, nil
}

// TestNewFileRepository_Normal はNewFileRepository()の正常系をテストします
func TestNewFileRepository_Normal(t *testing.T) {
	mockConverter := &MockEncodingConverter{}
	repo := NewFileRepository(mockConverter)

	if repo == nil {
		t.Error("NewFileRepository() should not return nil")
		return
	}

	// 型が正しいかチェック
	if _, ok := repo.(*FileRepositoryImpl); !ok {
		t.Error("NewFileRepository() should return *FileRepositoryImpl")
	}
}

// TestFileRepositoryImpl_ReadFile_Normal はReadFile()の正常系をテストします
func TestFileRepositoryImpl_ReadFile_Normal(t *testing.T) {
	// テスト用の一時ファイルを作成
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, World!"

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	mockConverter := &MockEncodingConverter{}
	repo := NewFileRepository(mockConverter)

	content, err := repo.ReadFile(testFile, domain.EncodingUTF8)
	if err != nil {
		t.Errorf("ReadFile() returned unexpected error: %v", err)
	}

	if content != testContent {
		t.Errorf("ReadFile() content = %s, expected %s", content, testContent)
	}
}

// TestFileRepositoryImpl_ReadFile_NonUTF8 はReadFile()の非UTF-8エンコーディングをテストします
func TestFileRepositoryImpl_ReadFile_NonUTF8(t *testing.T) {
	// テスト用の一時ファイルを作成
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, World!"

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	mockConverter := &MockEncodingConverter{
		convertToUTF8Func: func(content []byte, encoding domain.EncodingType) (string, error) {
			if encoding == domain.EncodingShiftJIS {
				return "Converted: " + string(content), nil
			}
			return string(content), nil
		},
	}
	repo := NewFileRepository(mockConverter)

	content, err := repo.ReadFile(testFile, domain.EncodingShiftJIS)
	if err != nil {
		t.Errorf("ReadFile() returned unexpected error: %v", err)
	}

	expectedContent := "Converted: " + testContent
	if content != expectedContent {
		t.Errorf("ReadFile() content = %s, expected %s", content, expectedContent)
	}
}

// TestFileRepositoryImpl_ReadFile_FileNotExists はReadFile()のファイル存在しないケースをテストします
func TestFileRepositoryImpl_ReadFile_FileNotExists(t *testing.T) {
	mockConverter := &MockEncodingConverter{}
	repo := NewFileRepository(mockConverter)

	_, err := repo.ReadFile("/nonexistent/file.txt", domain.EncodingUTF8)
	if err == nil {
		t.Error("ReadFile() should return error for non-existent file")
	}
}

// TestFileRepositoryImpl_ReadFile_ConversionError はReadFile()の変換エラーをテストします
func TestFileRepositoryImpl_ReadFile_ConversionError(t *testing.T) {
	// テスト用の一時ファイルを作成
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")

	err := os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	mockConverter := &MockEncodingConverter{
		convertToUTF8Func: func(content []byte, encoding domain.EncodingType) (string, error) {
			return "", fmt.Errorf("conversion error")
		},
	}
	repo := NewFileRepository(mockConverter)

	_, err = repo.ReadFile(testFile, domain.EncodingShiftJIS)
	if err == nil {
		t.Error("ReadFile() should return error when conversion fails")
	}
}

// TestFileRepositoryImpl_WriteFile_Normal はWriteFile()の正常系をテストします
func TestFileRepositoryImpl_WriteFile_Normal(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, World!"

	mockConverter := &MockEncodingConverter{}
	repo := NewFileRepository(mockConverter)

	err := repo.WriteFile(testFile, testContent, domain.EncodingUTF8)
	if err != nil {
		t.Errorf("WriteFile() returned unexpected error: %v", err)
	}

	// ファイルが正しく書き込まれたかチェック
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read written file: %v", err)
	}

	if string(content) != testContent {
		t.Errorf("Written content = %s, expected %s", string(content), testContent)
	}
}

// TestFileRepositoryImpl_WriteFile_NonUTF8 はWriteFile()の非UTF-8エンコーディングをテストします
func TestFileRepositoryImpl_WriteFile_NonUTF8(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, World!"

	mockConverter := &MockEncodingConverter{
		convertFromUTF8Func: func(content string, encoding domain.EncodingType) ([]byte, error) {
			if encoding == domain.EncodingShiftJIS {
				return []byte("Converted: " + content), nil
			}
			return []byte(content), nil
		},
	}
	repo := NewFileRepository(mockConverter)

	err := repo.WriteFile(testFile, testContent, domain.EncodingShiftJIS)
	if err != nil {
		t.Errorf("WriteFile() returned unexpected error: %v", err)
	}

	// ファイルが正しく書き込まれたかチェック
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read written file: %v", err)
	}

	expectedContent := "Converted: " + testContent
	if string(content) != expectedContent {
		t.Errorf("Written content = %s, expected %s", string(content), expectedContent)
	}
}

// TestFileRepositoryImpl_WriteFile_ConversionError はWriteFile()の変換エラーをテストします
func TestFileRepositoryImpl_WriteFile_ConversionError(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")

	mockConverter := &MockEncodingConverter{
		convertFromUTF8Func: func(content string, encoding domain.EncodingType) ([]byte, error) {
			return nil, fmt.Errorf("conversion error")
		},
	}
	repo := NewFileRepository(mockConverter)

	err := repo.WriteFile(testFile, "test content", domain.EncodingShiftJIS)
	if err == nil {
		t.Error("WriteFile() should return error when conversion fails")
	}
}

// TestFileRepositoryImpl_ListFiles_Normal はListFiles()の正常系をテストします
func TestFileRepositoryImpl_ListFiles_Normal(t *testing.T) {
	tempDir := t.TempDir()

	// テスト用ファイルを作成
	testFiles := []string{"file1.txt", "file2.go", "file3.md"}
	for _, fileName := range testFiles {
		filePath := filepath.Join(tempDir, fileName)
		err := os.WriteFile(filePath, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", fileName, err)
		}
	}

	mockConverter := &MockEncodingConverter{}
	repo := NewFileRepository(mockConverter)

	files, err := repo.ListFiles(tempDir, false)
	if err != nil {
		t.Errorf("ListFiles() returned unexpected error: %v", err)
	}

	if len(files) != len(testFiles) {
		t.Errorf("ListFiles() returned %d files, expected %d", len(files), len(testFiles))
	}

	// ファイル名をチェック
	for _, file := range files {
		fileName := filepath.Base(file)
		found := false
		for _, expectedFile := range testFiles {
			if fileName == expectedFile {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Unexpected file in list: %s", fileName)
		}
	}
}

// TestFileRepositoryImpl_ListFiles_Recursive はListFiles()の再帰処理をテストします
func TestFileRepositoryImpl_ListFiles_Recursive(t *testing.T) {
	tempDir := t.TempDir()

	// ディレクトリ構造を作成
	subDir := filepath.Join(tempDir, "subdir")
	err := os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// ファイルを作成
	files := map[string]string{
		"file1.txt":        "content1",
		"subdir/file2.txt": "content2",
		"subdir/file3.go":  "content3",
	}

	for filePath, content := range files {
		fullPath := filepath.Join(tempDir, filePath)
		err := os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", filePath, err)
		}
	}

	mockConverter := &MockEncodingConverter{}
	repo := NewFileRepository(mockConverter)

	listedFiles, err := repo.ListFiles(tempDir, true)
	if err != nil {
		t.Errorf("ListFiles() returned unexpected error: %v", err)
	}

	if len(listedFiles) != len(files) {
		t.Errorf("ListFiles() returned %d files, expected %d", len(listedFiles), len(files))
	}
}

// TestFileRepositoryImpl_ListFiles_DirectoryNotExists はListFiles()のディレクトリ存在しないケースをテストします
func TestFileRepositoryImpl_ListFiles_DirectoryNotExists(t *testing.T) {
	mockConverter := &MockEncodingConverter{}
	repo := NewFileRepository(mockConverter)

	_, err := repo.ListFiles("/nonexistent/directory", false)
	if err == nil {
		t.Error("ListFiles() should return error for non-existent directory")
	}
}

// TestFileRepositoryImpl_CreateBackup_Normal はCreateBackup()の正常系をテストします
func TestFileRepositoryImpl_CreateBackup_Normal(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, World!"

	// 元ファイルを作成
	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	mockConverter := &MockEncodingConverter{}
	repo := NewFileRepository(mockConverter)

	err = repo.CreateBackup(testFile, "")
	if err != nil {
		t.Errorf("CreateBackup() returned unexpected error: %v", err)
	}

	// バックアップファイルが作成されたかチェック
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to read temp directory: %v", err)
	}

	backupFound := false
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "test.txt.backup_") {
			backupFound = true
			// バックアップファイルの内容をチェック
			backupPath := filepath.Join(tempDir, entry.Name())
			backupContent, err := os.ReadFile(backupPath)
			if err != nil {
				t.Errorf("Failed to read backup file: %v", err)
			}
			if string(backupContent) != testContent {
				t.Errorf("Backup content = %s, expected %s", string(backupContent), testContent)
			}
			break
		}
	}

	if !backupFound {
		t.Error("Backup file was not created")
	}
}

// TestFileRepositoryImpl_CreateBackup_WithBackupDir はCreateBackup()のバックアップディレクトリ指定をテストします
func TestFileRepositoryImpl_CreateBackup_WithBackupDir(t *testing.T) {
	tempDir := t.TempDir()
	backupDir := filepath.Join(tempDir, "backups")
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, World!"

	// 元ファイルを作成
	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	mockConverter := &MockEncodingConverter{}
	repo := NewFileRepository(mockConverter)

	err = repo.CreateBackup(testFile, backupDir)
	if err != nil {
		t.Errorf("CreateBackup() returned unexpected error: %v", err)
	}

	// バックアップディレクトリが作成されたかチェック
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		t.Error("Backup directory was not created")
	}

	// バックアップファイルが作成されたかチェック
	expectedBackupSubDir := filepath.Join(backupDir, tempDir)
	entries, err := os.ReadDir(expectedBackupSubDir)
	if err != nil {
		t.Fatalf("Failed to read backup subdirectory: %v", err)
	}

	backupFound := false
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "test.txt.backup_") {
			backupFound = true
			break
		}
	}

	if !backupFound {
		t.Error("Backup file was not created in backup directory")
	}
}

// TestFileRepositoryImpl_CreateBackup_FileNotExists はCreateBackup()のファイル存在しないケースをテストします
func TestFileRepositoryImpl_CreateBackup_FileNotExists(t *testing.T) {
	mockConverter := &MockEncodingConverter{}
	repo := NewFileRepository(mockConverter)

	err := repo.CreateBackup("/nonexistent/file.txt", "")
	if err == nil {
		t.Error("CreateBackup() should return error for non-existent file")
	}
}

// TestFileRepositoryImpl_FileExists_Normal はFileExists()の正常系をテストします
func TestFileRepositoryImpl_FileExists_Normal(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")

	// ファイルを作成
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	mockConverter := &MockEncodingConverter{}
	repo := NewFileRepository(mockConverter)

	// 存在するファイル
	if !repo.FileExists(testFile) {
		t.Error("FileExists() should return true for existing file")
	}

	// 存在しないファイル
	if repo.FileExists("/nonexistent/file.txt") {
		t.Error("FileExists() should return false for non-existent file")
	}
}

// TestFileRepositoryImpl_IsDirectory_Normal はIsDirectory()の正常系をテストします
func TestFileRepositoryImpl_IsDirectory_Normal(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")

	// ファイルを作成
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	mockConverter := &MockEncodingConverter{}
	repo := NewFileRepository(mockConverter)

	// ディレクトリ
	if !repo.IsDirectory(tempDir) {
		t.Error("IsDirectory() should return true for directory")
	}

	// ファイル
	if repo.IsDirectory(testFile) {
		t.Error("IsDirectory() should return false for file")
	}

	// 存在しないパス
	if repo.IsDirectory("/nonexistent/path") {
		t.Error("IsDirectory() should return false for non-existent path")
	}
}

// TestFileRepositoryImpl_GetFileInfo_Normal はGetFileInfo()の正常系をテストします
func TestFileRepositoryImpl_GetFileInfo_Normal(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "Hello, World!"

	// ファイルを作成
	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	mockConverter := &MockEncodingConverter{
		detectEncodingFunc: func(content []byte) (domain.EncodingType, error) {
			return domain.EncodingShiftJIS, nil
		},
	}
	repo := NewFileRepository(mockConverter)

	fileInfo, err := repo.GetFileInfo(testFile)
	if err != nil {
		t.Errorf("GetFileInfo() returned unexpected error: %v", err)
	}

	if fileInfo == nil {
		t.Error("GetFileInfo() should not return nil")
		return
	}

	if fileInfo.Path != testFile {
		t.Errorf("FileInfo.Path = %s, expected %s", fileInfo.Path, testFile)
	}

	if fileInfo.IsDir {
		t.Error("FileInfo.IsDir should be false for file")
	}

	if fileInfo.Size != int64(len(testContent)) {
		t.Errorf("FileInfo.Size = %d, expected %d", fileInfo.Size, len(testContent))
	}

	if fileInfo.Encoding != domain.EncodingShiftJIS {
		t.Errorf("FileInfo.Encoding = %s, expected %s", fileInfo.Encoding, domain.EncodingShiftJIS)
	}
}

// TestFileRepositoryImpl_GetFileInfo_Directory はGetFileInfo()のディレクトリをテストします
func TestFileRepositoryImpl_GetFileInfo_Directory(t *testing.T) {
	tempDir := t.TempDir()

	mockConverter := &MockEncodingConverter{}
	repo := NewFileRepository(mockConverter)

	fileInfo, err := repo.GetFileInfo(tempDir)
	if err != nil {
		t.Errorf("GetFileInfo() returned unexpected error: %v", err)
	}

	if fileInfo == nil {
		t.Error("GetFileInfo() should not return nil")
		return
	}

	if !fileInfo.IsDir {
		t.Error("FileInfo.IsDir should be true for directory")
	}

	if fileInfo.Encoding != domain.EncodingUTF8 {
		t.Errorf("FileInfo.Encoding = %s, expected %s", fileInfo.Encoding, domain.EncodingUTF8)
	}
}

// TestFileRepositoryImpl_GetFileInfo_NotExists はGetFileInfo()のファイル存在しないケースをテストします
func TestFileRepositoryImpl_GetFileInfo_NotExists(t *testing.T) {
	mockConverter := &MockEncodingConverter{}
	repo := NewFileRepository(mockConverter)

	_, err := repo.GetFileInfo("/nonexistent/file.txt")
	if err == nil {
		t.Error("GetFileInfo() should return error for non-existent file")
	}
}

// TestFileRepositoryImpl_IsTextFile_Normal はIsTextFile()の正常系をテストします
func TestFileRepositoryImpl_IsTextFile_Normal(t *testing.T) {
	mockConverter := &MockEncodingConverter{}
	repo := NewFileRepository(mockConverter)

	// 具体的な実装型にキャスト
	repoImpl, ok := repo.(*FileRepositoryImpl)
	if !ok {
		t.Fatal("Failed to cast to *FileRepositoryImpl")
	}

	tests := []struct {
		name     string
		filePath string
		expected bool
	}{
		{
			name:     "テキストファイル(.txt)",
			filePath: "/test/file.txt",
			expected: true,
		},
		{
			name:     "Goファイル(.go)",
			filePath: "/test/file.go",
			expected: true,
		},
		{
			name:     "Markdownファイル(.md)",
			filePath: "/test/file.md",
			expected: true,
		},
		{
			name:     "JSONファイル(.json)",
			filePath: "/test/file.json",
			expected: true,
		},
		{
			name:     "バイナリファイル(.exe)",
			filePath: "/test/file.exe",
			expected: false,
		},
		{
			name:     "画像ファイル(.jpg)",
			filePath: "/test/file.jpg",
			expected: false,
		},
		{
			name:     "拡張子なし",
			filePath: "/test/file",
			expected: false,
		},
		{
			name:     "大文字拡張子(.TXT)",
			filePath: "/test/file.TXT",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repoImpl.IsTextFile(tt.filePath)
			if result != tt.expected {
				t.Errorf("IsTextFile(%s) = %v, expected %v", tt.filePath, result, tt.expected)
			}
		})
	}
}
