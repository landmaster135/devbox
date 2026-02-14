package usecases

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	stdhtml "html"
	"os"
	"sort"
	"strings"

	"golang.org/x/net/html"
	"gopkg.in/yaml.v3"
)

const (
	formatJSON = "json"
	formatYAML = "yaml"
	formatCSV  = "csv"
	formatHTML = "html"
)

// NormalizedData は全入力形式を統一した key-value リスト表現です。
type NormalizedData struct {
	KeyValueList []map[string]string
	Keys         []string
}

// Service はデータ変換処理を提供します。
type Service struct{}

// NewService は Service を生成します。
func NewService() *Service {
	return &Service{}
}

// ConvertFile は入力ファイルを読み込み、key-value リストへ正規化した上で指定形式へ変換して出力します。
func (s *Service) ConvertFile(inputPath, outputPath, inputFormat, outputFormat string) (string, error) {
	inputBytes, err := os.ReadFile(strings.TrimSpace(inputPath))
	if err != nil {
		return "", fmt.Errorf("入力ファイルの読み込みに失敗しました: %w", err)
	}

	normalized, err := s.NormalizeToKeyValueList(inputBytes, inputFormat)
	if err != nil {
		return "", err
	}

	outputBytes, err := s.SerializeFromKeyValueList(normalized, outputFormat)
	if err != nil {
		return "", err
	}

	perm := os.FileMode(0o644)
	if stat, statErr := os.Stat(outputPath); statErr == nil {
		perm = stat.Mode()
	}
	if err := os.WriteFile(outputPath, outputBytes, perm); err != nil {
		return "", fmt.Errorf("出力ファイルの書き込みに失敗しました: %w", err)
	}

	return fmt.Sprintf("変換完了: %s (%s) -> %s (%s)", inputPath, normalizeFormat(inputFormat), outputPath, normalizeFormat(outputFormat)), nil
}

// NormalizeToKeyValueList は任意形式の入力データを key-value リストへ正規化します。
func (s *Service) NormalizeToKeyValueList(content []byte, format string) (*NormalizedData, error) {
	switch normalizeFormat(format) {
	case formatJSON:
		return parseJSONToKeyValueList(content)
	case formatYAML:
		return parseYAMLToKeyValueList(content)
	case formatCSV:
		return parseCSVToKeyValueList(content)
	case formatHTML:
		return parseHTMLToKeyValueList(content)
	default:
		return nil, fmt.Errorf("未対応の入力形式です: %s", format)
	}
}

// SerializeFromKeyValueList は key-value リストを指定形式へ変換します。
func (s *Service) SerializeFromKeyValueList(data *NormalizedData, format string) ([]byte, error) {
	if data == nil {
		return nil, fmt.Errorf("変換元データが nil です")
	}

	keys := data.Keys
	if len(keys) == 0 {
		keys = collectSortedKeys(data.KeyValueList)
	}

	switch normalizeFormat(format) {
	case formatJSON:
		return serializeToJSON(data.KeyValueList)
	case formatYAML:
		return serializeToYAML(data.KeyValueList)
	case formatCSV:
		return serializeToCSV(data.KeyValueList, keys)
	case formatHTML:
		return serializeToHTML(data.KeyValueList, keys)
	default:
		return nil, fmt.Errorf("未対応の出力形式です: %s", format)
	}
}

func parseJSONToKeyValueList(content []byte) (*NormalizedData, error) {
	var raw any
	if err := json.Unmarshal(content, &raw); err != nil {
		return nil, fmt.Errorf("JSONの解析に失敗しました: %w", err)
	}

	recordsAny, err := parseRecordObjects(raw)
	if err != nil {
		return nil, fmt.Errorf("JSONの構造が不正です: %w", err)
	}

	records := normalizeRecords(recordsAny)
	return &NormalizedData{KeyValueList: records, Keys: collectSortedKeys(records)}, nil
}

func parseYAMLToKeyValueList(content []byte) (*NormalizedData, error) {
	var raw any
	if err := yaml.Unmarshal(content, &raw); err != nil {
		return nil, fmt.Errorf("YAMLの解析に失敗しました: %w", err)
	}

	recordsAny, err := parseRecordObjects(raw)
	if err != nil {
		return nil, fmt.Errorf("YAMLの構造が不正です: %w", err)
	}

	records := normalizeRecords(recordsAny)
	return &NormalizedData{KeyValueList: records, Keys: collectSortedKeys(records)}, nil
}

func parseCSVToKeyValueList(content []byte) (*NormalizedData, error) {
	reader := csv.NewReader(bytes.NewReader(content))
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("CSVの解析に失敗しました: %w", err)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("CSVが空です")
	}

	headers := make([]string, 0, len(rows[0]))
	headerExists := map[string]struct{}{}
	for i, h := range rows[0] {
		header := strings.TrimSpace(h)
		if header == "" {
			return nil, fmt.Errorf("CSVヘッダーの%d列目が空です", i+1)
		}
		if _, exists := headerExists[header]; exists {
			return nil, fmt.Errorf("CSVヘッダーに重複があります: %s", header)
		}
		headerExists[header] = struct{}{}
		headers = append(headers, header)
	}

	records := make([]map[string]string, 0, len(rows)-1)
	for rowIndex := 1; rowIndex < len(rows); rowIndex++ {
		row := rows[rowIndex]
		if len(row) > len(headers) {
			return nil, fmt.Errorf("CSVの%d行目の列数がヘッダーより多いです", rowIndex+1)
		}

		record := make(map[string]string, len(headers))
		for i, key := range headers {
			if i < len(row) {
				record[key] = row[i]
				continue
			}
			record[key] = ""
		}
		records = append(records, record)
	}

	return &NormalizedData{KeyValueList: records, Keys: headers}, nil
}

func parseHTMLToKeyValueList(content []byte) (*NormalizedData, error) {
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

		record := make(map[string]string, len(headers))
		for i, key := range headers {
			if i < len(cells) {
				record[key] = cells[i]
				continue
			}
			record[key] = ""
		}
		records = append(records, record)
	}

	return &NormalizedData{KeyValueList: records, Keys: headers}, nil
}

func serializeToJSON(records []map[string]string) ([]byte, error) {
	b, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("JSON生成に失敗しました: %w", err)
	}
	return append(b, '\n'), nil
}

func serializeToYAML(records []map[string]string) ([]byte, error) {
	b, err := yaml.Marshal(records)
	if err != nil {
		return nil, fmt.Errorf("YAML生成に失敗しました: %w", err)
	}
	return b, nil
}

func serializeToCSV(records []map[string]string, keys []string) ([]byte, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("CSV出力に必要なキー情報がありません")
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	if err := writer.Write(keys); err != nil {
		return nil, fmt.Errorf("CSVヘッダーの書き込みに失敗しました: %w", err)
	}

	for _, record := range records {
		row := make([]string, 0, len(keys))
		for _, key := range keys {
			row = append(row, record[key])
		}
		if err := writer.Write(row); err != nil {
			return nil, fmt.Errorf("CSV行の書き込みに失敗しました: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("CSV生成に失敗しました: %w", err)
	}

	return buf.Bytes(), nil
}

func serializeToHTML(records []map[string]string, keys []string) ([]byte, error) {
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

func normalizeFormat(format string) string {
	return strings.ToLower(strings.TrimSpace(format))
}

func parseRecordObjects(raw any) ([]map[string]any, error) {
	switch value := raw.(type) {
	case []any:
		records := make([]map[string]any, 0, len(value))
		for i, elem := range value {
			record, ok := toStringKeyMap(elem)
			if !ok {
				return nil, fmt.Errorf("%d件目がオブジェクトではありません", i+1)
			}
			records = append(records, record)
		}
		return records, nil
	case map[string]any, map[any]any:
		record, ok := toStringKeyMap(value)
		if !ok {
			return nil, fmt.Errorf("オブジェクト形式として解釈できません")
		}
		return []map[string]any{record}, nil
	default:
		return nil, fmt.Errorf("配列またはオブジェクト形式を指定してください")
	}
}

func toStringKeyMap(value any) (map[string]any, bool) {
	switch v := value.(type) {
	case map[string]any:
		converted := make(map[string]any, len(v))
		for key, val := range v {
			converted[key] = normalizeAnyValue(val)
		}
		return converted, true
	case map[any]any:
		converted := make(map[string]any, len(v))
		for key, val := range v {
			converted[fmt.Sprint(key)] = normalizeAnyValue(val)
		}
		return converted, true
	default:
		return nil, false
	}
}

func normalizeAnyValue(value any) any {
	switch v := value.(type) {
	case map[string]any:
		converted := make(map[string]any, len(v))
		for key, child := range v {
			converted[key] = normalizeAnyValue(child)
		}
		return converted
	case map[any]any:
		converted := make(map[string]any, len(v))
		for key, child := range v {
			converted[fmt.Sprint(key)] = normalizeAnyValue(child)
		}
		return converted
	case []any:
		items := make([]any, 0, len(v))
		for _, item := range v {
			items = append(items, normalizeAnyValue(item))
		}
		return items
	default:
		return v
	}
}

func normalizeRecords(recordsAny []map[string]any) []map[string]string {
	records := make([]map[string]string, 0, len(recordsAny))
	for _, recordAny := range recordsAny {
		record := make(map[string]string, len(recordAny))
		for key, val := range recordAny {
			record[key] = stringifyValue(val)
		}
		records = append(records, record)
	}
	return records
}

func stringifyValue(value any) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return v
	case bool,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return fmt.Sprint(v)
	default:
		b, err := json.Marshal(v)
		if err == nil {
			return string(b)
		}
		return fmt.Sprint(v)
	}
}

func collectSortedKeys(records []map[string]string) []string {
	seen := map[string]struct{}{}
	keys := make([]string, 0)

	for _, record := range records {
		recordKeys := make([]string, 0, len(record))
		for key := range record {
			recordKeys = append(recordKeys, key)
		}
		sort.Strings(recordKeys)
		for _, key := range recordKeys {
			if _, exists := seen[key]; exists {
				continue
			}
			seen[key] = struct{}{}
			keys = append(keys, key)
		}
	}

	return keys
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
