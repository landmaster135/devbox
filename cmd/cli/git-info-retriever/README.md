# Git情報取得・アーカイブツール

GitHubからリポジトリ情報を取得し、Bash関数を生成してリポジトリのクローンとZip圧縮を行うCLIツールです。

## 概要

このツールは、指定されたサービス（現在はGitHubのみサポート）からリポジトリ情報を取得し、以下の機能を提供します：

1. **リポジトリ情報取得**: JSON形式でリポジトリ情報を出力
2. **アーカイブ機能**: git cloneとzip圧縮を行うBash関数を生成

## 機能

### リポジトリ情報取得 (`retrieve`)
- GitHubの認証されたユーザーが所有するリポジトリ一覧を取得
- 各リポジトリの詳細情報を取得（言語、スター数、フォーク数、プルリクエスト数など）
- 並行処理による高速な情報取得
- JSON形式での結果出力

### アーカイブ機能 (`archive`)
- リポジトリ情報からBash関数を生成
- **実際のリポジトリのアーカイブ処理は実行されません**
- `archive_repos()`: git cloneとzip圧縮を実行
- `display_zipinfo()`: zipファイルの情報を表示
- `unzip_repos()`: zipファイルを展開
- GitHubから直接取得または既存JSONファイルから読み込み

## 使用方法

### 基本的な使用方法

```bash
# リポジトリ情報取得
go run ./cmd/cli/git-info-retriever/main.go -operation retrieve -service github -token <GitHubアクセストークン>

# アーカイブ機能（GitHubから取得）
go run ./cmd/cli/git-info-retriever/main.go -operation archive -service github -token <GitHubアクセストークン>

# アーカイブ機能（既存ファイルから）
go run ./cmd/cli/git-info-retriever/main.go -operation archive -service github -src-file ./repos.json
```

### パラメータ

| パラメータ | 必須 | 説明 |
|-----------|------|------|
| `-operation` | * | 操作タイプ（`retrieve`: リポジトリ情報取得、`archive`: Bash関数生成） |
| `-service` | * | サービスタイプ（現在は`github`のみサポート） |
| `-token` | 条件付き | GitHubアクセストークン（`retrieve`操作では必須、`archive`操作で`src-file`未指定の場合は必須） |
| `-save-file-path` | - | 結果を保存するファイルパス（指定しない場合は標準出力） |
| `-output-command-file-path` | - | Bash関数出力ファイルパス（`archive`操作で使用、指定しない場合は標準出力） |
| `-archive-dir` | - | アーカイブディレクトリ（`archive`操作で使用、デフォルト: `./archives`） |
| `-src-file` | - | 既存ファイルからリポジトリ情報を読み込み（`archive`操作で使用、指定時は`token`は不要） |
| `-help` | - | ヘルプを表示 |

## 使用例

### リポジトリ情報取得

```bash
# GitHubリポジトリ情報を標準出力に表示
go run ./cmd/cli/git-info-retriever/main.go -operation retrieve -service github -token ghp_xxxxxxxxxxxxxxxxxxxx

# GitHubリポジトリ情報をファイルに保存
go run ./cmd/cli/git-info-retriever/main.go -operation retrieve -service github -token ghp_xxxxxxxxxxxxxxxxxxxx -save-file-path ./repos.json

# 深いディレクトリに保存（ディレクトリは自動作成）
go run ./cmd/cli/git-info-retriever/main.go -operation retrieve -service github -token ghp_xxxxxxxxxxxxxxxxxxxx -save-file-path ./results/github-repos.json
```

### アーカイブ機能

```bash
# GitHubから取得してBash関数を標準出力に表示
go run ./cmd/cli/git-info-retriever/main.go -operation archive -service github -token ghp_xxxxxxxxxxxxxxxxxxxx

# GitHubから取得してBash関数をファイルに保存
go run ./cmd/cli/git-info-retriever/main.go -operation archive -service github -token ghp_xxxxxxxxxxxxxxxxxxxx -output-command-file-path ./archive_commands.sh

# カスタムアーカイブディレクトリを指定して、、ターミナルにBash関数を出力
go run ./cmd/cli/git-info-retriever/main.go -operation archive -service github -token ghp_xxxxxxxxxxxxxxxxxxxx -archive-dir ./my-archives

# 既存JSONファイルからBash関数を生成
go run ./cmd/cli/git-info-retriever/main.go -operation archive -service github -src-file ./repos.json -output-command-file-path ./archive_commands.sh

# ヘルプを表示
go run ./cmd/cli/git-info-retriever/main.go -help
```

### 生成されたBash関数の使用例

```bash
# 生成されたBash関数を読み込み
source ./archive_commands.sh

# リポジトリをクローンしてZip圧縮
archive_repos

# Zipファイルの情報を表示
display_zipinfo

# Zipファイルを展開
unzip_repos
```

## GitHubアクセストークンの取得方法

1. GitHubにログインし、Settings > Developer settings > Personal access tokens > Tokens (classic) に移動
2. "Generate new token (classic)" をクリック
3. 必要なスコープを選択（`repo`スコープが必要）
4. トークンを生成し、安全に保管

## 出力形式

### リポジトリ情報（`retrieve`操作）

ツールは以下の形式でJSON配列を出力します：

```json
[
  {
    "name": "repository-name",
    "description": "Repository description",
    "is_private": false,
    "http_url": "https://github.com/user/repository-name",
    "language": "Go",
    "languages": {
      "Go": 12345,
      "JavaScript": 6789
    },
    "repo_created_at": "2023-01-01T00:00:00Z",
    "repo_updated_at": "2023-12-31T23:59:59Z",
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

#### 出力フィールドの説明

| フィールド | 型 | 説明 |
|-----------|----|----|
| `name` | string | リポジトリ名 |
| `description` | string | リポジトリの説明 |
| `is_private` | boolean | プライベートリポジトリかどうか |
| `http_url` | string | リポジトリのHTTP URL |
| `language` | string | メイン言語 |
| `languages` | object | 使用言語とバイト数のマップ |
| `repo_created_at` | string | 作成日時（ISO 8601形式） |
| `repo_updated_at` | string | 最終更新日時（ISO 8601形式） |
| `stargazers_count` | number | スター数 |
| `forks_count` | number | フォーク数 |
| `issues_count` | number | オープンなイシュー数 |
| `pulls_count` | number | プルリクエスト数（全状態） |
| `size` | number | リポジトリサイズ（KB） |
| `subscribers_count` | number | ウォッチャー数 |
| `is_archived` | boolean | アーカイブされているかどうか |

### Bash関数（`archive`操作）

以下の形式でBash関数を生成します：

```bash
function archive_repos() {
	git clone https://github.com/user/repo1
	zip -rq ./archives/repo1.zip ./repo1
	git clone https://github.com/user/repo2
	zip -rq ./archives/repo2.zip ./repo2
}

function display_zipinfo() {
	zipinfo ./archives/repo1.zip
	zipinfo ./archives/repo2.zip
}

function unzip_repos() {
	unzip ./archives/repo1.zip -d ./unarchived
	unzip ./archives/repo2.zip -d ./unarchived
}
```

#### 生成される関数の説明

| 関数名 | 説明 |
|-------|------|
| `archive_repos()` | 各リポジトリをgit cloneし、指定されたアーカイブディレクトリにZip圧縮 |
| `display_zipinfo()` | 各Zipファイルの情報を表示 |
| `unzip_repos()` | 各Zipファイルを`./unarchived`ディレクトリに展開 |

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
- ファイル読み込みエラーやJSON解析エラーも適切に処理
- ネットワークエラーやその他のAPIエラーも適切に処理

## 注意事項

- GitHubアクセストークンは機密情報です。安全に管理してください
- GitHub APIのレート制限（認証済みユーザーは1時間に5000リクエスト）に注意してください
- 大量のリポジトリがある場合、処理に時間がかかる場合があります
- 生成されたBash関数を実行する前に、十分なディスク容量があることを確認してください
- git cloneを実行するため、適切なGit設定（認証情報など）が必要です

## ワークフロー例

1. **リポジトリ情報を取得してファイルに保存**
   ```bash
   go run ./cmd/cli/git-info-retriever/main.go -operation retrieve -service github -token <token> -save-file-path ./repos.json
   ```

2. **保存したファイルからBash関数を生成**
   ```bash
   go run ./cmd/cli/git-info-retriever/main.go -operation archive -service github -src-file ./repos.json -output-command-file-path ./archive_commands.sh
   ```

3. **生成されたBash関数を実行**
   ```bash
   source ./archive_commands.sh
   archive_repos
   ```

## 技術仕様

- **言語**: Go
- **依存関係**:
  - `github.com/google/go-github`: GitHub API クライアント
  - `golang.org/x/oauth2`: OAuth2 認証
- **アーキテクチャ**: Clean Architecture
  - `config`: 設定とCLIパラメータ解析
  - `usecases`: ビジネスロジック
  - `main.go`: エントリーポイント

## ライセンス

このプロジェクトのライセンスについては、プロジェクトルートのLICENSEファイルを参照してください。
