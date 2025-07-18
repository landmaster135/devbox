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
	config := &AppConfig{}

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

// TestAppConfig_SetDefaults_WithExistingValues は既存値がある場合のテストです
func (t *TestAppConfig) TestAppConfig_SetDefaults_WithExistingValues(test *testing.T) {
	// Arrange
	config := &AppConfig{
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

// TestAppConfig_SetDefaults_BaseDirFallback はBaseDir設定時のフォールバック処理テストです
func (t *TestAppConfig) TestAppConfig_SetDefaults_BaseDirFallback(test *testing.T) {
	// Arrange
	config := &AppConfig{}

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

// TestAppConfig_GetCLIPath_Normal はGetCLIPathメソッドの正常系テストです
func (t *TestAppConfig) TestAppConfig_GetCLIPath_Normal(test *testing.T) {
	// Arrange
	config := &AppConfig{
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
			config := &AppConfig{
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

// TestAppConfig_GetScriptsPath_Normal はGetScriptsPathメソッドの正常系テストです
func (t *TestAppConfig) TestAppConfig_GetScriptsPath_Normal(test *testing.T) {
	// Arrange
	config := &AppConfig{
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
			config := &AppConfig{
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

// TestAppConfig_Integration は統合テストです
func (t *TestAppConfig) TestAppConfig_Integration(test *testing.T) {
	// Arrange
	config := &AppConfig{
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

// 標準のテスト関数（Goのテストランナーが認識する）

func TestAppConfig_SetDefaults_Normal(t *testing.T) {
	testInstance := &TestAppConfig{}
	testInstance.TestAppConfig_SetDefaults_Normal(t)
}

func TestAppConfig_SetDefaults_WithExistingValues(t *testing.T) {
	testInstance := &TestAppConfig{}
	testInstance.TestAppConfig_SetDefaults_WithExistingValues(t)
}

func TestAppConfig_SetDefaults_BaseDirFallback(t *testing.T) {
	testInstance := &TestAppConfig{}
	testInstance.TestAppConfig_SetDefaults_BaseDirFallback(t)
}

func TestAppConfig_GetCLIPath_Normal(t *testing.T) {
	testInstance := &TestAppConfig{}
	testInstance.TestAppConfig_GetCLIPath_Normal(t)
}

func TestAppConfig_GetCLIPath_WithDifferentPaths(t *testing.T) {
	testInstance := &TestAppConfig{}
	testInstance.TestAppConfig_GetCLIPath_WithDifferentPaths(t)
}

func TestAppConfig_GetScriptsPath_Normal(t *testing.T) {
	testInstance := &TestAppConfig{}
	testInstance.TestAppConfig_GetScriptsPath_Normal(t)
}

func TestAppConfig_GetScriptsPath_WithDifferentPaths(t *testing.T) {
	testInstance := &TestAppConfig{}
	testInstance.TestAppConfig_GetScriptsPath_WithDifferentPaths(t)
}

func TestAppConfig_Integration(t *testing.T) {
	testInstance := &TestAppConfig{}
	testInstance.TestAppConfig_Integration(t)
}
