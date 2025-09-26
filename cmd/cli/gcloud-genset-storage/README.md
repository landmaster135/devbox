# Gcloud Genset Storage CLI

Google Cloud Storage (GCS) 操作用の `gsutil` / `gcloud` コマンドを生成する CLI ツールです。Cloud Storage へのアップロード / ダウンロード / バケット作成 / ACL 操作など、日常的な運用コマンドを安全に組み立てます。

## サポートする操作

| operation             | 説明 |
|-----------------------|------|
| `upload-files`        | ローカルディレクトリを GCS バケットへアップロード (`gsutil -m cp -r`) |
| `download-files`      | GCS 上の複数オブジェクトをローカルへダウンロード (`gsutil -m cp`) |
| `create-bucket`       | 指定クラス / ロケーションでバケットを作成 (`gsutil mb`) |
| `list-contents`       | GCS バケット (またはローカルディレクトリ) の内容を列挙 (`gsutil ls` / `ls`) |
| `show-details`        | バケットやオブジェクトの詳細情報を表示 (`gsutil ls -Lb/-L`) |
| `list-names`          | バケット配下またはオブジェクトのパス一覧を取得 (`gsutil ls`) |
| `delete-object`       | オブジェクト / フォルダを削除 (`gsutil rm`) |
| `get-acl`             | ACL を取得 (`gsutil acl get`) |
| `set-acl`             | ACL をファイルから設定 (`gsutil acl set`) |
| `grant-read-all`      | 全ユーザに READ 権限を付与 (`gsutil acl ch -u AllUsers:R`) |
| `remove-read-all`     | 全ユーザの READ 権限を解除 (`gsutil acl ch -d AllUsers`) |

## 使用例

### 1. ディレクトリをバケットへアップロード

```bash
go run ./cmd/cli/gcloud-genset-storage \
  -operation upload-files \
  -local-path ./build/output \
  -bucket-url gs://my-artifact-bucket/releases/
```

Output:
```bash
gsutil -m cp -r './build/output' 'gs://my-artifact-bucket/releases/'
```

### 2. 複数オブジェクトをダウンロード

```bash
go run ./cmd/cli/gcloud-genset-storage \
  -operation download-files \
  -sources gs://data/file1.csv,gs://data/file2.csv \
  -destination ./downloads
```

Output:
```bash
gsutil -m cp 'gs://data/file1.csv' 'gs://data/file2.csv' './downloads'
```

### 3. バケットの詳細を確認

```bash
go run ./cmd/cli/gcloud-genset-storage \
  -operation show-details \
  -target gs://my-artifact-bucket/
```

Output:
```bash
gsutil ls -Lb 'gs://my-artifact-bucket/'
```

## CLI オプション

| operation                | 必須オプション |
|--------------------------|----------------|
| `upload-files`           | `-local-path`, `-bucket-url` |
| `download-files`         | `-sources`, `-destination` |
| `create-bucket`          | `-bucket-url`, `-storage-class`, `-location` |
| `list-contents`          | `-target` |
| `show-details` / `list-names` | `-target (gs://...)` |
| `delete-object`          | `-target (gs://...)` |
| `get-acl`                | `-target (gs://...)` |
| `set-acl`                | `-acl-file`, `-target (gs://...)` |
| `grant-read-all`         | `-target (gs://...)` |
| `remove-read-all`        | `-target (gs://...)` |

`sources` はカンマ区切りで複数指定できます（例: `gs://bucket/file1,gs://bucket/file2`）。

## ビルド

```bash
go build -o bin/gcloud-genset-storage ./cmd/cli/gcloud-genset-storage
```

## 注意事項

- 生成されたコマンドは **実行されません**。必要に応じてコピーし、作業環境で実行してください。
- `gs://` で始まらない操作はローカルディレクトリ／ファイルに対して `ls` を使用します。
- ACL 操作はバケットレベル／オブジェクトレベル共に可能ですが、Uniform bucket-level access を有効にしている場合は Google Cloud の制約に従ってください。
