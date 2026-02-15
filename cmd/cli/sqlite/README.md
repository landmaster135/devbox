# sqlite CLI

SQLite ファイルに対して操作を行う CLI です。

## 対応 operation

| operation | 説明 |
| --- | --- |
| `list-tables` | SQLite ファイル内のテーブル一覧を取得する |

## フラグ

| フラグ | 必須 | 説明 | デフォルト |
| --- | --- | --- | --- |
| `--operation`, `-o` | 必須 | 実行する操作（`list-tables`） | なし |
| `--opearation` | 任意 | `--operation` の誤記互換フラグ | なし |
| `--db-path` | 必須 | SQLite ファイルのパス | なし |
| `--format` | 任意 | 出力形式（`text`, `json`） | `text` |
| `--help`, `-h` | 任意 | ヘルプ表示 | `false` |

## 使用例

```bash
go run ./cmd/cli/sqlite --operation=list-tables --db-path=./sample.db
```

```bash
go run ./cmd/cli/sqlite --opearation=list-tables --db-path=./sample.db --format=json
```

## 出力例

成功時（`--format=text`）:

```text
users
posts
```

成功時（`--format=json`）:

```json
[
  "posts",
  "users"
]
```

エラー時:

```text
エラー: --db-path は必須です
```
