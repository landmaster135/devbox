package usecases

import (
	"errors"
	"testing"
)

// TestSecretDetectorService はSecretDetectorServiceのテストクラス
type TestSecretDetectorService struct{}

// TestSecretDetectorService_StripProtocolPrefix_Normal はプロトコル識別子除去の正常系テスト
func TestSecretDetectorService_StripProtocolPrefix_Normal(t *testing.T) {
	mockExecutor := &MockCommandExecutor{}
	service := NewSecretDetectorService(false, false, mockExecutor)

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "PostgreSQL URL with placeholder",
			input:    "postgresql://YOUR_URL",
			expected: "YOUR_URL",
		},
		{
			name:     "HTTPS URL with placeholder",
			input:    "https://YOUR_API_ENDPOINT",
			expected: "YOUR_API_ENDPOINT",
		},
		{
			name:     "Redis URL with placeholder",
			input:    "redis://YOUR_REDIS_URL",
			expected: "YOUR_REDIS_URL",
		},
		{
			name:     "MySQL URL with placeholder",
			input:    "mysql://YOUR_DATABASE_URL",
			expected: "YOUR_DATABASE_URL",
		},
		{
			name:     "MongoDB URL with placeholder",
			input:    "mongodb://YOUR_MONGO_URL",
			expected: "YOUR_MONGO_URL",
		},
		{
			name:     "No protocol prefix",
			input:    "YOUR_API_KEY",
			expected: "YOUR_API_KEY",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Real PostgreSQL URL",
			input:    "postgresql://user:pass@localhost:5432/mydb",
			expected: "user:pass@localhost:5432/mydb",
		},
		{
			name:     "Real HTTPS URL",
			input:    "https://api.example.com/v1",
			expected: "api.example.com/v1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := service.stripProtocolPrefix(tc.input)
			if result != tc.expected {
				t.Errorf("stripProtocolPrefix(%q) = %q, expected %q", tc.input, result, tc.expected)
			}
		})
	}
}

// TestSecretDetectorService_IsPlaceholder_Normal はプレースホルダー判定の正常系テスト
func TestSecretDetectorService_IsPlaceholder_Normal(t *testing.T) {

	mockExecutor := &MockCommandExecutor{}
	service := NewSecretDetectorService(false, false, mockExecutor)

	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		// プロトコル付きプレースホルダー（改修対象）
		{
			name:     "PostgreSQL with YOUR_URL",
			input:    "postgresql://YOUR_URL",
			expected: true,
		},
		{
			name:     "HTTPS with YOUR_API_ENDPOINT",
			input:    "https://YOUR_API_ENDPOINT",
			expected: true,
		},
		{
			name:     "Redis with YOUR_REDIS_URL",
			input:    "redis://YOUR_REDIS_URL",
			expected: true,
		},
		{
			name:     "MySQL with YOUR_DATABASE_URL",
			input:    "mysql://YOUR_DATABASE_URL",
			expected: true,
		},
		// 従来のプレースホルダー
		{
			name:     "YOUR_API_KEY",
			input:    "YOUR_API_KEY",
			expected: true,
		},
		{
			name:     "YOUR_SECRET_KEY",
			input:    "YOUR_SECRET_KEY",
			expected: true,
		},
		{
			name:     "PLACEHOLDER",
			input:    "PLACEHOLDER",
			expected: true,
		},
		{
			name:     "EXAMPLE_KEY",
			input:    "EXAMPLE_KEY",
			expected: true,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: true,
		},
		{
			name:     "Short value",
			input:    "test",
			expected: true,
		},
		{
			name:     "Test value",
			input:    "test-api-key",
			expected: true,
		},
		{
			name:     "Demo value",
			input:    "demo-secret",
			expected: true,
		},
		// 実際のシークレット（false期待）
		{
			name:     "Real PostgreSQL URL",
			input:    "postgresql://user:pass@localhost:5432/mydb",
			expected: false,
		},
		{
			name:     "Real shorten PostgreSQL URL",
			input:    "postgres://user:pass@localhost:5432/mydb",
			expected: false,
		},
		{
			name:     "Real short HTTPS URL",
			input:    "https://api.example.com/v1/users",
			expected: true,
		},
		{
			name:     "Real long HTTPS URL",
			input:    "https://api.production123456789.example.com/v1/users/12345",
			expected: true, // "example"がテストパターンにマッチするためプレースホルダーとして判定される
		},
		{
			name:     "Real long HTTPS URL without test keywords",
			input:    "https://api.production123456789.mycompany.com/v1/users/12345",
			expected: false,
		},
		{
			name:     "Real API key pattern",
			input:    "sk-1234567890abcdef1234567890abcdef1234567890abcdef",
			expected: false,
		},
		{
			name:     "Real GitHub token",
			input:    "ghp_1234567890abcdef1234567890abcdef123456",
			expected: false,
		},
		{
			name:     "Long random string",
			input:    "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := service.IsPlaceholder(tc.input)
			if result != tc.expected {
				t.Errorf("IsPlaceholder(%q) = %v, expected %v", tc.input, result, tc.expected)
			}
		})
	}
}

// TestSecretDetectorService_IsSuspiciousKey_Normal は疑わしいキー名判定の正常系テスト
func TestSecretDetectorService_IsSuspiciousKey_Normal(t *testing.T) {

	mockExecutor := &MockCommandExecutor{}
	service := NewSecretDetectorService(false, false, mockExecutor)

	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		// 疑わしいキー名（true期待）
		{
			name:     "API_KEY",
			input:    "API_KEY",
			expected: true,
		},
		{
			name:     "SECRET_KEY",
			input:    "SECRET_KEY",
			expected: true,
		},
		{
			name:     "ACCESS_TOKEN",
			input:    "ACCESS_TOKEN",
			expected: true,
		},
		{
			name:     "DATABASE_URL",
			input:    "DATABASE_URL",
			expected: true,
		},
		{
			name:     "WEBHOOK_URL",
			input:    "WEBHOOK_URL",
			expected: true,
		},
		{
			name:     "PASSWORD",
			input:    "PASSWORD",
			expected: true,
		},
		{
			name:     "api-key with hyphen",
			input:    "api-key",
			expected: true,
		},
		{
			name:     "secret_key with underscore",
			input:    "secret_key",
			expected: true,
		},
		// 通常のキー名（false期待）
		{
			name:     "COMMAND",
			input:    "COMMAND",
			expected: false,
		},
		{
			name:     "ARGS",
			input:    "ARGS",
			expected: false,
		},
		{
			name:     "PORT",
			input:    "PORT",
			expected: false,
		},
		{
			name:     "HOST",
			input:    "HOST",
			expected: false,
		},
		{
			name:     "DEBUG",
			input:    "DEBUG",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := service.IsSuspiciousKey(tc.input)
			if result != tc.expected {
				t.Errorf("IsSuspiciousKey(%q) = %v, expected %v", tc.input, result, tc.expected)
			}
		})
	}
}

// TestSecretDetectorService_MatchesSecretPattern_Normal はシークレットパターン判定の正常系テスト
func TestSecretDetectorService_MatchesSecretPattern_Normal(t *testing.T) {

	mockExecutor := &MockCommandExecutor{}
	service := NewSecretDetectorService(false, false, mockExecutor)

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		// 実際のシークレットパターン
		{
			name:     "OpenAI API key",
			input:    "sk-1234567890abcdef1234567890abcdef1234567890abcdef",
			expected: `sk-[a-zA-Z0-9]{48}`,
		},
		{
			name:     "GitHub PAT",
			input:    "ghp_1234567890abcdef1234567890abcdef123456",
			expected: `ghp_[a-zA-Z0-9]{36}`,
		},
		{
			name:     "Google API key",
			input:    "AIza1234567890abcdefghijklmnopqrstuvwxyz",
			expected: `AIza[0-9A-Za-z_-]{35}`,
		},
		{
			name:     "AWS Access Key",
			input:    "AKIA1234567890ABCDEF",
			expected: `AKIA[0-9A-Z]{16}`,
		},
		{
			name:     "UUID format",
			input:    "12345678-1234-1234-1234-123456789abc",
			expected: `[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`,
		},
		// パターンにマッチしない値
		{
			name:     "Regular string",
			input:    "regular-string-value",
			expected: "",
		},
		{
			name:     "YOUR_API_KEY placeholder",
			input:    "YOUR_API_KEY",
			expected: "",
		},
		{
			name:     "Short string",
			input:    "short",
			expected: "",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := service.MatchesSecretPattern(tc.input)
			if result != tc.expected {
				t.Errorf("MatchesSecretPattern(%q) = %q, expected %q", tc.input, result, tc.expected)
			}
		})
	}
}

// TestSecretDetectorService_FilterConfigFiles_Normal は設定ファイルフィルタリングの正常系テスト
func TestSecretDetectorService_FilterConfigFiles_Normal(t *testing.T) {

	mockExecutor := &MockCommandExecutor{}
	service := NewSecretDetectorService(false, false, mockExecutor)

	testCases := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name: "Mixed files with config files",
			input: []string{
				"config.json",
				"settings.json",
				"main.go",
				"README.md",
				"cline_mcp_settings.json",
				"package.json",
				"app.config.js",
			},
			expected: []string{
				"config.json",
				"settings.json",
				"cline_mcp_settings.json",
				"package.json",
				"app.config.js",
			},
		},
		{
			name:     "No config files",
			input:    []string{"main.go", "README.md", "Dockerfile"},
			expected: []string{},
		},
		{
			name:     "Empty input",
			input:    []string{},
			expected: []string{},
		},
		{
			name: "Only config files",
			input: []string{
				"claude_desktop_config.json",
				"mcp_settings.json",
				"app.config.ts",
			},
			expected: []string{
				"claude_desktop_config.json",
				"mcp_settings.json",
				"app.config.ts",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := service.FilterConfigFiles(tc.input)
			if len(result) != len(tc.expected) {
				t.Errorf("FilterConfigFiles(%v) returned %d files, expected %d", tc.input, len(result), len(tc.expected))
				return
			}
			for i, expected := range tc.expected {
				if result[i] != expected {
					t.Errorf("FilterConfigFiles(%v)[%d] = %q, expected %q", tc.input, i, result[i], expected)
				}
			}
		})
	}
}

// TestSecretDetectorService_CalculateEntropy_Normal はエントロピー計算の正常系テスト
func TestSecretDetectorService_CalculateEntropy_Normal(t *testing.T) {

	mockExecutor := &MockCommandExecutor{}
	service := NewSecretDetectorService(false, false, mockExecutor)

	testCases := []struct {
		name     string
		input    string
		minValue float64
		maxValue float64
	}{
		{
			name:     "Empty string",
			input:    "",
			minValue: 0.0,
			maxValue: 0.0,
		},
		{
			name:     "Single character repeated",
			input:    "aaaaaaaaaa",
			minValue: 0.0,
			maxValue: 0.1,
		},
		{
			name:     "Low entropy string",
			input:    "password123",
			minValue: 2.0,
			maxValue: 4.0,
		},
		{
			name:     "High entropy string",
			input:    "sk-1234567890abcdef1234567890abcdef1234567890abcdef",
			minValue: 4.0,
			maxValue: 6.0,
		},
		{
			name:     "Random-like string",
			input:    "a1B2c3D4e5F6g7H8i9J0k1L2m3N4o5P6",
			minValue: 4.5,
			maxValue: 6.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := service.CalculateEntropy(tc.input)
			if result < tc.minValue || result > tc.maxValue {
				t.Errorf("CalculateEntropy(%q) = %f, expected between %f and %f", tc.input, result, tc.minValue, tc.maxValue)
			}
		})
	}
}

// TestSecretDetectorService_GetStagedFiles_Normal はGitステージファイル取得の正常系テスト
func TestSecretDetectorService_GetStagedFiles_Normal(t *testing.T) {
	testCases := []struct {
		name          string
		dryRun        bool
		verbose       bool
		mockOutput    string
		mockError     error
		expectedFiles []string
		expectedError bool
	}{
		{
			name:          "Normal case with multiple files",
			dryRun:        false,
			verbose:       false,
			mockOutput:    "config.json\nsettings.json\napp.config.js\n",
			mockError:     nil,
			expectedFiles: []string{"config.json", "settings.json", "app.config.js"},
			expectedError: false,
		},
		{
			name:          "Single file",
			dryRun:        false,
			verbose:       false,
			mockOutput:    "config.json\n",
			mockError:     nil,
			expectedFiles: []string{"config.json"},
			expectedError: false,
		},
		{
			name:          "No files",
			dryRun:        false,
			verbose:       false,
			mockOutput:    "",
			mockError:     nil,
			expectedFiles: []string{},
			expectedError: false,
		},
		{
			name:          "Empty lines filtered out",
			dryRun:        false,
			verbose:       false,
			mockOutput:    "config.json\n\nsettings.json\n\n",
			mockError:     nil,
			expectedFiles: []string{"config.json", "settings.json"},
			expectedError: false,
		},
		{
			name:          "Dry run mode",
			dryRun:        true,
			verbose:       false,
			mockOutput:    "",
			mockError:     nil,
			expectedFiles: []string{},
			expectedError: false,
		},
		{
			name:          "Command execution error",
			dryRun:        false,
			verbose:       false,
			mockOutput:    "",
			mockError:     errors.New("git command failed"),
			expectedFiles: nil,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockExecutor := &MockCommandExecutor{
				ExecuteFunc: func(name string, args ...string) ([]byte, error) {
					if tc.mockError != nil {
						return nil, tc.mockError
					}
					return []byte(tc.mockOutput), nil
				},
			}
			service := NewSecretDetectorService(tc.verbose, tc.dryRun, mockExecutor)

			result, err := service.GetStagedFiles()

			if tc.expectedError {
				if err == nil {
					t.Errorf("GetStagedFiles() expected error, but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("GetStagedFiles() unexpected error: %v", err)
				return
			}

			if len(result) != len(tc.expectedFiles) {
				t.Errorf("GetStagedFiles() returned %d files, expected %d", len(result), len(tc.expectedFiles))
				return
			}

			for i, expected := range tc.expectedFiles {
				if result[i] != expected {
					t.Errorf("GetStagedFiles()[%d] = %q, expected %q", i, result[i], expected)
				}
			}
		})
	}
}
