# Git Commit History Retriever

Gitリポジトリのコミット履歴を取得し、フィルタリングして表示するCLIツールです。

## 機能

- 指定されたGitディレクトリからコミット履歴を取得
- キーワードによるコミットメッセージの検索
- 日付範囲によるコミットの絞り込み
- グラフ形式での履歴表示
- **コミット詳細情報の自動取得**: 各コミットの変更ファイル統計を表示
- 「=== Commit History ===」と「=== Commit Details ===」見出し付きの整理された出力
- マージコミットやブランチ分岐にも対応

## インストール

```bash
# プロジェクトルートから
go build -o bin/git-commit-history-retriever ./cmd/cli/git-commit-history-retriever
```

## 使用方法

### 基本的な使用方法

```bash
./bin/git-commit-history-retriever -git-dir=/path/to/repo
```

出力例:
```
=== Commit History ===
* a1b2c3d - (HEAD -> main) feat: add new feature (2 hours ago) <Developer Name>
* d4e5f6g - fix: resolve bug in authentication (1 day ago) <Developer Name>
* h7i8j9k - refactor: improve code structure (3 days ago) <Developer Name>

=== Commit Details ===
commit a1b2c3d
Author: Developer Name <dev@example.com>
Date:   Fri Jul 11 08:00:00 2025 +0900

    feat: add new feature

 src/feature.go | 25 +++++++++++++++++++++++++
 src/main.go    |  3 +++
 2 files changed, 28 insertions(+)

commit d4e5f6g
Author: Developer Name <dev@example.com>
Date:   Thu Jul 10 08:00:00 2025 +0900

    fix: resolve bug in authentication

 src/auth.go | 10 +++++-----
 1 file changed, 5 insertions(+), 5 deletions(-)

commit h7i8j9k
Author: Developer Name <dev@example.com>
Date:   Wed Jul 9 08:00:00 2025 +0900

    refactor: improve code structure

 src/utils.go | 15 ++++++++-------
 1 file changed, 8 insertions(+), 7 deletions(-)
```

### オプション

| オプション | 説明 | 必須 | 例 |
|-----------|------|------|-----|
| `-git-dir` | 対象Gitディレクトリ | ✓ | `-git-dir=.` |
| `-keyword` | 検索キーワード | | `-keyword="feat:"` |
| `-since` | 開始年月日（YYYY-MM-DD形式） | | `-since="2025-01-01"` |
| `-until` | 終了年月日（YYYY-MM-DD形式） | | `-until="2025-01-31"` |

## 使用例

```bash
# 基本的なコミット履歴の取得
./bin/git-commit-history-retriever -git-dir=.

# キーワード検索
./bin/git-commit-history-retriever -git-dir=. -keyword="feat:"

# 日付範囲指定
./bin/git-commit-history-retriever -git-dir=. -since="2025-01-01" -until="2025-01-31"

# 複合条件
./bin/git-commit-history-retriever -git-dir=. -keyword="feat:" -since="2025-01-01" -until="2025-01-31"
```

## 出力フォーマット

出力は以下の形式で表示されます：

```
=== Commit History ===
* <短縮ハッシュ> - <ブランチ情報> <コミットメッセージ> (<相対時間>) <<作成者>>
```

条件に一致するコミットが見つからない場合：

```
=== Commit History ===
指定された条件に一致するコミットが見つかりませんでした。
```

## エラーハンドリング

### 無効なGitディレクトリ

```bash
./bin/git-commit-history-retriever -git-dir=/invalid/path
```

```
エラー: 指定されたディレクトリは有効なGitリポジトリではありません: /invalid/path
```

### 無効な日付フォーマット

```bash
./bin/git-commit-history-retriever -git-dir=. -since="2025/01/01"
```

```
エラー: --since の日付フォーマットが正しくありません。YYYY-MM-DD形式で入力してください
```

### 必須パラメータの不足

```bash
./bin/git-commit-history-retriever
```

```
エラー: --git-dir は必須パラメータです
使用方法: ./bin/git-commit-history-retriever [オプション]
...
```

## 技術仕様

### アーキテクチャ

- **Clean Architecture**: ドメイン、ユースケース、インフラストラクチャの分離
- **依存性注入**: テスタビリティの向上
- **インターフェース駆動**: モックによるテストの容易化

### ディレクトリ構造

```
internal/git_commit_history_retriever/
├── config/          # 設定管理
│   ├── config.go    # 設定構造体とパーサー
│   ├── flag_parser.go # フラグ解析
│   └── interfaces.go  # インターフェース定義
├── git/             # Git操作
│   ├── client.go    # Gitクライアント
│   └── executor.go  # コマンド実行
└── usecases/        # ビジネスロジック
    └── services.go  # サービス層
```

### 使用技術

- **Go**: プログラミング言語
- **標準ライブラリ**: `flag`, `os/exec`, `fmt`など
- **Git**: バージョン管理システム（外部依存）

## 開発者向け情報

### ビルド

```bash
# 開発用ビルド
go build -o bin/git-commit-history-retriever ./cmd/cli/git-commit-history-retriever

# リリース用ビルド（複数プラットフォーム）
GOOS=linux GOARCH=amd64 go build -o bin/git-commit-history-retriever-linux ./cmd/cli/git-commit-history-retriever
GOOS=windows GOARCH=amd64 go build -o bin/git-commit-history-retriever.exe ./cmd/cli/git-commit-history-retriever
GOOS=darwin GOARCH=amd64 go build -o bin/git-commit-history-retriever-mac ./cmd/cli/git-commit-history-retriever
```

### テスト

```bash
# 単体テスト
go test ./internal/git_commit_history_retriever/...

# カバレッジ付きテスト
go test -coverprofile=coverage.out ./internal/git_commit_history_retriever/...
go tool cover -html=coverage.out -o coverage.html
```

## ライセンス

このプロジェクトのライセンスについては、プロジェクトルートのLICENSEファイルを参照してください。
