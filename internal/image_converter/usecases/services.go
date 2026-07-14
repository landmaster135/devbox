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
	Decode func(io.Reader) (image.Image, error)
	Encode func(image.Image) ([]byte, error)
}

// Individual decoder helpers …
func decodeJPEG(r io.Reader) (image.Image, error) { return jpeg.Decode(r) }
func decodePNG(r io.Reader) (image.Image, error)  { return png.Decode(r) }
func decodeWebP(r io.Reader) (image.Image, error) { return xwebp.Decode(r) }
func decodeAVIF(r io.Reader) (image.Image, error) { return avif.Decode(r) }
func decodeSVG(r io.Reader) (image.Image, error) {
	icon, err := oksvg.ReadIconStream(r)
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
func MakeCodecTable(q int, lossless bool) map[string]codec {
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
			Lossless: lossless,
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
func ConvertFile(path, sourceDir, outputDir, outExt string, table map[string]codec) error {
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
	inFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = inFile.Close() }()

	info, err := inFile.Stat()
	if err != nil {
		return err
	}
	origSize := info.Size()

	// --- decode & encode ---
	img, err := cIn.Decode(inFile)
	if err != nil {
		return err
	}
	_ = inFile.Close() // Windowsの排他制御を避けるため明示的にクローズ
	encBytes, err := cOut.Encode(img)
	if err != nil {
		return err
	}

	// --- decide which bytes to save ---
	dstExt := outExt
	note := "converted"

	// サイズ比較は「SVG → 非SVG」、「Avif → Webp」、「any → Jpg」以外にのみ適用
	requiresForcedConvert := (inExt == "svg" && outExt != "svg") || (inExt == "avif" && outExt == "webp") || (outExt == "jpg")

	if !requiresForcedConvert && origSize > 0 && int64(len(encBytes)) >= origSize {
		dstExt = inExt
		note = "kept original"
		log.Printf("  kept original (%s): %d → %d bytes", path, origSize, len(encBytes))
	}

	// --- write ---
	rel, _ := filepath.Rel(sourceDir, path)
	base := strings.TrimSuffix(rel, filepath.Ext(rel))
	outPath := filepath.Join(outputDir, base+"."+dstExt)

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}

	var savedSize int64
	if note == "kept original" {
		if err := copyFile(path, outPath); err != nil {
			return err
		}
		savedSize = origSize
	} else {
		if err := os.WriteFile(outPath, encBytes, 0o644); err != nil {
			return err
		}
		savedSize = int64(len(encBytes))
	}

	var ratio float64
	if origSize > 0 {
		ratio = float64(savedSize) / float64(origSize) * 100
	}
	log.Printf("%s → %s (%s): %d → %d bytes | compressed to %.1f%%",
		path, outPath, note, origSize, savedSize, ratio)

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

// preserving the relative directory structure under sourceDir.
func ArchiveOriginal(srcPath, sourceDir, archiveDir string, move bool) error {
	if archiveDir == "" {
		return nil
	}

	rel, err := filepath.Rel(sourceDir, srcPath)
	if err != nil {
		return err
	}
	dest := filepath.Join(archiveDir, rel)

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}

	if !move {
		return copyFile(srcPath, dest)
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
