package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ConfigTestSuite struct {
	tempDir string
}

func (suite *ConfigTestSuite) SetupTest(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "config_test")
	require.NoError(t, err)
	suite.tempDir = tempDir
}

func (suite *ConfigTestSuite) TeardownTest(t *testing.T) {
	// テスト用の一時ディレクトリを削除
	if suite.tempDir != "" {
		os.RemoveAll(suite.tempDir)
	}
}

func TestNewConfig_Normal(t *testing.T) {
	suite := &ConfigTestSuite{}
	suite.SetupTest(t)
	defer suite.TeardownTest(t)

	// Arrange
	path := suite.tempDir
	recursive := true
	language := "jpn,eng"
	outputFormat := "json"
	outputDir := filepath.Join(suite.tempDir, "output")

	// Act
	config, err := NewConfig(path, recursive, language, outputFormat, outputDir)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, path, config.Path)
	assert.Equal(t, recursive, config.Recursive)
	assert.Equal(t, language, config.Language)
	assert.Equal(t, outputFormat, config.OutputFormat)
	assert.Equal(t, outputDir, config.OutputDir)
}

func TestNewConfig_InvalidPath_Normal(t *testing.T) {
	// Arrange
	path := "/non/existent/path"
	recursive := false
	language := "eng"
	outputFormat := "text"
	outputDir := ""

	// Act
	config, err := NewConfig(path, recursive, language, outputFormat, outputDir)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "指定されたパスが存在しません")
}

func TestNewConfig_EmptyPath_Normal(t *testing.T) {
	// Arrange
	path := ""
	recursive := false
	language := "eng"
	outputFormat := "text"
	outputDir := ""

	// Act
	config, err := NewConfig(path, recursive, language, outputFormat, outputDir)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "--path は必須パラメータです")
}

func TestNewConfig_InvalidOutputFormat_Normal(t *testing.T) {
	suite := &ConfigTestSuite{}
	suite.SetupTest(t)
	defer suite.TeardownTest(t)

	// Arrange
	path := suite.tempDir
	recursive := false
	language := "eng"
	outputFormat := "invalid"
	outputDir := ""

	// Act
	config, err := NewConfig(path, recursive, language, outputFormat, outputDir)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "--output-format は 'text' または 'json' を指定してください")
}

func TestValidateLanguages_Normal(t *testing.T) {
	testCases := []struct {
		name      string
		languages string
		expectErr bool
	}{
		{"Single language", "eng", false},
		{"Multiple languages", "jpn,eng,fra", false},
		{"With spaces", "jpn, eng, fra", false},
		{"Empty string", "", true},
		{"Empty language code", "jpn,,eng", true},
		{"Invalid long code", "jpn,english,fra", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			err := validateLanguages(tc.languages)

			// Assert
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfig_GetTesseractLanguages_Normal(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"Single language", "eng", "eng"},
		{"Multiple languages", "jpn,eng", "jpn+eng"},
		{"With spaces", "jpn, eng, fra", "jpn+eng+fra"},
		{"Three languages", "jpn,eng,fra", "jpn+eng+fra"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			config := &Config{Language: tc.input}

			// Act
			result := config.GetTesseractLanguages()

			// Assert
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestConfig_OutputDirCreation_Normal(t *testing.T) {
	suite := &ConfigTestSuite{}
	suite.SetupTest(t)
	defer suite.TeardownTest(t)

	// Arrange
	path := suite.tempDir
	recursive := false
	language := "eng"
	outputFormat := "text"
	outputDir := filepath.Join(suite.tempDir, "new", "output", "dir")

	// Act
	config, err := NewConfig(path, recursive, language, outputFormat, outputDir)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, config)

	// 出力ディレクトリが作成されたことを確認
	assert.DirExists(t, outputDir)
}
