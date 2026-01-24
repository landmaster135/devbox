package usecases

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/landmaster135/devbox/internal/data_converter/converters"
	"github.com/landmaster135/devbox/internal/data_converter/parsers"
	yaml "gopkg.in/yaml.v3"
)

// DataConverterService はデータ変換サービス
type DataConverterService struct {
	parser            *parsers.DataParser
	htmlConverter     *converters.HTMLConverter
	csvConverter      *converters.CSVConverter
	markdownConverter *converters.MarkdownConverter
}

// NewDataConverterService は新しいDataConverterServiceを作成する
func NewDataConverterService() *DataConverterService {
	return &DataConverterService{
		parser:            parsers.NewDataParser(),
		htmlConverter:     converters.NewHTMLConverter(),
		csvConverter:      converters.NewCSVConverter(),
		markdownConverter: converters.NewMarkdownConverter(),
	}
}

// ConvertData はデータを指定された形式に変換する
func (s *DataConverterService) ConvertData(inputFormat, outputFormat, input, inputFilePath string) (string, error) {
	// 入力データを解析
	data, err := s.parser.ParseInput(inputFormat, input, inputFilePath)
	if err != nil {
		return "", fmt.Errorf("入力データの解析に失敗しました: %v", err)
	}

	// 出力形式に応じて変換
	switch outputFormat {
	case "html":
		return s.convertToHTML(data)
	case "csv":
		return s.convertToCSV(data)
	case "tsv":
		return s.convertToTSV(data)
	case "array":
		return s.convertToArray(data)
	case "json":
		return s.convertToJSON(data)
	case "yaml":
		return s.convertToYAML(data)
	case "list":
		return s.convertToList(data)
	case "ordered-list":
		return s.convertToOrderedList(data)
	case "table":
		return s.convertToTable(data)
	default:
		return "", fmt.Errorf("未対応の出力形式です: %s", outputFormat)
	}
}

// convertToHTML はデータをHTML形式に変換する
func (s *DataConverterService) convertToHTML(data [][]string) (string, error) {
	// ヘッダー行があるかどうかを判定（簡単な実装として、常にtrueとする）
	isTheadContained := len(data) > 0
	textReplacingIfBlank := "💩"

	result := s.htmlConverter.ConvertToHTML(data, isTheadContained, textReplacingIfBlank)
	return result, nil
}

// convertToCSV はデータをCSV形式に変換する
func (s *DataConverterService) convertToCSV(data [][]string) (string, error) {
	result := s.csvConverter.ConvertToCSV(data)
	return result, nil
}

// convertToTSV はデータをTSV形式に変換する
func (s *DataConverterService) convertToTSV(data [][]string) (string, error) {
	result := s.csvConverter.ConvertToTSV(data)
	return result, nil
}

// convertToArray はデータを配列形式（2次元配列）に変換する
func (s *DataConverterService) convertToArray(data [][]string) (string, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("array変換エラー: %v", err)
	}
	return string(jsonBytes), nil
}

// convertToYAML はデータをYAML形式に変換する
func (s *DataConverterService) convertToYAML(data [][]string) (string, error) {
	records, err := s.convertToRecordObjects(data)
	if err != nil {
		return "", err
	}
	yamlBytes, err := yaml.Marshal(records)
	if err != nil {
		return "", fmt.Errorf("YAML変換エラー: %v", err)
	}
	return string(yamlBytes), nil
}

// convertToTable はデータをMarkdownテーブル形式に変換する
func (s *DataConverterService) convertToTable(data [][]string) (string, error) {
	result := s.markdownConverter.ConvertToTable(data)
	return result, nil
}

// convertToList はデータを箇条書きリストに変換する
func (s *DataConverterService) convertToList(data [][]string) (string, error) {
	result := s.markdownConverter.ConvertToList(data)
	return result, nil
}

// convertToOrderedList はデータを順序付きリストに変換する
func (s *DataConverterService) convertToOrderedList(data [][]string) (string, error) {
	result := s.markdownConverter.ConvertToOrderedList(data)
	return result, nil
}

// convertToJSON はデータをJSON形式（オブジェクトの配列）に変換する
func (s *DataConverterService) convertToJSON(data [][]string) (string, error) {
	records, err := s.convertToRecordObjects(data)
	if err != nil {
		return "", err
	}
	jsonBytes, err := json.Marshal(records)
	if err != nil {
		return "", fmt.Errorf("JSON変換エラー: %v", err)
	}
	return string(jsonBytes), nil
}

func (s *DataConverterService) convertToRecordObjects(data [][]string) ([]map[string]any, error) {
	if len(data) == 0 {
		return []map[string]any{}, nil
	}

	headers := data[0]
	records := make([]map[string]any, 0, max(len(data)-1, 0))

	for i := 1; i < len(data); i++ {
		row := data[i]
		obj := make(map[string]any)

		for j, header := range headers {
			var value any
			if j < len(row) {
				cellValue := strings.TrimSpace(row[j])
				if cellValue == "" {
					value = nil
				} else if num, err := strconv.Atoi(cellValue); err == nil {
					value = num
				} else if num, err := strconv.ParseFloat(cellValue, 64); err == nil {
					value = num
				} else {
					value = cellValue
				}
			} else {
				value = nil
			}
			obj[header] = value
		}
		records = append(records, obj)
	}

	return records, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
