package usecases

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// TestSecretDetectorService_StripProtocolPrefix_Normal はプロトコル識別子除去の正常系テスト
func TestSecretDetectorService_StripProtocolPrefix_Normal(t *testing.T) {
	mockExecutor := &MockCommandExecutor{}
	mockOutputWriter := &MockOutputWriter{}
	service := NewSecretDetectorService(false, false, "", mockExecutor, mockOutputWriter)

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

// TestSecretDetectorService_IsBinaryFile_Normal はバイナリファイル判定の正常系テスト
func TestSecretDetectorService_IsBinaryFile_Normal(t *testing.T) {
	mockExecutor := &MockCommandExecutor{}
	mockOutputWriter := &MockOutputWriter{}
	service := NewSecretDetectorService(false, false, "", mockExecutor, mockOutputWriter)

	testCases := []struct {
		name     string
		filename string
		expected bool
	}{
		// バイナリファイル（true期待）
		{
			name:     "JPEG image",
			filename: "image.jpg",
			expected: true,
		},
		{
			name:     "PNG image",
			filename: "logo.png",
			expected: true,
		},
		{
			name:     "PDF document",
			filename: "document.pdf",
			expected: true,
		},
		{
			name:     "Executable file",
			filename: "program.exe",
			expected: true,
		},
		{
			name:     "Shared library",
			filename: "library.so",
			expected: true,
		},
		{
			name:     "ZIP archive",
			filename: "archive.zip",
			expected: true,
		},
		{
			name:     "SQLite database",
			filename: "data.sqlite",
			expected: true,
		},
		// テキストファイル（false期待）
		{
			name:     "JSON config",
			filename: "config.json",
			expected: false,
		},
		{
			name:     "Go source",
			filename: "main.go",
			expected: false,
		},
		{
			name:     "Text file",
			filename: "README.txt",
			expected: false,
		},
		{
			name:     "Markdown file",
			filename: "README.md",
			expected: false,
		},
		{
			name:     "JavaScript file",
			filename: "app.js",
			expected: false,
		},
		{
			name:     "No extension",
			filename: "Dockerfile",
			expected: false,
		},
		{
			name:     "Empty filename",
			filename: "",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := service.IsBinaryFile(tc.filename)
			if result != tc.expected {
				t.Errorf("IsBinaryFile(%q) = %v, expected %v", tc.filename, result, tc.expected)
			}
		})
	}
}

// TestSecretDetectorService_CheckFileForHomePath_Normal はホームパス検知の正常系テスト
func TestSecretDetectorService_CheckFileForHomePath_Normal(t *testing.T) {
	// テスト用の一時ファイルを作成
	tempDir := t.TempDir()

	testCases := []struct {
		name              string
		filename          string
		content           string
		expectedCount     int
		expectedAllowed   int
		expectedForbidden int
	}{
		{
			name:              "File with forbidden home path",
			filename:          "test1.txt",
			content:           "This is a path: /home" + "/john/documents/file.txt\nAnother line",
			expectedCount:     1,
			expectedAllowed:   0,
			expectedForbidden: 1,
		},
		{
			name:              "File with allowed home/user path",
			filename:          "test2.txt",
			content:           "Default path: /home/user/config\nAnother line",
			expectedCount:     1,
			expectedAllowed:   1,
			expectedForbidden: 0,
		},
		{
			name:              "File with allowed [username] path",
			filename:          "test2b.txt",
			content:           "Template path: /home/[username]/config\nAnother line",
			expectedCount:     1,
			expectedAllowed:   1,
			expectedForbidden: 0,
		},
		{
			name:              "File with mixed home paths",
			filename:          "test3.txt",
			content:           "Forbidden: /home/alice/data\nAllowed: /home/user/settings\nAllowed: /home/[username]/template\nForbidden: /home" + "/bob/files",
			expectedCount:     4,
			expectedAllowed:   3,
			expectedForbidden: 1,
		},
		{
			name:              "File with no home paths",
			filename:          "test4.txt",
			content:           "Just regular content\nNo home paths here\n/var/log/app.log",
			expectedCount:     0,
			expectedAllowed:   0,
			expectedForbidden: 0,
		},
		{
			name:              "Empty file",
			filename:          "test5.txt",
			content:           "",
			expectedCount:     0,
			expectedAllowed:   0,
			expectedForbidden: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// テストファイルを作成
			testFile := filepath.Join(tempDir, tc.filename)
			err := os.WriteFile(testFile, []byte(tc.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			mockExecutor := &MockCommandExecutor{}
			mockOutputWriter := &MockOutputWriter{}
			service := NewSecretDetectorService(false, false, "", mockExecutor, mockOutputWriter)

			results, err := service.CheckFileForHomePath(testFile)
			if err != nil {
				t.Errorf("CheckFileForHomePath() unexpected error: %v", err)
				return
			}

			if len(results) != tc.expectedCount {
				t.Errorf("CheckFileForHomePath() returned %d results, expected %d", len(results), tc.expectedCount)
				return
			}

			allowedCount := 0
			forbiddenCount := 0
			for _, result := range results {
				if result.IsAllowed {
					allowedCount++
				} else {
					forbiddenCount++
				}
			}

			if allowedCount != tc.expectedAllowed {
				t.Errorf("CheckFileForHomePath() found %d allowed paths, expected %d", allowedCount, tc.expectedAllowed)
			}

			if forbiddenCount != tc.expectedForbidden {
				t.Errorf("CheckFileForHomePath() found %d forbidden paths, expected %d", forbiddenCount, tc.expectedForbidden)
			}
		})
	}
}

// TestSecretDetectorService_ExecuteSecretDetection_Normal はExecuteSecretDetectionの正常系テスト
func TestSecretDetectorService_ExecuteSecretDetection_Normal(t *testing.T) {
	testCases := []struct {
		name          string
		configFile    string
		mockOutput    string
		mockError     error
		expectedExit  int
		expectedError bool
	}{
		{
			name:          "No config files found",
			configFile:    "",
			mockOutput:    "",
			mockError:     nil,
			expectedExit:  0,
			expectedError: false,
		},
		{
			name:          "Specific config file specified - file not found",
			configFile:    "test_config.json",
			mockOutput:    "",
			mockError:     nil,
			expectedExit:  1,
			expectedError: true,
		},
		{
			name:          "Git command error",
			configFile:    "",
			mockOutput:    "",
			mockError:     errors.New("git command failed"),
			expectedExit:  1,
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
			mockOutputWriter := &MockOutputWriter{}
			service := NewSecretDetectorService(false, false, tc.configFile, mockExecutor, mockOutputWriter)

			exitCode, err := service.ExecuteSecretDetection()

			if tc.expectedError {
				if err == nil {
					t.Errorf("ExecuteSecretDetection() expected error, but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("ExecuteSecretDetection() unexpected error: %v", err)
				return
			}

			if exitCode != tc.expectedExit {
				t.Errorf("ExecuteSecretDetection() returned exit code %d, expected %d", exitCode, tc.expectedExit)
			}
		})
	}
}

// TestSecretDetectorService_IsPlaceholder_Normal はプレースホルダー判定の正常系テスト
func TestSecretDetectorService_IsPlaceholder_Normal(t *testing.T) {

	mockExecutor := &MockCommandExecutor{}
	mockOutputWriter := &MockOutputWriter{}
	service := NewSecretDetectorService(false, false, "", mockExecutor, mockOutputWriter)

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
	mockOutputWriter := &MockOutputWriter{}
	service := NewSecretDetectorService(false, false, "", mockExecutor, mockOutputWriter)

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
	mockOutputWriter := &MockOutputWriter{}
	service := NewSecretDetectorService(false, false, "", mockExecutor, mockOutputWriter)

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
	mockOutputWriter := &MockOutputWriter{}
	service := NewSecretDetectorService(false, false, "", mockExecutor, mockOutputWriter)

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
		{
			name: "Discord bot MCP settings file with full path",
			input: []string{
				"/home/user/devbox/.config/discord/mcp_settings_discord_bot_ayuayu.json",
				"other_file.txt",
				"README.md",
			},
			expected: []string{
				"/home/user/devbox/.config/discord/mcp_settings_discord_bot_ayuayu.json",
			},
		},
		{
			name: "MCP settings pattern matching variations",
			input: []string{
				"mcp_settings_discord_bot_ayuayu.json",
				"mcp_settings_another_bot.json",
				"mcp_settings.json",
				"not_mcp_settings.json",
				"regular_file.txt",
			},
			expected: []string{
				"mcp_settings_discord_bot_ayuayu.json",
				"mcp_settings_another_bot.json",
				"mcp_settings.json",
				"not_mcp_settings.json",
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
	mockOutputWriter := &MockOutputWriter{}
	service := NewSecretDetectorService(false, false, "", mockExecutor, mockOutputWriter)

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
			mockOutputWriter := &MockOutputWriter{}
			service := NewSecretDetectorService(tc.verbose, tc.dryRun, "", mockExecutor, mockOutputWriter)

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
