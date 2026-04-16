package usecases

import (
	"fmt"

	domain "github.com/landmaster135/devbox/internal/color_code_converter/domain"
)

// ColorConverterService はカラーコード変換のビジネスロジックを提供する
type ColorConverterService struct{}

// NewColorConverterService は新しいColorConverterServiceを作成する
func NewColorConverterService() *ColorConverterService {
	return &ColorConverterService{}
}

// ConvertColor はカラーコードを指定された形式に変換する
func (s *ColorConverterService) ConvertColor(srcFormat, destFormat, value string) (string, error) {
	// 入力値の検証
	if srcFormat == "" {
		return "", fmt.Errorf("変換元形式が指定されていません")
	}
	if destFormat == "" {
		return "", fmt.Errorf("変換先形式が指定されていません")
	}
	if value == "" {
		return "", fmt.Errorf("変換するカラーコード値が指定されていません")
	}

	// 変換元形式に応じてColorオブジェクトを作成
	color, err := s.parseColorFromFormat(srcFormat, value)
	if err != nil {
		return "", fmt.Errorf("カラーコードの解析に失敗しました: %v", err)
	}

	// 変換先形式に応じて文字列に変換
	result, err := s.formatColorToFormat(destFormat, color)
	if err != nil {
		return "", fmt.Errorf("カラーコードの変換に失敗しました: %v", err)
	}

	return result, nil
}

// parseColorFromFormat は指定された形式の文字列からColorオブジェクトを作成する
func (s *ColorConverterService) parseColorFromFormat(format, value string) (*domain.Color, error) {
	switch format {
	case "hex":
		return domain.ParseFromHex(value)
	case "rgb":
		return domain.ParseFromRGB(value)
	case "hsl":
		return domain.ParseFromHSL(value)
	case "hsv":
		return domain.ParseFromHSV(value)
	case "dec":
		return domain.ParseFromDecimal(value)
	default:
		return nil, fmt.Errorf("サポートされていない変換元形式です: %s", format)
	}
}

// formatColorToFormat はColorオブジェクトを指定された形式の文字列に変換する
func (s *ColorConverterService) formatColorToFormat(format string, color *domain.Color) (string, error) {
	switch format {
	case "hex":
		return color.ToHex(), nil
	case "rgb":
		return color.ToRGB(), nil
	case "hsl":
		return color.ToHSL(), nil
	case "hsv":
		return color.ToHSV(), nil
	case "dec":
		return color.ToDecimal(), nil
	default:
		return "", fmt.Errorf("サポートされていない変換先形式です: %s", format)
	}
}

// ValidateColorFormat はカラーコード形式が有効かどうかを検証する
func (s *ColorConverterService) ValidateColorFormat(format, value string) error {
	_, err := s.parseColorFromFormat(format, value)
	return err
}

// GetSupportedFormats はサポートされているカラーコード形式のリストを返す
func (s *ColorConverterService) GetSupportedFormats() []string {
	return []string{"hex", "rgb", "hsl", "hsv", "dec"}
}

// ConvertColorWithValidation はバリデーション付きでカラーコードを変換する
func (s *ColorConverterService) ConvertColorWithValidation(srcFormat, destFormat, value string) (string, error) {
	// サポートされている形式かチェック
	supportedFormats := s.GetSupportedFormats()

	if !s.isFormatSupported(srcFormat, supportedFormats) {
		return "", fmt.Errorf("サポートされていない変換元形式です: %s (サポート形式: %v)", srcFormat, supportedFormats)
	}

	if !s.isFormatSupported(destFormat, supportedFormats) {
		return "", fmt.Errorf("サポートされていない変換先形式です: %s (サポート形式: %v)", destFormat, supportedFormats)
	}

	// 入力値の形式検証
	if err := s.ValidateColorFormat(srcFormat, value); err != nil {
		return "", fmt.Errorf("入力値の形式が正しくありません: %v", err)
	}

	// 変換実行
	return s.ConvertColor(srcFormat, destFormat, value)
}

// isFormatSupported は指定された形式がサポートされているかチェックする
func (s *ColorConverterService) isFormatSupported(format string, supportedFormats []string) bool {
	for _, supported := range supportedFormats {
		if format == supported {
			return true
		}
	}
	return false
}
