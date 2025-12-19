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
	case "html":
		return p.parseHTML(data)
	case "list":
		return p.parseMarkdownList(data)
	case "ordered-list":
		return p.parseMarkdownOrderedList(data)
	case "table":
		return p.parseMarkdownTable(data)
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

// parseMarkdownTable はMarkdownテーブル形式のデータを解析する
func (p *DataParser) parseMarkdownTable(data string) ([][]string, error) {
	if data == "" {
		return [][]string{}, nil
	}

	lines := strings.Split(strings.TrimSpace(data), "\n")
	var result [][]string
	separatorFound := false

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// パイプで始まり、パイプで終わる行をチェック
		if !strings.HasPrefix(line, "|") || !strings.HasSuffix(line, "|") {
			return nil, fmt.Errorf("行 %d: Markdownテーブルの形式が正しくありません: %s", i+1, line)
		}

		// セパレーター行をチェック（|---|---|のような形式）
		if strings.Contains(line, "---") {
			separatorFound = true
			continue // セパレーター行はスキップ
		}

		// セルを抽出
		cells := p.parseMarkdownTableRow(line)
		if len(cells) > 0 {
			result = append(result, cells)
		}
	}

	// セパレーター行が見つからない場合はエラー
	if !separatorFound && len(result) > 1 {
		return nil, fmt.Errorf("Markdownテーブルのセパレーター行が見つかりません")
	}

	return result, nil
}

// parseMarkdownTableRow はMarkdownテーブルの1行を解析する
func (p *DataParser) parseMarkdownTableRow(line string) []string {
	// 先頭と末尾のパイプを除去
	line = strings.TrimPrefix(line, "|")
	line = strings.TrimSuffix(line, "|")

	// パイプで分割
	cells := strings.Split(line, "|")

	// 各セルをトリム
	for i, cell := range cells {
		cells[i] = strings.TrimSpace(cell)
	}

	return cells
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

// parseHTML はHTML形式のデータを解析する
func (p *DataParser) parseHTML(data string) ([][]string, error) {
	// 簡単なHTMLテーブル解析の実装
	// <table>タグ内の<tr>と<td>/<th>を抽出

	// テーブルタグを探す
	tableStart := strings.Index(strings.ToLower(data), "<table")
	if tableStart == -1 {
		return nil, fmt.Errorf("HTMLテーブルが見つかりません")
	}

	tableEnd := strings.Index(strings.ToLower(data[tableStart:]), "</table>")
	if tableEnd == -1 {
		return nil, fmt.Errorf("HTMLテーブルの終了タグが見つかりません")
	}

	tableContent := data[tableStart : tableStart+tableEnd+8] // "</table>"を含む

	// 行を抽出
	rows := p.extractTableRows(tableContent)

	var result [][]string
	for _, row := range rows {
		cells := p.extractTableCells(row)
		if len(cells) > 0 {
			result = append(result, cells)
		}
	}

	return result, nil
}

// extractTableRows はHTMLテーブルから行を抽出する
func (p *DataParser) extractTableRows(tableHTML string) []string {
	var rows []string
	content := strings.ToLower(tableHTML)

	start := 0
	for {
		trStart := strings.Index(content[start:], "<tr")
		if trStart == -1 {
			break
		}
		trStart += start

		// <tr>タグの終了を探す
		tagEnd := strings.Index(content[trStart:], ">")
		if tagEnd == -1 {
			break
		}
		trContentStart := trStart + tagEnd + 1

		// </tr>タグを探す
		trEnd := strings.Index(content[trContentStart:], "</tr>")
		if trEnd == -1 {
			break
		}
		trEnd += trContentStart

		// 元のHTMLから行内容を抽出（大文字小文字を保持）
		rowContent := tableHTML[trContentStart:trEnd]
		rows = append(rows, rowContent)

		start = trEnd + 5 // "</tr>"の長さ
	}

	return rows
}

// extractTableCells は行からセルを抽出する
func (p *DataParser) extractTableCells(rowHTML string) []string {
	var cells []string
	content := strings.ToLower(rowHTML)
	original := rowHTML

	start := 0
	for {
		// <td>または<th>タグを探す
		tdStart := strings.Index(content[start:], "<td")
		thStart := strings.Index(content[start:], "<th")

		var cellStart int
		var isHeader bool

		if tdStart == -1 && thStart == -1 {
			break
		} else if tdStart == -1 {
			cellStart = thStart + start
			isHeader = true
		} else if thStart == -1 {
			cellStart = tdStart + start
			isHeader = false
		} else if tdStart < thStart {
			cellStart = tdStart + start
			isHeader = false
		} else {
			cellStart = thStart + start
			isHeader = true
		}

		// タグの終了を探す
		tagEnd := strings.Index(content[cellStart:], ">")
		if tagEnd == -1 {
			break
		}
		cellContentStart := cellStart + tagEnd + 1

		// 終了タグを探す
		var endTag string
		if isHeader {
			endTag = "</th>"
		} else {
			endTag = "</td>"
		}

		cellEnd := strings.Index(content[cellContentStart:], endTag)
		if cellEnd == -1 {
			break
		}
		cellEnd += cellContentStart

		// セル内容を抽出（HTMLタグを除去）
		cellContent := original[cellContentStart:cellEnd]
		cleanContent := p.stripHTMLTags(cellContent)
		cells = append(cells, strings.TrimSpace(cleanContent))

		start = cellEnd + len(endTag)
	}

	return cells
}

// stripHTMLTags はHTMLタグを除去する
func (p *DataParser) stripHTMLTags(html string) string {
	result := html

	// 簡単なHTMLタグ除去
	for {
		start := strings.Index(result, "<")
		if start == -1 {
			break
		}
		end := strings.Index(result[start:], ">")
		if end == -1 {
			break
		}
		end += start + 1
		result = result[:start] + result[end:]
	}

	// HTMLエンティティのデコード（基本的なもののみ）
	result = strings.ReplaceAll(result, "&lt;", "<")
	result = strings.ReplaceAll(result, "&gt;", ">")
	result = strings.ReplaceAll(result, "&amp;", "&")
	result = strings.ReplaceAll(result, "&quot;", "\"")
	result = strings.ReplaceAll(result, "&#39;", "'")
	result = strings.ReplaceAll(result, "&nbsp;", " ")

	return result
}

// parseMarkdownList は箇条書きリスト形式のデータを解析する
func (p *DataParser) parseMarkdownList(data string) ([][]string, error) {
	if data == "" {
		return [][]string{}, nil
	}

	lines := strings.Split(strings.TrimSpace(data), "\n")
	var result [][]string

	// ヘッダー行を追加
	result = append(result, []string{"項目"})

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 箇条書きリストの形式をチェック（- または * で始まる）
		if strings.HasPrefix(line, "- ") {
			item := strings.TrimSpace(line[2:])
			if item != "" {
				result = append(result, []string{item})
			}
		} else if strings.HasPrefix(line, "* ") {
			item := strings.TrimSpace(line[2:])
			if item != "" {
				result = append(result, []string{item})
			}
		} else {
			return nil, fmt.Errorf("行 %d: 箇条書きリストの形式が正しくありません: %s", i+1, line)
		}
	}

	return result, nil
}

// parseMarkdownOrderedList は順序付きリスト形式のデータを解析する
func (p *DataParser) parseMarkdownOrderedList(data string) ([][]string, error) {
	if data == "" {
		return [][]string{}, nil
	}

	lines := strings.Split(strings.TrimSpace(data), "\n")
	var result [][]string

	// ヘッダー行を追加
	result = append(result, []string{"番号", "項目"})

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 順序付きリストの形式をチェック（数字. で始まる）
		dotIndex := strings.Index(line, ". ")
		if dotIndex == -1 {
			return nil, fmt.Errorf("行 %d: 順序付きリストの形式が正しくありません: %s", i+1, line)
		}

		numberStr := line[:dotIndex]
		item := strings.TrimSpace(line[dotIndex+2:])

		// 数字の妥当性をチェック
		if _, err := fmt.Sscanf(numberStr, "%d", new(int)); err != nil {
			return nil, fmt.Errorf("行 %d: 番号が正しくありません: %s", i+1, numberStr)
		}

		if item != "" {
			result = append(result, []string{numberStr, item})
		}
	}

	return result, nil
}
