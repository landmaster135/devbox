package common

import (
	"fmt"
	"strings"
)

const defaultSSHKeyPath = "$HOME/.ssh/google_compute_engine"

// IsBlank は文字列が空白のみかどうかを判定する。
func IsBlank(value string) bool {
	return strings.TrimSpace(value) == ""
}

// ShellQuote は値をシェル向けに単一引用符で安全に囲む。
func ShellQuote(value string) string {
	escaped := strings.ReplaceAll(value, "'", "'\"'\"'")
	return "'" + escaped + "'"
}

// ShellQuoteSSHKeyPath は SSH 鍵パスをシェル向けに引用する。
// 既定パスは $HOME 展開を維持するため二重引用符を使う。
func ShellQuoteSSHKeyPath(value string) string {
	if value == defaultSSHKeyPath {
		return "\"$HOME/.ssh/google_compute_engine\""
	}
	return ShellQuote(value)
}

// BuildSSHAgentSetupCommand は ssh-agent 起動と ssh-add を行うコマンドを生成する。
func BuildSSHAgentSetupCommand(sshKeyPath string) string {
	return fmt.Sprintf(
		"if [ -z \"${SSH_AUTH_SOCK:-}\" ]; then eval \"$(ssh-agent -s)\" >/dev/null; fi && ssh-add %s",
		ShellQuoteSSHKeyPath(sshKeyPath),
	)
}
