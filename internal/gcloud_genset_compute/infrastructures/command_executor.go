package infrastructures

import (
	"os/exec"
)

// CommandExecutor は外部コマンド実行を抽象化する。
type CommandExecutor interface {
	Execute(name string, args ...string) ([]byte, error)
}

// OSCommandExecutor は OS コマンドを実行する。
type OSCommandExecutor struct{}

// NewOSCommandExecutor は OSCommandExecutor を返す。
func NewOSCommandExecutor() *OSCommandExecutor {
	return &OSCommandExecutor{}
}

// Execute はシェルを介さずにコマンドを実行する。
func (e *OSCommandExecutor) Execute(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.CombinedOutput()
}
