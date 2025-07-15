# Context7 CLI

Context7 CLIは、最新のライブラリドキュメントを取得するためのコマンドラインツールです。Context7 APIを使用して、ライブラリの検索とドキュメントの取得を行います。

## 機能

- **ライブラリ検索**: ライブラリ名からContext7互換のライブラリIDを検索
- **ドキュメント取得**: ライブラリIDを使用して最新のドキュメントを取得
- **トピック指定**: 特定のトピックに焦点を当てたドキュメント取得
- **トークン数制御**: 取得するドキュメントの最大トークン数を指定

## インストール

```bash
# リポジトリをクローン
git clone https://github.com/landmaster135/devbox.git
cd devbox

# ビルド
go build -o context7 ./cmd/cli/context7/
```

## 使用方法

### 基本的な使用方法

```bash
context7 <command> [arguments]
```

### 利用可能なコマンド

#### `search` - ライブラリ検索

ライブラリ名からContext7互換のライブラリIDを検索します。

```bash
context7 search <library_name>
```

**例:**
```bash
# Reactライブラリを検索
context7 search react

# Next.jsライブラリを検索
context7 search "next.js"
```

#### `docs` - ドキュメント取得

ライブラリIDを使用してドキュメントを取得します。

```bash
context7 docs <library_id> [flags]
```

**フラグ:**
- `-topic string`: 特定のトピックに焦点を当てる（例: hooks, routing）
- `-tokens int`: 取得する最大トークン数（デフォルト: 10000）

**例:**
```bash
# Reactの基本ドキュメントを取得
context7 docs /context7/react_dev

# Next.jsのルーティングに関するドキュメントを取得
context7 docs /vercel/next.js -topic=routing

# MongoDBのドキュメントを5000トークンで取得
context7 docs /mongodb/docs -tokens=5000

# Supabaseの認証に関するドキュメントを8000トークンで取得
context7 docs /supabase/supabase -topic=auth -tokens=8000
```

#### `help` - ヘルプ表示

使用方法とコマンドの詳細を表示します。

```bash
context7 help
```

## 使用例

```bash
# Reactライブラリを検索
./context7 search react

# Vue.jsライブラリを検索
./context7 search vue

# TypeScriptライブラリを検索
./context7 search typescript

# Reactの基本ドキュメントを取得
./context7 docs /context7/react_dev

# 特定のトピックに焦点を当てたドキュメント取得
./context7 docs /context7/react_dev -topic=hooks

# トークン数を制限してドキュメント取得
./context7 docs /context7/react_dev -tokens=2000

# トピックとトークン数を両方指定
./context7 docs /vercel/next.js -topic=routing -tokens=3000

# 1. 使いたいライブラリを検索
./context7 search "react router"

# 2. 検索結果からライブラリIDを確認
# 例: /context7/reactrouter が見つかった場合

# 3. そのライブラリのドキュメントを取得
./context7 docs /context7/reactrouter

# 4. 特定の機能について詳しく調べる
./context7 docs /context7/reactrouter -topic=navigation -tokens=5000

# 存在しないライブラリIDを指定した場合
./context7 docs /nonexistent/library
# 出力: 指定されたライブラリIDのドキュメントが見つかりません。

# 無効なライブラリID形式を指定した場合
./context7 docs invalid-format
# 出力: ライブラリID検証エラー: ライブラリIDの形式が正しくありません。期待される形式: org/project または org/project/version
```

## 実行例

### 1. Reactライブラリの検索

```bash
$ context7 search react
検索結果:

1. React
   ID: /context7/react_dev
   説明: React is a JavaScript library for building user interfaces...
   コードスニペット数: 2053
   信頼スコア: 10.0
   最終更新: 2025-07-15T05:47:06.332Z
   状態: finalized

2. React Flow
   ID: /xyflow/xyflow
   説明: React Flow | Svelte Flow - Powerful open source libraries...
   コードスニペット数: 21
   信頼スコア: 9.5
   スター数: 29143
   最終更新: 2025-04-21T00:00:00.000Z
   状態: finalized
...
```

### 2. Reactドキュメントの取得

```bash
$ context7 docs /context7/react_dev -tokens=1000
=== /context7/react_dev のドキュメント ===

トークン数: 1000

TITLE: React DOM API Reference (react-dom@19.1)
DESCRIPTION: This section outlines the various APIs, hooks, and components...
SOURCE: https://react.dev/reference/react-dom/components/style

LANGUAGE: APIDOC
CODE:
```
react-dom@19.1 API Reference:

Hooks:
  - useFormStatus: Provides status information about the nearest pending form submission.

Components:
  - Common HTML Elements (e.g., <div>, <form>, <input>...
```
...
```

## アーキテクチャ

Context7 CLIは以下のコンポーネントで構成されています：

### ディレクトリ構造

```
cmd/cli/context7/
├── main.go                          # CLIエントリーポイント
└── README.md                        # このファイル

internal/context7/
├── domain/
│   └── models/
│       └── context7.go             # ドメインモデル
├── interfaces/
│   └── http_client.go              # HTTPクライアントインターフェース
└── usecases/
    ├── services.go                 # Context7サービス
    └── services_test.go            # テストコード
```

### 主要コンポーネント

#### 1. Domain Models (`internal/context7/domain/models/`)

- **SearchResult**: ライブラリ検索結果を表現
- **SearchResponse**: 検索APIレスポンス
- **DocOptions**: ドキュメント取得オプション
- **Context7Request/Response**: API通信用モデル

#### 2. Interfaces (`internal/context7/interfaces/`)

- **HTTPClient**: HTTP通信の抽象化インターフェース
- **DefaultHTTPClient**: 標準HTTP実装

#### 3. Use Cases (`internal/context7/usecases/`)

- **Context7Service**: Context7 APIとの通信を担当
  - `ResolveLibraryID()`: ライブラリ検索
  - `GetLibraryDocs()`: ドキュメント取得
  - `FormatSearchResults()`: 検索結果の整形
  - `ValidateLibraryID()`: ライブラリID検証

#### 4. CLI (`cmd/cli/context7/`)

- **main.go**: コマンドライン引数の解析とコマンド実行
  - `handleSearchCommand()`: 検索コマンド処理
  - `handleDocsCommand()`: ドキュメント取得コマンド処理
  - `printUsage()`: ヘルプ表示

## API仕様

Context7 CLIは以下のContext7 APIエンドポイントを使用します：

### 検索API

```
GET https://context7.com/api/v1/search?query={library_name}
```

### ドキュメント取得API

```
GET https://context7.com/api/v1/{library_id}?tokens={tokens}&topic={topic}&type=txt
```

## エラーハンドリング

- **レート制限**: 429ステータスコードを適切に処理し、ユーザーに分かりやすいメッセージを表示
- **404エラー**: ライブラリが見つからない場合の適切なメッセージ表示
- **ネットワークエラー**: 接続エラーやタイムアウトの処理
- **JSONパースエラー**: APIレスポンスの形式エラー処理

## テスト

包括的なテストスイートが含まれています：

```bash
# テスト実行
go test ./internal/context7/usecases/ -v

# カバレッジ確認
go test -coverprofile=coverage.out ./internal/context7/usecases/
go tool cover -html=coverage.out
```

### テストカバレッジ

- **Context7Service**: 全メソッドの正常系・異常系テスト
- **HTTPクライアント**: モックを使用したHTTP通信テスト
- **エラーハンドリング**: 各種エラーケースのテスト
- **データ整形**: 検索結果とドキュメントの整形テスト

## 開発

### 前提条件

- Go 1.21以上
- インターネット接続（Context7 APIアクセス用）

### 開発環境セットアップ

```bash
# 依存関係のインストール
go mod tidy

# テスト実行
go test ./...

# ビルド
go build -o context7 ./cmd/cli/context7/
```

### コーディング規約

- **命名規則**: Go標準の命名規則に従う
- **エラーハンドリング**: 適切なエラーラッピングとメッセージ
- **テスト**: 新機能には必ずテストを追加
- **ドキュメント**: 公開関数・メソッドにはGoDocコメントを記述

## ライセンス

このプロジェクトのライセンスについては、リポジトリのルートディレクトリのLICENSEファイルを参照してください。

## 貢献

バグ報告や機能要求は、GitHubのIssueで受け付けています。プルリクエストも歓迎します。

## 関連リンク

- [Context7公式サイト](https://context7.com/)
- [Context7 API ドキュメント](https://context7.com/api)
- [devboxリポジトリ](https://github.com/landmaster135/devbox)
