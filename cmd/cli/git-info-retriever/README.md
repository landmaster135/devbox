# Git情報取得ツール

GitHubからリポジトリ情報を取得するCLIツールです。

## 概要

このツールは、指定されたサービス（現在はGitHubのみサポート）からリポジトリ情報を取得し、JSON形式で出力します。

## 機能

- GitHubの認証されたユーザーが所有するリポジトリ一覧を取得
- 各リポジトリの詳細情報を取得（言語、スター数、フォーク数、プルリクエスト数など）
- 並行処理による高速な情報取得
- JSON形式での結果出力

## 使用方法

### 基本的な使用方法

```bash
go run ./cmd/cli/git-info-retriever/main.go -service github -token <GitHubアクセストークン>
```

### パラメータ

| パラメータ         | 必須 | 説明 |
|-------------------|------|------|
| `-service`        | ✓ | サービスタイプ（現在は`github`のみサポート） |
| `-token`          | ✓ | GitHubアクセストークン |
| `-save-file-path` | - | 結果を保存するファイルパス（指定しない場合は標準出力） |
| `-help`           | - | ヘルプを表示 |

## 使用例

```bash
# GitHubリポジトリ情報を標準出力に表示
go run ./cmd/cli/git-info-retriever/main.go -service github -token ghp_xxxxxxxxxxxxxxxxxxxx

# GitHubリポジトリ情報をファイルに保存
go run ./cmd/cli/git-info-retriever/main.go -service github -token ghp_xxxxxxxxxxxxxxxxxxxx -save-file-path ./output.json

# 深いディレクトリに保存（ディレクトリは自動作成）
go run ./cmd/cli/git-info-retriever/main.go -service github -token ghp_xxxxxxxxxxxxxxxxxxxx -save-file-path ./results/github-repos.json

# ヘルプを表示
go run ./cmd/cli/git-info-retriever/main.go -help
```

## GitHubアクセストークンの取得方法

1. GitHubにログインし、Settings > Developer settings > Personal access tokens > Tokens (classic) に移動
2. "Generate new token (classic)" をクリック
3. 必要なスコープを選択（`repo`スコープが必要）
4. トークンを生成し、安全に保管

## 出力形式

ツールは以下の形式でJSON配列を出力します：

```json
[
  {
    "name": "repository-name",
    "description": "Repository description",
    "is_private": false,
    "html_url": "https://github.com/user/repository-name",
    "language": "Go",
    "languages": {
      "Go": 12345,
      "JavaScript": 6789
    },
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-12-31T23:59:59Z",
    "stargazers_count": 10,
    "forks_count": 5,
    "issues_count": 2,
    "pulls_count": 3,
    "size": 1024,
    "subscribers_count": 8,
    "is_archived": false
  }
]
```

### 出力フィールドの説明

| フィールド | 型 | 説明 |
|-----------|----|----|
| `name` | string | リポジトリ名 |
| `description` | string | リポジトリの説明 |
| `is_private` | boolean | プライベートリポジトリかどうか |
| `html_url` | string | リポジトリのURL |
| `language` | string | メイン言語 |
| `languages` | object | 使用言語とバイト数のマップ |
| `created_at` | string | 作成日時（ISO 8601形式） |
| `updated_at` | string | 最終更新日時（ISO 8601形式） |
| `stargazers_count` | number | スター数 |
| `forks_count` | number | フォーク数 |
| `issues_count` | number | オープンなイシュー数 |
| `pulls_count` | number | プルリクエスト数（全状態） |
| `size` | number | リポジトリサイズ（KB） |
| `subscribers_count` | number | ウォッチャー数 |
| `is_archived` | boolean | アーカイブされているかどうか |

## ビルド方法

```bash
# devboxディレクトリで実行
cd devbox
go build -o bin/git-info-retriever ./cmd/cli/git-info-retriever
```

## エラーハンドリング

- 必須パラメータが不足している場合、エラーメッセージとヘルプを表示
- 無効なアクセストークンの場合、認証エラーを表示
- GitHub APIのレート制限に達した場合、適切なエラーメッセージを表示
- ネットワークエラーやその他のAPIエラーも適切に処理

## 注意事項

- GitHubアクセストークンは機密情報です。安全に管理してください
- GitHub APIのレート制限（認証済みユーザーは1時間に5000リクエスト）に注意してください
- 大量のリポジトリがある場合、処理に時間がかかる場合があります

## 技術仕様

- **言語**: Go 1.23.5
- **依存関係**:
  - `github.com/google/go-github`: GitHub API クライアント
  - `golang.org/x/oauth2`: OAuth2 認証
- **アーキテクチャ**: Clean Architecture
  - `config`: 設定とCLIパラメータ解析
  - `usecases`: ビジネスロジック
  - `main.go`: エントリーポイント

## ライセンス

このプロジェクトのライセンスについては、プロジェクトルートのLICENSEファイルを参照してください。
