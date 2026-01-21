# Docker CLI

env.yml に定義した情報を `docker-compose.yml` へ同期するユーティリティです。`npm run docker:env` 相当の環境変数反映に加え、指定サービスのポート番号やボリューム定義を env.yml ベースで更新できます。

## 主な機能

- YAML 形式 (`KEY: VALUE`) の環境変数を順序を保ったまま読み込み
- `environment:` ブロックを指定ファイル内で検出し、一括で差し替え
- 既存のアンカー（`&foo`/`*foo`）構造を維持したままエントリのみ更新
- 値に空白が含まれる場合は自動でダブルクォートを付与
- `ports:` ブロックおよび `labels.tsdproxy.container_port` を任意のサービスに対して更新
- `volumes:` ブロックを env.yml 内の構造化値（配列／マップ）で差し替え、bind mount などの設定を共有
- `user:` フィールドを env.yml の UID/GID 定義から同期し、実行ユーザーを一括変更

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
| `-operation` | * | なし | `env-into-compose` / `ports-into-compose` / `volumes-into-compose` / `user-into-compose` |
| `-compose-path` | | `docker-compose.yml` | 読み書き対象の docker-compose.yml パス |
| `-env-yaml-path` | | `env.yml` | 参照する環境変数YAMLのパス |
| `-port-key` | `ports-into-compose` 時 | なし | env.yml 内から参照するキー（例: `VITE_FRONT_URL_PORT`） |
| `-volume-key` | `volumes-into-compose` 時 | なし | env.yml 内から参照するボリューム定義のキー（例: `MOUNT_VOLUME`） |
| `-user-key` | `user-into-compose` 時 | なし | env.yml 内から参照するユーザー値のキー（例: `COMPOSE_USER`） |
| `-service` | `ports-into-compose` / `volumes-into-compose` / `user-into-compose` 時 | なし | 更新対象サービス名（例: `devbox`） |

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

### volumes-into-compose の例

```bash
go run ./cmd/cli/docker \
  -operation=volumes-into-compose \
  -compose-path=./docker-compose.yml \
  -env-yaml-path=./env.yml \
  -volume-key=MOUNT_VOLUME \
  -service=devbox
```

`env.yml` の `MOUNT_VOLUME` が次のように定義されていれば、その配列が `services.devbox.volumes` に展開されます。

```yaml
MOUNT_VOLUME:
  - type: bind
    source: /home/user/cron_output
    target: /app/volume
```

結果例:

```yaml
services:
  devbox:
    volumes:
      - type: bind
        source: /home/user/cron_output
        target: /app/volume
```

### user-into-compose の例

```bash
go run ./cmd/cli/docker \
  -operation=user-into-compose \
  -compose-path=./docker-compose.yml \
  -env-yaml-path=./env.yml \
  -user-key=COMPOSE_USER \
  -service=devbox
```

`env.yml` の `COMPOSE_USER` に `"8888:8888"` のような UID:GID 文字列を記述しておけば、`services.devbox.user` に反映されます。Compose 実行ユーザーをステージング／本番で切り替える用途に便利です。

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
