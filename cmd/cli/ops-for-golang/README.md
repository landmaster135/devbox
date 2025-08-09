# Ops for Golang

Go開発でよく実行するコマンドを自動化するCLIツールです。テストカバレッジの取得、カバレッジ分析、CLIツールの実行を効率化し、出力フィルタリング機能も提供します。

## 機能

- **test-coverage**: `go test -cover ./...`でカレントディレクトリのテストカバレッジを取得
- **test-coverage-project**: プロジェクト全体のテストカバレッジ取得とHTMLレポート生成
- **coverage-func**: カバレッジファイルから関数レベルのカバレッジ情報を表示
- **run**: `go run`でCLIツールを実行
- **grepパターン**: 正規表現による出力フィルタリング機能
- **短縮オプション**: 全てのオプションに短縮形を提供
- **柔軟な出力制御**: 実行コマンドと進行状況の詳細表示

## インストール

```bash
# プロジェクトルートから
go build -o bin/ops-for-golang ./cmd/cli/ops-for-golang
```

## 使用方法

### 基本的な使用方法

```bash
./bin/ops-for-golang -ops test-coverage -directory /path/to/project
```

### オプション

| オプション | 短縮形 | 説明 | 必須 | 例 |
|-----------|--------|------|------|-----|
| `-ops` | `-o` | 操作タイプ (test-coverage, test-coverage-project, coverage-func, run) | ✓ | `-o test-coverage` |
| `-directory` | `-d` | 対象ディレクトリパス (test-coverage, test-coverage-project用) | | `-d /path/to/project` |
| `-execution_file` | `-e` | 実行ファイルパス (run用) | | `-e ./main.go` |
| `-parameters` | `-p` | 実行パラメータ (run用) | | `-p "-dry-run -token 'test'"` |
| `-coverage_file` | `-c` | カバレッジファイルパス (coverage-func用) | | `-c coverage.out` |
| `-grep_pattern` | `-g` | 出力フィルタリング用のgrepパターン (全操作共通) | | `-g "100.0%"` |
| `-help` | `-h` | ヘルプを表示 | | `-h` |

## 使用例

### テストカバレッジ取得

```bash
# 基本的なテストカバレッジ取得
./bin/ops-for-golang -ops test-coverage -directory /path/to/project
./bin/ops-for-golang -o test-coverage -d /path/to/project

# PASSした結果のみ表示
./bin/ops-for-golang -ops test-coverage -directory /path/to/project -grep_pattern "PASS"
./bin/ops-for-golang -o test-coverage -d /path/to/project -g "PASS"
```

### プロジェクト全体のテストカバレッジ取得

```bash
# プロジェクト全体のカバレッジ取得とHTMLレポート生成
./bin/ops-for-golang -ops test-coverage-project -directory /path/to/project
./bin/ops-for-golang -o test-coverage-project -d /path/to/project
```

### カバレッジファイルから関数情報取得

```bash
# カバレッジファイルから関数レベルの情報を表示
./bin/ops-for-golang -ops coverage-func -coverage_file coverage.out
./bin/ops-for-golang -o coverage-func -c coverage.out

# 100%カバレッジの関数のみ表示
./bin/ops-for-golang -ops coverage-func -coverage_file coverage.out -grep_pattern "100.0%"
./bin/ops-for-golang -o coverage-func -c coverage.out -g "100.0%"

# 特定のファイルの関数のみ表示
./bin/ops-for-golang -ops coverage-func -coverage_file coverage.out -grep_pattern "main.go"
```

### go run実行

```bash
# CLIツールの実行
./bin/ops-for-golang -ops run -execution_file ./main.go -parameters "-dry-run -token 'test_token'"
./bin/ops-for-golang -o run -e ./main.go -p "-dry-run -token 'test_token'"

# パラメータなしで実行
./bin/ops-for-golang -ops run -execution_file ./main.go
./bin/ops-for-golang -o run -e ./main.go
```

### grepパターンでフィルタリング

```bash
# エラーのみ表示
./bin/ops-for-golang -o test-coverage -d /path/to/project -g "FAIL"

# 特定のパッケージのみ表示
./bin/ops-for-golang -o coverage-func -c coverage.out -g "internal/.*"

# 正規表現パターンを使用
./bin/ops-for-golang -o coverage-func -c coverage.out -g "[0-9]+\\.0%"
```

## 出力フォーマット

### test-coverage操作の出力

```
テストカバレッジを実行中: /path/to/project
実行コマンド: go test -cover ./...
grepパターン: PASS

?       github.com/example/project/cmd/cli  [no test files]
ok      github.com/example/project/internal/config     0.002s  coverage: 85.7% of statements
PASS
ok      github.com/example/project/internal/usecases   0.003s  coverage: 92.3% of statements
PASS

テストカバレッジの実行が完了しました。
```

### test-coverage-project操作の出力

```
プロジェクト全体のテストカバレッジを実行中: /path/to/project
実行コマンド: go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out && go tool cover -html=coverage.out -o coverage.html

Step 1: テストカバレッジプロファイルを生成中...
ok      github.com/example/project/internal/config     0.002s  coverage: 85.7% of statements
ok      github.com/example/project/internal/usecases   0.003s  coverage: 92.3% of statements

Step 2: カバレッジ関数レポートを生成中...
github.com/example/project/internal/config/config.go:23:        NewConfig               100.0%
github.com/example/project/internal/usecases/services.go:22:    NewService              85.7%
total:                                                          (statements)            89.2%

Step 3: HTMLカバレッジレポートを生成中...

プロジェクト全体のテストカバレッジの実行が完了しました。
HTMLレポートが生成されました: /path/to/project/coverage.html
```

### coverage-func操作の出力

```
カバレッジファイルから関数情報を取得中: coverage.out
実行コマンド: go tool cover -func=coverage.out
grepパターン: 100.0%

github.com/example/project/internal/config/config.go:23:        NewConfig               100.0%
github.com/example/project/internal/usecases/services.go:45:    HandleRequest           100.0%
total:                                                          (statements)            89.2%

カバレッジ関数情報の取得が完了しました。
```

### run操作の出力

```
go runを実行中: ./main.go
パラメータ: -dry-run -token 'test_token'
実行コマンド: go run ./main.go -dry-run -token 'test_token'

[実行されたプログラムの出力]

go runの実行が完了しました。
```

## エラーハンドリング

### 無効な操作タイプ

```bash
./bin/ops-for-golang -ops invalid
```

```
エラー: 無効な操作タイプです: invalid
Go開発用操作CLIツール

使用方法:
  テストカバレッジ取得:
    ./bin/ops-for-golang -ops test-coverage -directory /path/to/project
...
```

### 必須パラメータの不足

```bash
./bin/ops-for-golang -ops test-coverage
```

```
エラー: test-coverage操作にはディレクトリパスが必要です
```

### ディレクトリが存在しない場合

```bash
./bin/ops-for-golang -ops test-coverage -directory /nonexistent/path
```

```
エラー: 指定されたディレクトリが存在しません: /nonexistent/path
```

### カバレッジファイルが存在しない場合

```bash
./bin/ops-for-golang -ops coverage-func -coverage_file nonexistent.out
```

```
エラー: 指定されたカバレッジファイルが存在しません: nonexistent.out
```

### 無効な正規表現パターン

```bash
./bin/ops-for-golang -ops coverage-func -coverage_file coverage.out -grep_pattern "["
```

```
エラー: 出力のフィルタリングに失敗しました: 無効な正規表現パターンです: error parsing regexp: missing closing ]: `[`
```

### 実行ファイルが存在しない場合

```bash
./bin/ops-for-golang -ops run -execution_file ./nonexistent.go
```

```
エラー: 指定された実行ファイルが存在しません: ./nonexistent.go
```

## 技術仕様

### アーキテクチャ

- **Clean Architecture**: ドメイン、ユースケース、インフラストラクチャの分離
- **SOLID原則**: インターフェースを活用した疎結合な設計
- **依存性注入**: テスタビリティを考慮した設計

### ディレクトリ構造

```
internal/ops_for_golang/
├── config/          # 設定管理
│   ├── config.go    # 設定構造体とパーサー
│   └── flag_parser.go # フラグ解析
└── usecases/        # ビジネスロジック
    └── services.go  # サービス層
```

### 使用技術

- **Go**: プログラミング言語
- **標準ライブラリ**: `flag`, `fmt`, `os/exec`, `regexp`など
- **テストフレームワーク**: Go標準のtestingパッケージ

## 開発者向け情報

### ビルド

```bash
# 開発用ビルド
go build -o bin/ops-for-golang ./cmd/cli/ops-for-golang

# リリース用ビルド（複数プラットフォーム）
GOOS=linux GOARCH=amd64 go build -o bin/ops-for-golang-linux ./cmd/cli/ops-for-golang
GOOS=windows GOARCH=amd64 go build -o bin/ops-for-golang.exe ./cmd/cli/ops-for-golang
GOOS=darwin GOARCH=amd64 go build -o bin/ops-for-golang-mac ./cmd/cli/ops-for-golang
```

### テスト

```bash
# 単体テスト
go test ./internal/ops_for_golang/...

# カバレッジ付きテスト
go test -coverprofile=coverage.out ./internal/ops_for_golang/...
go tool cover -html=coverage.out -o coverage.html

# このツール自体を使ってテスト（メタ的な使用）
./bin/ops-for-golang -ops test-coverage-project -directory ./internal/ops_for_golang
```

### 実行例（開発時）

```bash
# プロジェクトルートから直接実行
cd cmd/cli && go run ./ops-for-golang/ -ops test-coverage -directory ../../internal/ops_for_golang

# 短縮形での実行
cd cmd/cli && go run ./ops-for-golang/ -o coverage-func -c ../../coverage.out -g "100.0%"
```

## 関連ツール

このCLIツールは、Go開発ワークフローの効率化を目的として設計されています。以下のような場面で特に有用です：

- **CI/CDパイプライン**: テストカバレッジの自動取得
- **開発中の品質チェック**: 関数レベルでのカバレッジ分析
- **CLIツールの開発**: パラメータ付きでの実行テスト
- **レポート生成**: HTMLカバレッジレポートの自動生成

## ライセンス

このプロジェクトのライセンスについては、プロジェクトルートのLICENSEファイルを参照してください。
