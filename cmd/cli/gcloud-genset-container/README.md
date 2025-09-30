# Gcloud Genset Container

Cloud Run コンテナ／Cloud Functions (Gen2)／Pub/Sub の日常運用コマンドを生成するgcloud コマンドジェネレーターです。

## 概要

- **Cloud Run (コンテナ)**: デプロイ、環境変数更新、サービス削除などを一発で組み立て
- **Cloud Functions (Gen2)**: HTTP/PubSub トリガーのデプロイや環境変数更新コマンドを生成
- **Pub/Sub**: トピック/サブスクリプション作成・一覧・削除をサポート
- **出力強調表示**: 生成されたコマンドを枠付きで表示し、そのままコピーできます

## インストール

```bash
# プロジェクトルートで実行
go build -o bin/gcloud-genset-container ./cmd/cli/gcloud-genset-container
```

## 使用方法

```bash
go run ./cmd/cli/gcloud-genset-container \
  -operation deploy-cloud-run-container \
  -service-name MY_SERVICE \
  -project-id MY_PROJECT_ID \
  -region us-central1 \
  -timeout 40m
```

標準出力に生成された `gcloud run deploy ...` コマンドが表示されます。

## 主なオプション

| オプション | 説明 | 必須となる主なオペレーション例 |
|------------|------|------------------------------|
| `-operation` | 実行する操作を指定 | 全て (必須) |
| `-service-name` | Cloud Run サービス名 (関数名と共通) | deploy/update/delete 系 |
| `-function-name` | Cloud Functions (Gen2) の関数名 | deploy-cloud-run-function 系 |
| `-project-id` | 対象プロジェクト ID | Cloud Run コンテナ操作、PubSub 作成系など |
| `-region` | リージョン (`us-central1` など) | 多くの Cloud Run/Functions 操作 |
| `-env-file` | 環境変数ファイル (デフォルト `env.yml`) | update-cloud-run-*-env (ファイル) |
| `-env-vars` | `KEY=VALUE,...` 形式の環境変数 | update-cloud-run-function-env |
| `-topic-name` | Pub/Sub トピック名 | Pub/Sub 作成/一覧 |
| `-subscription-name` | Pub/Sub サブスクリプション名 | Pub/Sub 作成/一覧 |
| `-subscription-names` | 削除対象サブスクリプション (カンマ区切り) | delete-cloud-pubsub-subscriptions* |
| `-topic-names` | 削除対象トピック (カンマ区切り) | delete-cloud-pubsub-topics* |
| `-push-service-account` | Push サブスクリプションのサービスアカウント | create-cloud-pubsub-subscription |
| `-trigger-service-account` | Functions Pub/Sub トリガー用サービスアカウント | deploy-cloud-run-function-triggered-by-pubsub |
| `-trigger-topic` | Functions Pub/Sub トリガートピック ID | deploy-cloud-run-function-triggered-by-pubsub |

詳細なフラグ一覧は `go run ./cmd/cli/gcloud-genset-container -help` で確認できます。

## オペレーション一覧

| Operation | 説明 | 主な必須パラメータ |
|-----------|------|--------------------|
| `deploy-cloud-run-container` | Cloud Run コンテナのデプロイ (`gcloud run deploy --source .`) | `-service-name`, `-project-id` |
| `update-cloud-run-container-env` | コンテナの環境変数ファイル更新 | `-service-name`, `-project-id`, `-region` |
| `deploy-cloud-run-function` | HTTP トリガー Cloud Functions (Gen2) のデプロイ | `-function-name`, `-region`, `-entry-point` |
| `deploy-cloud-run-function-triggered-by-pubsub` | Pub/Sub トリガー関数のデプロイ | `-function-name`, `-project-id`, `-trigger-service-account`, `-trigger-topic` |
| `update-cloud-run-function-env` | Cloud Functions(Gen2) の env 更新 (`--update-env-vars`) | `-service-name`, `-region`, `-env-vars` |
| `update-cloud-run-service-env` | Cloud Run サービスの env ファイル更新 | `-service-name`, `-project-id`, `-region` |
| `create-cloud-pubsub-topic` | トピック作成 | `-topic-name` |
| `list-cloud-pubsub-topics` | トピック一覧 (名前フィルタ付き) | `-topic-name` |
| `list-cloud-pubsub-subscriptions` | サブスクリプション一覧 (`--uri` オプション対応) | (任意) `-subscription-name` |
| `create-cloud-pubsub-subscription` | Push サブスクリプション作成 | `-subscription-name`, `-topic-name`, `-topic-project`, `-push-service-account` |
| `delete-cloud-pubsub-subscriptions-and-topics` | サブスク/トピックの一括削除 | `-subscription-names` または `-topic-names` |
| `delete-cloud-pubsub-subscriptions` | サブスクリプション削除 | `-subscription-names` |
| `delete-cloud-pubsub-topics` | トピック削除 | `-topic-names` |
| `delete-cloud-run-function` | Cloud Functions(Gen2) が利用する Cloud Run サービス削除 | `-service-name`, `-region` |

## 使用例

### Cloud Run コンテナのデプロイコマンドを生成

```bash
go run ./cmd/cli/gcloud-genset-container \
  -operation deploy-cloud-run-container \
  -service-name MY_SERVICE \
  -project-id MY_PROJECT_ID \
  -region asia-northeast1 \
  -timeout 45m \
  -run-service-account my-svc@my-project.iam.gserviceaccount.com \
  -allow-unauthenticated=false
```

出力例:
```
==============================
生成された gcloud コマンド
==============================
gcloud run deploy 'my-svc' --source . --project='my-project' --region='asia-northeast1' --timeout='45m' --service-account='my-svc@my-project.iam.gserviceaccount.com' --no-allow-unauthenticated
==============================
```

### Pub/Sub Push サブスクリプション生成コマンド

```bash
go run ./cmd/cli/gcloud-genset-container \
  -operation create-cloud-pubsub-subscription \
  -subscription-name my-sub \
  -topic-name my-topic \
  -topic-project MY_PROJECT_ID \
  -push-service-account push@my-project.iam.gserviceaccount.com
```

## テスト

```bash
go test ./internal/gcloud_genset_container/...
```

## ディレクトリ構成

```
internal/gcloud_genset_container/
├── config/       # CLI フラグ解析とバリデーション
│   ├── config.go
│   └── flag_parser.go
└── usecases/     # gcloud コマンド生成ロジック
    ├── services.go
    └── services_test.go
```
