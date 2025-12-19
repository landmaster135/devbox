package converters

import (
	"strings"
)

// CSVConverter はCSV変換を行う構造体
type CSVConverter struct{}

// NewCSVConverter は新しいCSVConverterを作成する
func NewCSVConverter() *CSVConverter {
	return &CSVConverter{}
}

// ConvertToCSV は2次元配列をCSV形式に変換する
func (c *CSVConverter) ConvertToCSV(values [][]string) string {
	return c.convertWithSeparator(values, ",")
}

// ConvertToTSV は2次元配列をTSV形式に変換する
func (c *CSVConverter) ConvertToTSV(values [][]string) string {
	return c.convertWithSeparator(values, "\t")
}

// convertWithSeparator は指定された区切り文字で2次元配列を変換する
func (c *CSVConverter) convertWithSeparator(values [][]string, separator string) string {
	if len(values) == 0 {
		return ""
	}

	var result strings.Builder
	for i, row := range values {
		rowStr := c.convertRowWithSeparator(row, separator)
		result.WriteString(rowStr)
		if i < len(values)-1 {
			result.WriteString("\n")
		}
	}
	return result.String()
}

// convertRowWithSeparator は行を指定された区切り文字で変換する
func (c *CSVConverter) convertRowWithSeparator(row []string, separator string) string {
	if len(row) == 0 {
		return ""
	}

	var result strings.Builder
	for i, cell := range row {
		// CSVの場合、カンマやダブルクォートが含まれる場合はエスケープが必要
		// ここでは簡単な実装として、ダブルクォートで囲む処理を追加
		processedCell := c.escapeCell(cell, separator)
		result.WriteString(processedCell)
		if i < len(row)-1 {
			result.WriteString(separator)
		}
	}
	return result.String()
}

// escapeCell はセルの内容をエスケープする
func (c *CSVConverter) escapeCell(cell, separator string) string {
	// 区切り文字、改行、ダブルクォートが含まれる場合はダブルクォートで囲む
	needsQuoting := strings.Contains(cell, separator) ||
		strings.Contains(cell, "\n") ||
		strings.Contains(cell, "\r") ||
		strings.Contains(cell, "\"")

	if needsQuoting {
		// ダブルクォートをエスケープ（""に変換）
		escaped := strings.ReplaceAll(cell, "\"", "\"\"")
		return "\"" + escaped + "\""
	}

	return cell
}
