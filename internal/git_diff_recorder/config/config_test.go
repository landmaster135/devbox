package config

import (
	"fmt"
	"testing"
)

// TestConfigParser_NewConfigParser_Normal はConfigParser作成の正常系テスト
func TestConfigParser_NewConfigParser_Normal(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockOSArgs := NewMockOSArgs([]string{"test"})

	// Act
	configParser := NewConfigParser(mockParser, mockOSArgs)

	// Assert
	if configParser == nil {
		t.Error("ConfigParserの作成に失敗しました")
		return
	}
	if configParser.flagParser != mockParser {
		t.Error("flagParserが正しく設定されていません")
	}
	if configParser.osArgs != mockOSArgs {
		t.Error("osArgsが正しく設定されていません")
	}
}

// TestConfigParser_ParseFlags_RecordMode は記録モードの正常系テスト
func TestConfigParser_ParseFlags_RecordMode(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockOSArgs := NewMockOSArgs([]string{"git-diff-recorder", "--output-dir", "/tmp/output"})
	configParser := NewConfigParser(mockParser, mockOSArgs)

	mockParser.SetStringValue("output-dir", "/tmp/output")

	// Act
	cfg, err := configParser.ParseFlags()

	// Assert
	if err != nil {
		t.Errorf("引数解析でエラーが発生しました: %v", err)
	}
	if cfg.OutputDir != "/tmp/output" {
		t.Errorf("OutputDirが期待値と異なります。期待値: /tmp/output, 実際: %s", cfg.OutputDir)
	}
	if cfg.ReadMode {
		t.Error("ReadModeがfalseであるべきです")
	}
	if cfg.GenMode {
		t.Error("GenModeがfalseであるべきです")
	}
}

// TestConfigParser_ParseFlags_AllFlags は全フラグ設定のテスト
func TestConfigParser_ParseFlags_AllFlags(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockOSArgs := NewMockOSArgs([]string{"git-diff-recorder"})
	configParser := NewConfigParser(mockParser, mockOSArgs)

	// 全てのフラグを設定
	mockParser.SetStringValue("output-dir", "/tmp/output")
	mockParser.SetBoolValue("staged-only", true)
	mockParser.SetBoolValue("read-mode", false)
	mockParser.SetBoolValue("gen-mode", false)
	mockParser.SetStringValue("source-dir", "/tmp/source")
	mockParser.SetStringValue("repository", "test-repo")
	mockParser.SetStringValue("git-dir", "/tmp/git")

	// Act
	cfg, err := configParser.ParseFlags()

	// Assert
	if err != nil {
		t.Errorf("引数解析でエラーが発生しました: %v", err)
	}
	if cfg.OutputDir != "/tmp/output" {
		t.Errorf("OutputDirが期待値と異なります。期待値: /tmp/output, 実際: %s", cfg.OutputDir)
	}
	if !cfg.StagedOnly {
		t.Error("StagedOnlyがtrueであるべきです")
	}
	if cfg.SourceDir != "/tmp/source" {
		t.Errorf("SourceDirが期待値と異なります。期待値: /tmp/source, 実際: %s", cfg.SourceDir)
	}
	if cfg.Repository != "test-repo" {
		t.Errorf("Repositoryが期待値と異なります。期待値: test-repo, 実際: %s", cfg.Repository)
	}
	if cfg.GitDir != "/tmp/git" {
		t.Errorf("GitDirが期待値と異なります。期待値: /tmp/git, 実際: %s", cfg.GitDir)
	}
}

// TestConfigParser_validateConfig_ReadModeMissingRepository は読み取りモードでリポジトリ不足のテスト
func TestConfigParser_validateConfig_ReadModeMissingRepository(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockOSArgs := NewMockOSArgs([]string{"git-diff-recorder"})
	configParser := NewConfigParser(mockParser, mockOSArgs)

	mockParser.SetBoolValue("read-mode", true)
	mockParser.SetStringValue("source-dir", "/tmp/source")
	// repositoryを設定しない

	// Act
	cfg, err := configParser.ParseFlags()

	// Assert
	if err == nil {
		t.Error("必須パラメータが不足しているのにエラーが発生しませんでした")
	}
	if cfg != nil {
		t.Error("エラー時にはconfigはnilであるべきです")
	}
}

// TestMockFlagParser_SetValues はモックの値設定テスト
func TestMockFlagParser_SetValues(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()

	// Act
	mockParser.SetStringValue("test-string", "test-value")
	mockParser.SetBoolValue("test-bool", true)

	// Assert
	if mockParser.testValues["test-string"] != "test-value" {
		t.Error("文字列値が正しく設定されていません")
	}
	if mockParser.testValues["test-bool"] != true {
		t.Error("ブール値が正しく設定されていません")
	}
}

// TestMockFlagParser_ParseError はモックのパースエラーテスト
func TestMockFlagParser_ParseError(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	expectedError := fmt.Errorf("test parse error")
	mockParser.SetParseError(expectedError)

	// Act
	err := mockParser.Parse()

	// Assert
	if err != expectedError {
		t.Errorf("期待されるエラーと異なります。期待値: %v, 実際: %v", expectedError, err)
	}
}

// TestMockOSArgs_SetArgs はモックのOS引数設定テスト
func TestMockOSArgs_SetArgs(t *testing.T) {
	// Arrange
	mockOSArgs := NewMockOSArgs([]string{"initial"})
	newArgs := []string{"new", "args"}

	// Act
	mockOSArgs.SetArgs(newArgs)

	// Assert
	args := mockOSArgs.Args()
	if len(args) != 2 {
		t.Errorf("引数の数が期待値と異なります。期待値: 2, 実際: %d", len(args))
	}
	if args[0] != "new" || args[1] != "args" {
		t.Errorf("引数の内容が期待値と異なります。期待値: [new, args], 実際: %v", args)
	}
}

// TestStandardFlagParser_Creation は標準フラグパーサー作成のテスト
func TestStandardFlagParser_Creation(t *testing.T) {
	// Arrange & Act
	parser := NewStandardFlagParser()

	// Assert
	if parser == nil {
		t.Error("StandardFlagParserの作成に失敗しました")
		return
	}
	if parser.flagSet == nil {
		t.Error("flagSetが設定されていません")
	}
}

// TestStandardOSArgs_Creation は標準OS引数作成のテスト
func TestStandardOSArgs_Creation(t *testing.T) {
	// Arrange & Act
	osArgs := NewStandardOSArgs()

	// Assert
	if osArgs == nil {
		t.Error("StandardOSArgsの作成に失敗しました")
	}

	// 実際のos.Argsを取得できることを確認
	args := osArgs.Args()
	if args == nil {
		t.Error("引数の取得に失敗しました")
	}
}
