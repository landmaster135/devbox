package usecases

import (
	"fmt"
	"testing"

	"github.com/landmaster135/devbox/internal/file_character_replacer/domain"
)

// MockFileRepository はテスト用のFileRepositoryモック
type MockFileRepository struct {
	readFileFunc     func(path string, encoding domain.EncodingType) (string, error)
	writeFileFunc    func(path string, content string, encoding domain.EncodingType) error
	listFilesFunc    func(dirPath string, recursive bool) ([]string, error)
	createBackupFunc func(filePath string, backupDir string) error
	fileExistsFunc   func(path string) bool
	isDirectoryFunc  func(path string) bool
	getFileInfoFunc  func(path string) (*domain.FileInfo, error)
}

// ReadFile はファイルを読み込みます
func (m *MockFileRepository) ReadFile(path string, encoding domain.EncodingType) (string, error) {
	if m.readFileFunc != nil {
		return m.readFileFunc(path, encoding)
	}
	return "", nil
}

// WriteFile はファイルに書き込みます
func (m *MockFileRepository) WriteFile(path string, content string, encoding domain.EncodingType) error {
	if m.writeFileFunc != nil {
		return m.writeFileFunc(path, content, encoding)
	}
	return nil
}

// ListFiles はディレクトリ内のファイル一覧を取得します
func (m *MockFileRepository) ListFiles(dirPath string, recursive bool) ([]string, error) {
	if m.listFilesFunc != nil {
		return m.listFilesFunc(dirPath, recursive)
	}
	return []string{}, nil
}

// CreateBackup はファイルのバックアップを作成します
func (m *MockFileRepository) CreateBackup(filePath string, backupDir string) error {
	if m.createBackupFunc != nil {
		return m.createBackupFunc(filePath, backupDir)
	}
	return nil
}

// FileExists はファイルまたはディレクトリが存在するかを確認します
func (m *MockFileRepository) FileExists(path string) bool {
	if m.fileExistsFunc != nil {
		return m.fileExistsFunc(path)
	}
	return true
}

// IsDirectory はパスがディレクトリかどうかを確認します
func (m *MockFileRepository) IsDirectory(path string) bool {
	if m.isDirectoryFunc != nil {
		return m.isDirectoryFunc(path)
	}
	return false
}

// GetFileInfo はファイル情報を取得します
func (m *MockFileRepository) GetFileInfo(path string) (*domain.FileInfo, error) {
	if m.getFileInfoFunc != nil {
		return m.getFileInfoFunc(path)
	}
	return &domain.FileInfo{
		Path:     path,
		IsDir:    false,
		Size:     100,
		Encoding: domain.EncodingUTF8,
	}, nil
}

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

// TestNewFileReplacerService_Normal はNewFileReplacerService()の正常系をテストします
func TestNewFileReplacerService_Normal(t *testing.T) {
	service := NewFileReplacerService()

	if service == nil {
		t.Error("NewFileReplacerService() should not return nil")
	}

	if service.fileRepo == nil {
		t.Error("fileRepo should not be nil")
	}

	if service.encodingConverter == nil {
		t.Error("encodingConverter should not be nil")
	}

	if service.config != nil {
		t.Error("config should be nil initially")
	}
}

// TestFileReplacerService_SetConfig_Normal はSetConfig()をテストします
func TestFileReplacerService_SetConfig_Normal(t *testing.T) {
	service := NewFileReplacerService()
	config := &domain.ReplacementConfig{
		Target:   "/test/path",
		From:     "old",
		To:       "new",
		Encoding: domain.EncodingUTF8,
	}

	service.SetConfig(config)

	if service.config != config {
		t.Error("SetConfig() should set the config")
	}
}

// TestFileReplacerService_ReplaceStrings_ValidationError はReplaceStrings()のバリデーションエラーをテストします
func TestFileReplacerService_ReplaceStrings_ValidationError(t *testing.T) {
	service := NewFileReplacerService()

	// 無効な設定（Targetが空）
	config := &domain.ReplacementConfig{
		Target:   "",
		From:     "old",
		To:       "new",
		Encoding: domain.EncodingUTF8,
	}
	service.SetConfig(config)

	result, err := service.ReplaceStrings()

	if err == nil {
		t.Error("ReplaceStrings() should return validation error")
	}

	if result == nil {
		t.Error("ReplaceStrings() should return result even on error")
	}

	if !result.HasErrors() {
		t.Error("Result should have errors")
	}
}

// TestFileReplacerService_ReplaceStrings_FileNotExists はReplaceStrings()のファイル存在チェックをテストします
func TestFileReplacerService_ReplaceStrings_FileNotExists(t *testing.T) {
	mockFileRepo := &MockFileRepository{
		fileExistsFunc: func(path string) bool {
			return false // ファイルが存在しない
		},
	}

	service := &FileReplacerService{
		fileRepo:          mockFileRepo,
		encodingConverter: &MockEncodingConverter{},
		config: &domain.ReplacementConfig{
			Target:   "/test/path",
			From:     "old",
			To:       "new",
			Encoding: domain.EncodingUTF8,
		},
	}

	result, err := service.ReplaceStrings()

	if err == nil {
		t.Error("ReplaceStrings() should return error when file does not exist")
	}

	if !result.HasErrors() {
		t.Error("Result should have errors")
	}
}

// TestFileReplacerService_ReplaceStrings_SingleFile_Normal はReplaceStrings()の単一ファイル処理をテストします
func TestFileReplacerService_ReplaceStrings_SingleFile_Normal(t *testing.T) {
	mockFileRepo := &MockFileRepository{
		fileExistsFunc: func(path string) bool {
			return true
		},
		isDirectoryFunc: func(path string) bool {
			return false // ファイル
		},
		readFileFunc: func(path string, encoding domain.EncodingType) (string, error) {
			return "This is old content with old text", nil
		},
		writeFileFunc: func(path string, content string, encoding domain.EncodingType) error {
			expectedContent := "This is new content with new text"
			if content != expectedContent {
				return fmt.Errorf("unexpected content: %s", content)
			}
			return nil
		},
	}

	service := &FileReplacerService{
		fileRepo:          mockFileRepo,
		encodingConverter: &MockEncodingConverter{},
		config: &domain.ReplacementConfig{
			Target:   "/test/file.txt",
			From:     "old",
			To:       "new",
			Encoding: domain.EncodingUTF8,
			Backup:   false,
			DryRun:   false,
		},
	}

	result, err := service.ReplaceStrings()

	if err != nil {
		t.Errorf("ReplaceStrings() returned unexpected error: %v", err)
	}

	if result.ProcessedFiles != 1 {
		t.Errorf("ProcessedFiles = %d, expected 1", result.ProcessedFiles)
	}

	if result.ReplacedCount != 2 {
		t.Errorf("ReplacedCount = %d, expected 2", result.ReplacedCount)
	}

	if result.HasErrors() {
		t.Error("Result should not have errors")
	}
}

// TestFileReplacerService_ReplaceStrings_SingleFile_DryRun はReplaceStrings()のドライランをテストします
func TestFileReplacerService_ReplaceStrings_SingleFile_DryRun(t *testing.T) {
	writeFileCalled := false
	mockFileRepo := &MockFileRepository{
		fileExistsFunc: func(path string) bool {
			return true
		},
		isDirectoryFunc: func(path string) bool {
			return false
		},
		readFileFunc: func(path string, encoding domain.EncodingType) (string, error) {
			return "This is old content", nil
		},
		writeFileFunc: func(path string, content string, encoding domain.EncodingType) error {
			writeFileCalled = true
			return nil
		},
	}

	service := &FileReplacerService{
		fileRepo:          mockFileRepo,
		encodingConverter: &MockEncodingConverter{},
		config: &domain.ReplacementConfig{
			Target:   "/test/file.txt",
			From:     "old",
			To:       "new",
			Encoding: domain.EncodingUTF8,
			DryRun:   true, // ドライラン
		},
	}

	result, err := service.ReplaceStrings()

	if err != nil {
		t.Errorf("ReplaceStrings() returned unexpected error: %v", err)
	}

	if writeFileCalled {
		t.Error("WriteFile should not be called in dry run mode")
	}

	if result.ProcessedFiles != 1 {
		t.Errorf("ProcessedFiles = %d, expected 1", result.ProcessedFiles)
	}

	if result.ReplacedCount != 1 {
		t.Errorf("ReplacedCount = %d, expected 1", result.ReplacedCount)
	}
}

// TestFileReplacerService_ReplaceStrings_SingleFile_WithBackup はReplaceStrings()のバックアップ機能をテストします
func TestFileReplacerService_ReplaceStrings_SingleFile_WithBackup(t *testing.T) {
	backupCalled := false
	mockFileRepo := &MockFileRepository{
		fileExistsFunc: func(path string) bool {
			return true
		},
		isDirectoryFunc: func(path string) bool {
			return false
		},
		readFileFunc: func(path string, encoding domain.EncodingType) (string, error) {
			return "This is old content", nil
		},
		writeFileFunc: func(path string, content string, encoding domain.EncodingType) error {
			return nil
		},
		createBackupFunc: func(filePath string, backupDir string) error {
			backupCalled = true
			if filePath != "/test/file.txt" {
				return fmt.Errorf("unexpected file path: %s", filePath)
			}
			return nil
		},
	}

	service := &FileReplacerService{
		fileRepo:          mockFileRepo,
		encodingConverter: &MockEncodingConverter{},
		config: &domain.ReplacementConfig{
			Target:   "/test/file.txt",
			From:     "old",
			To:       "new",
			Encoding: domain.EncodingUTF8,
			Backup:   true,
		},
	}

	result, err := service.ReplaceStrings()

	if err != nil {
		t.Errorf("ReplaceStrings() returned unexpected error: %v", err)
	}

	if !backupCalled {
		t.Error("CreateBackup should be called when backup is enabled")
	}

	if result.ProcessedFiles != 1 {
		t.Errorf("ProcessedFiles = %d, expected 1", result.ProcessedFiles)
	}
}

// TestFileReplacerService_ReplaceStrings_Directory_Normal はReplaceStrings()のディレクトリ処理をテストします
func TestFileReplacerService_ReplaceStrings_Directory_Normal(t *testing.T) {
	mockFileRepo := &MockFileRepository{
		fileExistsFunc: func(path string) bool {
			return true
		},
		isDirectoryFunc: func(path string) bool {
			return true // ディレクトリ
		},
		listFilesFunc: func(dirPath string, recursive bool) ([]string, error) {
			return []string{"/test/dir/file1.txt", "/test/dir/file2.go"}, nil
		},
		readFileFunc: func(path string, encoding domain.EncodingType) (string, error) {
			return "This is old content", nil
		},
		writeFileFunc: func(path string, content string, encoding domain.EncodingType) error {
			return nil
		},
	}

	service := &FileReplacerService{
		fileRepo:          mockFileRepo,
		encodingConverter: &MockEncodingConverter{},
		config: &domain.ReplacementConfig{
			Target:    "/test/dir",
			From:      "old",
			To:        "new",
			Encoding:  domain.EncodingUTF8,
			Recursive: true,
		},
	}

	result, err := service.ReplaceStrings()

	if err != nil {
		t.Errorf("ReplaceStrings() returned unexpected error: %v", err)
	}

	if result.ProcessedFiles != 2 {
		t.Errorf("ProcessedFiles = %d, expected 2", result.ProcessedFiles)
	}

	if result.ReplacedCount != 2 {
		t.Errorf("ReplacedCount = %d, expected 2", result.ReplacedCount)
	}
}

// TestFileReplacerService_isTextFile_Normal はisTextFile()をテストします
func TestFileReplacerService_isTextFile_Normal(t *testing.T) {
	service := NewFileReplacerService()

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.isTextFile(tt.filePath)
			if result != tt.expected {
				t.Errorf("isTextFile(%s) = %v, expected %v", tt.filePath, result, tt.expected)
			}
		})
	}
}

// TestFileReplacerService_GetSummary_Normal はGetSummary()をテストします
func TestFileReplacerService_GetSummary_Normal(t *testing.T) {
	service := &FileReplacerService{
		config: &domain.ReplacementConfig{
			DryRun: false,
		},
	}

	result := &domain.FileProcessResult{
		ProcessedFiles: 3,
		ReplacedCount:  5,
		Messages:       []string{"ファイル1処理完了", "ファイル2処理完了"},
		Errors:         []error{fmt.Errorf("テストエラー")},
	}

	summary := service.GetSummary(result)

	if summary == "" {
		t.Error("GetSummary() should not return empty string")
	}

	// サマリーに期待される内容が含まれているかチェック
	expectedContents := []string{
		"処理されたファイル数: 3",
		"置換された箇所数: 5",
		"ファイル1処理完了",
		"ファイル2処理完了",
		"テストエラー",
	}

	for _, expected := range expectedContents {
		if !contains(summary, expected) {
			t.Errorf("Summary should contain '%s'", expected)
		}
	}
}

// TestFileReplacerService_GetSummary_DryRun はGetSummary()のドライランをテストします
func TestFileReplacerService_GetSummary_DryRun(t *testing.T) {
	service := &FileReplacerService{
		config: &domain.ReplacementConfig{
			DryRun: true,
		},
	}

	result := &domain.FileProcessResult{
		ProcessedFiles: 1,
		ReplacedCount:  2,
	}

	summary := service.GetSummary(result)

	if !contains(summary, "ドライラン") {
		t.Error("Summary should contain 'ドライラン' for dry run mode")
	}
}

// contains は文字列に指定された部分文字列が含まれているかチェックします
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
		s[len(s)-len(substr):] == substr ||
		containsSubstring(s, substr))))
}

// containsSubstring は文字列内の部分文字列を検索します
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
