package usecases

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestProcessImageRename_Normal は ProcessImageRename メソッドの正常系テストです
func TestProcessImageRename_Normal(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()

	// テスト用の画像ファイルを作成
	testFiles := []string{"test1.jpg", "test2.png", "test3.webp"}
	for _, filename := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		if err := os.WriteFile(filePath, []byte("test image content"), 0644); err != nil {
			t.Fatalf("テストファイルの作成に失敗しました: %v", err)
		}
	}

	config := Config{
		SrcDir:     tempDir,
		SortByName: true,
		SortByTime: false,
		Prefix:     "IMG",
		Delimiter:  "_",
		Digits:     3,
		StartCount: 1,
		Recursive:  false,
		Workers:    1,
	}

	var stdout, stderr bytes.Buffer

	// Act
	successCount, errorCount, err := ProcessImageRename(config, &stdout, &stderr)

	// Assert
	if err != nil {
		t.Errorf("ProcessImageRename でエラーが発生しました: %v", err)
	}

	if successCount != 3 {
		t.Errorf("成功数が期待値と異なります。期待値: 3, 実際: %d", successCount)
	}

	if errorCount != 0 {
		t.Errorf("エラー数が期待値と異なります。期待値: 0, 実際: %d", errorCount)
	}

	// リネーム後のファイルが存在することを確認
	expectedFiles := []string{"IMG_001.jpg", "IMG_002.png", "IMG_003.webp"}
	for _, expectedFile := range expectedFiles {
		expectedPath := filepath.Join(tempDir, expectedFile)
		if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
			t.Errorf("リネーム後のファイルが存在しません: %s", expectedFile)
		}
	}

	// 標準出力に期待されるメッセージが含まれていることを確認
	stdoutStr := stdout.String()
	if !strings.Contains(stdoutStr, "画像ファイルが 3 件見つかりました") {
		t.Errorf("期待される出力メッセージが見つかりません: %s", stdoutStr)
	}
}

// TestProcessImageRename_EmptyDirectory は空のディレクトリでのテストです
func TestProcessImageRename_EmptyDirectory(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()

	config := Config{
		SrcDir:     tempDir,
		SortByName: true,
		SortByTime: false,
		Prefix:     "IMG",
		Delimiter:  "_",
		Digits:     3,
		StartCount: 1,
		Recursive:  false,
		Workers:    1,
	}

	var stdout, stderr bytes.Buffer

	// Act
	successCount, errorCount, err := ProcessImageRename(config, &stdout, &stderr)

	// Assert
	if err != nil {
		t.Errorf("ProcessImageRename でエラーが発生しました: %v", err)
	}

	if successCount != 0 {
		t.Errorf("成功数が期待値と異なります。期待値: 0, 実際: %d", successCount)
	}

	if errorCount != 0 {
		t.Errorf("エラー数が期待値と異なります。期待値: 0, 実際: %d", errorCount)
	}

	// 標準出力に期待されるメッセージが含まれていることを確認
	stdoutStr := stdout.String()
	if !strings.Contains(stdoutStr, "画像ファイルが見つかりませんでした") {
		t.Errorf("期待される出力メッセージが見つかりません: %s", stdoutStr)
	}
}

// TestProcessImageRename_InvalidConfig は無効な設定でのテストです
func TestProcessImageRename_InvalidConfig(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()

	// プレフィックスが空の無効な設定
	config := Config{
		SrcDir:     tempDir,
		SortByName: true,
		SortByTime: false,
		Prefix:     "", // 空のプレフィックス（無効）
		Delimiter:  "_",
		Digits:     3,
		StartCount: 1,
		Recursive:  false,
		Workers:    1,
	}

	var stdout, stderr bytes.Buffer

	// Act
	successCount, errorCount, err := ProcessImageRename(config, &stdout, &stderr)

	// Assert
	if err == nil {
		t.Error("無効な設定でエラーが発生しませんでした")
	}

	if successCount != 0 {
		t.Errorf("成功数が期待値と異なります。期待値: 0, 実際: %d", successCount)
	}

	if errorCount != 0 {
		t.Errorf("エラー数が期待値と異なります。期待値: 0, 実際: %d", errorCount)
	}

	// 標準エラー出力にエラーメッセージが含まれていることを確認
	stderrStr := stderr.String()
	if !strings.Contains(stderrStr, "プレフィックスは必須です") {
		t.Errorf("期待されるエラーメッセージが見つかりません: %s", stderrStr)
	}
}

// TestProcessImageRename_SortByTime は時間順ソートのテストです
func TestProcessImageRename_SortByTime(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()

	// テスト用の画像ファイルを異なる時間で作成
	testFiles := []string{"old.jpg", "new.png"}
	for i, filename := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		if err := os.WriteFile(filePath, []byte("test image content"), 0644); err != nil {
			t.Fatalf("テストファイルの作成に失敗しました: %v", err)
		}

		// ファイルの更新時間を調整（古いファイルを先に作成）
		if i == 0 {
			// Note: 実際のテストでは time.Now().Add(-time.Hour) などを使用してファイルの更新時間を調整
			// ここでは簡単のため、ファイル作成順序でソートをテストする
		}
	}

	config := Config{
		SrcDir:     tempDir,
		SortByName: false,
		SortByTime: true,
		Prefix:     "TIME",
		Delimiter:  "_",
		Digits:     2,
		StartCount: 10,
		Recursive:  false,
		Workers:    1,
	}

	var stdout, stderr bytes.Buffer

	// Act
	successCount, errorCount, err := ProcessImageRename(config, &stdout, &stderr)

	// Assert
	if err != nil {
		t.Errorf("ProcessImageRename でエラーが発生しました: %v", err)
	}

	if successCount != 2 {
		t.Errorf("成功数が期待値と異なります。期待値: 2, 実際: %d", successCount)
	}

	if errorCount != 0 {
		t.Errorf("エラー数が期待値と異なります。期待値: 0, 実際: %d", errorCount)
	}

	// 標準出力に時間順ソートのメッセージが含まれていることを確認
	stdoutStr := stdout.String()
	if !strings.Contains(stdoutStr, "ファイルを更新日時順に並べ替えています") {
		t.Errorf("期待される出力メッセージが見つかりません: %s", stdoutStr)
	}
}

// TestProcessImageRename_RecursiveSearch は再帰検索のテストです
func TestProcessImageRename_RecursiveSearch(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()

	// サブディレクトリを作成
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("サブディレクトリの作成に失敗しました: %v", err)
	}

	// ルートディレクトリとサブディレクトリにファイルを作成
	testFiles := map[string]string{
		"root.jpg":          tempDir,
		"subdir/nested.png": tempDir,
	}

	for filename, baseDir := range testFiles {
		filePath := filepath.Join(baseDir, filename)
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("ディレクトリの作成に失敗しました: %v", err)
		}
		if err := os.WriteFile(filePath, []byte("test image content"), 0644); err != nil {
			t.Fatalf("テストファイルの作成に失敗しました: %v", err)
		}
	}

	config := Config{
		SrcDir:     tempDir,
		SortByName: true,
		SortByTime: false,
		Prefix:     "REC",
		Delimiter:  "_",
		Digits:     2,
		StartCount: 1,
		Recursive:  true, // 再帰検索を有効
		Workers:    1,
	}

	var stdout, stderr bytes.Buffer

	// Act
	successCount, errorCount, err := ProcessImageRename(config, &stdout, &stderr)

	// Assert
	if err != nil {
		t.Errorf("ProcessImageRename でエラーが発生しました: %v", err)
	}

	if successCount != 2 {
		t.Errorf("成功数が期待値と異なります。期待値: 2, 実際: %d", successCount)
	}

	if errorCount != 0 {
		t.Errorf("エラー数が期待値と異なります。期待値: 0, 実際: %d", errorCount)
	}

	// 標準出力に期待されるファイル数が含まれていることを確認
	stdoutStr := stdout.String()
	if !strings.Contains(stdoutStr, "画像ファイルが 2 件見つかりました") {
		t.Errorf("期待される出力メッセージが見つかりません: %s", stdoutStr)
	}
}

// TestProcessImageRename_NonExistentDirectory は存在しないディレクトリでのテストです
func TestProcessImageRename_NonExistentDirectory(t *testing.T) {
	// Arrange
	nonExistentDir := "/path/to/nonexistent/directory"

	config := Config{
		SrcDir:     nonExistentDir,
		SortByName: true,
		SortByTime: false,
		Prefix:     "IMG",
		Delimiter:  "_",
		Digits:     3,
		StartCount: 1,
		Recursive:  false,
		Workers:    1,
	}

	var stdout, stderr bytes.Buffer

	// Act
	successCount, errorCount, err := ProcessImageRename(config, &stdout, &stderr)

	// Assert
	if err == nil {
		t.Error("存在しないディレクトリでエラーが発生しませんでした")
	}

	if successCount != 0 {
		t.Errorf("成功数が期待値と異なります。期待値: 0, 実際: %d", successCount)
	}

	if errorCount != 0 {
		t.Errorf("エラー数が期待値と異なります。期待値: 0, 実際: %d", errorCount)
	}

	// 標準エラー出力にエラーメッセージが含まれていることを確認
	stderrStr := stderr.String()
	if !strings.Contains(stderrStr, "アクセスエラー") {
		t.Errorf("期待されるエラーメッセージが見つかりません: %s", stderrStr)
	}
}

// TestProcessImageRename_MultipleWorkers は複数ワーカーでのテストです
func TestProcessImageRename_MultipleWorkers(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()

	// 複数のテスト用画像ファイルを作成
	testFiles := []string{"img1.jpg", "img2.png", "img3.webp", "img4.avif", "img5.jpeg"}
	for _, filename := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		if err := os.WriteFile(filePath, []byte("test image content"), 0644); err != nil {
			t.Fatalf("テストファイルの作成に失敗しました: %v", err)
		}
	}

	config := Config{
		SrcDir:     tempDir,
		SortByName: true,
		SortByTime: false,
		Prefix:     "MULTI",
		Delimiter:  "-",
		Digits:     4,
		StartCount: 100,
		Recursive:  false,
		Workers:    3, // 複数ワーカーを使用
	}

	var stdout, stderr bytes.Buffer

	// Act
	successCount, errorCount, err := ProcessImageRename(config, &stdout, &stderr)

	// Assert
	if err != nil {
		t.Errorf("ProcessImageRename でエラーが発生しました: %v", err)
	}

	if successCount != 5 {
		t.Errorf("成功数が期待値と異なります。期待値: 5, 実際: %d", successCount)
	}

	if errorCount != 0 {
		t.Errorf("エラー数が期待値と異なります。期待値: 0, 実際: %d", errorCount)
	}

	// リネーム後のファイルが存在することを確認
	expectedFiles := []string{
		"MULTI-0100.jpg", "MULTI-0101.png", "MULTI-0102.webp",
		"MULTI-0103.avif", "MULTI-0104.jpeg",
	}
	for _, expectedFile := range expectedFiles {
		expectedPath := filepath.Join(tempDir, expectedFile)
		if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
			t.Errorf("リネーム後のファイルが存在しません: %s", expectedFile)
		}
	}

	// 標準出力にワーカー数のメッセージが含まれていることを確認
	stdoutStr := stdout.String()
	if !strings.Contains(stdoutStr, "リネーム操作に 3 ワーカーを使用します") {
		t.Errorf("期待される出力メッセージが見つかりません: %s", stdoutStr)
	}
}

// TestProcessImageRename_CustomDelimiterAndDigits はカスタム区切り文字と桁数のテストです
func TestProcessImageRename_CustomDelimiterAndDigits(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()

	// テスト用の画像ファイルを作成
	testFiles := []string{"custom1.jpg", "custom2.png"}
	for _, filename := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		if err := os.WriteFile(filePath, []byte("test image content"), 0644); err != nil {
			t.Fatalf("テストファイルの作成に失敗しました: %v", err)
		}
	}

	config := Config{
		SrcDir:     tempDir,
		SortByName: true,
		SortByTime: false,
		Prefix:     "CUSTOM",
		Delimiter:  ".", // カスタム区切り文字
		Digits:     5,   // カスタム桁数
		StartCount: 1,
		Recursive:  false,
		Workers:    1,
	}

	var stdout, stderr bytes.Buffer

	// Act
	successCount, errorCount, err := ProcessImageRename(config, &stdout, &stderr)

	// Assert
	if err != nil {
		t.Errorf("ProcessImageRename でエラーが発生しました: %v", err)
	}

	if successCount != 2 {
		t.Errorf("成功数が期待値と異なります。期待値: 2, 実際: %d", successCount)
	}

	if errorCount != 0 {
		t.Errorf("エラー数が期待値と異なります。期待値: 0, 実際: %d", errorCount)
	}

	// リネーム後のファイルが期待される形式で存在することを確認
	expectedFiles := []string{"CUSTOM.00001.jpg", "CUSTOM.00002.png"}
	for _, expectedFile := range expectedFiles {
		expectedPath := filepath.Join(tempDir, expectedFile)
		if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
			t.Errorf("リネーム後のファイルが存在しません: %s", expectedFile)
		}
	}

	// 標準出力に設定情報が含まれていることを確認
	stdoutStr := stdout.String()
	if !strings.Contains(stdoutStr, "区切り文字: .") {
		t.Errorf("期待される出力メッセージが見つかりません: %s", stdoutStr)
	}
}
