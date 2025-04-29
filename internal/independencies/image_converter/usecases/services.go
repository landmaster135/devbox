package usecases

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	avif "github.com/gen2brain/avif"
	webp "github.com/gen2brain/webp"
	oksvg "github.com/srwiley/oksvg"
	rasterx "github.com/srwiley/rasterx"
	xwebp "golang.org/x/image/webp"
)

type codec struct {
	Decode func([]byte) (image.Image, error)
	Encode func(image.Image) ([]byte, error)
}

// Individual decoder helpers …
func decodeJPEG(b []byte) (image.Image, error) { return jpeg.Decode(bytes.NewReader(b)) }
func decodePNG(b []byte) (image.Image, error)  { return png.Decode(bytes.NewReader(b)) }
func decodeWebP(b []byte) (image.Image, error) { return xwebp.Decode(bytes.NewReader(b)) }
func decodeAVIF(b []byte) (image.Image, error) { return avif.Decode(bytes.NewReader(b)) }
func decodeSVG(b []byte) (image.Image, error) {
	icon, err := oksvg.ReadIconStream(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	w, h := int(icon.ViewBox.W), int(icon.ViewBox.H)
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	icon.SetTarget(0, 0, float64(w), float64(h))
	scanner := rasterx.NewScannerGV(w, h, rgba, rgba.Bounds())
	drawer := rasterx.NewDasher(w, h, scanner)
	icon.Draw(drawer, 1.0)
	return rgba, nil
}

// MakeCodecTable returns decoder/encoder map keyed by extension
func MakeCodecTable(q int) map[string]codec {
	encJPEG := func(img image.Image) ([]byte, error) {
		var buf bytes.Buffer
		err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: q})
		return buf.Bytes(), err
	}
	encPNG := func(img image.Image) ([]byte, error) {
		var buf bytes.Buffer
		err := png.Encode(&buf, img)
		return buf.Bytes(), err
	}
	encWebP := func(img image.Image) ([]byte, error) {
		var buf bytes.Buffer
		// quality=0–100, method=0–6 (速度↔品質)
		err := webp.Encode(&buf, img, webp.Options{
			Quality:  q,
			Method:   6,
			Lossless: false,
		})
		return buf.Bytes(), err
	}
	encAVIF := func(img image.Image) ([]byte, error) {
		var buf bytes.Buffer
		opt := avif.Options{Quality: q, Speed: 6} // 好みに応じて
		err := avif.Encode(&buf, img, opt)
		return buf.Bytes(), err
	}

	return map[string]codec{
		"jpg":  {decodeJPEG, encJPEG},
		"jpeg": {decodeJPEG, encJPEG},
		"png":  {decodePNG, encPNG},
		"webp": {decodeWebP, encWebP},
		"avif": {decodeAVIF, encAVIF},
		"svg":  {decodeSVG, nil}, // encode is nil because vector→vector not supported
	}
}

// ConvertFile runs a single conversion job
func ConvertFile(path, srcDir, outDir, outExt string, table map[string]codec) error {
	inExt := strings.TrimPrefix(strings.ToLower(filepath.Ext(path)), ".")
	cIn, ok := table[inExt]
	if !ok {
		return fmt.Errorf("no decoder for %s", inExt)
	}
	cOut := table[outExt]
	if cOut.Encode == nil {
		return fmt.Errorf("encode to .%s not supported", outExt)
	}

	// --- read original ---
	origBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// --- decode & encode ---
	img, err := cIn.Decode(origBytes)
	if err != nil {
		return err
	}
	encBytes, err := cOut.Encode(img)
	if err != nil {
		return err
	}

	// --- decide which bytes to save ---
	saveBytes := encBytes
	dstExt := outExt
	note := "converted"

	// ◎ サイズ比較は「SVG → 非SVG」、「Avif → Webp」、「any → Jpg」以外にのみ適用
	requiresForcedConvert := (inExt == "svg" && outExt != "svg") || (inExt == "avif" && outExt == "webp") || (outExt == "jpg")

	if !requiresForcedConvert && len(encBytes) >= len(origBytes) {
		saveBytes = origBytes
		dstExt = inExt
		note = "kept original"
		log.Printf("  kept original (%s): %d → %d bytes", path, len(origBytes), len(encBytes))
	}

	// --- write ---
	rel, _ := filepath.Rel(srcDir, path)
	base := strings.TrimSuffix(rel, filepath.Ext(rel))
	outPath := filepath.Join(outDir, base+"."+dstExt)

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(outPath, saveBytes, 0o644); err != nil {
		return err
	}

	ratio := float64(len(saveBytes)) / float64(len(origBytes)) * 100
	log.Printf("%s → %s (%s): %d → %d bytes | compressed to %.1f%%",
		path, outPath, note, len(origBytes), len(saveBytes), ratio)

	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

// preserving the relative directory structure under srcDir.
func ArchiveOriginal(srcPath, srcDir, archiveDir string, move bool) error {
	if !move {
		return nil // 触らない
	}

	rel, _ := filepath.Rel(srcDir, srcPath)
	dest := filepath.Join(archiveDir, rel)

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}

	// 1) rename で速く 2) cross-device fallback
	if err := os.Rename(srcPath, dest); err == nil {
		return nil
	}
	if err := copyFile(srcPath, dest); err != nil {
		return err
	}
	return os.Remove(srcPath)
}
