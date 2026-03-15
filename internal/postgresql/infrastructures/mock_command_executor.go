package infrastructures

// CommandCall は Execute の呼び出し引数を保持します。
type CommandCall struct {
	Name string
	Args []string
}

// MockCommandExecutor は CommandExecutor のモック実装です。
type MockCommandExecutor struct {
	ExecuteFunc func(name string, args ...string) ([]byte, error)
	Calls       []CommandCall
}

// Execute は呼び出しを記録して ExecuteFunc を実行します。
func (m *MockCommandExecutor) Execute(name string, args ...string) ([]byte, error) {
	m.Calls = append(m.Calls, CommandCall{
		Name: name,
		Args: append([]string(nil), args...),
	})

	if m.ExecuteFunc == nil {
		return []byte{}, nil
	}
	return m.ExecuteFunc(name, args...)
}
