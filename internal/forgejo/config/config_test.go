package config

import (
	"errors"
	"os"
	"path/filepath"
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
			name:      "project list",
			operation: "project list",
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
			wantMsg:   "forgejo-host が指定されていません",
		},
		{
			name:      "username missing",
			operation: "repo list",
			host:      "https://example.com",
			username:  "",
			token:     "token",
			wantMsg:   "forgejo-username が指定されていません",
		},
		{
			name:      "token missing",
			operation: "repo list",
			host:      "https://example.com",
			username:  "user",
			token:     "",
			wantMsg:   "forgejo-token が指定されていません",
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
	parser := flagParser.NewMockFlagParser()
	parser.SetArgs([]string{"repo", "list"})
	parser.SetBoolFlag("help", false)
	parser.SetBoolFlag("h", false)
	parser.SetBoolFlag("json", false)
	parser.SetStringFlag("forgejo-host", "https://example.com")
	parser.SetStringFlag("forgejo-username", "from-args")
	parser.SetStringFlag("forgejo-token", "token")

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("ParseFlagsWithParser() error = %v", err)
	}
	if cfg.Operation != "repo list" {
		t.Fatalf("Operation = %s, want repo list", cfg.Operation)
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
	tmpDir := t.TempDir()
	dotenvPath := filepath.Join(tmpDir, ".env")
	content := []byte("forgejo-host=https://example.com\nforgejo-username=env-user\nforgejo-token=env-token\n")
	if err := os.WriteFile(dotenvPath, content, 0o600); err != nil {
		t.Fatalf("failed to write .env: %v", err)
	}

	parser := flagParser.NewMockFlagParser()
	parser.SetArgs([]string{"project", "list"})
	parser.SetBoolFlag("json", false)
	cfg, err := ParseFlagsWithParserWithEnvFile(parser, dotenvPath)
	if err != nil {
		t.Fatalf("ParseFlagsWithParserWithEnvFile() error = %v", err)
	}
	if cfg.Operation != "project list" {
		t.Fatalf("Operation = %s, want project list", cfg.Operation)
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
}

func TestParseFlagsWithParser_InterspersedArgs(t *testing.T) {
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
		"-forgejo-username",
		"inter-user",
		"-forgejo-token",
		"inter-token",
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
}

func TestParseFlagsWithParser_ErrorOnParseFailure(t *testing.T) {
	parser := flagParser.NewMockFlagParser()
	parser.SetParseError(errors.New("parse failed"))

	if _, err := ParseFlagsWithParser(parser); err == nil {
		t.Fatalf("expected error")
	}
}
