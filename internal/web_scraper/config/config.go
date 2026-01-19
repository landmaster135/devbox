package config

import (
	"fmt"
	"strings"
	"time"
)

const (
	// OperationGetDOMTree は指定したURLのDOMツリーを取得する操作です。
	OperationGetDOMTree = "get_dom_tree"
)

// Config はCLIから受け取る設定値を保持します。
type Config struct {
	Operation   string
	URL         string
	WaitSeconds int
	OutputPath  string
}

// Validate は設定内容を検証し、不正な場合はエラーを返します。
func (c Config) Validate() error {
	op := normalizeOperation(c.Operation)
	if op == "" {
		return fmt.Errorf("operationを指定してください")
	}

	switch op {
	case OperationGetDOMTree:
		if strings.TrimSpace(c.URL) == "" {
			return fmt.Errorf("get_dom_treeではurlを指定してください")
		}
		if c.WaitSeconds < 0 {
			return fmt.Errorf("wait-secondsは0以上を指定してください")
		}
		if c.OutputPath != "" && strings.TrimSpace(c.OutputPath) == "" {
			return fmt.Errorf("output-fileは空文字では指定できません")
		}
	default:
		return fmt.Errorf("未対応のoperationです: %s", c.Operation)
	}

	return nil
}

// OperationName は整形済みのoperation名を返します。
func (c Config) OperationName() string {
	return normalizeOperation(c.Operation)
}

// WaitDuration は待機秒数をtime.Durationへ変換して返します。
func (c Config) WaitDuration() time.Duration {
	if c.WaitSeconds <= 0 {
		return 0
	}
	return time.Duration(c.WaitSeconds) * time.Second
}

// OutputFilePath は整形済みの出力ファイルパスを返します。
func (c Config) OutputFilePath() string {
	return strings.TrimSpace(c.OutputPath)
}

func normalizeOperation(op string) string {
	return strings.ToLower(strings.TrimSpace(op))
}
