# Notion Sync

Notionのページにブロックを追加するためのCLIツールです。notion-synchronizerサーバーのページパッチAPIエンドポイントにHTTPリクエストを送信します。

## 機能

- **Notionトークン認証**: Notionトークンを使用した安全な認証
- **柔軟なページ指定**: ページIDまたはコンテンツIDによるページ指定
- **マークダウンサポート**: マークダウンコンテンツの送信
- **ヘッダー制御**: ヘッダーレベル（H1、H2、H3）のトグル制御
- **カスタムエンドポイント**: APIエンドポイントURLの指定
- **短縮オプション**: 全てのオプションに短縮形を提供
- **位置引数サポート**: フラグ指定と位置引数の両方に対応

## インストール

```bash
# プロジェクトルートから
go build -o bin/notion-sync ./cmd/cli/notion-sync
```

## 使用方法

### 基本的な使用方法

```bash
./bin/notion-sync -operation "patch" -token "your_token" -page-id "page_id" -markdown "# Hello World" -endpoint-url "http://localhost:8080/task/patch"
```

### オプション

| オプション | 短縮形 | 説明 | 必須 | 例 |
|-----------|--------|------|------|-----|
| `-operation` | `-o` | 操作タイプ（現在サポート: patch） | * | `-o "patch"` |
| `-token` | `-t` | Notionトークン | * | `-t "your_token"` |
| `-con-id` | `-c` | コンテンツID（page-idと排他） | | `-c "TK000002873"` |
| `-page-id` | `-p` | ページID（con-idと排他） | | `-p "page_id"` |
| `-markdown` | `-m` | マークダウンコンテンツ | * | `-m "# Hello World"` |
| `-toggle-h1` | | H1ヘッダートグル | | `-toggle-h1` |
| `-toggle-h2` | | H2ヘッダートグル | | `-toggle-h2` |
| `-toggle-h3` | | H3ヘッダートグル | | `-toggle-h3` |
| `-endpoint-url` | `-u` | APIエンドポイントURL | * | `-u "http://localhost:8080/task/patch"` |
| `-help` | `-h` | ヘルプを表示 | | `-h` |

## 使用例

### ページID指定での基本使用

```bash
# 基本的な使用方法
./bin/notion-sync -operation "patch" -token "your_token" -page-id "page_id" -markdown "# Hello World" -endpoint-url "http://localhost:8080/task/patch"

# 短縮オプションを使用
./bin/notion-sync -o "patch" -t "your_token" -p "page_id" -m "# Hello World" -u "http://localhost:8080/task/patch"

# 位置引数での指定
./bin/notion-sync "your_token" "# Hello World" "http://localhost:8080/task/patch" -operation "patch" -page-id "page_id"
```

### コンテンツID指定での使用

```bash
# コンテンツIDを使用
./bin/notion-sync -operation "patch" -token "your_token" -con-id "TK000002873" -markdown "# Test Content" -endpoint-url "http://localhost:8080/task/patch"

# 短縮オプションを使用
./bin/notion-sync -o "patch" -t "your_token" -c "TK000002873" -m "# Test Content" -u "http://localhost:8080/task/patch"
```

### ヘッダートグルオプション付き

```bash
# H1とH2ヘッダーをトグル
./bin/notion-sync -operation "patch" -token "your_token" -page-id "page_id" -markdown "# Hello World" -toggle-h1 -toggle-h2 -endpoint-url "http://localhost:8080/task/patch"

# 全てのヘッダーレベルをトグル
./bin/notion-sync -o "patch" -t "your_token" -p "page_id" -m "# Hello World" -toggle-h1 -toggle-h2 -toggle-h3 -u "http://localhost:8080/task/patch"
```

### 複雑なマークダウンコンテンツ

```bash
# 複数行のマークダウン
./bin/notion-sync -o "patch" -t "your_token" -p "page_id" -m "# タイトル

## サブタイトル

- リスト項目1
- リスト項目2

**太字テキスト**と*斜体テキスト*" -u "http://localhost:8080/task/patch"
```

## 出力フォーマット

### 成功時の出力

```
Status: 200

Headers:
Content-Type: application/json
Date: Fri, 02 Aug 2025 06:00:00 GMT

Body:
{
  "success": true,
  "message": "Page patched successfully"
}
```

### エラー時の出力

```
Status: 400

Headers:
Content-Type: application/json

Body:
{
  "success": false,
  "error": "con_id または page_id のいずれかを指定してください"
}
```

## エラーハンドリング

### 必須パラメータの不足

```bash
./bin/notion-sync -token "your_token"
```

```
エラー: マークダウンコンテンツが指定されていません
Notion同期CLIツール

使用方法:
  基本的な使用方法（ページID指定）:
    ./bin/notion-sync -token "your_token" -page-id "page_id" -markdown "# Hello World" -endpoint-url "http://localhost:8080/task/patch"
...
```

### 排他制御エラー

```bash
./bin/notion-sync -token "your_token" -con-id "con_id" -page-id "page_id" -markdown "# Test" -endpoint-url "http://localhost:8080/task/patch"
```

```
エラー: con_id と page_id の両方を指定することはできません
```

### HTTPリクエストエラー

```bash
./bin/notion-sync -token "your_token" -page-id "page_id" -markdown "# Test" -endpoint-url "http://localhost:9999/invalid"
```

```
エラー: パッチリクエスト送信に失敗しました: HTTPリクエスト送信エラー: Get "http://localhost:9999/invalid": dial tcp [::1]:9999: connect: connection refused
```

## HTTPリクエスト形式

このツールは以下のJSON形式でPOSTリクエストを送信します：

```json
{
  "token": "your_token",
  "con_id": "TK000002873",
  "page_id": "",
  "markdown_content": "# Hello World",
  "toggle_h1": false,
  "toggle_h2": false,
  "toggle_h3": false
}
```

## 技術仕様

### アーキテクチャ

- **Clean Architecture**: ドメイン、ユースケース、インフラストラクチャの分離
- **SOLID原則**: インターフェースを活用した疎結合な設計
- **依存性注入**: テスト可能性を考慮した設計

### ディレクトリ構造

```
cmd/cli/notion-sync/
├── main.go                    # CLIエントリーポイント
└── README.md                  # このファイル

internal/notion_sync/
├── config/
│   ├── config.go              # 設定構造体とバリデーション
│   └── flag_parser.go         # コマンドライン引数解析
└── usecases/
    └── services.go            # NotionSyncService実装
```

### 使用技術

- **Go**: プログラミング言語
- **標準ライブラリ**: `flag`, `fmt`, `os`, `encoding/json`など
- **内部依存**: `devbox/internal/http_request` HTTPリクエスト処理

## 開発者向け情報

### ビルド

```bash
# 開発用ビルド
go build -o bin/notion-sync ./cmd/cli/notion-sync

# リリース用ビルド（複数プラットフォーム）
GOOS=linux GOARCH=amd64 go build -o bin/notion-sync-linux ./cmd/cli/notion-sync
GOOS=windows GOARCH=amd64 go build -o bin/notion-sync.exe ./cmd/cli/notion-sync
GOOS=darwin GOARCH=amd64 go build -o bin/notion-sync-mac ./cmd/cli/notion-sync
```

### テスト

```bash
# 単体テスト
go test ./internal/notion_sync/...

# カバレッジ付きテスト
go test -coverprofile=coverage.out ./internal/notion_sync/...
go tool cover -html=coverage.out -o coverage.html
```

### デバッグ実行

```bash
# go runでの実行
go run ./cmd/cli/notion-sync/main.go -operation "patch" -token "test" -page-id "test" -markdown "# Test" -endpoint-url "http://localhost:8080/task/patch"

# バイナリでの実行
./bin/notion-sync -operation "patch" -token "test" -page-id "test" -markdown "# Test" -endpoint-url "http://localhost:8080/task/patch"
```

## 関連プロジェクト

- [notion-synchronizer](https://github.com/landmaster135/notion-synchronizer) - バックエンドAPIサーバー
- [devbox](https://github.com/landmaster135/devbox) - 開発ツール集

## ライセンス

このプロジェクトのライセンスについては、プロジェクトルートのLICENSEファイルを参照してください。
