# Gcloud Genset DNS

Google Cloud DNS の `gcloud dns` コマンドを組み立てる CLI ツールです。`managed-zones list` 実行時に設定すべきフラグを整理し、再利用しやすいコマンド文字列として出力します。

## 概要

- **Managed Zones 一覧取得**: `gcloud dns managed-zones list` を必要なフラグ付きで生成
- **柔軟なフラグ指定**: プロジェクト・フォーマット・フィルター・ページングなどをオプションで指定
- **追加引数サポート**: 生成したコマンド末尾に任意の引数を付与可能
- **CLI 単体**: 標準出力に生成コマンドをハイライト表示

## インストール

```bash
# プロジェクトルートで実行
go build -o bin/gcloud-genset-dns ./cmd/cli/gcloud-genset-dns
```

## 使用方法

```bash
go run ./cmd/cli/gcloud-genset-dns -operation managed-zones-list -project my-project -format json
```

生成された `gcloud` コマンドが標準出力に表示されます。

## オプション

| オプション | 説明 | 必須 | 対応オペレーション | 例 |
|------------|------|------|---------------------|----|
| `-operation` | 実行する操作 (`managed-zones-list`) | * | 全て | `-operation managed-zones-list` |
| `-project` | 対象の GCP プロジェクト ID | | managed-zones-list | `-project my-project` |
| `-format` | 出力フォーマット (yaml, json, csv 等) | | managed-zones-list | `-format json` |
| `-filter` | 結果を絞り込むフィルター条件 | | managed-zones-list | `-filter 'name:example'` |
| `-limit` | 取得件数 (0 で指定なし) | | managed-zones-list | `-limit 20` |
| `-page-size` | 1 ページあたりの件数 (0 で指定なし) | | managed-zones-list | `-page-size 10` |
| `-sort-by` | ソート基準 | | managed-zones-list | `-sort-by NAME` |
| `-verbosity` | 詳細レベル (debug / info / warning など) | | managed-zones-list | `-verbosity debug` |
| `-uri` | URI 形式で出力 | | managed-zones-list | `-uri` |
| `-additional-args` | コマンド末尾に付与する追加引数 | | managed-zones-list | `-additional-args '--account=my-account'` |
| `-help` | ヘルプ表示 | | 全て | `-help` |

## オペレーション一覧

| オペレーション | 説明 | 必須パラメータ |
|----------------|------|----------------|
| `managed-zones-list` | `gcloud dns managed-zones list` コマンドを生成 | `-operation` |

## 使用例

```bash
$ go run ./cmd/cli/gcloud-genset-dns \
  -operation managed-zones-list \
  -project my-project \
  -filter "name:example" \
  -limit 50 \
  -uri
```

Output:
```bash
gcloud dns managed-zones list --project='my-project' --filter='name:example' --limit=50 --uri
```

## アーキテクチャ

```
internal/gcloud_genset_dns/
├── config/         # CLI フラグ解析と検証
│   ├── config.go
│   ├── config_test.go
│   └── flag_parser.go
└── usecases/       # gcloud dns コマンド生成ロジック
    ├── services.go
    └── services_test.go
```

## テスト

```bash
go test ./internal/gcloud_genset_dns/...
```
