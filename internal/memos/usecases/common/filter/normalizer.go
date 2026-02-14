package filter

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var timestampComparisonPattern = regexp.MustCompile(`\b(created_ts|updated_ts)\b\s*(==|!=|>=|<=|>|<)\s*('([^'\\]*)'|"([^"\\]*)")`)

// NormalizeFilter は created_ts/updated_ts の RFC3339 文字列比較を Unix 秒比較へ正規化する。
func NormalizeFilter(raw string) (string, error) {
	matches := timestampComparisonPattern.FindAllStringSubmatchIndex(raw, -1)
	if len(matches) == 0 {
		return raw, nil
	}

	buf := make([]byte, 0, len(raw))
	last := 0
	for _, match := range matches {
		field := raw[match[2]:match[3]]

		value, ok := extractQuotedValue(raw, match)
		if !ok {
			return "", fmt.Errorf("filter の %s 比較値を読み取れませんでした", field)
		}

		parsedTime, err := time.Parse(time.RFC3339Nano, value)
		if err != nil {
			return "", fmt.Errorf(
				"filter の %s 比較値 %q は RFC3339/RFC3339Nano（タイムゾーン必須）で指定してください: %w",
				field,
				value,
				err,
			)
		}

		buf = append(buf, raw[last:match[6]]...)
		buf = append(buf, strconv.FormatInt(parsedTime.Unix(), 10)...)
		buf = append(buf, raw[match[7]:match[1]]...)
		last = match[1]
	}
	buf = append(buf, raw[last:]...)

	return string(buf), nil
}

func extractQuotedValue(raw string, match []int) (string, bool) {
	singleQuotedStart := match[8]
	singleQuotedEnd := match[9]
	if singleQuotedStart >= 0 && singleQuotedEnd >= 0 {
		return raw[singleQuotedStart:singleQuotedEnd], true
	}

	doubleQuotedStart := match[10]
	doubleQuotedEnd := match[11]
	if doubleQuotedStart >= 0 && doubleQuotedEnd >= 0 {
		return raw[doubleQuotedStart:doubleQuotedEnd], true
	}

	return "", false
}
