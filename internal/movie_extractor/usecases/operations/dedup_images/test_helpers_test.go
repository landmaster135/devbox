package dedupimages

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"sort"
	"testing"
)

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

func mustWriteLineImagePNG(t *testing.T, path string, width int, height int, lineX int) {
	t.Helper()
	img := image.NewGray(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.SetGray(x, y, color.Gray{Y: 255})
		}
	}
	for y := 0; y < height; y++ {
		for dx := -1; dx <= 1; dx++ {
			x := lineX + dx
			if x >= 0 && x < width {
				img.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}

	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create file: %v", err)
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		t.Fatalf("failed to encode png: %v", err)
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
