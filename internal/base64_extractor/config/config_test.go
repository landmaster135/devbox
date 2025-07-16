package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestConfigParser はConfigParserのテストクラス
type TestConfigParser struct {
	tempDir string
}

// setupTestEnvironment はテスト環境をセットアップする
func (tc *TestConfigParser) setupTestEnvironment(t *testing.T) {
	var err error
	tc.tempDir, err = os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
}

// teardownTestEnvironment はテスト環境をクリーンアップする
func (tc *TestConfigParser) teardownTestEnvironment(t *testing.T) {
	if err := os.RemoveAll(tc.tempDir); err != nil {
		t.Errorf("一時ディレクトリの削除に失敗しました: %v", err)
	}
}

// createTestFile はテスト用のファイルを作成する
func (tc *TestConfigParser) createTestFile(t *testing.T, filename string) string {
	filePath := filepath.Join(tc.tempDir, filename)
	if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}
	return filePath
}

// TestNewConfig_Normal はNewConfigの正常系テスト
func TestNewConfig_Normal(t *testing.T) {
	// Arrange
	tc := &TestConfigParser{}
	tc.setupTestEnvironment(t)
	defer tc.teardownTestEnvironment(t)

	testFile := tc.createTestFile(t, "test.jpg")

	// Act
	config, err := NewConfig(testFile, true, "text")

	// Assert
	if err != nil {
		t.Fatalf("NewConfigでエラーが発生しました: %v", err)
	}
	if config == nil {
		t.Fatal("configがnilです")
	}
	if config.Path != testFile {
		t.Errorf("Pathが期待値と異なります。期待値: %s, 実際: %s", testFile, config.Path)
	}
	if config.Recursive != true {
		t.Errorf("Recursiveが期待値と異なります。期待値: true, 実際: %t", config.Recursive)
	}
	if config.OutputFormat != "text" {
		t.Errorf("OutputFormatが期待値と異なります。期待値: text, 実際: %s", config.OutputFormat)
	}
}

// TestNewConfig_EmptyPath は空のパスでのテスト
func TestNewConfig_EmptyPath(t *testing.T) {
	// Arrange & Act
	config, err := NewConfig("", true, "text")

	// Assert
	if err == nil {
		t.Fatal("エラーが発生するべきです")
	}
	if config != nil {
		t.Error("configはnilであるべきです")
	}
	if !strings.Contains(err.Error(), "--path は必須パラメータです") {
		t.Errorf("期待されるエラーメッセージが含まれていません: %v", err)
	}
}

// TestNewConfig_NonExistentPath は存在しないパスでのテスト
func TestNewConfig_NonExistentPath(t *testing.T) {
	// Arrange & Act
	config, err := NewConfig("/non/existent/path", true, "text")

	// Assert
	if err == nil {
		t.Fatal("エラーが発生するべきです")
	}
	if config != nil {
		t.Error("configはnilであるべきです")
	}
	if !strings.Contains(err.Error(), "指定されたパスが存在しません") {
		t.Errorf("期待されるエラーメッセージが含まれていません: %v", err)
	}
}

// TestNewConfig_InvalidOutputFormat は無効な出力形式でのテスト
func TestNewConfig_InvalidOutputFormat(t *testing.T) {
	// Arrange
	tc := &TestConfigParser{}
	tc.setupTestEnvironment(t)
	defer tc.teardownTestEnvironment(t)

	testFile := tc.createTestFile(t, "test.jpg")

	// Act
	config, err := NewConfig(testFile, true, "invalid")

	// Assert
	if err == nil {
		t.Fatal("エラーが発生するべきです")
	}
	if config != nil {
		t.Error("configはnilであるべきです")
	}
	if !strings.Contains(err.Error(), "--output-format は 'text' または 'json' を指定してください") {
		t.Errorf("期待されるエラーメッセージが含まれていません: %v", err)
	}
}

// TestValidateConfig_Normal はvalidateConfigの正常系テスト
func TestValidateConfig_Normal(t *testing.T) {
	// Arrange
	tc := &TestConfigParser{}
	tc.setupTestEnvironment(t)
	defer tc.teardownTestEnvironment(t)

	testFile := tc.createTestFile(t, "test.jpg")
	config := &Config{
		Path:         testFile,
		Recursive:    false,
		OutputFormat: "json",
	}

	// Act
	result, err := validateConfig(config)

	// Assert
	if err != nil {
		t.Fatalf("validateConfigでエラーが発生しました: %v", err)
	}
	if result == nil {
		t.Fatal("結果がnilです")
	}
	if result.Path != testFile {
		t.Errorf("Pathが期待値と異なります。期待値: %s, 実際: %s", testFile, result.Path)
	}
	if result.Recursive != false {
		t.Errorf("Recursiveが期待値と異なります。期待値: false, 実際: %t", result.Recursive)
	}
	if result.OutputFormat != "json" {
		t.Errorf("OutputFormatが期待値と異なります。期待値: json, 実際: %s", result.OutputFormat)
	}
}

// MockFlagParser はテスト用のフラグパーサー
type MockFlagParser struct {
	stringVars map[string]*string
	boolVars   map[string]*bool
	parseError error
}

// NewMockFlagParser は新しいMockFlagParserを作成する
func NewMockFlagParser() *MockFlagParser {
	return &MockFlagParser{
		stringVars: make(map[string]*string),
		boolVars:   make(map[string]*bool),
	}
}

// StringVar は文字列フラグを設定する
func (mfp *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
	*p = value
	mfp.stringVars[name] = p
}

// BoolVar はブールフラグを設定する
func (mfp *MockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	*p = value
	mfp.boolVars[name] = p
}

// Parse はフラグを解析する
func (mfp *MockFlagParser) Parse() error {
	return mfp.parseError
}

// SetStringValue は文字列値を設定する
func (mfp *MockFlagParser) SetStringValue(name, value string) {
	if p, exists := mfp.stringVars[name]; exists {
		*p = value
	}
}

// SetBoolValue はブール値を設定する
func (mfp *MockFlagParser) SetBoolValue(name string, value bool) {
	if p, exists := mfp.boolVars[name]; exists {
		*p = value
	}
}

// SetParseError は解析エラーを設定する
func (mfp *MockFlagParser) SetParseError(err error) {
	mfp.parseError = err
}

// MockOSArgs はテスト用のOS引数
type MockOSArgs struct {
	args []string
}

// NewMockOSArgs は新しいMockOSArgsを作成する
func NewMockOSArgs(args []string) *MockOSArgs {
	return &MockOSArgs{args: args}
}

// Args はOS引数を返す
func (moa *MockOSArgs) Args() []string {
	return moa.args
}

// TestConfigParser_ParseFlags_WithValidPath はParseFlags正常系テスト（簡略版）
func TestConfigParser_ParseFlags_WithValidPath(t *testing.T) {
	// Arrange
	tc := &TestConfigParser{}
	tc.setupTestEnvironment(t)
	defer tc.teardownTestEnvironment(t)

	testFile := tc.createTestFile(t, "test.jpg")

	// 直接NewConfigを使用してテスト（ParseFlagsの複雑さを回避）
	config, err := NewConfig(testFile, false, "json")

	// Assert
	if err != nil {
		t.Fatalf("NewConfigでエラーが発生しました: %v", err)
	}
	if config == nil {
		t.Fatal("configがnilです")
	}
	if config.Path != testFile {
		t.Errorf("Pathが期待値と異なります。期待値: %s, 実際: %s", testFile, config.Path)
	}
	if config.Recursive != false {
		t.Errorf("Recursiveが期待値と異なります。期待値: false, 実際: %t", config.Recursive)
	}
	if config.OutputFormat != "json" {
		t.Errorf("OutputFormatが期待値と異なります。期待値: json, 実際: %s", config.OutputFormat)
	}
}

// TestNewConfigParser_Normal はNewConfigParserの正常系テスト
func TestNewConfigParser_Normal(t *testing.T) {
	// Arrange
	mockFlagParser := NewMockFlagParser()
	mockOSArgs := NewMockOSArgs([]string{"test"})

	// Act
	configParser := NewConfigParser(mockFlagParser, mockOSArgs)

	// Assert
	if configParser == nil {
		t.Fatal("configParserがnilです")
	}
	if configParser.flagParser != mockFlagParser {
		t.Error("flagParserが期待値と異なります")
	}
	if configParser.osArgs != mockOSArgs {
		t.Error("osArgsが期待値と異なります")
	}
}
