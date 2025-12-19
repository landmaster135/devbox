package usecases

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/landmaster135/devbox/internal/image_filterer_v2/config"
)

func TestServiceProcessGrayscale(t *testing.T) {
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "input.png")
	outputPath := filepath.Join(tempDir, "output.png")

	src := image.NewNRGBA(image.Rect(0, 0, 2, 1))
	src.Set(0, 0, color.NRGBA{R: 255, G: 0, B: 0, A: 255})
	src.Set(1, 0, color.NRGBA{R: 0, G: 0, B: 255, A: 255})

	if err := savePNG(inputPath, src); err != nil {
		t.Fatalf("failed to save input: %v", err)
	}

	cfg := config.Config{
		InputPath:  inputPath,
		OutputPath: outputPath,
		Mode:       config.FilterModeGrayscale,
		Strength:   1.0,
	}

	svc, err := NewService(cfg)
	if err != nil {
		t.Fatalf("failed to build service: %v", err)
	}

	outPath, err := svc.Process()
	if err != nil {
		t.Fatalf("failed to process image: %v", err)
	}
	if outPath != outputPath {
		t.Fatalf("unexpected output path: got %s want %s", outPath, outputPath)
	}

	outImg, err := loadImage(outPath)
	if err != nil {
		t.Fatalf("failed to load output: %v", err)
	}

	gray0 := color.NRGBAModel.Convert(outImg.At(0, 0)).(color.NRGBA)
	if !approxEqual(int(gray0.R), 76, 3) {
		t.Fatalf("unexpected grayscale value (pixel 0): %d", gray0.R)
	}

	gray1 := color.NRGBAModel.Convert(outImg.At(1, 0)).(color.NRGBA)
	if !approxEqual(int(gray1.R), 29, 3) {
		t.Fatalf("unexpected grayscale value (pixel 1): %d", gray1.R)
	}
}

func TestNewServiceUnsupportedMode(t *testing.T) {
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "input.png")
	if err := savePNG(inputPath, image.NewNRGBA(image.Rect(0, 0, 1, 1))); err != nil {
		t.Fatalf("failed to save input: %v", err)
	}

	cfg := config.Config{
		InputPath: inputPath,
		Mode:      config.FilterMode("unknown"),
		Strength:  1.0,
	}

	if _, err := NewService(cfg); err == nil {
		t.Fatalf("expected error for unsupported mode")
	}
}

func approxEqual(value int, expected int, tol int) bool {
	return value >= expected-tol && value <= expected+tol
}

func savePNG(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return err
	}
	return f.Close()
}
