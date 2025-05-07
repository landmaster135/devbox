package usecases

import (
	"bytes"
	"os"
	"path/filepath"
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
	err := ValidateConfig(config, stderr)

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
	err := ValidateConfig(config, stderr)

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
	err := ValidateConfig(config, stderr)

	// Assert
	if err == nil {
		t.Error("ValidateConfig() should return error when both patterns are specified")
	}
}

func TestRenameVlcScreenshot_Normal(t *testing.T) {
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
	processScreenshotRename(file, true, false, &mu, &successCount, &errorCount, stdout, stderr)

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
	processScreenshotRename(file, false, true, &mu, &successCount, &errorCount, stdout, stderr)

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
