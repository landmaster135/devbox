package usecases

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestService_ExtractBlogContent_WithRealFiles(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "notion-blog-extractor-test-real-files")
	if err != nil {
		t.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// orgディレクトリのパス
	orgDir := "./test_data/org"

	// orgディレクトリの存在確認
	if _, err := os.Stat(orgDir); os.IsNotExist(err) {
		t.Fatalf("orgディレクトリが存在しません: %s", orgDir)
	}

	// orgディレクトリ内のファイルを一時ディレクトリにコピー
	srcTempDir := filepath.Join(tempDir, "src")
	if err := os.MkdirAll(srcTempDir, 0755); err != nil {
		t.Fatalf("ソース一時ディレクトリの作成に失敗しました: %v", err)
	}

	// orgディレクトリ内のMarkdownファイルをコピー
	orgFiles, err := filepath.Glob(filepath.Join(orgDir, "*.md"))
	if err != nil {
		t.Fatalf("orgディレクトリ内のファイル検索に失敗しました: %v", err)
	}

	if len(orgFiles) == 0 {
		t.Fatalf("orgディレクトリ内にMarkdownファイルが見つかりません")
	}

	for _, orgFile := range orgFiles {
		content, err := os.ReadFile(orgFile)
		if err != nil {
			t.Fatalf("ファイルの読み込みに失敗しました %s: %v", orgFile, err)
		}

		fileName := filepath.Base(orgFile)
		destFile := filepath.Join(srcTempDir, fileName)
		if err := os.WriteFile(destFile, content, 0644); err != nil {
			t.Fatalf("ファイルのコピーに失敗しました %s: %v", fileName, err)
		}
	}

	// 出力ディレクトリ
	destTempDir := filepath.Join(tempDir, "dest")

	// テスト対象のサービスを作成
	service := NewService()

	// テスト実行
	result, err := service.ExtractBlogContent(srcTempDir, destTempDir)

	// エラーの検証
	if err != nil {
		t.Errorf("ExtractBlogContent() でエラーが発生しました: %v", err)
		return
	}

	// 結果の検証
	if !strings.Contains(result, "処理完了:") {
		t.Errorf("期待される結果メッセージが含まれていません。結果: %s", result)
	}

	// 抽出されたファイル数の確認
	if !strings.Contains(result, "3件のファイルからコンテンツを抽出しました") {
		t.Errorf("期待される抽出ファイル数が含まれていません。結果: %s", result)
	}

	// 出力ディレクトリの存在確認
	if _, err := os.Stat(destTempDir); os.IsNotExist(err) {
		t.Error("出力ディレクトリが作成されていません")
		return
	}

	// 出力ファイルの確認
	outputFiles, err := filepath.Glob(filepath.Join(destTempDir, "*.md"))
	if err != nil {
		t.Fatalf("出力ファイルの検索に失敗しました: %v", err)
	}

	if len(outputFiles) != 3 {
		t.Errorf("期待される出力ファイル数は3ですが、実際は%d個でした", len(outputFiles))
	}

	// 各出力ファイルの内容を検証
	for _, outputFile := range outputFiles {
		content, err := os.ReadFile(outputFile)
		if err != nil {
			t.Errorf("出力ファイルの読み込みに失敗しました %s: %v", outputFile, err)
			continue
		}

		contentStr := string(content)

		// コンテンツマーカーで始まることを確認
		if !strings.HasPrefix(contentStr, "# Content") {
			t.Errorf("ファイル %s: コンテンツマーカーで始まっていません", filepath.Base(outputFile))
		}

		// "## はじまり" が含まれることを確認
		if !strings.Contains(contentStr, "## はじまり") {
			t.Errorf("ファイル %s: '## はじまり' マーカーが含まれていません", filepath.Base(outputFile))
		}

		// 元のメタデータ部分が除外されていることを確認（priority, status, tagsなど）
		if strings.Contains(contentStr, "priority:") || strings.Contains(contentStr, "status:") {
			t.Errorf("ファイル %s: メタデータ部分が除外されていません", filepath.Base(outputFile))
		}
	}
}

func TestService_ExtractBlogContent_WithRealFiles_SpecificContent(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "notion-blog-extractor-test-specific")
	if err != nil {
		t.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// orgディレクトリのパス
	orgDir := "./test_data/org"

	// orgディレクトリから実際に存在するファイルを取得
	orgFiles, err := filepath.Glob(filepath.Join(orgDir, "*.md"))
	if err != nil {
		t.Fatalf("orgディレクトリ内のファイル検索に失敗しました: %v", err)
	}

	if len(orgFiles) == 0 {
		t.Skip("orgディレクトリ内にMarkdownファイルが見つかりません")
	}

	// 最初のファイルを使用（Pythonファイルを想定）
	pythonFile := orgFiles[0]
	for _, file := range orgFiles {
		if strings.Contains(strings.ToLower(file), "python") {
			pythonFile = file
			break
		}
	}

	// ソース一時ディレクトリを作成
	srcTempDir := filepath.Join(tempDir, "src")
	if err := os.MkdirAll(srcTempDir, 0755); err != nil {
		t.Fatalf("ソース一時ディレクトリの作成に失敗しました: %v", err)
	}

	// ファイルをコピー
	content, err := os.ReadFile(pythonFile)
	if err != nil {
		t.Fatalf("ファイルの読み込みに失敗しました: %v", err)
	}

	testFileName := "python_test.md"
	testFilePath := filepath.Join(srcTempDir, testFileName)
	if err := os.WriteFile(testFilePath, content, 0644); err != nil {
		t.Fatalf("ファイルのコピーに失敗しました: %v", err)
	}

	// 出力ディレクトリ
	destTempDir := filepath.Join(tempDir, "dest")

	// テスト対象のサービスを作成
	service := NewService()

	// テスト実行
	result, err := service.ExtractBlogContent(srcTempDir, destTempDir)

	// エラーの検証
	if err != nil {
		t.Errorf("ExtractBlogContent() でエラーが発生しました: %v", err)
		return
	}

	// 結果の検証
	if !strings.Contains(result, "1件のファイルからコンテンツを抽出しました") {
		t.Errorf("期待される結果が得られませんでした。結果: %s", result)
	}

	// 出力ファイルの内容を詳細に検証
	outputFile := filepath.Join(destTempDir, testFileName)
	outputContent, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("出力ファイルの読み込みに失敗しました: %v", err)
	}

	outputStr := string(outputContent)

	// 期待される内容が含まれていることを確認
	expectedContents := []string{
		"# Content",
		"## はじまり",
		"Pythonのタブサイズ",
		"VSCode",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(outputStr, expected) {
			t.Errorf("期待される内容が含まれていません: %s", expected)
		}
	}

	// ファイル名に基づいた内容の確認
	fileName := filepath.Base(pythonFile)
	t.Logf("使用されたファイル: %s", fileName)

	// 出力内容の一部を表示（最大200文字）
	maxLen := 200
	if len(outputStr) < maxLen {
		maxLen = len(outputStr)
	}
	t.Logf("出力内容の一部: %s", outputStr[:maxLen])

	// 除外されるべき内容が含まれていないことを確認
	excludedContents := []string{
		"priority:",
		"status:",
		"tags:",
		"estimated_dev_minutes:",
		"# Draft",
	}

	for _, excluded := range excludedContents {
		if strings.Contains(outputStr, excluded) {
			t.Errorf("除外されるべき内容が含まれています: %s", excluded)
		}
	}
}

func TestService_ExtractBlogContent_WithRealFiles_NoMarkerFile(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "notion-blog-extractor-test-no-marker")
	if err != nil {
		t.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// ソース一時ディレクトリを作成
	srcTempDir := filepath.Join(tempDir, "src")
	if err := os.MkdirAll(srcTempDir, 0755); err != nil {
		t.Fatalf("ソース一時ディレクトリの作成に失敗しました: %v", err)
	}

	// マーカーを含まないテストファイルを作成
	noMarkerContent := `# テストファイル

これはマーカーを含まないファイルです。

## セクション1

通常のMarkdownコンテンツです。

## セクション2

このファイルは抽出対象になりません。
`

	testFilePath := filepath.Join(srcTempDir, "no_marker.md")
	if err := os.WriteFile(testFilePath, []byte(noMarkerContent), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}

	// 出力ディレクトリ
	destTempDir := filepath.Join(tempDir, "dest")

	// テスト対象のサービスを作成
	service := NewService()

	// テスト実行
	result, err := service.ExtractBlogContent(srcTempDir, destTempDir)

	// エラーの検証
	if err != nil {
		t.Errorf("ExtractBlogContent() でエラーが発生しました: %v", err)
		return
	}

	// 結果の検証：マーカーを含むファイルが見つからない場合のメッセージ
	expectedMessage := "指定されたディレクトリにコンテンツマーカーを含むMarkdownファイルが見つかりませんでした。"
	if !strings.Contains(result, expectedMessage) {
		t.Errorf("期待されるメッセージが含まれていません。結果: %s", result)
	}

	// 出力ディレクトリは作成されるが、ファイルは生成されないことを確認
	if _, err := os.Stat(destTempDir); os.IsNotExist(err) {
		t.Error("出力ディレクトリが作成されていません")
		return
	}

	outputFiles, err := filepath.Glob(filepath.Join(destTempDir, "*.md"))
	if err != nil {
		t.Fatalf("出力ファイルの検索に失敗しました: %v", err)
	}

	if len(outputFiles) != 0 {
		t.Errorf("出力ファイルが生成されるべきではありませんが、%d個のファイルが見つかりました", len(outputFiles))
	}
}
