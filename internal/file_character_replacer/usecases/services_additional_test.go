package usecases

import (
	"fmt"
	"testing"

	"github.com/landmaster135/devbox/internal/file_character_replacer/domain"
)

// TestFileReplacerService_ReplaceStrings_EncodingDetection はReplaceStrings()のエンコーディング検出をテストします
func TestFileReplacerService_ReplaceStrings_EncodingDetection(t *testing.T) {
	mockFileRepo := &MockFileRepository{
		fileExistsFunc: func(path string) bool {
			return true
		},
		isDirectoryFunc: func(path string) bool {
			return false
		},
		getFileInfoFunc: func(path string) (*domain.FileInfo, error) {
			return &domain.FileInfo{
				Path:     path,
				IsDir:    false,
				Size:     100,
				Encoding: domain.EncodingShiftJIS, // 検出されたエンコーディング
			}, nil
		},
		readFileFunc: func(path string, encoding domain.EncodingType) (string, error) {
			if encoding == domain.EncodingShiftJIS {
				return "This is old content", nil
			}
			return "", fmt.Errorf("unexpected encoding: %s", encoding)
		},
		writeFileFunc: func(path string, content string, encoding domain.EncodingType) error {
			if encoding != domain.EncodingShiftJIS {
				return fmt.Errorf("expected Shift_JIS encoding, got %s", encoding)
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
			Encoding: domain.EncodingShiftJIS, // 明示的にShift_JISを指定
		},
	}

	result, err := service.ReplaceStrings()
	if err != nil {
		t.Errorf("ReplaceStrings() returned unexpected error: %v", err)
	}

	if result.ProcessedFiles != 1 {
		t.Errorf("ProcessedFiles = %d, expected 1", result.ProcessedFiles)
	}

	if result.ReplacedCount != 1 {
		t.Errorf("ReplacedCount = %d, expected 1", result.ReplacedCount)
	}
}

// TestFileReplacerService_ReplaceStrings_MultipleReplacements はReplaceStrings()の複数置換をテストします
func TestFileReplacerService_ReplaceStrings_MultipleReplacements(t *testing.T) {
	mockFileRepo := &MockFileRepository{
		fileExistsFunc: func(path string) bool {
			return true
		},
		isDirectoryFunc: func(path string) bool {
			return false
		},
		readFileFunc: func(path string, encoding domain.EncodingType) (string, error) {
			return "old old old text with old words", nil
		},
		writeFileFunc: func(path string, content string, encoding domain.EncodingType) error {
			expectedContent := "new new new text with new words"
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
		},
	}

	result, err := service.ReplaceStrings()
	if err != nil {
		t.Errorf("ReplaceStrings() returned unexpected error: %v", err)
	}

	if result.ProcessedFiles != 1 {
		t.Errorf("ProcessedFiles = %d, expected 1", result.ProcessedFiles)
	}

	if result.ReplacedCount != 4 {
		t.Errorf("ReplacedCount = %d, expected 4", result.ReplacedCount)
	}
}

// TestFileReplacerService_ReplaceStrings_EmptyFile はReplaceStrings()の空ファイルをテストします
func TestFileReplacerService_ReplaceStrings_EmptyFile(t *testing.T) {
	mockFileRepo := &MockFileRepository{
		fileExistsFunc: func(path string) bool {
			return true
		},
		isDirectoryFunc: func(path string) bool {
			return false
		},
		readFileFunc: func(path string, encoding domain.EncodingType) (string, error) {
			return "", nil // 空ファイル
		},
	}

	service := &FileReplacerService{
		fileRepo:          mockFileRepo,
		encodingConverter: &MockEncodingConverter{},
		config: &domain.ReplacementConfig{
			Target:   "/test/empty.txt",
			From:     "old",
			To:       "new",
			Encoding: domain.EncodingUTF8,
		},
	}

	result, err := service.ReplaceStrings()
	if err != nil {
		t.Errorf("ReplaceStrings() returned unexpected error: %v", err)
	}

	if result.ProcessedFiles != 0 {
		t.Errorf("ProcessedFiles = %d, expected 0", result.ProcessedFiles)
	}

	if result.ReplacedCount != 0 {
		t.Errorf("ReplacedCount = %d, expected 0", result.ReplacedCount)
	}
}

// TestFileReplacerService_ReplaceStrings_LargeFile はReplaceStrings()の大きなファイルをテストします
func TestFileReplacerService_ReplaceStrings_LargeFile(t *testing.T) {
	// 大きなファイル内容をシミュレート
	largeContent := ""
	for i := 0; i < 1000; i++ {
		largeContent += fmt.Sprintf("Line %d with old content\n", i)
	}

	mockFileRepo := &MockFileRepository{
		fileExistsFunc: func(path string) bool {
			return true
		},
		isDirectoryFunc: func(path string) bool {
			return false
		},
		readFileFunc: func(path string, encoding domain.EncodingType) (string, error) {
			return largeContent, nil
		},
		writeFileFunc: func(path string, content string, encoding domain.EncodingType) error {
			// 置換が正しく行われたかチェック
			if !containsSubstring(content, "new content") {
				return fmt.Errorf("content should contain 'new content'")
			}
			if containsSubstring(content, "old content") {
				return fmt.Errorf("content should not contain 'old content'")
			}
			return nil
		},
	}

	service := &FileReplacerService{
		fileRepo:          mockFileRepo,
		encodingConverter: &MockEncodingConverter{},
		config: &domain.ReplacementConfig{
			Target:   "/test/large.txt",
			From:     "old",
			To:       "new",
			Encoding: domain.EncodingUTF8,
		},
	}

	result, err := service.ReplaceStrings()
	if err != nil {
		t.Errorf("ReplaceStrings() returned unexpected error: %v", err)
	}

	if result.ProcessedFiles != 1 {
		t.Errorf("ProcessedFiles = %d, expected 1", result.ProcessedFiles)
	}

	if result.ReplacedCount != 1000 {
		t.Errorf("ReplacedCount = %d, expected 1000", result.ReplacedCount)
	}
}

// TestFileReplacerService_ReplaceStrings_SpecialCharacters はReplaceStrings()の特殊文字をテストします
func TestFileReplacerService_ReplaceStrings_SpecialCharacters(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		from     string
		to       string
		expected string
	}{
		{
			name:     "改行文字",
			content:  "line1\nold\nline3",
			from:     "\n",
			to:       " ",
			expected: "line1 old line3",
		},
		{
			name:     "タブ文字",
			content:  "col1\told\tcol3",
			from:     "\t",
			to:       ",",
			expected: "col1,old,col3",
		},
		{
			name:     "Unicode文字",
			content:  "Hello 世界 old text",
			from:     "世界",
			to:       "world",
			expected: "Hello world old text",
		},
		{
			name:     "正規表現特殊文字",
			content:  "price: $100.50",
			from:     "$100.50",
			to:       "$200.00",
			expected: "price: $200.00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFileRepo := &MockFileRepository{
				fileExistsFunc: func(path string) bool {
					return true
				},
				isDirectoryFunc: func(path string) bool {
					return false
				},
				readFileFunc: func(path string, encoding domain.EncodingType) (string, error) {
					return tt.content, nil
				},
				writeFileFunc: func(path string, content string, encoding domain.EncodingType) error {
					if content != tt.expected {
						return fmt.Errorf("expected %q, got %q", tt.expected, content)
					}
					return nil
				},
			}

			service := &FileReplacerService{
				fileRepo:          mockFileRepo,
				encodingConverter: &MockEncodingConverter{},
				config: &domain.ReplacementConfig{
					Target:   "/test/file.txt",
					From:     tt.from,
					To:       tt.to,
					Encoding: domain.EncodingUTF8,
				},
			}

			result, err := service.ReplaceStrings()
			if err != nil {
				t.Errorf("ReplaceStrings() returned unexpected error: %v", err)
			}

			if result.ProcessedFiles != 1 {
				t.Errorf("ProcessedFiles = %d, expected 1", result.ProcessedFiles)
			}
		})
	}
}

// TestFileReplacerService_ReplaceStrings_CaseSensitive はReplaceStrings()の大文字小文字区別をテストします
func TestFileReplacerService_ReplaceStrings_CaseSensitive(t *testing.T) {
	mockFileRepo := &MockFileRepository{
		fileExistsFunc: func(path string) bool {
			return true
		},
		isDirectoryFunc: func(path string) bool {
			return false
		},
		readFileFunc: func(path string, encoding domain.EncodingType) (string, error) {
			return "Old old OLD content", nil
		},
		writeFileFunc: func(path string, content string, encoding domain.EncodingType) error {
			expectedContent := "Old new OLD content"
			if content != expectedContent {
				return fmt.Errorf("expected %q, got %q", expectedContent, content)
			}
			return nil
		},
	}

	service := &FileReplacerService{
		fileRepo:          mockFileRepo,
		encodingConverter: &MockEncodingConverter{},
		config: &domain.ReplacementConfig{
			Target:   "/test/file.txt",
			From:     "old", // 小文字のみ
			To:       "new",
			Encoding: domain.EncodingUTF8,
		},
	}

	result, err := service.ReplaceStrings()
	if err != nil {
		t.Errorf("ReplaceStrings() returned unexpected error: %v", err)
	}

	if result.ProcessedFiles != 1 {
		t.Errorf("ProcessedFiles = %d, expected 1", result.ProcessedFiles)
	}

	if result.ReplacedCount != 1 {
		t.Errorf("ReplacedCount = %d, expected 1 (only lowercase 'old' should be replaced)", result.ReplacedCount)
	}
}

// TestFileReplacerService_ReplaceStrings_SameFromTo はReplaceStrings()の同じfrom/toをテストします
func TestFileReplacerService_ReplaceStrings_SameFromTo(t *testing.T) {
	mockFileRepo := &MockFileRepository{
		fileExistsFunc: func(path string) bool {
			return true
		},
		isDirectoryFunc: func(path string) bool {
			return false
		},
		readFileFunc: func(path string, encoding domain.EncodingType) (string, error) {
			return "same same same", nil
		},
		writeFileFunc: func(path string, content string, encoding domain.EncodingType) error {
			expectedContent := "same same same" // 内容は同じだが置換は発生
			if content != expectedContent {
				return fmt.Errorf("expected %q, got %q", expectedContent, content)
			}
			return nil
		},
	}

	service := &FileReplacerService{
		fileRepo:          mockFileRepo,
		encodingConverter: &MockEncodingConverter{},
		config: &domain.ReplacementConfig{
			Target:   "/test/file.txt",
			From:     "same",
			To:       "same", // 同じ文字列
			Encoding: domain.EncodingUTF8,
		},
	}

	result, err := service.ReplaceStrings()
	if err != nil {
		t.Errorf("ReplaceStrings() returned unexpected error: %v", err)
	}

	if result.ProcessedFiles != 0 {
		t.Errorf("ProcessedFiles = %d, expected 0 (no actual change)", result.ProcessedFiles)
	}

	if result.ReplacedCount != 0 {
		t.Errorf("ReplacedCount = %d, expected 0 (no actual change)", result.ReplacedCount)
	}
}

// TestFileReplacerService_GetSummary_NoMessages はGetSummary()のメッセージなしケースをテストします
func TestFileReplacerService_GetSummary_NoMessages(t *testing.T) {
	service := &FileReplacerService{
		config: &domain.ReplacementConfig{
			DryRun: false,
		},
	}

	result := &domain.FileProcessResult{
		ProcessedFiles: 1,
		ReplacedCount:  2,
		Messages:       []string{}, // メッセージなし
		Errors:         []error{},  // エラーなし
	}

	summary := service.GetSummary(result)

	if summary == "" {
		t.Error("GetSummary() should not return empty string")
	}

	expectedContents := []string{
		"処理されたファイル数: 1",
		"置換された箇所数: 2",
	}

	for _, expected := range expectedContents {
		if !contains(summary, expected) {
			t.Errorf("Summary should contain '%s'", expected)
		}
	}

	// メッセージセクションが含まれていないことを確認
	if contains(summary, "=== 処理詳細 ===") {
		t.Error("Summary should not contain message section when no messages")
	}

	// エラーセクションが含まれていないことを確認
	if contains(summary, "=== エラー ===") {
		t.Error("Summary should not contain error section when no errors")
	}
}

// TestFileReplacerService_GetSummary_EmptyResult はGetSummary()の空結果をテストします
func TestFileReplacerService_GetSummary_EmptyResult(t *testing.T) {
	service := &FileReplacerService{
		config: &domain.ReplacementConfig{
			DryRun: false,
		},
	}

	result := &domain.FileProcessResult{
		ProcessedFiles: 0,
		ReplacedCount:  0,
		Messages:       []string{},
		Errors:         []error{},
	}

	summary := service.GetSummary(result)

	if summary == "" {
		t.Error("GetSummary() should not return empty string")
	}

	expectedContents := []string{
		"処理されたファイル数: 0",
		"置換された箇所数: 0",
	}

	for _, expected := range expectedContents {
		if !contains(summary, expected) {
			t.Errorf("Summary should contain '%s'", expected)
		}
	}
}

// TestFileReplacerService_isTextFile_EdgeCases はisTextFile()のエッジケースをテストします
func TestFileReplacerService_isTextFile_EdgeCases(t *testing.T) {
	service := NewFileReplacerService()

	tests := []struct {
		name     string
		filePath string
		expected bool
	}{
		{
			name:     "複数ドット(.tar.gz)",
			filePath: "/test/file.tar.gz",
			expected: false,
		},
		{
			name:     "隠しファイル(.gitignore)",
			filePath: "/test/.gitignore",
			expected: false,
		},
		{
			name:     "パスなし(file.txt)",
			filePath: "file.txt",
			expected: true,
		},
		{
			name:     "長いパス",
			filePath: "/very/long/path/to/deep/directory/structure/file.go",
			expected: true,
		},
		{
			name:     "Windows形式パス",
			filePath: "C:\\Users\\test\\file.py",
			expected: true,
		},
		{
			name:     "数字のみ拡張子(.123)",
			filePath: "/test/file.123",
			expected: false,
		},
		{
			name:     "空文字列",
			filePath: "",
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
