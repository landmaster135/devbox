package git

import (
	"os/exec"
)

// GitCommandExecutor はGitコマンドを実行するインターフェース
type GitCommandExecutor interface {
	Execute(workingDir string, args ...string) ([]byte, error)
}

// StandardGitExecutor は標準のGitコマンド実行器
type StandardGitExecutor struct{}

// NewStandardGitExecutor は新しいStandardGitExecutorを作成する
func NewStandardGitExecutor() *StandardGitExecutor {
	return &StandardGitExecutor{}
}

// Execute はGitコマンドを実行する
func (e *StandardGitExecutor) Execute(workingDir string, args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = workingDir
	return cmd.Output()
}
