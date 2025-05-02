# 環境変数ローダー

YAMLファイルから環境変数を設定するためのユーティリティです。

## 機能

- YAMLファイルから環境変数を読み込み、設定します
- コマンドラインからYAMLファイルのパスを指定できます
- YAMLファイルが指定されなかった場合、デフォルトで`env.yml`を使用します

## 使用方法

### 1. Go実行ファイルを使用する方法

```bash
# ビルド
go build -o bin/env-loader cmd/env-loader/main.go

# 実行（デフォルトのenv.ymlを使用）
./bin/env-loader

# 実行（カスタムYAMLファイルを指定）
./bin/env-loader -env path/to/custom_env.yml
```

**注意**: この方法では、環境変数は`env-loader`プロセス内でのみ設定されます。親プロセス（現在のシェル）には影響しません。

## YAMLファイルの形式

```yaml
# 環境変数の設定
DB_HOST: localhost
DB_PORT: 5432
DB_USER: devbox
DB_PASSWORD: password
DB_NAME: devbox_db
API_KEY: your_api_key
DEBUG: true
```

## 実装の詳細

このユーティリティは、クリーンアーキテクチャに基づいて実装されています。

- `domain/models`: ドメインモデル（`EnvConfig`）
- `domain/repositories`: リポジトリインターフェース（`EnvRepository`）
- `interfaces/repositories`: リポジトリの実装（`EnvRepositoryImpl`）
- `usecases/services`: ユースケースの実装（`EnvService`）
- `cmd/env-loader`: コマンドラインインターフェース\
