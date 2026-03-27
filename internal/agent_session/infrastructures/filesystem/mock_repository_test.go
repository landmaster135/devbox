package filesystem

import "testing"

type TestMockRepository struct{}

func TestMockRepository_RetrieveCodexSessions_DefaultReturn_Normal(t *testing.T) {
	t.Parallel()

	mock := &MockRepository{}
	sessions, err := mock.RetrieveCodexSessions("/tmp/.codex")
	if err != nil {
		t.Fatalf("RetrieveCodexSessions() error = %v", err)
	}
	if len(sessions) != 0 {
		t.Fatalf("RetrieveCodexSessions() len = %d, want 0", len(sessions))
	}
}
