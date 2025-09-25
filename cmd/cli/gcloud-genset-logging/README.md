# Gcloud Genset Logging

Google Cloud Logging の `gcloud` コマンドを生成する CLI ツールです。`logging read` や `logging sinks create` の呼び出しに必要なフラグを整理し、適切なフィルターや引数を組み立てます。

## 概要

- **フィルター組み立て**: 重要度・リソースタイプ・任意クエリから `logging read` 用フィルター文字列を生成
- **シンク作成**: シンク名・宛先・ログフィルターを基に `logging sinks create` コマンドを構築
- **追加引数対応**: 生成したコマンドに任意の追加引数を付与可能
- **CLI 単体**: 標準出力に生成結果を返すシンプルなツール

## インストール

```bash
# プロジェクトルートで実行
go build -o bin/gcloud-genset-logging ./cmd/cli/gcloud-genset-logging
```

## 使用方法

```bash
go run ./cmd/cli/gcloud-genset-logging -operation logging-read -severity ERROR -resource-type gce_instance
```

生成された `gcloud` コマンドが標準出力に表示されます。

## オプション

| オプション | 説明 | 必須 | 対応オペレーション | 例 |
|------------|------|------|---------------------|----|
| `-operation` | 実行する操作 (`logging-read` / `create-sink`) | * | 全て | `-operation logging-read` |
| `-severity` | ログの重要度 | | logging-read | `-severity ERROR` |
| `-limit` | 取得件数 (既定: 10) | | logging-read | `-limit 20` |
| `-query` | 追加クエリ | | logging-read | `-query textPayload:"Database"` |
| `-resource-type` | リソースタイプ | | logging-read | `-resource-type k8s_container` |
| `-filter` | 完全なフィルター文字列 | | logging-read | `-filter 'resource.type=gce_instance'` |
| `-additional-args` | 生成コマンドに付与する追加引数 | | 全て | `-additional-args '--format=json'` |
| `-sink-name` | 作成するシンク名 | * | create-sink | `-sink-name app-error-sink` |
| `-destination` | シンクの宛先 | * | create-sink | `-destination storage.googleapis.com/my-bucket` |
| `-log-filter` | シンク用フィルター | | create-sink | `-log-filter 'severity>=ERROR'` |
| `-help` | ヘルプ表示 | | 全て | `-help` |

## オペレーション一覧

| オペレーション | 説明 | 必須パラメータ |
|----------------|------|----------------|
| `logging-read` | `gcloud logging read` コマンドを生成 | `-operation`, 以下いずれか: `-filter` / `-severity` / `-resource-type` / `-query` |
| `create-sink` | `gcloud logging sinks create` コマンドを生成 | `-operation`, `-sink-name`, `-destination` |

## 使用例

### ログ取得コマンドの生成

```bash
$ go run ./cmd/cli/gcloud-genset-logging \
  -operation logging-read \
  -severity ERROR \
  -resource-type gce_instance \
  -limit 25 \
  -additional-args '--format=json'
```

Output:
```bash
gcloud logging read "severity>=ERROR AND resource.type=gce_instance" --limit=25 --format=json
```

### ログシンク作成コマンドの生成

```bash
$ go run ./cmd/cli/gcloud-genset-logging \
  -operation create-sink \
  -sink-name my-error-sink \
  -destination storage.googleapis.com/my-bucket \
  -log-filter 'severity>=ERROR'
```

Output:
```bash
gcloud logging sinks create my-error-sink storage.googleapis.com/my-bucket --log-filter="severity>=ERROR"
```

## アーキテクチャ

```
internal/gcloud_genset_logging/
├── config/       # CLI フラグ解析と設定検証
│   ├── config.go
│   ├── config_test.go
│   └── flag_parser.go
└── usecases/     # gcloud コマンド組み立てロジック
    ├── services.go
    └── services_test.go
```

## テスト

```bash
go test ./internal/gcloud_genset_logging/...
```
