package filesystem

import "github.com/landmaster135/devbox/internal/agent_session/domain"

// MockRepository はRepositoryのモック実装。
type MockRepository struct {
	RetrieveCodexSessionsFunc func(agentHomeDir string) ([]domain.Session, error)
}

// RetrieveCodexSessions はモック関数を呼び出す。
func (m *MockRepository) RetrieveCodexSessions(agentHomeDir string) ([]domain.Session, error) {
	if m.RetrieveCodexSessionsFunc != nil {
		return m.RetrieveCodexSessionsFunc(agentHomeDir)
	}
	return []domain.Session{}, nil
}
