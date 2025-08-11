# Google Cloud Workload Identity Federation CLIツール

Google Cloud Workload Identity FederationとGitHub Actions認証の設定を自動化するCLIツールです。

## 機能

- Google Cloud Workload Identity Poolの作成
- GitHub Actions用OIDCプロバイダーの作成
- サービスアカウントの作成と権限設定
- IAMポリシーバインディングの自動設定
- Workload Identityバインディングの設定
- GitHub Actions用ワークフロー設定の生成
- Discord Webhookによる設定完了通知

## 前提条件

- Google Cloud CLIがインストールされ、認証済みであること
- 対象のGoogle Cloudプロジェクトへの適切な権限があること
- Discord Webhook URLが準備されていること

## インストール

```bash
cd /home/nov/devbox
go build -o bin/gcloud-wrapper-workload-identity-federation ./cmd/cli/gcloud-wrapper-workload-identity-federation
```

## 使用例

### 基本的な使用方法

```bash
go run ./cmd/cli/gcloud-wrapper-workload-identity-federation \
  -project-id my-project \
  -pool-id github-pool \
  -provider-id github-provider \
  -service-account-id monitoring-sa \
  -repo-owner myorg \
  -repo-name myrepo \
  -webhook-url "https://discord.com/api/webhooks/YOUR_WEBHOOK_ID/YOUR_WEBHOOK_TOKEN"
```

### 全オプション指定

```bash
go run ./cmd/cli/gcloud-wrapper-workload-identity-federation \
  -project-id my-project \
  -pool-id github-pool \
  -provider-id github-provider \
  -service-account-id monitoring-sa \
  -location global \
  -pool-description "GitHub Actions用のWorkload Identity Pool" \
  -repo-owner myorg \
  -repo-name myrepo \
  -webhook-url "https://discord.com/api/webhooks/YOUR_WEBHOOK_ID/YOUR_WEBHOOK_TOKEN"
```

## オプション

### 必須オプション

| オプション | 説明 |
|-----------|------|
| `-project-id` | Google CloudプロジェクトID |
| `-pool-id` | Workload Identity Pool ID |
| `-provider-id` | OIDC Provider ID |
| `-service-account-id` | サービスアカウントID |
| `-repo-owner` | GitHubリポジトリのオーナー |
| `-repo-name` | GitHubリポジトリ名 |
| `-webhook-url` | Discord WebhookのURL |

### 任意オプション

| オプション | デフォルト値 | 説明 |
|-----------|-------------|------|
| `-location` | `global` | リソースのロケーション |
| `-pool-description` | (空) | Workload Identity Poolの説明 |
| `-help` | - | ヘルプを表示 |

## 作成されるリソース

このツールは以下のGoogle Cloudリソースを作成します：

### 1. Workload Identity Pool
- **名前**: 指定したpool-id
- **場所**: 指定したlocation（デフォルト: global）
- **説明**: 指定したpool-description（任意）

### 2. OIDC Provider (GitHub Actions用)
- **名前**: 指定したprovider-id
- **発行者URI**: `https://token.actions.githubusercontent.com/`
- **対象リポジトリ**: 指定したrepo-owner/repo-name

### 3. サービスアカウント
- **ID**: 指定したservice-account-id
- **メールアドレス**: `{service-account-id}@{project-id}.iam.gserviceaccount.com`

### 4. IAMポリシーバインディング
以下のロールがサービスアカウントに付与されます：
- `roles/monitoring.editor` - Cloud Monitoringの編集権限
- `roles/run.viewer` - Cloud Runの閲覧権限
- `roles/iam.serviceAccounts.getAccessToken` - サービスアカウントトークン取得権限

### 5. Workload Identityバインディング
- **プリンシパル**: `principalSet://iam.googleapis.com/projects/{PROJECT_NUMBER}/locations/{location}/workloadIdentityPools/{pool-id}/attribute.repository/{repo-owner}/{repo-name}`
- **ロール**: `roles/iam.workloadIdentityUser`

## 実行される処理

ツールは以下の順序でBashスクリプトを生成し、Discord通知を送信します：

1. Workload Identity Poolの作成
2. GitHub Actions用OIDCプロバイダーの作成
3. サービスアカウントの作成
4. IAMポリシーバインディングの追加（monitoring.editor）
5. IAMポリシーバインディングの追加（run.viewer）
6. IAMポリシーバインディングの追加（iam.serviceAccounts.getAccessToken）
7. プロジェクト番号の取得
8. Workload Identityバインディングの追加

## Discord通知内容

ツール実行後、Discordに以下の情報が通知されます：

1. **セットアップ完了通知**
   - プロジェクト情報
   - リポジトリ情報

2. **詳細な設定手順**
   - 作成されるリソース一覧
   - 実行用Bashスクリプト
   - GitHub Actions用ワークフロー設定
   - GitHub Secretsに設定する値

## GitHub Actions設定

### ワークフロー設定例

```yaml
on:
  push:

permissions:
  contents: write # リポジトリへの書き込み権限 (バッジ更新のため)
  id-token: write # OIDCトークンをリクエストする権限 (Google Cloud認証のため)

env:
  GOOGLE_CLOUD_PROJECT_ID_01: 'test-pj'
  GCLOUD_PROJECT_NUMBER: ${{ secrets.GCLOUD_PROJECT_NUMBER }}
  GCLOUD_POOL_ID: ${{ secrets.GCLOUD_POOL_ID }}
  GCLOUD_PROVIDER_ID: ${{ secrets.GCLOUD_PROVIDER_ID }}
  GCLOUD_SERVICE_ACCOUNT_EMAIL: ${{ secrets.GCLOUD_SERVICE_ACCOUNT_EMAIL }}

jobs:
  test:
    runs-on: ubuntu-latest
    if: contains(github.event.head_commit.message, '[skip ci]') == false
    steps:
      - name: Checkout
        uses: actions/checkout@v4
      - name: Set up Golang
        uses: actions/setup-go@v5
        with:
          go-version-file: 'go.mod'
      - id: 'gcloud_auth'
        name: 'Authenticate to Google Cloud'
        uses: 'google-github-actions/auth@v2'
        with:
          create_credentials_file: true
          workload_identity_provider: 'projects/${{ env.GCLOUD_PROJECT_NUMBER }}/locations/global/workloadIdentityPools/${{ env.GCLOUD_POOL_ID }}/providers/${{ env.GCLOUD_PROVIDER_ID }}'
          service_account: '${{ env.GCLOUD_SERVICE_ACCOUNT_EMAIL }}'
      - name: 'Set up Cloud SDK'
        uses: 'google-github-actions/setup-gcloud@v2'
      - name: 'Use gcloud CLI'
        run: 'gcloud info'
```

### GitHub Secretsの設定

リポジトリの Settings > Secrets and variables > Actions で以下のSecretsを設定してください：

| Secret名 | 値 | 説明 |
|---------|----|------------|
| `GCLOUD_PROJECT_NUMBER` | スクリプト実行後に表示される値 | Google Cloudプロジェクト番号 |
| `GCLOUD_POOL_ID` | 指定したpool-id | Workload Identity Pool ID |
| `GCLOUD_PROVIDER_ID` | 指定したprovider-id | OIDC Provider ID |
| `GCLOUD_SERVICE_ACCOUNT_EMAIL` | `{service-account-id}@{project-id}.iam.gserviceaccount.com` | サービスアカウントのメールアドレス |

## セットアップ手順

1. **ツールの実行**
   ```bash
   ./bin/gcloud-wrapper-workload-identity-federation [オプション]
   ```

2. **Discordで生成されたBashスクリプトを確認**

3. **Bashスクリプトの実行**
   - Discordに送信されたスクリプトをコピー
   - ローカル環境で実行
   - 必要なBash関数が事前に定義されていることを確認

4. **GitHub Secretsの設定**
   - スクリプト実行後に表示されるプロジェクト番号を記録
   - GitHub リポジトリのSecretsに必要な値を設定

5. **GitHub Actionsワークフローの設定**
   - Discordに送信されたワークフロー設定をコピー
   - `.github/workflows/` ディレクトリに配置

## エラーハンドリング

- 必須パラメータが不足している場合、エラーメッセージとヘルプが表示されます
- 設定の検証に失敗した場合、具体的なエラーメッセージが表示されます
- Discord通知の送信に失敗した場合、エラーメッセージが表示されます

## 注意事項

- このツールはBashスクリプトを生成するのみで、実際のGoogle Cloudリソースの作成は行いません
- 生成されたスクリプトを実行する前に、内容を確認してください
- 必要なBash関数（`create_workload_identity_pool`等）が事前に定義されている必要があります
- Google Cloud CLIの認証と適切な権限が必要です

## トラブルシューティング

### よくある問題

1. **Bash関数が見つからない**
   - `/home/nov/dotfiles/iac/gcloud/iam.sh` をsourceしてください

2. **権限エラー**
   - Google Cloudプロジェクトに対する適切な権限があることを確認してください

3. **Discord通知が送信されない**
   - Webhook URLが正しいことを確認してください
   - ネットワーク接続を確認してください

## ライセンス

このプロジェクトのライセンスについては、プロジェクトルートのライセンスファイルを参照してください。

## 関連ドキュメント

- [Google Cloud Workload Identity Federation](https://cloud.google.com/iam/docs/workload-identity-federation)
- [GitHub Actions OIDC](https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/about-security-hardening-with-openid-connect)
- [Discord Webhooks](https://discord.com/developers/docs/resources/webhook)
