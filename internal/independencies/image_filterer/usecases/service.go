package usecases

import (
	"fmt"
	"image"
	"image/draw"
	"os"
	"path/filepath"
	"strings"

	"github.com/anthonynsimon/bild/blur"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/gen2brain/webp"
)

// FilterMode は適用するフィルターの種類を表します
type FilterMode string

const (
	// BlurMode はぼかしフィルターを表します
	BlurMode FilterMode = "blur"
)

// 画像を読み込み、指定した領域にフィルターを適用して outDir に保存
func ApplyFilterAndSave(inPath, outDir string, x1, y1, x2, y2 int, suffix string, mode FilterMode, radius float64) error {
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

	// ── 元画像のクローンを作成 ──
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, img, bounds.Min, draw.Src)

	// ── 指定領域を切り出し ──
	subImg := transform.Crop(img, image.Rect(x1, y1, x2, y2))

	// ── フィルター適用 ──
	var filtered *image.RGBA
	switch mode {
	case BlurMode:
		filtered = blur.Gaussian(subImg, radius)
	default:
		return fmt.Errorf("unsupported filter mode: %s", mode)
	}

	// ── フィルター適用した領域を元の画像に合成 ──
	draw.Draw(dst, image.Rect(x1, y1, x2, y2), filtered, filtered.Bounds().Min, draw.Src)

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
		err = imgio.Save(outPath, dst, imgio.JPEGEncoder(95))
	case ".png":
		err = imgio.Save(outPath, dst, imgio.PNGEncoder())
	case ".webp":
		f, err := os.Create(outPath)
		if err != nil {
			return err
		}
		defer f.Close()
		err = webp.Encode(f, dst, webp.Options{Quality: 90})
		if err != nil {
			return err
		}
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
