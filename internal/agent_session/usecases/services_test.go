package usecases

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/landmaster135/devbox/internal/agent_session/domain"
	filesystem "github.com/landmaster135/devbox/internal/agent_session/infrastructures/filesystem"
)

type TestAgentSessionService struct{}

func TestNewAgentSessionService_Normal(t *testing.T) {
	t.Parallel()

	service := NewAgentSessionService()
	if service == nil {
		t.Fatal("NewAgentSessionService() returned nil")
	}
}

func TestAgentSessionService_RetrieveSessions_Normal(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 3, 28, 12, 0, 0, 0, time.Local)
	repository := &filesystem.MockRepository{
		RetrieveCodexSessionsFunc: func(agentHomeDir string) ([]domain.Session, error) {
			return []domain.Session{
				{
					UUID:         "11111111-2222-4333-8444-555555555555",
					CreatedAt:    now.Add(-2 * time.Hour),
					UpdatedAt:    now.Add(-1 * time.Hour),
					Branch:       "task",
					CWD:          "/workspace/devbox",
					Conversation: "Codexのセッション一覧ってどうやって取得出来るの？",
				},
				{
					UUID:         "aaaaaaaa-bbbb-4ccc-8ddd-eeeeeeeeeeee",
					CreatedAt:    now.Add(-30 * time.Hour),
					UpdatedAt:    now.Add(-20 * time.Hour),
					Branch:       "feature",
					CWD:          "/workspace/dotfiles",
					Conversation: "Create Git commit message in English.",
				},
			}, nil
		},
	}

	service := NewAgentSessionServiceWithDependencies(repository, func() time.Time { return now })
	result, err := service.RetrieveSessions(RetrieveSessionsInput{
		AgentType:    "codex",
		Limit:        10,
		AgentHomeDir: "/tmp/codex-home",
	})
	if err != nil {
		t.Fatalf("RetrieveSessions() error = %v", err)
	}

	if !strings.Contains(result, "UUID") {
		t.Fatalf("result does not contain header: %s", result)
	}
	if !strings.Contains(result, "11111111-2222-4333-8444-555555555555") {
		t.Fatalf("result does not contain expected session: %s", result)
	}
	if !strings.Contains(result, "task") {
		t.Fatalf("result does not contain branch: %s", result)
	}
}

func TestAgentSessionService_RetrieveSessions_FilterByDate_Normal(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 3, 28, 12, 0, 0, 0, time.Local)
	repository := &filesystem.MockRepository{
		RetrieveCodexSessionsFunc: func(agentHomeDir string) ([]domain.Session, error) {
			return []domain.Session{
				{
					UUID:         "11111111-2222-4333-8444-555555555555",
					CreatedAt:    time.Date(2026, 3, 10, 10, 0, 0, 0, time.Local),
					UpdatedAt:    now,
					Conversation: "in range",
				},
				{
					UUID:         "aaaaaaaa-bbbb-4ccc-8ddd-eeeeeeeeeeee",
					CreatedAt:    time.Date(2026, 2, 28, 10, 0, 0, 0, time.Local),
					UpdatedAt:    now.Add(-1 * time.Hour),
					Conversation: "out of range",
				},
			}, nil
		},
	}

	startDate := time.Date(2026, 3, 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(2026, 3, 31, 0, 0, 0, 0, time.Local)

	service := NewAgentSessionServiceWithDependencies(repository, func() time.Time { return now })
	result, err := service.RetrieveSessions(RetrieveSessionsInput{
		AgentType:    "codex",
		Limit:        10,
		StartDate:    &startDate,
		EndDate:      &endDate,
		AgentHomeDir: "/tmp/codex-home",
	})
	if err != nil {
		t.Fatalf("RetrieveSessions() error = %v", err)
	}

	if !strings.Contains(result, "11111111-2222-4333-8444-555555555555") {
		t.Fatalf("result should include in-range session: %s", result)
	}
	if strings.Contains(result, "aaaaaaaa-bbbb-4ccc-8ddd-eeeeeeeeeeee") {
		t.Fatalf("result should not include out-of-range session: %s", result)
	}
}

func TestAgentSessionService_RetrieveSessions_RepositoryError(t *testing.T) {
	t.Parallel()

	repository := &filesystem.MockRepository{
		RetrieveCodexSessionsFunc: func(agentHomeDir string) ([]domain.Session, error) {
			return nil, errors.New("repository failed")
		},
	}

	service := NewAgentSessionServiceWithDependencies(repository, time.Now)
	_, err := service.RetrieveSessions(RetrieveSessionsInput{
		AgentType:    "codex",
		Limit:        10,
		AgentHomeDir: "/tmp/codex-home",
	})
	if err == nil {
		t.Fatal("RetrieveSessions() error = nil")
	}
}

func TestAgentSessionService_RetrieveSessions_InvalidInput(t *testing.T) {
	t.Parallel()

	repository := &filesystem.MockRepository{}
	service := NewAgentSessionServiceWithDependencies(repository, time.Now)

	_, err := service.RetrieveSessions(RetrieveSessionsInput{
		AgentType:    "claude",
		Limit:        10,
		AgentHomeDir: "/tmp/codex-home",
	})
	if err == nil {
		t.Fatal("RetrieveSessions() error = nil")
	}

	_, err = service.RetrieveSessions(RetrieveSessionsInput{
		AgentType:    "codex",
		Limit:        0,
		AgentHomeDir: "/tmp/codex-home",
	})
	if err == nil {
		t.Fatal("RetrieveSessions() error = nil")
	}
}

func TestTruncateRunes_Normal(t *testing.T) {
	t.Parallel()

	if got := truncateRunes("abc", 3); got != "abc" {
		t.Fatalf("truncateRunes() = %s", got)
	}
	if got := truncateRunes("abcdef", 5); got != "ab..." {
		t.Fatalf("truncateRunes() = %s", got)
	}
	if got := truncateRunes("abcdef", 0); got != "" {
		t.Fatalf("truncateRunes() = %s", got)
	}
}

func TestFormatRelativeTime_Normal(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 3, 28, 12, 0, 0, 0, time.Local)
	testCases := []struct {
		name     string
		target   time.Time
		expected string
	}{
		{name: "justNow", target: now.Add(-10 * time.Second), expected: "just now"},
		{name: "minutes", target: now.Add(-5 * time.Minute), expected: "5 minutes ago"},
		{name: "hours", target: now.Add(-2 * time.Hour), expected: "2 hours ago"},
		{name: "days", target: now.Add(-48 * time.Hour), expected: "2 days ago"},
		{name: "olderThanMonth", target: now.Add(-40 * 24 * time.Hour), expected: "2026-02-16 12:00"},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			got := formatRelativeTime(now, tc.target)
			if got != tc.expected {
				t.Fatalf("formatRelativeTime() = %s, want %s", got, tc.expected)
			}
		})
	}
}
