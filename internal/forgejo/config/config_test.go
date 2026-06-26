package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	flagParser "github.com/landmaster135/devbox/internal/forgejo/infrastructures/flag_parser"
)

func TestNewConfig_Normal(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		host      string
		username  string
		token     string
		wantErr   bool
	}{
		{
			name:      "repo list",
			operation: "repo list",
			host:      "https://example.com",
			username:  "user",
			token:     "token",
			wantErr:   false,
		},
		{
			name:      "issue list",
			operation: "issue list",
			host:      "https://example.com",
			username:  "user",
			token:     "token",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := NewConfig(tt.operation, tt.host, tt.username, tt.token, false)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if cfg.Operation != tt.operation {
				t.Fatalf("Operation = %s, want %s", cfg.Operation, tt.operation)
			}
			if cfg.Host != tt.host {
				t.Fatalf("Host = %s, want %s", cfg.Host, tt.host)
			}
			if cfg.Username != tt.username {
				t.Fatalf("Username = %s, want %s", cfg.Username, tt.username)
			}
			if cfg.Token != tt.token {
				t.Fatalf("Token = %s, want %s", cfg.Token, tt.token)
			}
			if cfg.ReposWorkers != defaultWorkers {
				t.Fatalf("ReposWorkers = %d, want %d", cfg.ReposWorkers, defaultWorkers)
			}
		})
	}
}

func TestNewConfig_Error(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		host      string
		username  string
		token     string
		wantMsg   string
	}{
		{
			name:      "operation empty",
			operation: "",
			host:      "https://example.com",
			username:  "user",
			token:     "token",
			wantMsg:   "operationが指定されていません",
		},
		{
			name:      "invalid operation",
			operation: "unknown",
			host:      "https://example.com",
			username:  "user",
			token:     "token",
			wantMsg:   "未対応のoperationです: unknown",
		},
		{
			name:      "host missing",
			operation: "repo list",
			host:      "",
			username:  "user",
			token:     "token",
			wantMsg:   "FORGEJO_HOST が指定されていません",
		},
		{
			name:      "username missing",
			operation: "repo list",
			host:      "https://example.com",
			username:  "",
			token:     "token",
			wantMsg:   "FORGEJO_USERNAME が指定されていません",
		},
		{
			name:      "token missing",
			operation: "repo list",
			host:      "https://example.com",
			username:  "user",
			token:     "",
			wantMsg:   "FORGEJO_TOKEN が指定されていません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewConfig(tt.operation, tt.host, tt.username, tt.token, false)
			if err == nil {
				t.Fatalf("NewConfig() expected error")
			}
			if err.Error() != tt.wantMsg {
				t.Fatalf("unexpected error: %v", err.Error())
			}
		})
	}
}

func TestParseFlagsWithParser_ResolveOperationFromArgs(t *testing.T) {
	setForgejoEnv(t, "https://example.com", "from-env", "token")

	parser := flagParser.NewMockFlagParser()
	parser.SetArgs([]string{"repo", "list"})
	parser.SetBoolFlag("help", false)
	parser.SetBoolFlag("h", false)
	parser.SetBoolFlag("json", false)

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("ParseFlagsWithParser() error = %v", err)
	}
	if cfg.Operation != "repo list" {
		t.Fatalf("Operation = %s, want repo list", cfg.Operation)
	}
}

func TestParseFlagsWithParser_ResolveOperationFromFlag(t *testing.T) {
	setForgejoEnv(t, "https://example.com", "from-env", "token")

	parser := flagParser.NewMockFlagParser()
	parser.SetStringFlag("operation", "issue list")

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("ParseFlagsWithParser() error = %v", err)
	}
	if cfg.Operation != "issue list" {
		t.Fatalf("Operation = %s, want issue list", cfg.Operation)
	}
}

func TestParseFlagsWithParser_Help(t *testing.T) {
	parser := flagParser.NewMockFlagParser()
	parser.SetBoolFlag("help", true)

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("ParseFlagsWithParser() error = %v", err)
	}
	if !cfg.Help {
		t.Fatalf("Help should be true")
	}
}

func TestParseFlagsWithParser_UsesDotEnv(t *testing.T) {
	unsetForgejoEnv(t)
	tmpDir := t.TempDir()
	dotenvPath := filepath.Join(tmpDir, ".env")
	content := []byte("FORGEJO_HOST=https://example.com\nFORGEJO_USERNAME=env-user\nFORGEJO_TOKEN=env-token\n")
	if err := os.WriteFile(dotenvPath, content, 0o600); err != nil {
		t.Fatalf("failed to write .env: %v", err)
	}
	chdirForTest(t, tmpDir)

	parser := flagParser.NewMockFlagParser()
	parser.SetArgs([]string{"issue", "list"})
	parser.SetBoolFlag("json", false)
	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("ParseFlagsWithParser() error = %v", err)
	}
	if cfg.Operation != "issue list" {
		t.Fatalf("Operation = %s, want issue list", cfg.Operation)
	}
	if cfg.Host != "https://example.com" {
		t.Fatalf("Host = %s, want https://example.com", cfg.Host)
	}
	if cfg.Username != "env-user" {
		t.Fatalf("Username = %s, want env-user", cfg.Username)
	}
	if cfg.Token != "env-token" {
		t.Fatalf("Token = %s, want env-token", cfg.Token)
	}
	if cfg.ReposWorkers != defaultWorkers {
		t.Fatalf("ReposWorkers = %d, want %d", cfg.ReposWorkers, defaultWorkers)
	}
}

func TestParseFlagsWithParser_UsesEnvironmentVariables(t *testing.T) {
	setForgejoEnv(t, "https://env.example.com", "env-user", "env-token")

	parser := flagParser.NewMockFlagParser()
	parser.SetArgs([]string{"repo", "list"})
	parser.SetBoolFlag("json", true)

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("ParseFlagsWithParser() error = %v", err)
	}
	if cfg.Host != "https://env.example.com" {
		t.Fatalf("Host = %s, want https://env.example.com", cfg.Host)
	}
	if cfg.Username != "env-user" {
		t.Fatalf("Username = %s, want env-user", cfg.Username)
	}
	if cfg.Token != "env-token" {
		t.Fatalf("Token = %s, want env-token", cfg.Token)
	}
	if !cfg.JSON {
		t.Fatalf("JSON should be true")
	}
}

func TestParseFlagsWithParser_EnvironmentOverridesDotEnv(t *testing.T) {
	setForgejoEnv(t, "https://env.example.com", "env-user", "env-token")
	tmpDir := t.TempDir()
	dotenvPath := filepath.Join(tmpDir, ".env")
	content := []byte("FORGEJO_HOST=https://dotenv.example.com\nFORGEJO_USERNAME=dotenv-user\nFORGEJO_TOKEN=dotenv-token\n")
	if err := os.WriteFile(dotenvPath, content, 0o600); err != nil {
		t.Fatalf("failed to write .env: %v", err)
	}
	chdirForTest(t, tmpDir)

	parser := flagParser.NewMockFlagParser()
	parser.SetArgs([]string{"repo", "list"})

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("ParseFlagsWithParser() error = %v", err)
	}
	if cfg.Host != "https://env.example.com" {
		t.Fatalf("Host = %s, want https://env.example.com", cfg.Host)
	}
	if cfg.Username != "env-user" {
		t.Fatalf("Username = %s, want env-user", cfg.Username)
	}
	if cfg.Token != "env-token" {
		t.Fatalf("Token = %s, want env-token", cfg.Token)
	}
}

func TestParseFlagsWithParser_CustomReposWorkers(t *testing.T) {
	setForgejoEnv(t, "https://example.com", "user", "token")

	parser := flagParser.NewMockFlagParser()
	parser.SetArgs([]string{"repo", "list"})
	parser.SetBoolFlag("help", false)
	parser.SetBoolFlag("h", false)
	parser.SetBoolFlag("json", false)
	parser.SetStringFlag("repos-workers", "8")

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("ParseFlagsWithParser() error = %v", err)
	}
	if cfg.ReposWorkers != 8 {
		t.Fatalf("ReposWorkers = %d, want 8", cfg.ReposWorkers)
	}
}

func TestParseFlagsWithParser_InvalidReposWorkers(t *testing.T) {
	setForgejoEnv(t, "https://example.com", "user", "token")

	parser := flagParser.NewMockFlagParser()
	parser.SetArgs([]string{"repo", "list"})
	parser.SetBoolFlag("help", false)
	parser.SetBoolFlag("h", false)
	parser.SetBoolFlag("json", false)
	parser.SetStringFlag("repos-workers", "abc")

	if _, err := ParseFlagsWithParser(parser); err == nil {
		t.Fatal("expected error")
	}
}

func TestParseFlagsWithParser_NonPositiveReposWorkers_Error(t *testing.T) {
	setForgejoEnv(t, "https://example.com", "user", "token")

	parser := flagParser.NewMockFlagParser()
	parser.SetArgs([]string{"repo", "list"})
	parser.SetStringFlag("repos-workers", "0")

	_, err := ParseFlagsWithParser(parser)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "repos-workers は1以上を指定してください") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseFlagsWithParser_InterspersedArgs(t *testing.T) {
	setForgejoEnv(t, "https://example.com", "inter-user", "inter-token")

	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	os.Args = []string{
		"forgejo.exe",
		"repo",
		"list",
		"-json",
	}

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("ParseFlags() error = %v", err)
	}
	if cfg.Host != "https://example.com" {
		t.Fatalf("Host = %s, want https://example.com", cfg.Host)
	}
	if cfg.Username != "inter-user" {
		t.Fatalf("Username = %s, want inter-user", cfg.Username)
	}
	if cfg.Token != "inter-token" {
		t.Fatalf("Token = %s, want inter-token", cfg.Token)
	}
	if !cfg.JSON {
		t.Fatalf("JSON should be true")
	}
}

func TestParseFlagsWithParser_RemovedConnectionFlag_Error(t *testing.T) {
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	os.Args = []string{
		"forgejo.exe",
		"repo",
		"list",
		"-forgejo-host",
		"https://example.com",
	}

	_, err := ParseFlags()
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "フラグ解析に失敗しました") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseFlagsWithParser_MissingEnvironmentVariable_Error(t *testing.T) {
	unsetForgejoEnv(t)
	chdirForTest(t, t.TempDir())

	parser := flagParser.NewMockFlagParser()
	parser.SetArgs([]string{"repo", "list"})

	_, err := ParseFlagsWithParser(parser)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "環境変数 FORGEJO_HOST が設定されていません") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseFlagsWithParser_ErrorOnParseFailure(t *testing.T) {
	parser := flagParser.NewMockFlagParser()
	parser.SetParseError(errors.New("parse failed"))

	if _, err := ParseFlagsWithParser(parser); err == nil {
		t.Fatalf("expected error")
	}
}

func setForgejoEnv(t *testing.T, host, username, token string) {
	t.Helper()
	t.Setenv(envKeyHost, host)
	t.Setenv(envKeyUsername, username)
	t.Setenv(envKeyToken, token)
}

func unsetForgejoEnv(t *testing.T) {
	t.Helper()
	keys := []string{envKeyHost, envKeyUsername, envKeyToken}
	originals := make(map[string]string, len(keys))
	present := make(map[string]bool, len(keys))
	for _, key := range keys {
		value, ok := os.LookupEnv(key)
		originals[key] = value
		present[key] = ok
		if err := os.Unsetenv(key); err != nil {
			t.Fatalf("failed to unset %s: %v", key, err)
		}
	}
	t.Cleanup(func() {
		for _, key := range keys {
			if present[key] {
				_ = os.Setenv(key, originals[key])
				continue
			}
			_ = os.Unsetenv(key)
		}
	})
}

func chdirForTest(t *testing.T, dir string) {
	t.Helper()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalDir)
	})
}
