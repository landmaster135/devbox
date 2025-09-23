package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	security "github.com/landmaster135/devbox/internal/git_commit_history_retriever/security"
)

func createTempDirWithinWorkspace(t *testing.T, prefix string) string {
	t.Helper()

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	baseDir := filepath.Join(cwd, "test-temp")
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		t.Fatalf("failed to prepare base test directory: %v", err)
	}

	tempDir, err := os.MkdirTemp(baseDir, prefix)
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}

	t.Cleanup(func() { _ = os.RemoveAll(tempDir) })
	return tempDir
}

func createTempGitRepo(t *testing.T, prefix string) string {
	t.Helper()

	tempDir := createTempDirWithinWorkspace(t, prefix)

	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Skipf("git not available or failed to init: %v", err)
	}

	return tempDir
}

func TestNewStandardGitExecutor(t *testing.T) {
	executor := NewStandardGitExecutor()
	if executor == nil {
		t.Fatal("NewStandardGitExecutor returned nil")
	}
	if executor.pathValidator == nil {
		t.Error("pathValidator should not be nil")
	}

	var _ GitCommandExecutor = executor
}

func TestNewStandardGitExecutorWithValidator(t *testing.T) {
	validator := security.NewPathValidator([]string{"/custom/base"}, 2048)
	executor := NewStandardGitExecutorWithValidator(validator)
	if executor == nil {
		t.Fatal("NewStandardGitExecutorWithValidator returned nil")
	}
	if executor.pathValidator != validator {
		t.Error("pathValidator was not set from injected validator")
	}
}

func TestStandardGitExecutor_Execute_Normal(t *testing.T) {
	executor := NewStandardGitExecutor()
	repoDir := createTempGitRepo(t, "git-test-normal-")

	output, err := executor.Execute(repoDir, "status", "--porcelain")
	if err != nil {
		if strings.Contains(err.Error(), "無効なワーキングディレクトリ") {
			t.Fatalf("unexpected path validation error: %v", err)
		}
		t.Fatalf("expected no error, got %v", err)
	}

	if len(output) != 0 {
		t.Logf("status output: %s", string(output))
	}
}

func TestStandardGitExecutor_Execute_InvalidCommand(t *testing.T) {
	executor := NewStandardGitExecutor()
	repoDir := createTempGitRepo(t, "git-test-invalid-command-")

	_, err := executor.Execute(repoDir, "invalid-command")
	if err == nil {
		t.Fatal("expected error for invalid command, got nil")
	}
	if strings.Contains(err.Error(), "無効なワーキングディレクトリ") {
		t.Fatalf("received path validation error for valid repo: %v", err)
	}
}

func TestStandardGitExecutor_Execute_InvalidDirectory(t *testing.T) {
	executor := NewStandardGitExecutor()

	_, err := executor.Execute("/non/existent/directory", "status")
	if err == nil {
		t.Fatal("expected error for invalid directory, got nil")
	}
	if !strings.Contains(err.Error(), "無効なワーキングディレクトリ") {
		t.Fatalf("expected path validation error, got %v", err)
	}
}

func TestStandardGitExecutor_Execute_NonGitDirectory(t *testing.T) {
	executor := NewStandardGitExecutor()
	dir := createTempDirWithinWorkspace(t, "git-test-non-git-")

	_, err := executor.Execute(dir, "status")
	if err == nil {
		t.Fatal("expected error for non-git directory, got nil")
	}
	if !strings.Contains(err.Error(), "無効なワーキングディレクトリ") {
		t.Fatalf("expected path validation error, got %v", err)
	}
}

func TestStandardGitExecutor_Execute_EmptyArgs(t *testing.T) {
	executor := NewStandardGitExecutor()
	repoDir := createTempGitRepo(t, "git-test-empty-args-")

	_, err := executor.Execute(repoDir)
	if err != nil && strings.Contains(err.Error(), "無効なワーキングディレクトリ") {
		t.Fatalf("unexpected path validation error for repo: %v", err)
	}
}

func TestStandardGitExecutor_Execute_MultipleArgs(t *testing.T) {
	executor := NewStandardGitExecutor()
	repoDir := createTempGitRepo(t, "git-test-multi-args-")

	_, err := executor.Execute(repoDir, "log", "--oneline", "--max-count=1")
	if err != nil && strings.Contains(err.Error(), "無効なワーキングディレクトリ") {
		t.Fatalf("unexpected path validation error for repo: %v", err)
	}
}

func TestStandardGitExecutor_Execute_GitVersion(t *testing.T) {
	executor := NewStandardGitExecutor()
	repoDir := createTempGitRepo(t, "git-test-version-")

	output, err := executor.Execute(repoDir, "version")
	if err != nil {
		t.Skipf("git version command not available: %v", err)
	}

	if !strings.Contains(string(output), "git version") {
		t.Fatalf("expected output to contain 'git version', got %s", string(output))
	}
}

func TestStandardGitExecutor_Execute_Help(t *testing.T) {
	executor := NewStandardGitExecutor()
	repoDir := createTempGitRepo(t, "git-test-help-")

	output, err := executor.Execute(repoDir, "help", "-a")
	if err != nil {
		t.Skipf("git help not available: %v", err)
	}

	if len(output) == 0 {
		t.Fatal("expected non-empty output from git help")
	}
}

func TestStandardGitExecutor_Execute_Integration(t *testing.T) {
	executor := NewStandardGitExecutor()
	repoDir := createTempGitRepo(t, "git-test-integration-")

	if _, err := executor.Execute(repoDir, "config", "user.name", "Test User"); err != nil {
		t.Fatalf("failed to set user.name: %v", err)
	}
	if _, err := executor.Execute(repoDir, "config", "user.email", "test@example.com"); err != nil {
		t.Fatalf("failed to set user.email: %v", err)
	}

	testFile := filepath.Join(repoDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	if _, err := executor.Execute(repoDir, "add", "test.txt"); err != nil {
		t.Fatalf("failed to add file: %v", err)
	}

	if _, err := executor.Execute(repoDir, "commit", "-m", "Initial commit"); err != nil {
		t.Fatalf("failed to commit: %v", err)
	}

	output, err := executor.Execute(repoDir, "log", "--oneline")
	if err != nil {
		t.Fatalf("failed to get log: %v", err)
	}

	if !strings.Contains(string(output), "Initial commit") {
		t.Fatalf("expected log to contain 'Initial commit', got %s", string(output))
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
		t.Run(path, func(t *testing.T) {
			_, err := executor.Execute(path, "status")
			if err == nil {
				t.Fatalf("expected error for dangerous path %s, got nil", path)
			}
			if !strings.Contains(err.Error(), "無効なワーキングディレクトリ") {
				t.Fatalf("expected path validation error, got %v", err)
			}
		})
	}
}
