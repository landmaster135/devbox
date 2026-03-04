package infrastructures

// CommandCall は Execute の呼び出し引数を保持する。
type CommandCall struct {
	Name string
	Args []string
}

// MockCommandExecutor は CommandExecutor のモック実装。
type MockCommandExecutor struct {
	ExecuteFunc func(name string, args ...string) ([]byte, error)
	Calls       []CommandCall
}

// Execute は記録したうえで ExecuteFunc を呼び出す。
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
