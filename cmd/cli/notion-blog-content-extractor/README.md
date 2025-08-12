# Notion Blog Content Extractor

Markdownファイルからブログコンテンツを抽出するCLIツールです。特定のマーカーで区切られたコンテンツ部分のみを抽出し、新しいファイルとして保存します。

## 機能

- **コンテンツ抽出**: 特定のマーカー（`# Content` → 空行 → `## はじまり`）以降の内容を抽出
- **再帰的検索**: 指定されたディレクトリ内のすべての`.md`ファイルを再帰的に検索
- **自動ディレクトリ作成**: 出力先ディレクトリが存在しない場合は自動作成
- **エラーハンドリング**: 適切なエラーメッセージとステータスコードによる終了
- **ファイルフィルタリング**: マーカーを含まないファイルは自動的に除外

## インストール

```bash
# プロジェクトルートから
go build -o bin/notion-blog-content-extractor ./cmd/cli/notion-blog-content-extractor
```

## 使用方法

### 基本的な使用方法

```bash
./bin/notion-blog-content-extractor -src-dir=./blog-drafts -dest-dir=./extracted-content
```

### オプション

| オプション | 説明 | 必須 | 例 |
|-----------|------|------|-----|
| `-src-dir` | ソースディレクトリのパス（Markdownファイルが格納されているディレクトリ） | ✓ | `-src-dir=./blog-drafts` |
| `-dest-dir` | 出力先ディレクトリのパス（抽出したコンテンツを保存するディレクトリ） | ✓ | `-dest-dir=./extracted-content` |
| `-help` | ヘルプメッセージを表示 | | `-help` |

## 使用例

### 基本的な抽出

```bash
# カレントディレクトリの blog-drafts から extracted-content に抽出
./bin/notion-blog-content-extractor -src-dir=./blog-drafts -dest-dir=./extracted-content

# 絶対パスを使用した抽出
./bin/notion-blog-content-extractor -src-dir=/home/user/documents/blogs -dest-dir=/home/user/output
```

### ヘルプの表示

```bash
./bin/notion-blog-content-extractor -help
```

## マーカー形式

このツールは以下の形式のマーカーを検索します：

```markdown
# 記事のタイトル

この部分は抽出されません。

## 概要

この部分も抽出されません。

# Content

## はじまり

ここからが抽出対象のコンテンツです。

### セクション1

このセクションも含まれます。

- リスト項目1
- リスト項目2

### セクション2

```javascript
console.log("Hello, World!");
```

コードブロックも含まれます。

## まとめ

この部分もすべて抽出されます。
```

上記の例では、`# Content` → 空行 → `## はじまり` 以降のすべての内容が抽出されます。

## 出力フォーマット

### 成功時の出力

```
処理完了: 2件のファイルからコンテンツを抽出しました。
```

### マーカーを含むファイルが見つからない場合

```
指定されたディレクトリにコンテンツマーカーを含むMarkdownファイルが見つかりませんでした。
```

### 一部のファイルでエラーが発生した場合

```
処理完了: 1件のファイルからコンテンツを抽出しました。

エラーが発生したファイル:
ファイル broken-file.md: コンテンツの抽出に失敗しました: 指定されたマーカーが見つかりません
```

## エラーハンドリング

### ソースディレクトリが存在しない場合

```bash
./bin/notion-blog-content-extractor -src-dir=/nonexistent -dest-dir=./output
```

```
エラー: ソースディレクトリが存在しません: /nonexistent
exit status 1
```

### 必須パラメータの不足

```bash
./bin/notion-blog-content-extractor -src-dir=./blog-drafts
```

```
エラー: dest-dir パラメータは必須です
使用方法: ./bin/notion-blog-content-extractor [オプション]

Markdownファイルからブログコンテンツを抽出するツール
...
exit status 1
```

### ディレクトリではないパスを指定した場合

```bash
./bin/notion-blog-content-extractor -src-dir=./file.txt -dest-dir=./output
```

```
エラー: 指定されたパスはディレクトリではありません: ./file.txt
exit status 1
```

## 技術仕様

### アーキテクチャ

- **Clean Architecture**: ドメイン、ユースケース、インフラストラクチャの分離
- **SOLID原則**: インターフェースを活用した疎結合な設計
- **依存性注入**: テスト可能な設計のためのFileOperatorインターフェース

### ディレクトリ構造

```
internal/notion_blog_content_extractor/
├── config/          # 設定管理
│   ├── config.go    # 設定構造体とパーサー
│   └── flag_parser.go # フラグ解析インターフェース
└── usecases/        # ビジネスロジック
    └── services.go  # サービス層（コンテンツ抽出ロジック）
```

### 使用技術

- **Go**: プログラミング言語
- **標準ライブラリ**: `flag`, `fmt`, `os`, `path/filepath`, `regexp`, `strings`など
- **正規表現**: コンテンツマーカーの検索とパターンマッチング

### 処理フロー

1. **パラメータ検証**: 必須パラメータの存在確認
2. **ディレクトリ検証**: ソースディレクトリの存在と種別確認
3. **出力ディレクトリ作成**: 必要に応じて出力ディレクトリを作成
4. **ファイル検索**: 再帰的にMarkdownファイルを検索し、マーカーを含むファイルを特定
5. **コンテンツ抽出**: 各ファイルからマーカー以降の内容を抽出
6. **ファイル保存**: 抽出したコンテンツを新しいファイルとして保存
7. **結果報告**: 処理結果とエラー情報を出力

## 開発者向け情報

### ビルド

```bash
# 開発用ビルド
go build -o bin/notion-blog-content-extractor ./cmd/cli/notion-blog-content-extractor

# リリース用ビルド（複数プラットフォーム）
GOOS=linux GOARCH=amd64 go build -o bin/notion-blog-content-extractor-linux ./cmd/cli/notion-blog-content-extractor
GOOS=windows GOARCH=amd64 go build -o bin/notion-blog-content-extractor.exe ./cmd/cli/notion-blog-content-extractor
GOOS=darwin GOARCH=amd64 go build -o bin/notion-blog-content-extractor-mac ./cmd/cli/notion-blog-content-extractor
```

### テスト

```bash
# 単体テスト
go test ./internal/notion_blog_content_extractor/...

# カバレッジ付きテスト
go test -coverprofile=coverage.out ./internal/notion_blog_content_extractor/...
go tool cover -html=coverage.out -o coverage.html
```

### 開発時の実行

```bash
# go runを使用した開発時の実行
cd /path/to/devbox
go run ./cmd/cli/notion-blog-content-extractor -src-dir=./test-data -dest-dir=./output
```

## 実装の詳細

### 正規表現パターン

コンテンツマーカーの検索には以下の正規表現を使用しています：

```go
// 検索用パターン（マーカーの存在確認）
pattern := regexp.MustCompile(`# Content\s*\n\s*\n## はじまり\s*\n`)

// 抽出用パターン（マーカー以降のすべての内容を抽出）
pattern := regexp.MustCompile(`(?s)(# Content\s*\n\s*\n## はじまり\s*\n.*)`)
```

`(?s)`フラグにより、`.`が改行文字にもマッチするようになり、複数行にわたるコンテンツを正しく抽出できます。

### エラー処理

- **ファイル操作エラー**: 読み込み権限がない場合やディスク容量不足の場合
- **パターンマッチエラー**: マーカーが見つからない場合
- **ディレクトリ操作エラー**: 出力ディレクトリの作成に失敗した場合

## ライセンス

このプロジェクトのライセンスについては、プロジェクトルートのLICENSEファイルを参照してください。
