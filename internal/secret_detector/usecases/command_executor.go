package usecases

import (
	"os/exec"
)

// CommandExecutorRepository はコマンド実行のインターフェース
type CommandExecutorRepository interface {
	Execute(name string, args ...string) ([]byte, error)
}

// CommandExecutor は実際のexec.Commandを使用する実装
type CommandExecutor struct{}

// NewCommandExecutor は新しいCommandExecutorを作成
func NewCommandExecutor() *CommandExecutor {
	return &CommandExecutor{}
}

// Execute はコマンドを実行し、出力を返す
func (r *CommandExecutor) Execute(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.Output()
}

// MockCommandExecutor はテスト用のモック実装
type MockCommandExecutor struct {
	ExecuteFunc func(name string, args ...string) ([]byte, error)
}

// Execute はモック関数を実行
func (m *MockCommandExecutor) Execute(name string, args ...string) ([]byte, error) {
	return m.ExecuteFunc(name, args...)
}
