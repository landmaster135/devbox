package usecases

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/landmaster135/devbox/internal/data_converter/converters"
	"github.com/landmaster135/devbox/internal/data_converter/parsers"
)

// DataConverterService はデータ変換サービス
type DataConverterService struct {
	parser        *parsers.DataParser
	htmlConverter *converters.HTMLConverter
	csvConverter  *converters.CSVConverter
}

// NewDataConverterService は新しいDataConverterServiceを作成する
func NewDataConverterService() *DataConverterService {
	return &DataConverterService{
		parser:        parsers.NewDataParser(),
		htmlConverter: converters.NewHTMLConverter(),
		csvConverter:  converters.NewCSVConverter(),
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

// convertToJSON はデータをJSON形式（オブジェクトの配列）に変換する
func (s *DataConverterService) convertToJSON(data [][]string) (string, error) {
	if len(data) == 0 {
		return "[]", nil
	}

	// 1行目をヘッダーとして使用
	headers := data[0]
	var result []map[string]any

	// 2行目以降をオブジェクトに変換
	for i := 1; i < len(data); i++ {
		row := data[i]
		obj := make(map[string]any)

		for j, header := range headers {
			var value any
			if j < len(row) {
				cellValue := strings.TrimSpace(row[j])
				if cellValue == "" {
					value = nil // 空の値はnullに
				} else {
					// 数値判定を行い、可能なら数値に変換
					if num, err := strconv.Atoi(cellValue); err == nil {
						value = num
					} else if num, err := strconv.ParseFloat(cellValue, 64); err == nil {
						value = num
					} else {
						value = cellValue
					}
				}
			} else {
				value = nil // フィールドが存在しない場合もnull
			}
			obj[header] = value
		}
		result = append(result, obj)
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("JSON変換エラー: %v", err)
	}
	return string(jsonBytes), nil
}
