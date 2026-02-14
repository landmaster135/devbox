package mdunorderedlistfmt

import (
	"fmt"
	"strings"

	"github.com/landmaster135/devbox/internal/data_converter/usecases/formats/common"
)

const itemKey = "item"

// Parse は Markdown unordered-list を key-value リストへ正規化します。
func Parse(content []byte) (*common.NormalizedData, error) {
	lines := strings.Split(normalizeLineBreaks(string(content)), "\n")

	records := make([]map[string]string, 0, len(lines))
	for rowIndex, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if line == "" {
			continue
		}

		value, ok := parseUnorderedListLine(line)
		if !ok {
			return nil, fmt.Errorf("Markdown unordered-listの%d行目が不正です", rowIndex+1)
		}
		records = append(records, map[string]string{itemKey: value})
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("Markdown unordered-listが空です")
	}

	return common.NewNormalizedData(records, []string{itemKey}), nil
}

// Serialize は key-value リストを Markdown unordered-list へ変換します。
func Serialize(records []map[string]string, keys []string) ([]byte, error) {
	resolvedItemKey, err := resolveItemKey(keys, "md-unordered-list")
	if err != nil {
		return nil, err
	}

	var builder strings.Builder
	for _, record := range records {
		builder.WriteString(fmt.Sprintf("- %s\n", record[resolvedItemKey]))
	}
	return []byte(builder.String()), nil
}

func parseUnorderedListLine(line string) (string, bool) {
	if len(line) < 3 {
		return "", false
	}
	marker := line[0]
	if marker != '-' && marker != '*' && marker != '+' {
		return "", false
	}
	if line[1] != ' ' && line[1] != '\t' {
		return "", false
	}
	value := strings.TrimSpace(line[2:])
	if value == "" {
		return "", false
	}
	return value, true
}

func resolveItemKey(keys []string, formatName string) (string, error) {
	if len(keys) == 0 {
		return "", fmt.Errorf("%s出力に必要なキー情報がありません", formatName)
	}
	if len(keys) == 1 {
		return keys[0], nil
	}
	for _, key := range keys {
		if key == itemKey {
			return itemKey, nil
		}
	}
	return "", fmt.Errorf("%s は単一キーのデータのみ出力できます", formatName)
}

func normalizeLineBreaks(input string) string {
	return strings.ReplaceAll(input, "\r\n", "\n")
}
