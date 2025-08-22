# PostgreSQL CLIツール

PostgreSQLデータベースのテーブルダンプ機能を提供するCLIツールです。

## 機能

- テーブルの全レコードをダンプ
- 複数の出力フォーマット対応（JSON、CSV、SQL）
- レコード数制限機能
- 出力ディレクトリ指定

## ビルド方法

```bash
cd /home/nov/devbox/cmd/cli/postgresql
go build -o postgresql-cli .
```

## 使用例

### 基本的な使用例

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

### ヘルプ表示

```bash
go run ./cmd/cli/postgresql -help
```

## パラメータ

### 必須パラメータ

| パラメータ | 説明 | 例 |
|-----------|------|-----|
| `--operation` | 実行する操作（現在は"dump"のみ対応） | `dump` |
| `--database-url` | PostgreSQLデータベース接続URL | `postgres://user:pass@localhost/db` |
| `--table-name` | ダンプするテーブル名 | `users` |

### オプションパラメータ

| パラメータ | 説明 | デフォルト値 | 例 |
|-----------|------|-------------|-----|
| `--output-path` | 出力ディレクトリパス | カレントディレクトリ | `/tmp` |
| `--format` | 出力フォーマット（json, csv, sql） | `json` | `csv` |
| `--limit` | 最大レコード数 | 制限なし | `100` |
| `--help` | ヘルプを表示 | - | - |

## 出力形式

### JSON形式（デフォルト）

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

### 出力ファイル

ダンプされたデータは指定されたディレクトリに以下の命名規則でファイルが作成されます：

```
{テーブル名}_{実行日時}.{拡張子}
```

例：
- `users_20250822_180530.json`
- `products_20250822_180530.csv`
- `orders_20250822_180530.sql`

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
