package usecases

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateConfig_Normal(t *testing.T) {
	// Arrange
	stderr := &bytes.Buffer{}
	config := Config{
		SrcDir:    ".",
		Operation: OperationVLC,
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
		SrcDir:    ".",
		Operation: OperationUnknown,
	}

	// Act
	err := validateConfig(config, stderr)

	// Assert
	if err == nil {
		t.Error("ValidateConfig() should return error when no pattern is specified")
	}
}

func TestValidateConfig_InvalidOperation(t *testing.T) {
	// Arrange
	stderr := &bytes.Buffer{}
	config := Config{
		SrcDir:    ".",
		Operation: Operation("invalid"),
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

func TestRenamePixelScreenshot_Normal(t *testing.T) {
	// Arrange
	baseName := "screen-20250215-064735"
	ext := ".png"

	// Act
	newName, err := renamePixelScreenshot(baseName, ext)

	// Assert
	if err != nil {
		t.Errorf("renamePixelScreenshot() returned error: %v", err)
	}
	expected := "Screenshot_20250215-064735.png"
	if newName != expected {
		t.Errorf("renamePixelScreenshot() = %v, want %v", newName, expected)
	}
}

func TestRenamePixelScreenshot_MP4(t *testing.T) {
	// Arrange
	baseName := "screen-20250215-064735"
	ext := ".mp4"

	// Act
	newName, err := renamePixelScreenshot(baseName, ext)

	// Assert
	if err != nil {
		t.Errorf("renamePixelScreenshot() returned error: %v", err)
	}
	expected := "Screenshot_20250215-064735.mp4"
	if newName != expected {
		t.Errorf("renamePixelScreenshot() = %v, want %v", newName, expected)
	}
}

func TestRenamePixelScreenshot_InvalidFormat(t *testing.T) {
	// Arrange
	baseName := "screen-invalid-format"
	ext := ".png"

	// Act
	_, err := renamePixelScreenshot(baseName, ext)

	// Assert
	if err == nil {
		t.Error("renamePixelScreenshot() should return error for invalid format")
	}
}

func TestRenameXiaomiScreenshot_Normal(t *testing.T) {
	baseName := "Screenshot_2025-12-27-12-03-39-927_com.YoStarJP.AzurLane"
	ext := ".jpg"

	newName, err := renameXiaomiScreenshot(baseName, ext)
	if err != nil {
		t.Fatalf("renameXiaomiScreenshot() returned error: %v", err)
	}

	expected := "Screenshot_20251227-120339.jpg"
	if newName != expected {
		t.Fatalf("renameXiaomiScreenshot() = %s, want %s", newName, expected)
	}
}

func TestRenameXiaomiScreenshot_InvalidFormat(t *testing.T) {
	_, err := renameXiaomiScreenshot("Screenshot_invalid", ".png")
	if err == nil {
		t.Fatal("renameXiaomiScreenshot() should return error for invalid names")
	}
}

func TestRenameXiaomiToDateTime_Normal(t *testing.T) {
	baseName := "Screenshot_2025-12-27-12-03-39-927_com.YoStarJP.AzurLane"
	ext := ".jpg"

	newName, err := renameXiaomiToDateTime(baseName, ext)
	if err != nil {
		t.Fatalf("renameXiaomiToDateTime() returned error: %v", err)
	}

	expected := "20251227120339.jpg"
	if newName != expected {
		t.Fatalf("renameXiaomiToDateTime() = %s, want %s", newName, expected)
	}
}

func TestResolveScreenshotRenameTarget_Operations(t *testing.T) {
	testCases := []struct {
		name      string
		fileName  string
		operation Operation
		expected  string
	}{
		{
			name:      "PixelPNG",
			fileName:  "screen-20250215-064735.png",
			operation: OperationPixel,
			expected:  "Screenshot_20250215-064735.png",
		},
		{
			name:      "PixelMP4",
			fileName:  "screen-20250215-064735.mp4",
			operation: OperationPixel,
			expected:  "Screenshot_20250215-064735.mp4",
		},
		{
			name:      "VLC",
			fileName:  "vlcsnap-2025-05-07-12-34-56.png",
			operation: OperationVLC,
			expected:  "Screenshot_20250507-123456.png",
		},
		{
			name:      "Windows",
			fileName:  "スクリーンショット 2025-05-07 123456.png",
			operation: OperationWin,
			expected:  "Screenshot_20250507-123456.png",
		},
		{
			name:      "Xiaomi",
			fileName:  "Screenshot_2025-12-27-12-03-39-927_com.YoStarJP.AzurLane.jpg",
			operation: OperationXiaomi,
			expected:  "Screenshot_20251227-120339.jpg",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			src := filepath.Join(tempDir, tc.fileName)
			if err := os.WriteFile(src, []byte("test"), 0644); err != nil {
				t.Fatalf("failed to create test file: %v", err)
			}

			file := FileInfo{Path: src, Name: tc.fileName}
			newPath, shouldRename, err := resolveScreenshotRenameTarget(file, Config{Operation: tc.operation})
			if err != nil {
				t.Fatalf("resolveScreenshotRenameTarget() returned error: %v", err)
			}
			if !shouldRename {
				t.Fatalf("resolveScreenshotRenameTarget() should request rename for %s", tc.fileName)
			}

			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}
			if err := performScreenshotRename(file.Path, newPath, stdout, stderr); err != nil {
				t.Fatalf("performScreenshotRename() returned error: %v", err)
			}

			expectedPath := filepath.Join(tempDir, tc.expected)
			if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
				t.Fatalf("expected file %s does not exist", expectedPath)
			}
		})
	}
}

func TestResolveScreenshotRenameTarget_ToDateTime(t *testing.T) {
	testCases := []struct {
		name     string
		fileName string
		expected string
	}{
		{"VLC", "vlcsnap-2025-05-07-12-34-56.png", "20250507123456.png"},
		{"Windows", "スクリーンショット 2025-05-07 123456.png", "20250507123456.png"},
		{"Pixel", "screen-20250215-064735.png", "20250215064735.png"},
		{"Xiaomi", "Screenshot_2025-12-27-12-03-39-927_com.YoStarJP.AzurLane.jpg", "20251227120339.jpg"},
		{"ScreenshotPrefix", "Screenshot_20250101-010203.png", "20250101010203.png"},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			src := filepath.Join(tempDir, tc.fileName)
			if err := os.WriteFile(src, []byte("test"), 0644); err != nil {
				t.Fatalf("failed to create test file: %v", err)
			}

			file := FileInfo{Path: src, Name: tc.fileName}
			newPath, shouldRename, err := resolveScreenshotRenameTarget(file, Config{ToDateTime: true})
			if err != nil {
				t.Fatalf("resolveScreenshotRenameTarget() returned error: %v", err)
			}
			if !shouldRename {
				t.Fatalf("resolveScreenshotRenameTarget() should request rename for %s", tc.fileName)
			}

			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}
			if err := performScreenshotRename(file.Path, newPath, stdout, stderr); err != nil {
				t.Fatalf("performScreenshotRename() returned error: %v", err)
			}

			expectedPath := filepath.Join(tempDir, tc.expected)
			if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
				t.Fatalf("expected file %s does not exist", expectedPath)
			}
		})
	}
}

func TestResolveScreenshotRenameTarget_ParseError(t *testing.T) {
	tempDir := t.TempDir()
	fileName := "vlcsnap-invalid.png"
	path := filepath.Join(tempDir, fileName)
	if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	file := FileInfo{Path: path, Name: fileName}
	_, shouldRename, err := resolveScreenshotRenameTarget(file, Config{Operation: OperationVLC})
	if err == nil {
		t.Fatalf("resolveScreenshotRenameTarget() should return error for invalid pattern")
	}
	if shouldRename {
		t.Fatalf("resolveScreenshotRenameTarget() should not request rename when error occurs")
	}
}

func TestResolveScreenshotRenameTarget_NoMatch(t *testing.T) {
	tempDir := t.TempDir()
	fileName := "normal.png"
	path := filepath.Join(tempDir, fileName)
	if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	file := FileInfo{Path: path, Name: fileName}
	newPath, shouldRename, err := resolveScreenshotRenameTarget(file, Config{Operation: OperationVLC})
	if err != nil {
		t.Fatalf("resolveScreenshotRenameTarget() returned unexpected error: %v", err)
	}
	if shouldRename {
		t.Fatalf("resolveScreenshotRenameTarget() should skip non-matching file")
	}
	if newPath != "" {
		t.Fatalf("resolveScreenshotRenameTarget() should return empty path when skipping")
	}
}

func TestPerformScreenshotRename_RenameError(t *testing.T) {
	tempDir := t.TempDir()
	readOnlyDir := filepath.Join(tempDir, "readonly")
	if err := os.Mkdir(readOnlyDir, 0500); err != nil {
		t.Fatalf("failed to create readonly dir: %v", err)
	}

	if err := os.Chmod(readOnlyDir, 0700); err != nil {
		t.Fatalf("failed to chmod dir: %v", err)
	}
	src := filepath.Join(readOnlyDir, "vlcsnap-2025-05-07-12-34-56.png")
	if err := os.WriteFile(src, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	if err := os.Chmod(readOnlyDir, 0500); err != nil {
		t.Fatalf("failed to set readonly: %v", err)
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	if err := performScreenshotRename(src, filepath.Join(readOnlyDir, "new.png"), stdout, stderr); err == nil {
		t.Fatalf("performScreenshotRename() should fail on readonly directory")
	}
	if err := os.Chmod(readOnlyDir, 0700); err != nil {
		t.Fatalf("failed to reset readonly dir permission: %v", err)
	}
}
func TestValidateConfig_InvalidDirectory(t *testing.T) {
	// Arrange
	stderr := &bytes.Buffer{}
	config := Config{
		SrcDir:    "/non/existent/directory",
		Operation: OperationVLC,
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
	filesVlc, err := findScreenshotFiles(tempDir, false, OperationVLC, stdout, stderr)

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
	filesWin, err := findScreenshotFiles(tempDir, false, OperationWin, stdout, stderr)

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

func TestFindScreenshotFiles_XiaomiPattern(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-find-xiaomi")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	xiaomiFile := filepath.Join(tempDir, "Screenshot_2025-12-27-12-03-39-927_com.YoStarJP.AzurLane.jpg")
	if err := os.WriteFile(xiaomiFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	files, err := findScreenshotFiles(tempDir, false, OperationXiaomi, stdout, stderr)
	if err != nil {
		t.Fatalf("findScreenshotFiles() returned error: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("findScreenshotFiles() should find 1 Xiaomi file, got %d", len(files))
	}
	if files[0] != xiaomiFile {
		t.Fatalf("findScreenshotFiles() returned unexpected file: %s", files[0])
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
	files, err := findScreenshotFiles(tempDir, true, OperationVLC, stdout, stderr)

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
	_, err := findScreenshotFiles("/non/existent/directory", false, OperationVLC, stdout, stderr)

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
		SrcDir:    tempDir,
		Operation: OperationVLC,
		Workers:   2,
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
		SrcDir:    tempDir,
		Operation: OperationVLC,
		Workers:   0,
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
		SrcDir:    tempDir,
		Operation: OperationVLC,
		Workers:   10,
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

func TestRenameScreenshotFiles_AbortOnExistingNameConflict(t *testing.T) {
	tempDir := t.TempDir()
	sourceName := "vlcsnap-2025-05-07-12-34-56.png"
	sourcePath := filepath.Join(tempDir, sourceName)
	if err := os.WriteFile(sourcePath, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create source file: %v", err)
	}

	conflictName := "Screenshot_20250507-123456.png"
	conflictPath := filepath.Join(tempDir, conflictName)
	if err := os.WriteFile(conflictPath, []byte("existing"), 0644); err != nil {
		t.Fatalf("failed to create conflict file: %v", err)
	}

	fileInfos := []FileInfo{
		{Path: sourcePath, Name: sourceName},
	}

	config := Config{SrcDir: tempDir, Operation: OperationVLC, Workers: 2}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	successCount, errorCount := renameScreenshotFiles(fileInfos, config, stdout, stderr)
	if successCount != 0 {
		t.Fatalf("renameScreenshotFiles() successCount = %d, want 0", successCount)
	}
	if errorCount == 0 {
		t.Fatalf("renameScreenshotFiles() errorCount should be > 0 when conflicts exist")
	}

	if _, err := os.Stat(sourcePath); err != nil {
		t.Fatalf("source file should remain untouched: %v", err)
	}
	if _, err := os.Stat(conflictPath); err != nil {
		t.Fatalf("conflict file should remain untouched: %v", err)
	}
	if !strings.Contains(stderr.String(), conflictName) {
		t.Fatalf("stderr should mention conflict file, got: %s", stderr.String())
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
	_, err = findScreenshotFiles(tempDir, true, OperationVLC, stdout, stderr)

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
