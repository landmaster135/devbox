package mdtablefmt

import (
	"errors"
	"fmt"
	"strings"

	"github.com/landmaster135/devbox/internal/data_converter/usecases/formats/common"
)

// Parse は Markdown table を key-value リストへ正規化します。
func Parse(content []byte) (*common.NormalizedData, error) {
	lines := collectNonEmptyLines(normalizeLineBreaks(string(content)))
	if len(lines) < 2 {
		return nil, errors.New("Markdown tableの行が不足しています")
	}

	headers, err := parseHeaderRow(lines[0])
	if err != nil {
		return nil, err
	}

	separatorCells, err := splitTableRow(lines[1])
	if err != nil {
		return nil, fmt.Errorf("Markdown tableの区切り行が不正です: %w", err)
	}
	if len(separatorCells) != len(headers) {
		return nil, errors.New("Markdown tableの区切り行の列数がヘッダーと一致しません")
	}
	for i, cell := range separatorCells {
		if !isValidSeparatorCell(cell) {
			return nil, fmt.Errorf("Markdown tableの区切り行の%d列目が不正です", i+1)
		}
	}

	records := make([]map[string]string, 0, len(lines)-2)
	for rowIndex := 2; rowIndex < len(lines); rowIndex++ {
		cells, splitErr := splitTableRow(lines[rowIndex])
		if splitErr != nil {
			return nil, fmt.Errorf("Markdown tableの%d行目が不正です: %w", rowIndex+1, splitErr)
		}
		if len(cells) > len(headers) {
			return nil, fmt.Errorf("Markdown tableの%d行目の列数がヘッダーより多いです", rowIndex+1)
		}
		records = append(records, common.BuildRecordFromValues(headers, cells))
	}

	return common.NewNormalizedData(records, headers), nil
}

// Serialize は key-value リストを Markdown table へ変換します。
func Serialize(records []map[string]string, keys []string) ([]byte, error) {
	if len(keys) == 0 {
		return nil, errors.New("md-table出力に必要なキー情報がありません")
	}

	var builder strings.Builder
	builder.WriteString(renderRow(keys))
	builder.WriteString(renderSeparatorRow(len(keys)))
	for _, record := range records {
		values := common.PickValuesByKeys(record, keys)
		builder.WriteString(renderRow(values))
	}

	return []byte(builder.String()), nil
}

func parseHeaderRow(line string) ([]string, error) {
	headers, err := splitTableRow(line)
	if err != nil {
		return nil, fmt.Errorf("Markdown tableのヘッダー行が不正です: %w", err)
	}
	if len(headers) == 0 {
		return nil, errors.New("Markdown tableのヘッダーが空です")
	}

	seen := map[string]struct{}{}
	for i, header := range headers {
		if header == "" {
			return nil, fmt.Errorf("Markdown tableのヘッダー%d列目が空です", i+1)
		}
		if _, exists := seen[header]; exists {
			return nil, fmt.Errorf("Markdown tableのヘッダーに重複があります: %s", header)
		}
		seen[header] = struct{}{}
	}
	return headers, nil
}

func splitTableRow(line string) ([]string, error) {
	trimmed := strings.TrimSpace(line)
	if !strings.Contains(trimmed, "|") {
		return nil, errors.New("区切り文字 | がありません")
	}

	trimmed = strings.TrimSpace(strings.TrimPrefix(trimmed, "|"))
	trimmed = strings.TrimSpace(strings.TrimSuffix(trimmed, "|"))

	cells := make([]string, 0)
	var cellBuilder strings.Builder
	for i := 0; i < len(trimmed); i++ {
		ch := trimmed[i]
		if ch == '\\' && i+1 < len(trimmed) && trimmed[i+1] == '|' {
			cellBuilder.WriteByte('|')
			i++
			continue
		}
		if ch == '|' {
			cells = append(cells, strings.TrimSpace(cellBuilder.String()))
			cellBuilder.Reset()
			continue
		}
		cellBuilder.WriteByte(ch)
	}
	cells = append(cells, strings.TrimSpace(cellBuilder.String()))

	return cells, nil
}

func isValidSeparatorCell(cell string) bool {
	trimmed := strings.TrimSpace(cell)
	if trimmed == "" {
		return false
	}
	trimmed = strings.TrimPrefix(trimmed, ":")
	trimmed = strings.TrimSuffix(trimmed, ":")
	if len(trimmed) < 3 {
		return false
	}
	for _, ch := range trimmed {
		if ch != '-' {
			return false
		}
	}
	return true
}

func renderRow(values []string) string {
	var builder strings.Builder
	builder.WriteString("| ")
	for i, value := range values {
		if i > 0 {
			builder.WriteString(" | ")
		}
		builder.WriteString(escapeTableCell(value))
	}
	builder.WriteString(" |\n")
	return builder.String()
}

func renderSeparatorRow(columnCount int) string {
	cells := make([]string, 0, columnCount)
	for i := 0; i < columnCount; i++ {
		cells = append(cells, "---")
	}
	return renderRow(cells)
}

func escapeTableCell(value string) string {
	replaced := strings.ReplaceAll(value, "\n", " ")
	return strings.ReplaceAll(replaced, "|", "\\|")
}

func collectNonEmptyLines(content string) []string {
	lines := strings.Split(content, "\n")
	nonEmpty := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		nonEmpty = append(nonEmpty, trimmed)
	}
	return nonEmpty
}

func normalizeLineBreaks(input string) string {
	return strings.ReplaceAll(input, "\r\n", "\n")
}
