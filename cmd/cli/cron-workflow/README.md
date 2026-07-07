# Cron Workflow

あらかじめ定義したバックグラウンドジョブを継続的に動かし続ける CLI ツールです。バイナリは単一のスケジューラを起動し、`workflow/core.go` に宣言されたワークフローを登録してから `SIGINT` もしくは `SIGTERM` を受け取るまで待機します。

## 同梱ワークフロー

| 説明 | 実行間隔 | 目的 |
|---|---|---|
| 東京の天気通知 | `0 1 * * 0-6` | OpenWeatherMap から 3 日分の予報を取得し、Discord Webhook へ投稿します。 |
| 日次見出し通知 | `1 0 * * 0-6` | 日次テンプレート見出しを Discord Webhook へ投稿します。 |
| PostgreSQL ダンプ通知 | `0 2 * * 0-6` | staging / production DB をダンプし、`attachments` テーブルのデータだけ別 SQL ファイルへ分離して、結果を Discord Webhook へ投稿します。 |
| Memos PostgreSQL ダンプ通知 | `5 2 * * 0-6` | memos staging / production DB をダンプし、結果を Discord Webhook へ投稿します。 |
| PC情報スナップショット | `*/10 * * * 0-6` | `machine-info` ユースケースの `CollectAndSaveUbuntuInfo` を呼び出し、`PC_INFO_OUTPUT_DIR` で指定した配下に JSON ログを書き出します。 |

ワークフローを追加・更新する場合は `workflow/core.go` を変更し、CLI を再ビルドしてください。

## PostgreSQL ダンプ関連の必須環境変数

- `DATABASE_URL_01_STAGING`
- `DATABASE_DUMP_DIR_01_STAGING`
- `DATABASE_URL_01_PRODUCT`
- `DATABASE_DUMP_DIR_01_PRODUCT`
- `DATABASE_URL_01_MEMOS_STAGING`
- `DATABASE_DUMP_DIR_01_MEMOS_STAGING`
- `DATABASE_URL_01_MEMOS_PROD`
- `DATABASE_DUMP_DIR_01_MEMOS_PROD`

## ビルド

```bash
# リポジトリルートで実行
go build -o bin/cron-workflow ./cmd/cli/cron-workflow
```

## 使用例

```bash
# Ctrl+C が押されるまでスケジューラが動作
go run ./cmd/cli/cron-workflow
```

利用可能な CLI フラグは `-h` / `--help` のみです。

## 出力例

```
$ go run ./cmd/cli/cron-workflow
2026/01/20 10:00:00.123456 registered workflow "Daily Tokyo weather notification" (cron=CRON_TZ=Asia/Tokyo 0 1 * * 0-6)
2026/01/20 10:00:00.123789 registered workflow "Daily heading Discord notification" (cron=CRON_TZ=Asia/Tokyo 1 0 * * 0-6)
2026/01/20 10:00:00.124010 registered workflow "Daily PostgreSQL dump with notification" (cron=CRON_TZ=Asia/Tokyo 0 2 * * 0-6)
2026/01/20 10:00:00.124200 registered workflow "Daily PostgreSQL dump for memos with notification" (cron=CRON_TZ=Asia/Tokyo 5 2 * * 0-6)
2026/01/20 10:00:00.124350 registered workflow "Ubuntu PC info snapshot" (cron=CRON_TZ=Asia/Tokyo */10 * * * 0-6)
2026/01/20 10:00:00.124900 scheduler started (5 workflow(s)). waiting for termination signal...
^C
2026/01/20 10:02:10.000001 signal received: interrupt. shutting down...
2026/01/20 10:02:10.000400 scheduler stopped cleanly
```

スケジューラは無期限で稼働するため、Ctrl+C で停止するか管理プロセスから終了シグナルを送って終了させてください。
