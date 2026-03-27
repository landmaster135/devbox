package filesystem

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/landmaster135/devbox/internal/agent_session/domain"
)

const (
	sessionsDirectoryName     = "sessions"
	rolloutFilePrefix         = "rollout-"
	rolloutFileSuffix         = ".jsonl"
	rolloutTimestampLayout    = "2006-01-02T15-04-05"
	jsonLineScannerMaxBytes   = 8 * 1024 * 1024
	maxScanFiles              = 10000
	headRecordLimit           = 10
	userEventScanLimit        = 200
	maxReadableRecordsPerFile = headRecordLimit + userEventScanLimit
)

// SessionRepository はRepositoryの実装。
type SessionRepository struct{}

// NewRepository はRepositoryの実装を作成する。
func NewRepository() Repository {
	return &SessionRepository{}
}

// RetrieveCodexSessions はagentHomeDir配下からCodexセッション一覧を取得する。
func (r *SessionRepository) RetrieveCodexSessions(agentHomeDir string) ([]domain.Session, error) {
	sessionsRoot := filepath.Join(agentHomeDir, sessionsDirectoryName)
	info, err := os.Stat(sessionsRoot)
	if err != nil {
		if os.IsNotExist(err) {
			return []domain.Session{}, nil
		}
		return nil, fmt.Errorf("セッションディレクトリの確認に失敗しました: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("セッションディレクトリが存在しません: %s", sessionsRoot)
	}

	candidates, err := collectRolloutCandidates(sessionsRoot)
	if err != nil {
		return nil, fmt.Errorf("セッションファイルの走査に失敗しました: %w", err)
	}

	sessions := make([]domain.Session, 0, len(candidates))
	for _, candidate := range candidates {
		session, parseErr := r.parseRolloutFile(candidate.path, candidate.createdAt, candidate.uuid)
		if parseErr != nil {
			continue
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

type rolloutCandidate struct {
	path      string
	createdAt time.Time
	uuid      string
}

type rolloutLine struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type summary struct {
	sawSessionMeta bool
	sawUserEvent   bool
}

type sessionMetaPayload struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	CWD       string `json:"cwd"`
	Git       struct {
		Branch string `json:"branch"`
	} `json:"git"`
}

type eventMessagePayload struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Kind    string `json:"kind"`
}

func (r *SessionRepository) parseRolloutFile(path string, fallbackCreatedAt time.Time, fallbackUUID string) (domain.Session, error) {
	file, err := os.Open(path)
	if err != nil {
		return domain.Session{}, fmt.Errorf("ロールアウトファイルのオープンに失敗しました: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return domain.Session{}, fmt.Errorf("ロールアウトファイル情報の取得に失敗しました: %w", err)
	}

	session := domain.Session{
		UpdatedAt: fileInfo.ModTime(),
		Branch:    "-",
		CWD:       "-",
		UUID:      fallbackUUID,
		CreatedAt: fallbackCreatedAt,
	}

	scanner := bufio.NewScanner(file)
	buffer := make([]byte, 0, 64*1024)
	scanner.Buffer(buffer, jsonLineScannerMaxBytes)

	linesScanned := 0
	currentSummary := summary{}
	for scanner.Scan() {
		linesScanned++
		if linesScanned > maxReadableRecordsPerFile {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var parsedLine rolloutLine
		if err := json.Unmarshal([]byte(line), &parsedLine); err != nil {
			continue
		}

		switch parsedLine.Type {
		case "session_meta":
			r.parseSessionMeta(parsedLine.Payload, &session)
			currentSummary.sawSessionMeta = true
		case "event_msg":
			r.parseEventMessage(parsedLine.Payload, &session)
			if session.Conversation != "" {
				currentSummary.sawUserEvent = true
			}
		}

		if currentSummary.sawSessionMeta && currentSummary.sawUserEvent {
			break
		}
		if !currentSummary.sawSessionMeta && linesScanned >= headRecordLimit {
			break
		}
		if currentSummary.sawSessionMeta && !currentSummary.sawUserEvent && linesScanned >= maxReadableRecordsPerFile {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return domain.Session{}, fmt.Errorf("ロールアウトファイルの読み込みに失敗しました: %w", err)
	}

	if !currentSummary.sawSessionMeta || !currentSummary.sawUserEvent {
		return domain.Session{}, fmt.Errorf("セッションメタまたはユーザーメッセージが不足しています: %s", path)
	}
	if session.UUID == "" {
		return domain.Session{}, fmt.Errorf("セッションUUIDが取得できませんでした: %s", path)
	}
	if session.CreatedAt.IsZero() {
		session.CreatedAt = session.UpdatedAt
	}
	if session.Conversation == "" {
		session.Conversation = "-"
	}

	return session, nil
}

func (r *SessionRepository) parseSessionMeta(payload json.RawMessage, session *domain.Session) {
	var meta sessionMetaPayload
	if err := json.Unmarshal(payload, &meta); err != nil {
		return
	}

	if meta.ID != "" {
		session.UUID = meta.ID
	}
	if meta.CWD != "" {
		session.CWD = meta.CWD
	}
	if meta.Git.Branch != "" {
		session.Branch = meta.Git.Branch
	}
	if meta.Timestamp != "" {
		if parsed, err := time.Parse(time.RFC3339Nano, meta.Timestamp); err == nil {
			session.CreatedAt = parsed
		}
	}
}

func (r *SessionRepository) parseEventMessage(payload json.RawMessage, session *domain.Session) {
	if session.Conversation != "" {
		return
	}

	var eventPayload eventMessagePayload
	if err := json.Unmarshal(payload, &eventPayload); err != nil {
		return
	}
	if eventPayload.Type != "user_message" || strings.TrimSpace(eventPayload.Message) == "" {
		return
	}

	normalized := strings.TrimSpace(strings.ReplaceAll(eventPayload.Message, "\n", " "))
	session.Conversation = normalized
}

var uuidSuffixPattern = regexp.MustCompile(`^([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})$`)

func parseTimestampUUIDFromFilename(fileName string) (time.Time, string, bool) {
	if !strings.HasPrefix(fileName, rolloutFilePrefix) || !strings.HasSuffix(fileName, rolloutFileSuffix) {
		return time.Time{}, "", false
	}
	core := strings.TrimSuffix(strings.TrimPrefix(fileName, rolloutFilePrefix), rolloutFileSuffix)

	for idx := len(core) - 1; idx >= 0; idx-- {
		if core[idx] != '-' {
			continue
		}
		uuidCandidate := core[idx+1:]
		if !uuidSuffixPattern.MatchString(uuidCandidate) {
			continue
		}

		tsRaw := core[:idx]
		createdAt, err := time.Parse(rolloutTimestampLayout, tsRaw)
		if err != nil {
			return time.Time{}, "", false
		}
		return createdAt.UTC(), strings.ToLower(uuidCandidate), true
	}
	return time.Time{}, "", false
}

func collectRolloutCandidates(root string) ([]rolloutCandidate, error) {
	scannedFiles := 0
	candidates := make([]rolloutCandidate, 0)

	yearDirs, err := collectDirsDesc(root, 4)
	if err != nil {
		return nil, err
	}

	for _, yearDir := range yearDirs {
		if scannedFiles >= maxScanFiles {
			break
		}
		monthDirs, monthErr := collectDirsDesc(filepath.Join(root, yearDir), 2)
		if monthErr != nil {
			return nil, monthErr
		}
		for _, monthDir := range monthDirs {
			if scannedFiles >= maxScanFiles {
				break
			}
			monthPath := filepath.Join(root, yearDir, monthDir)
			dayDirs, dayErr := collectDirsDesc(monthPath, 2)
			if dayErr != nil {
				return nil, dayErr
			}
			for _, dayDir := range dayDirs {
				if scannedFiles >= maxScanFiles {
					break
				}
				dayPath := filepath.Join(monthPath, dayDir)
				entries, readErr := os.ReadDir(dayPath)
				if readErr != nil {
					return nil, readErr
				}
				dayCandidates := make([]rolloutCandidate, 0)
				for _, entry := range entries {
					if scannedFiles >= maxScanFiles {
						break
					}
					if entry.IsDir() {
						continue
					}
					createdAt, uuid, ok := parseTimestampUUIDFromFilename(entry.Name())
					if !ok {
						continue
					}
					scannedFiles++
					dayCandidates = append(dayCandidates, rolloutCandidate{
						path:      filepath.Join(dayPath, entry.Name()),
						createdAt: createdAt,
						uuid:      uuid,
					})
				}
				slices.SortFunc(dayCandidates, func(a, b rolloutCandidate) int {
					if a.createdAt.After(b.createdAt) {
						return -1
					}
					if a.createdAt.Before(b.createdAt) {
						return 1
					}
					return strings.Compare(b.uuid, a.uuid)
				})
				candidates = append(candidates, dayCandidates...)
			}
		}
	}

	return candidates, nil
}

func collectDirsDesc(parentPath string, width int) ([]string, error) {
	entries, err := os.ReadDir(parentPath)
	if err != nil {
		return nil, err
	}

	dirs := make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if len(name) != width {
			continue
		}
		if _, parseErr := strconv.Atoi(name); parseErr != nil {
			continue
		}
		dirs = append(dirs, name)
	}

	slices.SortFunc(dirs, func(a, b string) int {
		return strings.Compare(b, a)
	})
	return dirs, nil
}
