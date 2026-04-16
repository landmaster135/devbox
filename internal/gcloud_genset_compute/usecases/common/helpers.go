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

// BuildSSHKeyCreationCommand は SSH 鍵生成のコマンドを返す。
// force が false の場合、既存鍵を検知したら終了する。
// force が true の場合、既存鍵を検知したらログ出力後に上書き生成する。
func BuildSSHKeyCreationCommand(sshKeyPath string, force bool) string {
	pathAssignment := fmt.Sprintf("ssh_key_path=%s", ShellQuoteSSHKeyPath(sshKeyPath))
	if force {
		return fmt.Sprintf(
			"%s; if [ -f \"$ssh_key_path\" ]; then echo \"既存SSH秘密鍵を上書きしました: $ssh_key_path\" >&2; rm -f \"$ssh_key_path\" \"$ssh_key_path.pub\"; fi; ssh-keygen -t rsa -f \"$ssh_key_path\"",
			pathAssignment,
		)
	}

	return fmt.Sprintf(
		"%s; if [ -f \"$ssh_key_path\" ]; then echo \"SSH秘密鍵は既に存在します: $ssh_key_path。上書きするには -forces=true を指定してください\" >&2; exit 1; fi; ssh-keygen -t rsa -f \"$ssh_key_path\"",
		pathAssignment,
	)
}
