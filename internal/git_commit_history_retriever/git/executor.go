package git

import (
	"fmt"
	"os/exec"

	security "github.com/landmaster135/devbox/internal/git_commit_history_retriever/security"
)

// GitCommandExecutor はGitコマンドを実行するインターフェース
type GitCommandExecutor interface {
	Execute(workingDir string, args ...string) ([]byte, error)
}

// StandardGitExecutor は標準のGitコマンド実行器
type StandardGitExecutor struct {
	pathValidator *security.PathValidator
}

// NewStandardGitExecutor は新しいStandardGitExecutorを作成する
func NewStandardGitExecutor() *StandardGitExecutor {
	return &StandardGitExecutor{
		pathValidator: security.NewDefaultPathValidator(),
	}
}

// NewStandardGitExecutorWithValidator は依存性注入可能なStandardGitExecutorを作成する
func NewStandardGitExecutorWithValidator(validator *security.PathValidator) *StandardGitExecutor {
	return &StandardGitExecutor{
		pathValidator: validator,
	}
}

// Execute はGitコマンドを実行する
func (e *StandardGitExecutor) Execute(workingDir string, args ...string) ([]byte, error) {
	validatedPath, err := e.pathValidator.ValidateWorkingDirectory(workingDir)
	if err != nil {
		return nil, fmt.Errorf("無効なワーキングディレクトリ: %w", err)
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = validatedPath
	return cmd.Output()
}
