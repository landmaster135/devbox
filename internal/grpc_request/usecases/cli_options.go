package usecases

import (
	"fmt"
	"time"
)

// CLIOptions はCLIオプションの情報を表します
type CLIOptions struct {
	Server       string        `json:"server"`
	Method       string        `json:"method"`
	JSONFile     string        `json:"json_file"`
	UseTLS       bool          `json:"use_tls"`
	Token        string        `json:"token"`
	Timeout      time.Duration `json:"timeout"`
	TestConn     bool          `json:"test_conn"`
	ListServices bool          `json:"list_services"`
}

// NewCLIOptions は新しいCLIOptionsインスタンスを作成します
func NewCLIOptions() *CLIOptions {
	return &CLIOptions{
		Timeout: 30 * time.Second,
	}
}

// Validate はCLIオプションの妥当性を検証します
func (opts *CLIOptions) Validate() error {
	if opts.Server == "" {
		return fmt.Errorf("サーバーアドレスが指定されていません")
	}

	// 接続テストまたはサービス一覧表示の場合は、他のオプションは不要
	if opts.TestConn || opts.ListServices {
		return nil
	}

	// 通常のリクエストの場合は、メソッドとJSONファイルが必要
	if opts.Method == "" {
		return fmt.Errorf("メソッドが指定されていません")
	}

	if opts.JSONFile == "" {
		return fmt.Errorf("リクエストデータのJSONファイルが指定されていません")
	}

	return nil
}
