package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// TransportStrategy はトランスポート戦略を表す型
type TransportStrategy string

const (
	TransportSSEOnly   TransportStrategy = "sse-only"
	TransportHTTPOnly  TransportStrategy = "http-only"
	TransportSSEFirst  TransportStrategy = "sse-first"
	TransportHTTPFirst TransportStrategy = "http-first"
)

// Config はmcp-remote CLIの設定を保持する構造体
type Config struct {
	ServerURL                   string                       // リモートサーバーのURL（必須）
	CallbackPort                int                          // OAuthコールバック用ポート
	Headers                     map[string]string            // カスタムHTTPヘッダー
	TransportStrategy           TransportStrategy            // トランスポート戦略
	Host                        string                       // コールバックホスト
	Debug                       bool                         // デバッグモード
	AllowHTTP                   bool                         // HTTP接続を許可
	StaticOAuthClientMetadata   map[string]interface{}       // 静的OAuthクライアントメタデータ
	StaticOAuthClientInfo       map[string]interface{}       // 静的OAuthクライアント情報
	AuthorizeResource           string                       // 認可リソース
	Help                        bool                         // ヘルプ表示フラグ
}

// NewConfig は新しいConfigを作成する
func NewConfig(serverURL string) (*Config, error) {
	if serverURL == "" {
		return nil, fmt.Errorf("サーバーURLが指定されていません")
	}

	return &Config{
		ServerURL:         serverURL,
		CallbackPort:      0, // 0の場合は自動選択
		Headers:           make(map[string]string),
		TransportStrategy: TransportHTTPFirst, // デフォルト
		Host:              "localhost",
		Debug:             false,
		AllowHTTP:         false,
	}, nil
}

// ParseFlags はコマンドライン引数を解析してConfigを作成する
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

// ParseFlagsWithParser は指定されたFlagParserを使用してコマンドライン引数を解析する
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	var (
		serverURL                 = ""
		callbackPortStr           = "0"
		headers                   = ""
		transportStrategyStr      = "http-first"
		host                      = "localhost"
		debug                     = false
		allowHTTP                 = false
		staticOAuthClientMetadata = ""
		staticOAuthClientInfo     = ""
		authorizeResource         = ""
		help                      = false
	)

	parser.StringVar(&serverURL, "server-url", serverURL, "リモートMCPサーバーのURL（必須）")
	parser.StringVar(&serverURL, "s", serverURL, "サーバーURLの短縮形")

	parser.StringVar(&callbackPortStr, "callback-port", callbackPortStr, "OAuthコールバック用ポート（0で自動選択）")
	parser.StringVar(&callbackPortStr, "p", callbackPortStr, "ポートの短縮形")

	parser.StringVar(&headers, "header", headers, "カスタムHTTPヘッダー（Key:Value形式、複数指定可能）")
	parser.StringVar(&headers, "H", headers, "ヘッダーの短縮形")

	parser.StringVar(&transportStrategyStr, "transport", transportStrategyStr, "トランスポート戦略 (sse-only, http-only, sse-first, http-first)")
	parser.StringVar(&transportStrategyStr, "t", transportStrategyStr, "トランスポートの短縮形")

	parser.StringVar(&host, "host", host, "コールバックホスト")

	parser.BoolVar(&debug, "debug", debug, "デバッグモードを有効にする")
	parser.BoolVar(&debug, "d", debug, "デバッグの短縮形")

	parser.BoolVar(&allowHTTP, "allow-http", allowHTTP, "HTTP接続を許可する（信頼できるプライベートネットワークでのみ使用）")

	parser.StringVar(&staticOAuthClientMetadata, "static-oauth-client-metadata", staticOAuthClientMetadata, "静的OAuthクライアントメタデータ（JSON文字列または@ファイルパス）")
	parser.StringVar(&staticOAuthClientInfo, "static-oauth-client-info", staticOAuthClientInfo, "静的OAuthクライアント情報（JSON文字列または@ファイルパス）")
	parser.StringVar(&authorizeResource, "resource", authorizeResource, "認可リソース")

	parser.BoolVar(&help, "help", help, "ヘルプを表示")
	parser.BoolVar(&help, "h", help, "ヘルプの短縮形")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %v", err)
	}

	// ヘルプが要求された場合
	if help {
		return &Config{Help: true}, nil
	}

	// 位置引数からサーバーURLを取得
	args := parser.Args()
	if len(args) >= 1 && serverURL == "" {
		serverURL = args[0]
	}

	// サーバーURLの検証
	if serverURL == "" {
		return nil, fmt.Errorf("サーバーURLが指定されていません")
	}

	// ポート番号の変換
	callbackPort, err := strconv.Atoi(callbackPortStr)
	if err != nil {
		return nil, fmt.Errorf("無効なポート番号です: %s", callbackPortStr)
	}

	// トランスポート戦略の検証
	transportStrategy := TransportStrategy(transportStrategyStr)
	if !isValidTransportStrategy(transportStrategy) {
		return nil, fmt.Errorf("無効なトランスポート戦略です: %s", transportStrategyStr)
	}

	// ヘッダーの解析
	headerMap := make(map[string]string)
	if headers != "" {
		parts := strings.Split(headers, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if colonIndex := strings.Index(part, ":"); colonIndex != -1 {
				key := strings.TrimSpace(part[:colonIndex])
				value := strings.TrimSpace(part[colonIndex+1:])
				headerMap[key] = value
			} else {
				return nil, fmt.Errorf("無効なヘッダー形式です: %s", part)
			}
		}
	}

	config := &Config{
		ServerURL:         serverURL,
		CallbackPort:      callbackPort,
		Headers:           headerMap,
		TransportStrategy: transportStrategy,
		Host:              host,
		Debug:             debug,
		AllowHTTP:         allowHTTP,
		AuthorizeResource: authorizeResource,
	}

	// 静的OAuthクライアントメタデータの処理
	if staticOAuthClientMetadata != "" {
		// TODO: JSON解析の実装
		config.StaticOAuthClientMetadata = make(map[string]interface{})
	}

	// 静的OAuthクライアント情報の処理
	if staticOAuthClientInfo != "" {
		// TODO: JSON解析の実装
		config.StaticOAuthClientInfo = make(map[string]interface{})
	}

	return config, nil
}

// isValidTransportStrategy はトランスポート戦略が有効かどうかを確認する
func isValidTransportStrategy(strategy TransportStrategy) bool {
	switch strategy {
	case TransportSSEOnly, TransportHTTPOnly, TransportSSEFirst, TransportHTTPFirst:
		return true
	default:
		return false
	}
}

// PrintUsage は使用方法を表示する
func PrintUsage() {
	fmt.Fprintf(os.Stderr, `MCP Remote Proxy CLI

リモートMCPサーバーに接続するためのプロキシツール

使用方法:
  基本的な使用:
    %s https://example.com/mcp/sse
    %s -s https://example.com/mcp/sse -p 3334

  カスタムヘッダー付き:
    %s https://example.com/mcp/sse -H "Authorization:Bearer token123"

  トランスポート戦略指定:
    %s https://example.com/mcp/sse -t sse-only
    %s https://example.com/mcp/sse -t http-first

  デバッグモード:
    %s https://example.com/mcp/sse -d

オプション:
  -server-url, -s      リモートMCPサーバーのURL（必須）
  -callback-port, -p   OAuthコールバック用ポート（0で自動選択）
  -header, -H          カスタムHTTPヘッダー（Key:Value形式）
  -transport, -t       トランスポート戦略 (sse-only, http-only, sse-first, http-first)
  -host                コールバックホスト（デフォルト: localhost）
  -debug, -d           デバッグモードを有効にする
  -allow-http          HTTP接続を許可する（信頼できるプライベートネットワークでのみ使用）
  -static-oauth-client-metadata  静的OAuthクライアントメタデータ
  -static-oauth-client-info      静的OAuthクライアント情報
  -resource            認可リソース
  -help, -h            このヘルプを表示

例:
  # 基本的な接続
  %s https://remote.mcp.server/sse

  # カスタムポートとヘッダー
  %s https://remote.mcp.server/sse -p 9696 -H "Authorization:Bearer \${TOKEN}"

  # デバッグモードでSSE専用
  %s https://remote.mcp.server/sse -d -t sse-only

`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}
