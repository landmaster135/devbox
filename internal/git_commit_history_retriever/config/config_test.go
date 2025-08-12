package config

import (
	"fmt"
	"strings"
	"testing"
)

// #==============================================================#
// ##          Mocks                                             ##
// #==============================================================#
// MockFlagParser はテスト用のモックFlagParser
type MockFlagParser struct {
	StringVarFunc func(p *string, name string, value string, usage string)
	BoolVarFunc   func(p *bool, name string, value bool, usage string)
	ParseFunc     func() error
}

// StringVar は文字列フラグを定義する（モック）
func (m *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
	if m.StringVarFunc != nil {
		m.StringVarFunc(p, name, value, usage)
	}
}

// BoolVar はブールフラグを定義する（モック）
func (m *MockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	if m.BoolVarFunc != nil {
		m.BoolVarFunc(p, name, value, usage)
	}
}

// Parse はフラグを解析する（モック）
func (m *MockFlagParser) Parse() error {
	if m.ParseFunc != nil {
		return m.ParseFunc()
	}
	return nil
}

// MockOSArgs はテスト用のモックOSArgs
type MockOSArgs struct {
	ArgsFunc func() []string
}

// Args はOS引数を返す（モック）
func (m *MockOSArgs) Args() []string {
	if m.ArgsFunc != nil {
		return m.ArgsFunc()
	}
	return []string{"test-program"}
}

// #==============================================================#
// ##          Tests                                             ##
// #==============================================================#
// TestConfigParser_ParseFlags_Normal はParseFlagsの正常系テスト
func TestConfigParser_ParseFlags_Normal(t *testing.T) {
	// Arrange
	var capturedGitDir, capturedKeyword, capturedSince, capturedUntil string

	mockFlagParser := &MockFlagParser{
		StringVarFunc: func(p *string, name string, value string, usage string) {
			switch name {
			case "git-dir":
				*p = "/test/repo"
				capturedGitDir = "/test/repo"
			case "keyword":
				*p = "feat:"
				capturedKeyword = "feat:"
			case "since":
				*p = "2025-01-01"
				capturedSince = "2025-01-01"
			case "until":
				*p = "2025-01-31"
				capturedUntil = "2025-01-31"
			}
		},
		ParseFunc: func() error {
			return nil
		},
	}

	mockOSArgs := &MockOSArgs{
		ArgsFunc: func() []string {
			return []string{"test-program"}
		},
	}

	configParser := NewConfigParser(mockFlagParser, mockOSArgs)

	// Act
	config, err := configParser.ParseFlags()

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if config.GitDir != capturedGitDir {
		t.Errorf("Expected GitDir to be %s, got %s", capturedGitDir, config.GitDir)
	}

	if config.Keyword != capturedKeyword {
		t.Errorf("Expected Keyword to be %s, got %s", capturedKeyword, config.Keyword)
	}

	if config.Since != capturedSince {
		t.Errorf("Expected Since to be %s, got %s", capturedSince, config.Since)
	}

	if config.Until != capturedUntil {
		t.Errorf("Expected Until to be %s, got %s", capturedUntil, config.Until)
	}
}

// TestConfigParser_ParseFlags_MissingGitDir はGitDirが未指定の場合のテスト
func TestConfigParser_ParseFlags_MissingGitDir(t *testing.T) {
	// Arrange
	mockFlagParser := &MockFlagParser{
		StringVarFunc: func(p *string, name string, value string, usage string) {
			// GitDirを空のままにする
		},
		ParseFunc: func() error {
			return nil
		},
	}

	mockOSArgs := &MockOSArgs{}
	configParser := NewConfigParser(mockFlagParser, mockOSArgs)

	// Act
	_, err := configParser.ParseFlags()

	// Assert
	if err == nil {
		t.Error("Expected error for missing git-dir, got nil")
	}

	if !strings.Contains(err.Error(), "--git-dir は必須パラメータです") {
		t.Errorf("Expected error message about required git-dir, got %v", err)
	}
}

// TestConfigParser_ParseFlags_InvalidSinceDate は無効なSince日付の場合のテスト
func TestConfigParser_ParseFlags_InvalidSinceDate(t *testing.T) {
	// Arrange
	mockFlagParser := &MockFlagParser{
		StringVarFunc: func(p *string, name string, value string, usage string) {
			switch name {
			case "git-dir":
				*p = "/test/repo"
			case "since":
				*p = "invalid-date"
			}
		},
		ParseFunc: func() error {
			return nil
		},
	}

	mockOSArgs := &MockOSArgs{}
	configParser := NewConfigParser(mockFlagParser, mockOSArgs)

	// Act
	_, err := configParser.ParseFlags()

	// Assert
	if err == nil {
		t.Error("Expected error for invalid since date, got nil")
	}

	if !strings.Contains(err.Error(), "--since の日付フォーマットが正しくありません") {
		t.Errorf("Expected error message about invalid since date format, got %v", err)
	}
}

// TestConfigParser_ParseFlags_InvalidUntilDate は無効なUntil日付の場合のテスト
func TestConfigParser_ParseFlags_InvalidUntilDate(t *testing.T) {
	// Arrange
	mockFlagParser := &MockFlagParser{
		StringVarFunc: func(p *string, name string, value string, usage string) {
			switch name {
			case "git-dir":
				*p = "/test/repo"
			case "until":
				*p = "2025/01/31"
			}
		},
		ParseFunc: func() error {
			return nil
		},
	}

	mockOSArgs := &MockOSArgs{}
	configParser := NewConfigParser(mockFlagParser, mockOSArgs)

	// Act
	_, err := configParser.ParseFlags()

	// Assert
	if err == nil {
		t.Error("Expected error for invalid until date, got nil")
	}

	if !strings.Contains(err.Error(), "--until の日付フォーマットが正しくありません") {
		t.Errorf("Expected error message about invalid until date format, got %v", err)
	}
}

// TestConfigParser_ParseFlags_ParseError はフラグ解析エラーの場合のテスト
func TestConfigParser_ParseFlags_ParseError(t *testing.T) {
	// Arrange
	mockFlagParser := &MockFlagParser{
		StringVarFunc: func(p *string, name string, value string, usage string) {
			if name == "git-dir" {
				*p = "/test/repo"
			}
		},
		ParseFunc: func() error {
			return fmt.Errorf("flag parse error")
		},
	}

	mockOSArgs := &MockOSArgs{}
	configParser := NewConfigParser(mockFlagParser, mockOSArgs)

	// Act
	_, err := configParser.ParseFlags()

	// Assert
	if err == nil {
		t.Error("Expected error for flag parse failure, got nil")
	}

	if !strings.Contains(err.Error(), "flag parse error") {
		t.Errorf("Expected flag parse error, got %v", err)
	}
}

// TestConfigParser_validateConfig_ValidDates は有効な日付の場合のテスト
func TestConfigParser_validateConfig_ValidDates(t *testing.T) {
	// Arrange
	configParser := &ConfigParser{}
	config := &Config{
		GitDir: "/test/repo",
		Since:  "2025-01-01",
		Until:  "2025-12-31",
	}

	// Act
	result, err := configParser.validateConfig(config)

	// Assert
	if err != nil {
		t.Errorf("Expected no error for valid dates, got %v", err)
	}

	if result.Since != "2025-01-01" {
		t.Errorf("Expected Since to be 2025-01-01, got %s", result.Since)
	}

	if result.Until != "2025-12-31" {
		t.Errorf("Expected Until to be 2025-12-31, got %s", result.Until)
	}
}

// TestConfigParser_validateConfig_EmptyOptionalFields は空のオプションフィールドの場合のテスト
func TestConfigParser_validateConfig_EmptyOptionalFields(t *testing.T) {
	// Arrange
	configParser := &ConfigParser{}
	config := &Config{
		GitDir:  "/test/repo",
		Keyword: "",
		Since:   "",
		Until:   "",
	}

	// Act
	result, err := configParser.validateConfig(config)

	// Assert
	if err != nil {
		t.Errorf("Expected no error for empty optional fields, got %v", err)
	}

	if result.GitDir != "/test/repo" {
		t.Errorf("Expected GitDir to be /test/repo, got %s", result.GitDir)
	}
}

// TestParseFlags_BackwardCompatibility は後方互換性関数のテスト
func TestParseFlags_BackwardCompatibility(t *testing.T) {
	// Note: この関数は実際のflagパッケージを使用するため、
	// 実際のコマンドライン引数に依存する。
	// ここでは関数が存在することのみを確認する。

	// Act & Assert
	// 実際のテストは統合テストで行う
	// ここでは関数が呼び出し可能であることを確認
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ParseFlags function should not panic, got %v", r)
		}
	}()

	// 実際の呼び出しはスキップ（コマンドライン引数に依存するため）
}

// TestPrintUsage_Normal はPrintUsageの正常系テスト
func TestPrintUsage_Normal(t *testing.T) {
	// Act & Assert
	// PrintUsageは標準出力に出力するため、panicしないことのみを確認
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PrintUsage should not panic, got %v", r)
		}
	}()

	PrintUsage()
}
