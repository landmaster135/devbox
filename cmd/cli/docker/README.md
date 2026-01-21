# Docker CLI

env.yml に定義した環境変数を `docker-compose.yml` の `environment` セクションへ反映するユーティリティです。`npm run docker:env` 相当の同期処理を Go 製 CLI として提供します。

## 主な機能

- YAML 形式 (`KEY: VALUE`) の環境変数を順序を保ったまま読み込み
- `environment:` ブロックを指定ファイル内で検出し、一括で差し替え
- 既存のアンカー（`&foo`/`*foo`）構造を維持したままエントリのみ更新
- 値に空白が含まれる場合は自動でダブルクォートを付与

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
| `-operation` | * | なし | 実行する操作。現在は `env-into-compose` のみ対応 |
| `-compose-path` | | `docker-compose.yml` | 環境変数を書き込む docker-compose.yml のパス |
| `-env-yaml-path` | | `env.yml` | 参照する環境変数YAMLのパス |

### 典型的なワークフロー

1. `env.yml` を編集して `VITE_*` などの値を更新
2. `go run ./cmd/cli/docker -operation=env-into-compose` を実行して compose を同期
3. `pkg/docker/dockerize_dev.sh` などのスクリプトでビルド／デプロイを実施

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

## 関連スクリプト

- `pkg/docker/build_frontend_image.sh`: Vite 用のビルド引数を env.yml から生成して Docker イメージを作成
- `pkg/docker/dockerize_dev.sh`: 本CLIで同期後、開発用に docker compose を起動
- `pkg/docker/dockerize.sh`: 同期 & ビルド後にイメージをアーカイブ
