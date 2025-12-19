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

func TestGetServiceAccountEmail(t *testing.T) {
	tests := []struct {
		name             string
		project          string
		serviceAccountID string
		expected         string
	}{
		{
			name:             "With service account ID",
			project:          "test-project",
			serviceAccountID: "test-sa",
			expected:         "test-sa@test-project.iam.gserviceaccount.com",
		},
		{
			name:             "Without service account ID",
			project:          "test-project",
			serviceAccountID: "",
			expected:         "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Project:          tt.project,
				ServiceAccountID: tt.serviceAccountID,
			}
			result := cfg.GetServiceAccountEmail()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestParseFlagsWithParser_Normal(t *testing.T) {
	parser := NewMockFlagParser()

	// 事前に値を設定
	parser.SetStringFlag("operation", "create-dashboard-for-cloud-run")
	parser.SetStringFlag("project", "test-project")
	parser.SetStringFlag("location", "us-central1")
	parser.SetStringFlag("service", "test-service")
	parser.SetStringFlag("service-account-id", "test-sa")
	parser.SetBoolFlag("help", false)

	cfg, err := ParseFlagsWithParser(parser)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if cfg == nil {
		t.Fatal("Expected config to be non-nil")
		return
	}

	if cfg.Operation != "create-dashboard-for-cloud-run" {
		t.Errorf("Expected operation 'create-dashboard-for-cloud-run', got '%s'", cfg.Operation)
	}

	if cfg.Project != "test-project" {
		t.Errorf("Expected project 'test-project', got '%s'", cfg.Project)
	}

	if cfg.Location != "us-central1" {
		t.Errorf("Expected location 'us-central1', got '%s'", cfg.Location)
	}

	if cfg.Service != "test-service" {
		t.Errorf("Expected service 'test-service', got '%s'", cfg.Service)
	}

	if cfg.ServiceAccountID != "test-sa" {
		t.Errorf("Expected service account ID 'test-sa', got '%s'", cfg.ServiceAccountID)
	}

	if cfg.Help {
		t.Error("Expected Help to be false")
	}
}

func TestParseFlagsWithParser_MissingRequiredParams(t *testing.T) {
	tests := []struct {
		name        string
		operation   string
		project     string
		location    string
		service     string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Missing operation",
			operation:   "",
			project:     "test",
			location:    "us-central1",
			service:     "test-service",
			expectError: true,
			errorMsg:    "operation パラメータは必須です",
		},
		{
			name:        "Missing project",
			operation:   "create-dashboard-for-cloud-run",
			project:     "",
			location:    "us-central1",
			service:     "test-service",
			expectError: true,
			errorMsg:    "project パラメータは必須です",
		},
		{
			name:        "Missing location",
			operation:   "create-dashboard-for-cloud-run",
			project:     "test",
			location:    "",
			service:     "test-service",
			expectError: true,
			errorMsg:    "location パラメータは必須です",
		},
		{
			name:        "Missing service",
			operation:   "create-dashboard-for-cloud-run",
			project:     "test",
			location:    "us-central1",
			service:     "",
			expectError: true,
			errorMsg:    "service パラメータは必須です",
		},
		{
			name:        "Invalid operation",
			operation:   "invalid-operation",
			project:     "test",
			location:    "us-central1",
			service:     "test-service",
			expectError: true,
			errorMsg:    "未対応の操作です: invalid-operation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewMockFlagParser()
			parser.SetStringFlag("operation", tt.operation)
			parser.SetStringFlag("project", tt.project)
			parser.SetStringFlag("location", tt.location)
			parser.SetStringFlag("service", tt.service)
			parser.SetBoolFlag("help", false)

			_, err := ParseFlagsWithParser(parser)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
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

func TestParseFlagsWithParser_HelpFlag(t *testing.T) {
	parser := NewMockFlagParser()
	parser.SetBoolFlag("help", true)

	cfg, err := ParseFlagsWithParser(parser)

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

func TestParseFlags_Normal(t *testing.T) {
	// 元のos.Argsを保存
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// テスト用の引数を設定
	os.Args = []string{
		"test-program",
		"-operation=create-dashboard-for-cloud-run",
		"-project=test-project",
		"-location=us-central1",
		"-service=test-service",
		"-service-account-id=test-sa",
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

	if cfg.Operation != "create-dashboard-for-cloud-run" {
		t.Errorf("Expected operation 'create-dashboard-for-cloud-run', got '%s'", cfg.Operation)
	}

	if cfg.Project != "test-project" {
		t.Errorf("Expected project 'test-project', got '%s'", cfg.Project)
	}

	if cfg.Location != "us-central1" {
		t.Errorf("Expected location 'us-central1', got '%s'", cfg.Location)
	}

	if cfg.Service != "test-service" {
		t.Errorf("Expected service 'test-service', got '%s'", cfg.Service)
	}

	if cfg.ServiceAccountID != "test-sa" {
		t.Errorf("Expected service account ID 'test-sa', got '%s'", cfg.ServiceAccountID)
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
