package repositories

import (
	"bytes"
	"fmt"
	"io"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/landmaster135/devbox/internal/file_character_replacer/domain"
)

// EncodingConverterImpl はEncodingConverterインターフェースの具象実装です
type EncodingConverterImpl struct{}

// NewEncodingConverter は新しいEncodingConverterImplを作成します
func NewEncodingConverter() domain.EncodingConverter {
	return &EncodingConverterImpl{}
}

// getEncoder は指定されたエンコーディングタイプに対応するエンコーダーを取得します
func (c *EncodingConverterImpl) getEncoder(encodingType domain.EncodingType) (encoding.Encoding, error) {
	switch encodingType {
	case domain.EncodingUTF8:
		return nil, nil // UTF-8は変換不要
	case domain.EncodingShiftJIS:
		return japanese.ShiftJIS, nil
	case domain.EncodingEUCJP:
		return japanese.EUCJP, nil
	case domain.EncodingISO2022JP:
		return japanese.ISO2022JP, nil
	default:
		return nil, fmt.Errorf("サポートされていないエンコーディングです: %s", encodingType)
	}
}

// ConvertToUTF8 は指定されたエンコーディングのバイト列をUTF-8文字列に変換します
func (c *EncodingConverterImpl) ConvertToUTF8(content []byte, encodingType domain.EncodingType) (string, error) {
	// UTF-8の場合はそのまま返す
	if encodingType == domain.EncodingUTF8 {
		return string(content), nil
	}

	enc, err := c.getEncoder(encodingType)
	if err != nil {
		return "", err
	}

	// デコーダーを作成
	decoder := enc.NewDecoder()
	reader := transform.NewReader(bytes.NewReader(content), decoder)

	// UTF-8に変換
	utf8Content, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("エンコーディング変換に失敗しました (%s -> UTF-8): %w", encodingType, err)
	}

	return string(utf8Content), nil
}

// ConvertFromUTF8 はUTF-8文字列を指定されたエンコーディングのバイト列に変換します
func (c *EncodingConverterImpl) ConvertFromUTF8(content string, encodingType domain.EncodingType) ([]byte, error) {
	// UTF-8の場合はそのまま返す
	if encodingType == domain.EncodingUTF8 {
		return []byte(content), nil
	}

	enc, err := c.getEncoder(encodingType)
	if err != nil {
		return nil, err
	}

	// エンコーダーを作成
	encoder := enc.NewEncoder()
	reader := transform.NewReader(bytes.NewReader([]byte(content)), encoder)

	// 指定されたエンコーディングに変換
	encodedContent, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("エンコーディング変換に失敗しました (UTF-8 -> %s): %w", encodingType, err)
	}

	return encodedContent, nil
}

// DetectEncoding はバイト列から文字エンコーディングを推測します
func (c *EncodingConverterImpl) DetectEncoding(content []byte) (domain.EncodingType, error) {
	if len(content) == 0 {
		return domain.EncodingUTF8, nil
	}

	// UTF-8の検証
	if c.isValidUTF8(content) {
		return domain.EncodingUTF8, nil
	}

	// Shift_JISの検証
	if c.isValidShiftJIS(content) {
		return domain.EncodingShiftJIS, nil
	}

	// EUC-JPの検証
	if c.isValidEUCJP(content) {
		return domain.EncodingEUCJP, nil
	}

	// ISO-2022-JPの検証
	if c.isValidISO2022JP(content) {
		return domain.EncodingISO2022JP, nil
	}

	// デフォルトはUTF-8
	return domain.EncodingUTF8, nil
}

// isValidUTF8 はバイト列が有効なUTF-8かどうかを確認します
func (c *EncodingConverterImpl) isValidUTF8(content []byte) bool {
	// UTF-8のBOMをチェック
	if len(content) >= 3 && content[0] == 0xEF && content[1] == 0xBB && content[2] == 0xBF {
		return true
	}

	// UTF-8として変換を試行
	_, err := c.ConvertToUTF8(content, domain.EncodingUTF8)
	return err == nil
}

// isValidShiftJIS はバイト列が有効なShift_JISかどうかを確認します
func (c *EncodingConverterImpl) isValidShiftJIS(content []byte) bool {
	// Shift_JISとして変換を試行
	_, err := c.ConvertToUTF8(content, domain.EncodingShiftJIS)
	if err != nil {
		return false
	}

	// Shift_JISの特徴的なバイト範囲をチェック
	for i := 0; i < len(content); i++ {
		b := content[i]
		// 1バイト文字（ASCII + 半角カナ）
		if (b >= 0x20 && b <= 0x7E) || (b >= 0xA1 && b <= 0xDF) {
			continue
		}
		// 2バイト文字の1バイト目
		if (b >= 0x81 && b <= 0x9F) || (b >= 0xE0 && b <= 0xFC) {
			if i+1 < len(content) {
				b2 := content[i+1]
				// 2バイト文字の2バイト目
				if (b2 >= 0x40 && b2 <= 0x7E) || (b2 >= 0x80 && b2 <= 0xFC) {
					i++ // 2バイト目をスキップ
					continue
				}
			}
			return false
		}
		// 制御文字など
		if b < 0x20 && b != 0x09 && b != 0x0A && b != 0x0D {
			continue
		}
	}
	return true
}

// isValidEUCJP はバイト列が有効なEUC-JPかどうかを確認します
func (c *EncodingConverterImpl) isValidEUCJP(content []byte) bool {
	// EUC-JPとして変換を試行
	_, err := c.ConvertToUTF8(content, domain.EncodingEUCJP)
	if err != nil {
		return false
	}

	// EUC-JPの特徴的なバイト範囲をチェック
	for i := 0; i < len(content); i++ {
		b := content[i]
		// ASCII文字
		if b >= 0x20 && b <= 0x7E {
			continue
		}
		// 漢字（2バイト）
		if b >= 0xA1 && b <= 0xFE {
			if i+1 < len(content) {
				b2 := content[i+1]
				if b2 >= 0xA1 && b2 <= 0xFE {
					i++ // 2バイト目をスキップ
					continue
				}
			}
			return false
		}
		// 半角カナ（2バイト）
		if b == 0x8E {
			if i+1 < len(content) {
				b2 := content[i+1]
				if b2 >= 0xA1 && b2 <= 0xDF {
					i++ // 2バイト目をスキップ
					continue
				}
			}
			return false
		}
		// 制御文字など
		if b < 0x20 && b != 0x09 && b != 0x0A && b != 0x0D {
			continue
		}
	}
	return true
}

// isValidISO2022JP はバイト列が有効なISO-2022-JPかどうかを確認します
func (c *EncodingConverterImpl) isValidISO2022JP(content []byte) bool {
	// ISO-2022-JPとして変換を試行
	_, err := c.ConvertToUTF8(content, domain.EncodingISO2022JP)
	if err != nil {
		return false
	}

	// ISO-2022-JPのエスケープシーケンスをチェック
	contentStr := string(content)
	// ASCII
	if bytes.Contains(content, []byte("\x1b(B")) {
		return true
	}
	// JIS X 0208
	if bytes.Contains(content, []byte("\x1b$B")) || bytes.Contains(content, []byte("\x1b$@")) {
		return true
	}
	// JIS X 0201 カナ
	if bytes.Contains(content, []byte("\x1b(I")) {
		return true
	}

	// エスケープシーケンスが含まれていない場合はASCIIのみと判断
	for _, b := range content {
		if b >= 0x80 {
			return false
		}
	}

	// ASCIIのみの場合は可能性があるが、他の判定を優先
	_ = contentStr
	return false
}
