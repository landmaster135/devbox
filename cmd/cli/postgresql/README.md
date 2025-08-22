# PostgreSQL CLIツール

PostgreSQLデータベースのテーブルダンプ機能を提供するCLIツールです。

## 機能

- テーブルの全レコードをダンプ
- テーブル一覧の取得（最小限・詳細）
- 複数の出力フォーマット対応（JSON、CSV、SQL、テキスト）
- レコード数制限機能
- 出力ディレクトリ指定

## ビルド方法

```bash
cd /home/nov/devbox/cmd/cli/postgresql
go build -o postgresql-cli .
```

## 使用例

### テーブルダンプ

```bash
# 基本的なダンプ（JSON形式、カレントディレクトリに出力）
go run ./cmd/cli/postgresql --operation=dump --database-url="postgres://user:pass@localhost/db" --table-name=users

# CSV形式でダンプ
go run ./cmd/cli/postgresql --operation=dump --database-url="postgres://user:pass@localhost/db" --table-name=users --format=csv

# SQL形式でダンプ
go run ./cmd/cli/postgresql --operation=dump --database-url="postgres://user:pass@localhost/db" --table-name=users --format=sql

# 出力ディレクトリを指定
go run ./cmd/cli/postgresql --operation=dump --database-url="postgres://user:pass@localhost/db" --table-name=users --output-path=/tmp

# レコード数を制限
go run ./cmd/cli/postgresql --operation=dump --database-url="postgres://user:pass@localhost/db" --table-name=users --limit=100

# 全オプション指定
go run ./cmd/cli/postgresql --operation=dump --database-url="postgres://user:pass@localhost/db" --table-name=users --format=json --output-path=/tmp --limit=100
```

### テーブル一覧取得

```bash
# テーブル一覧（最小限）- テーブル名のみ
go run ./cmd/cli/postgresql --operation=list-tables-minimum --database-url="postgres://user:pass@localhost/db"

# テーブル一覧（詳細）- JSON形式
go run ./cmd/cli/postgresql --operation=list-tables --database-url="postgres://user:pass@localhost/db"

# テーブル一覧（詳細）- テキスト形式
go run ./cmd/cli/postgresql --operation=list-tables --database-url="postgres://user:pass@localhost/db" --format=text
```

### ヘルプ表示

```bash
go run ./cmd/cli/postgresql -help
```

## パラメータ

### 必須パラメータ

| パラメータ | 説明 | 例 |
|-----------|------|-----|
| `--operation` | 実行する操作 | `dump`, `list-tables-minimum`, `list-tables` |
| `--database-url` | PostgreSQLデータベース接続URL | `postgres://user:pass@localhost/db` |
| `--table-name` | ダンプするテーブル名（dump操作時のみ必須） | `users` |

### オプションパラメータ

| パラメータ | 説明 | デフォルト値 | 例 |
|-----------|------|-------------|-----|
| `--output-path` | 出力ディレクトリパス（dump操作時のみ） | カレントディレクトリ | `/tmp` |
| `--format` | 出力フォーマット | `json` | `csv`, `sql`, `text` |
| `--limit` | 最大レコード数（dump操作時のみ） | 制限なし | `100` |
| `--help` | ヘルプを表示 | - | - |

### 操作別パラメータ

dump操作
- **必須**: `--operation`, `--database-url`, `--table-name`
- **オプション**: `--output-path`, `--format` (json, csv, sql), `--limit`

list-tables-minimum操作
- **必須**: `--operation`, `--database-url`
- **オプション**: なし（formatは常にjson）

list-tables操作
- **必須**: `--operation`, `--database-url`
- **オプション**: `--format` (json, text)

## 出力形式

### dump操作の出力

**JSON形式（デフォルト）**

```json
{
  "table_name": "users",
  "record_count": 150,
  "output_path": "/tmp",
  "file_name": "users_20250822_180530.json",
  "format": "json",
  "executed_at": "2025-08-22 18:05:30"
}
```

**出力ファイル**

ダンプされたデータは指定されたディレクトリに以下の命名規則でファイルが作成されます：

```
{テーブル名}_{実行日時}.{拡張子}
```

例：
- `users_20250822_180530.json`
- `products_20250822_180530.csv`
- `orders_20250822_180530.sql`

### list-tables-minimum操作の出力

```json
[
  {"table_name": "users"},
  {"table_name": "products"},
  {"table_name": "orders"}
]
```

### list-tables操作の出力

**JSON形式**

```json
[
  {
    "table_name": "users",
    "table_comment": "ユーザー情報テーブル",
    "primary_keys": ["id"],
    "unique_keys": ["email"],
    "foreign_keys": []
  },
  {
    "table_name": "products",
    "table_comment": "商品情報テーブル",
    "primary_keys": ["product_id"],
    "unique_keys": ["product_code"],
    "foreign_keys": []
  }
]
```

**テキスト形式**

```
テーブル一覧:

テーブル名: users
コメント: ユーザー情報テーブル
主キー: id
一意キー: email
外部キー: なし

テーブル名: products
コメント: 商品情報テーブル
主キー: product_id
一意キー: product_code
外部キー: なし
```

## エラーハンドリング

- 必須パラメータが不足している場合、エラーメッセージとヘルプが表示されます
- データベース接続に失敗した場合、詳細なエラーメッセージが表示されます
- テーブルが存在しない場合、エラーメッセージが表示されます
- 出力ディレクトリの作成に失敗した場合、エラーメッセージが表示されます

## 注意事項

- データベース接続URLにはパスワードが含まれるため、コマンド履歴に注意してください
- 大量のデータをダンプする場合は、`--limit`オプションの使用を推奨します
- 出力ディレクトリが存在しない場合、自動的に作成されます
- 同名のファイルが存在する場合、上書きされます

## 依存関係

- Go 1.19以上
- PostgreSQLドライバー（github.com/lib/pq）

## テスト

```bash
# 設定のテスト
cd /home/nov/devbox/internal/postgresql/config
go test -v

# 全体のテスト
cd /home/nov/devbox
go test ./internal/postgresql/... -v
