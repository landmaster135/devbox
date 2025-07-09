package usecases

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// テスト用の画像ファイルを作成するヘルパー関数
func createTestImage(path string) error {
	// 10x10のテスト画像を作成
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			img.Set(x, y, color.RGBA{R: uint8(x * 25), G: uint8(y * 25), B: 128, A: 255})
		}
	}

	// ディレクトリを作成
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	// ファイルを作成
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	// PNGとして保存
	return png.Encode(file, img)
}

func TestIsImage(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "jpg extension",
			path:     "test.jpg",
			expected: true,
		},
		{
			name:     "jpeg extension",
			path:     "test.jpeg",
			expected: true,
		},
		{
			name:     "png extension",
			path:     "test.png",
			expected: true,
		},
		{
			name:     "uppercase JPG",
			path:     "test.JPG",
			expected: true,
		},
		{
			name:     "gif extension",
			path:     "test.gif",
			expected: false,
		},
		{
			name:     "no extension",
			path:     "test",
			expected: false,
		},
		{
			name:     "txt extension",
			path:     "test.txt",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsImage(tt.path)
			if got != tt.expected {
				t.Errorf("IsImage(%s) = %v, want %v", tt.path, got, tt.expected)
			}
		})
	}
}

func TestRotateAndSave(t *testing.T) {
	// テンポラリディレクトリを作成
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "input.png")
	outputDir := filepath.Join(tempDir, "output")

	// テスト用の画像を作成
	if err := createTestImage(inputPath); err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	// RotateAndSaveをテスト
	err := RotateAndSave(inputPath, outputDir, "rotated", 90.0)
	if err != nil {
		t.Errorf("RotateAndSave() error = %v", err)
		return
	}

	// 出力ファイルが存在することを確認
	outputPath := filepath.Join(outputDir, "input_rotated.png")
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Output file does not exist: %s", outputPath)
	}
}

func TestRotateAndSave_InvalidImage(t *testing.T) {
	tempDir := t.TempDir()
	invalidPath := filepath.Join(tempDir, "nonexistent.png")
	outputDir := filepath.Join(tempDir, "output")

	err := RotateAndSave(invalidPath, outputDir, "rotated", 90.0)
	if err == nil {
		t.Error("RotateAndSave() expected error for invalid image, got nil")
	}
}

func TestMoveFile(t *testing.T) {
	// テンポラリディレクトリを作成
	tempDir := t.TempDir()
	sourceDir := filepath.Join(tempDir, "source")
	destDir := filepath.Join(tempDir, "dest")

	// ソースファイルを作成
	srcFile := filepath.Join(sourceDir, "test.txt")
	if err := os.MkdirAll(sourceDir, 0o755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}
	
	file, err := os.Create(srcFile)
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	file.Close()

	// ファイルを移動
	err = MoveFile(srcFile, destDir)
	if err != nil {
		t.Errorf("MoveFile() error = %v", err)
		return
	}

	// 移動先にファイルが存在することを確認
	destFile := filepath.Join(destDir, "test.txt")
	if _, err := os.Stat(destFile); os.IsNotExist(err) {
		t.Errorf("Destination file does not exist: %s", destFile)
	}

	// 移動元にファイルが存在しないことを確認
	if _, err := os.Stat(srcFile); !os.IsNotExist(err) {
		t.Errorf("Source file still exists: %s", srcFile)
	}
}

func TestMoveFile_NonExistentFile(t *testing.T) {
	tempDir := t.TempDir()
	nonExistentFile := filepath.Join(tempDir, "nonexistent.txt")
	destDir := filepath.Join(tempDir, "dest")

	err := MoveFile(nonExistentFile, destDir)
	if err == nil {
		t.Error("MoveFile() expected error for non-existent file, got nil")
	}
}

func TestEnsureDir(t *testing.T) {
	tempDir := t.TempDir()
	testDir := filepath.Join(tempDir, "test", "nested", "directory")

	err := ensureDir(testDir)
	if err != nil {
		t.Errorf("ensureDir() error = %v", err)
		return
	}

	// ディレクトリが存在することを確認
	info, err := os.Stat(testDir)
	if os.IsNotExist(err) {
		t.Errorf("Directory does not exist: %s", testDir)
		return
	}

	if !info.IsDir() {
		t.Errorf("Path exists but is not a directory: %s", testDir)
	}

	// 権限を確認
	mode := info.Mode()
	expectedMode := os.FileMode(0o755)
	if mode.Perm() != expectedMode {
		t.Errorf("Directory has incorrect permissions: got %v, want %v", mode.Perm(), expectedMode)
	}
}

func TestEnsureDir_ExistingDirectory(t *testing.T) {
	tempDir := t.TempDir()
	
	// ディレクトリを2回作成しても問題ないことを確認
	err := ensureDir(tempDir)
	if err != nil {
		t.Errorf("ensureDir() error on existing directory = %v", err)
	}
}
