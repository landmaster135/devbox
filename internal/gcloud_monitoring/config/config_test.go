package config

import (
	"testing"
)

// MockFlagParser はテスト用のFlagParser実装
type MockFlagParser struct {
	stringVars map[string]*string
	boolVars   map[string]*bool
	args       []string
	parseError error
}

// NewMockFlagParser は新しいMockFlagParserを作成する
func NewMockFlagParser() *MockFlagParser {
	return &MockFlagParser{
		stringVars: make(map[string]*string),
		boolVars:   make(map[string]*bool),
		args:       []string{},
	}
}

// StringVar は文字列フラグを定義する
func (m *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
	*p = value
	m.stringVars[name] = p
}

// BoolVar はブールフラグを定義する
func (m *MockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	*p = value
	m.boolVars[name] = p
}

// Parse はフラグを解析する
func (m *MockFlagParser) Parse() error {
	return m.parseError
}

// Args は解析後の残りの引数を返す
func (m *MockFlagParser) Args() []string {
	return m.args
}

// SetStringValue はテスト用に文字列値を設定する
func (m *MockFlagParser) SetStringValue(name, value string) {
	if p, exists := m.stringVars[name]; exists {
		*p = value
	}
}

// SetBoolValue はテスト用にブール値を設定する
func (m *MockFlagParser) SetBoolValue(name string, value bool) {
	if p, exists := m.boolVars[name]; exists {
		*p = value
	}
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

func TestParseFlagsWithParser_ValidConfig(t *testing.T) {
	parser := NewMockFlagParser()

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Errorf("Unexpected error during initial parse: %v", err)
		return
	}

	// 有効なパラメータを設定
	parser.SetStringValue("operation", "create-dashboard-for-cloud-run")
	parser.SetStringValue("project", "test-project")
	parser.SetStringValue("location", "us-central1")
	parser.SetStringValue("service", "test-service")
	parser.SetStringValue("service-account-id", "test-sa")
	parser.SetBoolValue("help", false)

	cfg, err = ParseFlagsWithParser(parser)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
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
			parser.SetStringValue("operation", tt.operation)
			parser.SetStringValue("project", tt.project)
			parser.SetStringValue("location", tt.location)
			parser.SetStringValue("service", tt.service)
			parser.SetBoolValue("help", false)

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
	parser.SetBoolValue("help", true)

	cfg, err := ParseFlagsWithParser(parser)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !cfg.Help {
		t.Error("Expected Help to be true")
	}
}
