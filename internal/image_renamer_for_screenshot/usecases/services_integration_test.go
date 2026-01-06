package usecases

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProcessScreenshotRename_Normal(t *testing.T) {
	// Arrange
	tempDir, err := os.MkdirTemp("", "test-process-screenshot-rename")
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

	config := Config{
		SrcDir:     tempDir,
		VlcPattern: true,
		WinPattern: false,
		Workers:    2,
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	// Act
	successCount, errorCount, err := ProcessScreenshotRename(config, stdout, stderr)

	// Assert
	if err != nil {
		t.Errorf("ProcessScreenshotRename() returned error: %v", err)
	}
	if successCount != 1 {
		t.Errorf("ProcessScreenshotRename() successCount = %v, want %v", successCount, 1)
	}
	if errorCount != 0 {
		t.Errorf("ProcessScreenshotRename() errorCount = %v, want %v", errorCount, 0)
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

func TestProcessScreenshotRename_ToDateTime(t *testing.T) {
	// Arrange
	tempDir, err := os.MkdirTemp("", "test-process-screenshot-rename-datetime")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// VLCスクリーンショットファイルを作成
	vlcFile := filepath.Join(tempDir, "vlcsnap-2025-05-07-12-34-56.png")
	if err := os.WriteFile(vlcFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Screenshot_ファイルを作成
	screenshotFile := filepath.Join(tempDir, "Screenshot_20250507-123456.png")
	if err := os.WriteFile(screenshotFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	config := Config{
		SrcDir:     tempDir,
		ToDateTime: true,
		Workers:    2,
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	// Act
	successCount, errorCount, err := ProcessScreenshotRename(config, stdout, stderr)

	// Assert
	if err != nil {
		t.Errorf("ProcessScreenshotRename() returned error: %v", err)
	}
	if successCount != 2 {
		t.Errorf("ProcessScreenshotRename() successCount = %v, want %v", successCount, 2)
	}
	if errorCount != 0 {
		t.Errorf("ProcessScreenshotRename() errorCount = %v, want %v", errorCount, 0)
	}

	// VLCファイルがYYYYMMDDHHMMSS形式にリネームされたことを確認
	expectedVlcPath := filepath.Join(tempDir, "20250507123456.png")
	if _, err := os.Stat(expectedVlcPath); os.IsNotExist(err) {
		t.Errorf("VLC file was not renamed to %s", expectedVlcPath)
	}

	// Screenshot_ファイルがYYYYMMDDHHMMSS形式にリネームされたことを確認
	expectedScreenshotPath := filepath.Join(tempDir, "20250507123456.png")
	if _, err := os.Stat(expectedScreenshotPath); os.IsNotExist(err) {
		t.Errorf("Screenshot_ file was not renamed to %s", expectedScreenshotPath)
	}
}

func TestProcessScreenshotRename_NoFiles(t *testing.T) {
	// Arrange
	tempDir, err := os.MkdirTemp("", "test-process-screenshot-rename-nofiles")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 通常のファイルのみ作成（スクリーンショットファイルではない）
	normalFile := filepath.Join(tempDir, "normal.png")
	if err := os.WriteFile(normalFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	config := Config{
		SrcDir:     tempDir,
		VlcPattern: true,
		Workers:    1,
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	// Act
	successCount, errorCount, err := ProcessScreenshotRename(config, stdout, stderr)

	// Assert
	if err != nil {
		t.Errorf("ProcessScreenshotRename() returned error: %v", err)
	}
	if successCount != 0 {
		t.Errorf("ProcessScreenshotRename() successCount = %v, want %v", successCount, 0)
	}
	if errorCount != 0 {
		t.Errorf("ProcessScreenshotRename() errorCount = %v, want %v", errorCount, 0)
	}

	// 出力メッセージを確認
	output := stdout.String()
	if !strings.Contains(output, "スクリーンショットファイルが見つかりませんでした") {
		t.Errorf("Expected 'no files found' message in output: %s", output)
	}
}

func TestProcessScreenshotRename_InvalidConfig(t *testing.T) {
	// Arrange
	config := Config{
		SrcDir:     "/non/existent/directory",
		VlcPattern: true,
		Workers:    1,
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	// Act
	successCount, errorCount, err := ProcessScreenshotRename(config, stdout, stderr)

	// Assert
	if err == nil {
		t.Error("ProcessScreenshotRename() should return error for invalid config")
	}
	if successCount != 0 {
		t.Errorf("ProcessScreenshotRename() successCount = %v, want %v", successCount, 0)
	}
	if errorCount != 0 {
		t.Errorf("ProcessScreenshotRename() errorCount = %v, want %v", errorCount, 0)
	}
}

func TestProcessScreenshotRename_WithErrors(t *testing.T) {
	// Arrange
	tempDir, err := os.MkdirTemp("", "test-process-screenshot-rename-errors")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 正常なVLCスクリーンショットファイルを作成
	vlcFile := filepath.Join(tempDir, "vlcsnap-2025-05-07-12-34-56.png")
	if err := os.WriteFile(vlcFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// 不正なフォーマットのVLCスクリーンショットファイルを作成
	invalidVlcFile := filepath.Join(tempDir, "vlcsnap-invalid.png")
	if err := os.WriteFile(invalidVlcFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	config := Config{
		SrcDir:     tempDir,
		VlcPattern: true,
		Workers:    2,
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	// Act
	successCount, errorCount, err := ProcessScreenshotRename(config, stdout, stderr)

	// Assert
	if err == nil {
		t.Error("ProcessScreenshotRename() should return error when there are processing errors")
	}
	if successCount != 1 {
		t.Errorf("ProcessScreenshotRename() successCount = %v, want %v", successCount, 1)
	}
	if errorCount != 1 {
		t.Errorf("ProcessScreenshotRename() errorCount = %v, want %v", errorCount, 1)
	}

	// 正常なファイルがリネームされたことを確認
	expectedNewPath := filepath.Join(tempDir, "Screenshot_20250507-123456.png")
	if _, err := os.Stat(expectedNewPath); os.IsNotExist(err) {
		t.Errorf("Valid file was not renamed to %s", expectedNewPath)
	}

	// 不正なファイルはリネームされていないことを確認
	if _, err := os.Stat(invalidVlcFile); os.IsNotExist(err) {
		t.Errorf("Invalid file should not be renamed")
	}
}

func TestProcessScreenshotRename_MultiplePatterns(t *testing.T) {
	// Arrange
	tempDir, err := os.MkdirTemp("", "test-process-screenshot-rename-multiple")
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

	// Pixelスクリーンレコードファイルを作成
	pixelFile := filepath.Join(tempDir, "screen-20250507-123456.mp4")
	if err := os.WriteFile(pixelFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	config := Config{
		SrcDir:       tempDir,
		VlcPattern:   true,
		WinPattern:   true,
		PixelPattern: false,
		Workers:      3,
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	// Act - 複数パターンが指定されているため設定エラーになるはず
	successCount, errorCount, err := ProcessScreenshotRename(config, stdout, stderr)

	// Assert
	if err == nil {
		t.Error("ProcessScreenshotRename() should return error for multiple patterns")
	}
	if successCount != 0 {
		t.Errorf("ProcessScreenshotRename() successCount = %v, want %v", successCount, 0)
	}
	if errorCount != 0 {
		t.Errorf("ProcessScreenshotRename() errorCount = %v, want %v", errorCount, 0)
	}
}
