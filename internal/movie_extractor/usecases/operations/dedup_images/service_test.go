package dedupimages

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestHandle_ExampleCase(t *testing.T) {
	srcDir := t.TempDir()
	outDir := filepath.Join(t.TempDir(), "unique")

	group1 := [][]uint8{{10, 10}, {10, 10}}
	unique := [][]uint8{{120, 120}, {120, 120}}
	group2 := [][]uint8{{240, 240}, {240, 240}}

	mustWriteJPEG(t, filepath.Join(srcDir, "img01.jpg"), group1)
	mustWriteJPEG(t, filepath.Join(srcDir, "img02.jpg"), group1)
	mustWriteJPEG(t, filepath.Join(srcDir, "img03.jpg"), group1)
	mustWriteJPEG(t, filepath.Join(srcDir, "img11.jpg"), unique)
	mustWriteJPEG(t, filepath.Join(srcDir, "img21.jpg"), group2)
	mustWriteJPEG(t, filepath.Join(srcDir, "img22.jpg"), group2)

	service := NewService()
	result, err := service.Handle(Input{
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

func TestHandle_NoImages(t *testing.T) {
	service := NewService()
	_, err := service.Handle(Input{
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

func TestHandle_DecodeError(t *testing.T) {
	srcDir := t.TempDir()
	outDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(srcDir, "broken.png"), []byte("not-image"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	service := NewService()
	_, err := service.Handle(Input{
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

func TestHandle_CompareOnlyRecentSelectedImage(t *testing.T) {
	srcDir := t.TempDir()
	outDir := filepath.Join(t.TempDir(), "unique")

	imageA := [][]uint8{{5, 5}, {5, 5}}
	imageB := [][]uint8{{80, 80}, {80, 80}}

	mustWritePNG(t, filepath.Join(srcDir, "img01.png"), imageA)
	mustWritePNG(t, filepath.Join(srcDir, "img02.png"), imageB)
	mustWritePNG(t, filepath.Join(srcDir, "img03.png"), imageA)

	service := NewService()
	_, err := service.Handle(Input{
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

func TestHandle_LogOutput(t *testing.T) {
	srcDir := t.TempDir()
	outDir := filepath.Join(t.TempDir(), "unique")

	imageA := [][]uint8{{7, 7}, {7, 7}}
	imageB := [][]uint8{{7, 7}, {7, 8}}

	mustWritePNG(t, filepath.Join(srcDir, "img01.png"), imageA)
	mustWritePNG(t, filepath.Join(srcDir, "img02.png"), imageB)

	service := NewService()
	logBuffer := &bytes.Buffer{}
	result, err := service.Handle(Input{
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
