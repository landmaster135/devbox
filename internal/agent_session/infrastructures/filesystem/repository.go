package filesystem

import "github.com/landmaster135/devbox/internal/agent_session/domain"

// Repository はCodexセッション一覧の取得を担う。
type Repository interface {
	RetrieveCodexSessions(agentHomeDir string) ([]domain.Session, error)
}
