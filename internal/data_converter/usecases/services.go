package usecases

import (
	"fmt"
	"os"
	"strings"

	"github.com/landmaster135/devbox/internal/data_converter/usecases/formats/common"
	"github.com/landmaster135/devbox/internal/data_converter/usecases/formats/csvfmt"
	"github.com/landmaster135/devbox/internal/data_converter/usecases/formats/htmlfmt"
	"github.com/landmaster135/devbox/internal/data_converter/usecases/formats/jsonfmt"
	"github.com/landmaster135/devbox/internal/data_converter/usecases/formats/tsvfmt"
	"github.com/landmaster135/devbox/internal/data_converter/usecases/formats/yamlfmt"
)

const (
	formatJSON = "json"
	formatYAML = "yaml"
	formatCSV  = "csv"
	formatTSV  = "tsv"
	formatHTML = "html"
)

// NormalizedData は全入力形式を統一した key-value リスト表現です。
type NormalizedData = common.NormalizedData

// Service はデータ変換処理を提供します。
type Service struct{}

// NewService は Service を生成します。
func NewService() *Service {
	return &Service{}
}

// ConvertFile は入力ファイルを読み込み、key-value リストへ正規化した上で指定形式へ変換して出力します。
func (s *Service) ConvertFile(inputPath, outputPath, inputFormat, outputFormat string) (string, error) {
	inputBytes, err := os.ReadFile(strings.TrimSpace(inputPath))
	if err != nil {
		return "", fmt.Errorf("入力ファイルの読み込みに失敗しました: %w", err)
	}

	normalized, err := s.NormalizeToKeyValueList(inputBytes, inputFormat)
	if err != nil {
		return "", err
	}

	outputBytes, err := s.SerializeFromKeyValueList(normalized, outputFormat)
	if err != nil {
		return "", err
	}

	perm := resolveOutputFileMode(outputPath)
	if err := os.WriteFile(outputPath, outputBytes, perm); err != nil {
		return "", fmt.Errorf("出力ファイルの書き込みに失敗しました: %w", err)
	}

	return fmt.Sprintf("変換完了: %s (%s) -> %s (%s)", inputPath, normalizeFormat(inputFormat), outputPath, normalizeFormat(outputFormat)), nil
}

// NormalizeToKeyValueList は任意形式の入力データを key-value リストへ正規化します。
func (s *Service) NormalizeToKeyValueList(content []byte, format string) (*NormalizedData, error) {
	switch normalizeFormat(format) {
	case formatJSON:
		return jsonfmt.Parse(content)
	case formatYAML:
		return yamlfmt.Parse(content)
	case formatCSV:
		return csvfmt.Parse(content)
	case formatTSV:
		return tsvfmt.Parse(content)
	case formatHTML:
		return htmlfmt.Parse(content)
	default:
		return nil, fmt.Errorf("未対応の入力形式です: %s", format)
	}
}

// SerializeFromKeyValueList は key-value リストを指定形式へ変換します。
func (s *Service) SerializeFromKeyValueList(data *NormalizedData, format string) ([]byte, error) {
	if data == nil {
		return nil, fmt.Errorf("変換元データが nil です")
	}

	keys := data.Keys
	if len(keys) == 0 {
		keys = common.CollectSortedKeys(data.KeyValueList)
	}

	switch normalizeFormat(format) {
	case formatJSON:
		return jsonfmt.Serialize(data.KeyValueList)
	case formatYAML:
		return yamlfmt.Serialize(data.KeyValueList)
	case formatCSV:
		return csvfmt.Serialize(data.KeyValueList, keys)
	case formatTSV:
		return tsvfmt.Serialize(data.KeyValueList, keys)
	case formatHTML:
		return htmlfmt.Serialize(data.KeyValueList, keys)
	default:
		return nil, fmt.Errorf("未対応の出力形式です: %s", format)
	}
}

func normalizeFormat(format string) string {
	return strings.ToLower(strings.TrimSpace(format))
}

func resolveOutputFileMode(outputPath string) os.FileMode {
	perm := os.FileMode(0o644)
	if stat, statErr := os.Stat(outputPath); statErr == nil {
		perm = stat.Mode()
	}
	return perm
}
