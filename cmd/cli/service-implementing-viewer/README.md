# Service Implementing Viewer

複数のディレクトリ内にあるサービスの実装状況を表形式で表示するCLIツールです。

## 概要

このツールは、指定されたルートディレクトリ配下の複数の対象ディレクトリをスキャンし、各ディレクトリに実装されているサービスの一覧を表形式で表示します。サービス名は正規化され（「_」と「-」を統一）、昇順でソートされます。

## 機能

- **ディレクトリスキャン**: 指定されたルートディレクトリ配下の対象ディレクトリを再帰的にスキャン
- **サービス名正規化**: 「_」と「-」を同じものとして扱い、統一された形式で表示
- **表形式出力**: Markdown形式の表でサービス実装状況を視覚的に表示
- **統計情報表示**: 各ディレクトリの実装数、組み合わせパターンなどの詳細な統計情報を表示
- **ソート機能**: サービス名を昇順でソート
- **存在チェック**: 各サービスが各対象ディレクトリに存在するかを✅/❌で表示

## インストール

```bash
# devboxプロジェクトのルートディレクトリで実行
go build -o bin/service-implementing-viewer cmd/cli/service-implementing-viewer/main.go
```

## 使用方法

### 基本的な使用方法

```bash
go run cmd/cli/service-implementing-viewer/main.go -root-dir=<ルートディレクトリ> -target-dirs=<対象ディレクトリ> -operation=output
```

### オプション

- `-root-dir` (必須): スキャンするルートディレクトリのパス
- `-target-dirs` (必須): 対象ディレクトリ名をカンマ区切りで指定
- `-operation` (必須): 実行するオペレーション。既存機能は `output` を指定

## 使用例

```bash
# cliとmcpディレクトリをスキャン
go run cmd/cli/service-implementing-viewer/main.go -operation=output -root-dir=/home/user/devbox/cmd -target-dirs=cli,mcp

# 全ディレクトリをスキャン
go run cmd/cli/service-implementing-viewer/main.go -operation=output -root-dir=/home/user/devbox/cmd -target-dirs=cli,mcp,grpc/handlers,http/handlers

# 複数のディレクトリをスキャン
go run cmd/cli/service-implementing-viewer/main.go -operation=output -root-dir=/path/to/project -target-dirs="cli,mcp,powershell"
```

## 出力例

### 表形式出力
```
| service                        | cli | mcp | grpc/handlers | http/handlers |
| :----------------------------: | :-: | :-: | :-----------: | :-----------: |
| arithmetic-calculator          | ✅  | ✅  | ❌️           | ❌️           |
| base64-extractor               | ✅  | ❌️  | ❌️           | ❌️           |
| brave-search                   | ❌️  | ✅  | ❌️           | ❌️           |
| context7                       | ✅  | ✅  | ❌️           | ❌️           |
| git-commit-history-retriever   | ✅  | ✅  | ❌️           | ❌️           |
| github                         | ❌️  | ✅  | ❌️           | ❌️           |
| http-request                   | ✅  | ✅  | ❌️           | ❌️           |
| weather-notificator            | ✅  | ✅  | ✅           | ✅           |
```

### 統計情報
```
## 統計情報

- **総サービス数**: 8
- **CLIツール実装数**: 5
- **MCPツール実装数**: 6
- **gRPCハンドラ実装数**: 1
- **HTTPハンドラ実装数**: 1
- **CLIのみ実装**: 1
- **MCPのみ実装**: 2
- **gRPCハンドラのみ実装**: 0
- **HTTPハンドラのみ実装**: 0
- **CLI+MCP両方実装**: 4
- **全て実装済み**: 1
```

## アーキテクチャ

### ディレクトリ構造

```
internal/service_implementing_viewer/
├── config/                    # 設定管理
│   ├── config.go             # 設定構造体とバリデーション
│   ├── config_test.go        # 設定のテストコード
│   ├── flag_parser.go        # フラグパーサー実装
│   └── interfaces.go         # インターフェース定義
└── usecases/                 # ビジネスロジック
    ├── services.go           # サービス実装状況確認ロジック
    └── services_test.go      # ビジネスロジックのテストコード
```

## 主要コンポーネント

### Config パッケージ
- **Config**: CLI設定を保持する構造体
- **ConfigParser**: コマンドライン引数の解析を行う
- **FlagParser**: フラグ解析のインターフェース
- **OSArgs**: OS引数取得のインターフェース

### Usecases パッケージ
- **ServiceImplementingViewerService**: メインのサービスクラス
- **ServiceStatus**: サービスの実装状況を表す構造体
- **ServiceStatistics**: サービス実装の統計情報を表す構造体

## 主要メソッド

### ServiceImplementingViewerService
- `GetServiceImplementingStatus()`: サービス実装状況を取得し、表形式と統計情報で返す
- `getServicesInDirectory()`: 指定されたディレクトリ内のサービス名を取得
- `normalizeServiceName()`: サービス名を正規化（「_」→「-」）
- `isServiceImplementedInDirectory()`: サービスがディレクトリに実装されているかチェック
- `formatAsTable()`: 結果をMarkdown表形式でフォーマット
- `calculateStatistics()`: 統計情報を計算（各ディレクトリの実装数、組み合わせパターンなど）

## テスト

### テスト実行

```bash
# 全テスト実行
go test ./internal/service_implementing_viewer/...

# カバレッジ付きテスト実行
go test -coverprofile=coverage.out ./internal/service_implementing_viewer/...

# カバレッジレポート生成
go tool cover -html=coverage.out -o coverage.html
```

### テストカバレッジ

- **usecases**: 95.9%
- **config**: 56.2%

### テスト内容

Config パッケージ
- 設定の正常系・異常系テスト
- フラグ解析のテスト
- バリデーションのテスト
- モックを使用した依存関係のテスト

Usecases パッケージ
- サービス作成のテスト
- ディレクトリスキャンのテスト
- サービス名正規化のテスト
- 表形式フォーマットのテスト
- 一時ディレクトリを使用した統合テスト

## エラーハンドリング

### 一般的なエラー

1. **必須パラメータ不足**
   ```
   エラー: 設定の初期化に失敗しました: --root-dir は必須パラメータです
   ```

2. **ディレクトリ読み取りエラー**
   ```
   エラー: ディレクトリ /path/to/dir の読み取りに失敗しました: permission denied
   ```

### エラー対処法

- 必須パラメータ（`-root-dir`, `-target-dirs`）が正しく指定されているか確認
- 指定されたディレクトリが存在し、読み取り権限があるか確認
- パスが正しく指定されているか確認

## 開発

### 依存関係

- Go 1.19以上
- 標準ライブラリのみ使用（外部依存なし）

### コーディング規約

- Go標準のコーディング規約に準拠
- TDD（テスト駆動開発）アプローチを採用
- SOLID原則に基づいた設計
- インターフェースを活用したテスタブルな設計

### 貢献

1. フォークしてブランチを作成
2. 変更を実装
3. テストを追加・実行
4. プルリクエストを作成

## ライセンス

このプロジェクトのライセンスに従います。

## 更新履歴

### v1.1.0 (2025-09-15)
- 統計情報表示機能を追加
- gRPC/HTTPハンドラ対応を追加
- 4つのディレクトリ（cli, mcp, grpc/handlers, http/handlers）に対応
- 詳細な統計情報（実装数、組み合わせパターンなど）を表示
- テストケースを拡張して4つのディレクトリに対応

### v1.0.0 (2025-07-28)
- 初回リリース
- 基本的なサービス実装状況表示機能
- 表形式出力機能
- 包括的なテストスイート
