# Gcloud Genset Secret

Google Cloud Secret Manager の `gcloud` コマンドを生成する CLI ツールです。シークレットの作成・値の登録・バージョン取得・ラベル/エイリアス更新に必要なフラグを整理し、実行可能なコマンドラインを出力します。

## 概要

- **シークレット作成**: レプリケーションポリシーとロケーション指定に対応した `secrets create` コマンドを生成
- **値の登録**: `echo -n` によるパイプ構文を含む `secrets versions add` コマンドを作成
- **ワンショット作成**: シークレット作成と値登録を 1 本のコマンドで連結 (`&&`) 出力
- **参照/更新**: バージョンアクセス・ラベル更新・エイリアス更新コマンドを出力
- **CLI 単体**: 標準出力に生成結果を返すシンプルなツール
- **Discord通知コマンド**: Secret Manager操作の成功/失敗をDiscordに送るためのコマンドも併せて出力

## インストール

```bash
# プロジェクトルートで実行
go build -o bin/gcloud-genset-secret ./cmd/cli/gcloud-genset-secret
```

## 使用方法

```bash
go run ./cmd/cli/gcloud-genset-secret -operation create-secret -secret-name my-secret
```

生成された `gcloud` コマンドが標準出力に表示されます。

## オプション

| オプション | 説明 | 必須 | 対応オペレーション | 例 |
|------------|------|------|---------------------|----|
| `-operation` | 実行する操作 (`create-secret` など) | * | 全て | `-operation create-secret` |
| `-secret-name` | 対象となるシークレット名 | * | create-secret / add-secret-version / create-and-add-secret-version / access-secret-version / update-secret-labels / update-secret-version-aliases | `-secret-name app-secret` |
| `-replication-policy` | レプリケーションポリシー (`automatic` / `user-managed`) | | create-secret / create-and-add-secret-version | `-replication-policy user-managed` |
| `-locations` | `user-managed` のロケーション (カンマ区切り) | 条件付き | create-secret / create-and-add-secret-version | `-locations asia-northeast1,us-east1` |
| `-secret-value` | 登録するシークレットの値 | * | add-secret-version / create-and-add-secret-version | `-secret-value 'super-secret'` |
| `-version` | 参照するバージョン (既定: `latest`) | | access-secret-version | `-version 3` |
| `-labels` | ラベル更新値 (`KEY=VALUE,...`) | * | update-secret-labels | `-labels env=prod,team=platform` |
| `-alias-option` | エイリアス更新オプション (`--clear-version-aliases` 等) | * | update-secret-version-aliases | `-alias-option --update-version-aliases=prod=5` |
| `-help` | ヘルプを表示 | | 全て | `-help` |

## オペレーション一覧

| オペレーション | 説明 | 必須パラメータ |
|----------------|------|----------------|
| `create-secret` | `gcloud secrets create` コマンドを生成 | `-operation`, `-secret-name` (必要に応じ `-locations`) |
| `add-secret-version` | `gcloud secrets versions add` コマンドを生成 | `-operation`, `-secret-name`, `-secret-value` |
| `create-and-add-secret-version` | シークレット作成と値登録を 1 本のコマンドで生成 | `-operation`, `-secret-name`, `-secret-value` |
| `access-secret-version` | `gcloud secrets versions access` コマンドを生成 | `-operation`, `-secret-name` |
| `update-secret-labels` | `gcloud secrets update --update-labels` コマンドを生成 | `-operation`, `-secret-name`, `-labels` |
| `update-secret-version-aliases` | `gcloud secrets update` によるバージョンエイリアス操作コマンドを生成 | `-operation`, `-secret-name`, `-alias-option` |

## 使用例

### シークレットを作成し値まで登録

```bash
go run ./cmd/cli/gcloud-genset-secret \
  -operation create-and-add-secret-version \
  -secret-name SECRET_NAME \
  -secret-value 'SECRET_PASSWORD'
```

Output:
```bash
gcloud secrets create 'SECRET_NAME' --replication-policy='automatic' && echo -n 'SECRET_PASSWORD' | gcloud secrets versions add 'SECRET_NAME' --data-file=-
```

### user-managed でシークレットを作成

```bash
go run ./cmd/cli/gcloud-genset-secret \
  -operation create-secret \
  -secret-name multi-region-secret \
  -replication-policy user-managed \
  -locations asia-northeast1,us-east1
```

Output:
```bash
gcloud secrets create 'multi-region-secret' --replication-policy='user-managed' --locations='asia-northeast1,us-east1'
```

### バージョンを取得

```bash
go run ./cmd/cli/gcloud-genset-secret \
  -operation access-secret-version \
  -secret-name SECRET_NAME \
  -version 5
```

Output:
```bash
gcloud secrets versions access '5' --secret='SECRET_NAME'
```

### Discord通知コマンド

各 `gcloud` コマンドの直後に、Secret Manager 向け成功/失敗通知を送る Discord Webhook コマンドも表示されます。`DISCORD_WEBHOOK_URL_FOR_IAC_ON_GCLOUD` 環境変数を設定した上で、必要に応じて成功・失敗いずれかのコマンドを実行してください。

```bash
通知付きシェルコマンド
==============================
$HOME/devbox/pkg/bin/cli/linux_amd64/discord-webhook \
  -webhook-url "$DISCORD_WEBHOOK_URL_FOR_IAC_ON_GCLOUD" \
  -content-text 'シークレットを作るよ！' \
  -embed-type 'none'
if gcloud secrets create 'my-secret' --replication-policy='automatic'; then
  $HOME/devbox/pkg/bin/cli/linux_amd64/discord-webhook \
    -webhook-url "$DISCORD_WEBHOOK_URL_FOR_IAC_ON_GCLOUD" \
    -content-text '作ったよ！' \
    -embed-type 'google-secret-manager-success' \
    -embed-text 'シークレットを作ったよ！'
else
  $HOME/devbox/pkg/bin/cli/linux_amd64/discord-webhook \
    -webhook-url "$DISCORD_WEBHOOK_URL_FOR_IAC_ON_GCLOUD" \
    -content-text '失敗…' \
    -embed-type 'google-secret-manager-failed' \
    -embed-text 'シークレットが作れなかったよ…'
fi
```
通知内容の詳細カスタマイズや CLI の利用例については [`cmd/cli/discord-webhook/README.md`](../discord-webhook/README.md) も参照してください。

## アーキテクチャ

```
internal/gcloud_genset_secret/
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
go test ./internal/gcloud_genset_secret/...
```
