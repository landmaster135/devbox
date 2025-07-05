package usecases

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestConfigCreation はConfig作成のテストクラスです
type TestConfigCreation struct{}

// TestConfigCreation_Normal は正常なConfig作成をテストします
func (tc *TestConfigCreation) TestConfigCreation_Normal(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	srcDir1 := filepath.Join(tempDir, "src1")
	srcDir2 := filepath.Join(tempDir, "src2")
	destDir := filepath.Join(tempDir, "dest")

	// テスト用ディレクトリを作成
	os.MkdirAll(srcDir1, 0755)
	os.MkdirAll(srcDir2, 0755)
	os.MkdirAll(destDir, 0755)

	srcDirs := []string{srcDir1, srcDir2}
	extensions := []string{"jpg", "png"}

	// Act
	config, err := NewConfig(srcDirs, extensions, destDir, true, 4, false, false, false)

	// Assert
	if err != nil {
		t.Fatalf("設定作成に失敗しました: %v", err)
	}

	if len(config.SrcDirs) != 2 {
		t.Errorf("ソースディレクトリ数が期待値と異なります。期待値: 2, 実際: %d", len(config.SrcDirs))
	}

	if len(config.Extensions) != 2 {
		t.Errorf("拡張子数が期待値と異なります。期待値: 2, 実際: %d", len(config.Extensions))
	}

	// 拡張子の正規化確認
	if config.Extensions[0] != ".jpg" || config.Extensions[1] != ".png" {
		t.Errorf("拡張子の正規化が正しくありません。実際: %v", config.Extensions)
	}

	if config.DestDir != destDir {
		t.Errorf("宛先ディレクトリが期待値と異なります。期待値: %s, 実際: %s", destDir, config.DestDir)
	}

	if !config.Recursive {
		t.Error("再帰フラグが期待値と異なります")
	}

	if config.Workers != 4 {
		t.Errorf("ワーカー数が期待値と異なります。期待値: 4, 実際: %d", config.Workers)
	}
}

// TestConfigValidation はConfig検証のテストクラスです
type TestConfigValidation struct{}

// TestConfigValidation_EmptySrcDirs は空のソースディレクトリでエラーになることをテストします
func (tc *TestConfigValidation) TestConfigValidation_EmptySrcDirs(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	destDir := filepath.Join(tempDir, "dest")
	os.MkdirAll(destDir, 0755)

	// Act
	_, err := NewConfig([]string{}, []string{"jpg"}, destDir, false, 1, false, false, false)

	// Assert
	if err == nil {
		t.Error("空のソースディレクトリでエラーが発生しませんでした")
	}

	if !strings.Contains(err.Error(), "ソースディレクトリが指定されていません") {
		t.Errorf("期待されるエラーメッセージが含まれていません。実際: %s", err.Error())
	}
}

// TestConfigValidation_EmptyExtensions は空の拡張子でエラーになることをテストします
func (tc *TestConfigValidation) TestConfigValidation_EmptyExtensions(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	destDir := filepath.Join(tempDir, "dest")
	os.MkdirAll(srcDir, 0755)
	os.MkdirAll(destDir, 0755)

	// Act
	_, err := NewConfig([]string{srcDir}, []string{}, destDir, false, 1, false, false, false)

	// Assert
	if err == nil {
		t.Error("空の拡張子でエラーが発生しませんでした")
	}

	if !strings.Contains(err.Error(), "拡張子が指定されていません") {
		t.Errorf("期待されるエラーメッセージが含まれていません。実際: %s", err.Error())
	}
}

// TestConfigValidation_NonExistentSrcDir は存在しないソースディレクトリでエラーになることをテストします
func (tc *TestConfigValidation) TestConfigValidation_NonExistentSrcDir(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	destDir := filepath.Join(tempDir, "dest")
	os.MkdirAll(destDir, 0755)

	nonExistentDir := filepath.Join(tempDir, "nonexistent")

	// Act
	_, err := NewConfig([]string{nonExistentDir}, []string{"jpg"}, destDir, false, 1, false, false, false)

	// Assert
	if err == nil {
		t.Error("存在しないソースディレクトリでエラーが発生しませんでした")
	}

	if !strings.Contains(err.Error(), "にアクセスできません") {
		t.Errorf("期待されるエラーメッセージが含まれていません。実際: %s", err.Error())
	}
}

// TestConfigValidation_NonExistentDestDir は存在しない宛先ディレクトリでエラーになることをテストします
func (tc *TestConfigValidation) TestConfigValidation_NonExistentDestDir(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	os.MkdirAll(srcDir, 0755)

	nonExistentDir := filepath.Join(tempDir, "nonexistent")

	// Act
	_, err := NewConfig([]string{srcDir}, []string{"jpg"}, nonExistentDir, false, 1, false, false, false)

	// Assert
	if err == nil {
		t.Error("存在しない宛先ディレクトリでエラーが発生しませんでした")
	}

	if !strings.Contains(err.Error(), "にアクセスできません") {
		t.Errorf("期待されるエラーメッセージが含まれていません。実際: %s", err.Error())
	}
}

// TestWorkerNormalization はワーカー数正規化のテストクラスです
type TestWorkerNormalization struct{}

// TestWorkerNormalization_ZeroWorkers はワーカー数0の場合の正規化をテストします
func (tc *TestWorkerNormalization) TestWorkerNormalization_ZeroWorkers(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	destDir := filepath.Join(tempDir, "dest")
	os.MkdirAll(srcDir, 0755)
	os.MkdirAll(destDir, 0755)

	// Act
	config, err := NewConfig([]string{srcDir}, []string{"jpg"}, destDir, false, 0, false, false, false)

	// Assert
	if err != nil {
		t.Fatalf("設定作成に失敗しました: %v", err)
	}

	expectedWorkers := runtime.NumCPU()
	if config.Workers != expectedWorkers {
		t.Errorf("ワーカー数が期待値と異なります。期待値: %d, 実際: %d", expectedWorkers, config.Workers)
	}
}

// TestWorkerNormalization_ExcessiveWorkers は過剰なワーカー数の場合の正規化をテストします
func (tc *TestWorkerNormalization) TestWorkerNormalization_ExcessiveWorkers(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	destDir := filepath.Join(tempDir, "dest")
	os.MkdirAll(srcDir, 0755)
	os.MkdirAll(destDir, 0755)

	excessiveWorkers := runtime.NumCPU()*2 + 10

	// Act
	config, err := NewConfig([]string{srcDir}, []string{"jpg"}, destDir, false, excessiveWorkers, false, false, false)

	// Assert
	if err != nil {
		t.Fatalf("設定作成に失敗しました: %v", err)
	}

	expectedMaxWorkers := runtime.NumCPU() * 2
	if config.Workers != expectedMaxWorkers {
		t.Errorf("ワーカー数が期待値と異なります。期待値: %d, 実際: %d", expectedMaxWorkers, config.Workers)
	}
}

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

	config, err := NewConfig([]string{srcDir}, []string{"jpg", "png"}, destDir, false, 1, false, false, false)
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

	config, err := NewConfig([]string{srcDir}, []string{"jpg"}, destDir, true, 1, false, false, false)
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

	config, err := NewConfig([]string{srcDir}, []string{"jpg"}, destDir, false, 1, true, false, false) // ドライランモード
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

	config, err := NewConfig([]string{srcDir}, []string{"jpg"}, destDir, false, 1, false, false, false)
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

	config, err := NewConfig([]string{srcDir}, []string{"jpg", "png"}, destDir, false, 1, false, false, false)
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

// 各テストクラスのインスタンスを作成してテストを実行
func TestConfigCreation_Normal(t *testing.T) {
	tc := &TestConfigCreation{}
	tc.TestConfigCreation_Normal(t)
}

func TestConfigValidation_EmptySrcDirs(t *testing.T) {
	tc := &TestConfigValidation{}
	tc.TestConfigValidation_EmptySrcDirs(t)
}

func TestConfigValidation_EmptyExtensions(t *testing.T) {
	tc := &TestConfigValidation{}
	tc.TestConfigValidation_EmptyExtensions(t)
}

func TestConfigValidation_NonExistentSrcDir(t *testing.T) {
	tc := &TestConfigValidation{}
	tc.TestConfigValidation_NonExistentSrcDir(t)
}

func TestConfigValidation_NonExistentDestDir(t *testing.T) {
	tc := &TestConfigValidation{}
	tc.TestConfigValidation_NonExistentDestDir(t)
}

func TestWorkerNormalization_ZeroWorkers(t *testing.T) {
	tc := &TestWorkerNormalization{}
	tc.TestWorkerNormalization_ZeroWorkers(t)
}

func TestWorkerNormalization_ExcessiveWorkers(t *testing.T) {
	tc := &TestWorkerNormalization{}
	tc.TestWorkerNormalization_ExcessiveWorkers(t)
}

func TestFileManeuverService_FindTargetFiles(t *testing.T) {
	tc := &TestFileManeuverService{}
	tc.TestFileManeuverService_FindTargetFiles(t)
}

func TestFileManeuverService_FindTargetFilesRecursive(t *testing.T) {
	tc := &TestFileManeuverService{}
	tc.TestFileManeuverService_FindTargetFilesRecursive(t)
}

func TestFileManeuverService_DryRun(t *testing.T) {
	tc := &TestFileManeuverService{}
	tc.TestFileManeuverService_DryRun(t)
}

func TestFileManeuverService_FileConflict(t *testing.T) {
	tc := &TestFileManeuverService{}
	tc.TestFileManeuverService_FileConflict(t)
}

func TestFileManeuverService_ExecuteFileManeuver(t *testing.T) {
	tc := &TestFileManeuverService{}
	tc.TestFileManeuverService_ExecuteFileManeuver(t)
}
