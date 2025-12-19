# Gcloud Genset AI

Document AI の運用で利用する `gcloud`/`curl` コマンドを整形して出力する CLI ツールです。手動入力時に煩雑になりがちなアクセストークン取得やエンドポイントの組み立てを安全に自動化します。

## 概要

- **プロセッサバージョンのアンデプロイ支援**: `undeploy-processor-version` 操作で API 呼び出し用 `curl` コマンドを生成
- **必須パラメータのバリデーション**: リージョンやプロジェクト番号などの入力漏れを防止
- **コマンドハイライト出力**: 生成したコマンドを囲み枠で表示し、そのままコピーして利用可能

## インストール

```bash
# プロジェクトルートで実行
go build -o bin/gcloud-genset-ai ./cmd/cli/gcloud-genset-ai
```

## 使用方法

```bash
go run ./cmd/cli/gcloud-genset-ai \ 
  -operation undeploy-processor-version \ 
  -region us-central1 \ 
  -project-number 123456789012 \ 
  -processor-id PROC_ID \ 
  -version-id VERSION_ID
```

標準出力に生成されたコマンドが整形表示されます。

## オプション

| オプション | 説明 | 必須 | 対応オペレーション | 例 |
|------------|------|------|---------------------|----|
| `-operation` | 実行する操作 (`undeploy-processor-version`) | * | 全て | `-operation undeploy-processor-version` |
| `-region` | Document AI のリージョン | * | undeploy-processor-version | `-region us-central1` |
| `-project-number` | Google Cloud プロジェクト番号 | * | undeploy-processor-version | `-project-number 123456789012` |
| `-processor-id` | Document AI プロセッサ ID | * | undeploy-processor-version | `-processor-id a1b2c3d4` |
| `-version-id` | アンデプロイ対象のプロセッサバージョン ID | * | undeploy-processor-version | `-version-id 20240901` |
| `-help` | ヘルプを表示 |  | 全て | `-help` |

## オペレーション一覧

| オペレーション | 説明 | 必須パラメータ |
|----------------|------|----------------|
| `undeploy-processor-version` | プロセッサバージョンのアンデプロイ API 呼び出し用 `curl` コマンドを生成 | `-operation`, `-region`, `-project-number`, `-processor-id`, `-version-id` |

## 使用例

```bash
$ go run ./cmd/cli/gcloud-genset-ai \
  -operation undeploy-processor-version \
  -region us-central1 \
  -project-number 123456789012 \
  -processor-id 1234567890abcdef \
  -version-id 20240901
```

Output:
```bash
curl -s -X POST -H "Authorization: Bearer $(gcloud auth print-access-token)" -H "Content-Type: application/json" "https://us-central1-documentai.googleapis.com/v1beta3/123456789012/locations/us-central1/processors/1234567890abcdef/processorVersions/20240901:undeploy"
```

## テスト

```bash
go test ./internal/ai/...
```
