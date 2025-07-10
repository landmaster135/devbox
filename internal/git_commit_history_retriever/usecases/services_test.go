package usecases

import (
	"fmt"
	"strings"
	"testing"

	"github.com/landmaster135/devbox/internal/git_commit_history_retriever/config"
	"github.com/landmaster135/devbox/internal/git_commit_history_retriever/git"
)

// MockGitExecutor はテスト用のモックGitExecutor
type MockGitExecutor struct {
	ExecuteFunc func(workingDir string, args ...string) ([]byte, error)
}

// Execute はGitコマンドを実行する（モック）
func (m *MockGitExecutor) Execute(workingDir string, args ...string) ([]byte, error) {
	return m.ExecuteFunc(workingDir, args...)
}

// TestGitCommitHistoryService_GetCommitHistory_Normal はGetCommitHistoryの正常系テスト
func TestGitCommitHistoryService_GetCommitHistory_Normal(t *testing.T) {
	// Arrange
	mockExecutor := &MockGitExecutor{
		ExecuteFunc: func(workingDir string, args ...string) ([]byte, error) {
			if len(args) > 0 && args[0] == "rev-parse" {
				return []byte(".git"), nil
			}
			if len(args) > 0 && args[0] == "log" {
				return []byte("* abc1234 - feat: test commit (1 hour ago) <test@example.com>"), nil
			}
			return []byte(""), nil
		},
	}

	cfg := &config.Config{
		GitDir:  "/test/repo",
		Keyword: "",
		Since:   "",
		Until:   "",
	}

	gitClient := git.NewClientWithExecutor("/test/repo", mockExecutor)
	service := &GitCommitHistoryService{
		gitClient: gitClient,
		config:    cfg,
	}

	// Act
	result, err := service.GetCommitHistory()

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expectedHistory := "=== Commit History ===\n* abc1234 - feat: test commit (1 hour ago) <test@example.com>"
	if result != expectedHistory {
		t.Errorf("Expected %q, got %q", expectedHistory, result)
	}
}

// TestGitCommitHistoryService_GetCommitHistoryWithDetails_Normal はGetCommitHistoryWithDetailsの正常系テスト
func TestGitCommitHistoryService_GetCommitHistoryWithDetails_Normal(t *testing.T) {
	// Arrange
	mockExecutor := &MockGitExecutor{
		ExecuteFunc: func(workingDir string, args ...string) ([]byte, error) {
			if len(args) > 0 && args[0] == "rev-parse" {
				return []byte(".git"), nil
			}
			if len(args) > 0 && args[0] == "log" {
				return []byte("* abc1234 - feat: test commit (1 hour ago) <test@example.com>\n* def5678 - fix: bug fix (2 hours ago) <test@example.com>"), nil
			}
			if len(args) > 0 && args[0] == "show" {
				return []byte("commit abc1234\nAuthor: test@example.com\nDate: test date\n\n    feat: test commit\n\n file1.go | 10 ++++++++++\n 1 file changed, 10 insertions(+)\n\ncommit def5678\nAuthor: test@example.com\nDate: test date\n\n    fix: bug fix\n\n file2.go | 5 ++---\n 1 file changed, 2 insertions(+), 3 deletions(-)"), nil
			}
			return []byte(""), nil
		},
	}

	cfg := &config.Config{
		GitDir:  "/test/repo",
		Keyword: "",
		Since:   "",
		Until:   "",
	}

	gitClient := git.NewClientWithExecutor("/test/repo", mockExecutor)
	service := &GitCommitHistoryService{
		gitClient: gitClient,
		config:    cfg,
	}

	// Act
	result, err := service.GetCommitHistoryWithDetails()

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !strings.Contains(result, "=== Commit History ===") {
		t.Error("Expected result to contain commit history header")
	}

	if !strings.Contains(result, "=== Commit Details ===") {
		t.Error("Expected result to contain commit details header")
	}

	if !strings.Contains(result, "abc1234") {
		t.Error("Expected result to contain commit hash abc1234")
	}

	if !strings.Contains(result, "def5678") {
		t.Error("Expected result to contain commit hash def5678")
	}
}

// TestGitCommitHistoryService_GetCommitHistoryWithDetails_EmptyHistory は履歴が空の場合のテスト
func TestGitCommitHistoryService_GetCommitHistoryWithDetails_EmptyHistory(t *testing.T) {
	// Arrange
	mockExecutor := &MockGitExecutor{
		ExecuteFunc: func(workingDir string, args ...string) ([]byte, error) {
			if len(args) > 0 && args[0] == "rev-parse" {
				return []byte(".git"), nil
			}
			if len(args) > 0 && args[0] == "log" {
				return []byte(""), nil
			}
			return []byte(""), nil
		},
	}

	cfg := &config.Config{
		GitDir:  "/test/repo",
		Keyword: "",
		Since:   "",
		Until:   "",
	}

	gitClient := git.NewClientWithExecutor("/test/repo", mockExecutor)
	service := &GitCommitHistoryService{
		gitClient: gitClient,
		config:    cfg,
	}

	// Act
	result, err := service.GetCommitHistoryWithDetails()

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := "=== Commit History ===\n指定された条件に一致するコミットが見つかりませんでした。"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestGitCommitHistoryService_GetCommitHistoryWithDetails_InvalidRepo は無効なリポジトリの場合のテスト
func TestGitCommitHistoryService_GetCommitHistoryWithDetails_InvalidRepo(t *testing.T) {
	// Arrange
	mockExecutor := &MockGitExecutor{
		ExecuteFunc: func(workingDir string, args ...string) ([]byte, error) {
			if len(args) > 0 && args[0] == "rev-parse" {
				return nil, fmt.Errorf("not a git repository")
			}
			return []byte(""), nil
		},
	}

	cfg := &config.Config{
		GitDir:  "/invalid/repo",
		Keyword: "",
		Since:   "",
		Until:   "",
	}

	gitClient := git.NewClientWithExecutor("/invalid/repo", mockExecutor)
	service := &GitCommitHistoryService{
		gitClient: gitClient,
		config:    cfg,
	}

	// Act
	_, err := service.GetCommitHistoryWithDetails()

	// Assert
	if err == nil {
		t.Error("Expected error for invalid repository, got nil")
	}

	if !strings.Contains(err.Error(), "有効なGitリポジトリではありません") {
		t.Errorf("Expected error message about invalid git repository, got %v", err)
	}
}

// TestGitCommitHistoryService_GetCommitHistoryWithDetails_GitLogError はgit logエラーの場合のテスト
func TestGitCommitHistoryService_GetCommitHistoryWithDetails_GitLogError(t *testing.T) {
	// Arrange
	mockExecutor := &MockGitExecutor{
		ExecuteFunc: func(workingDir string, args ...string) ([]byte, error) {
			if len(args) > 0 && args[0] == "rev-parse" {
				return []byte(".git"), nil
			}
			if len(args) > 0 && args[0] == "log" {
				return nil, fmt.Errorf("git log failed")
			}
			return []byte(""), nil
		},
	}

	cfg := &config.Config{
		GitDir:  "/test/repo",
		Keyword: "",
		Since:   "",
		Until:   "",
	}

	gitClient := git.NewClientWithExecutor("/test/repo", mockExecutor)
	service := &GitCommitHistoryService{
		gitClient: gitClient,
		config:    cfg,
	}

	// Act
	_, err := service.GetCommitHistoryWithDetails()

	// Assert
	if err == nil {
		t.Error("Expected error for git log failure, got nil")
	}

	if !strings.Contains(err.Error(), "コミット履歴の取得に失敗しました") {
		t.Errorf("Expected error message about commit history failure, got %v", err)
	}
}

// TestGitCommitHistoryService_GetCommitHistoryWithDetails_GitShowError はgit showエラーの場合のテスト
func TestGitCommitHistoryService_GetCommitHistoryWithDetails_GitShowError(t *testing.T) {
	// Arrange
	mockExecutor := &MockGitExecutor{
		ExecuteFunc: func(workingDir string, args ...string) ([]byte, error) {
			if len(args) > 0 && args[0] == "rev-parse" {
				return []byte(".git"), nil
			}
			if len(args) > 0 && args[0] == "log" {
				return []byte("* abc1234 - feat: test commit (1 hour ago) <test@example.com>"), nil
			}
			if len(args) > 0 && args[0] == "show" {
				return nil, fmt.Errorf("git show failed")
			}
			return []byte(""), nil
		},
	}

	cfg := &config.Config{
		GitDir:  "/test/repo",
		Keyword: "",
		Since:   "",
		Until:   "",
	}

	gitClient := git.NewClientWithExecutor("/test/repo", mockExecutor)
	service := &GitCommitHistoryService{
		gitClient: gitClient,
		config:    cfg,
	}

	// Act
	_, err := service.GetCommitHistoryWithDetails()

	// Assert
	if err == nil {
		t.Error("Expected error for git show failure, got nil")
	}

	if !strings.Contains(err.Error(), "コミット詳細の取得に失敗しました") {
		t.Errorf("Expected error message about commit details failure, got %v", err)
	}
}
