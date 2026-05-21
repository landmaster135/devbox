package usecases

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	commandExecutor "github.com/landmaster135/devbox/internal/movie_extractor/infrastructures/command_executor"
)

func TestHandleDedupImages_ExampleCase(t *testing.T) {
	srcDir := t.TempDir()
	outDir := filepath.Join(t.TempDir(), "unique")

	group1 := [][]uint8{{10, 10}, {10, 10}}
	unique := [][]uint8{{20, 20}, {20, 20}}
	group2 := [][]uint8{{30, 30}, {30, 30}}

	mustWriteJPEG(t, filepath.Join(srcDir, "img01.jpg"), group1)
	mustWriteJPEG(t, filepath.Join(srcDir, "img02.jpg"), group1)
	mustWriteJPEG(t, filepath.Join(srcDir, "img03.jpg"), group1)
	mustWriteJPEG(t, filepath.Join(srcDir, "img11.jpg"), unique)
	mustWriteJPEG(t, filepath.Join(srcDir, "img21.jpg"), group2)
	mustWriteJPEG(t, filepath.Join(srcDir, "img22.jpg"), group2)

	service := NewServiceWithExecutor(&commandExecutor.MockRepository{})
	result, err := service.HandleDedupImages(DedupImagesInput{
		SrcDir:    srcDir,
		MatchRate: 100,
		OutDir:    outDir,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	names := mustReadFileNames(t, outDir)
	expected := []string{"img01.jpg", "img11.jpg", "img21.jpg"}
	if !reflect.DeepEqual(names, expected) {
		t.Fatalf("unexpected output files: got=%v want=%v", names, expected)
	}
	if !strings.Contains(result, "出力画像数: 3") {
		t.Fatalf("unexpected result: %s", result)
	}
}

func TestHandleDedupImages_MatchRateThreshold(t *testing.T) {
	srcDir := t.TempDir()
	outDirHigh := filepath.Join(t.TempDir(), "high")
	outDirLow := filepath.Join(t.TempDir(), "low")

	base := [][]uint8{{0, 0}, {0, 0}}
	partial := [][]uint8{{0, 0}, {0, 255}} // 一致率 75%

	mustWritePNG(t, filepath.Join(srcDir, "a.png"), base)
	mustWritePNG(t, filepath.Join(srcDir, "b.png"), partial)

	service := NewServiceWithExecutor(&commandExecutor.MockRepository{})
	if _, err := service.HandleDedupImages(DedupImagesInput{
		SrcDir:    srcDir,
		MatchRate: 80,
		OutDir:    outDirHigh,
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	highNames := mustReadFileNames(t, outDirHigh)
	if !reflect.DeepEqual(highNames, []string{"a.png", "b.png"}) {
		t.Fatalf("unexpected output with high threshold: %v", highNames)
	}

	if _, err := service.HandleDedupImages(DedupImagesInput{
		SrcDir:    srcDir,
		MatchRate: 70,
		OutDir:    outDirLow,
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lowNames := mustReadFileNames(t, outDirLow)
	if !reflect.DeepEqual(lowNames, []string{"a.png"}) {
		t.Fatalf("unexpected output with low threshold: %v", lowNames)
	}
}

func TestHandleDedupImages_NoImages(t *testing.T) {
	service := NewServiceWithExecutor(&commandExecutor.MockRepository{})
	_, err := service.HandleDedupImages(DedupImagesInput{
		SrcDir:    t.TempDir(),
		MatchRate: 90,
		OutDir:    t.TempDir(),
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "src-dir に対象画像が存在しません") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHandleDedupImages_DecodeError(t *testing.T) {
	srcDir := t.TempDir()
	outDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(srcDir, "broken.png"), []byte("not-image"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	service := NewServiceWithExecutor(&commandExecutor.MockRepository{})
	_, err := service.HandleDedupImages(DedupImagesInput{
		SrcDir:    srcDir,
		MatchRate: 90,
		OutDir:    outDir,
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "画像デコードに失敗しました") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHandleDedupImages_CompareOnlyRecentSelectedImage(t *testing.T) {
	srcDir := t.TempDir()
	outDir := filepath.Join(t.TempDir(), "unique")

	imageA := [][]uint8{{5, 5}, {5, 5}}
	imageB := [][]uint8{{20, 20}, {20, 20}}

	mustWritePNG(t, filepath.Join(srcDir, "img01.png"), imageA)
	mustWritePNG(t, filepath.Join(srcDir, "img02.png"), imageB)
	mustWritePNG(t, filepath.Join(srcDir, "img03.png"), imageA)

	service := NewServiceWithExecutor(&commandExecutor.MockRepository{})
	_, err := service.HandleDedupImages(DedupImagesInput{
		SrcDir:    srcDir,
		MatchRate: 100,
		OutDir:    outDir,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	names := mustReadFileNames(t, outDir)
	expected := []string{"img01.png", "img02.png", "img03.png"}
	if !reflect.DeepEqual(names, expected) {
		t.Fatalf("unexpected output files: got=%v want=%v", names, expected)
	}
}

func TestHandleDedupImages_LogOutput(t *testing.T) {
	srcDir := t.TempDir()
	outDir := filepath.Join(t.TempDir(), "unique")

	imageA := [][]uint8{{7, 7}, {7, 7}}
	imageB := [][]uint8{{7, 7}, {7, 8}}

	mustWritePNG(t, filepath.Join(srcDir, "img01.png"), imageA)
	mustWritePNG(t, filepath.Join(srcDir, "img02.png"), imageB)

	service := NewServiceWithExecutor(&commandExecutor.MockRepository{})
	logBuffer := &bytes.Buffer{}
	result, err := service.HandleDedupImages(DedupImagesInput{
		SrcDir:    srcDir,
		MatchRate: 100,
		Log:       true,
		LogWriter: logBuffer,
		OutDir:    outDir,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(logBuffer.String(), `照合率: "img02.png" vs "img01.png": 100.00%`) {
		t.Fatalf("unexpected log output: %s", logBuffer.String())
	}
	if strings.Contains(result, "照合率:") {
		t.Fatalf("unexpected result: %s", result)
	}
}

func mustWritePNG(t *testing.T, path string, pixels [][]uint8) {
	t.Helper()
	img := toGrayImage(t, pixels)
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create file: %v", err)
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		t.Fatalf("failed to encode png: %v", err)
	}
}

func mustWriteJPEG(t *testing.T, path string, pixels [][]uint8) {
	t.Helper()
	img := toGrayImage(t, pixels)
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create file: %v", err)
	}
	defer file.Close()
	if err := jpeg.Encode(file, img, &jpeg.Options{Quality: 100}); err != nil {
		t.Fatalf("failed to encode jpeg: %v", err)
	}
}

func toGrayImage(t *testing.T, pixels [][]uint8) *image.Gray {
	t.Helper()
	if len(pixels) == 0 || len(pixels[0]) == 0 {
		t.Fatal("pixels must not be empty")
	}
	height := len(pixels)
	width := len(pixels[0])
	img := image.NewGray(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		if len(pixels[y]) != width {
			t.Fatal("pixels width mismatch")
		}
		for x := 0; x < width; x++ {
			img.SetGray(x, y, color.Gray{Y: pixels[y][x]})
		}
	}
	return img
}

func mustReadFileNames(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("failed to read dir: %v", err)
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		names = append(names, entry.Name())
	}
	sort.Strings(names)
	return names
}
