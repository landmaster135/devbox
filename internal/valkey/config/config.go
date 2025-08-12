package config

import (
	"fmt"
	"os"
	"strings"
)

// Config はValkey CLIの設定を保持する構造体
type Config struct {
	Operation string   // 操作タイプ (get-keys, get-value, get-type, set-value, delete-key, delete-keys, select-keys, get-all-values, delete-data)
	Key       string   // 単一キー
	Keys      []string // 複数キー（カンマ区切り文字列から変換）
	Pattern   string   // パターン
	Value     string   // 値
	All       bool     // 全件フラグ
	DryRun    bool     // ドライランフラグ
	Host      string   // Valkeyホスト
	Port      int      // Valkeyポート
	Password  string   // パスワード
	Database  int      // データベース番号
	Help      bool     // ヘルプ表示フラグ
}

// NewConfig は新しいConfigを作成する
func NewConfig(operation string) (*Config, error) {
	if operation == "" {
		return nil, fmt.Errorf("操作タイプが指定されていません")
	}

	// 操作タイプの検証
	validOperations := []string{
		"get-keys", "get-value", "get-type", "set-value", "delete-key",
		"delete-keys", "select-keys", "get-all-values", "delete-data",
	}
	isValid := false
	for _, op := range validOperations {
		if operation == op {
			isValid = true
			break
		}
	}
	if !isValid {
		return nil, fmt.Errorf("無効な操作タイプです: %s", operation)
	}

	return &Config{
		Operation: operation,
		Host:      "localhost", // デフォルト値
		Port:      6379,        // デフォルト値
		Database:  0,           // デフォルト値
	}, nil
}

// BuildValkeyURL はConfigからValkey接続URLを構築する
func (c *Config) BuildValkeyURL() string {
	var url string

	if c.Password != "" {
		// パスワードありの場合
		url = fmt.Sprintf("valkey://:%s@%s:%d/%d",
			c.Password, c.Host, c.Port, c.Database)
	} else {
		// パスワードなしの場合
		url = fmt.Sprintf("valkey://%s:%d/%d",
			c.Host, c.Port, c.Database)
	}

	return url
}

// ParseFlagsWithParser は指定されたFlagParserを使用してコマンドライン引数を解析する
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	var (
		operation = ""
		key       = ""
		keys      = ""
		pattern   = ""
		value     = ""
		all       = false
		dryRun    = false
		host      = "localhost"
		port      = 6379
		password  = ""
		database  = 0
		help      = false
	)

	parser.StringVar(&operation, "operation", operation, "Valkey操作 (get-keys, get-value, get-type, set-value, delete-key, delete-keys, select-keys, get-all-values, delete-data)")
	parser.StringVar(&operation, "o", operation, "操作の短縮形")

	// データ操作のパラメータ
	parser.StringVar(&key, "key", key, "単一キー")
	parser.StringVar(&key, "k", key, "キーの短縮形")
	parser.StringVar(&keys, "keys", keys, "複数キー（カンマ区切り）")
	parser.StringVar(&pattern, "pattern", pattern, "パターン")
	parser.StringVar(&pattern, "p", pattern, "パターンの短縮形")
	parser.StringVar(&value, "value", value, "値")
	parser.StringVar(&value, "v", value, "値の短縮形")
	parser.BoolVar(&all, "all", all, "全件フラグ")
	parser.BoolVar(&all, "a", all, "全件フラグの短縮形")
	parser.BoolVar(&dryRun, "dry-run", dryRun, "ドライランフラグ")
	parser.BoolVar(&dryRun, "dr", dryRun, "ドライランフラグの短縮形")

	// 接続関連のパラメータ
	parser.StringVar(&host, "host", host, "Valkeyホスト")
	parser.StringVar(&host, "h", host, "ホストの短縮形")
	parser.IntVar(&port, "port", port, "Valkeyポート")
	parser.StringVar(&password, "password", password, "Valkeyパスワード")
	parser.StringVar(&password, "pass", password, "パスワードの短縮形")
	parser.IntVar(&database, "database", database, "データベース番号")
	parser.IntVar(&database, "db", database, "データベース番号の短縮形")

	parser.BoolVar(&help, "help", help, "ヘルプを表示")
	parser.BoolVar(&help, "help-flag", help, "ヘルプの短縮形")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %v", err)
	}

	// ヘルプが要求された場合
	if help {
		return &Config{Help: true}, nil
	}

	// 操作タイプが指定されていない場合のエラーチェック
	if operation == "" {
		return nil, fmt.Errorf("操作タイプが指定されていません")
	}

	config, err := NewConfig(operation)
	if err != nil {
		return nil, err
	}

	// パラメータを設定
	config.Key = key
	config.Pattern = pattern
	config.Value = value
	config.All = all
	config.DryRun = dryRun
	config.Host = host
	config.Port = port
	config.Password = password
	config.Database = database

	// keys文字列を[]stringに変換
	if keys != "" {
		parts := strings.Split(keys, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				config.Keys = append(config.Keys, part)
			}
		}
	}

	return config, nil
}

// ParseFlags はコマンドライン引数を解析してConfigを作成する
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

// PrintUsage は使用方法を表示する
func PrintUsage() {
	fmt.Fprintf(os.Stderr, `Valkey CLIツール

使用方法:
  キー取得:
    %s -operation get-keys -pattern "user:*"

  値取得:
    %s -operation get-value -key "user:123"

  型取得:
    %s -operation get-type -key "user:123"

  値設定:
    %s -operation set-value -key "user:123" -value '{"name":"John"}'

  キー削除:
    %s -operation delete-key -key "user:123"

  複数キー削除:
    %s -operation delete-keys -keys "user:123,user:456"

  値選択:
    %s -operation select-keys -key "user:123"
    %s -operation select-keys -keys "user:123,user:456"
    %s -operation select-keys -pattern "user:*"
    %s -operation select-keys -all

  全値取得:
    %s -operation get-all-values -keys "user:123,user:456"

  データ削除（ドライラン）:
    %s -operation delete-data -pattern "temp:*" -dry-run

  パスワード付き接続:
    %s -operation get-keys -pattern "*" -password "mypassword"

  特定データベース接続:
    %s -operation get-keys -pattern "*" -database 1

  短縮形:
    %s -o get-keys -p "user:*" -h "localhost" -port 6379 -pass "secret" -db 1

オプション:
  -operation, -o       Valkey操作 (get-keys, get-value, get-type, set-value, delete-key, delete-keys, select-keys, get-all-values, delete-data)
  -key, -k             単一キー
  -keys                複数キー（カンマ区切り）
  -pattern, -p         パターン
  -value, -v           値
  -all, -a             全件フラグ
  -dry-run, -dr        ドライランフラグ
  -host, -h            Valkeyホスト (デフォルト: localhost)
  -port                Valkeyポート (デフォルト: 6379)
  -password, -pass     Valkeyパスワード
  -database, -db       データベース番号 (デフォルト: 0)
  -help                このヘルプを表示

操作タイプ:
  get-keys             パターンに一致するキーを取得
  get-value            キーの値を取得
  get-type             キーの型を取得
  set-value            キーに値を設定
  delete-key           単一キーを削除
  delete-keys          複数キーを削除
  select-keys        条件に基づいて値を選択
  get-all-values       指定されたキーの全値を取得
  delete-data          条件に基づいてデータを削除

`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}
