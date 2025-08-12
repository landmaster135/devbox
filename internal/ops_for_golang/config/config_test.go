package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// #==============================================================#
// ##          Mock Structures                                   ##
// #==============================================================#

// MockFlagParser はFlagParserのモック実装です
type MockFlagParser struct {
	stringVars   map[string]*string
	boolVars     map[string]*bool
	stringValues map[string]string
	boolValues   map[string]bool
	parseError   error
	args         []string
}

// NewMockFlagParser は新しいMockFlagParserを作成します
func NewMockFlagParser() *MockFlagParser {
	return &MockFlagParser{
		stringVars:   make(map[string]*string),
		boolVars:     make(map[string]*bool),
		stringValues: make(map[string]string),
		boolValues:   make(map[string]bool),
		args:         []string{},
	}
}

// SetStringValue は文字列フラグの値を事前設定します
func (m *MockFlagParser) SetStringValue(name, value string) {
	m.stringValues[name] = value
}

// SetBoolValue はブールフラグの値を事前設定します
func (m *MockFlagParser) SetBoolValue(name string, value bool) {
	m.boolValues[name] = value
}

// SetParseError はParse時に返すエラーを設定します
func (m *MockFlagParser) SetParseError(err error) {
	m.parseError = err
}

// SetArgs は解析後の残り引数を設定します
func (m *MockFlagParser) SetArgs(args []string) {
	m.args = args
}

// StringVar は文字列フラグを定義します
func (m *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
	if presetValue, exists := m.stringValues[name]; exists {
		*p = presetValue // 事前設定値を優先
	} else {
		*p = value // デフォルト値
	}
	m.stringVars[name] = p
}

// BoolVar はブールフラグを定義します
func (m *MockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	if presetValue, exists := m.boolValues[name]; exists {
		*p = presetValue // 事前設定値を優先
	} else {
		*p = value // デフォルト値
	}
	m.boolVars[name] = p
}

// Parse はフラグを解析します
func (m *MockFlagParser) Parse() error {
	return m.parseError
}

// Args は解析後の残りの引数を返します
func (m *MockFlagParser) Args() []string {
	return m.args
}

// #==============================================================#
// ##          NewConfig Tests                                   ##
// #==============================================================#

// TestNewConfig_TestCoverage_Normal はNewConfigのtest-coverage正常系テストです
func TestNewConfig_TestCoverage_Normal(t *testing.T) {
	// Arrange
	const (
		operation   = "test-coverage"
		directory   = "/test/dir"
		grepPattern = "coverage"
	)

	// Act
	config, err := NewConfig(operation, directory, "", "", "", "", grepPattern)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, operation, config.Operation)
	assert.Equal(t, directory, config.Directory)
	assert.Equal(t, grepPattern, config.GrepPattern)
	assert.Empty(t, config.ExecutionFile)
	assert.Empty(t, config.RootDirectory)
	assert.Empty(t, config.Parameters)
	assert.Empty(t, config.CoverageFile)
	assert.False(t, config.Help)
}

// TestNewConfig_TestCoverageProject_Normal はNewConfigのtest-coverage-project正常系テストです
func TestNewConfig_TestCoverageProject_Normal(t *testing.T) {
	// Arrange
	const (
		operation = "test-coverage-project"
		directory = "/test/project"
	)

	// Act
	config, err := NewConfig(operation, directory, "", "", "", "", "")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, operation, config.Operation)
	assert.Equal(t, directory, config.Directory)
	assert.Empty(t, config.GrepPattern)
}

// TestNewConfig_Run_Normal はNewConfigのrun正常系テストです
func TestNewConfig_Run_Normal(t *testing.T) {
	// Arrange
	const (
		operation     = "run"
		executionFile = "/test/main.go"
		rootDirectory = "/test"
		parameters    = "-flag value"
	)

	// Act
	config, err := NewConfig(operation, "", executionFile, rootDirectory, parameters, "", "")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, operation, config.Operation)
	assert.Equal(t, executionFile, config.ExecutionFile)
	assert.Equal(t, rootDirectory, config.RootDirectory)
	assert.Equal(t, parameters, config.Parameters)
	assert.Empty(t, config.Directory)
	assert.Empty(t, config.CoverageFile)
}

// TestNewConfig_CoverageFunc_Normal はNewConfigのcoverage-func正常系テストです
func TestNewConfig_CoverageFunc_Normal(t *testing.T) {
	// Arrange
	const (
		operation    = "coverage-func"
		coverageFile = "/test/coverage.out"
		grepPattern  = "total"
	)

	// Act
	config, err := NewConfig(operation, "", "", "", "", coverageFile, grepPattern)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, operation, config.Operation)
	assert.Equal(t, coverageFile, config.CoverageFile)
	assert.Equal(t, grepPattern, config.GrepPattern)
	assert.Empty(t, config.Directory)
	assert.Empty(t, config.ExecutionFile)
	assert.Empty(t, config.RootDirectory)
	assert.Empty(t, config.Parameters)
}

// TestNewConfig_EmptyOperation はNewConfigの空操作タイプテストです
func TestNewConfig_EmptyOperation(t *testing.T) {
	// Act
	config, err := NewConfig("", "/test/dir", "", "", "", "", "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "操作タイプが指定されていません")
}

// TestNewConfig_InvalidOperation はNewConfigの無効操作タイプテストです
func TestNewConfig_InvalidOperation(t *testing.T) {
	// Arrange
	const invalidOperation = "invalid-operation"

	// Act
	config, err := NewConfig(invalidOperation, "/test/dir", "", "", "", "", "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "無効な操作タイプです: "+invalidOperation)
}

// TestNewConfig_TestCoverage_MissingDirectory はNewConfigのtest-coverageディレクトリ不足テストです
func TestNewConfig_TestCoverage_MissingDirectory(t *testing.T) {
	// Act
	config, err := NewConfig("test-coverage", "", "", "", "", "", "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "test-coverage操作にはディレクトリパスが必要です")
}

// TestNewConfig_TestCoverageProject_MissingDirectory はNewConfigのtest-coverage-projectディレクトリ不足テストです
func TestNewConfig_TestCoverageProject_MissingDirectory(t *testing.T) {
	// Act
	config, err := NewConfig("test-coverage-project", "", "", "", "", "", "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "test-coverage-project操作にはディレクトリパスが必要です")
}

// TestNewConfig_Run_MissingExecutionFile はNewConfigのrun実行ファイル不足テストです
func TestNewConfig_Run_MissingExecutionFile(t *testing.T) {
	// Act
	config, err := NewConfig("run", "", "", "/test", "", "", "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "run操作には実行ファイルが必要です")
}

// TestNewConfig_Run_MissingRootDirectory はNewConfigのrunルートディレクトリ不足テストです
func TestNewConfig_Run_MissingRootDirectory(t *testing.T) {
	// Act
	config, err := NewConfig("run", "", "/test/main.go", "", "", "", "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "run操作にはルートディレクトリが必要です")
}

// TestNewConfig_CoverageFunc_MissingCoverageFile はNewConfigのcoverage-funcカバレッジファイル不足テストです
func TestNewConfig_CoverageFunc_MissingCoverageFile(t *testing.T) {
	// Act
	config, err := NewConfig("coverage-func", "", "", "", "", "", "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "coverage-func操作にはカバレッジファイルが必要です")
}

// #==============================================================#
// ##          ParseFlags Tests                                  ##
// #==============================================================#

// TestParseFlags_TestCoverage_Normal はParseFlagsのtest-coverage正常系テストです
func TestParseFlags_TestCoverage_Normal(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockFlagParser)
		expectedConfig *Config
		expectError    bool
		errorMessage   string
	}{
		{
			name: "TestCoverage_LongFlags_Normal",
			setupMock: func(mock *MockFlagParser) {
				mock.SetStringValue("ops", "test-coverage")
				mock.SetStringValue("directory", "/test/dir")
				mock.SetStringValue("grep_pattern", "coverage")
			},
			expectedConfig: &Config{
				Operation:   "test-coverage",
				Directory:   "/test/dir",
				GrepPattern: "coverage",
			},
			expectError: false,
		},
		{
			name: "TestCoverage_ShortFlags_Normal",
			setupMock: func(mock *MockFlagParser) {
				mock.SetStringValue("o", "test-coverage")
				mock.SetStringValue("d", "/test/dir")
				mock.SetStringValue("g", "PASS")
			},
			expectedConfig: &Config{
				Operation:   "test-coverage",
				Directory:   "/test/dir",
				GrepPattern: "PASS",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockParser := NewMockFlagParser()
			tt.setupMock(mockParser)

			// Act
			config, err := ParseFlagsWithParser(mockParser)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, config)
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, config)
				assert.Equal(t, tt.expectedConfig.Operation, config.Operation)
				assert.Equal(t, tt.expectedConfig.Directory, config.Directory)
				assert.Equal(t, tt.expectedConfig.GrepPattern, config.GrepPattern)
			}
		})
	}
}

// TestParseFlags_TestCoverageProject_Normal はParseFlagsのtest-coverage-project正常系テストです
func TestParseFlags_TestCoverageProject_Normal(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringValue("ops", "test-coverage-project")
	mockParser.SetStringValue("directory", "/test/project")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "test-coverage-project", config.Operation)
	assert.Equal(t, "/test/project", config.Directory)
}

// TestParseFlags_Run_Normal はParseFlagsのrun正常系テストです
func TestParseFlags_Run_Normal(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringValue("ops", "run")
	mockParser.SetStringValue("execution_file", "/test/main.go")
	mockParser.SetStringValue("root_directory", "/test")
	mockParser.SetStringValue("parameters", "-flag value")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "run", config.Operation)
	assert.Equal(t, "/test/main.go", config.ExecutionFile)
	assert.Equal(t, "/test", config.RootDirectory)
	assert.Equal(t, "-flag value", config.Parameters)
}

// TestParseFlags_CoverageFunc_Normal はParseFlagsのcoverage-func正常系テストです
func TestParseFlags_CoverageFunc_Normal(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringValue("ops", "coverage-func")
	mockParser.SetStringValue("coverage_file", "/test/coverage.out")
	mockParser.SetStringValue("grep_pattern", "total")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "coverage-func", config.Operation)
	assert.Equal(t, "/test/coverage.out", config.CoverageFile)
	assert.Equal(t, "total", config.GrepPattern)
}

// TestParseFlags_Help_Normal はParseFlagsのヘルプ正常系テストです
func TestParseFlags_Help_Normal(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*MockFlagParser)
	}{
		{
			name: "Help_LongFlag_Normal",
			setupMock: func(mock *MockFlagParser) {
				mock.SetBoolValue("help", true)
			},
		},
		{
			name: "Help_ShortFlag_Normal",
			setupMock: func(mock *MockFlagParser) {
				mock.SetBoolValue("h", true)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockParser := NewMockFlagParser()
			tt.setupMock(mockParser)

			// Act
			config, err := ParseFlagsWithParser(mockParser)

			// Assert
			assert.NoError(t, err)
			assert.NotNil(t, config)
			assert.True(t, config.Help)
		})
	}
}

// TestParseFlags_ParseError はParseFlagsの解析エラーテストです
func TestParseFlags_ParseError(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	parseError := fmt.Errorf("invalid flag")
	mockParser.SetParseError(parseError)

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "フラグの解析に失敗しました")
}

// TestParseFlags_InvalidOperation はParseFlagsの無効操作タイプテストです
func TestParseFlags_InvalidOperation(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringValue("ops", "invalid-operation")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "無効な操作タイプです")
}

// TestParseFlags_MissingRequiredParams はParseFlagsの必須パラメータ不足テストです
func TestParseFlags_MissingRequiredParams(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func(*MockFlagParser)
		errorMessage string
	}{
		{
			name: "TestCoverage_MissingDirectory",
			setupMock: func(mock *MockFlagParser) {
				mock.SetStringValue("ops", "test-coverage")
			},
			errorMessage: "test-coverage操作にはディレクトリパスが必要です",
		},
		{
			name: "Run_MissingExecutionFile",
			setupMock: func(mock *MockFlagParser) {
				mock.SetStringValue("ops", "run")
				mock.SetStringValue("root_directory", "/test")
			},
			errorMessage: "run操作には実行ファイルが必要です",
		},
		{
			name: "CoverageFunc_MissingCoverageFile",
			setupMock: func(mock *MockFlagParser) {
				mock.SetStringValue("ops", "coverage-func")
			},
			errorMessage: "coverage-func操作にはカバレッジファイルが必要です",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockParser := NewMockFlagParser()
			tt.setupMock(mockParser)

			// Act
			config, err := ParseFlagsWithParser(mockParser)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, config)
			assert.Contains(t, err.Error(), tt.errorMessage)
		})
	}
}

// #==============================================================#
// ##          StandardFlagParser Tests                         ##
// #==============================================================#

// TestNewStandardFlagParser_Normal はNewStandardFlagParserの正常系テストです
func TestNewStandardFlagParser_Normal(t *testing.T) {
	// Act
	parser := NewStandardFlagParser()

	// Assert
	assert.NotNil(t, parser)
	assert.NotNil(t, parser.flagSet)
}

// TestStandardFlagParser_StringVar_Normal はStandardFlagParserのStringVar正常系テストです
func TestStandardFlagParser_StringVar_Normal(t *testing.T) {
	// Arrange
	parser := NewStandardFlagParser()
	var testString string

	// Act
	parser.StringVar(&testString, "test", "default", "test usage")

	// Assert
	assert.Equal(t, "default", testString)
}

// TestStandardFlagParser_BoolVar_Normal はStandardFlagParserのBoolVar正常系テストです
func TestStandardFlagParser_BoolVar_Normal(t *testing.T) {
	// Arrange
	parser := NewStandardFlagParser()
	var testBool bool

	// Act
	parser.BoolVar(&testBool, "test", true, "test usage")

	// Assert
	assert.True(t, testBool)
}

// TestStandardFlagParser_Args_Normal はStandardFlagParserのArgs正常系テストです
func TestStandardFlagParser_Args_Normal(t *testing.T) {
	// Arrange
	parser := NewStandardFlagParser()

	// Act
	args := parser.Args()

	// Assert
	assert.Nil(t, args)
}

// #==============================================================#
// ##          PrintUsage Function Tests                         ##
// #==============================================================#

// TestPrintUsage_Normal はPrintUsage関数の正常系テストです
func TestPrintUsage_Normal(t *testing.T) {
	// Note: PrintUsage関数はos.Stderrに出力するため、
	// 単体テストでは出力内容の検証が困難です。

	// 最低限の動作確認として、関数が存在し、パニックしないことを確認
	assert.NotPanics(t, func() {
		PrintUsage()
	})
}
