package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestAppConfig はAppConfigのテストクラスです
type TestAppConfig struct{}

// TestAppConfig_SetDefaults_Normal はSetDefaultsメソッドの正常系テストです
func (t *TestAppConfig) TestAppConfig_SetDefaults_Normal(test *testing.T) {
	// Arrange
	config := &ServiceConfig{}

	// Act
	config.SetDefaults()

	// Assert
	if config.BaseDir == "" {
		test.Error("BaseDir should not be empty after SetDefaults")
	}
	if config.CLIDir != "cmd/cli" {
		test.Errorf("Expected CLIDir to be 'cmd/cli', got '%s'", config.CLIDir)
	}
	if config.ScriptsDir != "scripts" {
		test.Errorf("Expected ScriptsDir to be 'scripts', got '%s'", config.ScriptsDir)
	}
	if config.OutputDir != "./pkg/bin/cli" {
		test.Errorf("Expected OutputDir to be './pkg/bin/cli', got '%s'", config.OutputDir)
	}
}

func TestAppConfig_SetDefaults_Normal(t *testing.T) {
	testInstance := &TestAppConfig{}
	testInstance.TestAppConfig_SetDefaults_Normal(t)
}

// TestAppConfig_SetDefaults_WithExistingValues は既存値がある場合のテストです
func (t *TestAppConfig) TestAppConfig_SetDefaults_WithExistingValues(test *testing.T) {
	// Arrange
	config := &ServiceConfig{
		BaseDir:    "/custom/base",
		CLIDir:     "custom/cli",
		ScriptsDir: "custom/scripts",
		OutputDir:  "custom/output",
	}

	// Act
	config.SetDefaults()

	// Assert
	if config.BaseDir != "/custom/base" {
		test.Errorf("Expected BaseDir to remain '/custom/base', got '%s'", config.BaseDir)
	}
	if config.CLIDir != "custom/cli" {
		test.Errorf("Expected CLIDir to remain 'custom/cli', got '%s'", config.CLIDir)
	}
	if config.ScriptsDir != "custom/scripts" {
		test.Errorf("Expected ScriptsDir to remain 'custom/scripts', got '%s'", config.ScriptsDir)
	}
	if config.OutputDir != "custom/output" {
		test.Errorf("Expected OutputDir to remain 'custom/output', got '%s'", config.OutputDir)
	}
}

func TestAppConfig_SetDefaults_WithExistingValues(t *testing.T) {
	testInstance := &TestAppConfig{}
	testInstance.TestAppConfig_SetDefaults_WithExistingValues(t)
}

// TestAppConfig_SetDefaults_BaseDirFallback はBaseDir設定時のフォールバック処理テストです
func (t *TestAppConfig) TestAppConfig_SetDefaults_BaseDirFallback(test *testing.T) {
	// Arrange
	config := &ServiceConfig{}

	// 現在のワーキングディレクトリを取得して期待値を設定
	expectedBaseDir, err := os.Getwd()
	if err != nil {
		expectedBaseDir = "."
	}

	// Act
	config.SetDefaults()

	// Assert
	if config.BaseDir != expectedBaseDir {
		test.Errorf("Expected BaseDir to be '%s', got '%s'", expectedBaseDir, config.BaseDir)
	}
}

func TestAppConfig_SetDefaults_BaseDirFallback(t *testing.T) {
	testInstance := &TestAppConfig{}
	testInstance.TestAppConfig_SetDefaults_BaseDirFallback(t)
}

// TestAppConfig_GetCLIPath_Normal はGetCLIPathメソッドの正常系テストです
func (t *TestAppConfig) TestAppConfig_GetCLIPath_Normal(test *testing.T) {
	// Arrange
	config := &ServiceConfig{
		BaseDir: "/home/user/project",
		CLIDir:  "cmd/cli",
	}
	expected := filepath.Join("/home/user/project", "cmd/cli")

	// Act
	result := config.GetCLIPath()

	// Assert
	if result != expected {
		test.Errorf("Expected GetCLIPath to return '%s', got '%s'", expected, result)
	}
}

func TestAppConfig_GetCLIPath_Normal(t *testing.T) {
	testInstance := &TestAppConfig{}
	testInstance.TestAppConfig_GetCLIPath_Normal(t)
}

// TestAppConfig_GetCLIPath_WithDifferentPaths は異なるパスでのGetCLIPathテストです
func (t *TestAppConfig) TestAppConfig_GetCLIPath_WithDifferentPaths(test *testing.T) {
	// Arrange
	testCases := []struct {
		name     string
		baseDir  string
		cliDir   string
		expected string
	}{
		{
			name:     "Unix style paths",
			baseDir:  "/usr/local/project",
			cliDir:   "bin/cli",
			expected: filepath.Join("/usr/local/project", "bin/cli"),
		},
		{
			name:     "Relative paths",
			baseDir:  ".",
			cliDir:   "cmd/cli",
			expected: filepath.Join(".", "cmd/cli"),
		},
		{
			name:     "Empty CLI dir",
			baseDir:  "/project",
			cliDir:   "",
			expected: filepath.Join("/project", ""),
		},
	}

	for _, tc := range testCases {
		test.Run(tc.name, func(t *testing.T) {
			// Arrange
			config := &ServiceConfig{
				BaseDir: tc.baseDir,
				CLIDir:  tc.cliDir,
			}

			// Act
			result := config.GetCLIPath()

			// Assert
			if result != tc.expected {
				t.Errorf("Expected GetCLIPath to return '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestAppConfig_GetCLIPath_WithDifferentPaths(t *testing.T) {
	testInstance := &TestAppConfig{}
	testInstance.TestAppConfig_GetCLIPath_WithDifferentPaths(t)
}

// TestAppConfig_GetScriptsPath_Normal はGetScriptsPathメソッドの正常系テストです
func (t *TestAppConfig) TestAppConfig_GetScriptsPath_Normal(test *testing.T) {
	// Arrange
	config := &ServiceConfig{
		BaseDir:    "/home/user/project",
		ScriptsDir: "scripts",
	}
	expected := filepath.Join("/home/user/project", "scripts")

	// Act
	result := config.GetScriptsPath()

	// Assert
	if result != expected {
		test.Errorf("Expected GetScriptsPath to return '%s', got '%s'", expected, result)
	}
}

func TestAppConfig_GetScriptsPath_Normal(t *testing.T) {
	testInstance := &TestAppConfig{}
	testInstance.TestAppConfig_GetScriptsPath_Normal(t)
}

// TestAppConfig_GetScriptsPath_WithDifferentPaths は異なるパスでのGetScriptsPathテストです
func (t *TestAppConfig) TestAppConfig_GetScriptsPath_WithDifferentPaths(test *testing.T) {
	// Arrange
	testCases := []struct {
		name       string
		baseDir    string
		scriptsDir string
		expected   string
	}{
		{
			name:       "Unix style paths",
			baseDir:    "/usr/local/project",
			scriptsDir: "bin/scripts",
			expected:   filepath.Join("/usr/local/project", "bin/scripts"),
		},
		{
			name:       "Relative paths",
			baseDir:    ".",
			scriptsDir: "scripts",
			expected:   filepath.Join(".", "scripts"),
		},
		{
			name:       "Empty scripts dir",
			baseDir:    "/project",
			scriptsDir: "",
			expected:   filepath.Join("/project", ""),
		},
	}

	for _, tc := range testCases {
		test.Run(tc.name, func(t *testing.T) {
			// Arrange
			config := &ServiceConfig{
				BaseDir:    tc.baseDir,
				ScriptsDir: tc.scriptsDir,
			}

			// Act
			result := config.GetScriptsPath()

			// Assert
			if result != tc.expected {
				t.Errorf("Expected GetScriptsPath to return '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestAppConfig_GetScriptsPath_WithDifferentPaths(t *testing.T) {
	testInstance := &TestAppConfig{}
	testInstance.TestAppConfig_GetScriptsPath_WithDifferentPaths(t)
}

// TestAppConfig_Integration は統合テストです
func (t *TestAppConfig) TestAppConfig_Integration(test *testing.T) {
	// Arrange
	config := &ServiceConfig{
		PackageName: "test-package",
		ShowHelp:    true,
	}

	// Act
	config.SetDefaults()
	cliPath := config.GetCLIPath()
	scriptsPath := config.GetScriptsPath()

	// Assert
	if config.PackageName != "test-package" {
		test.Errorf("Expected PackageName to remain 'test-package', got '%s'", config.PackageName)
	}
	if !config.ShowHelp {
		test.Error("Expected ShowHelp to remain true")
	}
	if config.BaseDir == "" {
		test.Error("BaseDir should be set after SetDefaults")
	}
	if cliPath == "" {
		test.Error("GetCLIPath should return non-empty path")
	}
	if scriptsPath == "" {
		test.Error("GetScriptsPath should return non-empty path")
	}
}

func TestAppConfig_Integration(t *testing.T) {
	testInstance := &TestAppConfig{}
	testInstance.TestAppConfig_Integration(t)
}

// TestAppConfig_SetDefaults_ErrorHandling はSetDefaultsのエラーハンドリングテストです
func (t *TestAppConfig) TestAppConfig_SetDefaults_ErrorHandling(test *testing.T) {
	// Arrange
	config := &ServiceConfig{
		BaseDir: "", // 空文字列でテスト
	}

	// Act
	config.SetDefaults()

	// Assert
	// エラーが発生してもフォールバック値が設定されることを確認
	if config.BaseDir == "" {
		test.Error("BaseDir should be set to fallback value even when os.Getwd() might fail")
	}
	// デフォルト値が正しく設定されることを確認
	if config.CLIDir != "cmd/cli" {
		test.Errorf("Expected CLIDir to be 'cmd/cli', got '%s'", config.CLIDir)
	}
	if config.ScriptsDir != "scripts" {
		test.Errorf("Expected ScriptsDir to be 'scripts', got '%s'", config.ScriptsDir)
	}
	if config.OutputDir != "./pkg/bin/cli" {
		test.Errorf("Expected OutputDir to be './pkg/bin/cli', got '%s'", config.OutputDir)
	}
}

// TestAppConfig_SetDefaults_ErrorHandling はSetDefaultsのエラーハンドリングテストです
func TestAppConfig_SetDefaults_ErrorHandling(t *testing.T) {
	testInstance := &TestAppConfig{}
	testInstance.TestAppConfig_SetDefaults_ErrorHandling(t)
}

// TestAppConfig_GetPaths_EdgeCases はパス取得メソッドのエッジケーステストです
func (t *TestAppConfig) TestAppConfig_GetPaths_EdgeCases(test *testing.T) {
	// Arrange
	testCases := []struct {
		name     string
		config   *ServiceConfig
		testType string
		expected string
	}{
		{
			name: "Special characters in BaseDir",
			config: &ServiceConfig{
				BaseDir: "/path/with spaces/and-special_chars",
				CLIDir:  "cmd/cli",
			},
			testType: "cli",
			expected: filepath.Join("/path/with spaces/and-special_chars", "cmd/cli"),
		},
		{
			name: "Special characters in ScriptsDir",
			config: &ServiceConfig{
				BaseDir:    "/base",
				ScriptsDir: "scripts-with_special.chars",
			},
			testType: "scripts",
			expected: filepath.Join("/base", "scripts-with_special.chars"),
		},
		{
			name: "Very long paths",
			config: &ServiceConfig{
				BaseDir: "/very/long/path/that/might/cause/issues/in/some/systems/but/should/work/fine",
				CLIDir:  "cmd/cli/with/nested/structure",
			},
			testType: "cli",
			expected: filepath.Join("/very/long/path/that/might/cause/issues/in/some/systems/but/should/work/fine", "cmd/cli/with/nested/structure"),
		},
		{
			name: "Dot paths",
			config: &ServiceConfig{
				BaseDir:    ".",
				ScriptsDir: "./scripts",
			},
			testType: "scripts",
			expected: filepath.Join(".", "./scripts"),
		},
	}

	for _, tc := range testCases {
		test.Run(tc.name, func(t *testing.T) {
			// Act & Assert
			var result string
			if tc.testType == "cli" {
				result = tc.config.GetCLIPath()
			} else {
				result = tc.config.GetScriptsPath()
			}

			if result != tc.expected {
				t.Errorf("Expected path '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

// TestAppConfig_GetPaths_EdgeCases はパス取得メソッドのエッジケーステストです
func TestAppConfig_GetPaths_EdgeCases(t *testing.T) {
	testInstance := &TestAppConfig{}
	testInstance.TestAppConfig_GetPaths_EdgeCases(t)
}
