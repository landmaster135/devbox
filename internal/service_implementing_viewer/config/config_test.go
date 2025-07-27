package config

import (
	"testing"
)

// TestConfig はConfigのテストクラス
type TestConfig struct{}

// TestNewConfig_Normal はNewConfigの正常系テスト
func TestNewConfig_Normal(t *testing.T) {
	// Arrange
	rootDir := "/test/root"
	targetDirs := "cli,mcp,powershell"

	// Act
	config, err := NewConfig(rootDir, targetDirs)

	// Assert
	if err != nil {
		t.Fatalf("エラーが発生しました: %v", err)
	}
	if config == nil {
		t.Fatal("configがnilです")
	}
	if config.RootDir != rootDir {
		t.Errorf("RootDirが期待値と異なります。期待値: %s, 実際: %s", rootDir, config.RootDir)
	}
	expectedTargetDirs := []string{"cli", "mcp", "powershell"}
	if len(config.TargetDirs) != len(expectedTargetDirs) {
		t.Errorf("TargetDirsの長さが期待値と異なります。期待値: %d, 実際: %d", len(expectedTargetDirs), len(config.TargetDirs))
	}
	for i, expected := range expectedTargetDirs {
		if config.TargetDirs[i] != expected {
			t.Errorf("TargetDirs[%d]が期待値と異なります。期待値: %s, 実際: %s", i, expected, config.TargetDirs[i])
		}
	}
}

// TestNewConfig_EmptyRootDir は空のRootDirのテスト
func TestNewConfig_EmptyRootDir(t *testing.T) {
	// Arrange
	rootDir := ""
	targetDirs := "cli,mcp"

	// Act
	config, err := NewConfig(rootDir, targetDirs)

	// Assert
	if err == nil {
		t.Fatal("エラーが発生しませんでした")
	}
	if config != nil {
		t.Error("configがnilではありません")
	}
	if err.Error() != "設定の初期化に失敗しました: --root-dir は必須パラメータです" {
		t.Errorf("エラーメッセージが期待値と異なります: %v", err)
	}
}

// TestNewConfig_EmptyTargetDirs は空のTargetDirsのテスト
func TestNewConfig_EmptyTargetDirs(t *testing.T) {
	// Arrange
	rootDir := "/test/root"
	targetDirs := ""

	// Act
	config, err := NewConfig(rootDir, targetDirs)

	// Assert
	if err == nil {
		t.Fatal("エラーが発生しませんでした")
	}
	if config != nil {
		t.Error("configがnilではありません")
	}
	if err.Error() != "設定の初期化に失敗しました: --target-dirs は必須パラメータです" {
		t.Errorf("エラーメッセージが期待値と異なります: %v", err)
	}
}

// TestParseTargetDirs_Normal はparseTargetDirsの正常系テスト
func TestParseTargetDirs_Normal(t *testing.T) {
	// Arrange
	testCases := []struct {
		input    string
		expected []string
	}{
		{"cli,mcp", []string{"cli", "mcp"}},
		{"cli, mcp, powershell", []string{"cli", "mcp", "powershell"}},
		{"  cli  ,  mcp  ", []string{"cli", "mcp"}},
		{"single", []string{"single"}},
		{"", []string{}},
		{"cli,,mcp", []string{"cli", "mcp"}},
	}

	for _, testCase := range testCases {
		// Act
		result := parseTargetDirs(testCase.input)

		// Assert
		if len(result) != len(testCase.expected) {
			t.Errorf("結果の長さが期待値と異なります。入力: %s, 期待値: %d, 実際: %d", testCase.input, len(testCase.expected), len(result))
			continue
		}
		for i, expected := range testCase.expected {
			if result[i] != expected {
				t.Errorf("結果[%d]が期待値と異なります。入力: %s, 期待値: %s, 実際: %s", i, testCase.input, expected, result[i])
			}
		}
	}
}

// TestValidateConfig_Normal はvalidateConfigの正常系テスト
func TestValidateConfig_Normal(t *testing.T) {
	// Arrange
	config := &Config{
		RootDir:    "/test/root",
		TargetDirs: []string{"cli", "mcp"},
	}

	// Act
	result, err := validateConfig(config)

	// Assert
	if err != nil {
		t.Fatalf("エラーが発生しました: %v", err)
	}
	if result == nil {
		t.Fatal("結果がnilです")
	}
	if result.RootDir != config.RootDir {
		t.Errorf("RootDirが期待値と異なります。期待値: %s, 実際: %s", config.RootDir, result.RootDir)
	}
}

// TestValidateConfig_EmptyRootDir は空のRootDirの検証テスト
func TestValidateConfig_EmptyRootDir(t *testing.T) {
	// Arrange
	config := &Config{
		RootDir:    "",
		TargetDirs: []string{"cli", "mcp"},
	}

	// Act
	result, err := validateConfig(config)

	// Assert
	if err == nil {
		t.Fatal("エラーが発生しませんでした")
	}
	if result != nil {
		t.Error("結果がnilではありません")
	}
	if err.Error() != "--root-dir は必須パラメータです" {
		t.Errorf("エラーメッセージが期待値と異なります: %v", err)
	}
}

// TestValidateConfig_EmptyTargetDirs は空のTargetDirsの検証テスト
func TestValidateConfig_EmptyTargetDirs(t *testing.T) {
	// Arrange
	config := &Config{
		RootDir:    "/test/root",
		TargetDirs: []string{},
	}

	// Act
	result, err := validateConfig(config)

	// Assert
	if err == nil {
		t.Fatal("エラーが発生しませんでした")
	}
	if result != nil {
		t.Error("結果がnilではありません")
	}
	if err.Error() != "--target-dirs は必須パラメータです" {
		t.Errorf("エラーメッセージが期待値と異なります: %v", err)
	}
}

// MockFlagParser はテスト用のFlagParser実装
type MockFlagParser struct {
	rootDir    string
	targetDirs string
	parseError error
}

// StringVar はモックの文字列フラグ定義
func (m *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
	switch name {
	case "root-dir":
		*p = m.rootDir
	case "target-dirs":
		*p = m.targetDirs
	}
}

// BoolVar はモックのブールフラグ定義
func (m *MockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	// テストでは使用しない
}

// Parse はモックのフラグ解析
func (m *MockFlagParser) Parse() error {
	return m.parseError
}

// MockOSArgs はテスト用のOSArgs実装
type MockOSArgs struct {
	args []string
}

// Args はモックのOS引数取得
func (m *MockOSArgs) Args() []string {
	return m.args
}

// TestConfigParser はConfigParserのテストクラス
type TestConfigParser struct{}

// TestConfigParser_ParseFlags_Normal はConfigParserのParseFlags正常系テスト
func TestConfigParser_ParseFlags_Normal(t *testing.T) {
	// Arrange
	mockFlagParser := &MockFlagParser{
		rootDir:    "/test/root",
		targetDirs: "cli,mcp",
		parseError: nil,
	}
	mockOSArgs := &MockOSArgs{
		args: []string{"program", "-root-dir=/test/root", "-target-dirs=cli,mcp"},
	}
	configParser := NewConfigParser(mockFlagParser, mockOSArgs)

	// Act
	config, err := configParser.ParseFlags()

	// Assert
	if err != nil {
		t.Fatalf("エラーが発生しました: %v", err)
	}
	if config == nil {
		t.Fatal("configがnilです")
	}
	if config.RootDir != "/test/root" {
		t.Errorf("RootDirが期待値と異なります。期待値: /test/root, 実際: %s", config.RootDir)
	}
	expectedTargetDirs := []string{"cli", "mcp"}
	if len(config.TargetDirs) != len(expectedTargetDirs) {
		t.Errorf("TargetDirsの長さが期待値と異なります。期待値: %d, 実際: %d", len(expectedTargetDirs), len(config.TargetDirs))
	}
}

// TestConfigParser_ParseFlags_ParseError はConfigParserのParseFlags解析エラーテスト
func TestConfigParser_ParseFlags_ParseError(t *testing.T) {
	// Arrange
	mockFlagParser := &MockFlagParser{
		rootDir:    "/test/root",
		targetDirs: "cli,mcp",
		parseError: &MockError{"parse error"},
	}
	mockOSArgs := &MockOSArgs{
		args: []string{"program"},
	}
	configParser := NewConfigParser(mockFlagParser, mockOSArgs)

	// Act
	config, err := configParser.ParseFlags()

	// Assert
	if err == nil {
		t.Fatal("エラーが発生しませんでした")
	}
	if config != nil {
		t.Error("configがnilではありません")
	}
}

// MockError はテスト用のエラー実装
type MockError struct {
	message string
}

// Error はエラーメッセージを返す
func (e *MockError) Error() string {
	return e.message
}
