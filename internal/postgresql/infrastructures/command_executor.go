package infrastructures

import "os/exec"

// CommandExecutor は外部コマンド実行を抽象化します。
type CommandExecutor interface {
	Execute(name string, args ...string) ([]byte, error)
}

// OSCommandExecutor は OS コマンドを実行する実装です。
type OSCommandExecutor struct{}

// NewOSCommandExecutor は OSCommandExecutor を返します。
func NewOSCommandExecutor() *OSCommandExecutor {
	return &OSCommandExecutor{}
}

// Execute はシェルを介さずコマンドを実行します。
func (e *OSCommandExecutor) Execute(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.CombinedOutput()
}
