# Gcloud Genset Deployment CLI

Google Cloud Deployment Manager の運用で利用する `gcloud deployment-manager` コマンドを安全に組み立てて実行する CLI ツールです。

## サポートする operation

| operation            | 説明 |
|----------------------|------|
| `list-deployments`   | Deployment Manager のデプロイメント一覧を取得し、フィルターやフォーマットを指定した `gcloud deployment-manager deployments list` を実行します。|

## 使用例

### 1. 基本的なデプロイメント一覧取得

```bash
go run ./cmd/cli/gcloud-genset-deployment \
  --operation=list-deployments \
  --project=my-sample-project
```

Output (抜粋):
```bash
[INFO] list-deployments: デプロイメント一覧の取得が完了しました
```

### 2. シンプル表示と最大取得件数の制限

```bash
go run ./cmd/cli/gcloud-genset-deployment \
  --operation=list-deployments \
  --project=my-sample-project \
  --limit=5 \
  --simple
```

Output 例 (`gcloud` の標準出力):
```text
NAME          INSERT_TIME
sample-001    2024-07-01T10:22:33.123-07:00
...
[INFO] list-deployments: デプロイメント一覧の取得が完了しました
```

### 3. 実行コマンドの確認とカスタムフォーマット

```bash
go run ./cmd/cli/gcloud-genset-deployment \
  --operation=list-deployments \
  --project=my-sample-project \
  --filter="name:my-deployment*" \
  --format="table(name,description,operation.status)" \
  --show-command
```

Output (冒頭):
```bash
[INFO] list-deployments: 実行コマンド: gcloud deployment-manager deployments list --project=my-sample-project --filter="name:my-deployment*" --format="table(name,description,operation.status)"
...
[INFO] list-deployments: デプロイメント一覧の取得が完了しました
```

## CLI オプション

| operation          | 主なオプション |
|--------------------|----------------|
| `list-deployments` | `--project` (対象プロジェクト ID)、`--filter`、`--format`、`--limit`、`--simple` (簡易フォーマットへ強制切替)、`--show-command` (実行前にコマンド表示) |

`--simple` を指定するとフォーマットは `table(name,insertTime)` に固定されます。`--show-command` を併用すると、実行される `gcloud` コマンドを確認した上で操作できます。`gcloud` バイナリが見つからない場合はエラーを返し、プロセスは実行されません。

## ビルド

```bash
./scripts/build_gcloud-genset-deployment.sh
```

ビルド済みバイナリは `pkg/bin/gcloud-genset-deployment/<platform>/` 配下に生成されます。実行時は Google Cloud SDK がインストール済みであり、対象プロジェクトにアクセスできる認証情報が設定されている必要があります。
