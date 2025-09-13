package converters

import (
	"fmt"
	"strconv"
	"strings"
)

// MarkdownConverter はMarkdownリスト変換を行う構造体
type MarkdownConverter struct{}

// NewMarkdownConverter は新しいMarkdownConverterを作成する
func NewMarkdownConverter() *MarkdownConverter {
	return &MarkdownConverter{}
}

// ConvertToList はデータを箇条書きリストに変換する
func (c *MarkdownConverter) ConvertToList(data [][]string) string {
	if len(data) == 0 {
		return ""
	}

	var result strings.Builder

	// 1列のテーブルの場合は単純な箇条書きリスト
	if len(data) > 0 && len(data[0]) == 1 {
		// ヘッダー行をスキップ（最初の行がヘッダーと仮定）
		startIndex := 1
		if len(data) == 1 {
			startIndex = 0 // データが1行しかない場合はヘッダーなしと判断
		}

		for i := startIndex; i < len(data); i++ {
			if len(data[i]) > 0 && strings.TrimSpace(data[i][0]) != "" {
				result.WriteString("- ")
				result.WriteString(strings.TrimSpace(data[i][0]))
				result.WriteString("\n")
			}
		}
	} else {
		// 複数列のテーブルの場合は各行を「項目名: 値」形式で変換
		if len(data) > 1 {
			headers := data[0]
			for i := 1; i < len(data); i++ {
				row := data[i]
				result.WriteString("- ")

				var parts []string
				for j, header := range headers {
					var value string
					if j < len(row) {
						value = strings.TrimSpace(row[j])
					}
					if value != "" {
						parts = append(parts, fmt.Sprintf("%s: %s", strings.TrimSpace(header), value))
					}
				}

				if len(parts) > 0 {
					result.WriteString(strings.Join(parts, ", "))
				}
				result.WriteString("\n")
			}
		}
	}

	return result.String()
}

// addEmptyColumnIfSingleColumn は1列のデータの場合、空の列を先頭に追加する
func (c *MarkdownConverter) addEmptyColumnIfSingleColumn(data [][]string) [][]string {
	if len(data) == 0 {
		return data
	}

	// 1列のデータかどうかをチェック
	isSingleColumn := true
	for _, row := range data {
		if len(row) != 1 {
			isSingleColumn = false
			break
		}
	}

	// 1列でない場合はそのまま返す
	if !isSingleColumn {
		return data
	}

	// 1列の場合、空の列を先頭に追加して2列にする
	result := make([][]string, len(data))
	for i, row := range data {
		if i == 0 {
			// ヘッダー行の場合
			result[i] = []string{"", row[0]}
		} else {
			// データ行の場合
			result[i] = []string{"", row[0]}
		}
	}

	return result
}

// ConvertToTable はデータをMarkdownテーブルに変換する
func (c *MarkdownConverter) ConvertToTable(data [][]string) string {
	if len(data) == 0 {
		return ""
	}

	// 1列のデータ（箇条書きリストから変換されたデータ）の場合、空の列を先頭に追加
	processedData := c.addEmptyColumnIfSingleColumn(data)

	var result strings.Builder

	// 各列の最大幅を計算
	colWidths := c.calculateColumnWidths(processedData)

	// ヘッダー行を出力
	if len(processedData) > 0 {
		result.WriteString("|")
		for i, cell := range processedData[0] {
			paddedCell := c.padCell(cell, colWidths[i])
			result.WriteString(" ")
			result.WriteString(paddedCell)
			result.WriteString(" |")
		}
		result.WriteString("\n")

		// セパレーター行を出力
		result.WriteString("|")
		for i := 0; i < len(processedData[0]); i++ {
			result.WriteString(strings.Repeat("-", colWidths[i]+2))
			result.WriteString("|")
		}
		result.WriteString("\n")
	}

	// データ行を出力
	for i := 1; i < len(processedData); i++ {
		row := processedData[i]
		result.WriteString("|")
		for j := 0; j < len(processedData[0]); j++ {
			var cell string
			if j < len(row) {
				cell = strings.TrimSpace(row[j])
			}
			paddedCell := c.padCell(cell, colWidths[j])
			result.WriteString(" ")
			result.WriteString(paddedCell)
			result.WriteString(" |")
		}
		result.WriteString("\n")
	}

	return result.String()
}

// calculateColumnWidths は各列の最大幅を計算する
func (c *MarkdownConverter) calculateColumnWidths(data [][]string) []int {
	if len(data) == 0 {
		return []int{}
	}

	maxCols := 0
	for _, row := range data {
		if len(row) > maxCols {
			maxCols = len(row)
		}
	}

	widths := make([]int, maxCols)

	for _, row := range data {
		for i, cell := range row {
			cellLen := len(strings.TrimSpace(cell))
			if cellLen > widths[i] {
				widths[i] = cellLen
			}
		}
	}

	// 最小幅を3に設定（セパレーター行の"---"のため）
	for i := range widths {
		if widths[i] < 3 {
			widths[i] = 3
		}
	}

	return widths
}

// padCell はセルを指定された幅にパディングする
func (c *MarkdownConverter) padCell(cell string, width int) string {
	cell = strings.TrimSpace(cell)
	if len(cell) >= width {
		return cell
	}
	return cell + strings.Repeat(" ", width-len(cell))
}

// ConvertToOrderedList はデータを順序付きリストに変換する
func (c *MarkdownConverter) ConvertToOrderedList(data [][]string) string {
	if len(data) == 0 {
		return ""
	}

	var result strings.Builder

	// 1列のテーブルの場合は単純な順序付きリスト
	if len(data) > 0 && len(data[0]) == 1 {
		// ヘッダー行をスキップ（最初の行がヘッダーと仮定）
		startIndex := 1
		if len(data) == 1 {
			startIndex = 0 // データが1行しかない場合はヘッダーなしと判断
		}

		counter := 1
		for i := startIndex; i < len(data); i++ {
			if len(data[i]) > 0 && strings.TrimSpace(data[i][0]) != "" {
				result.WriteString(strconv.Itoa(counter))
				result.WriteString(". ")
				result.WriteString(strings.TrimSpace(data[i][0]))
				result.WriteString("\n")
				counter++
			}
		}
	} else {
		// 複数列のテーブルの場合は各行を「項目名: 値」形式で変換
		if len(data) > 1 {
			headers := data[0]
			counter := 1
			for i := 1; i < len(data); i++ {
				row := data[i]
				result.WriteString(strconv.Itoa(counter))
				result.WriteString(". ")

				var parts []string
				for j, header := range headers {
					var value string
					if j < len(row) {
						value = strings.TrimSpace(row[j])
					}
					if value != "" {
						parts = append(parts, fmt.Sprintf("%s: %s", strings.TrimSpace(header), value))
					}
				}

				if len(parts) > 0 {
					result.WriteString(strings.Join(parts, ", "))
				}
				result.WriteString("\n")
				counter++
			}
		}
	}

	return result.String()
}
