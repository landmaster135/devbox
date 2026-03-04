package common

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"
)

// FormatTable はヘッダーと行データを表形式の文字列へ整形する。
func FormatTable(headers []string, rows [][]string) string {
	var buf bytes.Buffer
	writer := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)

	_, _ = fmt.Fprintln(writer, strings.Join(headers, "\t"))
	for _, row := range rows {
		_, _ = fmt.Fprintln(writer, strings.Join(row, "\t"))
	}
	_ = writer.Flush()

	return strings.TrimRight(buf.String(), "\n")
}

// ZoneBasename は zone URL から末尾の zone 名のみを返す。
func ZoneBasename(zone string) string {
	trimmed := strings.TrimSpace(zone)
	if trimmed == "" {
		return ""
	}

	idx := strings.LastIndex(trimmed, "/")
	if idx < 0 {
		return trimmed
	}
	return trimmed[idx+1:]
}
