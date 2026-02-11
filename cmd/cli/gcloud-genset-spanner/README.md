# Gcloud Genset Spanner

Cloud Spanner の日常運用で使う `gcloud spanner` コマンドを安全に組み立て、コピーしやすい形で提示する CLI ツールです。インスタンス・データベースの準備で毎回調べがちなフラグをテンプレ化し、入力値の検証と整形を自動化します。

## 特徴

- **操作を `--operation` で明示**: `instance-list` や `db-create` など、必要なフラグが自動的に要求されます。
- **必須フラグの事前チェック**: `instance-create` では `--config` や `--nodes` の入力漏れを起動前に検出します。

## インストール

```bash
# プロジェクトルートで実行
go build -o bin/gcloud-genset-spanner ./cmd/cli/gcloud-genset-spanner
```

## オプション

| オプション | 説明 | 必須 | 対応オペレーション |
|-----------|------|------|---------------------|
| `-operation` | 実行する操作（`instance-list` / `instance-create` / `db-create` / `db-list` / `db-describe`） | ✔︎ | 全て |
| `-instance-id` | 対象の Spanner インスタンス ID | ✔︎ | `instance-create`, `db-*` |
| `-config` | インスタンス構成 ID（例: `regional-asia-northeast1`） | ✔︎ | `instance-create` |
| `-description` | インスタンスの説明 | ✔︎ | `instance-create` |
| `-nodes` | 作成するノード数 (1 以上) | ✔︎ | `instance-create` |
| `-db-id` | 対象のデータベース ID | ✔︎ | `db-create`, `db-describe` |
| `-ddl-file-path` | 適用する DDL ファイルへのパス | ✔︎ | `db-create` |
| `-help` | ヘルプを表示 |  | 全て |

## オペレーション一覧

| オペレーション | 説明 | 生成されるコマンド例 |
|----------------|------|-----------------------|
| `instance-list` | すべての Spanner インスタンスを一覧表示 | `gcloud spanner instances list` |
| `instance-create` | インスタンスを新規作成 | `gcloud spanner instances create ... --config=... --nodes=...` |
| `db-create` | インスタンス配下にデータベースを新規作成 | `gcloud spanner databases create ... --instance=... --ddl-file=...` |
| `db-list` | 指定インスタンスのデータベース一覧を表示 | `gcloud spanner databases list --instance=...` |
| `db-describe` | 指定データベースの詳細を表示 | `gcloud spanner databases describe ... --instance=...` |

## 使用例
インスタンス一覧
```bash
go run ./cmd/cli/gcloud-genset-spanner -operation instance-list
```

インスタンス作成
```bash
go run ./cmd/cli/gcloud-genset-spanner \
  -operation instance-create \
  -instance-id payments-prod \
  -config regional-asia-northeast1 \
  -description "Payments Production" \
  -nodes 2
```

データベース作成
```bash
go run ./cmd/cli/gcloud-genset-spanner \
  -operation db-create \
  -instance-id payments-prod \
  -db-id ledger \
  -ddl-file-path schema.sql
```

## テスト

```bash
go test ./internal/gcloud_genset_spanner/...
```
