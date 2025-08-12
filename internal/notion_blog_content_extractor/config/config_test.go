package config

import (
	"fmt"
	"os"
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

func TestParseFlagsWithParser_Normal(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*MockFlagParser)
		expectedCfg *Config
		expectError bool
	}{
		{
			name: "ValidFlags_Normal",
			setupMock: func(parser *MockFlagParser) {
				parser.SetStringFlag("src-dir", "/test/src")
				parser.SetStringFlag("dest-dir", "/test/dest")
				parser.SetBoolFlag("help", false)
			},
			expectedCfg: &Config{
				SrcDir:  "/test/src",
				DestDir: "/test/dest",
				Help:    false,
			},
			expectError: false,
		},
		{
			name: "HelpFlag_Normal",
			setupMock: func(parser *MockFlagParser) {
				parser.SetBoolFlag("help", true)
			},
			expectedCfg: &Config{
				SrcDir:  "",
				DestDir: "",
				Help:    true,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewMockFlagParser()
			tt.setupMock(parser)

			cfg, err := ParseFlagsWithParser(parser)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if cfg == nil {
					t.Fatal("Expected config to be non-nil")
				}
				if cfg.SrcDir != tt.expectedCfg.SrcDir {
					t.Errorf("Expected SrcDir '%s', got '%s'", tt.expectedCfg.SrcDir, cfg.SrcDir)
				}
				if cfg.DestDir != tt.expectedCfg.DestDir {
					t.Errorf("Expected DestDir '%s', got '%s'", tt.expectedCfg.DestDir, cfg.DestDir)
				}
				if cfg.Help != tt.expectedCfg.Help {
					t.Errorf("Expected Help %v, got %v", tt.expectedCfg.Help, cfg.Help)
				}
			}
		})
	}
}

func TestParseFlagsWithParser_MissingRequiredParams(t *testing.T) {
	tests := []struct {
		name        string
		srcDir      string
		destDir     string
		help        bool
		expectError bool
		errorMsg    string
	}{
		{
			name:        "MissingSrcDir_Error",
			srcDir:      "",
			destDir:     "/test/dest",
			help:        false,
			expectError: true,
			errorMsg:    "src-dir パラメータは必須です",
		},
		{
			name:        "MissingDestDir_Error",
			srcDir:      "/test/src",
			destDir:     "",
			help:        false,
			expectError: true,
			errorMsg:    "dest-dir パラメータは必須です",
		},
		{
			name:        "BothMissing_Error",
			srcDir:      "",
			destDir:     "",
			help:        false,
			expectError: true,
			errorMsg:    "src-dir パラメータは必須です",
		},
		{
			name:        "HelpFlagSet_NoError",
			srcDir:      "",
			destDir:     "",
			help:        true,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewMockFlagParser()
			parser.SetStringFlag("src-dir", tt.srcDir)
			parser.SetStringFlag("dest-dir", tt.destDir)
			parser.SetBoolFlag("help", tt.help)

			_, err := ParseFlagsWithParser(parser)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestParseFlagsWithParser_ParseError(t *testing.T) {
	parser := NewMockFlagParser()
	parser.SetParseError(fmt.Errorf("parse error"))

	_, err := ParseFlagsWithParser(parser)

	if err == nil {
		t.Error("Expected error but got none")
	}

	expectedMsg := "フラグの解析に失敗しました: parse error"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestParseFlags_Normal(t *testing.T) {
	// 元のos.Argsを保存
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// テスト用の引数を設定
	os.Args = []string{
		"test-program",
		"-src-dir=/test/src",
		"-dest-dir=/test/dest",
	}

	cfg, err := ParseFlags()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if cfg == nil {
		t.Fatal("Expected config to be non-nil")
		return
	}

	if cfg.SrcDir != "/test/src" {
		t.Errorf("Expected SrcDir '/test/src', got '%s'", cfg.SrcDir)
	}

	if cfg.DestDir != "/test/dest" {
		t.Errorf("Expected DestDir '/test/dest', got '%s'", cfg.DestDir)
	}

	if cfg.Help {
		t.Error("Expected Help to be false")
	}
}

func TestParseFlags_HelpFlag(t *testing.T) {
	// 元のos.Argsを保存
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// テスト用の引数を設定
	os.Args = []string{
		"test-program",
		"-help",
	}

	cfg, err := ParseFlags()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if cfg == nil {
		t.Fatal("Expected config to be non-nil")
		return
	}

	if !cfg.Help {
		t.Error("Expected Help to be true")
	}
}

func TestParseFlags_MissingRequiredParams(t *testing.T) {
	// 元のos.Argsを保存
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// テスト用の引数を設定（必須パラメータなし）
	os.Args = []string{
		"test-program",
	}

	_, err := ParseFlags()

	if err == nil {
		t.Error("Expected error but got none")
	}

	expectedMsg := "src-dir パラメータは必須です"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestPrintUsage_Normal(t *testing.T) {
	// PrintUsage関数が正常に実行されることを確認
	// パニックが発生しないことをテスト
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PrintUsage panicked: %v", r)
		}
	}()

	PrintUsage()
	// 関数が正常に完了すればテスト成功
}

func TestConfig_Validation(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		valid  bool
	}{
		{
			name: "ValidConfig_Normal",
			config: &Config{
				SrcDir:  "/test/src",
				DestDir: "/test/dest",
				Help:    false,
			},
			valid: true,
		},
		{
			name: "HelpConfig_Valid",
			config: &Config{
				SrcDir:  "",
				DestDir: "",
				Help:    true,
			},
			valid: true,
		},
		{
			name: "EmptyConfig_Invalid",
			config: &Config{
				SrcDir:  "",
				DestDir: "",
				Help:    false,
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Configの妥当性を検証するロジックをテスト
			// 実際の実装では、ParseFlagsWithParserで検証が行われる
			isValid := tt.config.Help || (tt.config.SrcDir != "" && tt.config.DestDir != "")

			if isValid != tt.valid {
				t.Errorf("Expected config validity %v, got %v", tt.valid, isValid)
			}
		})
	}
}

func TestMockFlagParser_StringVar(t *testing.T) {
	parser := NewMockFlagParser()

	var testValue string
	parser.StringVar(&testValue, "test-flag", "default", "test usage")

	// デフォルト値が設定されることを確認
	if testValue != "default" {
		t.Errorf("Expected default value 'default', got '%s'", testValue)
	}

	// 事前設定値が優先されることを確認
	parser.SetStringFlag("test-flag", "preset")
	var testValue2 string
	parser.StringVar(&testValue2, "test-flag", "default", "test usage")

	if testValue2 != "preset" {
		t.Errorf("Expected preset value 'preset', got '%s'", testValue2)
	}
}

func TestMockFlagParser_BoolVar(t *testing.T) {
	parser := NewMockFlagParser()

	var testValue bool
	parser.BoolVar(&testValue, "test-flag", true, "test usage")

	// デフォルト値が設定されることを確認
	if !testValue {
		t.Error("Expected default value true, got false")
	}

	// 事前設定値が優先されることを確認
	parser.SetBoolFlag("test-flag", false)
	var testValue2 bool
	parser.BoolVar(&testValue2, "test-flag", true, "test usage")

	if testValue2 {
		t.Error("Expected preset value false, got true")
	}
}

func TestMockFlagParser_Parse(t *testing.T) {
	parser := NewMockFlagParser()

	// 正常なケース
	err := parser.Parse()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// エラーケース
	parser.SetParseError(fmt.Errorf("test error"))
	err = parser.Parse()
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "test error" {
		t.Errorf("Expected error 'test error', got '%s'", err.Error())
	}
}

func TestMockFlagParser_Args(t *testing.T) {
	parser := NewMockFlagParser()

	// デフォルトでは空のスライス
	args := parser.Args()
	if len(args) != 0 {
		t.Errorf("Expected empty args, got %v", args)
	}

	// 引数を設定
	testArgs := []string{"arg1", "arg2", "arg3"}
	parser.SetArgs(testArgs)
	args = parser.Args()

	if len(args) != len(testArgs) {
		t.Errorf("Expected %d args, got %d", len(testArgs), len(args))
	}

	for i, arg := range testArgs {
		if args[i] != arg {
			t.Errorf("Expected arg[%d] '%s', got '%s'", i, arg, args[i])
		}
	}
}
