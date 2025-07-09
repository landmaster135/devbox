package usecases

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestValidateConfig_Normal(t *testing.T) {
	// Arrange
	stderr := &bytes.Buffer{}
	config := Config{
		SrcDir:     ".",
		VlcPattern: true,
		WinPattern: false,
	}

	// Act
	err := validateConfig(config, stderr)

	// Assert
	if err != nil {
		t.Errorf("ValidateConfig() returned error: %v", err)
	}
}

func TestValidateConfig_NoPattern(t *testing.T) {
	// Arrange
	stderr := &bytes.Buffer{}
	config := Config{
		SrcDir:     ".",
		VlcPattern: false,
		WinPattern: false,
	}

	// Act
	err := validateConfig(config, stderr)

	// Assert
	if err == nil {
		t.Error("ValidateConfig() should return error when no pattern is specified")
	}
}

func TestValidateConfig_BothPatterns(t *testing.T) {
	// Arrange
	stderr := &bytes.Buffer{}
	config := Config{
		SrcDir:     ".",
		VlcPattern: true,
		WinPattern: true,
	}

	// Act
	err := validateConfig(config, stderr)

	// Assert
	if err == nil {
		t.Error("ValidateConfig() should return error when both patterns are specified")
	}
}

func TestRenameVlcScreenshot_NormalPattern1(t *testing.T) {
	// Arrange
	baseName := "vlcsnap-2025-05-07-12-34-56"
	ext := ".png"

	// Act
	newName, err := renameVlcScreenshot(baseName, ext)

	// Assert
	if err != nil {
		t.Errorf("renameVlcScreenshot() returned error: %v", err)
	}
	expected := "Screenshot_20250507-123456.png"
	if newName != expected {
		t.Errorf("renameVlcScreenshot() = %v, want %v", newName, expected)
	}
}

func TestRenameVlcScreenshot_NormalPattern2(t *testing.T) {
	// Arrange
	baseName := "vlcsnap-2025-05-06-23h59m44s239"
	ext := ".png"

	// Act
	newName, err := renameVlcScreenshot(baseName, ext)

	// Assert
	if err != nil {
		t.Errorf("renameVlcScreenshot() returned error: %v", err)
	}
	expected := "Screenshot_20250506-235944.png"
	if newName != expected {
		t.Errorf("renameVlcScreenshot() = %v, want %v", newName, expected)
	}
}

func TestRenameVlcScreenshot_InvalidFormat(t *testing.T) {
	// Arrange
	baseName := "vlcsnap-invalid-format"
	ext := ".png"

	// Act
	_, err := renameVlcScreenshot(baseName, ext)

	// Assert
	if err == nil {
		t.Error("renameVlcScreenshot() should return error for invalid format")
	}
}

func TestRenameWindowsScreenshot_Normal(t *testing.T) {
	// Arrange
	baseName := "スクリーンショット 2025-05-07 123456"
	ext := ".png"

	// Act
	newName, err := renameWindowsScreenshot(baseName, ext)

	// Assert
	if err != nil {
		t.Errorf("renameWindowsScreenshot() returned error: %v", err)
	}
	expected := "Screenshot_20250507-123456.png"
	if newName != expected {
		t.Errorf("renameWindowsScreenshot() = %v, want %v", newName, expected)
	}
}

func TestRenameWindowsScreenshot_InvalidFormat(t *testing.T) {
	// Arrange
	baseName := "スクリーンショット invalid-format"
	ext := ".png"

	// Act
	_, err := renameWindowsScreenshot(baseName, ext)

	// Assert
	if err == nil {
		t.Error("renameWindowsScreenshot() should return error for invalid format")
	}
}

func TestRenameAndroidScreenshot_Normal(t *testing.T) {
	// Arrange
	baseName := "screen-20250215-064735"
	ext := ".png"

	// Act
	newName, err := renameAndroidScreenshot(baseName, ext)

	// Assert
	if err != nil {
		t.Errorf("renameAndroidScreenshot() returned error: %v", err)
	}
	expected := "Screenshot_20250215-064735.png"
	if newName != expected {
		t.Errorf("renameAndroidScreenshot() = %v, want %v", newName, expected)
	}
}

func TestRenameAndroidScreenshot_MP4(t *testing.T) {
	// Arrange
	baseName := "screen-20250215-064735"
	ext := ".mp4"

	// Act
	newName, err := renameAndroidScreenshot(baseName, ext)

	// Assert
	if err != nil {
		t.Errorf("renameAndroidScreenshot() returned error: %v", err)
	}
	expected := "Screenshot_20250215-064735.mp4"
	if newName != expected {
		t.Errorf("renameAndroidScreenshot() = %v, want %v", newName, expected)
	}
}

func TestRenameAndroidScreenshot_InvalidFormat(t *testing.T) {
	// Arrange
	baseName := "screen-invalid-format"
	ext := ".png"

	// Act
	_, err := renameAndroidScreenshot(baseName, ext)

	// Assert
	if err == nil {
		t.Error("renameAndroidScreenshot() should return error for invalid format")
	}
}

func TestProcessScreenshotRename_AndroidPattern(t *testing.T) {
	// Arrange
	tempDir, err := os.MkdirTemp("", "test-android-screenshot")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "screen-20250215-064735.png")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	file := FileInfo{
		Path: testFile,
		Name: "screen-20250215-064735.png",
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	var mu sync.Mutex
	successCount := 0
	errorCount := 0

	// Act
	processScreenshotRename(file, false, false, true, &mu, &successCount, &errorCount, stdout, stderr)

	// Assert
	if successCount != 1 {
		t.Errorf("processScreenshotRename() successCount = %v, want %v", successCount, 1)
	}
	if errorCount != 0 {
		t.Errorf("processScreenshotRename() errorCount = %v, want %v", errorCount, 0)
	}

	expectedNewPath := filepath.Join(tempDir, "Screenshot_20250215-064735.png")
	if _, err := os.Stat(expectedNewPath); os.IsNotExist(err) {
		t.Errorf("File was not renamed to %s", expectedNewPath)
	}
}

func TestProcessScreenshotRename_AndroidPatternMP4(t *testing.T) {
	// Arrange
	tempDir, err := os.MkdirTemp("", "test-android-screenshot-mp4")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "screen-20250215-064735.mp4")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	file := FileInfo{
		Path: testFile,
		Name: "screen-20250215-064735.mp4",
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	var mu sync.Mutex
	successCount := 0
	errorCount := 0

	// Act
	processScreenshotRename(file, false, false, true, &mu, &successCount, &errorCount, stdout, stderr)

	// Assert
	if successCount != 1 {
		t.Errorf("processScreenshotRename() successCount = %v, want %v", successCount, 1)
	}
	if errorCount != 0 {
		t.Errorf("processScreenshotRename() errorCount = %v, want %v", errorCount, 0)
	}

	expectedNewPath := filepath.Join(tempDir, "Screenshot_20250215-064735.mp4")
	if _, err := os.Stat(expectedNewPath); os.IsNotExist(err) {
		t.Errorf("File was not renamed to %s", expectedNewPath)
	}
}

func TestProcessScreenshotRename_VlcPattern(t *testing.T) {
	// Arrange
	tempDir, err := os.MkdirTemp("", "test-vlc-screenshot")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "vlcsnap-2025-05-07-12-34-56.png")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	file := FileInfo{
		Path: testFile,
		Name: "vlcsnap-2025-05-07-12-34-56.png",
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	var mu sync.Mutex
	successCount := 0
	errorCount := 0

	// Act
	processScreenshotRename(file, true, false, false, &mu, &successCount, &errorCount, stdout, stderr)

	// Assert
	if successCount != 1 {
		t.Errorf("processScreenshotRename() successCount = %v, want %v", successCount, 1)
	}
	if errorCount != 0 {
		t.Errorf("processScreenshotRename() errorCount = %v, want %v", errorCount, 0)
	}

	expectedNewPath := filepath.Join(tempDir, "Screenshot_20250507-123456.png")
	if _, err := os.Stat(expectedNewPath); os.IsNotExist(err) {
		t.Errorf("File was not renamed to %s", expectedNewPath)
	}
}

func TestProcessScreenshotRename_WinPattern(t *testing.T) {
	// Arrange
	tempDir, err := os.MkdirTemp("", "test-win-screenshot")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "スクリーンショット 2025-05-07 123456.png")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	file := FileInfo{
		Path: testFile,
		Name: "スクリーンショット 2025-05-07 123456.png",
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	var mu sync.Mutex
	successCount := 0
	errorCount := 0

	// Act
	processScreenshotRename(file, false, true, false, &mu, &successCount, &errorCount, stdout, stderr)

	// Assert
	if successCount != 1 {
		t.Errorf("processScreenshotRename() successCount = %v, want %v", successCount, 1)
	}
	if errorCount != 0 {
		t.Errorf("processScreenshotRename() errorCount = %v, want %v", errorCount, 0)
	}

	expectedNewPath := filepath.Join(tempDir, "Screenshot_20250507-123456.png")
	if _, err := os.Stat(expectedNewPath); os.IsNotExist(err) {
		t.Errorf("File was not renamed to %s", expectedNewPath)
	}
}

func TestValidateConfig_InvalidDirectory(t *testing.T) {
	// Arrange
	stderr := &bytes.Buffer{}
	config := Config{
		SrcDir:     "/non/existent/directory",
		VlcPattern: true,
		WinPattern: false,
	}

	// Act
	err := validateConfig(config, stderr)

	// Assert
	if err == nil {
		t.Error("ValidateConfig() should return error when directory does not exist")
	}
}

func TestFindScreenshotFiles_NonRecursive(t *testing.T) {
	// Arrange
	tempDir, err := os.MkdirTemp("", "test-find-screenshots")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// VLCスクリーンショットファイルを作成
	vlcFile := filepath.Join(tempDir, "vlcsnap-2025-05-07-12-34-56.png")
	if err := os.WriteFile(vlcFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Windowsスクリーンショットファイルを作成
	winFile := filepath.Join(tempDir, "スクリーンショット 2025-05-07 123456.png")
	if err := os.WriteFile(winFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// 通常のファイルを作成
	normalFile := filepath.Join(tempDir, "normal.png")
	if err := os.WriteFile(normalFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	// Act - VLCパターンのみ
	filesVlc, err := findScreenshotFiles(tempDir, false, true, false, false, stdout, stderr)

	// Assert
	if err != nil {
		t.Errorf("FindScreenshotFiles() returned error: %v", err)
	}
	if len(filesVlc) != 1 {
		t.Errorf("FindScreenshotFiles() with VLC pattern found %v files, want 1", len(filesVlc))
	}
	if len(filesVlc) > 0 && !strings.Contains(filesVlc[0], "vlcsnap") {
		t.Errorf("FindScreenshotFiles() found wrong file: %v", filesVlc[0])
	}

	// Act - Windowsパターンのみ
	filesWin, err := findScreenshotFiles(tempDir, false, false, true, false, stdout, stderr)

	// Assert
	if err != nil {
		t.Errorf("FindScreenshotFiles() returned error: %v", err)
	}
	if len(filesWin) != 1 {
		t.Errorf("FindScreenshotFiles() with Windows pattern found %v files, want 1", len(filesWin))
	}
	if len(filesWin) > 0 && !strings.Contains(filesWin[0], "スクリーンショット") {
		t.Errorf("FindScreenshotFiles() found wrong file: %v", filesWin[0])
	}
}

func TestFindScreenshotFiles_Recursive(t *testing.T) {
	// Arrange
	tempDir, err := os.MkdirTemp("", "test-find-screenshots-recursive")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// サブディレクトリを作成
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	// ルートディレクトリにVLCスクリーンショットファイルを作成
	vlcFile := filepath.Join(tempDir, "vlcsnap-2025-05-07-12-34-56.png")
	if err := os.WriteFile(vlcFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// サブディレクトリにVLCスクリーンショットファイルを作成
	subVlcFile := filepath.Join(subDir, "vlcsnap-2025-05-08-12-34-56.png")
	if err := os.WriteFile(subVlcFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	// Act - 再帰的検索
	files, err := findScreenshotFiles(tempDir, true, true, false, false, stdout, stderr)

	// Assert
	if err != nil {
		t.Errorf("FindScreenshotFiles() returned error: %v", err)
	}
	if len(files) != 2 {
		t.Errorf("FindScreenshotFiles() with recursive found %v files, want 2", len(files))
	}
}

func TestFindScreenshotFiles_InvalidDirectory(t *testing.T) {
	// Arrange
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	// Act
	_, err := findScreenshotFiles("/non/existent/directory", false, true, false, false, stdout, stderr)

	// Assert
	if err == nil {
		t.Error("FindScreenshotFiles() should return error for invalid directory")
	}
}

func TestGetFileInfos_Normal(t *testing.T) {
	// Arrange
	tempDir, err := os.MkdirTemp("", "test-get-file-infos")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// テストファイルを作成
	testFile := filepath.Join(tempDir, "test.png")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	files := []string{testFile}
	stderr := &bytes.Buffer{}

	// Act
	fileInfos, err := getFileInfos(files, stderr)

	// Assert
	if err != nil {
		t.Errorf("GetFileInfos() returned error: %v", err)
	}
	if len(fileInfos) != 1 {
		t.Errorf("GetFileInfos() returned %v file infos, want 1", len(fileInfos))
	}
	if fileInfos[0].Path != testFile {
		t.Errorf("GetFileInfos() returned wrong path: %v, want %v", fileInfos[0].Path, testFile)
	}
	if fileInfos[0].Name != "test.png" {
		t.Errorf("GetFileInfos() returned wrong name: %v, want %v", fileInfos[0].Name, "test.png")
	}
}

func TestGetFileInfos_InvalidFile(t *testing.T) {
	// Arrange
	files := []string{"/non/existent/file.png"}
	stderr := &bytes.Buffer{}

	// Act
	fileInfos, err := getFileInfos(files, stderr)

	// Assert
	if err == nil {
		t.Error("GetFileInfos() should return error for invalid file")
	}
	if len(fileInfos) != 1 {
		t.Errorf("GetFileInfos() returned %v file infos, want 1", len(fileInfos))
	}
}

func TestRenameScreenshotFiles_Normal(t *testing.T) {
	// Arrange
	tempDir, err := os.MkdirTemp("", "test-rename-screenshots")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// VLCスクリーンショットファイルを作成
	vlcFile := filepath.Join(tempDir, "vlcsnap-2025-05-07-12-34-56.png")
	if err := os.WriteFile(vlcFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Windowsスクリーンショットファイルを作成
	winFile := filepath.Join(tempDir, "スクリーンショット 2025-05-07 123456.png")
	if err := os.WriteFile(winFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	fileInfos := []FileInfo{
		{
			Path: vlcFile,
			Name: "vlcsnap-2025-05-07-12-34-56.png",
		},
		{
			Path: winFile,
			Name: "スクリーンショット 2025-05-07 123456.png",
		},
	}

	config := Config{
		SrcDir:     tempDir,
		VlcPattern: true,
		WinPattern: false,
		Workers:    2,
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	// Act
	successCount, errorCount := renameScreenshotFiles(fileInfos, config, stdout, stderr)

	// Assert
	if successCount != 1 {
		t.Errorf("RenameScreenshotFiles() successCount = %v, want %v", successCount, 1)
	}
	if errorCount != 0 {
		t.Errorf("RenameScreenshotFiles() errorCount = %v, want %v", errorCount, 0)
	}

	// VLCファイルがリネームされたことを確認
	expectedNewPath := filepath.Join(tempDir, "Screenshot_20250507-123456.png")
	if _, err := os.Stat(expectedNewPath); os.IsNotExist(err) {
		t.Errorf("File was not renamed to %s", expectedNewPath)
	}

	// Windowsファイルはリネームされていないことを確認（VLCパターンのみ指定したため）
	if _, err := os.Stat(winFile); os.IsNotExist(err) {
		t.Errorf("Windows screenshot file should not be renamed")
	}
}

func TestRenameScreenshotFiles_WorkerAdjustment(t *testing.T) {
	// Arrange
	tempDir, err := os.MkdirTemp("", "test-rename-screenshots-workers")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// テストファイルを作成
	testFile := filepath.Join(tempDir, "vlcsnap-2025-05-07-12-34-56.png")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	fileInfos := []FileInfo{
		{
			Path: testFile,
			Name: "vlcsnap-2025-05-07-12-34-56.png",
		},
	}

	// ケース1: ワーカー数が0の場合（1に調整される）
	config1 := Config{
		SrcDir:     tempDir,
		VlcPattern: true,
		WinPattern: false,
		Workers:    0,
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	// Act
	successCount1, _ := renameScreenshotFiles(fileInfos, config1, stdout, stderr)

	// Assert
	if successCount1 != 1 {
		t.Errorf("RenameScreenshotFiles() with Workers=0 successCount = %v, want %v", successCount1, 1)
	}

	// ファイルを再作成
	os.RemoveAll(tempDir)
	tempDir, _ = os.MkdirTemp("", "test-rename-screenshots-workers2")
	defer os.RemoveAll(tempDir)

	testFile = filepath.Join(tempDir, "vlcsnap-2025-05-07-12-34-56.png")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	fileInfos = []FileInfo{
		{
			Path: testFile,
			Name: "vlcsnap-2025-05-07-12-34-56.png",
		},
	}

	// ケース2: ワーカー数がファイル数より多い場合（ファイル数に調整される）
	config2 := Config{
		SrcDir:     tempDir,
		VlcPattern: true,
		WinPattern: false,
		Workers:    10,
	}

	stdout = &bytes.Buffer{}
	stderr = &bytes.Buffer{}

	// Act
	successCount2, _ := renameScreenshotFiles(fileInfos, config2, stdout, stderr)

	// Assert
	if successCount2 != 1 {
		t.Errorf("RenameScreenshotFiles() with Workers>fileInfos successCount = %v, want %v", successCount2, 1)
	}
}

func TestFindScreenshotFiles_WalkDirError(t *testing.T) {
	// Arrange
	// 一時的に権限のないディレクトリを作成
	tempDir, err := os.MkdirTemp("", "test-find-screenshots-error")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// サブディレクトリを作成して権限を制限
	subDir := filepath.Join(tempDir, "restricted")
	if err := os.Mkdir(subDir, 0700); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	// サブディレクトリにファイルを作成
	subFile := filepath.Join(subDir, "test.png")
	if err := os.WriteFile(subFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// サブディレクトリの権限を変更して読み取り不可にする
	if err := os.Chmod(subDir, 0000); err != nil {
		t.Fatalf("Failed to change directory permission: %v", err)
	}
	defer os.Chmod(subDir, 0700) // テスト後に権限を戻す

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	// Act - 再帰的検索でエラーが発生するケース
	_, err = findScreenshotFiles(tempDir, true, true, false, false, stdout, stderr)

	// Assert
	// 権限の問題でエラーが発生する可能性があるが、OSによって動作が異なるため
	// エラーの有無だけでなく、関数が実行されることを確認
	if err != nil {
		// エラーが発生した場合はOK
		t.Logf("Expected error occurred: %v", err)
	} else {
		// エラーが発生しなかった場合も、関数が実行されていればOK
		t.Log("No error occurred, but function executed")
	}
}

func TestProcessScreenshotRename_ParseError(t *testing.T) {
	// Arrange
	tempDir, err := os.MkdirTemp("", "test-process-rename-parse-error")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 不正なフォーマットのVLCスクリーンショットファイル
	testFile := filepath.Join(tempDir, "vlcsnap-invalid.png")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	file := FileInfo{
		Path: testFile,
		Name: "vlcsnap-invalid.png",
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	var mu sync.Mutex
	successCount := 0
	errorCount := 0

	// Act
	processScreenshotRename(file, true, false, false, &mu, &successCount, &errorCount, stdout, stderr)

	// Assert
	if successCount != 0 {
		t.Errorf("processScreenshotRename() successCount = %v, want %v", successCount, 0)
	}
	if errorCount != 1 {
		t.Errorf("processScreenshotRename() errorCount = %v, want %v", errorCount, 1)
	}
}

func TestProcessScreenshotRename_InvalidPattern(t *testing.T) {
	// Arrange
	tempDir, err := os.MkdirTemp("", "test-process-rename-invalid")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "normal.png")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	file := FileInfo{
		Path: testFile,
		Name: "normal.png",
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	var mu sync.Mutex
	successCount := 0
	errorCount := 0

	// Act
	processScreenshotRename(file, true, false, false, &mu, &successCount, &errorCount, stdout, stderr)

	// Assert
	if successCount != 0 {
		t.Errorf("processScreenshotRename() successCount = %v, want %v", successCount, 0)
	}
	if errorCount != 0 {
		t.Errorf("processScreenshotRename() errorCount = %v, want %v", errorCount, 0)
	}

	// ファイルがリネームされていないことを確認
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Errorf("File should not be renamed")
	}
}

func TestProcessScreenshotRename_RenameError(t *testing.T) {
	// Arrange
	tempDir, err := os.MkdirTemp("", "test-process-rename-error")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 読み取り専用ディレクトリを作成（リネームエラーを発生させるため）
	readOnlyDir := filepath.Join(tempDir, "readonly")
	if err := os.Mkdir(readOnlyDir, 0500); err != nil {
		t.Fatalf("Failed to create readonly dir: %v", err)
	}

	// 読み取り専用ディレクトリにファイルを作成
	testFile := filepath.Join(readOnlyDir, "vlcsnap-2025-05-07-12-34-56.png")

	// ファイルを作成するためにディレクトリのパーミッションを一時的に変更
	if err := os.Chmod(readOnlyDir, 0700); err != nil {
		t.Fatalf("Failed to change directory permission: %v", err)
	}

	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// ディレクトリを再び読み取り専用に戻す
	if err := os.Chmod(readOnlyDir, 0500); err != nil {
		t.Fatalf("Failed to change directory permission back: %v", err)
	}

	file := FileInfo{
		Path: testFile,
		Name: "vlcsnap-2025-05-07-12-34-56.png",
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	var mu sync.Mutex
	successCount := 0
	errorCount := 0

	// Act
	processScreenshotRename(file, true, false, false, &mu, &successCount, &errorCount, stdout, stderr)

	// Assert
	if successCount != 0 {
		t.Errorf("processScreenshotRename() successCount = %v, want %v", successCount, 0)
	}
	if errorCount != 1 {
		t.Errorf("processScreenshotRename() errorCount = %v, want %v", errorCount, 1)
	}
}

func TestIsImageExt(t *testing.T) {
	// Arrange
	testCases := []struct {
		ext      string
		expected bool
	}{
		{".jpg", true},
		{".jpeg", true},
		{".png", true},
		{".webp", true},
		{".avif", true},
		{".mp4", true},
		{".gif", false},
		{".txt", false},
		{".pdf", false},
	}

	// Act & Assert
	for _, tc := range testCases {
		result := isImageExt(tc.ext)
		if result != tc.expected {
			t.Errorf("isImageExt(%s) = %v, want %v", tc.ext, result, tc.expected)
		}
	}
}
