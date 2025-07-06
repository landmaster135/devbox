package usecases

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// #==============================================================#
// ##          Tests                                             ##
// #==============================================================#
func TestReadTextFile_Normal(t *testing.T) {
	// テスト用の一時ファイルを作成
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	testContent := "Hello, World!\nThis is a test file."

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}

	// ファイルを読み取り
	result, err := readTextFile(testFile)
	if err != nil {
		t.Fatalf("ReadTextFile()でエラーが発生しました: %v", err)
	}

	if result != testContent {
		t.Errorf("期待される内容と異なります。期待値: %q, 実際の値: %q", testContent, result)
	}
}

func TestReadTextFile_EmptyPath(t *testing.T) {
	result, err := readTextFile("")
	if err != nil {
		t.Fatalf("空のパスでエラーが発生しました: %v", err)
	}

	if result != "" {
		t.Errorf("空のパスの場合は空文字列を返すべきです。実際の値: %q", result)
	}
}

func TestReadTextFile_NonExistentFile(t *testing.T) {
	_, err := readTextFile("non_existent_file.txt")
	if err == nil {
		t.Error("存在しないファイルの場合はエラーを返すべきです")
	}
}

func TestGenerateConfigJS_Normal(t *testing.T) {
	config := DiffConfig{
		LeftText:  "Hello\nWorld",
		RightText: "Hello\nUniverse",
	}

	result := generateConfigJS(config)

	// 期待される内容を確認
	if !strings.Contains(result, "const CONFIG = {") {
		t.Error("CONFIG オブジェクトが含まれていません")
	}

	if !strings.Contains(result, "leftText:") {
		t.Error("leftText プロパティが含まれていません")
	}

	if !strings.Contains(result, "rightText:") {
		t.Error("rightText プロパティが含まれていません")
	}

	// エスケープされた改行が含まれているか確認
	if !strings.Contains(result, "\\n") {
		t.Error("改行文字がエスケープされていません")
	}
}

func TestGenerateConfigJS_EmptyTexts(t *testing.T) {
	config := DiffConfig{
		LeftText:  "",
		RightText: "",
	}

	result := generateConfigJS(config)

	// 空文字列が正しく処理されているか確認
	if !strings.Contains(result, `leftText: ""`) {
		t.Error("空の leftText が正しく処理されていません")
	}

	if !strings.Contains(result, `rightText: ""`) {
		t.Error("空の rightText が正しく処理されていません")
	}
}

func TestEscapeJSString_Normal(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", `""`},
		{"Hello", `"Hello"`},
		{"Hello\nWorld", `"Hello\nWorld"`},
		{"Hello\tWorld", `"Hello\tWorld"`},
		{`Hello"World`, `"Hello\"World"`},
		{"Hello\\World", `"Hello\\World"`},
		{"Hello\r\nWorld", `"Hello\r\nWorld"`},
	}

	for _, test := range tests {
		result := escapeJSString(test.input)
		if result != test.expected {
			t.Errorf("escapeJSString(%q) = %q, 期待値: %q", test.input, result, test.expected)
		}
	}
}

func TestGenerateHTML_Normal(t *testing.T) {
	tmplStr := `<!DOCTYPE html>
<html>
<head>
    <link rel="stylesheet" href="style.css">
</head>
<body>
    <script src="config.js"></script>
    <script src="script.js"></script>
</body>
</html>`

	data := TemplateData{
		Style:    "body { color: red; }",
		Script:   "console.log('test');",
		ConfigJS: "const CONFIG = {};",
	}

	result, err := generateHTML(tmplStr, data)
	if err != nil {
		t.Fatalf("GenerateHTML()でエラーが発生しました: %v", err)
	}

	// CSS が埋め込まれているか確認
	if !strings.Contains(result, "<style>") {
		t.Error("CSS が埋め込まれていません")
	}

	if !strings.Contains(result, "body { color: red; }") {
		t.Error("CSS の内容が正しく埋め込まれていません")
	}

	// JavaScript が埋め込まれているか確認
	if !strings.Contains(result, "console.log('test');") {
		t.Error("JavaScript が埋め込まれていません")
	}

	if !strings.Contains(result, "const CONFIG = {};") {
		t.Error("設定 JavaScript が埋め込まれていません")
	}

	// 元のリンクタグが削除されているか確認
	if strings.Contains(result, `href="style.css"`) {
		t.Error("元の CSS リンクタグが削除されていません")
	}

	if strings.Contains(result, `src="script.js"`) {
		t.Error("元の JavaScript タグが削除されていません")
	}
}

func TestGenerateHTML_EmptyTemplate(t *testing.T) {
	data := TemplateData{
		Style:    "test",
		Script:   "test",
		ConfigJS: "test",
	}

	result, err := generateHTML("", data)
	if err != nil {
		t.Fatalf("空のテンプレートでエラーが発生しました: %v", err)
	}

	if result != "" {
		t.Error("空のテンプレートの場合は空文字列を返すべきです")
	}
}

func TestDiffConfig_Struct(t *testing.T) {
	config := DiffConfig{
		LeftText:  "test left",
		RightText: "test right",
	}

	if config.LeftText != "test left" {
		t.Error("LeftText が正しく設定されていません")
	}

	if config.RightText != "test right" {
		t.Error("RightText が正しく設定されていません")
	}
}

func TestTemplateData_Struct(t *testing.T) {
	data := TemplateData{
		Style:    "test style",
		Script:   "test script",
		ConfigJS: "test config",
	}

	if data.Style != "test style" {
		t.Error("Style が正しく設定されていません")
	}

	if data.Script != "test script" {
		t.Error("Script が正しく設定されていません")
	}

	if data.ConfigJS != "test config" {
		t.Error("ConfigJS が正しく設定されていません")
	}
}

func TestProcessDiffDreamer_Normal(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tmpDir := t.TempDir()

	// テスト用ファイルを作成
	leftFile := filepath.Join(tmpDir, "left.txt")
	rightFile := filepath.Join(tmpDir, "right.txt")
	outputFile := filepath.Join(tmpDir, "output.html")

	leftContent := "Hello\nWorld"
	rightContent := "Hello\nUniverse"

	err := os.WriteFile(leftFile, []byte(leftContent), 0644)
	if err != nil {
		t.Fatalf("左側テストファイルの作成に失敗しました: %v", err)
	}

	err = os.WriteFile(rightFile, []byte(rightContent), 0644)
	if err != nil {
		t.Fatalf("右側テストファイルの作成に失敗しました: %v", err)
	}

	// テスト用のHTMLテンプレート
	indexHTML := `<!DOCTYPE html>
<html>
<head>
    <link rel="stylesheet" href="style.css">
</head>
<body>
    <script src="config.js"></script>
    <script src="script.js"></script>
</body>
</html>`

	styleCSS := "body { color: red; }"
	scriptJS := "console.log('test');"

	// ProcessDiffDreamerを実行（ブラウザ起動はスキップするためにモック化が必要だが、ここでは基本的な処理のテストのみ）
	// 実際のテストではブラウザ起動部分をモック化する必要があります
	err = ProcessDiffDreamer(leftFile, rightFile, outputFile, indexHTML, styleCSS, scriptJS)
	if err != nil {
		t.Fatalf("HTMLファイルの作成に失敗しました: %v", err)
	}

	// ブラウザ起動でエラーが発生する可能性があるため、ファイル作成までの処理をテスト
	// HTMLファイルが作成されているかを確認
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		// ファイルが作成されていない場合、ブラウザ起動以外のエラーかもしれない
		if err != nil && !strings.Contains(err.Error(), "ブラウザの起動に失敗しました") {
			t.Fatalf("ProcessDiffDreamer()で予期しないエラーが発生しました: %v", err)
		}
	}

	// HTMLファイルの内容を確認
	if _, err := os.Stat(outputFile); err == nil {
		htmlContent, err := os.ReadFile(outputFile)
		if err != nil {
			t.Fatalf("出力HTMLファイルの読み取りに失敗しました: %v", err)
		}

		htmlStr := string(htmlContent)

		// CSS が埋め込まれているか確認
		if !strings.Contains(htmlStr, "body { color: red; }") {
			t.Error("CSS が正しく埋め込まれていません")
		}

		// JavaScript が埋め込まれているか確認
		if !strings.Contains(htmlStr, "console.log('test');") {
			t.Error("JavaScript が正しく埋め込まれていません")
		}

		// 設定が埋め込まれているか確認
		if !strings.Contains(htmlStr, "const CONFIG = {") {
			t.Error("設定 JavaScript が埋め込まれていません")
		}

		// テキスト内容が含まれているか確認
		if !strings.Contains(htmlStr, "Hello\\nWorld") {
			t.Error("左側テキストが正しく埋め込まれていません")
		}

		if !strings.Contains(htmlStr, "Hello\\nUniverse") {
			t.Error("右側テキストが正しく埋め込まれていません")
		}
	}
}

func TestProcessDiffDreamer_LeftFileError(t *testing.T) {
	tmpDir := t.TempDir()
	rightFile := filepath.Join(tmpDir, "right.txt")
	outputFile := filepath.Join(tmpDir, "output.html")

	// 右側ファイルのみ作成
	err := os.WriteFile(rightFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("右側テストファイルの作成に失敗しました: %v", err)
	}

	// 存在しない左側ファイルを指定
	err = ProcessDiffDreamer("non_existent_left.txt", rightFile, outputFile, "", "", "")
	if err == nil {
		t.Error("存在しない左側ファイルの場合はエラーを返すべきです")
	}

	if !strings.Contains(err.Error(), "左側ファイルの読み取りエラー") {
		t.Errorf("期待されるエラーメッセージが含まれていません: %v", err)
	}
}

func TestProcessDiffDreamer_RightFileError(t *testing.T) {
	tmpDir := t.TempDir()
	leftFile := filepath.Join(tmpDir, "left.txt")
	outputFile := filepath.Join(tmpDir, "output.html")

	// 左側ファイルのみ作成
	err := os.WriteFile(leftFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("左側テストファイルの作成に失敗しました: %v", err)
	}

	// 存在しない右側ファイルを指定
	err = ProcessDiffDreamer(leftFile, "non_existent_right.txt", outputFile, "", "", "")
	if err == nil {
		t.Error("存在しない右側ファイルの場合はエラーを返すべきです")
	}

	if !strings.Contains(err.Error(), "右側ファイルの読み取りエラー") {
		t.Errorf("期待されるエラーメッセージが含まれていません: %v", err)
	}
}

func TestProcessDiffDreamer_EmptyFiles(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.html")

	// 空のファイルパスを指定
	indexHTML := `<!DOCTYPE html>
<html>
<head>
    <link rel="stylesheet" href="style.css">
</head>
<body>
    <script src="config.js"></script>
    <script src="script.js"></script>
</body>
</html>`

	err := ProcessDiffDreamer("", "", outputFile, indexHTML, "", "")

	// ブラウザ起動でエラーが発生する可能性があるが、ファイル処理自体は成功するはず
	if err != nil && !strings.Contains(err.Error(), "ブラウザの起動に失敗しました") {
		t.Errorf("空のファイルパスでの処理で予期しないエラーが発生しました: %v", err)
	}

	// HTMLファイルが作成されているかを確認
	if _, err := os.Stat(outputFile); err == nil {
		htmlContent, err := os.ReadFile(outputFile)
		if err != nil {
			t.Fatalf("出力HTMLファイルの読み取りに失敗しました: %v", err)
		}

		htmlStr := string(htmlContent)

		// 空のテキストが正しく処理されているか確認
		if !strings.Contains(htmlStr, `leftText: ""`) {
			t.Error("空の左側テキストが正しく処理されていません")
		}

		if !strings.Contains(htmlStr, `rightText: ""`) {
			t.Error("空の右側テキストが正しく処理されていません")
		}
	}
}
