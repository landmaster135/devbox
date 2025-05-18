package usecases

import (
	"strings"

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
	// RemoveVoicedSoundMode は濁音を清音に変換するモードを表します
	RemoveVoicedSoundMode ConversionMode = "unvoiced"
	// AddVoicedSoundMode は清音を濁音に変換するモードを表します
	AddVoicedSoundMode ConversionMode = "voiced"
	// VoicedSoundPairsMode は濁音と半濁音の変換処理だけを行うモードを表します
	VoicedSoundPairsMode ConversionMode = "voiced-pairs"
)

// VoicedSoundMap は清音と濁音のマッピングを保持します
var VoicedSoundMap = map[rune]rune{
	// 平仮名の変換マップ（清音→濁音）
	'か': 'が', 'き': 'ぎ', 'く': 'ぐ', 'け': 'げ', 'こ': 'ご',
	'さ': 'ざ', 'し': 'じ', 'す': 'ず', 'せ': 'ぜ', 'そ': 'ぞ',
	'た': 'だ', 'ち': 'ぢ', 'つ': 'づ', 'て': 'で', 'と': 'ど',
	'は': 'ば', 'ひ': 'び', 'ふ': 'ぶ', 'へ': 'べ', 'ほ': 'ぼ',
	// 片仮名の変換マップ（清音→濁音）
	'カ': 'ガ', 'キ': 'ギ', 'ク': 'グ', 'ケ': 'ゲ', 'コ': 'ゴ',
	'サ': 'ザ', 'シ': 'ジ', 'ス': 'ズ', 'セ': 'ゼ', 'ソ': 'ゾ',
	'タ': 'ダ', 'チ': 'ヂ', 'ツ': 'ヅ', 'テ': 'デ', 'ト': 'ド',
	'ハ': 'バ', 'ヒ': 'ビ', 'フ': 'ブ', 'ヘ': 'ベ', 'ホ': 'ボ',
}

// SemiVoicedSoundMap は清音と半濁音のマッピングを保持します
var SemiVoicedSoundMap = map[rune]rune{
	// 平仮名の変換マップ（清音→半濁音）
	'は': 'ぱ', 'ひ': 'ぴ', 'ふ': 'ぷ', 'へ': 'ぺ', 'ほ': 'ぽ',
	// 片仮名の変換マップ（清音→半濁音）
	'ハ': 'パ', 'ヒ': 'ピ', 'フ': 'プ', 'ヘ': 'ペ', 'ホ': 'ポ',
}

// VoicedSoundPairsMap は濁音と半濁音の変換ペアを保持します
// Pythonのmake_voicedsound関数の動作を再現
var VoicedSoundPairsMap = map[string]string{
	// 平仮名の濁音ペア（2文字 → 1文字）
	"が": "が", "ぎ": "ぎ", "ぐ": "ぐ", "げ": "げ", "ご": "ご",
	"ざ": "ざ", "じ": "じ", "ず": "ず", "ぜ": "ぜ", "ぞ": "ぞ",
	"だ": "だ", "ぢ": "ぢ", "づ": "づ", "で": "で", "ど": "ど",
	"ば": "ば", "び": "び", "ぶ": "ぶ", "べ": "べ", "ぼ": "ぼ",
	"ぱ": "ぱ", "ぴ": "ぴ", "ぷ": "ぷ", "ぺ": "ぺ", "ぽ": "ぽ",
	// 片仮名の濁音ペア（2文字 → 1文字）
	"ガ": "ガ", "ギ": "ギ", "グ": "グ", "ゲ": "ゲ", "ゴ": "ゴ",
	"ザ": "ザ", "ジ": "ジ", "ズ": "ズ", "ゼ": "ゼ", "ゾ": "ゾ",
	"ダ": "ダ", "ヂ": "ヂ", "ヅ": "ヅ", "デ": "デ", "ド": "ド",
	"バ": "バ", "ビ": "ビ", "ブ": "ブ", "ベ": "ベ", "ボ": "ボ",
	"パ": "パ", "ピ": "ピ", "プ": "プ", "ペ": "ペ", "ポ": "ポ",
}

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
	case RemoveVoicedSoundMode:
		return k.RemoveVoicedSound(input)
	case AddVoicedSoundMode:
		return k.AddVoicedSound(input)
	case VoicedSoundPairsMode:
		return k.ConvertVoicedSoundPairs(input)
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

// RemoveVoicedSound は濁音と半濁音を清音に変換します
func (k *KanaConverter) RemoveVoicedSound(input string) string {
	// まず全角に正規化
	input = norm.NFKC.String(input)

	// 結果を格納するためのバッファを作成
	var result strings.Builder

	// 濁音と半濁音を清音に変換
	for _, r := range input {
		replaced := false

		// 濁音から清音への変換
		for voiced, unvoiced := range reverseMap(VoicedSoundMap) {
			if r == voiced {
				result.WriteRune(unvoiced)
				replaced = true
				break
			}
		}

		// 半濁音から清音への変換
		if !replaced {
			for semiVoiced, unvoiced := range reverseMap(SemiVoicedSoundMap) {
				if r == semiVoiced {
					result.WriteRune(unvoiced)
					replaced = true
					break
				}
			}
		}

		// 変換されなかった文字はそのまま追加
		if !replaced {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// AddVoicedSound は清音を濁音に変換します（可能な場合）
func (k *KanaConverter) AddVoicedSound(input string) string {
	// まず全角に正規化
	input = norm.NFKC.String(input)

	// 結果を格納するためのバッファを作成
	var result strings.Builder

	// 清音を濁音に変換
	for _, r := range input {
		replaced := false

		// 清音から濁音への変換
		if voiced, ok := VoicedSoundMap[r]; ok {
			result.WriteRune(voiced)
			replaced = true
		}

		// 変換されなかった文字はそのまま追加
		if !replaced {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// ConvertVoicedSoundPairs は濁音と半濁音の変換処理だけを行います
// Pythonのmake_voicedsound関数の動作を再現
func (k *KanaConverter) ConvertVoicedSoundPairs(input string) string {
	// まず全角に正規化
	input = norm.NFKC.String(input)

	// 濁音と半濁音のペアを処理
	for pair, replacement := range VoicedSoundPairsMap {
		input = strings.ReplaceAll(input, pair, replacement)
	}

	return input
}

// reverseMap はマップのキーと値を入れ替えた新しいマップを作成します
func reverseMap(m map[rune]rune) map[rune]rune {
	reversed := make(map[rune]rune, len(m))
	for k, v := range m {
		reversed[v] = k
	}
	return reversed
}

// IsValidMode はモード文字列が有効かどうかを判定します
func IsValidMode(mode string) bool {
	return mode == string(FullWidthMode) ||
	       mode == string(HalfWidthMode) ||
	       mode == string(RemoveVoicedSoundMode) ||
	       mode == string(AddVoicedSoundMode) ||
	       mode == string(VoicedSoundPairsMode)
}

// GetModeFromString は文字列からConversionModeを取得します
func GetModeFromString(mode string) ConversionMode {
	switch mode {
	case string(HalfWidthMode):
		return HalfWidthMode
	case string(RemoveVoicedSoundMode):
		return RemoveVoicedSoundMode
	case string(AddVoicedSoundMode):
		return AddVoicedSoundMode
	case string(VoicedSoundPairsMode):
		return VoicedSoundPairsMode
	default:
		return FullWidthMode
	}
}

// GetModeDescription はモードの説明文を取得します
func GetModeDescription(mode ConversionMode) string {
	switch mode {
	case FullWidthMode:
		return "full-width"
	case HalfWidthMode:
		return "half-width"
	case RemoveVoicedSoundMode:
		return "remove-voiced-sound"
	case AddVoicedSoundMode:
		return "add-voiced-sound"
	case VoicedSoundPairsMode:
		return "voiced-pairs"
	default:
		return "unknown"
	}
}
