package command_executor

// Repository は外部コマンド実行を抽象化します。
type Repository interface {
	Execute(name string, args ...string) ([]byte, error)
}
