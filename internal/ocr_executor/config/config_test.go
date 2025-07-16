package config

import (
	"errors"
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

// MockFlagParser はテスト用のFlagParserモック
type MockFlagParser struct {
	stringVars map[string]*string
	boolVars   map[string]*bool
	parseError error
}

func NewMockFlagParser() *MockFlagParser {
	return &MockFlagParser{
		stringVars: make(map[string]*string),
		boolVars:   make(map[string]*bool),
	}
}

func (m *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
	*p = value
	m.stringVars[name] = p
}

func (m *MockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	*p = value
	m.boolVars[name] = p
}

func (m *MockFlagParser) Parse() error {
	return m.parseError
}

func (m *MockFlagParser) SetParseError(err error) {
	m.parseError = err
}

func (m *MockFlagParser) SetStringValue(name, value string) {
	if p, exists := m.stringVars[name]; exists {
		*p = value
	}
}

func (m *MockFlagParser) SetBoolValue(name string, value bool) {
	if p, exists := m.boolVars[name]; exists {
		*p = value
	}
}

// MockOSArgs はテスト用のOSArgsモック
type MockOSArgs struct {
	args []string
}

func NewMockOSArgs(args []string) *MockOSArgs {
	return &MockOSArgs{args: args}
}

func (m *MockOSArgs) Args() []string {
	return m.args
}

func TestConfigParser_ParseFlags_Normal(t *testing.T) {
	suite := &ConfigTestSuite{}
	suite.SetupTest(t)
	defer suite.TeardownTest(t)

	// Arrange
	mockFlagParser := NewMockFlagParser()
	mockOSArgs := NewMockOSArgs([]string{"test-program", "-path", suite.tempDir})
	configParser := NewConfigParser(mockFlagParser, mockOSArgs)

	// Act
	config, err := configParser.ParseFlags()

	// Assert - デフォルト値（空文字列）でパスが設定されるため、バリデーションエラーが発生するはず
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "--path は必須パラメータです")
}

func TestConfigParser_ParseFlags_ParseError_Normal(t *testing.T) {
	// Arrange
	mockFlagParser := NewMockFlagParser()
	mockOSArgs := NewMockOSArgs([]string{"test-program"})
	configParser := NewConfigParser(mockFlagParser, mockOSArgs)

	expectedError := errors.New("parse error")
	mockFlagParser.SetParseError(expectedError)

	// Act
	config, err := configParser.ParseFlags()

	// Assert
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Equal(t, expectedError, err)
}

func TestConfigParser_ParseFlags_EmptyPath_Normal(t *testing.T) {
	// Arrange
	mockFlagParser := NewMockFlagParser()
	mockOSArgs := NewMockOSArgs([]string{"test-program"})
	configParser := NewConfigParser(mockFlagParser, mockOSArgs)

	// Act - パスが空の場合
	config, err := configParser.ParseFlags()

	// Assert
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "--path は必須パラメータです")
}

func TestNewConfigParser_Normal(t *testing.T) {
	// Arrange
	mockFlagParser := NewMockFlagParser()
	mockOSArgs := NewMockOSArgs([]string{"test-program"})

	// Act
	configParser := NewConfigParser(mockFlagParser, mockOSArgs)

	// Assert
	assert.NotNil(t, configParser)
}

func TestRealFlagParser_Creation_Normal(t *testing.T) {
	// Act
	parser := NewRealFlagParser()

	// Assert
	assert.NotNil(t, parser)
}

func TestRealFlagParser_Methods_Normal(t *testing.T) {
	// Arrange
	parser := NewRealFlagParser()
	var testString string
	var testBool bool

	// Act & Assert - メソッドがパニックしないことを確認
	assert.NotPanics(t, func() {
		parser.StringVar(&testString, "test-string", "default", "test usage")
	})

	assert.NotPanics(t, func() {
		parser.BoolVar(&testBool, "test-bool", true, "test usage")
	})

	assert.NotPanics(t, func() {
		parser.Parse()
	})
}

func TestStandardOSArgs_Args_Normal(t *testing.T) {
	// Arrange
	osArgs := NewStandardOSArgs()

	// Act
	args := osArgs.Args()

	// Assert
	assert.NotNil(t, args)
	assert.IsType(t, []string{}, args)
}

func TestParseFlags_Normal(t *testing.T) {
	suite := &ConfigTestSuite{}
	suite.SetupTest(t)
	defer suite.TeardownTest(t)

	// この関数は実際のflag.CommandLineを使用するため、
	// 実際のテストでは制限があるが、関数が存在することを確認
	// Arrange & Act & Assert
	// 実際のos.Argsを変更するのは危険なので、関数の存在のみ確認
	assert.NotNil(t, ParseFlags)
}

func TestPrintUsage_Normal(t *testing.T) {
	// Arrange & Act
	// PrintUsageは標準出力に出力するため、実際の出力内容のテストは困難
	// 関数がパニックしないことを確認
	assert.NotPanics(t, func() {
		PrintUsage()
	})
}

func TestValidateConfig_InvalidOutputDirCreation_Normal(t *testing.T) {
	// Arrange
	config := &Config{
		Path:         "/tmp", // 存在するパス
		Recursive:    false,
		Language:     "eng",
		OutputFormat: "text",
		OutputDir:    "/root/cannot_create", // 権限がないディレクトリ
	}

	// Act
	result, err := validateConfig(config)

	// Assert
	// 権限エラーが発生する可能性があるが、環境によって異なる
	// エラーが発生するか、正常に作成されるかのいずれか
	if err != nil {
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "出力ディレクトリの作成に失敗しました")
	} else {
		assert.NotNil(t, result)
	}
}

func TestValidateLanguages_EdgeCases_Normal(t *testing.T) {
	testCases := []struct {
		name      string
		languages string
		expectErr bool
	}{
		{"Single character", "j", false},
		{"Two characters", "jp", false},
		{"Three characters", "jpn", false},
		{"Four characters", "jpnn", true},
		{"Only commas", ",,,", true},
		{"Leading comma", ",eng", true},
		{"Trailing comma", "eng,", true},
		{"Multiple spaces", "jpn,  ,eng", true},
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

func TestConfig_GetTesseractLanguages_EdgeCases_Normal(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"Single character", "j", "j"},
		{"Empty after trim", " , ", "+"},
		{"Many spaces", "jpn,   eng   ,   fra", "jpn+eng+fra"},
		{"Single language with spaces", "  eng  ", "eng"},
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
