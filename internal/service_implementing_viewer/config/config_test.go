package config

import (
	"bytes"
	"flag"
	"io"
	"os"
	"strings"
	"testing"
)

// TestConfig はConfigのテストクラス
type TestConfig struct{}

// TestNewConfig_Normal はNewConfigの正常系テスト
func TestNewConfig_Normal(t *testing.T) {
	// Arrange
	rootDir := "/test/root"
	targetDirs := "cli,mcp,powershell"
	operation := "output"

	// Act
	config, err := NewConfig(rootDir, targetDirs, operation)

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
	if config.Operation != operation {
		t.Errorf("Operationが期待値と異なります。期待値: %s, 実際: %s", operation, config.Operation)
	}
}

// TestNewConfig_EmptyRootDir は空のRootDirのテスト
func TestNewConfig_EmptyRootDir(t *testing.T) {
	// Arrange
	rootDir := ""
	targetDirs := "cli,mcp"

	// Act
	config, err := NewConfig(rootDir, targetDirs, "output")

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
	config, err := NewConfig(rootDir, targetDirs, "output")

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
		Operation:  "output",
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
		Operation:  "output",
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
		Operation:  "output",
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

// TestValidateConfig_EmptyOperation は空Operationの検証テスト
func TestValidateConfig_EmptyOperation(t *testing.T) {
	// Arrange
	config := &Config{
		RootDir:    "/test/root",
		TargetDirs: []string{"cli", "mcp"},
		Operation:  "",
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
	if err.Error() != "--operation は必須です" {
		t.Errorf("エラーメッセージが期待値と異なります: %v", err)
	}
}

// TestValidateConfig_InvalidOperation は不正Operationの検証テスト
func TestValidateConfig_InvalidOperation(t *testing.T) {
	// Arrange
	config := &Config{
		RootDir:    "/test/root",
		TargetDirs: []string{"cli", "mcp"},
		Operation:  "invalid",
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
	if !strings.Contains(err.Error(), "--operation は次のいずれかを指定してください") {
		t.Errorf("エラーメッセージが期待値と異なります: %v", err)
	}
}

// MockFlagParser はテスト用のFlagParser実装
type MockFlagParser struct {
	rootDir    string
	targetDirs string
	operation  string
	parseError error
}

// StringVar はモックの文字列フラグ定義
func (m *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
	switch name {
	case "root-dir":
		*p = m.rootDir
	case "target-dirs":
		*p = m.targetDirs
	case "operation":
		*p = m.operation
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
		operation:  "output",
		parseError: nil,
	}
	mockOSArgs := &MockOSArgs{
		args: []string{"program", "-root-dir=/test/root", "-target-dirs=cli,mcp", "-operation=output"},
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
	if config.Operation != "output" {
		t.Errorf("Operationが期待値と異なります。期待値: output, 実際: %s", config.Operation)
	}
}

// TestConfigParser_ParseFlags_ParseError はConfigParserのParseFlags解析エラーテスト
func TestConfigParser_ParseFlags_ParseError(t *testing.T) {
	// Arrange
	mockFlagParser := &MockFlagParser{
		rootDir:    "/test/root",
		targetDirs: "cli,mcp",
		operation:  "output",
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

// captureStdout はPrintUsageなどの出力を捕捉するテストヘルパー
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	originalStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("標準出力のパイプ作成に失敗しました: %v", err)
	}
	os.Stdout = w
	fn()
	if err := w.Close(); err != nil {
		t.Fatalf("標準出力のクローズに失敗しました: %v", err)
	}
	os.Stdout = originalStdout
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("標準出力の読み取りに失敗しました: %v", err)
	}
	if err := r.Close(); err != nil {
		t.Fatalf("パイプのクローズに失敗しました: %v", err)
	}
	return buf.String()
}

// TestStandardFlagParser_StringAndBoolVar はStandardFlagParserの各メソッドを検証する
func TestStandardFlagParser_StringAndBoolVar(t *testing.T) {
	originalArgs := os.Args
	originalFlagSet := flag.CommandLine
	defer func() {
		os.Args = originalArgs
		flag.CommandLine = originalFlagSet
	}()

	os.Args = []string{"service-viewer", "-root-dir=/tmp/project", "-target-dirs=cli,mcp", "-dry-run=true", "extra"}
	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flagSet.SetOutput(io.Discard)
	flag.CommandLine = flagSet

	parser := NewStandardFlagParser()

	var rootDir, targetDirs string
	var dryRun bool
	parser.StringVar(&rootDir, "root-dir", "", "test root")
	parser.StringVar(&targetDirs, "target-dirs", "", "test target")
	parser.BoolVar(&dryRun, "dry-run", false, "test bool")

	if err := parser.Parse(); err != nil {
		t.Fatalf("フラグの解析に失敗しました: %v", err)
	}

	if rootDir != "/tmp/project" {
		t.Errorf("rootDirの値が期待値と異なります: %s", rootDir)
	}
	if targetDirs != "cli,mcp" {
		t.Errorf("targetDirsの値が期待値と異なります: %s", targetDirs)
	}
	if !dryRun {
		t.Errorf("dryRunの値が期待値と異なります: %v", dryRun)
	}

	args := parser.Args()
	if len(args) != 1 || args[0] != "extra" {
		t.Errorf("Argsの解析結果が期待値と異なります: %v", args)
	}
}

// TestNewStandardOSArgs は標準実装が現在のos.Argsを返すことを確認する
func TestNewStandardOSArgs(t *testing.T) {
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	expected := []string{"viewer", "-flag"}
	os.Args = expected

	osArgs := NewStandardOSArgs()
	result := osArgs.Args()

	if len(result) != len(expected) {
		t.Fatalf("Argsの要素数が期待値と異なります: %d", len(result))
	}
	for i, v := range expected {
		if result[i] != v {
			t.Errorf("Args[%d]が期待値と異なります。期待値: %s, 実際: %s", i, v, result[i])
		}
	}
}

// TestParseFlags はグローバル関数ParseFlagsの動作を確認する
func TestParseFlags(t *testing.T) {
	originalArgs := os.Args
	originalFlagSet := flag.CommandLine
	defer func() {
		os.Args = originalArgs
		flag.CommandLine = originalFlagSet
	}()

	os.Args = []string{"service-implementing-viewer", "-root-dir=/workspace", "-target-dirs=cli,mcp", "-operation=output"}
	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flagSet.SetOutput(io.Discard)
	flag.CommandLine = flagSet

	config, err := ParseFlags()
	if err != nil {
		t.Fatalf("ParseFlagsがエラーを返しました: %v", err)
	}

	if config.RootDir != "/workspace" {
		t.Errorf("RootDirが期待値と異なります: %s", config.RootDir)
	}
	if len(config.TargetDirs) != 2 || config.TargetDirs[0] != "cli" || config.TargetDirs[1] != "mcp" {
		t.Errorf("TargetDirsが期待値と異なります: %v", config.TargetDirs)
	}
	if config.Operation != "output" {
		t.Errorf("Operationが期待値と異なります: %s", config.Operation)
	}
}

// TestPrintUsage は使用方法メッセージを検証する
func TestPrintUsage(t *testing.T) {
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	os.Args = []string{"service-implementing-viewer"}

	output := captureStdout(t, PrintUsage)

	expectedFragments := []string{
		"使用方法: service-implementing-viewer [オプション]",
		"-root-dir string",
		"ルートディレクトリ（必須）",
		"-target-dirs string",
		"対象ディレクトリ（必須、カンマ区切り）",
		"-operation string",
		"実行する操作（必須: output",
		"使用例:",
		"-operation=output",
	}

	for _, fragment := range expectedFragments {
		if !strings.Contains(output, fragment) {
			t.Errorf("PrintUsageの出力に期待する文字列が含まれていません: %s", fragment)
		}
	}
}
