package usecases

import (
	"fmt"
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
	expectedCounts := 6
	if !strings.Contains(result, fmt.Sprintf("%d件のファイルからコンテンツを抽出しました", expectedCounts)) {
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

	if len(outputFiles) != expectedCounts {
		t.Errorf("期待される出力ファイル数は%dですが、実際は%d個でした", expectedCounts, len(outputFiles))
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

// TestService_ExtractBlogContent_WithSpecificPythonFile_Invalid は 日本語を使ったファイル名を取得することが出来ないことを確認するテスト
func TestService_ExtractBlogContent_WithSpecificPythonFile_Invalid(t *testing.T) {

	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "notion-blog-extractor-python-file-test")
	if err != nil {
		t.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// filepath.Globを使って実際のPythonファイルを見つける
	globPattern := "./test_data/org/Python*"
	matches, err := filepath.Glob(globPattern)
	if err != nil {
		t.Fatalf("filepath.Glob()でエラーが発生しました: %v", err)
	}

	if len(matches) == 0 {
		t.Fatalf("Pythonファイルが見つかりません: %s", globPattern)
	}

	// 最初に見つかったPythonファイルを使用
	originalFilePath := matches[0]
	t.Logf("使用するファイル: %s", originalFilePath)

	// ファイルの存在確認
	if _, err := os.Stat(originalFilePath); os.IsNotExist(err) {
		t.Fatalf("指定されたファイルが存在しません: %s", originalFilePath)
	}

	// 指定されたファイルを読み込み
	content, err := os.ReadFile(originalFilePath)
	if err != nil {
		t.Fatalf("ファイルの読み込みに失敗しました %s: %v", originalFilePath, err)
	}

	// ソース一時ディレクトリを作成
	srcTempDir := filepath.Join(tempDir, "src")
	if err := os.MkdirAll(srcTempDir, 0755); err != nil {
		t.Fatalf("ソース一時ディレクトリの作成に失敗しました: %v", err)
	}

	// 一時ディレクトリにファイルをコピー
	fileName := filepath.Base(originalFilePath)
	destFilePath := filepath.Join(srcTempDir, fileName)
	if err := os.WriteFile(destFilePath, content, 0644); err != nil {
		t.Fatalf("ファイルのコピーに失敗しました %s: %v", fileName, err)
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

	// 抽出されたファイル数の確認（1件のファイルを処理）
	if !strings.Contains(result, "1件のファイルからコンテンツを抽出しました") {
		t.Errorf("期待される抽出ファイル数が含まれていません。結果: %s", result)
	}

	// 出力ディレクトリの存在確認
	if _, err := os.Stat(destTempDir); os.IsNotExist(err) {
		t.Error("出力ディレクトリが作成されていません")
		return
	}

	// 出力ファイルの確認
	outputFilePath := filepath.Join(destTempDir, fileName)
	if _, err := os.Stat(outputFilePath); os.IsNotExist(err) {
		t.Error("出力ファイルが生成されていません")
		return
	}

	// 出力ファイルの内容を読み込んで検証
	outputContent, err := os.ReadFile(outputFilePath)
	if err != nil {
		t.Errorf("出力ファイルの読み込みに失敗しました %s: %v", outputFilePath, err)
		return
	}

	contentStr := string(outputContent)

	// コンテンツマーカーで始まることを確認
	if !strings.HasPrefix(contentStr, "# Content") {
		t.Errorf("ファイル %s: コンテンツマーカーで始まっていません", fileName)
	}

	// "## はじまり" が含まれることを確認
	if !strings.Contains(contentStr, "## はじまり") {
		t.Errorf("ファイル %s: '## はじまり' マーカーが含まれていません", fileName)
	}

	// 元のメタデータ部分が除外されていることを確認（priority, status, tagsなど）
	if strings.Contains(contentStr, "priority:") || strings.Contains(contentStr, "status:") {
		t.Errorf("ファイル %s: メタデータ部分が除外されていません", fileName)
	}

	// ファイル名に期待される文字が含まれていることを確認（特殊文字処理のテスト）
	if !strings.Contains(fileName, "Python") {
		t.Errorf("ファイル名に期待される文字が含まれていません: %s", fileName)
	}

	// 実際のファイル名に含まれる文字列を確認（日本語文字は1文字程度なら読み取れることを確認）
	if !strings.Contains(fileName, "の") {
		t.Errorf("ファイル名に期待される日本語文字が含まれていません: %s", fileName)
	}

	expectedChar := "タブサイズ"
	if strings.Contains(fileName, expectedChar) {
		t.Errorf("ファイル名に期待される「%s」が含まれないはずです（おそらくUnicodeのバイト数によるバグ）: %s", expectedChar, fileName)
	}
	expectedChar = "のタブサイズがどうしても"
	if strings.Contains(fileName, expectedChar) {
		t.Errorf("ファイル名に期待される「%s」が含まれないはずです（おそらくUnicodeのバイト数によるバグ）: %s", expectedChar, fileName)
	}

	// 長いファイル名の処理確認（cが連続している部分）
	if !strings.Contains(fileName, "ccccccccccccccccccccccccccccccc") {
		t.Errorf("ファイル名に期待される長い文字列が含まれていません: %s", fileName)
	}
}

func TestService_ExtractBlogContent_WithSpecificYouTubeFile(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "notion-blog-extractor-youtube-file-test")
	if err != nil {
		t.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 指定されたファイルのパス
	originalFilePath := "./test_data/org/【YouTubeなど】動画の更新状況_2025年01月_20250111_test_22222222222222222222222222222.md"

	// ファイルの存在確認
	if _, err := os.Stat(originalFilePath); os.IsNotExist(err) {
		t.Fatalf("指定されたファイルが存在しません: %s", originalFilePath)
	}

	// 指定されたファイルを読み込み
	content, err := os.ReadFile(originalFilePath)
	if err != nil {
		t.Fatalf("ファイルの読み込みに失敗しました %s: %v", originalFilePath, err)
	}

	// ソース一時ディレクトリを作成
	srcTempDir := filepath.Join(tempDir, "src")
	if err := os.MkdirAll(srcTempDir, 0755); err != nil {
		t.Fatalf("ソース一時ディレクトリの作成に失敗しました: %v", err)
	}

	// 一時ディレクトリにファイルをコピー
	fileName := filepath.Base(originalFilePath)
	destFilePath := filepath.Join(srcTempDir, fileName)
	if err := os.WriteFile(destFilePath, content, 0644); err != nil {
		t.Fatalf("ファイルのコピーに失敗しました %s: %v", fileName, err)
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

	// 抽出されたファイル数の確認（1件のファイルを処理）
	if !strings.Contains(result, "1件のファイルからコンテンツを抽出しました") {
		t.Errorf("期待される抽出ファイル数が含まれていません。結果: %s", result)
	}

	// 出力ディレクトリの存在確認
	if _, err := os.Stat(destTempDir); os.IsNotExist(err) {
		t.Error("出力ディレクトリが作成されていません")
		return
	}

	// 出力ファイルの確認
	outputFilePath := filepath.Join(destTempDir, fileName)
	if _, err := os.Stat(outputFilePath); os.IsNotExist(err) {
		t.Error("出力ファイルが生成されていません")
		return
	}

	// 出力ファイルの内容を読み込んで検証
	outputContent, err := os.ReadFile(outputFilePath)
	if err != nil {
		t.Errorf("出力ファイルの読み込みに失敗しました %s: %v", outputFilePath, err)
		return
	}

	contentStr := string(outputContent)

	// コンテンツマーカーで始まることを確認
	if !strings.HasPrefix(contentStr, "# Content") {
		t.Errorf("ファイル %s: コンテンツマーカーで始まっていません", fileName)
	}

	// "## はじまり" が含まれることを確認
	if !strings.Contains(contentStr, "## はじまり") {
		t.Errorf("ファイル %s: '## はじまり' マーカーが含まれていません", fileName)
	}

	// 元のメタデータ部分が除外されていることを確認（priority, status, tagsなど）
	if strings.Contains(contentStr, "priority:") || strings.Contains(contentStr, "status:") {
		t.Errorf("ファイル %s: メタデータ部分が除外されていません", fileName)
	}

	// ファイル名に日本語が含まれていることを確認（特殊文字処理のテスト）
	if !strings.Contains(fileName, "YouTube") || !strings.Contains(fileName, "動画の更新状況") {
		t.Errorf("ファイル名に期待される日本語文字が含まれていません: %s", fileName)
	}

	// 長いファイル名の処理確認（2が連続している部分）
	if !strings.Contains(fileName, "22222222222222222222222222222") {
		t.Errorf("ファイル名に期待される長い文字列が含まれていません: %s", fileName)
	}

	// 特殊文字（【】）の処理確認
	if !strings.Contains(fileName, "【") || !strings.Contains(fileName, "】") {
		t.Errorf("ファイル名に期待される特殊文字が含まれていません: %s", fileName)
	}
}

func TestService_ExtractBlogContent_FileAccessIssue(t *testing.T) {
	// 指定されたファイルのパス（問題のあるファイル名）
	originalFilePath := "./test_data/org/Pythonのタブサイズがどうしても「4」にならない_20241215_test ccccccccccccccccccccccccccccccc.md"

	// os.Statでファイルが取得できないことを検証
	_, err := os.Stat(originalFilePath)
	if err == nil {
		t.Errorf("os.Stat()でファイルが見つかりました。期待されるエラーが発生しませんでした: %s", originalFilePath)
	} else {
		t.Logf("期待通りos.Stat()でエラーが発生しました: %v", err)
	}

	// os.ReadFileでも試してみる
	_, err = os.ReadFile(originalFilePath)
	if err == nil {
		t.Errorf("os.ReadFile()でファイルが読み込めました。期待されるエラーが発生しませんでした: %s", originalFilePath)
	} else {
		t.Logf("期待通りos.ReadFile()でエラーが発生しました: %v", err)
	}

	// filepath.Globを使って実際のファイル名を取得してみる
	globPattern := "./test_data/org/Python*"
	matches, err := filepath.Glob(globPattern)
	if err != nil {
		t.Errorf("filepath.Glob()でエラーが発生しました: %v", err)
	} else {
		t.Logf("filepath.Glob()で見つかったファイル: %v", matches)

		// 実際のファイルが存在する場合、そのファイルでテストを実行
		if len(matches) > 0 {
			for _, match := range matches {
				if strings.Contains(match, "Python") {
					t.Logf("実際のPythonファイル: %s", match)

					// 実際のファイルでos.Statを試す
					_, err := os.Stat(match)
					if err != nil {
						t.Errorf("実際のファイルでもos.Stat()エラー: %v", err)
					} else {
						t.Logf("実際のファイルはos.Stat()で正常に取得できました: %s", match)
					}

					// 実際のファイルでos.ReadFileを試す
					content, err := os.ReadFile(match)
					if err != nil {
						t.Errorf("実際のファイルでもos.ReadFile()エラー: %v", err)
					} else {
						t.Logf("実際のファイルはos.ReadFile()で正常に読み込めました。サイズ: %d bytes", len(content))
					}
				}
			}
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
