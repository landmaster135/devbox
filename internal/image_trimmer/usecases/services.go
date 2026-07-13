package usecases

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

// 画像を読み込みトリミングして outputDir に保存
func CropAndSave(inputPath, outputDir string, x1, y1, x2, y2 int, suffix string) error {
	// ── 読み込み ──
	img, err := imgio.Open(inputPath) // 拡張子を気にせずデコード
	if err != nil {
		return fmt.Errorf("open %s: %w", inputPath, err)
	}

	// ── 矩形チェック ──
	bounds := img.Bounds()
	if x1 < 0 || y1 < 0 || x2 <= x1 || y2 <= y1 ||
		x2 > bounds.Dx() || y2 > bounds.Dy() {
		return fmt.Errorf("invalid rectangle %v for %s",
			image.Rect(x1, y1, x2, y2), inputPath)
	}

	// ── トリミング ──
	cropped := transform.Crop(img, image.Rect(x1, y1, x2, y2))

	// ── 保存パス準備 ──
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return err
	}
	base := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	outPath := filepath.Join(outputDir,
		fmt.Sprintf("%s_%s%s", base, suffix, filepath.Ext(inputPath)))

	// ── 形式に応じて保存 ──
	ext := strings.ToLower(filepath.Ext(outPath))
	switch ext {
	case ".jpg", ".jpeg":
		err = imgio.Save(outPath, cropped, imgio.JPEGEncoder(95))
	case ".png":
		err = imgio.Save(outPath, cropped, imgio.PNGEncoder())
	default:
		err = fmt.Errorf("unsupported extension: %s", ext)
	}
	return err
}

// 元画像を archiveDir に移動
func MoveOriginal(sourceFile, archiveDir string) error {
	if err := os.MkdirAll(archiveDir, 0o755); err != nil {
		return err
	}
	dst := filepath.Join(archiveDir, filepath.Base(sourceFile))
	return os.Rename(sourceFile, dst)
}
