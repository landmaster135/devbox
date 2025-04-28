// main_test.go
// Unit tests for convertFile logic.
// These tests use stub codecs so that no real image encoding/decoding libraries
// are required. Each test creates a temporary input file, invokes convertFile,
// and asserts that the correct output file (bytes & extension) is produced.
package usecases

import (
	"bytes"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"
)

// stubDecode returns a 1×1 pixel RGBA image (actual content is irrelevant).
func stubDecode(_ []byte) (image.Image, error) {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{255, 0, 0, 255})
	return img, nil
}

// makeStubCodecTable builds a codec table where the Encode function always
// returns a byte slice of the specified size (so we can control “compressed” size).
func makeStubCodecTable(outSize int) map[string]codec {
	stubEncode := func(_ image.Image) ([]byte, error) {
		return bytes.Repeat([]byte{0x42}, outSize), nil
	}
	c := codec{Decode: stubDecode, Encode: stubEncode}
	return map[string]codec{
		"jpg":  c,
		"png":  c,
		"webp": c,
		"avif": c,
		"svg":  c,
	}
}

// helper writes a dummy file of given size and returns its path.
func writeDummyFile(dir, name string, size int) (string, error) {
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, bytes.Repeat([]byte{0x24}, size), 0o644); err != nil {
		return "", err
	}
	return path, nil
}

func TestConvertKeepsOriginalWhenLarger(t *testing.T) {
	srcDir := t.TempDir()
	outDir := t.TempDir()

	origPath, err := writeDummyFile(srcDir, "photo.jpg", 100) // 100‑byte source
	if err != nil {
		t.Fatal(err)
	}

	codecs := makeStubCodecTable(200) // Encode yields 200‑byte output (larger)

	if err := ConvertFile(origPath, srcDir, outDir, "png", codecs); err != nil {
		t.Fatalf("convertFile failed: %v", err)
	}

	// Expect: kept original -> output file has .jpg and 100 bytes
	outPath := filepath.Join(outDir, "photo.jpg")
	info, err := os.Stat(outPath)
	if err != nil {
		t.Fatalf("expected output file %s: %v", outPath, err)
	}
	if size := info.Size(); size != 100 {
		t.Errorf("expected 100 bytes, got %d", size)
	}
}

func TestConvertConvertedWhenSmaller(t *testing.T) {
	srcDir := t.TempDir()
	outDir := t.TempDir()

	origPath, _ := writeDummyFile(srcDir, "flower.jpg", 200) // 200‑byte source

	codecs := makeStubCodecTable(50) // Encode yields 50‑byte output (smaller)

	if err := ConvertFile(origPath, srcDir, outDir, "webp", codecs); err != nil {
		t.Fatalf("convertFile failed: %v", err)
	}

	// Expect: converted -> output file has .webp and 50 bytes
	outPath := filepath.Join(outDir, "flower.webp")
	info, err := os.Stat(outPath)
	if err != nil {
		t.Fatalf("expected output file %s: %v", outPath, err)
	}
	if size := info.Size(); size != 50 {
		t.Errorf("expected 50 bytes, got %d", size)
	}
}

func TestConvertSVGForceConvert(t *testing.T) {
	srcDir := t.TempDir()
	outDir := t.TempDir()

	origPath, _ := writeDummyFile(srcDir, "icon.svg", 80) // 80‑byte SVG

	codecs := makeStubCodecTable(300) // Encode yields larger size

	// Rule: SVG input → non‑SVG output must convert regardless of size
	if err := ConvertFile(origPath, srcDir, outDir, "png", codecs); err != nil {
		t.Fatalf("convertFile failed: %v", err)
	}

	// Expect: converted -> .png exists (even though it's larger)
	outPath := filepath.Join(outDir, "icon.png")
	info, err := os.Stat(outPath)
	if err != nil {
		t.Fatalf("expected output file %s: %v", outPath, err)
	}
	if size := info.Size(); size != 300 {
		t.Errorf("expected 300 bytes, got %d", size)
	}
}
