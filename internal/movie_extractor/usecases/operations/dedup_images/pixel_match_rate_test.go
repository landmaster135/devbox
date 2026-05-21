package dedupimages

import (
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCalculatePixelMatchRate_GrayscaleFallback(t *testing.T) {
	dir := t.TempDir()
	basePath := filepath.Join(dir, "base.png")
	partialPath := filepath.Join(dir, "partial.png")

	base := [][]uint8{{0, 0}, {0, 0}}
	partial := [][]uint8{{0, 0}, {0, 255}} // 一致率 75%

	mustWritePNG(t, basePath, base)
	mustWritePNG(t, partialPath, partial)

	rate, err := calculatePixelMatchRate(basePath, partialPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := 0.75
	if math.Abs(rate-expected) > 0.01 {
		t.Fatalf("unexpected match rate: got=%f want=%f", rate, expected)
	}
}

func TestCalculatePixelMatchRate_DifferentDimensions(t *testing.T) {
	dir := t.TempDir()
	smallPath := filepath.Join(dir, "small.png")
	largePath := filepath.Join(dir, "large.png")

	mustWritePNG(t, smallPath, [][]uint8{{0, 0}, {0, 0}})
	mustWritePNG(t, largePath, [][]uint8{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}})

	rate, err := calculatePixelMatchRate(smallPath, largePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rate != 0 {
		t.Fatalf("unexpected match rate: got=%f want=0", rate)
	}
}

func TestCalculatePixelMatchRate_DecodeError(t *testing.T) {
	dir := t.TempDir()
	validPath := filepath.Join(dir, "valid.png")
	brokenPath := filepath.Join(dir, "broken.png")
	mustWritePNG(t, validPath, [][]uint8{{0, 0}, {0, 0}})
	if err := os.WriteFile(brokenPath, []byte("not-image"), 0644); err != nil {
		t.Fatalf("failed to write broken image: %v", err)
	}

	_, err := calculatePixelMatchRate(validPath, brokenPath)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "画像デコードに失敗しました") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCalculatePixelMatchRate_EdgeDifference(t *testing.T) {
	dir := t.TempDir()
	leftLinePath := filepath.Join(dir, "left.png")
	rightLinePath := filepath.Join(dir, "right.png")

	mustWriteLineImagePNG(t, leftLinePath, 64, 64, 14)
	mustWriteLineImagePNG(t, rightLinePath, 64, 64, 40)

	rate, err := calculatePixelMatchRate(leftLinePath, rightLinePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rate >= 0.80 {
		t.Fatalf("unexpected match rate: got=%f want<0.80", rate)
	}
}
