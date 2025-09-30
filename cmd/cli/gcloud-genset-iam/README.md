# Gcloud Genset IAM

Google Cloud IAM / Workload Identity Federation の `gcloud` コマンドを素早く組み立てる CLI ツールです。サービスアカウントの操作から Workload Identity Pool の構築・破棄まで、ツール内で管理しているスクリプトを Go 実装として再現し、コピペ可能なコマンド列を出力します。

## 概要

- **IAM バインド生成**: サービスアカウント/プロジェクト向けの `add-iam-policy-binding` コマンドを安全に生成
- **サービスアカウント管理**: create / list / enable / disable / delete / undelete / describe / update の各操作をサポート
- **Workload Identity Pool 管理**: プールの作成・更新・削除・復元・一覧・詳細取得を単一 CLI で呼び出し
- **Workload Identity Provider 管理**: OIDC プロバイダーの作成・更新・クリーンアップを支援 (GitHub Actions 用のエイリアスも提供)
- **Federation セットアップ/クリーンアップ**: リポジトリ連携一式を自動化するシェルスクリプトを生成
- **出力のみ**: CLI 自体は `gcloud` を実行せず、生成結果を標準出力に表示する安全設計

## インストール

```bash
# プロジェクトルートで実行
go build -o bin/gcloud-genset-iam ./cmd/cli/gcloud-genset-iam
```

## 使い方 (基本形)

```bash
go run ./cmd/cli/gcloud-genset-iam \
  -operation add-iam-policy-binding-to-project \
  -project-id MY_PROJECT \
  -service-account-id SERVICE_ACCOUNT_ID \
  -role roles/run.viewer
```

生成された `gcloud` コマンドが標準出力に装飾付きで表示されます。必要に応じてクリップボードへコピーし、そのまま実行できます。

## 主なオプション

| オプション | 説明 | 主な利用シナリオ |
|------------|------|------------------|
| `-operation` | 実行する操作 (必須) | 全操作 |
| `-project-id` / `-project-number` | プロジェクト ID / プロジェクト番号 | IAM バインド / プール / プロバイダー系 |
| `-service-account-id` / `-service-account-email` | サービスアカウント ID / メール | サービスアカウント操作全般 |
| `-role` / `-member` | IAM ロール / メンバー指定 | `add-iam-policy-binding-*` |
| `-pool-id` / `-provider-id` | Workload Identity Pool / Provider ID | プール・プロバイダー操作 |
| `-location` | ロケーション (デフォルト `global`) | プール・プロバイダー操作 |
| `-repository-owner` / `-repository-name` | GitHub リポジトリ情報 | Workload Identity Federation 関連 |
| `-condition` / `-condition-from-file` | IAM 条件をインライン / ファイルで指定 (排他) | `add-iam-policy-binding-to-service-account` 等 |
| `-description` / `-display-name` / `-disabled` | リソース更新項目 | `update-service-account`, `update-workload-identity-pool`, `update-oidc-workload-identity-pool-provider` |
| `-allowed-audiences` / `-attribute-mapping` / `-attribute-condition` | OIDC プロバイダー構成 | プロバイダーの作成・更新 |
| `-skip-confirmation` | クリーンアップ時の確認プロンプトを省略 | `cleanup-workload-identity-federation` |

詳しいシグネチャは `cfg.PrintUsage()` または本 README を参照してください。

## オペレーション一覧

### IAM バインド / サービスアカウント
- `add-iam-policy-binding-to-project`
- `add-iam-policy-binding-to-service-account`
- `add-workload-identity-binding-to-service-account`
- `create-service-account`
- `list-service-accounts`
- `enable-service-account` / `disable-service-account`
- `delete-service-account` / `undelete-service-account`
- `update-service-account`
- `describe-service-account`

### Workload Identity Pool
- `create-workload-identity-pool`
- `list-workload-identity-pools`
- `describe-workload-identity-pool`
- `update-workload-identity-pool`
- `delete-workload-identity-pool`
- `undelete-workload-identity-pool`

### Workload Identity Provider
- `create-oidc-workload-identity-pool-provider`
- `create-oidc-workload-identity-pool-provider-for-github-actions`
- `list-workload-identity-pool-providers`
- `describe-workload-identity-pool-provider`
- `update-oidc-workload-identity-pool-provider`
- `delete-workload-identity-pool-provider`
- `undelete-workload-identity-pool-provider`

### Federation シナリオ
- `setup-workload-identity-federation`
- `cleanup-workload-identity-federation`

## 使用例

### 1. サービスアカウントを作成しロール付与
```bash
$ go run ./cmd/cli/gcloud-genset-iam \
  -operation create-service-account \
  -service-account-id SERVICE_ACCOUNT_ID \
  -project-id MY_PROJECT \
  -role roles/run.admin
```

### 2. Workload Identity Pool を作成して GitHub Actions 用プロバイダーを登録
```bash
$ go run ./cmd/cli/gcloud-genset-iam \
  -operation setup-workload-identity-federation \
  -project-id MY_PROJECT \
  -pool-id github-pool \
  -provider-id github-provider \
  -service-account-id gha \
  -repository-owner landmaster135 \
  -repository-name devbox
```
出力されるスクリプトをシェルに貼り付けるだけで設定が完了します。

### 3. Federation リソースのクリーンアップ
```bash
$ go run ./cmd/cli/gcloud-genset-iam \
  -operation cleanup-workload-identity-federation \
  -project-id MY_PROJECT \
  -pool-id POOL_ID \
  -provider-id PROVIDER_ID \
  -service-account-id gha
```
`-skip-confirmation` を指定すると確認プロンプトを省略できます。

## テスト

```bash
go test ./internal/gcloud_genset_iam/...
```

## ディレクトリ構成

```
internal/gcloud_genset_iam/
├── config/     # CLI フラグ解析と検証
│   ├── config.go
│   ├── config_test.go
│   └── flag_parser.go
└── usecases/   # gcloud コマンド生成ロジック
    ├── services.go
    └── services_test.go
```
