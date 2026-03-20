package domain

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// Color はカラーコードの内部表現を保持する構造体
type Color struct {
	R float64 // Red (0-255)
	G float64 // Green (0-255)
	B float64 // Blue (0-255)
}

const maxDecimalColorValue = 16777215

// NewColor は新しいColorを作成する
func NewColor(r, g, b float64) *Color {
	return &Color{
		R: clamp(r, 0, 255),
		G: clamp(g, 0, 255),
		B: clamp(b, 0, 255),
	}
}

// clamp は値を指定された範囲内に制限する
func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// ParseFromHex はHEX形式の文字列からColorを作成する
func ParseFromHex(hexStr string) (*Color, error) {
	// #を除去し、大文字に変換
	hexStr = strings.TrimPrefix(hexStr, "#")
	hexStr = strings.ToUpper(hexStr)

	// 3桁の場合は6桁に展開
	if len(hexStr) == 3 {
		hexStr = string(hexStr[0]) + string(hexStr[0]) +
			string(hexStr[1]) + string(hexStr[1]) +
			string(hexStr[2]) + string(hexStr[2])
	}

	// 6桁でない場合はエラー
	if len(hexStr) != 6 {
		return nil, fmt.Errorf("無効なHEX形式です: %s", hexStr)
	}

	// 16進数の正規表現チェック
	matched, _ := regexp.MatchString("^[0-9A-F]{6}$", hexStr)
	if !matched {
		return nil, fmt.Errorf("無効なHEX形式です: %s", hexStr)
	}

	// RGB値に変換
	r, _ := strconv.ParseInt(hexStr[0:2], 16, 64)
	g, _ := strconv.ParseInt(hexStr[2:4], 16, 64)
	b, _ := strconv.ParseInt(hexStr[4:6], 16, 64)

	return NewColor(float64(r), float64(g), float64(b)), nil
}

// ParseFromRGB はRGB形式の文字列からColorを作成する
func ParseFromRGB(rgbStr string) (*Color, error) {
	// rgb(r,g,b) の形式をパース
	re := regexp.MustCompile(`rgb\s*\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*\)`)
	matches := re.FindStringSubmatch(strings.ToLower(rgbStr))

	if len(matches) != 4 {
		return nil, fmt.Errorf("無効なRGB形式です: %s", rgbStr)
	}

	r, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return nil, fmt.Errorf("無効なR値です: %s", matches[1])
	}

	g, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		return nil, fmt.Errorf("無効なG値です: %s", matches[2])
	}

	b, err := strconv.ParseFloat(matches[3], 64)
	if err != nil {
		return nil, fmt.Errorf("無効なB値です: %s", matches[3])
	}

	return NewColor(r, g, b), nil
}

// ParseFromHSL はHSL形式の文字列からColorを作成する
func ParseFromHSL(hslStr string) (*Color, error) {
	// hsl(h,s%,l%) の形式をパース
	re := regexp.MustCompile(`hsl\s*\(\s*(\d+(?:\.\d+)?)\s*,\s*(\d+(?:\.\d+)?)%\s*,\s*(\d+(?:\.\d+)?)%\s*\)`)
	matches := re.FindStringSubmatch(strings.ToLower(hslStr))

	if len(matches) != 4 {
		return nil, fmt.Errorf("無効なHSL形式です: %s", hslStr)
	}

	h, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return nil, fmt.Errorf("無効なH値です: %s", matches[1])
	}

	s, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		return nil, fmt.Errorf("無効なS値です: %s", matches[2])
	}

	l, err := strconv.ParseFloat(matches[3], 64)
	if err != nil {
		return nil, fmt.Errorf("無効なL値です: %s", matches[3])
	}

	return hslToRgb(h, s/100, l/100), nil
}

// ParseFromHSV はHSV形式の文字列からColorを作成する
func ParseFromHSV(hsvStr string) (*Color, error) {
	// hsv(h,s%,v%) の形式をパース
	re := regexp.MustCompile(`hsv\s*\(\s*(\d+(?:\.\d+)?)\s*,\s*(\d+(?:\.\d+)?)%\s*,\s*(\d+(?:\.\d+)?)%\s*\)`)
	matches := re.FindStringSubmatch(strings.ToLower(hsvStr))

	if len(matches) != 4 {
		return nil, fmt.Errorf("無効なHSV形式です: %s", hsvStr)
	}

	h, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return nil, fmt.Errorf("無効なH値です: %s", matches[1])
	}

	s, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		return nil, fmt.Errorf("無効なS値です: %s", matches[2])
	}

	v, err := strconv.ParseFloat(matches[3], 64)
	if err != nil {
		return nil, fmt.Errorf("無効なV値です: %s", matches[3])
	}

	return hsvToRgb(h, s/100, v/100), nil
}

// ParseFromDecimal は10進数形式の文字列からColorを作成する
func ParseFromDecimal(decStr string) (*Color, error) {
	decStr = strings.TrimSpace(decStr)
	if decStr == "" {
		return nil, fmt.Errorf("無効なDEC形式です: %s", decStr)
	}

	matched, _ := regexp.MatchString(`^\d+$`, decStr)
	if !matched {
		return nil, fmt.Errorf("無効なDEC形式です: %s", decStr)
	}

	decValue, err := strconv.ParseInt(decStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("無効なDEC形式です: %s", decStr)
	}

	if decValue < 0 || decValue > maxDecimalColorValue {
		return nil, fmt.Errorf("DEC値が範囲外です: %s", decStr)
	}

	r := (decValue >> 16) & 0xFF
	g := (decValue >> 8) & 0xFF
	b := decValue & 0xFF

	return NewColor(float64(r), float64(g), float64(b)), nil
}

// ToHex はColorをHEX形式の文字列に変換する
func (c *Color) ToHex() string {
	r := int(math.Round(c.R))
	g := int(math.Round(c.G))
	b := int(math.Round(c.B))
	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}

// ToRGB はColorをRGB形式の文字列に変換する
func (c *Color) ToRGB() string {
	r := int(math.Round(c.R))
	g := int(math.Round(c.G))
	b := int(math.Round(c.B))
	return fmt.Sprintf("rgb(%d,%d,%d)", r, g, b)
}

// ToHSL はColorをHSL形式の文字列に変換する
func (c *Color) ToHSL() string {
	h, s, l := rgbToHsl(c.R, c.G, c.B)
	return fmt.Sprintf("hsl(%.0f,%.0f%%,%.0f%%)", h, s*100, l*100)
}

// ToHSV はColorをHSV形式の文字列に変換する
func (c *Color) ToHSV() string {
	h, s, v := rgbToHsv(c.R, c.G, c.B)
	return fmt.Sprintf("hsv(%.0f,%.0f%%,%.0f%%)", h, s*100, v*100)
}

// ToDecimal はColorを10進数形式の文字列に変換する
func (c *Color) ToDecimal() string {
	r := int(math.Round(c.R))
	g := int(math.Round(c.G))
	b := int(math.Round(c.B))
	decValue := (r << 16) | (g << 8) | b
	return strconv.Itoa(decValue)
}

// rgbToHsl はRGB値をHSL値に変換する
func rgbToHsl(r, g, b float64) (h, s, l float64) {
	r /= 255
	g /= 255
	b /= 255

	max := math.Max(math.Max(r, g), b)
	min := math.Min(math.Min(r, g), b)
	l = (max + min) / 2

	if max == min {
		h = 0
		s = 0
	} else {
		d := max - min
		if l > 0.5 {
			s = d / (2 - max - min)
		} else {
			s = d / (max + min)
		}

		switch max {
		case r:
			h = (g-b)/d + func() float64 {
				if g < b {
					return 6
				}
				return 0
			}()
		case g:
			h = (b-r)/d + 2
		case b:
			h = (r-g)/d + 4
		}
		h /= 6
	}

	h *= 360
	if h < 0 {
		h += 360
	}

	return h, s, l
}

// hslToRgb はHSL値をRGB値に変換する
func hslToRgb(h, s, l float64) *Color {
	h /= 360

	var r, g, b float64

	if s == 0 {
		r = l
		g = l
		b = l
	} else {
		hue2rgb := func(p, q, t float64) float64 {
			if t < 0 {
				t += 1
			}
			if t > 1 {
				t -= 1
			}
			if t < 1.0/6 {
				return p + (q-p)*6*t
			}
			if t < 1.0/2 {
				return q
			}
			if t < 2.0/3 {
				return p + (q-p)*(2.0/3-t)*6
			}
			return p
		}

		var q float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - l*s
		}
		p := 2*l - q

		r = hue2rgb(p, q, h+1.0/3)
		g = hue2rgb(p, q, h)
		b = hue2rgb(p, q, h-1.0/3)
	}

	return NewColor(r*255, g*255, b*255)
}

// rgbToHsv はRGB値をHSV値に変換する
func rgbToHsv(r, g, b float64) (h, s, v float64) {
	r /= 255
	g /= 255
	b /= 255

	max := math.Max(math.Max(r, g), b)
	min := math.Min(math.Min(r, g), b)
	v = max

	d := max - min
	if max == 0 {
		s = 0
	} else {
		s = d / max
	}

	if max == min {
		h = 0
	} else {
		switch max {
		case r:
			h = (g-b)/d + func() float64 {
				if g < b {
					return 6
				}
				return 0
			}()
		case g:
			h = (b-r)/d + 2
		case b:
			h = (r-g)/d + 4
		}
		h /= 6
	}

	h *= 360
	if h < 0 {
		h += 360
	}

	return h, s, v
}

// hsvToRgb はHSV値をRGB値に変換する
func hsvToRgb(h, s, v float64) *Color {
	h /= 360

	c := v * s
	x := c * (1 - math.Abs(math.Mod(h*6, 2)-1))
	m := v - c

	var r, g, b float64

	switch {
	case h < 1.0/6:
		r, g, b = c, x, 0
	case h < 2.0/6:
		r, g, b = x, c, 0
	case h < 3.0/6:
		r, g, b = 0, c, x
	case h < 4.0/6:
		r, g, b = 0, x, c
	case h < 5.0/6:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}

	return NewColor((r+m)*255, (g+m)*255, (b+m)*255)
}
