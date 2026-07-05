# GitHub CLI ツール

GitHubのイシューを取得するCLIツールです。

## 機能

- GitHubリポジトリのイシュー一覧取得
- 特定のイシューを番号で取得

## 使用方法

### 基本的な使用方法

```bash
# イシュー一覧取得
go run ./cmd/cli/github -operation list-issues -owner OWNER -repo REPO

# 短縮形を使用
go run ./cmd/cli/github -o list-issues -ow OWNER -r REPO
```

### オプションパラメータ付き

```bash
# 状態とソートを指定
go run ./cmd/cli/github -o list-issues -ow octocat -r Hello-World -state open -sort created -direction desc

# ページネーション
go run ./cmd/cli/github -o list-issues -ow octocat -r Hello-World -per-page 10 -page 2

# 短縮形でオプション指定
go run ./cmd/cli/github -o list-issues -ow octocat -r Hello-World -s closed -so updated -d asc -pp 50 -p 2
```

## 環境変数

GitHubトークンは CLI フラグではなく、環境変数または `.env` から読み取ります。

```bash
export GITHUB_PERSONAL_ACCESS_TOKEN=ghp_xxxxxxxxxxxx
```

## パラメータ

### 必須パラメータ

| パラメータ | 短縮形 | 説明 |
|-----------|--------|------|
| `-operation` | `-o` | 操作タイプ (list-issues) |
| `-owner` | `-ow` | リポジトリオーナー |
| `-repo` | `-r` | リポジトリ名 |

### 任意パラメータ

| パラメータ | 短縮形 | 説明 | デフォルト値 |
|-----------|--------|------|-------------|
| `-state` | `-s` | イシューの状態 (open, closed, all) | "" |
| `-sort` | `-so` | ソート項目 (created, updated, comments) | "" |
| `-direction` | `-d` | ソート方向 (asc, desc) | "" |
| `-per-page` | `-pp` | ページあたりの件数 | 30 |
| `-page` | `-p` | ページ番号 | 1 |
| `-issue-number` | `-in` | イシュー番号（特定イシュー取得用） | 0 |
| `-help` | `-h` | ヘルプを表示 | false |

## 使用例

### 1. 基本的なイシュー一覧取得

```bash
go run ./cmd/cli/github -o list-issues -ow octocat -r Hello-World
```

### 2. オープンなイシューのみを作成日時の降順で取得

```bash
go run ./cmd/cli/github -o list-issues -ow octocat -r Hello-World -s open -so created -d desc
```

### 3. クローズされたイシューを更新日時の昇順で取得

```bash
go run ./cmd/cli/github -o list-issues -ow octocat -r Hello-World -s closed -so updated -d asc
```

### 4. ページネーションを使用して2ページ目を取得（1ページあたり10件）

```bash
go run ./cmd/cli/github -o list-issues -ow octocat -r Hello-World -pp 10 -p 2
```

### 5. 特定のイシューを番号で取得

```bash
go run ./cmd/cli/github -o list-issues -ow octocat -r Hello-World -issue-number 123
# または短縮形
go run ./cmd/cli/github -o list-issues -ow octocat -r Hello-World -in 123
```

### 6. ヘルプの表示

```bash
go run ./cmd/cli/github -h
# または
go run ./cmd/cli/github -help
```

## GitHubトークンの取得方法

1. GitHubにログインし、Settings > Developer settings > Personal access tokens > Tokens (classic) に移動
2. "Generate new token (classic)" をクリック
3. 必要なスコープを選択（リポジトリの読み取りには `repo` スコープが必要）
4. トークンを生成し、`GITHUB_PERSONAL_ACCESS_TOKEN` に設定

## 出力形式

イシュー一覧はJSON形式で出力されます。各イシューには以下の情報が含まれます：

- id: イシューID
- number: イシュー番号
- title: タイトル
- body: 本文
- state: 状態 (open/closed)
- created_at: 作成日時
- updated_at: 更新日時
- user: 作成者情報
- labels: ラベル情報
- assignees: アサイン先情報

## エラーハンドリング

- 必須パラメータが不足している場合、エラーメッセージとヘルプが表示されます
- `GITHUB_PERSONAL_ACCESS_TOKEN` が未設定の場合、エラーメッセージとヘルプが表示されます
- GitHubトークンが無効な場合、認証エラーが表示されます
- 存在しないリポジトリを指定した場合、404エラーが表示されます

## ビルド方法

```bash
cd $HOME/devbox/cmd/cli/github
go build -o github main.go
```

## テスト実行

```bash
cd $HOME/devbox/internal/github
go test ./...
