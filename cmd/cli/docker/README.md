# Docker CLI

env.yml に定義した情報を `docker-compose.yml` へ同期するユーティリティです。`npm run docker:env` 相当の環境変数反映に加え、指定サービスのポート番号を env.yml ベースで更新できます。

## 主な機能

- YAML 形式 (`KEY: VALUE`) の環境変数を順序を保ったまま読み込み
- `environment:` ブロックを指定ファイル内で検出し、一括で差し替え
- 既存のアンカー（`&foo`/`*foo`）構造を維持したままエントリのみ更新
- 値に空白が含まれる場合は自動でダブルクォートを付与
- `ports:` ブロックおよび `labels.tsdproxy.container_port` を任意のサービスに対して更新

## インストール

```bash
# プロジェクトルートからビルド
go build -o bin/docker-sync ./cmd/cli/docker
```

## 使用例

### 最低限の実行例

```bash
# プロジェクトルートで実行
go run ./cmd/cli/docker \
  -operation=env-into-compose \
  -compose-path=./docker-compose.yml \
  -env-yaml-path=./env.yml
```

出力例:
```
3件の環境変数を docker-compose.yml に反映しました
```

### フラグ一覧

| フラグ | 必須 | デフォルト | 説明 |
| --- | --- | --- | --- |
| `-operation` | * | なし | `env-into-compose` または `ports-into-compose` |
| `-compose-path` | | `docker-compose.yml` | 読み書き対象の docker-compose.yml パス |
| `-env-yaml-path` | | `env.yml` | 参照する環境変数YAMLのパス |
| `-port-key` | `ports-into-compose` 時 | なし | env.yml 内から参照するキー（例: `VITE_FRONT_URL_PORT`） |
| `-service` | `ports-into-compose` 時 | なし | 更新対象サービス名（例: `dathub`） |

### 典型的なワークフロー

1. `env.yml` を編集して `VITE_*` などの値を更新
2. `go run ./cmd/cli/docker -operation=env-into-compose` を実行して compose を同期
3. `pkg/docker/dockerize_dev.sh` などのスクリプトでビルド／デプロイを実施

### ports-into-compose の例

```bash
go run ./cmd/cli/docker \
  -operation=ports-into-compose \
  -compose-path=./docker-compose.yml \
  -env-yaml-path=./env.yml \
  -port-key=CRON_URL_PORT \
  -service=devbox
```

`env.yml` の `CRON_URL_PORT` を読み取り、`services.devbox.ports` と `services.devbox.labels.tsdproxy.container_port` が同期されます。

### YAML フォーマットについて

- `KEY: VALUE` 形式で記述します
- インラインコメント（`# comment`）は自動的に除去されます
- シングル／ダブルクォートを付けた値もサポートされます

```yaml
VITE_GIN_MODE: production          # コメントは無視されます
SIMPLE_VALUE: simple
VALUE_WITH_SPACE: "value with space"
```

## エラー例

| メッセージ | 意味 |
| --- | --- |
| `--operation は必須です` | フラグ未指定。`-operation=env-into-compose` を付与してください |
| `docker-compose.yml に environment セクションが見つかりませんでした` | 対象ファイルに `environment:` ブロックが存在しないため更新できません |
| `XXX には有効な環境変数が含まれていません` | `env.yml` にパース可能なエントリがない状態です |
| `YYY が env.yml に存在しません` | `ports-into-compose` で指定した `-port-key` が見つからない |
| `SERVICE の ports セクションを更新できませんでした` | `ports-into-compose` 対象サービスに `ports:` が無い |
| `SERVICE の labels 内に tsdproxy.container_port が見つかりません` | `ports-into-compose` 対象サービスに該当ラベルが無い |

## 関連スクリプト

- `pkg/docker/build_frontend_image.sh`: Vite 用のビルド引数を env.yml から生成して Docker イメージを作成
- `pkg/docker/dockerize_dev.sh`: 本CLIで同期後、開発用に docker compose を起動
- `pkg/docker/dockerize.sh`: 同期 & ビルド後にイメージをアーカイブ
