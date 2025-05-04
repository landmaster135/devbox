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

// 画像を読み込みトリミングして outDir に保存
func CropAndSave(inPath, outDir string, x1, y1, x2, y2 int, suffix string) error {
	// ── 読み込み ──
	img, err := imgio.Open(inPath) // 拡張子を気にせずデコード
	if err != nil {
		return fmt.Errorf("open %s: %w", inPath, err)
	}

	// ── 矩形チェック ──
	bounds := img.Bounds()
	if x1 < 0 || y1 < 0 || x2 <= x1 || y2 <= y1 ||
		x2 > bounds.Dx() || y2 > bounds.Dy() {
		return fmt.Errorf("invalid rectangle %v for %s",
			image.Rect(x1, y1, x2, y2), inPath)
	}

	// ── トリミング ──
	cropped := transform.Crop(img, image.Rect(x1, y1, x2, y2))

	// ── 保存パス準備 ──
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	base := strings.TrimSuffix(filepath.Base(inPath), filepath.Ext(inPath))
	outPath := filepath.Join(outDir,
		fmt.Sprintf("%s_%s%s", base, suffix, filepath.Ext(inPath)))

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

// 元画像を arcDir に移動
func MoveOriginal(src, arcDir string) error {
	if err := os.MkdirAll(arcDir, 0o755); err != nil {
		return err
	}
	dst := filepath.Join(arcDir, filepath.Base(src))
	return os.Rename(src, dst)
}
