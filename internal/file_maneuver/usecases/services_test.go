package usecases

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// #==============================================================#
// ##          Tests for FileManeuverService                     ##
// #==============================================================#
// TestFileManeuverService はFileManeuverServiceのテストクラスです
type TestFileManeuverService struct{}

// TestFileManeuverService_FindTargetFiles は対象ファイル検索をテストします
func (tc *TestFileManeuverService) TestFileManeuverService_FindTargetFiles(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	destDir := filepath.Join(tempDir, "dest")
	os.MkdirAll(srcDir, 0755)
	os.MkdirAll(destDir, 0755)

	// テストファイルを作成
	testFiles := []string{"test1.jpg", "test2.png", "test3.txt", "test4.JPG"}
	for _, file := range testFiles {
		filePath := filepath.Join(srcDir, file)
		os.WriteFile(filePath, []byte("test content"), 0644)
	}

	config, err := NewConfig([]string{srcDir}, []string{"jpg", "png"}, "", destDir, false, 1, false, false, false)
	if err != nil {
		t.Fatalf("設定作成に失敗しました: %v", err)
	}

	service := NewFileManeuverService(config)
	var stdout, stderr bytes.Buffer

	// Act
	files, err := service.FindTargetFiles(&stdout, &stderr)

	// Assert
	if err != nil {
		t.Fatalf("ファイル検索に失敗しました: %v", err)
	}

	// jpg, png, JPGファイルが見つかることを確認（txtは除外）
	expectedCount := 3
	if len(files) != expectedCount {
		t.Errorf("見つかったファイル数が期待値と異なります。期待値: %d, 実際: %d", expectedCount, len(files))
	}

	// ファイル名を確認
	foundFiles := make(map[string]bool)
	for _, file := range files {
		foundFiles[filepath.Base(file)] = true
	}

	if !foundFiles["test1.jpg"] || !foundFiles["test2.png"] || !foundFiles["test4.JPG"] {
		t.Errorf("期待されるファイルが見つかりませんでした。見つかったファイル: %v", files)
	}

	if foundFiles["test3.txt"] {
		t.Error("除外されるべきファイルが見つかりました: test3.txt")
	}
}

func TestFileManeuverService_FindTargetFiles(t *testing.T) {
	tc := &TestFileManeuverService{}
	tc.TestFileManeuverService_FindTargetFiles(t)
}

// TestFileManeuverService_FindTargetFiles_ByNameContains はファイル名部分一致のみで対象を検索できることをテストします
func (tc *TestFileManeuverService) TestFileManeuverService_FindTargetFiles_ByNameContains(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	destDir := filepath.Join(tempDir, "dest")
	os.MkdirAll(srcDir, 0755)
	os.MkdirAll(destDir, 0755)

	testFiles := []string{"report-alpha.jpg", "monthly_report.txt", "notes.doc"}
	for _, file := range testFiles {
		os.WriteFile(filepath.Join(srcDir, file), []byte("test"), 0644)
	}

	config, err := NewConfig([]string{srcDir}, []string{}, "REPORT", destDir, false, 1, false, false, false)
	if err != nil {
		t.Fatalf("設定作成に失敗しました: %v", err)
	}

	service := NewFileManeuverService(config)
	var stdout, stderr bytes.Buffer

	// Act
	files, err := service.FindTargetFiles(&stdout, &stderr)

	// Assert
	if err != nil {
		t.Fatalf("ファイル検索に失敗しました: %v", err)
	}

	if len(files) != 2 {
		t.Fatalf("ファイル名フィルターのみで想定外のファイル数です。期待値: 2, 実際: %d", len(files))
	}

	found := map[string]bool{}
	for _, file := range files {
		found[filepath.Base(file)] = true
	}

	if !found["report-alpha.jpg"] || !found["monthly_report.txt"] {
		t.Errorf("期待されるファイルが見つかりませんでした。found=%v", found)
	}

	if found["notes.doc"] {
		t.Error("除外されるべきファイルが検索結果に含まれています: notes.doc")
	}
}

func TestFileManeuverService_FindTargetFiles_ByNameContains(t *testing.T) {
	tc := &TestFileManeuverService{}
	tc.TestFileManeuverService_FindTargetFiles_ByNameContains(t)
}

// TestFileManeuverService_FindTargetFiles_NameContainsWithExtensions は拡張子とファイル名フィルターを同時に適用できることをテストします
func (tc *TestFileManeuverService) TestFileManeuverService_FindTargetFiles_NameContainsWithExtensions(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	destDir := filepath.Join(tempDir, "dest")
	os.MkdirAll(srcDir, 0755)
	os.MkdirAll(destDir, 0755)

	testFiles := []string{"report.jpg", "report.txt", "photo.jpg"}
	for _, file := range testFiles {
		os.WriteFile(filepath.Join(srcDir, file), []byte("test"), 0644)
	}

	config, err := NewConfig([]string{srcDir}, []string{"jpg"}, "report", destDir, false, 1, false, false, false)
	if err != nil {
		t.Fatalf("設定作成に失敗しました: %v", err)
	}

	service := NewFileManeuverService(config)
	var stdout, stderr bytes.Buffer

	// Act
	files, err := service.FindTargetFiles(&stdout, &stderr)

	// Assert
	if err != nil {
		t.Fatalf("ファイル検索に失敗しました: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("フィルターの掛け合わせ結果が期待値と異なります。期待値: 1, 実際: %d", len(files))
	}

	if filepath.Base(files[0]) != "report.jpg" {
		t.Errorf("期待されるファイルが見つかりませんでした。found=%v", files)
	}
}

func TestFileManeuverService_FindTargetFiles_NameContainsWithExtensions(t *testing.T) {
	tc := &TestFileManeuverService{}
	tc.TestFileManeuverService_FindTargetFiles_NameContainsWithExtensions(t)
}

// TestFileManeuverService_FindTargetFilesRecursive は再帰的ファイル検索をテストします
func (tc *TestFileManeuverService) TestFileManeuverService_FindTargetFilesRecursive(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	subDir := filepath.Join(srcDir, "subdir")
	destDir := filepath.Join(tempDir, "dest")
	os.MkdirAll(subDir, 0755)
	os.MkdirAll(destDir, 0755)

	// ルートディレクトリにファイル作成
	os.WriteFile(filepath.Join(srcDir, "root.jpg"), []byte("test"), 0644)
	// サブディレクトリにファイル作成
	os.WriteFile(filepath.Join(subDir, "sub.jpg"), []byte("test"), 0644)

	config, err := NewConfig([]string{srcDir}, []string{"jpg"}, "", destDir, true, 1, false, false, false)
	if err != nil {
		t.Fatalf("設定作成に失敗しました: %v", err)
	}

	service := NewFileManeuverService(config)
	var stdout, stderr bytes.Buffer

	// Act
	files, err := service.FindTargetFiles(&stdout, &stderr)

	// Assert
	if err != nil {
		t.Fatalf("ファイル検索に失敗しました: %v", err)
	}

	// 再帰的検索で両方のファイルが見つかることを確認
	expectedCount := 2
	if len(files) != expectedCount {
		t.Errorf("見つかったファイル数が期待値と異なります。期待値: %d, 実際: %d", expectedCount, len(files))
	}
}

func TestFileManeuverService_FindTargetFilesRecursive(t *testing.T) {
	tc := &TestFileManeuverService{}
	tc.TestFileManeuverService_FindTargetFilesRecursive(t)
}

// TestFileManeuverService_DryRun はドライランモードをテストします
func (tc *TestFileManeuverService) TestFileManeuverService_DryRun(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	destDir := filepath.Join(tempDir, "dest")
	os.MkdirAll(srcDir, 0755)
	os.MkdirAll(destDir, 0755)

	// テストファイルを作成
	testFile := filepath.Join(srcDir, "test.jpg")
	os.WriteFile(testFile, []byte("test content"), 0644)

	config, err := NewConfig([]string{srcDir}, []string{"jpg"}, "", destDir, false, 1, true, false, false) // ドライランモード
	if err != nil {
		t.Fatalf("設定作成に失敗しました: %v", err)
	}

	service := NewFileManeuverService(config)
	var stdout, stderr bytes.Buffer

	// Act
	successCount, errorCount, err := service.ExecuteFileManeuver(&stdout, &stderr)

	// Assert
	if err != nil {
		t.Fatalf("ファイル移動処理に失敗しました: %v", err)
	}

	if successCount != 1 {
		t.Errorf("成功カウントが期待値と異なります。期待値: 1, 実際: %d", successCount)
	}

	if errorCount != 0 {
		t.Errorf("エラーカウントが期待値と異なります。期待値: 0, 実際: %d", errorCount)
	}

	// ドライランなので元ファイルが残っていることを確認
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("ドライランモードで元ファイルが削除されました")
	}

	// 宛先にファイルが移動されていないことを確認
	destFile := filepath.Join(destDir, "test.jpg")
	if _, err := os.Stat(destFile); !os.IsNotExist(err) {
		t.Error("ドライランモードで宛先にファイルが作成されました")
	}

	// 出力にドライランメッセージが含まれることを確認
	output := stdout.String()
	if !strings.Contains(output, "ドライランモード") {
		t.Error("ドライランメッセージが出力されませんでした")
	}
}

func TestFileManeuverService_DryRun(t *testing.T) {
	tc := &TestFileManeuverService{}
	tc.TestFileManeuverService_DryRun(t)
}

// TestFileManeuverService_FileConflict はファイル衝突処理をテストします
func (tc *TestFileManeuverService) TestFileManeuverService_FileConflict(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	destDir := filepath.Join(tempDir, "dest")
	os.MkdirAll(srcDir, 0755)
	os.MkdirAll(destDir, 0755)

	// ソースファイルを作成
	srcFile := filepath.Join(srcDir, "test.jpg")
	os.WriteFile(srcFile, []byte("source content"), 0644)

	// 宛先に同名ファイルを作成（衝突状況）
	destFile := filepath.Join(destDir, "test.jpg")
	os.WriteFile(destFile, []byte("existing content"), 0644)

	config, err := NewConfig([]string{srcDir}, []string{"jpg"}, "", destDir, false, 1, false, false, false)
	if err != nil {
		t.Fatalf("設定作成に失敗しました: %v", err)
	}

	service := NewFileManeuverService(config)
	var stdout, stderr bytes.Buffer

	// Act
	successCount, errorCount, err := service.ExecuteFileManeuver(&stdout, &stderr)

	// Assert
	if err != nil {
		t.Fatalf("ファイル移動処理に失敗しました: %v", err)
	}

	// 衝突によりスキップされるため成功カウントは0
	if successCount != 0 {
		t.Errorf("成功カウントが期待値と異なります。期待値: 0, 実際: %d", successCount)
	}

	if errorCount != 0 {
		t.Errorf("エラーカウントが期待値と異なります。期待値: 0, 実際: %d", errorCount)
	}

	// 元ファイルが残っていることを確認
	if _, err := os.Stat(srcFile); os.IsNotExist(err) {
		t.Error("衝突時に元ファイルが削除されました")
	}

	// 宛先ファイルが元の内容のままであることを確認
	content, err := os.ReadFile(destFile)
	if err != nil {
		t.Fatalf("宛先ファイルの読み取りに失敗しました: %v", err)
	}

	if string(content) != "existing content" {
		t.Error("宛先ファイルの内容が変更されました")
	}

	// 警告メッセージが出力されることを確認
	stderrOutput := stderr.String()
	if !strings.Contains(stderrOutput, "警告") {
		t.Error("衝突警告メッセージが出力されませんでした")
	}
}

func TestFileManeuverService_FileConflict(t *testing.T) {
	tc := &TestFileManeuverService{}
	tc.TestFileManeuverService_FileConflict(t)
}

// TestFileManeuverService_ExecuteFileManeuver は統合テストを実行します
func (tc *TestFileManeuverService) TestFileManeuverService_ExecuteFileManeuver(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	destDir := filepath.Join(tempDir, "dest")
	os.MkdirAll(srcDir, 0755)
	os.MkdirAll(destDir, 0755)

	// テストファイルを作成
	testFiles := []string{"file1.jpg", "file2.png"}
	for _, file := range testFiles {
		filePath := filepath.Join(srcDir, file)
		os.WriteFile(filePath, []byte("test content"), 0644)
	}

	config, err := NewConfig([]string{srcDir}, []string{"jpg", "png"}, "", destDir, false, 1, false, false, false)
	if err != nil {
		t.Fatalf("設定作成に失敗しました: %v", err)
	}

	service := NewFileManeuverService(config)
	var stdout, stderr bytes.Buffer

	// Act
	successCount, errorCount, err := service.ExecuteFileManeuver(&stdout, &stderr)

	// Assert
	if err != nil {
		t.Fatalf("ファイル移動処理に失敗しました: %v", err)
	}

	if successCount != 2 {
		t.Errorf("成功カウントが期待値と異なります。期待値: 2, 実際: %d", successCount)
	}

	if errorCount != 0 {
		t.Errorf("エラーカウントが期待値と異なります。期待値: 0, 実際: %d", errorCount)
	}

	// ファイルが正しく移動されたことを確認
	for _, file := range testFiles {
		srcPath := filepath.Join(srcDir, file)
		destPath := filepath.Join(destDir, file)

		// 元ファイルが削除されていることを確認
		if _, err := os.Stat(srcPath); !os.IsNotExist(err) {
			t.Errorf("元ファイルが削除されていません: %s", srcPath)
		}

		// 宛先ファイルが存在することを確認
		if _, err := os.Stat(destPath); os.IsNotExist(err) {
			t.Errorf("宛先ファイルが作成されていません: %s", destPath)
		}
	}
}

func TestFileManeuverService_ExecuteFileManeuver(t *testing.T) {
	tc := &TestFileManeuverService{}
	tc.TestFileManeuverService_ExecuteFileManeuver(t)
}

// TestFileManeuverService_CopyMode はコピーモードをテストします
func (tc *TestFileManeuverService) TestFileManeuverService_CopyMode(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	destDir := filepath.Join(tempDir, "dest")
	os.MkdirAll(srcDir, 0755)
	os.MkdirAll(destDir, 0755)

	// テストファイルを作成
	testFile := filepath.Join(srcDir, "test.jpg")
	testContent := "test content for copy"
	os.WriteFile(testFile, []byte(testContent), 0644)

	config, err := NewConfig([]string{srcDir}, []string{"jpg"}, "", destDir, false, 1, false, true, false) // コピーモード
	if err != nil {
		t.Fatalf("設定作成に失敗しました: %v", err)
	}

	service := NewFileManeuverService(config)
	var stdout, stderr bytes.Buffer

	// Act
	successCount, errorCount, err := service.ExecuteFileManeuver(&stdout, &stderr)

	// Assert
	if err != nil {
		t.Fatalf("ファイルコピー処理に失敗しました: %v", err)
	}

	if successCount != 1 {
		t.Errorf("成功カウントが期待値と異なります。期待値: 1, 実際: %d", successCount)
	}

	if errorCount != 0 {
		t.Errorf("エラーカウントが期待値と異なります。期待値: 0, 実際: %d", errorCount)
	}

	// 元ファイルが残っていることを確認（コピーモード）
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("コピーモードで元ファイルが削除されました")
	}

	// 宛先ファイルが存在することを確認
	destFile := filepath.Join(destDir, "test.jpg")
	if _, err := os.Stat(destFile); os.IsNotExist(err) {
		t.Error("宛先ファイルが作成されていません")
	}

	// ファイル内容が正しくコピーされていることを確認
	copiedContent, err := os.ReadFile(destFile)
	if err != nil {
		t.Fatalf("コピーされたファイルの読み取りに失敗しました: %v", err)
	}

	if string(copiedContent) != testContent {
		t.Errorf("ファイル内容が正しくコピーされていません。期待値: %s, 実際: %s", testContent, string(copiedContent))
	}

	// 出力にコピーメッセージが含まれることを確認
	output := stdout.String()
	if !strings.Contains(output, "ファイルコピーを開始します") {
		t.Error("コピーメッセージが出力されませんでした")
	}
}

func TestFileManeuverService_CopyMode(t *testing.T) {
	tc := &TestFileManeuverService{}
	tc.TestFileManeuverService_CopyMode(t)
}

// TestFileManeuverService_OverwriteMode は上書きモードをテストします
func (tc *TestFileManeuverService) TestFileManeuverService_OverwriteMode(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	destDir := filepath.Join(tempDir, "dest")
	os.MkdirAll(srcDir, 0755)
	os.MkdirAll(destDir, 0755)

	// ソースファイルを作成
	srcFile := filepath.Join(srcDir, "test.jpg")
	srcContent := "new content"
	os.WriteFile(srcFile, []byte(srcContent), 0644)

	// 宛先に同名ファイルを作成（上書き対象）
	destFile := filepath.Join(destDir, "test.jpg")
	oldContent := "old content"
	os.WriteFile(destFile, []byte(oldContent), 0644)

	config, err := NewConfig([]string{srcDir}, []string{"jpg"}, "", destDir, false, 1, false, false, true) // 上書きモード
	if err != nil {
		t.Fatalf("設定作成に失敗しました: %v", err)
	}

	service := NewFileManeuverService(config)
	var stdout, stderr bytes.Buffer

	// Act
	successCount, errorCount, err := service.ExecuteFileManeuver(&stdout, &stderr)

	// Assert
	if err != nil {
		t.Fatalf("ファイル移動処理に失敗しました: %v", err)
	}

	if successCount != 1 {
		t.Errorf("成功カウントが期待値と異なります。期待値: 1, 実際: %d", successCount)
	}

	if errorCount != 0 {
		t.Errorf("エラーカウントが期待値と異なります。期待値: 0, 実際: %d", errorCount)
	}

	// 元ファイルが削除されていることを確認
	if _, err := os.Stat(srcFile); !os.IsNotExist(err) {
		t.Error("上書きモードで元ファイルが削除されていません")
	}

	// 宛先ファイルが新しい内容で上書きされていることを確認
	newContent, err := os.ReadFile(destFile)
	if err != nil {
		t.Fatalf("宛先ファイルの読み取りに失敗しました: %v", err)
	}

	if string(newContent) != srcContent {
		t.Errorf("ファイルが正しく上書きされていません。期待値: %s, 実際: %s", srcContent, string(newContent))
	}

	// 出力に上書きメッセージが含まれることを確認
	output := stdout.String()
	if !strings.Contains(output, "上書きモード") {
		t.Error("上書きメッセージが出力されませんでした")
	}
}

func TestFileManeuverService_OverwriteMode(t *testing.T) {
	tc := &TestFileManeuverService{}
	tc.TestFileManeuverService_OverwriteMode(t)
}

// TestFileManeuverService_CopyModeWithOverwrite はコピーモード＋上書きモードをテストします
func (tc *TestFileManeuverService) TestFileManeuverService_CopyModeWithOverwrite(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	destDir := filepath.Join(tempDir, "dest")
	os.MkdirAll(srcDir, 0755)
	os.MkdirAll(destDir, 0755)

	// ソースファイルを作成
	srcFile := filepath.Join(srcDir, "test.jpg")
	srcContent := "new content for copy"
	os.WriteFile(srcFile, []byte(srcContent), 0644)

	// 宛先に同名ファイルを作成（上書き対象）
	destFile := filepath.Join(destDir, "test.jpg")
	oldContent := "old content"
	os.WriteFile(destFile, []byte(oldContent), 0644)

	config, err := NewConfig([]string{srcDir}, []string{"jpg"}, "", destDir, false, 1, false, true, true) // コピー＋上書きモード
	if err != nil {
		t.Fatalf("設定作成に失敗しました: %v", err)
	}

	service := NewFileManeuverService(config)
	var stdout, stderr bytes.Buffer

	// Act
	successCount, errorCount, err := service.ExecuteFileManeuver(&stdout, &stderr)

	// Assert
	if err != nil {
		t.Fatalf("ファイルコピー処理に失敗しました: %v", err)
	}

	if successCount != 1 {
		t.Errorf("成功カウントが期待値と異なります。期待値: 1, 実際: %d", successCount)
	}

	if errorCount != 0 {
		t.Errorf("エラーカウントが期待値と異なります。期待値: 0, 実際: %d", errorCount)
	}

	// 元ファイルが残っていることを確認（コピーモード）
	if _, err := os.Stat(srcFile); os.IsNotExist(err) {
		t.Error("コピーモードで元ファイルが削除されました")
	}

	// 宛先ファイルが新しい内容で上書きされていることを確認
	newContent, err := os.ReadFile(destFile)
	if err != nil {
		t.Fatalf("宛先ファイルの読み取りに失敗しました: %v", err)
	}

	if string(newContent) != srcContent {
		t.Errorf("ファイルが正しく上書きコピーされていません。期待値: %s, 実際: %s", srcContent, string(newContent))
	}

	// 出力にコピー＋上書きメッセージが含まれることを確認
	output := stdout.String()
	if !strings.Contains(output, "ファイルコピーを開始します") {
		t.Error("コピーメッセージが出力されませんでした")
	}
	if !strings.Contains(output, "上書きモード") {
		t.Error("上書きメッセージが出力されませんでした")
	}
}

func TestFileManeuverService_CopyModeWithOverwrite(t *testing.T) {
	tc := &TestFileManeuverService{}
	tc.TestFileManeuverService_CopyModeWithOverwrite(t)
}

// TestFileManeuverService_MultipleWorkers は複数ワーカーでの処理をテストします
func (tc *TestFileManeuverService) TestFileManeuverService_MultipleWorkers(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	destDir := filepath.Join(tempDir, "dest")
	os.MkdirAll(srcDir, 0755)
	os.MkdirAll(destDir, 0755)

	// 複数のテストファイルを作成
	testFiles := []string{"file1.jpg", "file2.jpg", "file3.jpg", "file4.jpg", "file5.jpg"}
	for _, file := range testFiles {
		filePath := filepath.Join(srcDir, file)
		os.WriteFile(filePath, []byte("test content"), 0644)
	}

	config, err := NewConfig([]string{srcDir}, []string{"jpg"}, "", destDir, false, 3, false, false, false) // 3ワーカー
	if err != nil {
		t.Fatalf("設定作成に失敗しました: %v", err)
	}

	service := NewFileManeuverService(config)
	var stdout, stderr bytes.Buffer

	// Act
	successCount, errorCount, err := service.ExecuteFileManeuver(&stdout, &stderr)

	// Assert
	if err != nil {
		t.Fatalf("ファイル移動処理に失敗しました: %v", err)
	}

	if successCount != len(testFiles) {
		t.Errorf("成功カウントが期待値と異なります。期待値: %d, 実際: %d", len(testFiles), successCount)
	}

	if errorCount != 0 {
		t.Errorf("エラーカウントが期待値と異なります。期待値: 0, 実際: %d", errorCount)
	}

	// 出力にワーカー数メッセージが含まれることを確認
	output := stdout.String()
	if !strings.Contains(output, "3 ワーカーを使用します") {
		t.Error("ワーカー数メッセージが出力されませんでした")
	}

	// 全ファイルが正しく移動されたことを確認
	for _, file := range testFiles {
		srcPath := filepath.Join(srcDir, file)
		destPath := filepath.Join(destDir, file)

		// 元ファイルが削除されていることを確認
		if _, err := os.Stat(srcPath); !os.IsNotExist(err) {
			t.Errorf("元ファイルが削除されていません: %s", srcPath)
		}

		// 宛先ファイルが存在することを確認
		if _, err := os.Stat(destPath); os.IsNotExist(err) {
			t.Errorf("宛先ファイルが作成されていません: %s", destPath)
		}
	}
}

func TestFileManeuverService_MultipleWorkers(t *testing.T) {
	tc := &TestFileManeuverService{}
	tc.TestFileManeuverService_MultipleWorkers(t)
}

// TestFileManeuverService_WorkersMoreThanFiles はワーカー数がファイル数より多い場合をテストします
func (tc *TestFileManeuverService) TestFileManeuverService_WorkersMoreThanFiles(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	destDir := filepath.Join(tempDir, "dest")
	os.MkdirAll(srcDir, 0755)
	os.MkdirAll(destDir, 0755)

	// 少数のテストファイルを作成
	testFiles := []string{"file1.jpg", "file2.jpg"}
	for _, file := range testFiles {
		filePath := filepath.Join(srcDir, file)
		os.WriteFile(filePath, []byte("test content"), 0644)
	}

	config, err := NewConfig([]string{srcDir}, []string{"jpg"}, "", destDir, false, 10, false, false, false) // 10ワーカー（ファイル数より多い）
	if err != nil {
		t.Fatalf("設定作成に失敗しました: %v", err)
	}

	service := NewFileManeuverService(config)
	var stdout, stderr bytes.Buffer

	// Act
	successCount, errorCount, err := service.ExecuteFileManeuver(&stdout, &stderr)

	// Assert
	if err != nil {
		t.Fatalf("ファイル移動処理に失敗しました: %v", err)
	}

	if successCount != len(testFiles) {
		t.Errorf("成功カウントが期待値と異なります。期待値: %d, 実際: %d", len(testFiles), successCount)
	}

	if errorCount != 0 {
		t.Errorf("エラーカウントが期待値と異なります。期待値: 0, 実際: %d", errorCount)
	}

	// 出力にワーカー数がファイル数に調整されたメッセージが含まれることを確認
	output := stdout.String()
	if !strings.Contains(output, "2 ワーカーを使用します") {
		t.Error("調整されたワーカー数メッセージが出力されませんでした")
	}
}

func TestFileManeuverService_WorkersMoreThanFiles(t *testing.T) {
	tc := &TestFileManeuverService{}
	tc.TestFileManeuverService_WorkersMoreThanFiles(t)
}

// TestFileManeuverService_MultipleSrcDirs は複数ソースディレクトリをテストします
func (tc *TestFileManeuverService) TestFileManeuverService_MultipleSrcDirs(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	srcDir1 := filepath.Join(tempDir, "src1")
	srcDir2 := filepath.Join(tempDir, "src2")
	destDir := filepath.Join(tempDir, "dest")
	os.MkdirAll(srcDir1, 0755)
	os.MkdirAll(srcDir2, 0755)
	os.MkdirAll(destDir, 0755)

	// 各ソースディレクトリにファイルを作成
	os.WriteFile(filepath.Join(srcDir1, "file1.jpg"), []byte("content1"), 0644)
	os.WriteFile(filepath.Join(srcDir1, "file2.png"), []byte("content2"), 0644)
	os.WriteFile(filepath.Join(srcDir2, "file3.jpg"), []byte("content3"), 0644)
	os.WriteFile(filepath.Join(srcDir2, "file4.png"), []byte("content4"), 0644)

	config, err := NewConfig([]string{srcDir1, srcDir2}, []string{"jpg", "png"}, "", destDir, false, 2, false, false, false)
	if err != nil {
		t.Fatalf("設定作成に失敗しました: %v", err)
	}

	service := NewFileManeuverService(config)
	var stdout, stderr bytes.Buffer

	// Act
	successCount, errorCount, err := service.ExecuteFileManeuver(&stdout, &stderr)

	// Assert
	if err != nil {
		t.Fatalf("ファイル移動処理に失敗しました: %v", err)
	}

	if successCount != 4 {
		t.Errorf("成功カウントが期待値と異なります。期待値: 4, 実際: %d", successCount)
	}

	if errorCount != 0 {
		t.Errorf("エラーカウントが期待値と異なります。期待値: 0, 実際: %d", errorCount)
	}

	// 各ディレクトリからファイルが発見されたメッセージを確認
	output := stdout.String()
	if !strings.Contains(output, "から 2 ファイルを発見しました") {
		t.Error("各ディレクトリからのファイル発見メッセージが出力されませんでした")
	}

	// 全ファイルが正しく移動されたことを確認
	expectedFiles := []string{"file1.jpg", "file2.png", "file3.jpg", "file4.png"}
	for _, file := range expectedFiles {
		destPath := filepath.Join(destDir, file)
		if _, err := os.Stat(destPath); os.IsNotExist(err) {
			t.Errorf("宛先ファイルが作成されていません: %s", destPath)
		}
	}
}

func TestFileManeuverService_MultipleSrcDirs(t *testing.T) {
	tc := &TestFileManeuverService{}
	tc.TestFileManeuverService_MultipleSrcDirs(t)
}

// TestFileManeuverService_NonRecursiveMode は非再帰モードをテストします
func (tc *TestFileManeuverService) TestFileManeuverService_NonRecursiveMode(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	subDir := filepath.Join(srcDir, "subdir")
	destDir := filepath.Join(tempDir, "dest")
	os.MkdirAll(subDir, 0755)
	os.MkdirAll(destDir, 0755)

	// ルートディレクトリにファイル作成
	os.WriteFile(filepath.Join(srcDir, "root.jpg"), []byte("root content"), 0644)
	// サブディレクトリにファイル作成（非再帰モードでは無視される）
	os.WriteFile(filepath.Join(subDir, "sub.jpg"), []byte("sub content"), 0644)

	config, err := NewConfig([]string{srcDir}, []string{"jpg"}, "", destDir, false, 1, false, false, false) // 非再帰モード
	if err != nil {
		t.Fatalf("設定作成に失敗しました: %v", err)
	}

	service := NewFileManeuverService(config)
	var stdout, stderr bytes.Buffer

	// Act
	files, err := service.FindTargetFiles(&stdout, &stderr)

	// Assert
	if err != nil {
		t.Fatalf("ファイル検索に失敗しました: %v", err)
	}

	// 非再帰モードではルートディレクトリのファイルのみが見つかる
	expectedCount := 1
	if len(files) != expectedCount {
		t.Errorf("見つかったファイル数が期待値と異なります。期待値: %d, 実際: %d", expectedCount, len(files))
	}

	// ルートファイルのみが含まれることを確認
	foundFiles := make(map[string]bool)
	for _, file := range files {
		foundFiles[filepath.Base(file)] = true
	}

	if !foundFiles["root.jpg"] {
		t.Error("ルートディレクトリのファイルが見つかりませんでした")
	}

	if foundFiles["sub.jpg"] {
		t.Error("非再帰モードでサブディレクトリのファイルが見つかりました")
	}
}

func TestFileManeuverService_NonRecursiveMode(t *testing.T) {
	tc := &TestFileManeuverService{}
	tc.TestFileManeuverService_NonRecursiveMode(t)
}

// TestFileManeuverService_EmptyDirectory は空のディレクトリ処理をテストします
func (tc *TestFileManeuverService) TestFileManeuverService_EmptyDirectory(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	destDir := filepath.Join(tempDir, "dest")
	os.MkdirAll(srcDir, 0755)
	os.MkdirAll(destDir, 0755)

	// 空のディレクトリ（対象ファイルなし）

	config, err := NewConfig([]string{srcDir}, []string{"jpg"}, "", destDir, false, 1, false, false, false)
	if err != nil {
		t.Fatalf("設定作成に失敗しました: %v", err)
	}

	service := NewFileManeuverService(config)
	var stdout, stderr bytes.Buffer

	// Act
	successCount, errorCount, err := service.ExecuteFileManeuver(&stdout, &stderr)

	// Assert
	if err != nil {
		t.Fatalf("ファイル移動処理に失敗しました: %v", err)
	}

	if successCount != 0 {
		t.Errorf("成功カウントが期待値と異なります。期待値: 0, 実際: %d", successCount)
	}

	if errorCount != 0 {
		t.Errorf("エラーカウントが期待値と異なります。期待値: 0, 実際: %d", errorCount)
	}

	// 出力に対象ファイルなしメッセージが含まれることを確認
	output := stdout.String()
	if !strings.Contains(output, "移動対象のファイルが見つかりませんでした") {
		t.Error("対象ファイルなしメッセージが出力されませんでした")
	}
}

func TestFileManeuverService_EmptyDirectory(t *testing.T) {
	tc := &TestFileManeuverService{}
	tc.TestFileManeuverService_EmptyDirectory(t)
}

// TestFileManeuverService_DryRunCopyMode はドライランコピーモードをテストします
func (tc *TestFileManeuverService) TestFileManeuverService_DryRunCopyMode(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	destDir := filepath.Join(tempDir, "dest")
	os.MkdirAll(srcDir, 0755)
	os.MkdirAll(destDir, 0755)

	// テストファイルを作成
	testFile := filepath.Join(srcDir, "test.jpg")
	os.WriteFile(testFile, []byte("test content"), 0644)

	config, err := NewConfig([]string{srcDir}, []string{"jpg"}, "", destDir, false, 1, true, true, false) // ドライラン＋コピーモード
	if err != nil {
		t.Fatalf("設定作成に失敗しました: %v", err)
	}

	service := NewFileManeuverService(config)
	var stdout, stderr bytes.Buffer

	// Act
	successCount, errorCount, err := service.ExecuteFileManeuver(&stdout, &stderr)

	// Assert
	if err != nil {
		t.Fatalf("ファイルコピー処理に失敗しました: %v", err)
	}

	if successCount != 1 {
		t.Errorf("成功カウントが期待値と異なります。期待値: 1, 実際: %d", successCount)
	}

	if errorCount != 0 {
		t.Errorf("エラーカウントが期待値と異なります。期待値: 0, 実際: %d", errorCount)
	}

	// ドライランなので元ファイルが残っていることを確認
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("ドライランモードで元ファイルが削除されました")
	}

	// 宛先にファイルがコピーされていないことを確認
	destFile := filepath.Join(destDir, "test.jpg")
	if _, err := os.Stat(destFile); !os.IsNotExist(err) {
		t.Error("ドライランモードで宛先にファイルが作成されました")
	}

	// 出力にドライラン＋コピーメッセージが含まれることを確認
	output := stdout.String()
	if !strings.Contains(output, "ドライランモード") {
		t.Error("ドライランメッセージが出力されませんでした")
	}
	if !strings.Contains(output, "コピー予定") {
		t.Error("コピー予定メッセージが出力されませんでした")
	}
}

func TestFileManeuverService_DryRunCopyMode(t *testing.T) {
	tc := &TestFileManeuverService{}
	tc.TestFileManeuverService_DryRunCopyMode(t)
}
