# Gcloud Genset CloudSQL

Cloud SQL インスタンスの運用で利用する `gcloud sql` コマンドを安全に組み立てて提示する CLI ツールです。`dotfiles/iac/gcloud/db.sh` に定義されている運用フローを Go で再現し、入力値の検証やコマンド整形をサポートします。

## 概要

- **削除フローの自動生成**: 削除前に起動・削除保護解除を順番に実行する複合コマンドを出力
- **起動ポリシーや削除保護の切り替え**: `patch` 操作用コマンドを迷わず取得
- **コピーしやすい出力**: 生成されたコマンドを枠で囲って表示し、そのままコピー可能

## インストール

```bash
# プロジェクトルートで実行
go build -o bin/gcloud-genset-cloudsql ./cmd/cli/gcloud-genset-cloudsql
```

## 使用方法

```bash
go run ./cmd/cli/gcloud-genset-cloudsql \
  -operation delete-instance \
  -instance-name my-sql-instance
```

標準出力に整形されたコマンドが表示されます。

## オプション

| オプション | 説明 | 必須 | 対応オペレーション | 例 |
|------------|------|------|---------------------|----|
| `-operation` | 実行する操作（`delete-instance` / `patch-deletion-protection` / `patch-activation-policy` / `start-instance` / `stop-instance`） | * | 全て | `-operation delete-instance` |
| `-instance-name` | 対象の Cloud SQL インスタンス名 | * | 全て | `-instance-name my-sql-instance` |
| `-deletion-protection-mode` | 削除保護の設定（`enable` or `disable`） | * | `patch-deletion-protection` | `-deletion-protection-mode disable` |
| `-activation-policy` | 起動ポリシー（`always` or `never`） | * | `patch-activation-policy` | `-activation-policy always` |
| `-help` | ヘルプを表示 |  | 全て | `-help` |

## オペレーション一覧
`-operation`に指定出来る値の一覧は下記になります。

| オペレーション | 説明 | 必須パラメータ |
|----------------|------|----------------|
| `delete-instance` | インスタンスを削除する複合コマンドを生成（起動 → 削除保護解除 → 削除） | `-operation`, `-instance-name` |
| `patch-deletion-protection` | 削除保護の有効化/無効化コマンドを生成 | `-operation`, `-instance-name`, `-deletion-protection-mode` |
| `patch-activation-policy` | 起動ポリシーを変更するコマンドを生成 | `-operation`, `-instance-name`, `-activation-policy` |
| `start-instance` | 起動ポリシーを `ALWAYS` に変更するコマンドを生成 | `-operation`, `-instance-name` |
| `stop-instance` | 起動ポリシーを `never` に変更するコマンドを生成 | `-operation`, `-instance-name` |

## 使用例

```bash
go run ./cmd/cli/gcloud-genset-cloudsql \
  -operation patch-deletion-protection \
  -instance-name my-sql-instance \
  -deletion-protection-mode disable
```

```bash
go run ./cmd/cli/gcloud-genset-cloudsql \
  -operation patch-activation-policy \
  -instance-name my-sql-instance \
  -activation-policy always
```

```bash
go run ./cmd/cli/gcloud-genset-cloudsql \
  -operation delete-instance \
  -instance-name my-sql-instance
```

## テスト

```bash
go test ./internal/gcloud_genset_cloudsql/...
```
