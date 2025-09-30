package usecases

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/landmaster135/devbox/internal/image_filterer_v2/config"
)

type headlessProcessor struct {
	cfg config.Config
}

func newProcessor(cfg config.Config) (processor, error) {
	return &headlessProcessor{cfg: cfg}, nil
}

func (p *headlessProcessor) Process() (string, error) {
	srcImg, err := loadImage(p.cfg.InputPath)
	if err != nil {
		return "", fmt.Errorf("failed to load input: %w", err)
	}

	result, err := p.applyFilter(srcImg)
	if err != nil {
		return "", err
	}

	if err := saveImage(p.cfg.OutputPath, p.cfg.OutputFormat, result); err != nil {
		return "", fmt.Errorf("failed to save output: %w", err)
	}

	return p.cfg.OutputPath, nil
}

func (p *headlessProcessor) applyFilter(src image.Image) (image.Image, error) {
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	tint, err := config.ParseHexColor(p.cfg.TintHex)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tint color: %w", err)
	}

	dst := image.NewNRGBA(bounds)
	strength := clampFloat(float32(p.cfg.Strength), 0, 1)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.NRGBAModel.Convert(src.At(x, y)).(color.NRGBA)
			rf := float32(c.R) / 255.0
			gf := float32(c.G) / 255.0
			bf := float32(c.B) / 255.0

			switch p.cfg.Mode {
			case config.FilterModeGrayscale:
				gray := rf*0.299 + gf*0.587 + bf*0.114
				rf = mix(rf, gray, strength)
				gf = mix(gf, gray, strength)
				bf = mix(bf, gray, strength)
			case config.FilterModeColorize:
				rf = mix(rf, tint[0], strength)
				gf = mix(gf, tint[1], strength)
				bf = mix(bf, tint[2], strength)
			case config.FilterModeVignette:
				uvx := (float32(x-bounds.Min.X) + 0.5) / float32(width)
				uvy := (float32(y-bounds.Min.Y) + 0.5) / float32(height)
				dx := uvx - 0.5
				dy := uvy - 0.5
				dist := float32(math.Sqrt(float64(dx*dx + dy*dy)))
				start := float32(0.35)
				end := float32(0.7) - strength*0.3
				if end < 0.25 {
					end = 0.25
				}
				vign := smoothstep(end, start, dist)
				rf *= vign
				gf *= vign
				bf *= vign
			default:
				return nil, fmt.Errorf("unsupported filter mode: %s", p.cfg.Mode)
			}

			dst.SetNRGBA(x, y, color.NRGBA{
				R: floatToByte(rf),
				G: floatToByte(gf),
				B: floatToByte(bf),
				A: c.A,
			})
		}
	}

	return dst, nil
}

func mix(a, b, t float32) float32 {
	return a*(1-t) + b*t
}

func smoothstep(edge0, edge1, x float32) float32 {
	if edge0 == edge1 {
		if x < edge0 {
			return 0
		}
		return 1
	}
	t := (x - edge0) / (edge1 - edge0)
	t = clampFloat(t, 0, 1)
	return t * t * (3 - 2*t)
}

func clampFloat(v, min, max float32) float32 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func floatToByte(v float32) uint8 {
	v = clampFloat(v, 0, 1)
	return uint8(v*255 + 0.5)
}
