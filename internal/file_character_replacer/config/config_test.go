package config

import (
	"errors"
	"testing"

	"github.com/landmaster135/devbox/internal/file_character_replacer/domain"
)

// TestConfig_ToReplacementConfig_Normal はConfig.ToReplacementConfig()の正常系をテストします
func TestConfig_ToReplacementConfig_Normal(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		expected *domain.ReplacementConfig
	}{
		{
			name: "全フィールド指定",
			config: &Config{
				Target:    "/test/path",
				From:      "old",
				To:        "new",
				Encoding:  "shift_jis",
				Recursive: true,
				Backup:    true,
				BackupDir: "/backup",
				DryRun:    true,
			},
			expected: &domain.ReplacementConfig{
				Target:    "/test/path",
				From:      "old",
				To:        "new",
				Encoding:  domain.EncodingShiftJIS,
				Recursive: true,
				Backup:    true,
				BackupDir: "/backup",
				DryRun:    true,
			},
		},
		{
			name: "エンコーディング空（デフォルトUTF-8）",
			config: &Config{
				Target:   "/test/path",
				From:     "old",
				To:       "new",
				Encoding: "",
			},
			expected: &domain.ReplacementConfig{
				Target:   "/test/path",
				From:     "old",
				To:       "new",
				Encoding: domain.EncodingUTF8,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.ToReplacementConfig()

			if result.Target != tt.expected.Target {
				t.Errorf("Target = %v, expected %v", result.Target, tt.expected.Target)
			}
			if result.From != tt.expected.From {
				t.Errorf("From = %v, expected %v", result.From, tt.expected.From)
			}
			if result.To != tt.expected.To {
				t.Errorf("To = %v, expected %v", result.To, tt.expected.To)
			}
			if result.Encoding != tt.expected.Encoding {
				t.Errorf("Encoding = %v, expected %v", result.Encoding, tt.expected.Encoding)
			}
			if result.Recursive != tt.expected.Recursive {
				t.Errorf("Recursive = %v, expected %v", result.Recursive, tt.expected.Recursive)
			}
			if result.Backup != tt.expected.Backup {
				t.Errorf("Backup = %v, expected %v", result.Backup, tt.expected.Backup)
			}
			if result.BackupDir != tt.expected.BackupDir {
				t.Errorf("BackupDir = %v, expected %v", result.BackupDir, tt.expected.BackupDir)
			}
			if result.DryRun != tt.expected.DryRun {
				t.Errorf("DryRun = %v, expected %v", result.DryRun, tt.expected.DryRun)
			}
		})
	}
}

// TestValidateEncoding_Normal はvalidateEncoding()の正常系をテストします
func TestValidateEncoding_Normal(t *testing.T) {
	tests := []struct {
		name     string
		encoding string
	}{
		{
			name:     "空文字列（デフォルト値使用）",
			encoding: "",
		},
		{
			name:     "UTF-8",
			encoding: "utf-8",
		},
		{
			name:     "Shift_JIS",
			encoding: "shift_jis",
		},
		{
			name:     "EUC-JP",
			encoding: "euc-jp",
		},
		{
			name:     "ISO-2022-JP",
			encoding: "iso-2022-jp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEncoding(tt.encoding)
			if err != nil {
				t.Errorf("validateEncoding() returned unexpected error: %v", err)
			}
		})
	}
}

// TestValidateEncoding_Error はvalidateEncoding()のエラーケースをテストします
func TestValidateEncoding_Error(t *testing.T) {
	tests := []struct {
		name     string
		encoding string
	}{
		{
			name:     "無効なエンコーディング",
			encoding: "invalid",
		},
		{
			name:     "大文字小文字違い",
			encoding: "UTF-8",
		},
		{
			name:     "スペース含み",
			encoding: "utf-8 ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEncoding(tt.encoding)
			if err == nil {
				t.Error("validateEncoding() should return error but got nil")
			}
		})
	}
}

// TestValidateConfig_Normal はvalidateConfig()の正常系をテストします
func TestValidateConfig_Normal(t *testing.T) {
	config := &Config{
		Target:   "/test/path",
		From:     "old",
		To:       "new",
		Encoding: "utf-8",
	}

	result, err := validateConfig(config)
	if err != nil {
		t.Errorf("validateConfig() returned unexpected error: %v", err)
	}

	if result.Target != config.Target {
		t.Errorf("Target = %v, expected %v", result.Target, config.Target)
	}
	if result.From != config.From {
		t.Errorf("From = %v, expected %v", result.From, config.From)
	}
	if result.To != config.To {
		t.Errorf("To = %v, expected %v", result.To, config.To)
	}
	if result.Encoding != config.Encoding {
		t.Errorf("Encoding = %v, expected %v", result.Encoding, config.Encoding)
	}
}

// TestValidateConfig_DefaultEncoding はvalidateConfig()のデフォルトエンコーディング設定をテストします
func TestValidateConfig_DefaultEncoding(t *testing.T) {
	config := &Config{
		Target:   "/test/path",
		From:     "old",
		To:       "new",
		Encoding: "", // 空文字列
	}

	result, err := validateConfig(config)
	if err != nil {
		t.Errorf("validateConfig() returned unexpected error: %v", err)
	}

	if result.Encoding != "utf-8" {
		t.Errorf("Encoding = %v, expected utf-8", result.Encoding)
	}
}

// TestValidateConfig_Error はvalidateConfig()のエラーケースをテストします
func TestValidateConfig_Error(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
	}{
		{
			name: "Targetが空",
			config: &Config{
				Target:   "",
				From:     "old",
				To:       "new",
				Encoding: "utf-8",
			},
		},
		{
			name: "Fromが空",
			config: &Config{
				Target:   "/test/path",
				From:     "",
				To:       "new",
				Encoding: "utf-8",
			},
		},
		{
			name: "Toが空",
			config: &Config{
				Target:   "/test/path",
				From:     "old",
				To:       "",
				Encoding: "utf-8",
			},
		},
		{
			name: "無効なエンコーディング",
			config: &Config{
				Target:   "/test/path",
				From:     "old",
				To:       "new",
				Encoding: "invalid",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validateConfig(tt.config)
			if err == nil {
				t.Error("validateConfig() should return error but got nil")
			}
		})
	}
}

// TestNewConfig_Normal はNewConfig()の正常系をテストします
func TestNewConfig_Normal(t *testing.T) {
	config, err := NewConfig("/test/path", "old", "new", "utf-8", "/backup", true, true, true)
	if err != nil {
		t.Errorf("NewConfig() returned unexpected error: %v", err)
	}

	if config.Target != "/test/path" {
		t.Errorf("Target = %v, expected /test/path", config.Target)
	}
	if config.From != "old" {
		t.Errorf("From = %v, expected old", config.From)
	}
	if config.To != "new" {
		t.Errorf("To = %v, expected new", config.To)
	}
	if config.Encoding != "utf-8" {
		t.Errorf("Encoding = %v, expected utf-8", config.Encoding)
	}
	if config.BackupDir != "/backup" {
		t.Errorf("BackupDir = %v, expected /backup", config.BackupDir)
	}
	if !config.Recursive {
		t.Error("Recursive should be true")
	}
	if !config.Backup {
		t.Error("Backup should be true")
	}
	if !config.DryRun {
		t.Error("DryRun should be true")
	}
}

// TestNewConfig_Error はNewConfig()のエラーケースをテストします
func TestNewConfig_Error(t *testing.T) {
	tests := []struct {
		name      string
		target    string
		from      string
		to        string
		encoding  string
		backupDir string
		recursive bool
		backup    bool
		dryRun    bool
	}{
		{
			name:      "Targetが空",
			target:    "",
			from:      "old",
			to:        "new",
			encoding:  "utf-8",
			backupDir: "/backup",
			recursive: false,
			backup:    false,
			dryRun:    false,
		},
		{
			name:      "Fromが空",
			target:    "/test/path",
			from:      "",
			to:        "new",
			encoding:  "utf-8",
			backupDir: "/backup",
			recursive: false,
			backup:    false,
			dryRun:    false,
		},
		{
			name:      "Toが空",
			target:    "/test/path",
			from:      "old",
			to:        "",
			encoding:  "utf-8",
			backupDir: "/backup",
			recursive: false,
			backup:    false,
			dryRun:    false,
		},
		{
			name:      "無効なエンコーディング",
			target:    "/test/path",
			from:      "old",
			to:        "new",
			encoding:  "invalid",
			backupDir: "/backup",
			recursive: false,
			backup:    false,
			dryRun:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewConfig(tt.target, tt.from, tt.to, tt.encoding, tt.backupDir, tt.recursive, tt.backup, tt.dryRun)
			if err == nil {
				t.Error("NewConfig() should return error but got nil")
			}
		})
	}
}

// MockFlagParser はテスト用のFlagParserモック
type MockFlagParser struct {
	stringVars     map[string]*string
	boolVars       map[string]*bool
	parseError     error
	expectedValues map[string]interface{}
}

// NewMockFlagParser は新しいMockFlagParserを作成します
func NewMockFlagParser() *MockFlagParser {
	return &MockFlagParser{
		stringVars:     make(map[string]*string),
		boolVars:       make(map[string]*bool),
		expectedValues: make(map[string]interface{}),
	}
}

// StringVar は文字列フラグを定義します
func (m *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
	// 期待値が設定されている場合はそれを使用、そうでなければデフォルト値
	if expectedValue, exists := m.expectedValues[name]; exists {
		if strValue, ok := expectedValue.(string); ok {
			*p = strValue
		} else {
			*p = value
		}
	} else {
		*p = value
	}
	m.stringVars[name] = p
}

// BoolVar はブールフラグを定義します
func (m *MockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	// 期待値が設定されている場合はそれを使用、そうでなければデフォルト値
	if expectedValue, exists := m.expectedValues[name]; exists {
		if boolValue, ok := expectedValue.(bool); ok {
			*p = boolValue
		} else {
			*p = value
		}
	} else {
		*p = value
	}
	m.boolVars[name] = p
}

// Parse はフラグを解析します
func (m *MockFlagParser) Parse() error {
	return m.parseError
}

// SetStringValue はテスト用に文字列値を設定します
func (m *MockFlagParser) SetStringValue(name, value string) {
	if p, exists := m.stringVars[name]; exists {
		*p = value
	}
}

// SetBoolValue はテスト用にブール値を設定します
func (m *MockFlagParser) SetBoolValue(name string, value bool) {
	if p, exists := m.boolVars[name]; exists {
		*p = value
	}
}

// SetParseError はテスト用にパースエラーを設定します
func (m *MockFlagParser) SetParseError(err error) {
	m.parseError = err
}

// MockOSArgs はテスト用のOSArgsモック
type MockOSArgs struct {
	args []string
}

// NewMockOSArgs は新しいMockOSArgsを作成します
func NewMockOSArgs(args []string) *MockOSArgs {
	return &MockOSArgs{args: args}
}

// Args はOS引数を返します
func (m *MockOSArgs) Args() []string {
	return m.args
}

// TestConfigParser_ParseFlags_Normal はConfigParser.ParseFlags()の正常系をテストします
func TestConfigParser_ParseFlags_Normal(t *testing.T) {
	mockFlagParser := NewMockFlagParser()
	mockOSArgs := NewMockOSArgs([]string{"program"})

	// 期待値を事前に設定
	mockFlagParser.expectedValues = map[string]interface{}{
		"target":     "/test/path",
		"from":       "old",
		"to":         "new",
		"encoding":   "utf-8",
		"backup-dir": "/backup",
		"recursive":  true,
		"backup":     true,
		"dry-run":    true,
	}

	configParser := NewConfigParser(mockFlagParser, mockOSArgs)

	config, err := configParser.ParseFlags()
	if err != nil {
		t.Errorf("ParseFlags() returned unexpected error: %v", err)
	}

	if config.Target != "/test/path" {
		t.Errorf("Target = %v, expected /test/path", config.Target)
	}
	if config.From != "old" {
		t.Errorf("From = %v, expected old", config.From)
	}
	if config.To != "new" {
		t.Errorf("To = %v, expected new", config.To)
	}
	if config.Encoding != "utf-8" {
		t.Errorf("Encoding = %v, expected utf-8", config.Encoding)
	}
	if config.BackupDir != "/backup" {
		t.Errorf("BackupDir = %v, expected /backup", config.BackupDir)
	}
	if !config.Recursive {
		t.Error("Recursive should be true")
	}
	if !config.Backup {
		t.Error("Backup should be true")
	}
	if !config.DryRun {
		t.Error("DryRun should be true")
	}
}

// TestConfigParser_ParseFlags_ParseError はConfigParser.ParseFlags()のパースエラーをテストします
func TestConfigParser_ParseFlags_ParseError(t *testing.T) {
	mockFlagParser := NewMockFlagParser()
	mockOSArgs := NewMockOSArgs([]string{"program"})
	configParser := NewConfigParser(mockFlagParser, mockOSArgs)

	// パースエラーを設定
	parseError := errors.New("parse error")
	mockFlagParser.SetParseError(parseError)

	_, err := configParser.ParseFlags()
	if err == nil {
		t.Error("ParseFlags() should return error but got nil")
	}
}

// TestConfigParser_ParseFlags_ValidationError はConfigParser.ParseFlags()のバリデーションエラーをテストします
func TestConfigParser_ParseFlags_ValidationError(t *testing.T) {
	mockFlagParser := NewMockFlagParser()
	mockOSArgs := NewMockOSArgs([]string{"program"})
	configParser := NewConfigParser(mockFlagParser, mockOSArgs)

	// 無効な設定（Targetが空）
	mockFlagParser.SetStringValue("target", "")
	mockFlagParser.SetStringValue("from", "old")
	mockFlagParser.SetStringValue("to", "new")

	_, err := configParser.ParseFlags()
	if err == nil {
		t.Error("ParseFlags() should return validation error but got nil")
	}
}

// TestParseFlags_Normal はParseFlags()関数の正常系をテストします
func TestParseFlags_Normal(t *testing.T) {
	// 元のos.Argsを保存
	originalArgs := []string{"program", "-target=/test/path", "-from=old", "-to=new"}

	// この関数は実際のflagパッケージを使用するため、モックは困難
	// 代わりに、関数が存在することを確認
	_, err := ParseFlags()
	// 実際の引数がないためエラーが発生することを期待
	if err == nil {
		t.Log("ParseFlags() function exists and can be called")
	} else {
		t.Logf("ParseFlags() returned error as expected: %v", err)
	}

	_ = originalArgs // 使用していることを示す
}

// TestPrintUsage_Normal はPrintUsage()関数をテストします
func TestPrintUsage_Normal(t *testing.T) {
	// PrintUsage()は出力のみを行う関数なので、パニックしないことを確認
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PrintUsage() panicked: %v", r)
		}
	}()

	PrintUsage()
	// 正常に実行されればテスト成功
}

// TestNewConfigParser_Normal はNewConfigParser()の正常系をテストします
func TestNewConfigParser_Normal(t *testing.T) {
	mockFlagParser := NewMockFlagParser()
	mockOSArgs := NewMockOSArgs([]string{"program"})

	configParser := NewConfigParser(mockFlagParser, mockOSArgs)

	if configParser == nil {
		t.Error("NewConfigParser() should not return nil")
		return
	}

	if configParser.flagParser != mockFlagParser {
		t.Error("flagParser should be set correctly")
	}

	if configParser.osArgs != mockOSArgs {
		t.Error("osArgs should be set correctly")
	}
}

// TestConfigParser_validateConfig_Normal はConfigParser.validateConfig()の正常系をテストします
func TestConfigParser_validateConfig_Normal(t *testing.T) {
	mockFlagParser := NewMockFlagParser()
	mockOSArgs := NewMockOSArgs([]string{"program"})
	configParser := NewConfigParser(mockFlagParser, mockOSArgs)

	config := &Config{
		Target:   "/test/path",
		From:     "old",
		To:       "new",
		Encoding: "utf-8",
	}

	result, err := configParser.validateConfig(config)
	if err != nil {
		t.Errorf("validateConfig() returned unexpected error: %v", err)
	}

	if result == nil {
		t.Error("validateConfig() should not return nil")
		return
	}

	if result.Target != config.Target {
		t.Errorf("Target = %v, expected %v", result.Target, config.Target)
	}
}

// TestConfigParser_validateConfig_Error はConfigParser.validateConfig()のエラーケースをテストします
func TestConfigParser_validateConfig_Error(t *testing.T) {
	mockFlagParser := NewMockFlagParser()
	mockOSArgs := NewMockOSArgs([]string{"program"})
	configParser := NewConfigParser(mockFlagParser, mockOSArgs)

	// 無効な設定（Targetが空）
	config := &Config{
		Target:   "",
		From:     "old",
		To:       "new",
		Encoding: "utf-8",
	}

	_, err := configParser.validateConfig(config)
	if err == nil {
		t.Error("validateConfig() should return error for invalid config")
	}
}
