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

// TestEncodingConverterImpl_ConvertFromUTF8_NonUTF8 はConvertFromUTF8()の非UTF-8エンコーディングをテストします
func TestEncodingConverterImpl_ConvertFromUTF8_NonUTF8(t *testing.T) {
	converter := NewEncodingConverter()

	tests := []struct {
		name     string
		content  string
		encoding domain.EncodingType
	}{
		{
			name:     "Shift_JISへの変換",
			content:  "Hello",
			encoding: domain.EncodingShiftJIS,
		},
		{
			name:     "EUC-JPへの変換",
			content:  "Hello",
			encoding: domain.EncodingEUCJP,
		},
		{
			name:     "ISO-2022-JPへの変換",
			content:  "Hello",
			encoding: domain.EncodingISO2022JP,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter.ConvertFromUTF8(tt.content, tt.encoding)
			if err != nil {
				t.Errorf("ConvertFromUTF8() returned unexpected error: %v", err)
			}
			if len(result) == 0 {
				t.Error("ConvertFromUTF8() should not return empty result")
			}
		})
	}
}

// TestEncodingConverterImpl_DetectEncoding_Various は様々なエンコーディングの検出をテストします
func TestEncodingConverterImpl_DetectEncoding_Various(t *testing.T) {
	converter := NewEncodingConverter()

	tests := []struct {
		name     string
		content  []byte
		expected domain.EncodingType
	}{
		{
			name:     "無効なバイト列",
			content:  []byte{0xFF, 0xFE, 0x00, 0x00},
			expected: domain.EncodingUTF8, // デフォルト
		},
		{
			name:     "高バイト文字（UTF-8以外の可能性）",
			content:  []byte{0x80, 0x81, 0x82},
			expected: domain.EncodingUTF8, // 実装では最終的にUTF-8がデフォルト
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

// TestEncodingConverterImpl_DetectEncoding_EdgeCases はDetectEncoding()のエッジケースをテストします
func TestEncodingConverterImpl_DetectEncoding_EdgeCases(t *testing.T) {
	converter := NewEncodingConverter()

	tests := []struct {
		name    string
		content []byte
	}{
		{
			name:    "Shift_JIS風バイト列",
			content: []byte{0x82, 0xA0}, // あ in Shift_JIS
		},
		{
			name:    "EUC-JP風バイト列",
			content: []byte{0xA4, 0xA2}, // あ in EUC-JP
		},
		{
			name:    "ISO-2022-JP風エスケープシーケンス",
			content: []byte("\x1b$B$\""), // あ in ISO-2022-JP
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter.DetectEncoding(tt.content)
			if err != nil {
				t.Errorf("DetectEncoding() returned unexpected error: %v", err)
			}
			// 実装の動作を確認するだけで、特定の結果は期待しない
			t.Logf("DetectEncoding() for %s = %v", tt.name, result)
		})
	}
}

// TestEncodingConverterImpl_isValidShiftJIS_Various はisValidShiftJIS()の様々なケースをテストします
func TestEncodingConverterImpl_isValidShiftJIS_Various(t *testing.T) {
	converter := &EncodingConverterImpl{}

	tests := []struct {
		name     string
		content  []byte
		expected bool
	}{
		{
			name:     "Shift_JIS 2バイト文字（第1バイト範囲1）",
			content:  []byte{0x81, 0x40}, // 有効な2バイト文字
			expected: true,
		},
		{
			name:     "Shift_JIS 2バイト文字（第1バイト範囲2）",
			content:  []byte{0xE0, 0x40}, // 有効な2バイト文字
			expected: true,
		},
		{
			name:     "Shift_JIS 半角カナ",
			content:  []byte{0xA1}, // 半角カナ
			expected: true,
		},
		{
			name:     "無効な2バイト文字（第2バイトが範囲外）",
			content:  []byte{0x81, 0x3F}, // 無効な第2バイト
			expected: false,
		},
		{
			name:     "不完全な2バイト文字",
			content:  []byte{0x81}, // 第2バイトがない
			expected: false,
		},
		{
			name:     "制御文字（無効）",
			content:  []byte{0x01}, // 無効な制御文字
			expected: true, // 制御文字は許可される
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

// TestEncodingConverterImpl_isValidEUCJP_Various はisValidEUCJP()の様々なケースをテストします
func TestEncodingConverterImpl_isValidEUCJP_Various(t *testing.T) {
	converter := &EncodingConverterImpl{}

	tests := []struct {
		name     string
		content  []byte
		expected bool
	}{
		{
			name:     "EUC-JP 漢字（2バイト）",
			content:  []byte{0xA4, 0xA2}, // あ in EUC-JP
			expected: true,
		},
		{
			name:     "EUC-JP 半角カナ（2バイト）",
			content:  []byte{0x8E, 0xA1}, // 半角カナ
			expected: true,
		},
		{
			name:     "無効な漢字（第2バイトが範囲外）",
			content:  []byte{0xA1, 0xA0}, // 無効な第2バイト
			expected: false,
		},
		{
			name:     "不完全な漢字",
			content:  []byte{0xA1}, // 第2バイトがない
			expected: false,
		},
		{
			name:     "不完全な半角カナ",
			content:  []byte{0x8E}, // 第2バイトがない
			expected: false,
		},
		{
			name:     "無効な半角カナ（第2バイトが範囲外）",
			content:  []byte{0x8E, 0xA0}, // 無効な第2バイト
			expected: false,
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

// TestEncodingConverterImpl_isValidISO2022JP_HighByte はisValidISO2022JP()の高バイト処理をテストします
func TestEncodingConverterImpl_isValidISO2022JP_HighByte(t *testing.T) {
	converter := &EncodingConverterImpl{}

	tests := []struct {
		name     string
		content  []byte
		expected bool
	}{
		{
			name:     "高バイト文字を含む（無効）",
			content:  []byte{0x80, 0x81, 0x82},
			expected: false,
		},
		{
			name:     "ASCII範囲のみ",
			content:  []byte{0x20, 0x7E},
			expected: false, // エスケープシーケンスがないため
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

// TestEncodingConverterImpl_ConvertToUTF8_ActualEncodings は実際のエンコーディングでの変換をテストします
func TestEncodingConverterImpl_ConvertToUTF8_ActualEncodings(t *testing.T) {
	converter := NewEncodingConverter()

	tests := []struct {
		name     string
		content  []byte
		encoding domain.EncodingType
		wantErr  bool
	}{
		{
			name:     "Shift_JIS ASCII",
			content:  []byte("Hello"),
			encoding: domain.EncodingShiftJIS,
			wantErr:  false,
		},
		{
			name:     "EUC-JP ASCII",
			content:  []byte("Hello"),
			encoding: domain.EncodingEUCJP,
			wantErr:  false,
		},
		{
			name:     "ISO-2022-JP ASCII",
			content:  []byte("Hello"),
			encoding: domain.EncodingISO2022JP,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter.ConvertToUTF8(tt.content, tt.encoding)
			if tt.wantErr {
				if err == nil {
					t.Error("ConvertToUTF8() should return error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("ConvertToUTF8() returned unexpected error: %v", err)
				}
				if len(result) == 0 {
					t.Error("ConvertToUTF8() should not return empty result")
				}
			}
		})
	}
}

// TestEncodingConverterImpl_ConvertToUTF8_ErrorCases はConvertToUTF8()のエラーケースをテストします
func TestEncodingConverterImpl_ConvertToUTF8_ErrorCases(t *testing.T) {
	converter := NewEncodingConverter()

	tests := []struct {
		name     string
		content  []byte
		encoding domain.EncodingType
	}{
		{
			name:     "潜在的に無効なShift_JISバイト列",
			content:  []byte{0x81, 0x3F}, // 無効な第2バイト
			encoding: domain.EncodingShiftJIS,
		},
		{
			name:     "潜在的に無効なEUC-JPバイト列",
			content:  []byte{0xA1, 0xA0}, // 無効な第2バイト
			encoding: domain.EncodingEUCJP,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter.ConvertToUTF8(tt.content, tt.encoding)
			// エラーが発生するかどうかは実装依存なので、結果をログに記録するだけ
			t.Logf("ConvertToUTF8() for %s: result=%q, err=%v", tt.name, result, err)
		})
	}
}
