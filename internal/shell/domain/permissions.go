package domain

import (
	"fmt"
	"strings"
)

// SandboxPermissions はサンドボックス実行ポリシーを表す
// Cmd/CLIでもMCPツールでも同じ語彙を使えるようにドメイン層で定義する
// 値はCodex shellツールに合わせた文字列表現とする
//
//	use_default       : 既定のサンドボックスで実行
//	require_escalated : サンドボックス外実行を明示的に要求
type SandboxPermissions string

const (
	// SandboxPermissionsUseDefault は既定のサンドボックスを意味する
	SandboxPermissionsUseDefault SandboxPermissions = "use_default"
	// SandboxPermissionsRequireEscalated はサンドボックス外実行を意味する
	SandboxPermissionsRequireEscalated SandboxPermissions = "require_escalated"
)

// ParseSandboxPermissions は文字列をSandboxPermissionsに変換する
func ParseSandboxPermissions(value string) (SandboxPermissions, error) {
	if value == "" {
		return SandboxPermissionsUseDefault, nil
	}

	normalized := strings.ToLower(strings.TrimSpace(value))
	switch normalized {
	case "use_default", "use-default", "default", "default_sandbox":
		return SandboxPermissionsUseDefault, nil
	case "require_escalated", "require-escalated", "require-escalation", "require_escalation", "escalated":
		return SandboxPermissionsRequireEscalated, nil
	default:
		return "", fmt.Errorf("無効なsandbox_permissionsです: %s", value)
	}
}

// RequiresJustification はエスカレーション要求かどうかを返す
func (s SandboxPermissions) RequiresJustification() bool {
	return s == SandboxPermissionsRequireEscalated
}

// String はフレンドリーな文字列表現を返す
func (s SandboxPermissions) String() string {
	if s == "" {
		return string(SandboxPermissionsUseDefault)
	}
	return string(s)
}

// Validate は有効な値かどうかを判定する
func (s SandboxPermissions) Validate() error {
	switch s {
	case SandboxPermissionsUseDefault, SandboxPermissionsRequireEscalated:
		return nil
	case "":
		return fmt.Errorf("sandbox_permissionsが指定されていません")
	default:
		return fmt.Errorf("サポートされていないsandbox_permissionsです: %s", string(s))
	}
}
