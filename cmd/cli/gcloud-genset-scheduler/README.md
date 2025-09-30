# gcloud-genset-scheduler

Google Cloud Scheduler 向けの gcloud コマンドを生成する CLI ツールです。`operation` を指定すると、必要なオプションを組み立てたコマンドを表示します。Cloud Run / Cloud Functions / Cloud SQL 向けのジョブ生成から、ジョブの更新・停止・削除まで幅広い操作をカバーします。

## 主な機能

- **Pub/Sub ジョブ生成**: Cloud Functions やワークフローをトリガーするジョブを作成
- **HTTP ジョブ生成**: Cloud Run コンテナや任意の HTTP エンドポイントを呼び出すジョブを作成
- **Cloud SQL 自動化**: インスタンス起動・停止用の Pub/Sub ジョブを自動生成し、メッセージ本文を JSON で作成
- **ジョブ管理**: 一覧取得、スケジュールやメッセージ本文の更新、一時停止 / 再開 / 削除 / 強制実行

## ビルド

```bash
# プロジェクトルートから
./scripts/build_gcloud_genset_scheduler.sh
```

## 基本的な使い方

```bash
go run ./cmd/cli/gcloud-genset-scheduler -operation create-pubsub-job \
  -job-name exec-cloud-run \
  -project-id SAMPLE_PROJECT \
  -pubsub-topic projects/sample/topics/trigger
```

出力例:
```
==============================
生成された gcloud コマンド
==============================
gcloud scheduler jobs create pubsub 'exec-cloud-run' \
  --schedule='0 4 * * 0-6' \
  --description='Trigger Cloud Functions to start Cloud SQL instance.' \
  --project='SAMPLE_PROJECT' \
  --location='us-central1' \
  --time-zone='Asia/Tokyo' \
  --topic='projects/sample/topics/trigger'
==============================
```

## サポートされている operation

| operation | 説明 | 主な必須フラグ |
|-----------|------|----------------|
| `create-pubsub-job` | Cloud Functions などを Pub/Sub で起動するジョブを作成 | `-job-name`, `-project-id`, `-pubsub-topic` |
| `create-http-job` | Cloud Run コンテナなどを HTTP で呼び出すジョブを作成 | `-job-name`, `-project-id`, `-http-method`, `-service-url` |
| `create-cloud-sql-job` | Cloud SQL 操作用 Pub/Sub ジョブを作成 (メッセージに JSON を組み込み) | `-job-name`, `-project-id`, `-pubsub-topic`, `-db-instance-id` |
| `create-start-cloud-sql-job` | Cloud SQL 起動ジョブを生成 (ジョブ名とスケジュールを自動補完) | `-project-id`, `-pubsub-topic`, `-db-instance-id` |
| `create-stop-cloud-sql-job` | Cloud SQL 停止ジョブを生成 (ジョブ名とスケジュールを自動補完) | `-project-id`, `-pubsub-topic`, `-db-instance-id` |
| `list-jobs` | ジョブ一覧を表示 | *任意* `-location`, `-limit` |
| `update-http-job` | HTTP ジョブのスケジュールや本文を更新 | `-job-name` |
| `update-pubsub-job` | Pub/Sub ジョブのスケジュールや本文を更新 | `-job-name` |
| `pause-job` / `resume-job` / `delete-job` / `run-job` | 各種ジョブ操作 | `-job-name`, `-location` |

## 使用例

### Cloud Run (HTTP) ジョブを作成
```bash
go run ./cmd/cli/gcloud-genset-scheduler -operation create-http-job \
  -job-name nightly-http-call \
  -project-id SAMPLE_PROJECT \
  -http-method POST \
  -service-url https://asia-northeast1-example.run.app/hook \
  -headers "Content-Type=application/json" \
  -message-body '{"task":"nightly"}' \
  -oidc-service-account-email cron-runner@example.iam.gserviceaccount.com
```

### Cloud SQL 起動ジョブを自動生成
```bash
go run ./cmd/cli/gcloud-genset-scheduler -operation create-start-cloud-sql-job \
  -project-id SAMPLE_PROJECT \
  -pubsub-topic projects/sample/topics/cloud-sql \
  -db-instance-id DB_INSTANCE_ID \
  -discord-webhook-url https://discordapp.com/api/webhooks/... \
  -cloud-sql-icon-url https://example.com/icons/cloud-sql.png
```

### ジョブを一時停止
```bash
go run ./cmd/cli/gcloud-genset-scheduler -operation pause-job \
  -job-name nightly-http-call \
  -location asia-northeast1
```

### ジョブのスケジュールを更新
```bash
go run ./cmd/cli/gcloud-genset-scheduler -operation update-pubsub-job \
  -job-name exec-cloud-run \
  -schedule '0 3 * * *'
```

## 補足

- `schedule`, `description`, `time-zone`, `location` などは省略するとスクリプト版 (`scheduler.sh`) と同様のデフォルト値が自動的に付与されます。
- Cloud SQL 系の operation では、Discord Webhook やアイコン URL を指定するとメッセージ本文の JSON に反映されます。
- 生成されたコマンドはコピー＆ペーストでそのまま実行できる形式（シングルクォート＋バックスラッシュによる改行）で出力されます。
