package usecases

import (
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/landmaster135/devbox/internal/agent_session/domain"
	filesystem "github.com/landmaster135/devbox/internal/agent_session/infrastructures/filesystem"
)

const (
	codexAgentType            = "codex"
	conversationMaxLengthRune = 120
)

// RetrieveSessionsInput はセッション取得時の入力値を表す。
type RetrieveSessionsInput struct {
	AgentType    string
	Limit        int
	StartDate    *time.Time
	EndDate      *time.Time
	AgentHomeDir string
}

// AgentSessionService はエージェントセッション一覧取得のユースケースを提供する。
type AgentSessionService struct {
	repository filesystem.Repository
	nowFunc    func() time.Time
}

// NewAgentSessionService はAgentSessionServiceを生成する。
func NewAgentSessionService() *AgentSessionService {
	return NewAgentSessionServiceWithDependencies(filesystem.NewRepository(), time.Now)
}

// NewAgentSessionServiceWithDependencies は依存性注入付きでAgentSessionServiceを生成する。
func NewAgentSessionServiceWithDependencies(repository filesystem.Repository, nowFunc func() time.Time) *AgentSessionService {
	if repository == nil {
		repository = filesystem.NewRepository()
	}
	if nowFunc == nil {
		nowFunc = time.Now
	}

	return &AgentSessionService{
		repository: repository,
		nowFunc:    nowFunc,
	}
}

// RetrieveSessions は条件に一致するセッション一覧を取得し、表示文字列を返す。
func (s *AgentSessionService) RetrieveSessions(input RetrieveSessionsInput) (string, error) {
	if input.AgentType != codexAgentType {
		return "", fmt.Errorf("未対応のagent-typeです: %s", input.AgentType)
	}
	if input.Limit <= 0 {
		return "", fmt.Errorf("limitは1以上を指定してください")
	}

	sessions, err := s.repository.RetrieveCodexSessions(input.AgentHomeDir)
	if err != nil {
		return "", err
	}

	filteredSessions := s.filterByDateRange(sessions, input.StartDate, input.EndDate)
	sort.Slice(filteredSessions, func(i, j int) bool {
		if filteredSessions[i].UpdatedAt.Equal(filteredSessions[j].UpdatedAt) {
			return filteredSessions[i].CreatedAt.After(filteredSessions[j].CreatedAt)
		}
		return filteredSessions[i].UpdatedAt.After(filteredSessions[j].UpdatedAt)
	})

	if len(filteredSessions) > input.Limit {
		filteredSessions = filteredSessions[:input.Limit]
	}

	return s.formatSessions(filteredSessions), nil
}

func (s *AgentSessionService) filterByDateRange(sessions []domain.Session, startDate, endDate *time.Time) []domain.Session {
	if startDate == nil && endDate == nil {
		return sessions
	}

	filtered := make([]domain.Session, 0, len(sessions))
	for _, session := range sessions {
		targetDate := session.CreatedAt
		if targetDate.IsZero() {
			targetDate = session.UpdatedAt
		}
		normalizedDate := normalizeDate(targetDate)

		if startDate != nil && normalizedDate.Before(normalizeDate(*startDate)) {
			continue
		}
		if endDate != nil && normalizedDate.After(normalizeDate(*endDate)) {
			continue
		}
		filtered = append(filtered, session)
	}
	return filtered
}

func (s *AgentSessionService) formatSessions(sessions []domain.Session) string {
	if len(sessions) == 0 {
		return "対象セッションは見つかりませんでした。\n"
	}

	now := s.nowFunc()
	var builder strings.Builder
	writer := tabwriter.NewWriter(&builder, 0, 0, 2, ' ', 0)

	fmt.Fprintln(writer, "UUID\tCreated at\tUpdated at\tBranch\tCWD\tConversation")
	for _, session := range sessions {
		createdAtText := formatRelativeTime(now, session.CreatedAt)
		updatedAtText := formatRelativeTime(now, session.UpdatedAt)
		branchText := normalizeForTable(session.Branch)
		cwdText := normalizeForTable(session.CWD)
		conversationText := truncateRunes(normalizeForTable(compactWhitespace(session.Conversation)), conversationMaxLengthRune)

		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\n",
			normalizeForTable(session.UUID),
			createdAtText,
			updatedAtText,
			branchText,
			cwdText,
			conversationText,
		)
	}
	_ = writer.Flush()

	return builder.String()
}

func normalizeDate(value time.Time) time.Time {
	local := value.In(time.Local)
	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, time.Local)
}

func normalizeForTable(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "-"
	}
	return trimmed
}

func compactWhitespace(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

func truncateRunes(value string, limit int) string {
	if limit <= 0 {
		return ""
	}
	runes := []rune(value)
	if len(runes) <= limit {
		return value
	}
	if limit <= 3 {
		return string(runes[:limit])
	}
	return string(runes[:limit-3]) + "..."
}

func formatRelativeTime(now time.Time, target time.Time) string {
	if target.IsZero() {
		return "-"
	}

	delta := now.Sub(target)
	if delta < 0 {
		delta = -delta
	}

	if delta < time.Minute {
		return "just now"
	}
	minutes := int(delta / time.Minute)
	if delta < time.Hour {
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	}
	hours := int(delta / time.Hour)
	if delta < 24*time.Hour {
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	}
	days := int(delta / (24 * time.Hour))
	if days == 1 {
		return "1 day ago"
	}
	if days <= 30 {
		return fmt.Sprintf("%d days ago", days)
	}
	return target.In(time.Local).Format("2006-01-02 15:04")
}
