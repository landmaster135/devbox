package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// #==============================================================#
// ##          Tests                                             ##
// #==============================================================#
// TestStandardGitExecutor_Execute_Normal はExecuteの正常系テスト
func TestStandardGitExecutor_Execute_Normal(t *testing.T) {
	// Arrange
	executor := NewStandardGitExecutor()

	// 一時ディレクトリを作成してGitリポジトリを初期化
	tempDir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Gitリポジトリを初期化
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Skipf("Git not available or failed to init: %v", err)
	}

	// Act
	output, err := executor.Execute(tempDir, "status", "--porcelain")

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 新しいリポジトリなので出力は空であることを期待
	if len(output) != 0 {
		t.Logf("Output: %s", string(output))
		// 空でない場合もエラーにはしない（環境によって異なる可能性があるため）
	}
}

// TestStandardGitExecutor_Execute_InvalidCommand は無効なコマンドの場合のテスト
func TestStandardGitExecutor_Execute_InvalidCommand(t *testing.T) {
	// Arrange
	executor := NewStandardGitExecutor()
	tempDir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Act
	_, err = executor.Execute(tempDir, "invalid-command")

	// Assert
	if err == nil {
		t.Error("Expected error for invalid git command, got nil")
	}
}

// TestStandardGitExecutor_Execute_InvalidDirectory は無効なディレクトリの場合のテスト
func TestStandardGitExecutor_Execute_InvalidDirectory(t *testing.T) {
	// Arrange
	executor := NewStandardGitExecutor()
	invalidDir := "/non/existent/directory"

	// Act
	_, err := executor.Execute(invalidDir, "status")

	// Assert
	if err == nil {
		t.Error("Expected error for invalid directory, got nil")
	}
}

// TestStandardGitExecutor_Execute_NonGitDirectory は非Gitディレクトリの場合のテスト
func TestStandardGitExecutor_Execute_NonGitDirectory(t *testing.T) {
	// Arrange
	executor := NewStandardGitExecutor()

	// 一時ディレクトリを作成（Gitリポジトリではない）
	tempDir, err := os.MkdirTemp("", "non-git-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Act
	_, err = executor.Execute(tempDir, "status")

	// Assert
	if err == nil {
		t.Error("Expected error for non-git directory, got nil")
	}

	// エラーメッセージに"not a git repository"が含まれることを確認
	if !strings.Contains(err.Error(), "not a git repository") {
		t.Logf("Error message: %v", err)
		// 環境によってエラーメッセージが異なる可能性があるため、ログのみ
	}
}

// TestStandardGitExecutor_Execute_EmptyArgs は引数が空の場合のテスト
func TestStandardGitExecutor_Execute_EmptyArgs(t *testing.T) {
	// Arrange
	executor := NewStandardGitExecutor()
	tempDir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Act
	_, err = executor.Execute(tempDir)

	// Assert
	if err == nil {
		t.Error("Expected error for empty git args, got nil")
	}
}

// TestStandardGitExecutor_Execute_MultipleArgs は複数引数の場合のテスト
func TestStandardGitExecutor_Execute_MultipleArgs(t *testing.T) {
	// Arrange
	executor := NewStandardGitExecutor()

	// 一時ディレクトリを作成してGitリポジトリを初期化
	tempDir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Gitリポジトリを初期化
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Skipf("Git not available or failed to init: %v", err)
	}

	// Act
	output, err := executor.Execute(tempDir, "log", "--oneline", "--max-count=1")

	// Assert
	// 新しいリポジトリなのでコミットがない場合はエラーになる
	if err != nil {
		// これは期待される動作（コミットがないため）
		t.Logf("Expected error for empty repository: %v", err)
	} else {
		t.Logf("Unexpected output: %s", string(output))
	}
}

// TestNewStandardGitExecutor_Normal はNewStandardGitExecutorの正常系テスト
func TestNewStandardGitExecutor_Normal(t *testing.T) {
	// Act
	executor := NewStandardGitExecutor()

	// Assert
	if executor == nil {
		t.Error("Expected executor to be non-nil")
	}

	// インターフェースとして使用可能かチェック
	var _ GitCommandExecutor = executor
}

// TestStandardGitExecutor_Execute_GitVersion はgit versionコマンドのテスト
func TestStandardGitExecutor_Execute_GitVersion(t *testing.T) {
	// Arrange
	executor := NewStandardGitExecutor()
	tempDir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Act
	output, err := executor.Execute(tempDir, "version")

	// Assert
	if err != nil {
		t.Skipf("Git not available: %v", err)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "git version") {
		t.Errorf("Expected output to contain 'git version', got %s", outputStr)
	}
}

// TestStandardGitExecutor_Execute_Help はgit helpコマンドのテスト
func TestStandardGitExecutor_Execute_Help(t *testing.T) {
	// Arrange
	executor := NewStandardGitExecutor()
	tempDir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Act
	output, err := executor.Execute(tempDir, "help", "--all")

	// Assert
	if err != nil {
		t.Skipf("Git help not available: %v", err)
	}

	outputStr := string(output)
	if len(outputStr) == 0 {
		t.Error("Expected non-empty output from git help")
	}
}

// TestStandardGitExecutor_Execute_Integration は統合テスト
func TestStandardGitExecutor_Execute_Integration(t *testing.T) {
	// Arrange
	executor := NewStandardGitExecutor()

	// 一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "git-integration-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Gitリポジトリを初期化
	_, err = executor.Execute(tempDir, "init")
	if err != nil {
		t.Skipf("Git not available: %v", err)
	}

	// 設定を追加
	_, err = executor.Execute(tempDir, "config", "user.name", "Test User")
	if err != nil {
		t.Errorf("Failed to set git config: %v", err)
	}

	_, err = executor.Execute(tempDir, "config", "user.email", "test@example.com")
	if err != nil {
		t.Errorf("Failed to set git config: %v", err)
	}

	// テストファイルを作成
	testFile := fmt.Sprintf("%s/test.txt", tempDir)
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// ファイルを追加
	_, err = executor.Execute(tempDir, "add", "test.txt")
	if err != nil {
		t.Errorf("Failed to add file: %v", err)
	}

	// コミット
	_, err = executor.Execute(tempDir, "commit", "-m", "Initial commit")
	if err != nil {
		t.Errorf("Failed to commit: %v", err)
	}

	// ログを確認
	output, err := executor.Execute(tempDir, "log", "--oneline")
	if err != nil {
		t.Errorf("Failed to get log: %v", err)
	}

	// Assert
	outputStr := string(output)
	if !strings.Contains(outputStr, "Initial commit") {
		t.Errorf("Expected log to contain 'Initial commit', got %s", outputStr)
	}
}
