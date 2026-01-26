# PostgreSQL CLIツール

PostgreSQLデータベースのテーブルダンプ機能を提供するCLIツールです。

## 機能

- 単一テーブルの全レコードをダンプ
- **データベース内の全テーブルを一括ダンプ**
- テーブル一覧の取得（最小限・詳細）
- 複数の出力フォーマット対応（JSON、CSV、SQL、テキスト）
- ダンプ結果サマリをJSONまたはMarkdownで取得
- ダンプ結果サマリの見出しをカスタマイズ
- レコード数制限機能
- 出力ディレクトリ指定
- タイムスタンプと出力ファイル名に使用するタイムゾーンを指定可能

## ビルド方法

```bash
cd /home/user/devbox/cmd/cli/postgresql
go build -o postgresql-cli .
```

## 使用例

### 単一テーブルダンプ

```bash
# 基本的なダンプ（JSON形式、カレントディレクトリに出力）
go run ./cmd/cli/postgresql --operation=dump --database-url="postgres://user:pass@localhost/db" --table-name=users

# CSV形式でダンプ
go run ./cmd/cli/postgresql --operation=dump --database-url="postgres://user:pass@localhost/db" --table-name=users --format=csv

# SQL形式でダンプ
go run ./cmd/cli/postgresql --operation=dump --database-url="postgres://user:pass@localhost/db" --table-name=users --format=sql

# 出力ディレクトリを指定
go run ./cmd/cli/postgresql --operation=dump --database-url="postgres://user:pass@localhost/db" --table-name=users --output-path=/tmp

# タイムゾーンを指定（Asia/Tokyo）
go run ./cmd/cli/postgresql --operation=dump --database-url="postgres://user:pass@localhost/db" --table-name=users --timezone=Asia/Tokyo

# レコード数を制限
go run ./cmd/cli/postgresql --operation=dump --database-url="postgres://user:pass@localhost/db" --table-name=users --limit=100

# 全オプション指定
go run ./cmd/cli/postgresql --operation=dump --database-url="postgres://user:pass@localhost/db" --table-name=users --format=json --output-path=/tmp --limit=100
```

### 全テーブル一括ダンプ

```bash
# 基本的な全テーブルダンプ（JSON形式、カレントディレクトリに出力）
go run ./cmd/cli/postgresql --operation=dump-all-tables --database-url="postgres://user:pass@localhost/db"

# CSV形式で全テーブルダンプ
go run ./cmd/cli/postgresql --operation=dump-all-tables --database-url="postgres://user:pass@localhost/db" --format=csv

# SQL形式で全テーブルダンプ
go run ./cmd/cli/postgresql --operation=dump-all-tables --database-url="postgres://user:pass@localhost/db" --format=sql

# 出力ディレクトリを指定して全テーブルダンプ
go run ./cmd/cli/postgresql --operation=dump-all-tables --database-url="postgres://user:pass@localhost/db" --output-path=/tmp/dumps

# 各テーブルのレコード数を制限して全テーブルダンプ
go run ./cmd/cli/postgresql --operation=dump-all-tables --database-url="postgres://user:pass@localhost/db" --limit=1000

# 並行処理数を指定して全テーブルダンプ
go run ./cmd/cli/postgresql --operation=dump-all-tables --database-url="postgres://user:pass@localhost/db" --concurrency=3

# タイムゾーンを指定して全テーブルダンプ
go run ./cmd/cli/postgresql --operation=dump-all-tables --database-url="postgres://user:pass@localhost/db" --timezone=UTC

# 全オプション指定で全テーブルダンプ
go run ./cmd/cli/postgresql --operation=dump-all-tables --database-url="postgres://user:pass@localhost/db" --format=csv --output-path=/tmp/dumps --limit=500 --concurrency=5

# ダンプ結果サマリをMarkdownで取得
go run ./cmd/cli/postgresql --operation=dump-all-tables --database-url="postgres://user:pass@localhost/db" --result-format=markdown

# Markdown出力の見出しを変更
go run ./cmd/cli/postgresql --operation=dump-all-tables --database-url="postgres://user:pass@localhost/db" --result-format=markdown --result-heading="Production Dump"
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
| `--operation` | 実行する操作 | `dump`, `dump-all-tables`, `list-tables-minimum`, `list-tables` |
| `--database-url` | PostgreSQLデータベース接続URL | `postgres://user:pass@localhost/db` |
| `--table-name` | ダンプするテーブル名（dump操作時のみ必須） | `users` |

### オプションパラメータ

| パラメータ | 説明 | デフォルト値 | 例 |
|-----------|------|-------------|-----|
| `--output-path` | 出力ディレクトリパス（dump操作時のみ） | カレントディレクトリ | `/tmp` |
| `--format` | 出力フォーマット | `json` | `csv`, `sql`, `text` |
| `--result-format` | ダンプ結果サマリのフォーマット（dump-all-tables向け） | `json` | `markdown` |
| `--result-heading` | ダンプ結果サマリの見出し（dump-all-tables + Markdown向け） | なし | `Production Dump` |
| `--timezone` | タイムスタンプ/ファイル名に使用するタイムゾーン | システムローカル | `Asia/Tokyo` |
| `--limit` | 最大レコード数（dump操作時のみ） | 制限なし | `100` |
| `--concurrency` | 並行処理数（dump-all-tables操作時のみ） | CPUコア数（最大10） | `3` |
| `--help` | ヘルプを表示 | - | - |

### 操作別パラメータ

dump操作
- **必須**: `--operation`, `--database-url`, `--table-name`
- **オプション**: `--output-path`, `--format` (json, csv, sql), `--limit`, `--timezone`

dump-all-tables操作
- **必須**: `--operation`, `--database-url`
- **オプション**: `--output-path`, `--format` (json, csv, sql), `--limit`, `--concurrency`, `--result-format` (json, markdown), `--result-heading`, `--timezone`

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

### dump-all-tables操作の出力

**JSON形式（デフォルト）**

```json
{
  "database_name": "mydb",
  "total_tables": 3,
  "results": [
    {
      "table_name": "users",
      "record_count": 150,
      "output_path": "/tmp/dumps",
      "file_name": "users_20250822_180530.json",
      "format": "json",
      "executed_at": "2025-08-22 18:05:30"
    },
    {
      "table_name": "products",
      "record_count": 75,
      "output_path": "/tmp/dumps",
      "file_name": "products_20250822_180531.json",
      "format": "json",
      "executed_at": "2025-08-22 18:05:31"
    },
    {
      "table_name": "orders",
      "record_count": 200,
      "output_path": "/tmp/dumps",
      "file_name": "orders_20250822_180532.json",
      "format": "json",
      "executed_at": "2025-08-22 18:05:32"
    }
  ],
  "failed_tables": [],
  "executed_at": "2025-08-22 18:05:30"
}
```

**Markdown形式（`--result-format=markdown`）**

```
## PostgreSQL Dump Report

| 項目 | 値 |
| --- | --- |
| Database | `mydb` |
| Total tables discovered | 3 |
| Successful dumps | 3 |
| Failed dumps | 0 |
| Executed at | 2025-08-22 18:05:30 |

### Successful Tables
| Table | Rows | File | Format |
| --- | --- | --- | --- |
| `users` | 150 | users_20250822_180530.json | json |
| `products` | 75 | products_20250822_180531.json | json |
| `orders` | 200 | orders_20250822_180532.json | json |

### Failed Tables
| Table | Error |
| --- | --- |
| なし | - |
```

**エラーが発生した場合の出力例**

```json
{
  "database_name": "mydb",
  "total_tables": 3,
  "results": [
    {
      "table_name": "users",
      "record_count": 150,
      "output_path": "/tmp/dumps",
      "file_name": "users_20250822_180530.json",
      "format": "json",
      "executed_at": "2025-08-22 18:05:30"
    }
  ],
  "failed_tables": [
    {
      "table_name": "products",
      "error": "permission denied for table products"
    },
    {
      "table_name": "orders",
      "error": "disk full"
    }
  ],
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

`--timezone`を指定した場合、これらのファイル名や結果中の`executed_at`は指定タイムゾーンを基準とした日時になります。

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
cd /home/user/devbox/internal/postgresql/config
go test -v

# 全体のテスト
cd /home/user/devbox
go test ./internal/postgresql/... -v
