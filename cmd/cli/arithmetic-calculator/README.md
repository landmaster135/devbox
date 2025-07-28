# Arithmetic Calculator

算術計算を行うCLIツールです。基本的な四則演算、配列の合計計算、ファイルの行数評価機能を提供します。

## 機能

- **基本計算**: 加算、減算、乗算、除算の四則演算
- **配列計算**: 複数の数値の合計計算
- **ファイル評価**: ファイルの行数が指定された閾値を超えているかの評価
- **柔軟な入力方式**: フラグ指定と位置引数の両方に対応
- **短縮オプション**: 全てのオプションに短縮形を提供
- **JSON出力**: ファイル評価結果を構造化されたJSON形式で出力

## インストール

```bash
# プロジェクトルートから
go build -o bin/arithmetic-calculator ./cmd/cli/arithmetic-calculator
```

## 使用方法

### 基本的な使用方法

```bash
./bin/arithmetic-calculator -operation add -x 10 -y 5
```

出力例:
```
10.00 + 5.00 = 15.00
```

### オプション

| オプション | 短縮形 | 説明 | 必須 | 例 |
|-----------|--------|------|------|-----|
| `-operation` | `-o` | 実行する操作 (add, subtract, multiply, divide, sum, evaluate_line_count) | ✓ | `-o add` |
| `-x` | | 第一オペランド (基本計算用) | | `-x 10` |
| `-y` | | 第二オペランド (基本計算用) | | `-y 5` |
| `-numbers` | `-n` | カンマ区切りの数値リスト (sum操作用) | | `-n 1,2,3,4,5` |
| `-file` | `-f` | 評価するファイルのパス (evaluate_line_count操作用) | | `-f /path/to/file.txt` |
| `-threshold` | `-t` | 行数の閾値 (evaluate_line_count操作用) | | `-t 100` |
| `-help` | `-h` | ヘルプを表示 | | `-h` |

## 使用例

### 基本計算

```bash
# 加算
./bin/arithmetic-calculator -operation add -x 10 -y 5
./bin/arithmetic-calculator -o add -x 10 -y 5
./bin/arithmetic-calculator -o add 10 5  # 位置引数

# 減算
./bin/arithmetic-calculator -operation subtract -x 10 -y 3

# 乗算
./bin/arithmetic-calculator -operation multiply -x 4 -y 7

# 除算
./bin/arithmetic-calculator -operation divide -x 20 -y 4
```

### 配列計算

```bash
# 複数の数値の合計
./bin/arithmetic-calculator -operation sum -numbers 1,2,3,4,5
./bin/arithmetic-calculator -o sum -n 1,2,3,4,5
```

### ファイル行数評価

```bash
# ファイルの行数が閾値を超えているかチェック
./bin/arithmetic-calculator -operation evaluate_line_count -file /path/to/file.txt -threshold 100
./bin/arithmetic-calculator -o evaluate_line_count -f /path/to/file.txt -t 100
./bin/arithmetic-calculator -o evaluate_line_count -f /path/to/file.txt 100  # 位置引数
```

## 出力フォーマット

### 基本計算の出力

```
10.00 + 5.00 = 15.00
20.00 - 8.00 = 12.00
4.00 * 7.00 = 28.00
20.00 / 4.00 = 5.00
```

### 配列計算の出力

```
sum([1 2 3 4 5]) = 15.00
```

### ファイル評価の出力

```json
{"file_path":"/path/to/file.txt","line_count":150,"threshold":100,"exceeds_threshold":true}
```

## エラーハンドリング

### 無効な操作タイプ

```bash
./bin/arithmetic-calculator -operation invalid
```

```
エラー: 無効な操作タイプです: invalid
```

### 必須パラメータの不足

```bash
./bin/arithmetic-calculator
```

```
エラー: 操作タイプが指定されていません
算術計算CLIツール

使用方法:
  基本計算:
    ./bin/arithmetic-calculator -operation add -x 10 -y 5
    ./bin/arithmetic-calculator -o add 10 5
...
```

### 無効な数値

```bash
./bin/arithmetic-calculator -o add -x invalid -y 5
```

```
エラー: 無効なx値です: invalid
```

### ゼロ除算

```bash
./bin/arithmetic-calculator -o divide -x 10 -y 0
```

```
エラー: ゼロで除算することはできません
```

### ファイルが存在しない場合

```bash
./bin/arithmetic-calculator -o evaluate_line_count -f /nonexistent/file.txt -t 100
```

```
エラー: ファイルを開けませんでした: /nonexistent/file.txt
```

## 技術仕様

### アーキテクチャ

- **Clean Architecture**: ドメイン、ユースケース、インフラストラクチャの分離
- **SOLID原則**: インターフェースを活用した疎結合な設計
- **TDD**: テスト駆動開発によるテストファースト実装

### ディレクトリ構造

```
internal/arithmetic_calculator/
├── config/          # 設定管理
│   ├── config.go    # 設定構造体とパーサー
│   ├── flag_parser.go # フラグ解析
│   └── interfaces.go  # インターフェース定義
└── usecases/        # ビジネスロジック
    └── services.go  # サービス層
```

### 使用技術

- **Go**: プログラミング言語
- **標準ライブラリ**: `flag`, `fmt`, `strconv`, `strings`など
- **テストフレームワーク**: Go標準のtestingパッケージ

## 開発者向け情報

### ビルド

```bash
# 開発用ビルド
go build -o bin/arithmetic-calculator ./cmd/cli/arithmetic-calculator

# リリース用ビルド（複数プラットフォーム）
GOOS=linux GOARCH=amd64 go build -o bin/arithmetic-calculator-linux ./cmd/cli/arithmetic-calculator
GOOS=windows GOARCH=amd64 go build -o bin/arithmetic-calculator.exe ./cmd/cli/arithmetic-calculator
GOOS=darwin GOARCH=amd64 go build -o bin/arithmetic-calculator-mac ./cmd/cli/arithmetic-calculator
```

### テスト

```bash
# 単体テスト
go test ./internal/arithmetic_calculator/...

# カバレッジ付きテスト
go test -coverprofile=coverage.out ./internal/arithmetic_calculator/...
go tool cover -html=coverage.out -o coverage.html
```

### 関連ツール

このCLIツールと同じサービスを使用するMCPツールも利用可能です：
- `devbox/cmd/mcp/arithmetic_calculator/mcp.go`

MCPツールを使用することで、Model Context Protocol経由で同じ計算機能にアクセスできます。

## ライセンス

このプロジェクトのライセンスについては、プロジェクトルートのLICENSEファイルを参照してください。
