package mdorderedlistfmt

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/landmaster135/devbox/internal/data_converter/usecases/formats/common"
)

const itemKey = "item"

// Parse は Markdown ordered-list を key-value リストへ正規化します。
func Parse(content []byte) (*common.NormalizedData, error) {
	lines := strings.Split(normalizeLineBreaks(string(content)), "\n")

	records := make([]map[string]string, 0, len(lines))
	expectedNumber := 1
	for rowIndex, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if line == "" {
			continue
		}

		number, value, ok := parseOrderedListLine(line)
		if !ok {
			return nil, fmt.Errorf("Markdown ordered-listの%d行目が不正です", rowIndex+1)
		}
		if number != expectedNumber {
			return nil, fmt.Errorf("Markdown ordered-listの%d行目の番号が連番ではありません", rowIndex+1)
		}

		records = append(records, map[string]string{itemKey: value})
		expectedNumber++
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("Markdown ordered-listが空です")
	}

	return common.NewNormalizedData(records, []string{itemKey}), nil
}

// Serialize は key-value リストを Markdown ordered-list へ変換します。
func Serialize(records []map[string]string, keys []string) ([]byte, error) {
	resolvedItemKey, err := resolveItemKey(keys, "md-ordered-list")
	if err != nil {
		return nil, err
	}

	var builder strings.Builder
	for i, record := range records {
		builder.WriteString(fmt.Sprintf("%d. %s\n", i+1, record[resolvedItemKey]))
	}
	return []byte(builder.String()), nil
}

func parseOrderedListLine(line string) (int, string, bool) {
	digitEnd := 0
	for digitEnd < len(line) && line[digitEnd] >= '0' && line[digitEnd] <= '9' {
		digitEnd++
	}
	if digitEnd == 0 || digitEnd >= len(line) {
		return 0, "", false
	}
	if line[digitEnd] != '.' {
		return 0, "", false
	}

	valueStart := digitEnd + 1
	if valueStart >= len(line) || (line[valueStart] != ' ' && line[valueStart] != '\t') {
		return 0, "", false
	}
	for valueStart < len(line) && (line[valueStart] == ' ' || line[valueStart] == '\t') {
		valueStart++
	}
	if valueStart >= len(line) {
		return 0, "", false
	}

	number, err := strconv.Atoi(line[:digitEnd])
	if err != nil {
		return 0, "", false
	}
	value := strings.TrimSpace(line[valueStart:])
	if value == "" {
		return 0, "", false
	}
	return number, value, true
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
