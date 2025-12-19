package usecases

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

// loadImage は標準ライブラリで画像を読み込みます。
func loadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// saveImage は PNG/JPEG で書き出します。
func saveImage(path, format string, img image.Image) (err error) {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := out.Close(); err == nil && cerr != nil {
			err = cerr
		}
	}()

	encoder := strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))
	if encoder == "" {
		encoder = strings.ToLower(format)
	}

	switch encoder {
	case "", "png":
		err = png.Encode(out, img)
	case "jpg", "jpeg":
		rgba := image.NewRGBA(img.Bounds())
		drawOpaque(rgba, img)
		err = jpeg.Encode(out, rgba, &jpeg.Options{Quality: 95})
	default:
		err = fmt.Errorf("unsupported output format: %s", encoder)
	}

	return err
}

// drawOpaque はアルファ付き画像を白背景に変換します。
func drawOpaque(dst *image.RGBA, src image.Image) {
	bounds := src.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			ix := x - bounds.Min.X
			iy := y - bounds.Min.Y
			c := color.NRGBAModel.Convert(src.At(x, y)).(color.NRGBA)
			if c.A == 0xff {
				dst.Set(ix, iy, c)
				continue
			}
			alpha := float64(c.A) / 255.0
			blended := color.RGBA{
				R: uint8(alpha*float64(c.R) + (1-alpha)*255.0 + 0.5),
				G: uint8(alpha*float64(c.G) + (1-alpha)*255.0 + 0.5),
				B: uint8(alpha*float64(c.B) + (1-alpha)*255.0 + 0.5),
				A: 0xff,
			}
			dst.Set(ix, iy, blended)
		}
	}
}
