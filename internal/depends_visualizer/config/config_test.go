package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// TestGetLanguageConfig_Normal はGetLanguageConfig関数の正常系テストです
func TestGetLanguageConfig_Normal(t *testing.T) {
	// テスト用のデフォルト設定を一時的に保存
	originalConfig := DefaultConfig
	defer func() {
		DefaultConfig = originalConfig
	}()

	// テストケース
	testCases := []struct {
		name      string
		extension string
		wantOK    bool
	}{
		{"Go", ".go", true},
		{"Python", ".py", true},
		{"JavaScript", ".js", true},
		{"Unsupported", ".txt", false},
		{"Empty", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lang, ok := GetLanguageConfig(tc.extension)
			if ok != tc.wantOK {
				t.Errorf("GetLanguageConfig(%s) ok = %v, wantOK %v", tc.extension, ok, tc.wantOK)
				return
			}

			if ok {
				// 言語設定が取得できた場合、その内容を検証
				if lang.FunctionHeader == "" {
					t.Errorf("GetLanguageConfig(%s) 関数ヘッダーが空です", tc.extension)
				}
				if lang.FunctionTail == "" {
					t.Errorf("GetLanguageConfig(%s) 関数テイルが空です", tc.extension)
				}
				if lang.CommentPrefix == "" {
					t.Errorf("GetLanguageConfig(%s) コメントプレフィックスが空です", tc.extension)
				}
			}
		})
	}
}

// TestGetSpaces_Normal はGetSpaces関数の正常系テストです
func TestGetSpaces_Normal(t *testing.T) {
	// テスト用のデフォルト設定を一時的に保存
	originalConfig := DefaultConfig
	defer func() {
		DefaultConfig = originalConfig
	}()

	// テスト用の設定
	DefaultConfig.Spaces = []string{" ", "\t"}

	// テスト実行
	spaces := GetSpaces()

	// 結果を検証
	if !reflect.DeepEqual(spaces, DefaultConfig.Spaces) {
		t.Errorf("GetSpaces() = %v, want %v", spaces, DefaultConfig.Spaces)
	}
}

// TestLoadFromFile_Normal はLoadFromFile関数の正常系テストです
func TestLoadFromFile_Normal(t *testing.T) {
	// テスト用のデフォルト設定を一時的に保存
	originalConfig := DefaultConfig
	defer func() {
		DefaultConfig = originalConfig
	}()

	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("テスト用ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// テスト用の設定ファイルを作成
	configPath := filepath.Join(tempDir, "config.json")
	configContent := `{
		"spaces": [" ", "\t", "  "],
		"languages": {
			".go": {
				"function_header": "func ",
				"function_tail": "(",
				"main_marker": "func main() {",
				"comment_prefix": "//",
				"multiline_comment": true
			},
			".custom": {
				"function_header": "function ",
				"function_tail": "(",
				"main_marker": "main() {",
				"comment_prefix": "#",
				"multiline_comment": false
			}
		},
		"output_format": "dot"
	}`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("テスト用設定ファイルの作成に失敗しました: %v", err)
	}

	// テスト実行
	err = LoadFromFile(configPath)
	if err != nil {
		t.Fatalf("LoadFromFile(%s) error = %v", configPath, err)
	}

	// 結果を検証
	if len(DefaultConfig.Spaces) != 3 {
		t.Errorf("LoadFromFile後のSpacesの長さが期待と異なります: got %d, want 3", len(DefaultConfig.Spaces))
	}
	if DefaultConfig.OutputFormat != "dot" {
		t.Errorf("LoadFromFile後のOutputFormatが期待と異なります: got %s, want dot", DefaultConfig.OutputFormat)
	}
	if len(DefaultConfig.Languages) != 2 {
		t.Errorf("LoadFromFile後のLanguagesの長さが期待と異なります: got %d, want 2", len(DefaultConfig.Languages))
	}
	if _, ok := DefaultConfig.Languages[".custom"]; !ok {
		t.Errorf("LoadFromFile後にカスタム言語設定が見つかりません")
	}
}

// TestLoadFromFile_FileNotFound はLoadFromFile関数のファイル不在テストです
func TestLoadFromFile_FileNotFound(t *testing.T) {
	// 存在しないファイルパス
	nonExistentPath := "/path/to/nonexistent/config.json"

	// テスト実行
	err := LoadFromFile(nonExistentPath)
	if err == nil {
		t.Errorf("存在しないファイルでもエラーが発生しませんでした")
	}
}

// TestLoadFromFile_InvalidJSON はLoadFromFile関数の無効なJSON形式テストです
func TestLoadFromFile_InvalidJSON(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("テスト用ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// テスト用の無効なJSON形式の設定ファイルを作成
	configPath := filepath.Join(tempDir, "invalid_config.json")
	configContent := `{
		"spaces": [" ", "\t"],
		"languages": {
			".go": {
				"function_header": "func ",
				"function_tail": "(",
				"main_marker": "func main() {",
				"comment_prefix": "//",
				"multiline_comment": true
			},
		}, // 無効なJSON形式（余分なカンマ）
		"output_format": "dot"
	}`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("テスト用設定ファイルの作成に失敗しました: %v", err)
	}

	// テスト実行
	err = LoadFromFile(configPath)
	if err == nil {
		t.Errorf("無効なJSON形式でもエラーが発生しませんでした")
	}
}

// TestAppConfig_Normal はAppConfig構造体の正常系テストです
func TestAppConfig_Normal(t *testing.T) {
	// テスト用の設定
	cfg := AppConfig{
		ConfigPath:  "/path/to/config.json",
		SourceFile:  "/path/to/source.go",
		Extension:   ".go",
		OutputPath:  "/path/to/output.md",
		Format:      "mermaid",
		Recursive:   true,
		Verbose:     true,
		Directory:   "/path/to/dir",
	}

	// 各フィールドを検証
	if cfg.ConfigPath != "/path/to/config.json" {
		t.Errorf("AppConfig.ConfigPath = %s, want /path/to/config.json", cfg.ConfigPath)
	}
	if cfg.SourceFile != "/path/to/source.go" {
		t.Errorf("AppConfig.SourceFile = %s, want /path/to/source.go", cfg.SourceFile)
	}
	if cfg.Extension != ".go" {
		t.Errorf("AppConfig.Extension = %s, want .go", cfg.Extension)
	}
	if cfg.OutputPath != "/path/to/output.md" {
		t.Errorf("AppConfig.OutputPath = %s, want /path/to/output.md", cfg.OutputPath)
	}
	if cfg.Format != "mermaid" {
		t.Errorf("AppConfig.Format = %s, want mermaid", cfg.Format)
	}
	if !cfg.Recursive {
		t.Errorf("AppConfig.Recursive = %v, want true", cfg.Recursive)
	}
	if !cfg.Verbose {
		t.Errorf("AppConfig.Verbose = %v, want true", cfg.Verbose)
	}
	if cfg.Directory != "/path/to/dir" {
		t.Errorf("AppConfig.Directory = %s, want /path/to/dir", cfg.Directory)
	}
}
