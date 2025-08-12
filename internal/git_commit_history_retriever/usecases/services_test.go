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

// TestNewGitCommitHistoryService_Normal はNewGitCommitHistoryServiceの正常系テスト
func TestNewGitCommitHistoryService_Normal(t *testing.T) {
	// Arrange
	cfg := &config.Config{
		GitDir:  "/test/repo",
		Keyword: "feat:",
		Since:   "2025-01-01",
		Until:   "2025-01-31",
	}

	// Act
	service := NewGitCommitHistoryService("/test/repo", cfg)

	// Assert
	if service == nil {
		t.Error("Expected service to be non-nil")
		return
	}

	if service.config != cfg {
		t.Error("Expected config to be the provided config")
	}

	if service.gitClient == nil {
		t.Error("Expected gitClient to be non-nil")
	}
}

// TestNewGitCommitHistoryService_EmptyGitDir は空のGitDirの場合のテスト
func TestNewGitCommitHistoryService_EmptyGitDir(t *testing.T) {
	// Arrange
	cfg := &config.Config{
		GitDir:  "",
		Keyword: "",
		Since:   "",
		Until:   "",
	}

	// Act
	service := NewGitCommitHistoryService("", cfg)

	// Assert
	if service == nil {
		t.Error("Expected service to be non-nil")
		return
	}

	if service.config != cfg {
		t.Error("Expected config to be the provided config")
	}

	if service.gitClient == nil {
		t.Error("Expected gitClient to be non-nil")
	}
}

// TestNewGitCommitHistoryService_NilConfig はnilConfigの場合のテスト
func TestNewGitCommitHistoryService_NilConfig(t *testing.T) {
	// Act
	service := NewGitCommitHistoryService("/test/repo", nil)

	// Assert
	if service == nil {
		t.Error("Expected service to be non-nil")
		return
	}

	if service.config != nil {
		t.Error("Expected config to be nil")
	}

	if service.gitClient == nil {
		t.Error("Expected gitClient to be non-nil")
	}
}

// TestGitCommitHistoryService_GetCommitHistory_WithKeyword はキーワード付きGetCommitHistoryのテスト
func TestGitCommitHistoryService_GetCommitHistory_WithKeyword(t *testing.T) {
	// Arrange
	mockExecutor := &MockGitExecutor{
		ExecuteFunc: func(workingDir string, args ...string) ([]byte, error) {
			if len(args) > 0 && args[0] == "rev-parse" {
				return []byte(".git"), nil
			}
			if len(args) > 0 && args[0] == "log" {
				// キーワード検索が含まれていることを確認
				for _, arg := range args {
					if strings.Contains(arg, "--grep=feat:") {
						return []byte("* abc1234 - feat: test commit (1 hour ago) <test@example.com>"), nil
					}
				}
			}
			return []byte(""), nil
		},
	}

	cfg := &config.Config{
		GitDir:  "/test/repo",
		Keyword: "feat:",
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

	if !strings.Contains(result, "abc1234") {
		t.Error("Expected result to contain commit with keyword")
	}

	if !strings.Contains(result, "=== Commit History ===") {
		t.Error("Expected result to contain header")
	}
}

// TestGitCommitHistoryService_GetCommitHistory_WithDateRange は日付範囲付きGetCommitHistoryのテスト
func TestGitCommitHistoryService_GetCommitHistory_WithDateRange(t *testing.T) {
	// Arrange
	mockExecutor := &MockGitExecutor{
		ExecuteFunc: func(workingDir string, args ...string) ([]byte, error) {
			if len(args) > 0 && args[0] == "rev-parse" {
				return []byte(".git"), nil
			}
			if len(args) > 0 && args[0] == "log" {
				// 日付範囲が含まれていることを確認
				hasSince := false
				hasUntil := false
				for _, arg := range args {
					if strings.Contains(arg, "--since=2025-01-01") {
						hasSince = true
					}
					if strings.Contains(arg, "--until=2025-01-31") {
						hasUntil = true
					}
				}
				if hasSince && hasUntil {
					return []byte("* abc1234 - feat: test commit (1 hour ago) <test@example.com>"), nil
				}
			}
			return []byte(""), nil
		},
	}

	cfg := &config.Config{
		GitDir:  "/test/repo",
		Keyword: "",
		Since:   "2025-01-01",
		Until:   "2025-01-31",
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

	if !strings.Contains(result, "abc1234") {
		t.Error("Expected result to contain commit within date range")
	}

	if !strings.Contains(result, "=== Commit History ===") {
		t.Error("Expected result to contain header")
	}
}

// TestGitCommitHistoryService_GetCommitHistoryWithDetails_WithKeywordAndDateRange はキーワードと日付範囲付きGetCommitHistoryWithDetailsのテスト
func TestGitCommitHistoryService_GetCommitHistoryWithDetails_WithKeywordAndDateRange(t *testing.T) {
	// Arrange
	mockExecutor := &MockGitExecutor{
		ExecuteFunc: func(workingDir string, args ...string) ([]byte, error) {
			if len(args) > 0 && args[0] == "rev-parse" {
				return []byte(".git"), nil
			}
			if len(args) > 0 && args[0] == "log" {
				// キーワードと日付範囲が含まれていることを確認
				hasKeyword := false
				hasSince := false
				hasUntil := false
				for _, arg := range args {
					if strings.Contains(arg, "--grep=feat:") {
						hasKeyword = true
					}
					if strings.Contains(arg, "--since=2025-01-01") {
						hasSince = true
					}
					if strings.Contains(arg, "--until=2025-01-31") {
						hasUntil = true
					}
				}
				if hasKeyword && hasSince && hasUntil {
					return []byte("* abc1234 - feat: test commit (1 hour ago) <test@example.com>"), nil
				}
			}
			if len(args) > 0 && args[0] == "show" {
				return []byte("commit abc1234\nAuthor: test@example.com\nDate: test date\n\n    feat: test commit\n\n file1.go | 10 ++++++++++\n 1 file changed, 10 insertions(+)"), nil
			}
			return []byte(""), nil
		},
	}

	cfg := &config.Config{
		GitDir:  "/test/repo",
		Keyword: "feat:",
		Since:   "2025-01-01",
		Until:   "2025-01-31",
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
		t.Error("Expected result to contain commit hash")
	}

	if !strings.Contains(result, "feat: test commit") {
		t.Error("Expected result to contain commit message")
	}
}

// TestGitCommitHistoryService_GetCommitHistoryWithDetails_NoCommitDetails はコミット詳細がない場合のテスト
func TestGitCommitHistoryService_GetCommitHistoryWithDetails_NoCommitDetails(t *testing.T) {
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
				return []byte(""), nil // 空の詳細
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

	// 詳細が空の場合は詳細ヘッダーが含まれないことを確認
	if strings.Contains(result, "=== Commit Details ===") {
		t.Error("Expected result to not contain commit details header when details are empty")
	}

	if !strings.Contains(result, "abc1234") {
		t.Error("Expected result to contain commit hash")
	}
}
