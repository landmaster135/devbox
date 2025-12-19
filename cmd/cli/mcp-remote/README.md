# MCP Remote - Go Implementation

Node.js版[mcp-remote](https://github.com/mark3labs/mcp-remote)を参考にしたGoでのMCP Remote CLIツール実装です。

## 概要

MCP Remoteは、リモートMCPサーバーに接続するためのプロキシツールです。STDIOトランスポートとリモートサーバー間で双方向通信を行い、リモートサーバーのツールを透明にローカルで利用できるようにします。

## 特徴

- **双方向プロキシ機能**: STDIOとリモートサーバー間の完全な双方向通信
- **動的ツール中継**: リモートサーバーのツールを自動取得し、透明にプロキシ
- **マルチトランスポート対応**: HTTP、SSE（Server-Sent Events）をサポート
- **Node.js版互換**: 同じコマンドライン引数とオプションをサポート
- **単一バイナリ**: Goの利点を活かした単一バイナリでの配布
- **Clean Architecture**: 保守性の高い設計

## インストール

```bash
# リポジトリをクローン
git clone <repository-url>
cd devbox

# ビルド
go build -o bin/mcp-remote ./cmd/cli/mcp-remote

# 実行可能ファイルをPATHに追加（オプション）
sudo cp bin/mcp-remote /usr/local/bin/
```

## 使用方法

### 基本的な使用

```bash
# HTTPサーバーに接続
./bin/mcp-remote http://localhost:8080/mcp --allow-http

# HTTPSサーバーに接続
./bin/mcp-remote https://remote.mcp.server/mcp
```

### オプション

```
使用方法:
  基本的な使用:
    ./bin/mcp-remote https://example.com/mcp/sse
    ./bin/mcp-remote -s https://example.com/mcp/sse -p 3334

  カスタムヘッダー付き:
    ./bin/mcp-remote https://example.com/mcp/sse -H "Authorization:Bearer token123"

  トランスポート戦略指定:
    ./bin/mcp-remote https://example.com/mcp/sse -t sse-only
    ./bin/mcp-remote https://example.com/mcp/sse -t http-first

  デバッグモード:
    ./bin/mcp-remote https://example.com/mcp/sse -d

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
```

## 使用例

```bash
# 基本的な接続
./bin/mcp-remote https://remote.mcp.server/sse

# カスタムポートとヘッダー
./bin/mcp-remote https://remote.mcp.server/sse -p 9696 -H "Authorization:Bearer ${TOKEN}"

# デバッグモードでSSE専用
./bin/mcp-remote https://remote.mcp.server/sse -d -t sse-only

# ローカルHTTPサーバーに接続（開発用）
./bin/mcp-remote http://localhost:8080/mcp -d --allow-http
```

## アーキテクチャ

### プロジェクト構造

```
cmd/cli/mcp-remote/
├── main.go                 # メインエントリーポイント
└── README.md               # このファイル

internal/mcp_remote/
├── config/                 # 設定管理
│   ├── config.go           # 設定構造体とフラグ解析
│   ├── config_test.go      # 設定のテスト
│   └── flag_parser.go      # フラグパーサーインターフェース
├── usecases/               # ビジネスロジック層
│   ├── services.go         # プロキシサービス
│   └── services_test.go    # サービスのテスト
├── interface/transport/    # トランスポート層
│   ├── proxy.go            # 双方向プロキシ実装
│   ├── sse_client.go       # SSEクライアント
│   └── http_client.go      # HTTPクライアント
└── infrastructure/utils/   # インフラストラクチャ層
    ├── logger.go           # ログ機能
    ├── url_hash.go         # URL ハッシュ生成
    └── port_finder.go      # 利用可能ポート検索
```

### 動作原理

1. **接続確立**: リモートMCPサーバーに接続し、利用可能なツール一覧を取得
2. **ツール登録**: 取得したツールをローカルMCPサーバーにプロキシツールとして登録
3. **双方向中継**: STDIOからのリクエストをリモートサーバーに転送し、レスポンスを返却
4. **透明性**: クライアントからはローカルツールと同じように見える

## 開発

### 依存関係

- Go
- `github.com/mark3labs/mcp-go` - MCPプロトコル実装
- `github.com/stretchr/testify` - テストフレームワーク

### テスト実行

```bash
# 全テスト実行
go test ./internal/mcp_remote/...

# カバレッジ付きテスト
go test -coverprofile=coverage.out ./internal/mcp_remote/...
go tool cover -html=coverage.out -o coverage.html
```

### ビルド

```bash
# 開発用ビルド
go build -o bin/mcp-remote ./cmd/cli/mcp-remote

# リリース用ビルド（複数プラットフォーム）
GOOS=linux GOARCH=amd64 go build -o bin/mcp-remote-linux-amd64 ./cmd/cli/mcp-remote
GOOS=darwin GOARCH=amd64 go build -o bin/mcp-remote-darwin-amd64 ./cmd/cli/mcp-remote
GOOS=windows GOARCH=amd64 go build -o bin/mcp-remote-windows-amd64.exe ./cmd/cli/mcp-remote
```

## トラブルシューティング

### よくある問題

**1. 接続エラー**
```
Error: HTTP接続テストに失敗しました
```
- サーバーURLが正しいか確認
- HTTPサーバーの場合は`--allow-http`フラグを使用
- ネットワーク接続を確認

**2. ツールが見つからない**
```
Error: リモートサーバーからのツール一覧取得に失敗しました
```
- リモートサーバーがMCPプロトコルに対応しているか確認
- デバッグモード（`-d`）でログを確認

**3. 認証エラー**
```
Error: 401 Unauthorized
```
- 認証ヘッダーが正しく設定されているか確認
- トークンの有効期限を確認

### デバッグ

デバッグモードを有効にすると、詳細なログが出力されます：

```bash
./bin/mcp-remote https://example.com/mcp -d
```

ログ出力例：
```
[mcp-remote] MCP Remote Proxy を開始します
[proxy] プロキシを開始します: https://example.com/mcp
[proxy] HTTPクライアントを作成中...
[proxy] HTTP接続が確認されました
[proxy] リモートサーバーから4個のツールを取得しました
[proxy] ツールを追加中: patch-page
```

## Node.js版との違い

### 互換性

- **コマンドライン引数**: 完全互換
- **トランスポート戦略**: 完全互換
- **動作**: 透明なプロキシとして同じ動作

### Go版の利点

- **パフォーマンス**: より高速な実行
- **配布**: 単一バイナリでの配布
- **メモリ使用量**: より効率的なメモリ使用
- **型安全性**: コンパイル時の型チェック

## ライセンス

このプロジェクトは[元のmcp-remote](https://github.com/mark3labs/mcp-remote)と同じライセンスに従います。

## 関連リンク

- [MCP (Model Context Protocol)](https://modelcontextprotocol.io/)
- [mcp-go](https://github.com/mark3labs/mcp-go)
- [元のmcp-remote (Node.js版)](https://github.com/mark3labs/mcp-remote)
