package usecases

import (
	"golang.org/x/text/unicode/norm"
	"golang.org/x/text/width"
)

// ConversionMode は変換モードを表す型です
type ConversionMode string

const (
	// FullWidthMode は全角変換モードを表します
	FullWidthMode ConversionMode = "full"
	// HalfWidthMode は半角変換モードを表します
	HalfWidthMode ConversionMode = "half"
)

// KanaConverter はカナ文字の変換に関する機能を提供します
type KanaConverter struct {
	Mode ConversionMode
}

// NewKanaConverter は新しいKanaConverterインスタンスを作成します
func NewKanaConverter(mode ConversionMode) *KanaConverter {
	return &KanaConverter{
		Mode: mode,
	}
}

// Convert は入力文字列を指定されたモードに基づいて変換します
func (k *KanaConverter) Convert(input string) string {
	switch k.Mode {
	case FullWidthMode:
		return k.ToFullWidth(input)
	case HalfWidthMode:
		return k.ToHalfWidth(input)
	default:
		// デフォルトは全角変換
		return k.ToFullWidth(input)
	}
}

// ToFullWidth は入力文字列を全角カナに変換します
func (k *KanaConverter) ToFullWidth(input string) string {
	return norm.NFKC.String(input)
}

// ToHalfWidth は入力文字列を半角カナに変換します
func (k *KanaConverter) ToHalfWidth(input string) string {
	// まず全角に正規化してから半角に変換
	fullWidth := norm.NFKC.String(input)
	return width.Narrow.String(fullWidth)
}

// IsValidMode はモード文字列が有効かどうかを判定します
func IsValidMode(mode string) bool {
	return mode == string(FullWidthMode) || mode == string(HalfWidthMode)
}

// GetModeFromString は文字列からConversionModeを取得します
func GetModeFromString(mode string) ConversionMode {
	if mode == string(HalfWidthMode) {
		return HalfWidthMode
	}
	return FullWidthMode
}

// GetModeDescription はモードの説明文を取得します
func GetModeDescription(mode ConversionMode) string {
	switch mode {
	case FullWidthMode:
		return "full-width"
	case HalfWidthMode:
		return "half-width"
	default:
		return "unknown"
	}
}
