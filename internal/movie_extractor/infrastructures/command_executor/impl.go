package command_executor

import "os/exec"

// OSCommandExecutor は OS コマンド実行の実装です。
type OSCommandExecutor struct{}

// NewOSCommandExecutor は OSCommandExecutor を生成します。
func NewOSCommandExecutor() *OSCommandExecutor {
	return &OSCommandExecutor{}
}

// Execute はシェルを介さずにコマンドを実行します。
func (e *OSCommandExecutor) Execute(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.CombinedOutput()
}
