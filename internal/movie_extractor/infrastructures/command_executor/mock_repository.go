package command_executor

// CommandCall は Execute の呼び出し情報です。
type CommandCall struct {
	Name string
	Args []string
}

// MockRepository は Repository のモック実装です。
type MockRepository struct {
	ExecuteFunc func(name string, args ...string) ([]byte, error)
	Calls       []CommandCall
}

// Execute は呼び出しを記録して ExecuteFunc を呼び出します。
func (m *MockRepository) Execute(name string, args ...string) ([]byte, error) {
	m.Calls = append(m.Calls, CommandCall{
		Name: name,
		Args: append([]string(nil), args...),
	})

	if m.ExecuteFunc == nil {
		return []byte{}, nil
	}
	return m.ExecuteFunc(name, args...)
}
