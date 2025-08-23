package parsers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// DataParser はデータ解析を行う構造体
type DataParser struct{}

// NewDataParser は新しいDataParserを作成する
func NewDataParser() *DataParser {
	return &DataParser{}
}

// ParseInput は入力データを2次元配列に解析する
func (p *DataParser) ParseInput(inputFormat, input, inputFilePath string) ([][]string, error) {
	var data string
	var err error

	// 入力データの取得
	if inputFilePath != "" {
		data, err = p.readFile(inputFilePath)
		if err != nil {
			return nil, fmt.Errorf("ファイル読み込みエラー: %v", err)
		}
	} else {
		data = input
	}

	// 形式に応じて解析
	switch inputFormat {
	case "json":
		return p.parseJSON(data)
	case "csv":
		return p.parseCSV(data)
	case "tsv":
		return p.parseTSV(data)
	default:
		return nil, fmt.Errorf("未対応の入力形式です: %s", inputFormat)
	}
}

// readFile はファイルを読み込む
func (p *DataParser) readFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// parseJSON はJSON形式のデータを解析する
func (p *DataParser) parseJSON(data string) ([][]string, error) {
	var result [][]string
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil, fmt.Errorf("JSON解析エラー: %v", err)
	}
	return result, nil
}

// parseCSV はCSV形式のデータを解析する
func (p *DataParser) parseCSV(data string) ([][]string, error) {
	return p.parseDelimited(data, ",")
}

// parseTSV はTSV形式のデータを解析する
func (p *DataParser) parseTSV(data string) ([][]string, error) {
	return p.parseDelimited(data, "\t")
}

// parseDelimited は区切り文字で区切られたデータを解析する
func (p *DataParser) parseDelimited(data, delimiter string) ([][]string, error) {
	if data == "" {
		return [][]string{}, nil
	}

	lines := strings.Split(strings.TrimSpace(data), "\n")
	result := make([][]string, 0, len(lines))

	for i, line := range lines {
		if line == "" {
			continue
		}

		fields, err := p.parseDelimitedLine(line, delimiter)
		if err != nil {
			return nil, fmt.Errorf("行 %d の解析エラー: %v", i+1, err)
		}
		result = append(result, fields)
	}

	return result, nil
}

// parseDelimitedLine は1行の区切り文字データを解析する
func (p *DataParser) parseDelimitedLine(line, delimiter string) ([]string, error) {
	var fields []string
	var current strings.Builder
	inQuotes := false
	i := 0

	for i < len(line) {
		char := line[i]

		if char == '"' {
			if inQuotes && i+1 < len(line) && line[i+1] == '"' {
				// エスケープされたダブルクォート
				current.WriteByte('"')
				i += 2
			} else {
				// クォートの開始または終了
				inQuotes = !inQuotes
				i++
			}
		} else if !inQuotes && strings.HasPrefix(line[i:], delimiter) {
			// 区切り文字
			fields = append(fields, current.String())
			current.Reset()
			i += len(delimiter)
		} else {
			current.WriteByte(char)
			i++
		}
	}

	// 最後のフィールドを追加
	fields = append(fields, current.String())

	return fields, nil
}
