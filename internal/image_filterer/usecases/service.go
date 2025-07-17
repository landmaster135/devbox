package usecases

import (
	"fmt"
	"image"
	"image/draw"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/anthonynsimon/bild/blur"
	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/gen2brain/webp"
)

// FilterMode は適用するフィルターの種類を表します
type FilterMode string

const (
	// BlurMode はぼかしフィルターを表します
	BlurMode FilterMode = "blur"
	// GrayscaleMode はグレースケールフィルターを表します
	GrayscaleMode FilterMode = "grayscale"
)

func multiplyAndRound(value int, multiplier float64) int {
	// 整数値をfloat64に変換して指定された倍率で掛ける
	multiplied := float64(value) * multiplier
	// 四捨五入して整数に変換
	rounded := int(math.Round(multiplied))
	return rounded
}
// 画像を読み込み、指定した領域にフィルターを適用して outDir に保存
func ApplyFilterAndSave(inPath, outDir string, x1, y1, x2, y2 int, suffix string, mode FilterMode, radius float64, rWeight, gWeight, bWeight float64) error {
	// ログ出力用のフォーマット文字列
	logFormat := "処理情報: ファイル=%s, 範囲=(%d,%d)-(%d,%d), モード=%s, 半径=%.1f"
	fmt.Printf(logFormat+"\n", filepath.Base(inPath), x1, y1, x2, y2, mode, radius)

	// ── 読み込み ──
	img, err := imgio.Open(inPath) // 拡張子を気にせずデコード
	if err != nil {
		return fmt.Errorf("open %s: %w", inPath, err)
	}

	// ── 矩形チェック ──
	bounds := img.Bounds()
	// 画像の実際の境界を取得
	minX, minY := bounds.Min.X, bounds.Min.Y
	maxX, maxY := bounds.Max.X, bounds.Max.Y

	fmt.Printf("画像情報: サイズ=%dx%d, 境界=(%d,%d)-(%d,%d)\n",
		bounds.Dx(), bounds.Dy(), minX, minY, maxX, maxY)

	// 座標が画像の範囲内かチェック
	if x1 < minX || y1 < minY || x2 > maxX || y2 > maxY || x2 <= x1 || y2 <= y1 {
		return fmt.Errorf("invalid rectangle (%d,%d)-(%d,%d) for %s with bounds (%d,%d)-(%d,%d)",
			x1, y1, x2, y2, inPath, minX, minY, maxX, maxY)
	}

	// ── 元画像のクローンを作成 ──
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, img, bounds.Min, draw.Src)

	// ── 指定領域を切り出し ──
	// x1, y1, x2, y2 = multiplyAndRound(x1, 1.5), multiplyAndRound(y1, 1.5), multiplyAndRound(x2, 1.5), multiplyAndRound(y2, 1.5)
	cropRect := image.Rect(x1, y1, x2, y2)
	fmt.Printf("切り出し範囲: (%d,%d)-(%d,%d)\n", cropRect.Min.X, cropRect.Min.Y, cropRect.Max.X, cropRect.Max.Y)
	subImg := transform.Crop(img, cropRect)

	// ── フィルター適用 ──
	var filtered *image.RGBA
	switch mode {
	case BlurMode:
		filtered = blur.Gaussian(subImg, radius)
	case GrayscaleMode:
		filtered = effect.GrayscaleWithWeights(subImg, rWeight, gWeight, bWeight)
	default:
		return fmt.Errorf("unsupported filter mode: %s", mode)
	}

	// ── フィルター適用した領域を元の画像に合成 ──
	targetRect := image.Rect(x1, y1, x2, y2)
	fmt.Printf("合成先範囲: (%d,%d)-(%d,%d)\n", targetRect.Min.X, targetRect.Min.Y, targetRect.Max.X, targetRect.Max.Y)
	fmt.Printf("フィルター適用画像の境界: (%d,%d)-(%d,%d)\n",
		filtered.Bounds().Min.X, filtered.Bounds().Min.Y,
		filtered.Bounds().Max.X, filtered.Bounds().Max.Y)

	draw.Draw(dst, targetRect, filtered, filtered.Bounds().Min, draw.Src)

	// ── 保存パス準備 ──
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	base := strings.TrimSuffix(filepath.Base(inPath), filepath.Ext(inPath))
	outPath := filepath.Join(outDir,
		fmt.Sprintf("%s_%s%s", base, suffix, filepath.Ext(inPath)))

	// ── 形式に応じて保存 ──
	ext := strings.ToLower(filepath.Ext(outPath))
	fmt.Printf("保存形式: %s, 保存先: %s\n", ext, outPath)

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
