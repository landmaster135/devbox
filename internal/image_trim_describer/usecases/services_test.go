package usecases

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetImageFiles(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "test_images")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// テスト用ファイルを作成
	testFiles := []string{
		"image1.jpg",
		"image2.png",
		"image3.jpeg",
		"image4.PNG",  // 大文字の拡張子
		"image5.gif",  // サポートしていない形式
		"text.txt",    // 非画像ファイル
		"document.pdf", // 非画像ファイル
	}

	for _, filename := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		file, err := os.Create(filePath)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
		file.Close()
	}

	// テストの実行
	imageFiles, err := GetImageFiles(tempDir)
	if err != nil {
		t.Errorf("GetImageFiles() error = %v", err)
		return
	}

	// 期待される結果
	expected := []string{"image1.jpg", "image2.png", "image3.jpeg", "image4.PNG"}

	// 結果の検証
	if len(imageFiles) != len(expected) {
		t.Errorf("GetImageFiles() returned %d files, expected %d", len(imageFiles), len(expected))
	}

	// 各ファイルが含まれているか確認
	for _, expectedFile := range expected {
		found := false
		for _, file := range imageFiles {
			if file == expectedFile {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected file %s not found in result", expectedFile)
		}
	}
}

func TestFormatJSArray(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "Empty array",
			input:    []string{},
			expected: "[]",
		},
		{
			name:     "Single item",
			input:    []string{"file1.jpg"},
			expected: "[\n    \"file1.jpg\"\n  ]",
		},
		{
			name:     "Multiple items",
			input:    []string{"file1.jpg", "file2.png", "file3.jpeg"},
			expected: "[\n    \"file1.jpg\",\n    \"file2.png\",\n    \"file3.jpeg\"\n  ]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatJSArray(tt.input)
			if result != tt.expected {
				t.Errorf("formatJSArray() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestGenerateConfigJS(t *testing.T) {
	config := Config{
		ImageFiles: []string{"image1.jpg", "image2.png"},
		SrcDir:     "/src/dir",
		OutDir:     "/out/dir",
		ArcDir:     "/archive/dir",
		Suffix:     "_trimmed",
		Move:       true,
	}

	result := GenerateConfigJS(config)

	// 生成されたJavaScriptコードの構造を確認
	if !strings.Contains(result, "const CONFIG = {") {
		t.Error("Generated config does not contain 'const CONFIG = {'")
	}
	if !strings.Contains(result, "imageFiles: [") {
		t.Error("Generated config does not contain 'imageFiles: ['")
	}
	if !strings.Contains(result, "\"image1.jpg\"") {
		t.Error("Generated config does not contain first image file")
	}
	if !strings.Contains(result, "srcDir: \"/src/dir\"") {
		t.Error("Generated config does not contain correct srcDir")
	}
	if !strings.Contains(result, "move: true") {
		t.Error("Generated config does not contain correct move flag")
	}
}

func TestGenerateHTML(t *testing.T) {
	tmplStr := `<html>
<head>
<link rel="stylesheet" href="style.css">
<script src="config.js"></script>
</head>
<body>
<script src="script.js"></script>
</body>
</html>`

	data := TemplateData{
		Style:    "body { color: red; }",
		Script:   "console.log('script');",
		ConfigJS: "const config = {};",
	}

	result, err := GenerateHTML(tmplStr, data)
	if err != nil {
		t.Errorf("GenerateHTML() error = %v", err)
		return
	}

	// 置換が正しく行われているか確認
	if !strings.Contains(result, "<style>") {
		t.Error("HTML does not contain <style> tag")
	}
	if !strings.Contains(result, "body { color: red; }") {
		t.Error("Style content was not properly replaced")
	}
	if !strings.Contains(result, "console.log('script');") {
		t.Error("Script content was not properly replaced")
	}
	if !strings.Contains(result, "const config = {};") {
		t.Error("Config JS was not properly replaced")
	}
}

func TestOpenHTMLPage(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "test_html")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// テスト用のHTMLコンテンツ
	htmlContent := "<html><body><h1>Test</h1></body></html>"
	outputFile := filepath.Join(tempDir, "test.html")

	// 注意: ブラウザを実際に開くテストは環境に依存するため、
	// ここではファイル作成部分のみをテストします
	err = OpenHTMLPage(htmlContent, outputFile)

	// ブラウザが開かない環境でもファイルは作成される
	if err != nil && !strings.Contains(err.Error(), "ブラウザの起動に失敗しました") {
		t.Errorf("OpenHTMLPage() error = %v", err)
	}

	// ファイルが作成されたか確認
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("HTML file was not created")
	}

	// ファイル内容の確認
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Errorf("Failed to read created file: %v", err)
	}
	if string(content) != htmlContent {
		t.Errorf("File content does not match expected content")
	}
}

func TestCopyFile(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "test_copy")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// ソースファイルを作成
	srcContent := "Test content for copying"
	srcFile := filepath.Join(tempDir, "source.txt")
	err = os.WriteFile(srcFile, []byte(srcContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// コピー先のパス
	dstFile := filepath.Join(tempDir, "destination.txt")

	// ファイルをコピー
	err = CopyFile(srcFile, dstFile)
	if err != nil {
		t.Errorf("CopyFile() error = %v", err)
		return
	}

	// コピー先ファイルの存在確認
	if _, err := os.Stat(dstFile); os.IsNotExist(err) {
		t.Error("Destination file does not exist")
		return
	}

	// 内容の確認
	dstContent, err := os.ReadFile(dstFile)
	if err != nil {
		t.Errorf("Failed to read destination file: %v", err)
		return
	}

	if string(dstContent) != srcContent {
		t.Errorf("Copied content does not match: got %q, want %q", dstContent, srcContent)
	}
}

func TestCopyFileErrors(t *testing.T) {
	// 存在しないファイルをコピーしようとする
	err := CopyFile("nonexistent.txt", "destination.txt")
	if err == nil {
		t.Error("Expected error when copying non-existent file")
	}

	// 無効なディレクトリへのコピー
	tempDir, err := os.MkdirTemp("", "test_copy_error")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	srcFile := filepath.Join(tempDir, "source.txt")
	err = os.WriteFile(srcFile, []byte("content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// 読み取り専用ディレクトリにコピー
	roDir := filepath.Join(tempDir, "readonly")
	err = os.Mkdir(roDir, 0555) // 読み取り専用
	if err != nil {
		t.Fatalf("Failed to create readonly dir: %v", err)
	}

	dstFile := filepath.Join(roDir, "destination.txt")
	err = CopyFile(srcFile, dstFile)
	if err == nil {
		t.Error("Expected error when copying to readonly directory")
	}
}

func TestGetImageFilesErrors(t *testing.T) {
	// 存在しないディレクトリ
	_, err := GetImageFiles("nonexistent_directory")
	if err == nil {
		t.Error("Expected error when reading non-existent directory")
	}
}

// ベンチマークテスト（オプション）
func BenchmarkFormatJSArray(b *testing.B) {
	items := make([]string, 100)
	for i := 0; i < 100; i++ {
		items[i] = "file" + string(rune(i)) + ".jpg"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		formatJSArray(items)
	}
}

func BenchmarkGenerateConfigJS(b *testing.B) {
	config := Config{
		ImageFiles: make([]string, 100),
		SrcDir:     "/test/src",
		OutDir:     "/test/out",
		ArcDir:     "/test/arc",
		Suffix:     "_test",
		Move:       true,
	}

	for i := 0; i < 100; i++ {
		config.ImageFiles[i] = "file" + string(rune(i)) + ".jpg"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateConfigJS(config)
	}
}
