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
