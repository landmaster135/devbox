package usecases

import (
	"fmt"
	"testing"

	domain "github.com/landmaster135/devbox/internal/color_code_converter/domain"
)

func TestNewColorConverterService_Normal(t *testing.T) {
	service := NewColorConverterService()
	if service == nil {
		t.Errorf("NewColorConverterService() returned nil")
	}
}

func TestConvertColor_Normal(t *testing.T) {
	const (
		testHexValue = "#FF0000"
		testRgbValue = "rgb(255,0,0)"
		testHslValue = "hsl(0,100%,50%)"
		testHsvValue = "hsv(0,100%,100%)"
		testDecValue = "16711680"
	)

	tests := []struct {
		name        string
		srcFormat   string
		destFormat  string
		value       string
		expected    string
		expectError bool
	}{
		{
			name:       "HexToRgb_Normal",
			srcFormat:  "hex",
			destFormat: "rgb",
			value:      testHexValue,
			expected:   testRgbValue,
		},
		{
			name:       "RgbToHex_Normal",
			srcFormat:  "rgb",
			destFormat: "hex",
			value:      testRgbValue,
			expected:   testHexValue,
		},
		{
			name:       "HexToHsl_Normal",
			srcFormat:  "hex",
			destFormat: "hsl",
			value:      testHexValue,
			expected:   testHslValue,
		},
		{
			name:       "HexToHsv_Normal",
			srcFormat:  "hex",
			destFormat: "hsv",
			value:      testHexValue,
			expected:   testHsvValue,
		},
		{
			name:       "SameFormat_Normal",
			srcFormat:  "hex",
			destFormat: "hex",
			value:      testHexValue,
			expected:   testHexValue,
		},
		{
			name:       "HexToDec_Normal",
			srcFormat:  "hex",
			destFormat: "dec",
			value:      testHexValue,
			expected:   testDecValue,
		},
		{
			name:       "DecToHex_Normal",
			srcFormat:  "dec",
			destFormat: "hex",
			value:      testDecValue,
			expected:   testHexValue,
		},
		{
			name:        "EmptySrcFormat_Error",
			srcFormat:   "",
			destFormat:  "rgb",
			value:       testHexValue,
			expectError: true,
		},
		{
			name:        "EmptyDestFormat_Error",
			srcFormat:   "hex",
			destFormat:  "",
			value:       testHexValue,
			expectError: true,
		},
		{
			name:        "EmptyValue_Error",
			srcFormat:   "hex",
			destFormat:  "rgb",
			value:       "",
			expectError: true,
		},
		{
			name:        "InvalidSrcFormat_Error",
			srcFormat:   "invalid",
			destFormat:  "rgb",
			value:       testHexValue,
			expectError: true,
		},
		{
			name:        "InvalidDestFormat_Error",
			srcFormat:   "hex",
			destFormat:  "invalid",
			value:       testHexValue,
			expectError: true,
		},
		{
			name:        "InvalidValue_Error",
			srcFormat:   "hex",
			destFormat:  "rgb",
			value:       "invalid",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewColorConverterService()
			result, err := service.ConvertColor(tt.srcFormat, tt.destFormat, tt.value)

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

			if result != tt.expected {
				t.Errorf("ConvertColor(%v, %v, %v) = %v, want %v", tt.srcFormat, tt.destFormat, tt.value, result, tt.expected)
			}
		})
	}
}

func TestParseColorFromFormat_Normal(t *testing.T) {
	const (
		testHexValue = "#FF0000"
		testRgbValue = "rgb(255,0,0)"
		testHslValue = "hsl(0,100%,50%)"
		testHsvValue = "hsv(0,100%,100%)"
		testDecValue = "16711680"
	)

	tests := []struct {
		name        string
		format      string
		value       string
		expectError bool
	}{
		{
			name:   "ValidHex_Normal",
			format: "hex",
			value:  testHexValue,
		},
		{
			name:   "ValidRgb_Normal",
			format: "rgb",
			value:  testRgbValue,
		},
		{
			name:   "ValidHsl_Normal",
			format: "hsl",
			value:  testHslValue,
		},
		{
			name:   "ValidHsv_Normal",
			format: "hsv",
			value:  testHsvValue,
		},
		{
			name:   "ValidDec_Normal",
			format: "dec",
			value:  testDecValue,
		},
		{
			name:        "InvalidFormat_Error",
			format:      "invalid",
			value:       testHexValue,
			expectError: true,
		},
		{
			name:        "InvalidHexValue_Error",
			format:      "hex",
			value:       "invalid",
			expectError: true,
		},
		{
			name:        "InvalidRgbValue_Error",
			format:      "rgb",
			value:       "invalid",
			expectError: true,
		},
		{
			name:        "InvalidDecValue_Error",
			format:      "dec",
			value:       "invalid",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewColorConverterService()
			color, err := service.parseColorFromFormat(tt.format, tt.value)

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

			if color == nil {
				t.Errorf("parseColorFromFormat returned nil color")
			}

			// RGB値が0-255の範囲内にあることを確認
			if color.R < 0 || color.R > 255 || color.G < 0 || color.G > 255 || color.B < 0 || color.B > 255 {
				t.Errorf("parseColorFromFormat returned invalid RGB values: {%v, %v, %v}", color.R, color.G, color.B)
			}
		})
	}
}

func TestFormatColorToFormat_Normal(t *testing.T) {
	testColor := &domain.Color{R: 255, G: 0, B: 0}

	tests := []struct {
		name        string
		format      string
		expected    string
		expectError bool
	}{
		{
			name:     "ValidHex_Normal",
			format:   "hex",
			expected: "#FF0000",
		},
		{
			name:     "ValidRgb_Normal",
			format:   "rgb",
			expected: "rgb(255,0,0)",
		},
		{
			name:     "ValidHsl_Normal",
			format:   "hsl",
			expected: "hsl(0,100%,50%)",
		},
		{
			name:     "ValidHsv_Normal",
			format:   "hsv",
			expected: "hsv(0,100%,100%)",
		},
		{
			name:     "ValidDec_Normal",
			format:   "dec",
			expected: "16711680",
		},
		{
			name:        "InvalidFormat_Error",
			format:      "invalid",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewColorConverterService()
			result, err := service.formatColorToFormat(tt.format, testColor)

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

			if result != tt.expected {
				t.Errorf("formatColorToFormat(%v, color) = %v, want %v", tt.format, result, tt.expected)
			}
		})
	}
}

func TestValidateColorFormat_Normal(t *testing.T) {
	const (
		testHexValue = "#FF0000"
		testRgbValue = "rgb(255,0,0)"
		testHslValue = "hsl(0,100%,50%)"
		testHsvValue = "hsv(0,100%,100%)"
		testDecValue = "16711680"
	)

	tests := []struct {
		name        string
		format      string
		value       string
		expectError bool
	}{
		{
			name:   "ValidHex_Normal",
			format: "hex",
			value:  testHexValue,
		},
		{
			name:   "ValidRgb_Normal",
			format: "rgb",
			value:  testRgbValue,
		},
		{
			name:   "ValidHsl_Normal",
			format: "hsl",
			value:  testHslValue,
		},
		{
			name:   "ValidHsv_Normal",
			format: "hsv",
			value:  testHsvValue,
		},
		{
			name:   "ValidDec_Normal",
			format: "dec",
			value:  testDecValue,
		},
		{
			name:        "InvalidHexValue_Error",
			format:      "hex",
			value:       "invalid",
			expectError: true,
		},
		{
			name:        "InvalidRgbValue_Error",
			format:      "rgb",
			value:       "invalid",
			expectError: true,
		},
		{
			name:        "InvalidDecValue_Error",
			format:      "dec",
			value:       "99999999",
			expectError: true,
		},
		{
			name:        "InvalidFormat_Error",
			format:      "invalid",
			value:       testHexValue,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewColorConverterService()
			err := service.ValidateColorFormat(tt.format, tt.value)

			if tt.expectError {
				if err == nil {
					t.Errorf("期待されたエラーが発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラーが発生しました: %v", err)
			}
		})
	}
}

func TestGetSupportedFormats_Normal(t *testing.T) {
	service := NewColorConverterService()
	formats := service.GetSupportedFormats()

	expectedFormats := []string{"hex", "rgb", "hsl", "hsv", "dec"}

	if len(formats) != len(expectedFormats) {
		t.Errorf("GetSupportedFormats() returned %d formats, want %d", len(formats), len(expectedFormats))
		return
	}

	for i, format := range formats {
		if format != expectedFormats[i] {
			t.Errorf("GetSupportedFormats()[%d] = %v, want %v", i, format, expectedFormats[i])
		}
	}
}

func TestConvertColorWithValidation_Normal(t *testing.T) {
	const (
		testHexValue = "#FF0000"
		testRgbValue = "rgb(255,0,0)"
		testDecValue = "16711680"
	)

	tests := []struct {
		name        string
		srcFormat   string
		destFormat  string
		value       string
		expected    string
		expectError bool
	}{
		{
			name:       "ValidConversion_Normal",
			srcFormat:  "hex",
			destFormat: "rgb",
			value:      testHexValue,
			expected:   testRgbValue,
		},
		{
			name:       "ValidDecConversion_Normal",
			srcFormat:  "dec",
			destFormat: "hex",
			value:      testDecValue,
			expected:   testHexValue,
		},
		{
			name:        "UnsupportedSrcFormat_Error",
			srcFormat:   "invalid",
			destFormat:  "rgb",
			value:       testHexValue,
			expectError: true,
		},
		{
			name:        "UnsupportedDestFormat_Error",
			srcFormat:   "hex",
			destFormat:  "invalid",
			value:       testHexValue,
			expectError: true,
		},
		{
			name:        "InvalidValue_Error",
			srcFormat:   "hex",
			destFormat:  "rgb",
			value:       "invalid",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewColorConverterService()
			result, err := service.ConvertColorWithValidation(tt.srcFormat, tt.destFormat, tt.value)

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

			if result != tt.expected {
				t.Errorf("ConvertColorWithValidation(%v, %v, %v) = %v, want %v", tt.srcFormat, tt.destFormat, tt.value, result, tt.expected)
			}
		})
	}
}

func TestIsFormatSupported_Normal(t *testing.T) {
	supportedFormats := []string{"hex", "rgb", "hsl", "hsv", "dec"}

	tests := []struct {
		name     string
		format   string
		expected bool
	}{
		{
			name:     "SupportedHex_Normal",
			format:   "hex",
			expected: true,
		},
		{
			name:     "SupportedRgb_Normal",
			format:   "rgb",
			expected: true,
		},
		{
			name:     "SupportedHsl_Normal",
			format:   "hsl",
			expected: true,
		},
		{
			name:     "SupportedHsv_Normal",
			format:   "hsv",
			expected: true,
		},
		{
			name:     "SupportedDec_Normal",
			format:   "dec",
			expected: true,
		},
		{
			name:     "UnsupportedFormat_Normal",
			format:   "invalid",
			expected: false,
		},
		{
			name:     "EmptyFormat_Normal",
			format:   "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewColorConverterService()
			result := service.isFormatSupported(tt.format, supportedFormats)

			if result != tt.expected {
				t.Errorf("isFormatSupported(%v) = %v, want %v", tt.format, result, tt.expected)
			}
		})
	}
}

func TestColorConversionIntegration_Normal(t *testing.T) {
	// 統合テスト：複数の形式間での変換を確認
	service := NewColorConverterService()

	// テストケース：赤色の変換
	testCases := []struct {
		format string
		value  string
	}{
		{"hex", "#FF0000"},
		{"rgb", "rgb(255,0,0)"},
		{"hsl", "hsl(0,100%,50%)"},
		{"hsv", "hsv(0,100%,100%)"},
		{"dec", "16711680"},
	}

	// 各形式から他の全ての形式への変換をテスト
	for _, src := range testCases {
		for _, dest := range testCases {
			t.Run(fmt.Sprintf("%s_to_%s", src.format, dest.format), func(t *testing.T) {
				result, err := service.ConvertColorWithValidation(src.format, dest.format, src.value)
				if err != nil {
					t.Errorf("変換に失敗しました: %s -> %s: %v", src.format, dest.format, err)
					return
				}

				// 結果が空でないことを確認
				if result == "" {
					t.Errorf("変換結果が空です: %s -> %s", src.format, dest.format)
				}

				// 変換結果が有効な形式であることを確認
				err = service.ValidateColorFormat(dest.format, result)
				if err != nil {
					t.Errorf("変換結果が無効です: %s -> %s, result: %s, error: %v", src.format, dest.format, result, err)
				}
			})
		}
	}
}
