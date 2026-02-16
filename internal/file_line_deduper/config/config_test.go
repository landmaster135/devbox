package config

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

type mockFlagParser struct {
	stringVars map[string]*string
	intVars    map[string]*int
	boolVars   map[string]*bool
	args       []string
	parseError error
	parseHook  func()
}

func newMockFlagParser() *mockFlagParser {
	return &mockFlagParser{
		stringVars: make(map[string]*string),
		intVars:    make(map[string]*int),
		boolVars:   make(map[string]*bool),
		args:       []string{},
	}
}

func (m *mockFlagParser) StringVar(p *string, name string, value string, usage string) {
	*p = value
	m.stringVars[name] = p
}

func (m *mockFlagParser) IntVar(p *int, name string, value int, usage string) {
	*p = value
	m.intVars[name] = p
}

func (m *mockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	*p = value
	m.boolVars[name] = p
}

func (m *mockFlagParser) Parse() error {
	if m.parseHook != nil {
		m.parseHook()
	}
	return m.parseError
}

func (m *mockFlagParser) Args() []string {
	return m.args
}

func (m *mockFlagParser) setStringValue(name, value string) {
	if ptr, ok := m.stringVars[name]; ok {
		*ptr = value
	}
}

func (m *mockFlagParser) setIntValue(name string, value int) {
	if ptr, ok := m.intVars[name]; ok {
		*ptr = value
	}
}

func (m *mockFlagParser) setBoolValue(name string, value bool) {
	if ptr, ok := m.boolVars[name]; ok {
		*ptr = value
	}
}

func (m *mockFlagParser) setParseError(err error) {
	m.parseError = err
}

func (m *mockFlagParser) setParseHook(parseHook func()) {
	m.parseHook = parseHook
}

func TestNewConfig_Normal(t *testing.T) {
	cfg, err := NewConfig("test.txt", 1, 5)
	if err != nil {
		t.Fatalf("NewConfig() error = %v, want nil", err)
	}

	if cfg.FilePath != "test.txt" {
		t.Errorf("FilePath = %q, want %q", cfg.FilePath, "test.txt")
	}

	if cfg.StartPos != 1 {
		t.Errorf("StartPos = %d, want %d", cfg.StartPos, 1)
	}

	if cfg.EndPos != 5 {
		t.Errorf("EndPos = %d, want %d", cfg.EndPos, 5)
	}
}

func TestNewConfig_Error(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		start   int
		end     int
		wantErr string
	}{
		{
			name:    "ファイル未指定",
			file:    "",
			start:   1,
			end:     5,
			wantErr: "ファイルパスを指定してください（-file オプション）",
		},
		{
			name:    "start_end未指定",
			file:    "test.txt",
			start:   0,
			end:     0,
			wantErr: "開始位置と終了位置を指定してください（-start と -end オプション）",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewConfig(tt.file, tt.start, tt.end)
			if err == nil {
				t.Fatal("NewConfig() error = nil, want error")
			}

			if err.Error() != tt.wantErr {
				t.Errorf("NewConfig() error = %q, want %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestParseFlagsWithParser_Normal(t *testing.T) {
	parser := newMockFlagParser()
	parser.setParseHook(func() {
		parser.setStringValue("file", "data.txt")
		parser.setIntValue("start", 3)
		parser.setIntValue("end", 6)
	})

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("ParseFlagsWithParser() error = %v, want nil", err)
	}

	if cfg.FilePath != "data.txt" || cfg.StartPos != 3 || cfg.EndPos != 6 {
		t.Errorf("ParseFlagsWithParser() returned unexpected config: %+v", cfg)
	}
}

func TestParseFlagsWithParser_Help(t *testing.T) {
	parser := newMockFlagParser()
	parser.setParseHook(func() {
		parser.setBoolValue("help", true)
	})

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("ParseFlagsWithParser() error = %v, want nil", err)
	}

	if !cfg.Help {
		t.Error("Config.Help = false, want true")
	}
}

func TestParseFlagsWithParser_ParseError(t *testing.T) {
	parser := newMockFlagParser()
	parser.setParseError(errors.New("parse failed"))

	_, err := ParseFlagsWithParser(parser)
	if err == nil {
		t.Fatal("ParseFlagsWithParser() error = nil, want error")
	}

	if !strings.Contains(err.Error(), "フラグの解析に失敗しました") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestParseFlagsWithParser_ConfigError(t *testing.T) {
	parser := newMockFlagParser()
	parser.setParseHook(func() {
		parser.setStringValue("file", "")
		parser.setIntValue("start", 1)
		parser.setIntValue("end", 2)
	})

	_, err := ParseFlagsWithParser(parser)
	if err == nil {
		t.Fatal("ParseFlagsWithParser() error = nil, want error")
	}

	if !strings.Contains(err.Error(), "ファイルパスを指定してください") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestParseFlagsWithArgs_Normal(t *testing.T) {
	cfg, err := ParseFlagsWithArgs([]string{"-file", "data.txt", "-start", "3", "-end", "6"})
	if err != nil {
		t.Fatalf("ParseFlagsWithArgs() error = %v, want nil", err)
	}

	if cfg.FilePath != "data.txt" || cfg.StartPos != 3 || cfg.EndPos != 6 {
		t.Errorf("ParseFlagsWithArgs() returned unexpected config: %+v", cfg)
	}
}

func TestParseFlags_Normal(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"cmd", "-file", "sample.txt", "-start", "2", "-end", "4"}

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("ParseFlags() error = %v, want nil", err)
	}

	if cfg.FilePath != "sample.txt" || cfg.StartPos != 2 || cfg.EndPos != 4 {
		t.Errorf("ParseFlags() returned unexpected config: %+v", cfg)
	}
}

func TestParseFlagsWithArgs_InvalidFlag(t *testing.T) {
	_, err := ParseFlagsWithArgs([]string{"-invalid"})
	if err == nil {
		t.Fatal("ParseFlagsWithArgs() error = nil, want error")
	}

	if !strings.Contains(err.Error(), "フラグの解析に失敗しました") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestPrintUsage(t *testing.T) {
	origStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	defer r.Close()

	os.Stderr = w
	PrintUsage()
	w.Close()
	os.Stderr = origStderr

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("io.Copy() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "file-line-deduper") {
		t.Errorf("PrintUsage() output does not contain tool name: %s", output)
	}
	if !strings.Contains(output, "-file") {
		t.Errorf("PrintUsage() output does not contain -file option: %s", output)
	}
}
