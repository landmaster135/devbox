package security

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewPathValidator(t *testing.T) {
	tests := []struct {
		name             string
		allowedBasePaths []string
		maxPathLength    int
		expectedMaxLen   int
	}{
		{
			name:             "正常なパラメータ",
			allowedBasePaths: []string{"/home/user"},
			maxPathLength:    1000,
			expectedMaxLen:   1000,
		},
		{
			name:             "最大長が0の場合デフォルト値が設定される",
			allowedBasePaths: []string{"/home/user"},
			maxPathLength:    0,
			expectedMaxLen:   4096,
		},
		{
			name:             "最大長が負の値の場合デフォルト値が設定される",
			allowedBasePaths: []string{"/home/user"},
			maxPathLength:    -1,
			expectedMaxLen:   4096,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewPathValidator(tt.allowedBasePaths, tt.maxPathLength)
			if validator.maxPathLength != tt.expectedMaxLen {
				t.Errorf("maxPathLength = %d, want %d", validator.maxPathLength, tt.expectedMaxLen)
			}
		})
	}
}

func TestNewDefaultPathValidator(t *testing.T) {
	validator := NewDefaultPathValidator()
	if validator == nil {
		t.Fatal("NewDefaultPathValidator() returned nil")
	}
	if validator.maxPathLength != 4096 {
		t.Errorf("maxPathLength = %d, want 4096", validator.maxPathLength)
	}
}

func TestPathValidator_ValidateWorkingDirectory_Normal(t *testing.T) {
	// ホームディレクトリ配下にテスト用ディレクトリを作成
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home directory: %v", err)
	}

	testDir := generateUniqueTestDir(homeDir, "test_normal")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	// テスト用のGitリポジトリを作成
	gitDir := filepath.Join(testDir, ".git")
	if err := os.Mkdir(gitDir, 0755); err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}

	validator := NewPathValidator([]string{homeDir}, 4096)

	validatedPath, err := validator.ValidateWorkingDirectory(testDir)
	if err != nil {
		t.Errorf("ValidateWorkingDirectory() error = %v, want nil", err)
	}

	expectedPath := filepath.Clean(testDir)
	if validatedPath != expectedPath {
		t.Errorf("ValidateWorkingDirectory() = %v, want %v", validatedPath, expectedPath)
	}
}

func TestPathValidator_ValidateWorkingDirectory_EmptyPath(t *testing.T) {
	validator := NewDefaultPathValidator()

	_, err := validator.ValidateWorkingDirectory("")
	if err == nil {
		t.Error("ValidateWorkingDirectory() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "パスが空です") {
		t.Errorf("ValidateWorkingDirectory() error = %v, want error containing 'パスが空です'", err)
	}
}

func TestPathValidator_ValidateWorkingDirectory_TooLongPath(t *testing.T) {
	validator := NewPathValidator([]string{}, 10) // 短い最大長を設定

	longPath := strings.Repeat("a", 20)
	_, err := validator.ValidateWorkingDirectory(longPath)
	if err == nil {
		t.Error("ValidateWorkingDirectory() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "パスが長すぎます") {
		t.Errorf("ValidateWorkingDirectory() error = %v, want error containing 'パスが長すぎます'", err)
	}
}

func TestPathValidator_checkDangerousPatterns(t *testing.T) {
	validator := NewDefaultPathValidator()

	dangerousPaths := []string{
		"/path/with/../traversal",
		"/path/with~home",
		"/path/with$var",
		"/path/with`command`",
		"/path/with;command",
		"/path/with|pipe",
		"/path/with&background",
		"/path/with>redirect",
		"/path/with<redirect",
		"/path/with*wildcard",
		"/path/with?wildcard",
		"/path/with\nnewline",
		"/path/with\rcarriage",
		"/path/with\ttab",
		"/path/with\x00null",
	}

	for _, path := range dangerousPaths {
		t.Run("dangerous_path_"+path, func(t *testing.T) {
			err := validator.checkDangerousPatterns(path)
			if err == nil {
				t.Errorf("checkDangerousPatterns(%s) error = nil, want error", path)
			}
			if !strings.Contains(err.Error(), "危険な文字が含まれています") {
				t.Errorf("checkDangerousPatterns(%s) error = %v, want error containing '危険な文字が含まれています'", path, err)
			}
		})
	}
}

func TestPathValidator_checkDangerousPatterns_SafePath(t *testing.T) {
	validator := NewDefaultPathValidator()

	safePaths := []string{
		"/home/user/project",
		"/var/log/application",
		"/opt/myapp/data",
		"./relative/path",
		"simple_path",
	}

	for _, path := range safePaths {
		t.Run("safe_path_"+path, func(t *testing.T) {
			err := validator.checkDangerousPatterns(path)
			if err != nil {
				t.Errorf("checkDangerousPatterns(%s) error = %v, want nil", path, err)
			}
		})
	}
}

func TestPathValidator_checkAllowedBasePaths(t *testing.T) {
	tempDir := t.TempDir()
	validator := NewPathValidator([]string{tempDir}, 4096)

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "許可されたベースパス配下",
			path:    filepath.Join(tempDir, "subdir"),
			wantErr: false,
		},
		{
			name:    "ベースパスそのもの",
			path:    tempDir,
			wantErr: false,
		},
		{
			name:    "許可されていないパス",
			path:    "/unauthorized/path",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.checkAllowedBasePaths(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkAllowedBasePaths() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPathValidator_checkAllowedBasePaths_NoRestriction(t *testing.T) {
	// 許可されたベースパスが設定されていない場合
	validator := NewPathValidator([]string{}, 4096)

	err := validator.checkAllowedBasePaths("/any/path")
	if err != nil {
		t.Errorf("checkAllowedBasePaths() error = %v, want nil (no restriction)", err)
	}
}

func TestPathValidator_checkDirectoryExists(t *testing.T) {
	validator := NewDefaultPathValidator()

	// 存在するディレクトリ
	tempDir := t.TempDir()
	err := validator.checkDirectoryExists(tempDir)
	if err != nil {
		t.Errorf("checkDirectoryExists() error = %v, want nil", err)
	}

	// 存在しないディレクトリ
	nonExistentDir := filepath.Join(tempDir, "nonexistent")
	err = validator.checkDirectoryExists(nonExistentDir)
	if err == nil {
		t.Error("checkDirectoryExists() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "ディレクトリが存在しません") {
		t.Errorf("checkDirectoryExists() error = %v, want error containing 'ディレクトリが存在しません'", err)
	}

	// ファイルを指定した場合
	tempFile := filepath.Join(tempDir, "testfile")
	if err := os.WriteFile(tempFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	err = validator.checkDirectoryExists(tempFile)
	if err == nil {
		t.Error("checkDirectoryExists() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "指定されたパスはディレクトリではありません") {
		t.Errorf("checkDirectoryExists() error = %v, want error containing 'ディレクトリではありません'", err)
	}
}

func TestPathValidator_checkIsGitRepository(t *testing.T) {
	validator := NewDefaultPathValidator()

	// Gitリポジトリ（.gitディレクトリあり）
	tempDir := t.TempDir()
	gitDir := filepath.Join(tempDir, ".git")
	if err := os.Mkdir(gitDir, 0755); err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}

	err := validator.checkIsGitRepository(tempDir)
	if err != nil {
		t.Errorf("checkIsGitRepository() error = %v, want nil", err)
	}

	// 非Gitリポジトリ
	nonGitDir := t.TempDir()
	err = validator.checkIsGitRepository(nonGitDir)
	if err == nil {
		t.Error("checkIsGitRepository() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "gitリポジトリではありません") {
		t.Errorf("checkIsGitRepository() error = %v, want error containing 'gitリポジトリではありません'", err)
	}
}

func TestPathValidator_ValidateWorkingDirectory_Integration(t *testing.T) {
	// 統合テスト：実際のワークフローをテスト
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home directory: %v", err)
	}

	testDir := generateUniqueTestDir(homeDir, "test_integration")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	// テスト用のGitリポジトリを作成
	gitDir := filepath.Join(testDir, ".git")
	if err := os.Mkdir(gitDir, 0755); err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}

	validator := NewPathValidator([]string{homeDir}, 4096)

	// 相対パスでテスト
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(filepath.Dir(testDir)); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	relativePath := filepath.Base(testDir)
	validatedPath, err := validator.ValidateWorkingDirectory(relativePath)
	if err != nil {
		t.Errorf("ValidateWorkingDirectory() error = %v, want nil", err)
	}

	expectedPath := filepath.Clean(testDir)
	if validatedPath != expectedPath {
		t.Errorf("ValidateWorkingDirectory() = %v, want %v", validatedPath, expectedPath)
	}
}
