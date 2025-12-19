package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// FilterMode は使用可能なフィルター種別を表します。
type FilterMode string

const (
	// FilterModeGrayscale はグレースケール変換です。
	FilterModeGrayscale FilterMode = "grayscale"
	// FilterModeColorize はカラー調整（ティント付与）です。
	FilterModeColorize FilterMode = "colorize"
	// FilterModeVignette はビネット（周辺減光）です。
	FilterModeVignette FilterMode = "vignette"
)

// SupportedFilterModes はサポートされているフィルターモード一覧です。
var SupportedFilterModes = []FilterMode{
	FilterModeGrayscale,
	FilterModeColorize,
	FilterModeVignette,
}

// Config は CLI から渡される画像処理設定を表します。
type Config struct {
	InputPath    string
	OutputPath   string
	OutputFormat string
	Mode         FilterMode
	Strength     float64
	TintHex      string
}

// Validate は設定の妥当性を検証します。
func (c *Config) Validate() error {
	if c.InputPath == "" {
		return fmt.Errorf("input path is required")
	}

	info, err := os.Stat(c.InputPath)
	if err != nil {
		return fmt.Errorf("failed to stat input path: %w", err)
	}
	if info.IsDir() {
		return fmt.Errorf("input path must be a file: %s", c.InputPath)
	}

	if c.Mode == "" {
		return fmt.Errorf("filter mode is required")
	}
	if !isSupportedMode(c.Mode) {
		return fmt.Errorf("unsupported filter mode: %s", c.Mode)
	}

	if c.Strength < 0.0 || c.Strength > 1.0 {
		return fmt.Errorf("strength must be between 0 and 1: %.2f", c.Strength)
	}

	if c.OutputFormat != "" {
		switch strings.ToLower(c.OutputFormat) {
		case "png", "jpg", "jpeg":
		default:
			return fmt.Errorf("unsupported output format: %s", c.OutputFormat)
		}
	}

	if c.TintHex != "" {
		if _, err := ParseHexColor(c.TintHex); err != nil {
			return fmt.Errorf("invalid tint color: %w", err)
		}
	}

	return nil
}

// Normalise は不足している値を補完します。
func (c *Config) Normalise() {
	if c.TintHex == "" {
		c.TintHex = "#ffffff"
	}

	if c.OutputPath == "" {
		ext := c.OutputFormat
		if ext == "" {
			ext = strings.TrimPrefix(strings.ToLower(filepath.Ext(c.InputPath)), ".")
			if ext == "" {
				ext = "png"
			}
		}
		base := strings.TrimSuffix(filepath.Base(c.InputPath), filepath.Ext(c.InputPath))
		c.OutputPath = filepath.Join(filepath.Dir(c.InputPath), fmt.Sprintf("%s_filtered.%s", base, ext))
	}
}

func isSupportedMode(mode FilterMode) bool {
	for _, m := range SupportedFilterModes {
		if m == mode {
			return true
		}
	}
	return false
}

// ParseHexColor は "#RRGGBB" または "RRGGBB" 形式のカラーコードを正規化します。
func ParseHexColor(hex string) ([3]float32, error) {
	var rgb [3]float32

	h := strings.TrimPrefix(strings.TrimSpace(hex), "#")
	if len(h) != 6 {
		return rgb, fmt.Errorf("hex color must be 6 characters: %s", hex)
	}

	for i := 0; i < 3; i++ {
		part := h[i*2 : i*2+2]
		v, err := strconv.ParseUint(part, 16, 8)
		if err != nil {
			return rgb, fmt.Errorf("invalid hex component %q: %w", part, err)
		}
		rgb[i] = float32(v) / 255.0
	}

	return rgb, nil
}
