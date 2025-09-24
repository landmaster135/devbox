package analyzer

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/landmaster135/devbox/internal/depends_visualizer/config"
)

// TestAnalyzeFile_Normal はAnalyzeFile関数の正常系テストです
func TestAnalyzeFile_Normal(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "analyzer_test")
	if err != nil {
		t.Fatalf("テスト用ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// テスト用のGoファイルを作成
	testFilePath := filepath.Join(tempDir, "test.go")
	testContent := `package test

func main() {
	hello()
	world()
}

func hello() {
	fmt.Println("Hello")
}

func world() {
	hello()
	fmt.Println("World")
}
`
	if err := os.WriteFile(testFilePath, []byte(testContent), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}

	// デフォルト設定を一時的に上書き
	originalConfig := config.DefaultConfig
	defer func() {
		config.DefaultConfig = originalConfig
	}()

	// テスト実行
	deps, err := AnalyzeFile(testFilePath, ".go")
	if err != nil {
		t.Fatalf("AnalyzeFile関数の実行に失敗しました: %v", err)
	}

	// 期待される依存関係
	expected := map[string][]string{
		"main":  {"hello", "world"},
		"hello": {},
		"world": {"hello"},
	}

	// 結果を検証
	if !reflect.DeepEqual(deps, expected) {
		t.Errorf("期待する依存関係: %v, 実際の依存関係: %v", expected, deps)
	}
}

// TestAnalyzeFile_InvalidExtension はAnalyzeFile関数の拡張子エラーテストです
func TestAnalyzeFile_InvalidExtension(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "analyzer_test")
	if err != nil {
		t.Fatalf("テスト用ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// テスト用のテキストファイルを作成
	testFilePath := filepath.Join(tempDir, "test.txt")
	testContent := "This is not a Go file"
	if err := os.WriteFile(testFilePath, []byte(testContent), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}

	// テスト実行
	_, err = AnalyzeFile(testFilePath, ".txt")
	if err == nil {
		t.Errorf("サポートされていない拡張子でもエラーが発生しませんでした")
	}
}

// TestAnalyzeFile_FileNotFound はAnalyzeFile関数のファイル不在テストです
func TestAnalyzeFile_FileNotFound(t *testing.T) {
	// 存在しないファイルパス
	nonExistentPath := "/path/to/nonexistent/file.go"

	// テスト実行
	_, err := AnalyzeFile(nonExistentPath, ".go")
	if err == nil {
		t.Errorf("存在しないファイルでもエラーが発生しませんでした")
	}
}

// TestAnalyzeDependencies_Normal はAnalyzeDependencies関数の正常系テストです
func TestAnalyzeDependencies_Normal(t *testing.T) {
	// テスト用のコード行
	lines := []string{
		"package test",
		"",
		"func main() {",
		"	hello()",
		"	world()",
		"}",
		"",
		"func hello() {",
		"	fmt.Println(\"Hello\")",
		"}",
		"",
		"func world() {",
		"	hello()",
		"	fmt.Println(\"World\")",
		"}",
	}

	// 関数リスト
	functions := []string{"main", "hello", "world"}

	// デフォルト設定を一時的に上書き
	originalConfig := config.DefaultConfig
	defer func() {
		config.DefaultConfig = originalConfig
	}()

	// テスト実行
	deps, err := AnalyzeDependencies(lines, functions, ".go")
	if err != nil {
		t.Fatalf("AnalyzeDependencies関数の実行に失敗しました: %v", err)
	}

	// 期待される依存関係
	expected := map[string][]string{
		"main":  {"hello", "world"},
		"hello": {},
		"world": {"hello"},
	}

	// 結果を検証
	if !reflect.DeepEqual(deps, expected) {
		t.Errorf("期待する依存関係: %v, 実際の依存関係: %v", expected, deps)
	}
}

// TestAnalyzeDependencies_InvalidExtension はAnalyzeDependencies関数の拡張子エラーテストです
func TestAnalyzeDependencies_InvalidExtension(t *testing.T) {
	// テスト用のコード行
	lines := []string{"This is not a valid code"}

	// 関数リスト
	functions := []string{"func1"}

	// テスト実行
	_, err := AnalyzeDependencies(lines, functions, ".invalid")
	if err == nil {
		t.Errorf("サポートされていない拡張子でもエラーが発生しませんでした")
	}
}

// TestAnalyzeDependencies_Deduplicates は重複参照が1回だけ登録されることを確認します
func TestAnalyzeDependencies_Deduplicates(t *testing.T) {
	lines := []string{
		"package test",
		"",
		"func main() {",
		"	foo()",
		"	foo()",
		"}",
		"",
		"func foo() {}",
	}

	functions := []string{"main", "foo"}

	deps, err := AnalyzeDependencies(lines, functions, ".go")
	if err != nil {
		t.Fatalf("AnalyzeDependenciesの実行に失敗しました: %v", err)
	}

	expected := map[string][]string{
		"main": {"foo"},
		"foo":  {},
	}

	if !reflect.DeepEqual(deps, expected) {
		t.Errorf("期待する依存関係: %v, 実際の依存関係: %v", expected, deps)
	}
}
