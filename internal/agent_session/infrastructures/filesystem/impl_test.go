package filesystem

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/landmaster135/devbox/internal/agent_session/domain"
)

type TestSessionRepository struct{}

func TestSessionRepository_RetrieveCodexSessions_Normal(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	sessionsDir := filepath.Join(tempDir, "sessions", "2026", "03", "28")
	if err := os.MkdirAll(sessionsDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	firstPath := filepath.Join(sessionsDir, "rollout-2026-03-28T01-02-03-11111111-2222-4333-8444-555555555555.jsonl")
	firstContent := strings.Join([]string{
		`{"type":"session_meta","payload":{"id":"11111111-2222-4333-8444-555555555555","timestamp":"2026-03-28T01:02:03Z","cwd":"/workspace/devbox","git":{"branch":"task"}}}`,
		`{"type":"event_msg","payload":{"type":"user_message","message":"Codexのセッション一覧を取りたい","kind":"plain"}}`,
	}, "\n")
	if err := os.WriteFile(firstPath, []byte(firstContent), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	updatedAt := time.Date(2026, 3, 28, 2, 0, 0, 0, time.Local)
	if err := os.Chtimes(firstPath, updatedAt, updatedAt); err != nil {
		t.Fatalf("Chtimes() error = %v", err)
	}

	repository := NewRepository()
	sessions, err := repository.RetrieveCodexSessions(tempDir)
	if err != nil {
		t.Fatalf("RetrieveCodexSessions() error = %v", err)
	}
	if len(sessions) != 1 {
		t.Fatalf("RetrieveCodexSessions() len = %d, want 1", len(sessions))
	}

	got := sessions[0]
	if got.UUID != "11111111-2222-4333-8444-555555555555" {
		t.Fatalf("UUID = %s", got.UUID)
	}
	if got.Branch != "task" {
		t.Fatalf("Branch = %s", got.Branch)
	}
	if got.CWD != "/workspace/devbox" {
		t.Fatalf("CWD = %s", got.CWD)
	}
	if got.Conversation != "Codexのセッション一覧を取りたい" {
		t.Fatalf("Conversation = %s", got.Conversation)
	}
	if got.UpdatedAt.Unix() != updatedAt.Unix() {
		t.Fatalf("UpdatedAt = %v, want %v", got.UpdatedAt, updatedAt)
	}
	if got.CreatedAt.IsZero() {
		t.Fatal("CreatedAt is zero")
	}
}

func TestSessionRepository_RetrieveCodexSessions_NoSessionsDirectory_Normal(t *testing.T) {
	t.Parallel()

	repository := NewRepository()
	sessions, err := repository.RetrieveCodexSessions(t.TempDir())
	if err != nil {
		t.Fatalf("RetrieveCodexSessions() error = %v", err)
	}
	if len(sessions) != 0 {
		t.Fatalf("RetrieveCodexSessions() len = %d, want 0", len(sessions))
	}
}

func TestSessionRepository_RetrieveCodexSessions_UUIDFallbackFromFilename_Normal(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	sessionsDir := filepath.Join(tempDir, "sessions", "2026", "03", "28")
	if err := os.MkdirAll(sessionsDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	filePath := filepath.Join(sessionsDir, "rollout-2026-03-28T01-02-03-aaaaaaaa-bbbb-4ccc-8ddd-eeeeeeeeeeee.jsonl")
	content := strings.Join([]string{
		`{"type":"session_meta","payload":{"timestamp":"2026-03-28T01:02:03Z","cwd":"/workspace/devbox","git":{"branch":"task"}}}`,
		`{"type":"event_msg","payload":{"type":"user_message","message":"hello","kind":"plain"}}`,
	}, "\n")
	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	repository := NewRepository()
	sessions, err := repository.RetrieveCodexSessions(tempDir)
	if err != nil {
		t.Fatalf("RetrieveCodexSessions() error = %v", err)
	}
	if len(sessions) != 1 {
		t.Fatalf("RetrieveCodexSessions() len = %d, want 1", len(sessions))
	}
	if sessions[0].UUID != "aaaaaaaa-bbbb-4ccc-8ddd-eeeeeeeeeeee" {
		t.Fatalf("UUID = %s", sessions[0].UUID)
	}
}

func TestSessionRepository_RetrieveCodexSessions_SessionsRootIsFile(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	sessionsRoot := filepath.Join(tempDir, "sessions")
	if err := os.WriteFile(sessionsRoot, []byte("not directory"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	repository := NewRepository()
	_, err := repository.RetrieveCodexSessions(tempDir)
	if err == nil {
		t.Fatal("RetrieveCodexSessions() error = nil")
	}
}

func TestSessionRepository_RetrieveCodexSessions_InvalidRolloutFile_Normal(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	sessionsDir := filepath.Join(tempDir, "sessions", "2026", "03", "28")
	if err := os.MkdirAll(sessionsDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	invalidPath := filepath.Join(sessionsDir, "rollout-2026-03-28T01-02-03-invalid.jsonl")
	if err := os.WriteFile(invalidPath, []byte("{\"type\":\"event_msg\",\"payload\":{\"type\":\"user_message\",\"message\":\"x\"}}"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	repository := NewRepository()
	sessions, err := repository.RetrieveCodexSessions(tempDir)
	if err != nil {
		t.Fatalf("RetrieveCodexSessions() error = %v", err)
	}
	if len(sessions) != 0 {
		t.Fatalf("RetrieveCodexSessions() len = %d, want 0", len(sessions))
	}
}

func TestExtractUUIDFromFilename_Normal(t *testing.T) {
	t.Parallel()

	createdAt, uuid, ok := parseTimestampUUIDFromFilename("rollout-2026-03-28T01-02-03-aaaaaaaa-bbbb-4ccc-8ddd-eeeeeeeeeeee.jsonl")
	if !ok {
		t.Fatal("parseTimestampUUIDFromFilename() should be ok")
	}
	if uuid != "aaaaaaaa-bbbb-4ccc-8ddd-eeeeeeeeeeee" {
		t.Fatalf("uuid = %s", uuid)
	}
	if createdAt.IsZero() {
		t.Fatal("createdAt is zero")
	}

	_, invalidUUID, ok := parseTimestampUUIDFromFilename("rollout-2026-03-28T01-02-03-invalid.jsonl")
	if ok || invalidUUID != "" {
		t.Fatalf("parseTimestampUUIDFromFilename() should fail")
	}
}

func TestExtractCreatedAtFromFilename_Normal(t *testing.T) {
	t.Parallel()

	createdAt, _, ok := parseTimestampUUIDFromFilename("rollout-2026-03-28T01-02-03-aaaaaaaa-bbbb-4ccc-8ddd-eeeeeeeeeeee.jsonl")
	if !ok || createdAt.IsZero() {
		t.Fatal("parseTimestampUUIDFromFilename() returned invalid createdAt")
	}

	invalidCreatedAt, _, ok := parseTimestampUUIDFromFilename("invalid-file-name.jsonl")
	if ok || !invalidCreatedAt.IsZero() {
		t.Fatalf("parseTimestampUUIDFromFilename() should fail")
	}
}

func TestSessionRepository_ParseEventMessage_Normal(t *testing.T) {
	t.Parallel()

	repository := &SessionRepository{}
	session := &domain.Session{}

	repository.parseEventMessage(json.RawMessage(`{"type":"tool_call","message":"ignore","kind":"plain"}`), session)
	if session.Conversation != "" {
		t.Fatalf("Conversation = %s", session.Conversation)
	}

	repository.parseEventMessage(json.RawMessage(`{"type":"user_message","message":"hello\nworld","kind":"plain"}`), session)
	if session.Conversation != "hello world" {
		t.Fatalf("Conversation = %s", session.Conversation)
	}
}

func TestSessionRepository_ParseSessionMeta_Normal(t *testing.T) {
	t.Parallel()

	repository := &SessionRepository{}
	session := &domain.Session{}

	repository.parseSessionMeta(json.RawMessage(`{"id":"11111111-2222-4333-8444-555555555555","timestamp":"2026-03-28T01:02:03Z","cwd":"/workspace/devbox","git":{"branch":"task"}}`), session)
	if session.UUID != "11111111-2222-4333-8444-555555555555" {
		t.Fatalf("UUID = %s", session.UUID)
	}
	if session.Branch != "task" {
		t.Fatalf("Branch = %s", session.Branch)
	}
	if session.CWD != "/workspace/devbox" {
		t.Fatalf("CWD = %s", session.CWD)
	}
	if session.CreatedAt.IsZero() {
		t.Fatal("CreatedAt is zero")
	}
}
