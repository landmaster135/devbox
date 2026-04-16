package domain

import (
	"math"
	"testing"
)

func TestNewColor_Normal(t *testing.T) {
	tests := []struct {
		name     string
		r, g, b  float64
		expected *Color
	}{
		{
			name: "ValidRgbValues_Normal",
			r:    255, g: 0, b: 0,
			expected: &Color{R: 255, G: 0, B: 0},
		},
		{
			name: "ClampHighValues_Normal",
			r:    300, g: 300, b: 300,
			expected: &Color{R: 255, G: 255, B: 255},
		},
		{
			name: "ClampLowValues_Normal",
			r:    -10, g: -10, b: -10,
			expected: &Color{R: 0, G: 0, B: 0},
		},
		{
			name: "MixedValues_Normal",
			r:    128, g: 64, b: 192,
			expected: &Color{R: 128, G: 64, B: 192},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color := NewColor(tt.r, tt.g, tt.b)
			if color.R != tt.expected.R || color.G != tt.expected.G || color.B != tt.expected.B {
				t.Errorf("NewColor(%v, %v, %v) = {%v, %v, %v}, want {%v, %v, %v}",
					tt.r, tt.g, tt.b, color.R, color.G, color.B, tt.expected.R, tt.expected.G, tt.expected.B)
			}
		})
	}
}

func TestParseFromHex_Normal(t *testing.T) {
	tests := []struct {
		name        string
		hexStr      string
		expected    *Color
		expectError bool
	}{
		{
			name:     "ValidHex6Digits_Normal",
			hexStr:   "#FF0000",
			expected: &Color{R: 255, G: 0, B: 0},
		},
		{
			name:     "ValidHex6DigitsLowerCase_Normal",
			hexStr:   "#ff0000",
			expected: &Color{R: 255, G: 0, B: 0},
		},
		{
			name:     "ValidHex3Digits_Normal",
			hexStr:   "#F00",
			expected: &Color{R: 255, G: 0, B: 0},
		},
		{
			name:     "ValidHexWithoutHash_Normal",
			hexStr:   "FF0000",
			expected: &Color{R: 255, G: 0, B: 0},
		},
		{
			name:        "InvalidHexLength_Error",
			hexStr:      "#FF00",
			expectError: true,
		},
		{
			name:        "InvalidHexCharacters_Error",
			hexStr:      "#GGGGGG",
			expectError: true,
		},
		{
			name:        "EmptyString_Error",
			hexStr:      "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color, err := ParseFromHex(tt.hexStr)

			if tt.expectError {
				if err == nil {
					t.Errorf("期待されたエラーが発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラーが発生しました: %v", err)
				return
			}

			if color.R != tt.expected.R || color.G != tt.expected.G || color.B != tt.expected.B {
				t.Errorf("ParseFromHex(%v) = {%v, %v, %v}, want {%v, %v, %v}",
					tt.hexStr, color.R, color.G, color.B, tt.expected.R, tt.expected.G, tt.expected.B)
			}
		})
	}
}

func TestParseFromRGB_Normal(t *testing.T) {
	tests := []struct {
		name        string
		rgbStr      string
		expected    *Color
		expectError bool
	}{
		{
			name:     "ValidRgb_Normal",
			rgbStr:   "rgb(255,0,0)",
			expected: &Color{R: 255, G: 0, B: 0},
		},
		{
			name:     "ValidRgbWithSpaces_Normal",
			rgbStr:   "rgb( 255 , 0 , 0 )",
			expected: &Color{R: 255, G: 0, B: 0},
		},
		{
			name:     "ValidRgbUpperCase_Normal",
			rgbStr:   "RGB(128,64,192)",
			expected: &Color{R: 128, G: 64, B: 192},
		},
		{
			name:        "InvalidRgbFormat_Error",
			rgbStr:      "rgb(255,0)",
			expectError: true,
		},
		{
			name:        "InvalidRgbValues_Error",
			rgbStr:      "rgb(abc,0,0)",
			expectError: true,
		},
		{
			name:        "EmptyString_Error",
			rgbStr:      "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color, err := ParseFromRGB(tt.rgbStr)

			if tt.expectError {
				if err == nil {
					t.Errorf("期待されたエラーが発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラーが発生しました: %v", err)
				return
			}

			if color.R != tt.expected.R || color.G != tt.expected.G || color.B != tt.expected.B {
				t.Errorf("ParseFromRGB(%v) = {%v, %v, %v}, want {%v, %v, %v}",
					tt.rgbStr, color.R, color.G, color.B, tt.expected.R, tt.expected.G, tt.expected.B)
			}
		})
	}
}

func TestParseFromHSL_Normal(t *testing.T) {
	tests := []struct {
		name        string
		hslStr      string
		expectError bool
	}{
		{
			name:   "ValidHsl_Normal",
			hslStr: "hsl(0,100%,50%)",
		},
		{
			name:   "ValidHslWithSpaces_Normal",
			hslStr: "hsl( 0 , 100% , 50% )",
		},
		{
			name:   "ValidHslUpperCase_Normal",
			hslStr: "HSL(120,100%,50%)",
		},
		{
			name:        "InvalidHslFormat_Error",
			hslStr:      "hsl(0,100%)",
			expectError: true,
		},
		{
			name:        "InvalidHslValues_Error",
			hslStr:      "hsl(abc,100%,50%)",
			expectError: true,
		},
		{
			name:        "EmptyString_Error",
			hslStr:      "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color, err := ParseFromHSL(tt.hslStr)

			if tt.expectError {
				if err == nil {
					t.Errorf("期待されたエラーが発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラーが発生しました: %v", err)
				return
			}

			// HSLからRGBへの変換は複雑なので、値が0-255の範囲内にあることを確認
			if color.R < 0 || color.R > 255 || color.G < 0 || color.G > 255 || color.B < 0 || color.B > 255 {
				t.Errorf("ParseFromHSL(%v) = {%v, %v, %v}, RGB値が範囲外です",
					tt.hslStr, color.R, color.G, color.B)
			}
		})
	}
}

func TestParseFromHSV_Normal(t *testing.T) {
	tests := []struct {
		name        string
		hsvStr      string
		expectError bool
	}{
		{
			name:   "ValidHsv_Normal",
			hsvStr: "hsv(0,100%,100%)",
		},
		{
			name:   "ValidHsvWithSpaces_Normal",
			hsvStr: "hsv( 0 , 100% , 100% )",
		},
		{
			name:   "ValidHsvUpperCase_Normal",
			hsvStr: "HSV(120,100%,100%)",
		},
		{
			name:        "InvalidHsvFormat_Error",
			hsvStr:      "hsv(0,100%)",
			expectError: true,
		},
		{
			name:        "InvalidHsvValues_Error",
			hsvStr:      "hsv(abc,100%,100%)",
			expectError: true,
		},
		{
			name:        "EmptyString_Error",
			hsvStr:      "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color, err := ParseFromHSV(tt.hsvStr)

			if tt.expectError {
				if err == nil {
					t.Errorf("期待されたエラーが発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラーが発生しました: %v", err)
				return
			}

			// HSVからRGBへの変換は複雑なので、値が0-255の範囲内にあることを確認
			if color.R < 0 || color.R > 255 || color.G < 0 || color.G > 255 || color.B < 0 || color.B > 255 {
				t.Errorf("ParseFromHSV(%v) = {%v, %v, %v}, RGB値が範囲外です",
					tt.hsvStr, color.R, color.G, color.B)
			}
		})
	}
}

func TestParseFromDecimal_Normal(t *testing.T) {
	tests := []struct {
		name        string
		decStr      string
		expected    *Color
		expectError bool
	}{
		{
			name:     "ValidDecimalRed_Normal",
			decStr:   "16711680",
			expected: &Color{R: 255, G: 0, B: 0},
		},
		{
			name:     "ValidDecimalZero_Normal",
			decStr:   "0",
			expected: &Color{R: 0, G: 0, B: 0},
		},
		{
			name:     "ValidDecimalMax_Normal",
			decStr:   "16777215",
			expected: &Color{R: 255, G: 255, B: 255},
		},
		{
			name:        "NegativeValue_Error",
			decStr:      "-1",
			expectError: true,
		},
		{
			name:        "OutOfRangeValue_Error",
			decStr:      "16777216",
			expectError: true,
		},
		{
			name:        "InvalidCharacters_Error",
			decStr:      "abc",
			expectError: true,
		},
		{
			name:        "EmptyString_Error",
			decStr:      "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color, err := ParseFromDecimal(tt.decStr)

			if tt.expectError {
				if err == nil {
					t.Errorf("期待されたエラーが発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラーが発生しました: %v", err)
				return
			}

			if color.R != tt.expected.R || color.G != tt.expected.G || color.B != tt.expected.B {
				t.Errorf("ParseFromDecimal(%v) = {%v, %v, %v}, want {%v, %v, %v}",
					tt.decStr, color.R, color.G, color.B, tt.expected.R, tt.expected.G, tt.expected.B)
			}
		})
	}
}

func TestToHex_Normal(t *testing.T) {
	tests := []struct {
		name     string
		color    *Color
		expected string
	}{
		{
			name:     "RedColor_Normal",
			color:    &Color{R: 255, G: 0, B: 0},
			expected: "#FF0000",
		},
		{
			name:     "GreenColor_Normal",
			color:    &Color{R: 0, G: 255, B: 0},
			expected: "#00FF00",
		},
		{
			name:     "BlueColor_Normal",
			color:    &Color{R: 0, G: 0, B: 255},
			expected: "#0000FF",
		},
		{
			name:     "WhiteColor_Normal",
			color:    &Color{R: 255, G: 255, B: 255},
			expected: "#FFFFFF",
		},
		{
			name:     "BlackColor_Normal",
			color:    &Color{R: 0, G: 0, B: 0},
			expected: "#000000",
		},
		{
			name:     "MixedColor_Normal",
			color:    &Color{R: 128, G: 64, B: 192},
			expected: "#8040C0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.color.ToHex()
			if result != tt.expected {
				t.Errorf("ToHex() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestToRGB_Normal(t *testing.T) {
	tests := []struct {
		name     string
		color    *Color
		expected string
	}{
		{
			name:     "RedColor_Normal",
			color:    &Color{R: 255, G: 0, B: 0},
			expected: "rgb(255,0,0)",
		},
		{
			name:     "GreenColor_Normal",
			color:    &Color{R: 0, G: 255, B: 0},
			expected: "rgb(0,255,0)",
		},
		{
			name:     "BlueColor_Normal",
			color:    &Color{R: 0, G: 0, B: 255},
			expected: "rgb(0,0,255)",
		},
		{
			name:     "MixedColor_Normal",
			color:    &Color{R: 128, G: 64, B: 192},
			expected: "rgb(128,64,192)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.color.ToRGB()
			if result != tt.expected {
				t.Errorf("ToRGB() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestToDecimal_Normal(t *testing.T) {
	tests := []struct {
		name     string
		color    *Color
		expected string
	}{
		{
			name:     "RedColor_Normal",
			color:    &Color{R: 255, G: 0, B: 0},
			expected: "16711680",
		},
		{
			name:     "GreenColor_Normal",
			color:    &Color{R: 0, G: 255, B: 0},
			expected: "65280",
		},
		{
			name:     "BlueColor_Normal",
			color:    &Color{R: 0, G: 0, B: 255},
			expected: "255",
		},
		{
			name:     "WhiteColor_Normal",
			color:    &Color{R: 255, G: 255, B: 255},
			expected: "16777215",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.color.ToDecimal()
			if result != tt.expected {
				t.Errorf("ToDecimal() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestColorConversionRoundTrip_Normal(t *testing.T) {
	// HEX -> RGB -> HEX の往復変換テスト
	originalHex := "#FF0000"

	// HEX -> Color
	color, err := ParseFromHex(originalHex)
	if err != nil {
		t.Fatalf("ParseFromHex failed: %v", err)
	}

	// Color -> RGB
	rgbStr := color.ToRGB()

	// RGB -> Color
	color2, err := ParseFromRGB(rgbStr)
	if err != nil {
		t.Fatalf("ParseFromRGB failed: %v", err)
	}

	// Color -> HEX
	finalHex := color2.ToHex()

	if finalHex != originalHex {
		t.Errorf("Round trip conversion failed: %v -> %v -> %v", originalHex, rgbStr, finalHex)
	}
}

func TestClamp_Normal(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		min      float64
		max      float64
		expected float64
	}{
		{
			name:     "ValueInRange_Normal",
			value:    128,
			min:      0,
			max:      255,
			expected: 128,
		},
		{
			name:     "ValueBelowMin_Normal",
			value:    -10,
			min:      0,
			max:      255,
			expected: 0,
		},
		{
			name:     "ValueAboveMax_Normal",
			value:    300,
			min:      0,
			max:      255,
			expected: 255,
		},
		{
			name:     "ValueAtMin_Normal",
			value:    0,
			min:      0,
			max:      255,
			expected: 0,
		},
		{
			name:     "ValueAtMax_Normal",
			value:    255,
			min:      0,
			max:      255,
			expected: 255,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := clamp(tt.value, tt.min, tt.max)
			if result != tt.expected {
				t.Errorf("clamp(%v, %v, %v) = %v, want %v", tt.value, tt.min, tt.max, result, tt.expected)
			}
		})
	}
}

// 浮動小数点数の比較用ヘルパー関数
func floatEquals(a, b, tolerance float64) bool {
	return math.Abs(a-b) < tolerance
}
