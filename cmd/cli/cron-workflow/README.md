# Cron Workflow

あらかじめ定義したバックグラウンドジョブを継続的に動かし続ける CLI ツールです。バイナリは単一のスケジューラを起動し、`workflow/core.go` に宣言されたワークフローを登録してから `SIGINT` もしくは `SIGTERM` を受け取るまで待機します。

## 同梱ワークフロー

| 説明 | 実行間隔 | タイムゾーン | 目的 |
|------|-----------|--------------|------|
| ハートビート監視 (毎分) | `*/1 * * * *` | Asia/Tokyo (既定) | 外形監視が生存確認できるよう短いログを継続的に出力します。 |
| 毎時の状態スナップショット | `15 * * * *` | Asia/Tokyo | `context.Context` を介したキャンセルを尊重する集計ジョブの例となります。 |
| 東京の天気通知 (毎朝) | `0 1 * * 0-6` | Asia/Tokyo | OpenWeatherMap から 3 日分の予報を取得し、Discord Webhook へ投稿します。 |
| PC情報スナップショット (10分毎) | `*/10 * * * *` | Asia/Tokyo | `machine-info` ユースケースの `CollectAndSaveUbuntuInfo` を呼び出し、`PC_INFO_OUTPUT_DIR` で指定した配下に JSON ログを書き出します。 |

ワークフローを追加・更新する場合は `workflow/core.go` を変更し、CLI を再ビルドしてください。

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
2026/01/20 10:00:00.123456 registered workflow "Heartbeat monitor (every minute)" (cron=CRON_TZ=Asia/Tokyo */1 * * * *)
2026/01/20 10:00:00.123789 registered workflow "Hourly state snapshot" (cron=CRON_TZ=UTC 0 * * * *)
2026/01/20 10:00:00.123900 scheduler started (2 workflow(s)). waiting for termination signal...
2026/01/20 10:01:00.125678 workflow "Heartbeat monitor (every minute)" completed
2026/01/20 10:02:00.130201 workflow "Heartbeat monitor (every minute)" completed
^C
2026/01/20 10:02:10.000001 signal received: interrupt. shutting down...
2026/01/20 10:02:10.000400 scheduler stopped cleanly
```

スケジューラは無期限で稼働するため、Ctrl+C で停止するか管理プロセスから終了シグナルを送って終了させてください。
