package usecases

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

// ---------------------------------------------------------
// 画像判定（拡張子）
// ---------------------------------------------------------
func IsImage(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".jpg", ".jpeg", ".png":
		return true
	default:
		return false
	}
}

// ---------------------------------------------------------
// actual rotate & save
// ---------------------------------------------------------
func RotateAndSave(srcPath, outDir, suffix string, angle float64) error {
	img, err := imgio.Open(srcPath)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}

	rot := transform.Rotate(img, angle, &transform.RotationOptions{
		ResizeBounds: true, // 回転後に画像全体が入るキャンバスへ拡大
		// Pivot は nil なら自動で中心
	})

	ext := strings.ToLower(filepath.Ext(srcPath))
	name := strings.TrimSuffix(filepath.Base(srcPath), ext)
	outPath := filepath.Join(outDir, fmt.Sprintf("%s_%s%s", name, suffix, ext))

	if err := ensureDir(outDir); err != nil {
		return err
	}

	var enc imgio.Encoder
	switch ext {
	case ".jpg", ".jpeg":
		enc = imgio.JPEGEncoder(95) // 品質 95 でエンコーダー作成
	case ".png":
		enc = imgio.PNGEncoder()
	}

	if err := imgio.Save(outPath, rot, enc); err != nil {
		return fmt.Errorf("save: %w", err)
	}
	return nil
}

// ---------------------------------------------------------
// originals をアーカイブへ移動
// ---------------------------------------------------------
func MoveFile(src, dstDir string) error {
	if err := ensureDir(dstDir); err != nil {
		return err
	}
	dst := filepath.Join(dstDir, filepath.Base(src))
	return os.Rename(src, dst)
}

// dstDir が無ければ作成
func ensureDir(dir string) error {
	return os.MkdirAll(dir, 0o755)
}
