# Gcloud Genset Billing CLI

Google Cloud Billing の運用で利用する `gcloud billing` コマンドを安全に組み立てる CLI ツールです。請求アカウントや予算、プロジェクトに関する情報取得コマンドを素早く生成できます。

## サポートする操作

| operation          | 説明 |
|--------------------|------|
| `list-budgets`     | 請求アカウント配下の予算一覧を表示するコマンドを生成します。未指定時は最初の請求アカウントを自動取得します。|
| `list-projects`    | 請求アカウントに紐付くプロジェクト一覧を表示するコマンドを生成します。|
| `describe-project` | 指定した、もしくは現在設定されているプロジェクトの請求状態を確認するコマンドを生成します。|
| `describe-budget`  | 特定の予算詳細を確認するコマンドを生成します。未指定時は請求アカウントを自動取得します。|

## 使用例

### 1. 予算一覧を取得するコマンドを生成

```bash
go run ./cmd/cli/gcloud-genset-billing \
  -operation list-budgets \
  -limit 20 \
  -billing-account 0000-AAAA-BBBB
```

Output:
```bash
gcloud billing budgets list --billing-account='0000-AAAA-BBBB' --limit=20
```

### 2. 請求プロジェクトをフィルター付きで取得

```bash
go run ./cmd/cli/gcloud-genset-billing \
  -operation list-projects \
  -limit 5 \
  -filter "project_id ~ ^test-"
```

Output (請求アカウントは自動取得):
```bash
billing_account=$(gcloud billing accounts list --format='value(name)' 2>/dev/null | head -n 1); if [ -z "$billing_account" ]; then echo "請求アカウントが見つかりません" >&2; exit 1; fi; gcloud billing projects list --billing-account="$billing_account" --limit=5 --filter='project_id ~ ^test-'
```

### 3. プロジェクト請求情報を確認

```bash
go run ./cmd/cli/gcloud-genset-billing \
  -operation describe-project \
  -project-id MY_PROJECT_123
```

Output:
```bash
gcloud billing projects describe 'MY_PROJECT_123'
```

### 4. 特定予算の詳細を確認

```bash
go run ./cmd/cli/gcloud-genset-billing \
  -operation describe-budget \
  -budget-id 00AA00-123456-FFFF \
  -billing-account 0000-AAAA-BBBB
```

Output:
```bash
gcloud billing budgets describe '00AA00-123456-FFFF' --billing-account='0000-AAAA-BBBB'
```

## CLI オプション

| operation            | 主なオプション |
|----------------------|----------------|
| `list-budgets`       | `-limit` (デフォルト: 10), `-billing-account` |
| `list-projects`      | `-limit` (デフォルト: 10), `-filter`, `-billing-account` |
| `describe-project`   | `-project-id` (未指定時は `gcloud config get-value project` を利用) |
| `describe-budget`    | `-budget-id` (必須), `-billing-account` |

生成される文字列は **自動実行されません**。内容を確認し、必要に応じてコピーして利用してください。請求アカウントを指定しない場合は、`gcloud billing accounts list` の最初の結果を用いてコマンドが組み立てられます。

## ビルド

```bash
go build -o bin/gcloud-genset-billing ./cmd/cli/gcloud-genset-billing
```
