package common

import "strings"

// IsBlank は文字列が空白のみかどうかを判定する。
func IsBlank(value string) bool {
	return strings.TrimSpace(value) == ""
}

// ShellQuote は値をシェル向けに単一引用符で安全に囲む。
func ShellQuote(value string) string {
	escaped := strings.ReplaceAll(value, "'", "'\"'\"'")
	return "'" + escaped + "'"
}
