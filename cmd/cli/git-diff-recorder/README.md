# Git Diff Recorder

Git差分を記録・読み取りするCLIツールです。リポジトリの差分情報を構造化されたフォーマットでファイルに出力し、後で読み取ることができます。

## 機能

- **記録モード**: Git差分の記録（ステージング済み/全差分）
- **読み取りモード**: 保存されたdiffファイルから詳細差分を表示
- **出力モード**: 指定Gitディレクトリから直接差分を取得・出力
- リポジトリ情報の記録（リポジトリ名、ブランチ、コミットハッシュ）
- ファイル変更統計の記録
- 構造化された出力フォーマット（将来のコミットメッセージ生成に対応）

## 使用方法

### 記録モード（既存機能）

```bash
# 全ての差分を記録
go run ./cmd/cli/git-diff-recorder --output-dir /path/to/output

# ステージング済み差分のみを記録
go run ./cmd/cli/git-diff-recorder --output-dir /path/to/output --staged-only
```

### 読み取りモード

```bash
# 指定リポジトリの最新diffファイルから詳細差分を表示
go run ./cmd/cli/git-diff-recorder --read-mode --source-dir /path/to/diffs --repository repo-name
```

### 出力モード

```bash
# 指定Gitディレクトリから直接差分を取得
go run ./cmd/cli/git-diff-recorder --output-mode --git-dir /path/to/git/repository

# ステージング済み差分のみを取得
go run ./cmd/cli/git-diff-recorder --output-mode --git-dir /path/to/git/repository --staged-only
```

## 使用例
```bash
# 全ての差分を記録（現在のディレクトリ）
go run ./cmd/cli/git-diff-recorder --output-dir /tmp/diffs

# ステージング済み差分のみを記録
go run ./cmd/cli/git-diff-recorder --output-dir /tmp/diffs --staged-only

# 指定Gitディレクトリの差分を記録
go run ./cmd/cli/git-diff-recorder --output-dir /tmp/diffs --git-dir /home/user/my-project

# devboxリポジトリの最新差分を表示
go run ./cmd/cli/git-diff-recorder --read-mode --source-dir /tmp/diffs --repository devbox

# 別のリポジトリの差分を表示
go run ./cmd/cli/git-diff-recorder --read-mode --source-dir /home/user/git-diffs --repository my-project

# 指定Gitディレクトリから直接差分を取得
go run ./cmd/cli/git-diff-recorder --output-mode --git-dir /home/user/my-project

# 指定Gitディレクトリのステージング済み差分のみを取得
go run ./cmd/cli/git-diff-recorder --output-mode --git-dir /home/user/my-project --staged-only
```

### パラメータ

#### 記録モード
- `--output-dir` (必須): 出力先ディレクトリ
- `--git-dir` (オプション): 対象Gitディレクトリ（未指定時は現在のディレクトリ）
- `--staged-only` (オプション): ステージング済み差分のみ記録 (デフォルト: false)

#### 読み取りモード
- `--read-mode` (必須): 読み取りモードを有効にする
- `--source-dir` (必須): 読み取り対象のディレクトリ
- `--repository` (必須): 対象リポジトリ名

#### 出力モード
- `--output-mode` (必須): 出力モードを有効にする
- `--git-dir` (必須): 対象Gitディレクトリ
- `--staged-only` (オプション): ステージング済み差分のみ取得 (デフォルト: false)

### 出力形式

出力ファイルは以下の形式で作成されます：

```
{output-dir}/
└── {repository-name}/
    └── diff_yyyyMMddhhmmss.txt
```

### 出力内容の例

```
=== Git Diff Record ===
Generated at: 2025-01-07 12:30:45
Repository: devbox
Branch: webp
Latest commit: a1b2c3d4
Options: --staged-only=false

=== File Changes Summary ===
Modified files: 2
New files: 4
Deleted files: 0

=== New Files ===
cmd/cli/git-diff-recorder/main.go
internal/git_diff_recorder/usecases/services.go
internal/git_diff_recorder/config/config.go
internal/git_diff_recorder/git/client.go

=== Detailed Diff ===
diff --git a/.github/workflows/test_integration.yml b/.github/workflows/test_integration.yml
index 2561e13..e215710 100644
--- a/.github/workflows/test_integration.yml
+++ b/.github/workflows/test_integration.yml
...
```

## ビルド方法

```bash
# devboxディレクトリで実行
cd devbox
go build -o bin/git-diff-recorder ./cmd/cli/git-diff-recorder
```

## テスト実行

```bash
# devboxディレクトリで実行
cd devbox
go test ./internal/git_diff_recorder/usecases/... -v
```

## 構造化された出力について

出力ファイルには以下の見出しが含まれており、将来のコミットメッセージ生成機能で活用できます：

- `=== Git Diff Record ===`: 基本情報セクション
- `=== File Changes Summary ===`: ファイル変更統計セクション
- `=== New Files ===`: 新規ファイル一覧セクション
- `=== Deleted Files ===`: 削除ファイル一覧セクション
- `=== Detailed Diff ===`: 詳細差分セクション

これらの見出しは定数として定義されており、パーサーで簡単に識別できます。

## 注意事項

- このツールはGitリポジトリ内で実行する必要があります
- 出力ディレクトリが存在しない場合は自動的に作成されます
- 同じ秒に複数回実行すると、ファイルが上書きされる可能性があります
