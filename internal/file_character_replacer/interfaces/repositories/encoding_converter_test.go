package repositories

import (
	"testing"

	"github.com/landmaster135/devbox/internal/file_character_replacer/domain"
)

// TestEncodingConverterImpl_ConvertToUTF8_Normal はConvertToUTF8()の正常系をテストします
func TestEncodingConverterImpl_ConvertToUTF8_Normal(t *testing.T) {
	converter := NewEncodingConverter()

	tests := []struct {
		name     string
		content  []byte
		encoding domain.EncodingType
		expected string
	}{
		{
			name:     "UTF-8からUTF-8（変換なし）",
			content:  []byte("Hello, 世界"),
			encoding: domain.EncodingUTF8,
			expected: "Hello, 世界",
		},
		{
			name:     "空のバイト列",
			content:  []byte{},
			encoding: domain.EncodingUTF8,
			expected: "",
		},
		{
			name:     "ASCII文字のみ",
			content:  []byte("Hello World"),
			encoding: domain.EncodingUTF8,
			expected: "Hello World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter.ConvertToUTF8(tt.content, tt.encoding)
			if err != nil {
				t.Errorf("ConvertToUTF8() returned unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("ConvertToUTF8() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestEncodingConverterImpl_ConvertToUTF8_Error はConvertToUTF8()のエラーケースをテストします
func TestEncodingConverterImpl_ConvertToUTF8_Error(t *testing.T) {
	converter := NewEncodingConverter()

	tests := []struct {
		name     string
		content  []byte
		encoding domain.EncodingType
	}{
		{
			name:     "無効なエンコーディング",
			content:  []byte("test"),
			encoding: domain.EncodingType("invalid"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := converter.ConvertToUTF8(tt.content, tt.encoding)
			if err == nil {
				t.Error("ConvertToUTF8() should return error but got nil")
			}
		})
	}
}

// TestEncodingConverterImpl_ConvertFromUTF8_Normal はConvertFromUTF8()の正常系をテストします
func TestEncodingConverterImpl_ConvertFromUTF8_Normal(t *testing.T) {
	converter := NewEncodingConverter()

	tests := []struct {
		name     string
		content  string
		encoding domain.EncodingType
	}{
		{
			name:     "UTF-8からUTF-8（変換なし）",
			content:  "Hello, 世界",
			encoding: domain.EncodingUTF8,
		},
		{
			name:     "空文字列",
			content:  "",
			encoding: domain.EncodingUTF8,
		},
		{
			name:     "ASCII文字のみ",
			content:  "Hello World",
			encoding: domain.EncodingUTF8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter.ConvertFromUTF8(tt.content, tt.encoding)
			if err != nil {
				t.Errorf("ConvertFromUTF8() returned unexpected error: %v", err)
			}
			if string(result) != tt.content {
				t.Errorf("ConvertFromUTF8() = %v, expected %v", string(result), tt.content)
			}
		})
	}
}

// TestEncodingConverterImpl_ConvertFromUTF8_Error はConvertFromUTF8()のエラーケースをテストします
func TestEncodingConverterImpl_ConvertFromUTF8_Error(t *testing.T) {
	converter := NewEncodingConverter()

	tests := []struct {
		name     string
		content  string
		encoding domain.EncodingType
	}{
		{
			name:     "無効なエンコーディング",
			content:  "test",
			encoding: domain.EncodingType("invalid"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := converter.ConvertFromUTF8(tt.content, tt.encoding)
			if err == nil {
				t.Error("ConvertFromUTF8() should return error but got nil")
			}
		})
	}
}

// TestEncodingConverterImpl_DetectEncoding_Normal はDetectEncoding()の正常系をテストします
func TestEncodingConverterImpl_DetectEncoding_Normal(t *testing.T) {
	converter := NewEncodingConverter()

	tests := []struct {
		name     string
		content  []byte
		expected domain.EncodingType
	}{
		{
			name:     "空のバイト列",
			content:  []byte{},
			expected: domain.EncodingUTF8,
		},
		{
			name:     "UTF-8 BOM付き",
			content:  []byte{0xEF, 0xBB, 0xBF, 'H', 'e', 'l', 'l', 'o'},
			expected: domain.EncodingUTF8,
		},
		{
			name:     "ASCII文字のみ",
			content:  []byte("Hello World"),
			expected: domain.EncodingUTF8,
		},
		{
			name:     "UTF-8文字列",
			content:  []byte("Hello, 世界"),
			expected: domain.EncodingUTF8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter.DetectEncoding(tt.content)
			if err != nil {
				t.Errorf("DetectEncoding() returned unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("DetectEncoding() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestEncodingConverterImpl_getEncoder_Normal はgetEncoder()の正常系をテストします
func TestEncodingConverterImpl_getEncoder_Normal(t *testing.T) {
	converter := &EncodingConverterImpl{}

	tests := []struct {
		name     string
		encoding domain.EncodingType
	}{
		{
			name:     "UTF-8",
			encoding: domain.EncodingUTF8,
		},
		{
			name:     "Shift_JIS",
			encoding: domain.EncodingShiftJIS,
		},
		{
			name:     "EUC-JP",
			encoding: domain.EncodingEUCJP,
		},
		{
			name:     "ISO-2022-JP",
			encoding: domain.EncodingISO2022JP,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoder, err := converter.getEncoder(tt.encoding)
			if err != nil {
				t.Errorf("getEncoder() returned unexpected error: %v", err)
			}
			if tt.encoding == domain.EncodingUTF8 {
				if encoder != nil {
					t.Error("getEncoder() should return nil for UTF-8")
				}
			} else {
				if encoder == nil {
					t.Error("getEncoder() should not return nil for non-UTF-8 encoding")
				}
			}
		})
	}
}

// TestEncodingConverterImpl_getEncoder_Error はgetEncoder()のエラーケースをテストします
func TestEncodingConverterImpl_getEncoder_Error(t *testing.T) {
	converter := &EncodingConverterImpl{}

	tests := []struct {
		name     string
		encoding domain.EncodingType
	}{
		{
			name:     "無効なエンコーディング",
			encoding: domain.EncodingType("invalid"),
		},
		{
			name:     "空文字列",
			encoding: domain.EncodingType(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := converter.getEncoder(tt.encoding)
			if err == nil {
				t.Error("getEncoder() should return error but got nil")
			}
		})
	}
}

// TestEncodingConverterImpl_isValidUTF8_Normal はisValidUTF8()をテストします
func TestEncodingConverterImpl_isValidUTF8_Normal(t *testing.T) {
	converter := &EncodingConverterImpl{}

	tests := []struct {
		name     string
		content  []byte
		expected bool
	}{
		{
			name:     "UTF-8 BOM付き",
			content:  []byte{0xEF, 0xBB, 0xBF, 'H', 'e', 'l', 'l', 'o'},
			expected: true,
		},
		{
			name:     "有効なUTF-8文字列",
			content:  []byte("Hello, 世界"),
			expected: true,
		},
		{
			name:     "ASCII文字のみ",
			content:  []byte("Hello World"),
			expected: true,
		},
		{
			name:     "空のバイト列",
			content:  []byte{},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.isValidUTF8(tt.content)
			if result != tt.expected {
				t.Errorf("isValidUTF8() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestEncodingConverterImpl_isValidShiftJIS_Normal はisValidShiftJIS()をテストします
func TestEncodingConverterImpl_isValidShiftJIS_Normal(t *testing.T) {
	converter := &EncodingConverterImpl{}

	tests := []struct {
		name     string
		content  []byte
		expected bool
	}{
		{
			name:     "ASCII文字のみ",
			content:  []byte("Hello World"),
			expected: true,
		},
		{
			name:     "空のバイト列",
			content:  []byte{},
			expected: true,
		},
		{
			name:     "制御文字（タブ）",
			content:  []byte{0x09},
			expected: true,
		},
		{
			name:     "制御文字（改行）",
			content:  []byte{0x0A},
			expected: true,
		},
		{
			name:     "制御文字（復帰）",
			content:  []byte{0x0D},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.isValidShiftJIS(tt.content)
			if result != tt.expected {
				t.Errorf("isValidShiftJIS() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestEncodingConverterImpl_isValidEUCJP_Normal はisValidEUCJP()をテストします
func TestEncodingConverterImpl_isValidEUCJP_Normal(t *testing.T) {
	converter := &EncodingConverterImpl{}

	tests := []struct {
		name     string
		content  []byte
		expected bool
	}{
		{
			name:     "ASCII文字のみ",
			content:  []byte("Hello World"),
			expected: true,
		},
		{
			name:     "空のバイト列",
			content:  []byte{},
			expected: true,
		},
		{
			name:     "制御文字（タブ）",
			content:  []byte{0x09},
			expected: true,
		},
		{
			name:     "制御文字（改行）",
			content:  []byte{0x0A},
			expected: true,
		},
		{
			name:     "制御文字（復帰）",
			content:  []byte{0x0D},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.isValidEUCJP(tt.content)
			if result != tt.expected {
				t.Errorf("isValidEUCJP() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestEncodingConverterImpl_isValidISO2022JP_Normal はisValidISO2022JP()をテストします
func TestEncodingConverterImpl_isValidISO2022JP_Normal(t *testing.T) {
	converter := &EncodingConverterImpl{}

	tests := []struct {
		name     string
		content  []byte
		expected bool
	}{
		{
			name:     "ASCII文字のみ",
			content:  []byte("Hello World"),
			expected: false, // エスケープシーケンスがないため
		},
		{
			name:     "空のバイト列",
			content:  []byte{},
			expected: false,
		},
		{
			name:     "ASCIIエスケープシーケンス",
			content:  []byte("\x1b(BHello"),
			expected: true,
		},
		{
			name:     "JIS X 0208エスケープシーケンス",
			content:  []byte("\x1b$BHello"),
			expected: true,
		},
		{
			name:     "JIS X 0208エスケープシーケンス（旧）",
			content:  []byte("\x1b$@Hello"),
			expected: true,
		},
		{
			name:     "JIS X 0201カナエスケープシーケンス",
			content:  []byte("\x1b(IHello"),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.isValidISO2022JP(tt.content)
			if result != tt.expected {
				t.Errorf("isValidISO2022JP() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
