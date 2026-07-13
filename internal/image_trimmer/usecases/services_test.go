package usecases

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anthonynsimon/bild/imgio"
)

// テスト用の画像を作成するヘルパー関数
func createTestImage(width, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// 赤い背景
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}
	return img
}

func TestCropAndSave(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tmpDir := t.TempDir()
	testInputDir := filepath.Join(tmpDir, "input")
	testOutputDir := filepath.Join(tmpDir, "output")

	// テスト用入力ディレクトリを作成
	if err := os.MkdirAll(testInputDir, 0o755); err != nil {
		t.Fatalf("failed to create test input dir: %v", err)
	}

	// テスト画像を作成して保存
	testImg := createTestImage(100, 100)
	testInputPath := filepath.Join(testInputDir, "test.png")
	if err := imgio.Save(testInputPath, testImg, imgio.PNGEncoder()); err != nil {
		t.Fatalf("failed to save test image: %v", err)
	}

	tests := []struct {
		name        string
		inputPath   string
		outputDir   string
		x1, y1      int
		x2, y2      int
		suffix      string
		expectError bool
		errContains string
	}{
		{
			name:        "正常系: 有効な矩形でトリミング",
			inputPath:   testInputPath,
			outputDir:   testOutputDir,
			x1:          10,
			y1:          10,
			x2:          50,
			y2:          50,
			suffix:      "cropped",
			expectError: false,
		},
		{
			name:        "異常系: 負の座標",
			inputPath:   testInputPath,
			outputDir:   testOutputDir,
			x1:          -1,
			y1:          10,
			x2:          50,
			y2:          50,
			suffix:      "cropped",
			expectError: true,
			errContains: "invalid rectangle",
		},
		{
			name:        "異常系: x2 < x1",
			inputPath:   testInputPath,
			outputDir:   testOutputDir,
			x1:          50,
			y1:          10,
			x2:          10,
			y2:          50,
			suffix:      "cropped",
			expectError: true,
			errContains: "invalid rectangle",
		},
		{
			name:        "異常系: 画像範囲外",
			inputPath:   testInputPath,
			outputDir:   testOutputDir,
			x1:          10,
			y1:          10,
			x2:          150,
			y2:          50,
			suffix:      "cropped",
			expectError: true,
			errContains: "invalid rectangle",
		},
		{
			name:        "異常系: 存在しないファイル",
			inputPath:   filepath.Join(testInputDir, "nonexistent.png"),
			outputDir:   testOutputDir,
			x1:          10,
			y1:          10,
			x2:          50,
			y2:          50,
			suffix:      "cropped",
			expectError: true,
			errContains: "open",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CropAndSave(tt.inputPath, tt.outputDir, tt.x1, tt.y1, tt.x2, tt.y2, tt.suffix)

			if tt.expectError {
				if err == nil {
					t.Errorf("CropAndSave() expected error but got nil")
				} else if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("CropAndSave() error = %v, want error containing %s", err, tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("CropAndSave() unexpected error = %v", err)
				}

				// 出力ファイルが存在するか確認
				expectedOutputPath := filepath.Join(tt.outputDir, "test_cropped.png")
				if _, err := os.Stat(expectedOutputPath); os.IsNotExist(err) {
					t.Errorf("CropAndSave() output file not created: %s", expectedOutputPath)
				}

				// 出力画像を読み込んで確認
				outputImg, err := imgio.Open(expectedOutputPath)
				if err != nil {
					t.Errorf("Failed to open output image: %v", err)
				}

				expectedWidth := tt.x2 - tt.x1
				expectedHeight := tt.y2 - tt.y1
				if bounds := outputImg.Bounds(); bounds.Dx() != expectedWidth || bounds.Dy() != expectedHeight {
					t.Errorf("CropAndSave() output image size = %dx%d, want %dx%d",
						bounds.Dx(), bounds.Dy(), expectedWidth, expectedHeight)
				}
			}
		})
	}
}

// JPEGテストを別途実施（PNGとは異なるエンコード処理のため）
func TestCropAndSaveJPEG(t *testing.T) {
	tmpDir := t.TempDir()
	testInputDir := filepath.Join(tmpDir, "input")
	testOutputDir := filepath.Join(tmpDir, "output")

	if err := os.MkdirAll(testInputDir, 0o755); err != nil {
		t.Fatalf("failed to create test input dir: %v", err)
	}

	// JPEG形式のテスト画像
	testImg := createTestImage(100, 100)
	testInputPath := filepath.Join(testInputDir, "test.jpg")
	if err := imgio.Save(testInputPath, testImg, imgio.JPEGEncoder(95)); err != nil {
		t.Fatalf("failed to save test JPEG image: %v", err)
	}

	err := CropAndSave(testInputPath, testOutputDir, 10, 10, 50, 50, "cropped")
	if err != nil {
		t.Errorf("CropAndSave() unexpected error = %v", err)
	}

	// 出力ファイルが存在するか確認
	expectedOutputPath := filepath.Join(testOutputDir, "test_cropped.jpg")
	if _, err := os.Stat(expectedOutputPath); os.IsNotExist(err) {
		t.Errorf("CropAndSave() JPEG output file not created: %s", expectedOutputPath)
	}
}

func TestCropAndSaveUnsupportedFormat(t *testing.T) {
	tmpDir := t.TempDir()
	testInputDir := filepath.Join(tmpDir, "input")
	testOutputDir := filepath.Join(tmpDir, "output")

	if err := os.MkdirAll(testInputDir, 0o755); err != nil {
		t.Fatalf("failed to create test input dir: %v", err)
	}

	// サポートされていない形式（.bmp）
	testImg := createTestImage(100, 100)
	testInputPath := filepath.Join(testInputDir, "test.bmp")
	if err := imgio.Save(testInputPath, testImg, imgio.PNGEncoder()); err != nil {
		t.Fatalf("failed to save test image: %v", err)
	}

	err := CropAndSave(testInputPath, testOutputDir, 10, 10, 50, 50, "cropped")
	if err == nil || !strings.Contains(err.Error(), "unsupported extension") {
		t.Errorf("CropAndSave() expected unsupported extension error, got %v", err)
	}
}

func TestMoveOriginal(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tmpDir := t.TempDir()
	testSourceDir := filepath.Join(tmpDir, "source")
	testArchiveDir := filepath.Join(tmpDir, "archive")

	// ソースディレクトリを作成
	if err := os.MkdirAll(testSourceDir, 0o755); err != nil {
		t.Fatalf("failed to create test source dir: %v", err)
	}

	// テストファイルを作成
	testFile := filepath.Join(testSourceDir, "test.png")
	content := []byte("test content")
	if err := os.WriteFile(testFile, content, 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	tests := []struct {
		name        string
		sourceFile  string
		archiveDir  string
		expectError bool
		errContains string
		setupErr    bool
	}{
		{
			name:        "正常系: ファイルを正常に移動",
			sourceFile:  testFile,
			archiveDir:  testArchiveDir,
			expectError: false,
		},
		{
			name:        "異常系: 存在しないファイル",
			sourceFile:  filepath.Join(testSourceDir, "nonexistent.png"),
			archiveDir:  testArchiveDir,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MoveOriginal(tt.sourceFile, tt.archiveDir)

			if tt.expectError {
				if err == nil {
					t.Errorf("MoveOriginal() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("MoveOriginal() unexpected error = %v", err)
				}

				// 移動先にファイルが存在するか確認
				expectedPath := filepath.Join(tt.archiveDir, filepath.Base(tt.sourceFile))
				if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
					t.Errorf("MoveOriginal() file not found at %s", expectedPath)
				}

				// 元の場所にファイルがないか確認
				if _, err := os.Stat(tt.sourceFile); !os.IsNotExist(err) {
					t.Errorf("MoveOriginal() source file still exists at %s", tt.sourceFile)
				}
			}
		})
	}
}

// MoveOriginalでディレクトリが作成されることをテスト
func TestMoveOriginalCreateDir(t *testing.T) {
	tmpDir := t.TempDir()
	testSourceDir := filepath.Join(tmpDir, "source")
	testArchiveDir := filepath.Join(tmpDir, "archive", "subdir") // 存在しないネストディレクトリ

	// ソースディレクトリとファイルを作成
	if err := os.MkdirAll(testSourceDir, 0o755); err != nil {
		t.Fatalf("failed to create test source dir: %v", err)
	}

	testFile := filepath.Join(testSourceDir, "test.png")
	if err := os.WriteFile(testFile, []byte("test"), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	err := MoveOriginal(testFile, testArchiveDir)
	if err != nil {
		t.Errorf("MoveOriginal() unexpected error = %v", err)
	}

	// ネストディレクトリが作成されているか確認
	if _, err := os.Stat(testArchiveDir); os.IsNotExist(err) {
		t.Errorf("MoveOriginal() did not create archive directory: %s", testArchiveDir)
	}
}
