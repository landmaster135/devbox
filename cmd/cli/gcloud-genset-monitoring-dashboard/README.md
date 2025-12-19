# Google Cloud Run モニタリングダッシュボード作成ツール

Google Cloud Run サービス用のモニタリングダッシュボードを自動作成するCLIツールです。

## 機能

- Cloud Run サービスの存在確認
- 包括的なモニタリングダッシュボードの作成
- サービスアカウントによる認証サポート
- Application Default Credentials (ADC) サポート

## ダッシュボード構成

作成されるダッシュボードには以下の16個のウィジェットが含まれます：

### 行1: リクエスト概要
- リクエスト数 (req/sec)
- ステータス別リクエスト (2xx/4xx/5xx)
- 累積リクエスト数 (24h total)
- リクエスト時間分布 (heatmap)

### 行2: パフォーマンス指標
- リクエストレイテンシ (P50,P95,P99)
- エラー率
- 最大同時リクエスト
- レスポンス時間 (Mean/Max)

### 行3: コンテナ指標
- インスタンス数
- 起動レイテンシ
- 課金時間

### 行4: リソース使用状況
- CPU使用率
- メモリ使用率
- メモリ使用量

### 行5: ネットワーク
- 送信バイト数
- 受信バイト数

## 使用方法

### 必須パラメータ

```bash
go run ./cmd/cli/gcloud-genset-monitoring-dashboard \
  -operation=create-dashboard-for-cloud-run \
  -project=<Google Cloud プロジェクトID> \
  -location=<Cloud Run サービスのロケーション> \
  -service=<Cloud Run サービス名>
```

### オプションパラメータ

```bash
go run ./cmd/cli/gcloud-genset-monitoring-dashboard \
  -operation=create-dashboard-for-cloud-run \
  -project=my-project \
  -location=us-central1 \
  -service=my-service \
  -service-account-id=monitoring-sa
```

### パラメータ説明

- `-operation`: 実行する操作 (現在は `create-dashboard-for-cloud-run` のみサポート)
- `-project`: Google Cloud プロジェクトID
- `-location`: Cloud Run サービスのロケーション (例: us-central1, asia-northeast1)
- `-service`: Cloud Run サービス名
- `-service-account-id`: サービスアカウントID (オプション)
- `-help`: ヘルプメッセージを表示

## 認証

### Application Default Credentials (ADC) を使用する場合

```bash
# 開発環境での認証設定
gcloud auth application-default login

# 本番環境では環境変数やCompute Engineのサービスアカウントが自動的に使用されます
```

### サービスアカウントを指定する場合

```bash
go run ./cmd/cli/gcloud-genset-monitoring-dashboard \
  -operation=create-dashboard-for-cloud-run \
  -project=my-project \
  -location=us-central1 \
  -service=my-service \
  -service-account-id=monitoring-service-account
```

この場合、現在の認証情報を使って指定されたサービスアカウントにimpersonateします。

## 必要な権限

### 実行者に必要な権限

- `run.services.get` (Cloud Run サービスの確認用)
- `monitoring.dashboards.create` (ダッシュボード作成用)

### サービスアカウントを指定する場合の追加権限

- `iam.serviceAccountTokenCreator` (指定されたサービスアカウントのImpersonationに必要)

## ビルド

```bash
cd devbox/cmd/cli/gcloud-genset-monitoring-dashboard
go build -o gcloud-genset-monitoring-dashboard .
```

## テスト

```bash
# 設定のテスト
cd devbox/internal/gcloud_monitoring/config
go test -v

# サービスのテスト
cd devbox/internal/gcloud_monitoring/usecases
go test -v

# 全体のテスト
cd devbox
go test ./internal/gcloud_monitoring/... -v
```

## 使用例

### 基本的な使用例
```bash
go run ./cmd/cli/gcloud-genset-monitoring-dashboard \
  -operation=create-dashboard-for-cloud-run \
  -project=my-gcp-project \
  -location=us-central1 \
  -service=my-web-app
```

### サービスアカウントを使用する例

```bash
go run ./cmd/cli/gcloud-genset-monitoring-dashboard \
  -operation=create-dashboard-for-cloud-run \
  -project=my-gcp-project \
  -location=asia-northeast1 \
  -service=my-api-service \
  -service-account-id=monitoring-dashboard-creator
```

### ヘルプの表示

```bash
go run ./cmd/cli/gcloud-genset-monitoring-dashboard -help
```

## エラーハンドリング

- Cloud Run サービスが見つからない場合、エラーメッセージが表示されます
- 認証に失敗した場合、適切なエラーメッセージが表示されます
- 必須パラメータが不足している場合、ヘルプメッセージが表示されます

## トラブルシューティング

### 認証エラーが発生する場合

1. ADCが設定されているか確認:
   ```bash
   gcloud auth application-default print-access-token
   ```

2. 必要な権限があるか確認:
   ```bash
   gcloud projects get-iam-policy PROJECT_ID
   ```

### Cloud Run サービスが見つからない場合

1. サービス名とロケーションが正しいか確認:
   ```bash
   gcloud run services list --project=PROJECT_ID
   ```

2. プロジェクトIDが正しいか確認:
   ```bash
   gcloud config get-value project
   ```

## 依存関係

- Go
- Google Cloud Go SDK
- Google Cloud Run API
- Google Cloud Monitoring API

## ライセンス

このプロジェクトのライセンスに従います。
