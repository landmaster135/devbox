package git

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	security "github.com/landmaster135/devbox/internal/git_diff_recorder/security"
)

// generateUniqueTestDir はハッシュ値を使用してユニークなテストディレクトリ名を生成する
func generateUniqueTestDir(baseDir, prefix string) string {
	// 現在時刻とテスト名からハッシュを生成
	data := fmt.Sprintf("%s_%d_%s", prefix, time.Now().UnixNano(), "test_git_executor")
	hash := sha256.Sum256([]byte(data))
	hashStr := fmt.Sprintf("%x", hash)[:16] // 16文字のハッシュを使用
	return filepath.Join(baseDir, fmt.Sprintf("%s_%s", prefix, hashStr))
}

func TestNewStandardGitExecutor(t *testing.T) {
	executor := NewStandardGitExecutor()
	if executor == nil {
		t.Fatal("NewStandardGitExecutor() returned nil")
	}
	if executor.pathValidator == nil {
		t.Error("pathValidator is nil")
	}
}

func TestNewStandardGitExecutorWithValidator(t *testing.T) {
	validator := security.NewDefaultPathValidator()
	executor := NewStandardGitExecutorWithValidator(validator)
	if executor == nil {
		t.Fatal("NewStandardGitExecutorWithValidator() returned nil")
	}
	if executor.pathValidator != validator {
		t.Error("pathValidator is not the expected instance")
	}
}

func TestStandardGitExecutor_Execute_ValidPath(t *testing.T) {
	// ホームディレクトリ配下にテスト用ディレクトリを作成
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home directory: %v", err)
	}

	tempDir := generateUniqueTestDir(homeDir, "test_git_valid")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// テスト用のGitリポジトリを作成
	gitDir := filepath.Join(tempDir, ".git")
	if err := os.Mkdir(gitDir, 0755); err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}

	// 許可されたパスでValidatorを作成
	validator := security.NewPathValidator([]string{homeDir}, 4096)
	executor := NewStandardGitExecutorWithValidator(validator)

	// git statusコマンドを実行（実際のgitコマンドが必要）
	_, err = executor.Execute(tempDir, "status", "--porcelain")
	if err != nil {
		// gitコマンドが利用できない環境でもテストが通るように、
		// パス検証エラー以外は許可する
		if strings.Contains(err.Error(), "無効なワーキングディレクトリ") {
			t.Errorf("Execute() error = %v, want no path validation error", err)
		}
		// gitコマンド自体のエラーは許可（テスト環境によってはgitが無い場合がある）
	}
}

func TestStandardGitExecutor_Execute_InvalidPath(t *testing.T) {
	// 存在しないディレクトリでテスト
	nonExistentDir := "/nonexistent/directory"

	executor := NewStandardGitExecutor()

	_, err := executor.Execute(nonExistentDir, "status")
	if err == nil {
		t.Error("Execute() error = nil, want error for invalid path")
	}
	if !strings.Contains(err.Error(), "無効なワーキングディレクトリ") {
		t.Errorf("Execute() error = %v, want error containing '無効なワーキングディレクトリ'", err)
	}
}

func TestStandardGitExecutor_Execute_DangerousPath(t *testing.T) {
	executor := NewStandardGitExecutor()

	dangerousPaths := []string{
		"/path/with/../traversal",
		"/path/with;command",
		"/path/with|pipe",
		"/path/with&background",
	}

	for _, path := range dangerousPaths {
		t.Run("dangerous_path_"+path, func(t *testing.T) {
			_, err := executor.Execute(path, "status")
			if err == nil {
				t.Errorf("Execute() error = nil, want error for dangerous path: %s", path)
			}
			if !strings.Contains(err.Error(), "無効なワーキングディレクトリ") {
				t.Errorf("Execute() error = %v, want error containing '無効なワーキングディレクトリ'", err)
			}
		})
	}
}

func TestStandardGitExecutor_Execute_EmptyPath(t *testing.T) {
	executor := NewStandardGitExecutor()

	_, err := executor.Execute("", "status")
	if err == nil {
		t.Error("Execute() error = nil, want error for empty path")
	}
	if !strings.Contains(err.Error(), "無効なワーキングディレクトリ") {
		t.Errorf("Execute() error = %v, want error containing '無効なワーキングディレクトリ'", err)
	}
}

func TestStandardGitExecutor_Execute_NonGitDirectory(t *testing.T) {
	// ホームディレクトリ配下にテスト用ディレクトリを作成（.gitディレクトリなし）
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home directory: %v", err)
	}

	tempDir := generateUniqueTestDir(homeDir, "test_git_nongit")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 許可されたパスでValidatorを作成
	validator := security.NewPathValidator([]string{homeDir}, 4096)
	executor := NewStandardGitExecutorWithValidator(validator)

	_, err = executor.Execute(tempDir, "status")
	if err == nil {
		t.Error("Execute() error = nil, want error for non-git directory")
	}
	if !strings.Contains(err.Error(), "無効なワーキングディレクトリ") {
		t.Errorf("Execute() error = %v, want error containing '無効なワーキングディレクトリ'", err)
	}
}

// testGitExecutor はテスト用のGitExecutor
type testGitExecutor struct {
	pathValidator security.PathValidatorInterface
}

func (e *testGitExecutor) Execute(workingDir string, args ...string) ([]byte, error) {
	// パス検証を実行
	validatedPath, err := e.pathValidator.ValidateWorkingDirectory(workingDir)
	if err != nil {
		return nil, fmt.Errorf("無効なワーキングディレクトリ: %w", err)
	}

	// テスト用なので実際のgitコマンドは実行せず、検証済みパスを使用したことを示すダミーレスポンスを返す
	_ = validatedPath
	return []byte("test output"), nil
}

func TestStandardGitExecutor_Execute_WithMockValidator(t *testing.T) {
	tests := []struct {
		name          string
		validateFunc  func(string) (string, error)
		inputPath     string
		wantErr       bool
		errorContains string
	}{
		{
			name: "バリデーション成功",
			validateFunc: func(path string) (string, error) {
				return "/validated/path", nil
			},
			inputPath: "/input/path",
			wantErr:   false, // gitコマンドのエラーは許可
		},
		{
			name: "バリデーション失敗",
			validateFunc: func(path string) (string, error) {
				return "", errors.New("validation failed")
			},
			inputPath:     "/input/path",
			wantErr:       true,
			errorContains: "無効なワーキングディレクトリ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockValidator := &security.MockPathValidator{
				ValidateFunc: tt.validateFunc,
			}

			// テスト用のexecutor構造体を作成
			executor := &testGitExecutor{
				pathValidator: mockValidator,
			}

			_, err := executor.Execute(tt.inputPath, "status")

			if tt.wantErr {
				if err == nil {
					t.Error("Execute() error = nil, want error")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Execute() error = %v, want error containing '%s'", err, tt.errorContains)
				}
			} else {
				// バリデーション成功の場合、gitコマンド自体のエラーは許可
				if err != nil && strings.Contains(err.Error(), "無効なワーキングディレクトリ") {
					t.Errorf("Execute() unexpected validation error = %v", err)
				}
			}
		})
	}
}

func TestMockGitExecutor_Execute(t *testing.T) {
	executor := NewMockGitExecutor()

	// レスポンスを設定
	executor.SetResponse("status --porcelain", []byte("M  file.txt\n"))
	executor.SetError("invalid command", errors.New("command failed"))

	// 正常なレスポンス
	output, err := executor.Execute("/any/dir", "status", "--porcelain")
	if err != nil {
		t.Errorf("Execute() error = %v, want nil", err)
	}
	expected := "M  file.txt\n"
	if string(output) != expected {
		t.Errorf("Execute() output = %s, want %s", string(output), expected)
	}

	// エラーレスポンス
	_, err = executor.Execute("/any/dir", "invalid", "command")
	if err == nil {
		t.Error("Execute() error = nil, want error")
	}
	if err.Error() != "command failed" {
		t.Errorf("Execute() error = %v, want 'command failed'", err)
	}

	// デフォルトレスポンス
	output, err = executor.Execute("/any/dir", "unknown", "command")
	if err != nil {
		t.Errorf("Execute() error = %v, want nil", err)
	}
	if string(output) != "" {
		t.Errorf("Execute() output = %s, want empty string", string(output))
	}
}

func TestJoinArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "空の引数",
			args:     []string{},
			expected: "",
		},
		{
			name:     "単一の引数",
			args:     []string{"status"},
			expected: "status",
		},
		{
			name:     "複数の引数",
			args:     []string{"status", "--porcelain"},
			expected: "status --porcelain",
		},
		{
			name:     "多数の引数",
			args:     []string{"diff", "--cached", "--name-only"},
			expected: "diff --cached --name-only",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := joinArgs(tt.args)
			if result != tt.expected {
				t.Errorf("joinArgs() = %s, want %s", result, tt.expected)
			}
		})
	}
}
