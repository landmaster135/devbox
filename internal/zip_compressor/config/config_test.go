package config

import (
	"testing"
)

// MockFlagParser はテスト用のFlagParserモック
type MockFlagParser struct {
	stringVars   map[string]*string
	boolVars     map[string]*bool
	stringValues map[string]string // 事前設定された文字列値
	boolValues   map[string]bool   // 事前設定されたブール値
	args         []string
	parseError   error
}

// NewMockFlagParser は新しいMockFlagParserを作成する
func NewMockFlagParser() *MockFlagParser {
	return &MockFlagParser{
		stringVars:   make(map[string]*string),
		boolVars:     make(map[string]*bool),
		stringValues: make(map[string]string),
		boolValues:   make(map[string]bool),
		args:         []string{},
	}
}

// StringVar は文字列フラグを定義する（モック）
func (m *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
	// 事前設定された値があるかチェック
	if presetValue, exists := m.stringValues[name]; exists {
		*p = presetValue
	} else {
		*p = value // デフォルト値を設定
	}
	m.stringVars[name] = p
}

// BoolVar はブールフラグを定義する（モック）
func (m *MockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	// 事前設定された値があるかチェック
	if presetValue, exists := m.boolValues[name]; exists {
		*p = presetValue
	} else {
		*p = value // デフォルト値を設定
	}
	m.boolVars[name] = p
}

func (m *MockFlagParser) Parse() error {
	return m.parseError
}

func (m *MockFlagParser) Args() []string {
	return m.args
}

// SetStringFlag はテスト用に文字列フラグの値を設定する
func (m *MockFlagParser) SetStringFlag(name, value string) {
	m.stringValues[name] = value
	if p, exists := m.stringVars[name]; exists {
		*p = value
	}
}

// SetBoolFlag はテスト用にブールフラグの値を設定する
func (m *MockFlagParser) SetBoolFlag(name string, value bool) {
	m.boolValues[name] = value
	if p, exists := m.boolVars[name]; exists {
		*p = value
	}
}

// SetArgs はテスト用に残りの引数を設定する
func (m *MockFlagParser) SetArgs(args []string) {
	m.args = args
}

// SetParseError はテスト用に解析エラーを設定する
func (m *MockFlagParser) SetParseError(err error) {
	m.parseError = err
}

func TestNewConfig_Normal(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		path      string
		wantErr   bool
	}{
		{
			name:      "正常ケース_compress",
			operation: "compress",
			path:      "/path/to/file",
			wantErr:   false,
		},
		{
			name:      "正常ケース_decompress",
			operation: "decompress",
			path:      "/path/to/archive.zip",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := NewConfig(tt.operation, tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if cfg.Operation != tt.operation {
					t.Errorf("NewConfig() Operation = %v, want %v", cfg.Operation, tt.operation)
				}
				if cfg.Path != tt.path {
					t.Errorf("NewConfig() Path = %v, want %v", cfg.Path, tt.path)
				}
			}
		})
	}
}

func TestNewConfig_Error(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		path      string
		wantErr   string
	}{
		{
			name:      "操作タイプが空",
			operation: "",
			path:      "/path/to/file",
			wantErr:   "操作タイプが指定されていません",
		},
		{
			name:      "無効な操作タイプ",
			operation: "invalid",
			path:      "/path/to/file",
			wantErr:   "無効な操作タイプです: invalid",
		},
		{
			name:      "パスが空",
			operation: "compress",
			path:      "",
			wantErr:   "パスが指定されていません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewConfig(tt.operation, tt.path)
			if err == nil {
				t.Errorf("NewConfig() expected error but got nil")
				return
			}
			if err.Error() != tt.wantErr {
				t.Errorf("NewConfig() error = %v, want %v", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestParseFlagsWithParser_Normal(t *testing.T) {
	tests := []struct {
		name         string
		operation    string
		path         string
		help         bool
		args         []string
		expectedOp   string
		expectedPath string
		expectedHelp bool
	}{
		{
			name:         "フラグ指定_compress",
			operation:    "compress",
			path:         "/path/to/file",
			help:         false,
			args:         []string{},
			expectedOp:   "compress",
			expectedPath: "/path/to/file",
			expectedHelp: false,
		},
		{
			name:         "フラグ指定_decompress",
			operation:    "decompress",
			path:         "/path/to/archive.zip",
			help:         false,
			args:         []string{},
			expectedOp:   "decompress",
			expectedPath: "/path/to/archive.zip",
			expectedHelp: false,
		},
		{
			name:         "位置引数指定",
			operation:    "",
			path:         "",
			help:         false,
			args:         []string{"compress", "/path/to/file"},
			expectedOp:   "compress",
			expectedPath: "/path/to/file",
			expectedHelp: false,
		},
		{
			name:         "ヘルプ指定",
			operation:    "",
			path:         "",
			help:         true,
			args:         []string{},
			expectedOp:   "",
			expectedPath: "",
			expectedHelp: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewMockFlagParser()
			parser.SetArgs(tt.args)

			// 事前に値を設定
			parser.SetStringFlag("operation", tt.operation)
			parser.SetStringFlag("o", tt.operation)
			parser.SetStringFlag("path", tt.path)
			parser.SetStringFlag("p", tt.path)
			parser.SetBoolFlag("help", tt.help)
			parser.SetBoolFlag("h", tt.help)

			cfg, err := ParseFlagsWithParser(parser)

			if err != nil {
				t.Errorf("ParseFlagsWithParser() error = %v", err)
				return
			}

			if tt.expectedHelp {
				if !cfg.Help {
					t.Errorf("ParseFlagsWithParser() Help = %v, want %v", cfg.Help, tt.expectedHelp)
				}
				return
			}

			if cfg.Operation != tt.expectedOp {
				t.Errorf("ParseFlagsWithParser() Operation = %v, want %v", cfg.Operation, tt.expectedOp)
			}
			if cfg.Path != tt.expectedPath {
				t.Errorf("ParseFlagsWithParser() Path = %v, want %v", cfg.Path, tt.expectedPath)
			}
		})
	}
}

func TestStartsWith_Normal(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		prefix string
		want   bool
	}{
		{
			name:   "前方一致_true",
			s:      "hello world",
			prefix: "hello",
			want:   true,
		},
		{
			name:   "前方一致_false",
			s:      "hello world",
			prefix: "world",
			want:   false,
		},
		{
			name:   "完全一致",
			s:      "test",
			prefix: "test",
			want:   true,
		},
		{
			name:   "空文字列",
			s:      "test",
			prefix: "",
			want:   true,
		},
		{
			name:   "プレフィックスが長い",
			s:      "hi",
			prefix: "hello",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := startsWith(tt.s, tt.prefix)
			if got != tt.want {
				t.Errorf("startsWith() = %v, want %v", got, tt.want)
			}
		})
	}
}
