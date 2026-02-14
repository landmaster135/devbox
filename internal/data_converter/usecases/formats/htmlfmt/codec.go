package htmlfmt

import (
	"bytes"
	"fmt"
	stdhtml "html"
	"strings"

	"github.com/landmaster135/devbox/internal/data_converter/usecases/formats/common"
	"golang.org/x/net/html"
)

// Parse は HTML テーブルを key-value リストへ正規化します。
func Parse(content []byte) (*common.NormalizedData, error) {
	doc, err := html.Parse(bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("HTMLの解析に失敗しました: %w", err)
	}

	table := findFirstElementByTag(doc, "table")
	if table == nil {
		return nil, fmt.Errorf("HTML内に table 要素がありません")
	}

	rows := collectTableRows(table)
	if len(rows) == 0 {
		return nil, fmt.Errorf("HTMLテーブルに行がありません")
	}

	firstCells, firstRowHasHeader := extractRowCells(rows[0])
	if len(firstCells) == 0 {
		return nil, fmt.Errorf("HTMLテーブルの先頭行にセルがありません")
	}

	headers := make([]string, 0, len(firstCells))
	dataStartIndex := 0
	if firstRowHasHeader {
		headers = ensureUniqueHeaderNames(firstCells)
		dataStartIndex = 1
	} else {
		for i := range firstCells {
			headers = append(headers, fmt.Sprintf("column_%d", i+1))
		}
	}

	records := make([]map[string]string, 0, len(rows)-dataStartIndex)
	for rowIndex := dataStartIndex; rowIndex < len(rows); rowIndex++ {
		cells, _ := extractRowCells(rows[rowIndex])
		if len(cells) == 0 {
			continue
		}
		if len(cells) > len(headers) {
			return nil, fmt.Errorf("HTMLテーブルの%d行目のセル数がヘッダーより多いです", rowIndex+1)
		}
		records = append(records, common.BuildRecordFromValues(headers, cells))
	}

	return common.NewNormalizedData(records, headers), nil
}

// Serialize は key-value リストを HTML テーブルへ変換します。
func Serialize(records []map[string]string, keys []string) ([]byte, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("HTML出力に必要なキー情報がありません")
	}

	var builder strings.Builder
	builder.WriteString("<table>\n")
	builder.WriteString("  <thead>\n")
	builder.WriteString("    <tr>")
	for _, key := range keys {
		builder.WriteString("<th>")
		builder.WriteString(stdhtml.EscapeString(key))
		builder.WriteString("</th>")
	}
	builder.WriteString("</tr>\n")
	builder.WriteString("  </thead>\n")
	builder.WriteString("  <tbody>\n")

	for _, record := range records {
		builder.WriteString("    <tr>")
		for _, key := range keys {
			builder.WriteString("<td>")
			builder.WriteString(stdhtml.EscapeString(record[key]))
			builder.WriteString("</td>")
		}
		builder.WriteString("</tr>\n")
	}

	builder.WriteString("  </tbody>\n")
	builder.WriteString("</table>\n")

	return []byte(builder.String()), nil
}

func findFirstElementByTag(node *html.Node, tag string) *html.Node {
	if node == nil {
		return nil
	}
	if node.Type == html.ElementNode && node.Data == tag {
		return node
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if found := findFirstElementByTag(child, tag); found != nil {
			return found
		}
	}
	return nil
}

func collectTableRows(tableNode *html.Node) []*html.Node {
	rows := make([]*html.Node, 0)
	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "tr" {
			rows = append(rows, node)
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}
	walk(tableNode)
	return rows
}

func extractRowCells(row *html.Node) ([]string, bool) {
	cells := make([]string, 0)
	hasHeader := false
	for child := row.FirstChild; child != nil; child = child.NextSibling {
		if child.Type != html.ElementNode {
			continue
		}
		if child.Data != "th" && child.Data != "td" {
			continue
		}
		if child.Data == "th" {
			hasHeader = true
		}
		cells = append(cells, strings.TrimSpace(extractTextContent(child)))
	}
	return cells, hasHeader
}

func extractTextContent(node *html.Node) string {
	var builder strings.Builder
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.TextNode {
			builder.WriteString(n.Data)
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}
	walk(node)
	return builder.String()
}

func ensureUniqueHeaderNames(headers []string) []string {
	results := make([]string, 0, len(headers))
	used := map[string]int{}
	for i, raw := range headers {
		header := strings.TrimSpace(raw)
		if header == "" {
			header = fmt.Sprintf("column_%d", i+1)
		}
		if count, exists := used[header]; exists {
			count++
			used[header] = count
			header = fmt.Sprintf("%s_%d", header, count)
		} else {
			used[header] = 1
		}
		results = append(results, header)
	}
	return results
}
