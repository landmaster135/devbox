# Gcloud Genset Init

Google Cloud プロジェクトの初期設定で利用する `gcloud` コマンドを生成する CLI ツールです。`auth login` と `config set project` を安全に組み立て、都度のコマンド入力を省力化します。

## 概要

- **認証コマンド生成**: プロジェクト ID 指定で `gcloud auth login` コマンドを生成
- **プロジェクト設定**: `gcloud config set project` のフラグを整理して出力
- **追加引数対応**: 必要に応じて任意のフラグをコマンドへ付与可能

## インストール

```bash
# プロジェクトルートで実行
go build -o bin/gcloud-genset-init ./cmd/cli/gcloud-genset-init
```

## 使用方法

```bash
go run ./cmd/cli/gcloud-genset-init -operation auth-login -project-id MY_PROJECT_123
```

生成された `gcloud` コマンドが標準出力に表示されます。

## オプション

| オプション | 説明 | 必須 | 対応オペレーション | 例 |
|------------|------|------|---------------------|----|
| `-operation` | 実行する操作 (`auth-login` / `set-project-config`) | * | 全て | `-operation auth-login` |
| `-project-id` | 対象となるプロジェクト ID | * | auth-login, set-project-config | `-project-id MY_PROJECT_123` |
| `-additional-args` | 生成コマンドへ付与する追加引数 |  | 全て | `-additional-args '--quiet'` |
| `-help` | ヘルプの表示 |  | 全て | `-help` |

## オペレーション一覧

| オペレーション | 説明 | 必須パラメータ |
|----------------|------|----------------|
| `auth-login` | `gcloud auth login` コマンドを生成 | `-operation`, `-project-id` |
| `set-project-config` | `gcloud config set project` コマンドを生成 | `-operation`, `-project-id` |

## 使用例

### 認証コマンドの生成

```bash
$ go run ./cmd/cli/gcloud-genset-init \
  -operation auth-login \
  -project-id MY_PROJECT_123 \
  -additional-args '--quiet'
```

Output:
```bash
gcloud auth login 'MY_PROJECT_123' --quiet
```

### プロジェクト設定コマンドの生成

```bash
$ go run ./cmd/cli/gcloud-genset-init \
  -operation set-project-config \
  -project-id MY_PROJECT_123
```

Output:
```bash
gcloud config set project 'MY_PROJECT_123'
```

## テスト

```bash
go test ./internal/gcloud_genset_init/...
```
