package katextablefmt

import (
	"errors"
	"fmt"
	"strings"

	"github.com/landmaster135/devbox/internal/data_converter/usecases/formats/common"
)

const (
	arrayBeginPrefix = `\begin{array}{`
	arrayEndToken    = `\end{array}`
)

// Parse は KaTeX table を key-value リストへ正規化します。
func Parse(content []byte) (*common.NormalizedData, error) {
	body, err := extractArrayBody(normalizeLineBreaks(string(content)))
	if err != nil {
		return nil, err
	}

	rows, err := splitKaTeXRows(body)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, errors.New("KaTeX tableに行がありません")
	}

	headers, err := parseHeaderRow(rows[0])
	if err != nil {
		return nil, err
	}

	records := make([]map[string]string, 0, len(rows)-1)
	for rowIndex := 1; rowIndex < len(rows); rowIndex++ {
		values, parseErr := parseRowCells(rows[rowIndex])
		if parseErr != nil {
			return nil, fmt.Errorf("KaTeX tableの%d行目が不正です: %w", rowIndex+1, parseErr)
		}
		if len(values) > len(headers) {
			return nil, fmt.Errorf("KaTeX tableの%d行目の列数がヘッダーより多いです", rowIndex+1)
		}
		records = append(records, common.BuildRecordFromValues(headers, values))
	}

	return common.NewNormalizedData(records, headers), nil
}

// Serialize は key-value リストを KaTeX table へ変換します。
func Serialize(records []map[string]string, keys []string) ([]byte, error) {
	if len(keys) == 0 {
		return nil, errors.New("katex-table出力に必要なキー情報がありません")
	}

	var builder strings.Builder
	builder.WriteString("$$ \\def\\arraystretch{1.4}\n")
	builder.WriteString("\\small\n")
	builder.WriteString("\\begin{array}{")
	builder.WriteString(buildAlignment(keys))
	builder.WriteString("}\n")
	builder.WriteString("\\hline\n")
	builder.WriteString(renderRow(keys))
	builder.WriteString("\\hline\n")
	for _, record := range records {
		values := common.PickValuesByKeys(record, keys)
		builder.WriteString(renderRow(values))
		builder.WriteString("\\hline\n")
	}
	builder.WriteString("\\end{array}\n")
	builder.WriteString("$$\n")

	return []byte(builder.String()), nil
}

func parseHeaderRow(row string) ([]string, error) {
	headers, err := parseRowCells(row)
	if err != nil {
		return nil, fmt.Errorf("KaTeX tableのヘッダー行が不正です: %w", err)
	}
	if len(headers) == 0 {
		return nil, errors.New("KaTeX tableのヘッダーが空です")
	}

	seen := map[string]struct{}{}
	for i, header := range headers {
		if strings.TrimSpace(header) == "" {
			return nil, fmt.Errorf("KaTeX tableのヘッダー%d列目が空です", i+1)
		}
		if _, exists := seen[header]; exists {
			return nil, fmt.Errorf("KaTeX tableのヘッダーに重複があります: %s", header)
		}
		seen[header] = struct{}{}
	}
	return headers, nil
}

func parseRowCells(row string) ([]string, error) {
	cells, err := splitCells(row)
	if err != nil {
		return nil, err
	}

	values := make([]string, 0, len(cells))
	for _, cell := range cells {
		value, parseErr := parseCellValue(cell)
		if parseErr != nil {
			return nil, parseErr
		}
		values = append(values, value)
	}
	return values, nil
}

func parseCellValue(cell string) (string, error) {
	trimmed := strings.TrimSpace(cell)
	if trimmed == "" {
		return "", nil
	}

	if strings.HasPrefix(trimmed, `\text{`) {
		value, err := parseTextMacro(trimmed)
		if err != nil {
			return "", err
		}
		return unescapeTextContent(value), nil
	}

	if strings.HasPrefix(trimmed, `\text`) {
		return "", errors.New(`\text{...} の形式が不正です`)
	}

	return unescapeTextContent(trimmed), nil
}

func parseTextMacro(cell string) (string, error) {
	const prefix = `\text{`
	if !strings.HasPrefix(cell, prefix) {
		return "", errors.New(`\text{...} で開始していません`)
	}

	var builder strings.Builder
	depth := 1
	for i := len(prefix); i < len(cell); i++ {
		ch := cell[i]
		if ch == '\\' {
			if i+1 < len(cell) {
				builder.WriteByte(ch)
				builder.WriteByte(cell[i+1])
				i++
				continue
			}
			builder.WriteByte(ch)
			continue
		}
		if ch == '{' {
			depth++
			builder.WriteByte(ch)
			continue
		}
		if ch == '}' {
			depth--
			if depth == 0 {
				if strings.TrimSpace(cell[i+1:]) != "" {
					return "", errors.New(`\text{...} の後ろに余分な文字があります`)
				}
				return builder.String(), nil
			}
			builder.WriteByte(ch)
			continue
		}
		builder.WriteByte(ch)
	}

	return "", errors.New(`\text{...} が閉じられていません`)
}

func extractArrayBody(content string) (string, error) {
	start := strings.Index(content, arrayBeginPrefix)
	if start < 0 {
		return "", errors.New(`\begin{array}{...} が見つかりません`)
	}

	alignStart := start + len(arrayBeginPrefix)
	alignEnd := strings.Index(content[alignStart:], "}")
	if alignEnd < 0 {
		return "", errors.New("KaTeX tableの列定義が閉じられていません")
	}

	bodyStart := alignStart + alignEnd + 1
	endRel := strings.Index(content[bodyStart:], arrayEndToken)
	if endRel < 0 {
		return "", errors.New(`\end{array} が見つかりません`)
	}

	return content[bodyStart : bodyStart+endRel], nil
}

func splitKaTeXRows(body string) ([]string, error) {
	rows := make([]string, 0)
	var current strings.Builder

	for i := 0; i < len(body); i++ {
		ch := body[i]
		if ch == '\\' && i+1 < len(body) && body[i+1] == '\\' {
			row := normalizeRow(current.String())
			if row != "" {
				rows = append(rows, row)
			}
			current.Reset()
			i++
			continue
		}
		current.WriteByte(ch)
	}

	rest := normalizeRow(current.String())
	if rest != "" {
		return nil, errors.New(`行終端 \\ が不足しています`)
	}

	return rows, nil
}

func normalizeRow(raw string) string {
	replaced := strings.ReplaceAll(raw, `\hline`, "")
	return strings.TrimSpace(replaced)
}

func splitCells(row string) ([]string, error) {
	cells := make([]string, 0)
	var current strings.Builder

	for i := 0; i < len(row); i++ {
		ch := row[i]
		if ch == '\\' && i+1 < len(row) && row[i+1] == '&' {
			current.WriteByte('&')
			i++
			continue
		}
		if ch == '&' {
			cells = append(cells, strings.TrimSpace(current.String()))
			current.Reset()
			continue
		}
		current.WriteByte(ch)
	}
	cells = append(cells, strings.TrimSpace(current.String()))

	if len(cells) == 0 {
		return nil, errors.New("セルがありません")
	}
	return cells, nil
}

func buildAlignment(keys []string) string {
	var builder strings.Builder
	for range keys {
		builder.WriteString("|l")
	}
	builder.WriteString("|")
	return builder.String()
}

func renderRow(values []string) string {
	escaped := make([]string, 0, len(values))
	for _, value := range values {
		escaped = append(escaped, `\text{`+escapeTextContent(value)+`}`)
	}
	return strings.Join(escaped, " & ") + " \\\\\n"
}

func escapeTextContent(value string) string {
	var builder strings.Builder
	for _, r := range value {
		switch r {
		case '\\':
			builder.WriteString(`\textbackslash{}`)
		case '{':
			builder.WriteString(`\{`)
		case '}':
			builder.WriteString(`\}`)
		case '&':
			builder.WriteString(`\&`)
		default:
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func unescapeTextContent(value string) string {
	replaced := strings.ReplaceAll(value, `\textbackslash{}`, `\`)
	var builder strings.Builder
	for i := 0; i < len(replaced); i++ {
		ch := replaced[i]
		if ch == '\\' && i+1 < len(replaced) {
			next := replaced[i+1]
			switch next {
			case '\\', '{', '}', '&':
				builder.WriteByte(next)
				i++
				continue
			}
		}
		builder.WriteByte(ch)
	}
	return builder.String()
}

func normalizeLineBreaks(input string) string {
	return strings.ReplaceAll(input, "\r\n", "\n")
}
