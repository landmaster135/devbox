# Gcloud Genset Monitoring

Google Cloud Monitoring の `gcloud` コマンドを生成する CLI ツールです。

## 概要

- **ダッシュボード一覧取得**: `gcloud monitoring dashboards list` のフィルターや並び替えをまとめて構築
- **ダッシュボード詳細取得**: ID を指定して `gcloud monitoring dashboards describe` コマンドを生成
- **スヌーズ/アップタイム設定一覧**: URI 表示やページサイズなどの指定を含めたコマンドを出力
- **ヘルプ対応**: `-help` フラグで利用可能なフラグと操作を確認

## インストール

```bash
# プロジェクトルートで実行
go build -o bin/gcloud-genset-monitoring ./cmd/cli/gcloud-genset-monitoring
```

## 使用方法

```bash
go run ./cmd/cli/gcloud-genset-monitoring -operation list-dashboards -project MY_PROJECT -format json
```

標準出力に生成された `gcloud` コマンドが強調表示で出力されます。そのままコピーして実行できます。

## オプション

| オプション | 説明 | 必須 | 対応オペレーション | 例 |
|------------|------|------|---------------------|----|
| `-operation` | 実行する操作 (`list-dashboards` / `describe-dashboard` / `list-snoozes` / `list-uptime-configs`) | * | 全て | `-operation list-dashboards` |
| `-project` | Google Cloud プロジェクト ID |  | 全て | `-project MY_PROJECT` |
| `-filter` | 結果をフィルタリングする式 |  | list-dashboards, list-snoozes, list-uptime-configs | `-filter displayName:test` |
| `-format` | 出力形式 (table, json, yaml など) |  | 全て | `-format json` |
| `-page-size` | 1ページあたりの取得件数 |  | list-dashboards, list-snoozes, list-uptime-configs | `-page-size 50` |
| `-sort-by` | 並び替えに使用するフィールド |  | list-dashboards, list-snoozes, list-uptime-configs | `-sort-by displayName` |
| `-limit` | 取得する最大件数 |  | list-dashboards, list-snoozes, list-uptime-configs | `-limit 20` |
| `-uri` | リソース URI を表示 |  | list-snoozes, list-uptime-configs | `-uri` |
| `-dashboard-id` | 参照するダッシュボード ID | * | describe-dashboard | `-dashboard-id dashboards/123456789` |
| `-help` | ヘルプを表示 |  | 全て | `-help` |

## オペレーション一覧

| オペレーション | 説明 | 必須パラメータ |
|----------------|------|----------------|
| `list-dashboards` | `gcloud monitoring dashboards list` コマンドを生成 | `-operation` |
| `describe-dashboard` | `gcloud monitoring dashboards describe` コマンドを生成 | `-operation`, `-dashboard-id` |
| `list-snoozes` | `gcloud monitoring snoozes list` コマンドを生成 | `-operation` |
| `list-uptime-configs` | `gcloud monitoring uptime list-configs` コマンドを生成 | `-operation` |

## 使用例

### ダッシュボード一覧を JSON で取得するコマンドの生成

```bash
$ go run ./cmd/cli/gcloud-genset-monitoring \
  -operation list-dashboards \
  -project MY_PROJECT \
  -filter 'displayName:test' \
  -page-size 20 \
  -sort-by displayName \
  -limit 50 \
  -format json
```

出力:
```
gcloud monitoring dashboards list --project=MY_PROJECT --filter=displayName:test --format="json" --page-size=20 --sort-by=displayName --limit=50
```

### ダッシュボード詳細コマンドの生成

```bash
$ go run ./cmd/cli/gcloud-genset-monitoring \
  -operation describe-dashboard \
  -project MY_PROJECT \
  -dashboard-id dashboards/123456789 \
  -format yaml
```

出力:
```
gcloud monitoring dashboards describe dashboards/123456789 --project=MY_PROJECT --format="yaml"
```

### Snooze 設定一覧コマンドの生成

```bash
$ go run ./cmd/cli/gcloud-genset-monitoring \
  -operation list-snoozes \
  -project MY_PROJECT \
  -filter 'displayName:maintenance' \
  -page-size 10 \
  -limit 20 \
  -uri
```

出力:
```
gcloud monitoring snoozes list --project=MY_PROJECT --filter=displayName:maintenance --page-size=10 --limit=20 --uri
```

## テスト

```bash
go test ./internal/gcloud_genset_monitoring/...
```

## ディレクトリ構成

```
internal/gcloud_genset_monitoring/
├── config/       # CLI フラグ解析とバリデーション
│   ├── config.go
│   ├── config_test.go
│   └── flag_parser.go
└── usecases/     # gcloud コマンド組み立てロジック
    ├── services.go
    └── services_test.go
```
