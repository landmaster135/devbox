package usecases

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/landmaster135/devbox/internal/image_converter_by_libwebp/config"
	"github.com/landmaster135/devbox/internal/image_converter_by_libwebp/infrastructures/libwebp"
)

func TestService_Convert_Normal(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	outDir := filepath.Join(tempDir, "out")
	createTestFile(t, filepath.Join(srcDir, "a.jpg"))
	createTestFile(t, filepath.Join(srcDir, "b.png"))
	createTestFile(t, filepath.Join(srcDir, "skip.txt"))

	mockConverter := &libwebp.MockConverter{}
	service := NewServiceWithConverter(mockConverter)
	cfg := newTestConfig(t, srcDir, outDir, "", false, false, 1)

	result, err := service.Convert(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Convert() error = %v", err)
	}
	if result.SuccessCount != 2 || result.ErrorCount != 0 {
		t.Fatalf("result = %#v", result)
	}
	if mockConverter.CallCount() != 2 {
		t.Fatalf("CallCount() = %d, want 2", mockConverter.CallCount())
	}

	gotOutputs := map[string]bool{}
	for _, call := range mockConverter.Calls {
		gotOutputs[call.OutputPath] = true
		if call.Quality != 99 {
			t.Fatalf("Quality = %d, want 99", call.Quality)
		}
		if call.Lossless {
			t.Fatalf("Lossless = true, want false")
		}
	}
	if !gotOutputs[filepath.Join(outDir, "a.webp")] {
		t.Fatalf("missing output path for a.jpg: %#v", gotOutputs)
	}
	if !gotOutputs[filepath.Join(outDir, "b.webp")] {
		t.Fatalf("missing output path for b.png: %#v", gotOutputs)
	}
}

func TestService_Convert_Recursive_Normal(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	outDir := filepath.Join(tempDir, "out")
	createTestFile(t, filepath.Join(srcDir, "root.jpg"))
	createTestFile(t, filepath.Join(srcDir, "sub", "child.png"))

	mockConverter := &libwebp.MockConverter{}
	service := NewServiceWithConverter(mockConverter)
	cfg := newTestConfig(t, srcDir, outDir, "", false, true, 1)

	result, err := service.Convert(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Convert() error = %v", err)
	}
	if result.SuccessCount != 2 {
		t.Fatalf("SuccessCount = %d, want 2", result.SuccessCount)
	}

	gotOutputs := map[string]bool{}
	for _, call := range mockConverter.Calls {
		gotOutputs[call.OutputPath] = true
	}
	if !gotOutputs[filepath.Join(outDir, "root.webp")] {
		t.Fatalf("missing root output: %#v", gotOutputs)
	}
	if !gotOutputs[filepath.Join(outDir, "sub", "child.webp")] {
		t.Fatalf("missing nested output: %#v", gotOutputs)
	}
}

func TestService_Convert_NotRecursive_Normal(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	outDir := filepath.Join(tempDir, "out")
	createTestFile(t, filepath.Join(srcDir, "root.jpg"))
	createTestFile(t, filepath.Join(srcDir, "sub", "child.png"))

	mockConverter := &libwebp.MockConverter{}
	service := NewServiceWithConverter(mockConverter)
	cfg := newTestConfig(t, srcDir, outDir, "", false, false, 1)

	result, err := service.Convert(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Convert() error = %v", err)
	}
	if result.SuccessCount != 1 {
		t.Fatalf("SuccessCount = %d, want 1", result.SuccessCount)
	}
	if mockConverter.CallCount() != 1 {
		t.Fatalf("CallCount() = %d, want 1", mockConverter.CallCount())
	}
}

func TestService_Convert_CheckAvailableError(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	outDir := filepath.Join(tempDir, "out")
	createTestFile(t, filepath.Join(srcDir, "a.jpg"))

	mockConverter := &libwebp.MockConverter{
		CheckAvailableFunc: func() error {
			return errors.New("cwebp not found")
		},
	}
	service := NewServiceWithConverter(mockConverter)
	cfg := newTestConfig(t, srcDir, outDir, "", false, false, 1)

	_, err := service.Convert(context.Background(), cfg)
	if err == nil {
		t.Fatalf("Convert() error = nil")
	}
	if err.Error() != "cwebp not found" {
		t.Fatalf("Convert() error = %q", err.Error())
	}
}

func TestService_Convert_PartialError(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	outDir := filepath.Join(tempDir, "out")
	createTestFile(t, filepath.Join(srcDir, "ok.jpg"))
	createTestFile(t, filepath.Join(srcDir, "ng.jpg"))

	mockConverter := &libwebp.MockConverter{
		ConvertToWebPFunc: func(ctx context.Context, inputPath string, outputPath string, quality int, lossless bool) error {
			if strings.Contains(inputPath, "ng.jpg") {
				return errors.New("convert failed")
			}
			return nil
		},
	}
	service := NewServiceWithConverter(mockConverter)
	cfg := newTestConfig(t, srcDir, outDir, "", false, false, 1)

	result, err := service.Convert(context.Background(), cfg)
	if err == nil {
		t.Fatalf("Convert() error = nil")
	}
	if result.SuccessCount != 1 || result.ErrorCount != 1 {
		t.Fatalf("result = %#v", result)
	}
}

func TestService_Convert_ArchiveCopy_Normal(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	outDir := filepath.Join(tempDir, "out")
	archiveDir := filepath.Join(tempDir, "archive")
	srcPath := filepath.Join(srcDir, "a.jpg")
	createTestFile(t, srcPath)

	mockConverter := &libwebp.MockConverter{}
	service := NewServiceWithConverter(mockConverter)
	cfg := newTestConfig(t, srcDir, outDir, archiveDir, false, false, 1)

	result, err := service.Convert(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Convert() error = %v", err)
	}
	if result.SuccessCount != 1 {
		t.Fatalf("SuccessCount = %d, want 1", result.SuccessCount)
	}
	if _, err := os.Stat(srcPath); err != nil {
		t.Fatalf("source should remain: %v", err)
	}
	if _, err := os.Stat(filepath.Join(archiveDir, "a.jpg")); err != nil {
		t.Fatalf("archive copy missing: %v", err)
	}
}

func TestService_Convert_ArchiveMove_Normal(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	outDir := filepath.Join(tempDir, "out")
	archiveDir := filepath.Join(tempDir, "archive")
	srcPath := filepath.Join(srcDir, "a.jpg")
	createTestFile(t, srcPath)

	mockConverter := &libwebp.MockConverter{}
	service := NewServiceWithConverter(mockConverter)
	cfg := newTestConfig(t, srcDir, outDir, archiveDir, true, false, 1)

	result, err := service.Convert(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Convert() error = %v", err)
	}
	if result.SuccessCount != 1 {
		t.Fatalf("SuccessCount = %d, want 1", result.SuccessCount)
	}
	if _, err := os.Stat(srcPath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("source should be moved, stat err = %v", err)
	}
	if _, err := os.Stat(filepath.Join(archiveDir, "a.jpg")); err != nil {
		t.Fatalf("archive moved file missing: %v", err)
	}
}

func TestService_Convert_SourceDirectoryError(t *testing.T) {
	tempDir := t.TempDir()
	outDir := filepath.Join(tempDir, "out")
	mockConverter := &libwebp.MockConverter{}
	service := NewServiceWithConverter(mockConverter)
	cfg := newTestConfig(t, filepath.Join(tempDir, "missing"), outDir, "", false, false, 1)

	_, err := service.Convert(context.Background(), cfg)
	if err == nil {
		t.Fatalf("Convert() error = nil")
	}
	if !strings.Contains(err.Error(), "入力ディレクトリを確認できません") {
		t.Fatalf("Convert() error = %q", err.Error())
	}
}

func TestService_Convert_NilConfig(t *testing.T) {
	service := NewServiceWithConverter(&libwebp.MockConverter{})

	_, err := service.Convert(context.Background(), nil)
	if err == nil {
		t.Fatalf("Convert() error = nil")
	}
	if err.Error() != "設定が指定されていません" {
		t.Fatalf("Convert() error = %q", err.Error())
	}
}

func newTestConfig(t *testing.T, srcDir string, outDir string, archiveDir string, move bool, recursive bool, workers int) *config.Config {
	t.Helper()
	cfg, err := config.NewConfig(srcDir, outDir, archiveDir, "webp", move, 99, workers, recursive, false)
	if err != nil {
		t.Fatalf("NewConfig() error = %v", err)
	}
	return cfg
}

func createTestFile(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(path, []byte("test image"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
}
