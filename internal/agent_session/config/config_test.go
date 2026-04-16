package config

import (
	"errors"
	"flag"
	"io"
	"os"
	"strings"
	"testing"
)

type TestConfig struct{}

type MockFlagParser struct {
	stringValues map[string]string
	intValues    map[string]int
	parseError   error
}

func NewMockFlagParser() *MockFlagParser {
	return &MockFlagParser{
		stringValues: make(map[string]string),
		intValues:    make(map[string]int),
	}
}

func (m *MockFlagParser) SetStringFlag(name, value string) {
	m.stringValues[name] = value
}

func (m *MockFlagParser) SetIntFlag(name string, value int) {
	m.intValues[name] = value
}

func (m *MockFlagParser) SetParseError(err error) {
	m.parseError = err
}

func (m *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
	if preset, ok := m.stringValues[name]; ok {
		*p = preset
		return
	}
	*p = value
}

func (m *MockFlagParser) IntVar(p *int, name string, value int, usage string) {
	if preset, ok := m.intValues[name]; ok {
		*p = preset
		return
	}
	*p = value
}

func (m *MockFlagParser) Parse() error {
	return m.parseError
}

type MockOSArgs struct {
	args []string
}

func (m *MockOSArgs) Args() []string {
	return m.args
}

func TestNewConfig_Normal(t *testing.T) {
	t.Parallel()

	cfg, err := NewConfig("retrieve-session", "codex", 10, "20260301", "20260331", "/tmp/codex-home")
	if err != nil {
		t.Fatalf("NewConfig() error = %v", err)
	}
	if cfg.Operation != "retrieve-session" {
		t.Fatalf("Operation = %s", cfg.Operation)
	}
	if cfg.AgentType != "codex" {
		t.Fatalf("AgentType = %s", cfg.AgentType)
	}
	if cfg.Limit != 10 {
		t.Fatalf("Limit = %d", cfg.Limit)
	}
	if cfg.StartDateValue == nil || cfg.EndDateValue == nil {
		t.Fatal("date values should not be nil")
	}
}

func TestNewConfig_Error(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		operation string
		agentType string
		limit     int
		startDate string
		endDate   string
		homeDir   string
		errorPart string
	}{
		{name: "operation未指定", operation: "", agentType: "codex", limit: 10, homeDir: "/tmp", errorPart: "--operation は必須です"},
		{name: "operation不正", operation: "output", agentType: "codex", limit: 10, homeDir: "/tmp", errorPart: "--operation は retrieve-session のみ対応しています"},
		{name: "agentType未指定", operation: "retrieve-session", agentType: "", limit: 10, homeDir: "/tmp", errorPart: "--agent-type は必須です"},
		{name: "agentType不正", operation: "retrieve-session", agentType: "claude", limit: 10, homeDir: "/tmp", errorPart: "--agent-type は codex のみ対応しています"},
		{name: "limit不正", operation: "retrieve-session", agentType: "codex", limit: 0, homeDir: "/tmp", errorPart: "--limit は1以上を指定してください"},
		{name: "startDate不正", operation: "retrieve-session", agentType: "codex", limit: 10, startDate: "2026-01-01", homeDir: "/tmp", errorPart: "--start-date の形式が不正です"},
		{name: "endDate不正", operation: "retrieve-session", agentType: "codex", limit: 10, endDate: "2026-01-31", homeDir: "/tmp", errorPart: "--end-date の形式が不正です"},
		{name: "日付範囲逆転", operation: "retrieve-session", agentType: "codex", limit: 10, startDate: "20260331", endDate: "20260301", homeDir: "/tmp", errorPart: "--start-date は --end-date 以下を指定してください"},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewConfig(tc.operation, tc.agentType, tc.limit, tc.startDate, tc.endDate, tc.homeDir)
			if err == nil {
				t.Fatal("NewConfig() error = nil")
			}
			if err != nil && !strings.Contains(err.Error(), tc.errorPart) {
				t.Fatalf("error = %q, want contains %q", err.Error(), tc.errorPart)
			}
		})
	}
}

func TestConfigParser_ParseFlags_Normal(t *testing.T) {
	t.Parallel()

	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "retrieve-session")
	mockParser.SetStringFlag("agent-type", "codex")
	mockParser.SetIntFlag("limit", 25)
	mockParser.SetStringFlag("start-date", "20260301")
	mockParser.SetStringFlag("end-date", "20260331")
	mockParser.SetStringFlag("agent-home-dir", "/tmp/codex-home")

	parser := NewConfigParser(mockParser, &MockOSArgs{args: []string{"agent-session"}})
	cfg, err := parser.ParseFlags()
	if err != nil {
		t.Fatalf("ParseFlags() error = %v", err)
	}

	if cfg.Operation != "retrieve-session" || cfg.AgentType != "codex" {
		t.Fatalf("cfg = %+v", cfg)
	}
	if cfg.Limit != 25 {
		t.Fatalf("Limit = %d", cfg.Limit)
	}
	if cfg.AgentHomeDir != "/tmp/codex-home" {
		t.Fatalf("AgentHomeDir = %s", cfg.AgentHomeDir)
	}
}

func TestConfigParser_ParseFlags_ParseError(t *testing.T) {
	t.Parallel()

	mockParser := NewMockFlagParser()
	mockParser.SetParseError(errors.New("parse failed"))

	parser := NewConfigParser(mockParser, &MockOSArgs{args: []string{"agent-session"}})
	_, err := parser.ParseFlags()
	if err == nil {
		t.Fatal("ParseFlags() error = nil")
	}
	if !strings.Contains(err.Error(), "フラグの解析に失敗しました") {
		t.Fatalf("error = %q", err.Error())
	}
}

func TestConfig_ParseFlags_GlobalParser_Normal(t *testing.T) {
	oldArgs := os.Args
	oldCommandLine := flag.CommandLine
	defer func() {
		os.Args = oldArgs
		flag.CommandLine = oldCommandLine
	}()

	os.Args = []string{
		"agent-session",
		"-operation=retrieve-session",
		"-agent-type=codex",
		"-limit=3",
		"-agent-home-dir=/tmp/.codex",
	}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flag.CommandLine.SetOutput(io.Discard)

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("ParseFlags() error = %v", err)
	}
	if cfg.Operation != "retrieve-session" {
		t.Fatalf("Operation = %s", cfg.Operation)
	}
	if cfg.AgentType != "codex" {
		t.Fatalf("AgentType = %s", cfg.AgentType)
	}
	if cfg.Limit != 3 {
		t.Fatalf("Limit = %d", cfg.Limit)
	}
	if cfg.AgentHomeDir != "/tmp/.codex" {
		t.Fatalf("AgentHomeDir = %s", cfg.AgentHomeDir)
	}
}

func TestConfig_PrintUsage_Normal(t *testing.T) {
	output := captureStdout(func() {
		PrintUsage()
	})

	if !strings.Contains(output, "使用方法:") {
		t.Fatalf("usage output is unexpected: %s", output)
	}
	if !strings.Contains(output, "-operation string") {
		t.Fatalf("usage output is unexpected: %s", output)
	}
}

func captureStdout(run func()) string {
	original := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return ""
	}
	defer r.Close()

	os.Stdout = w
	run()
	_ = w.Close()
	os.Stdout = original

	buf, _ := io.ReadAll(r)
	return string(buf)
}
