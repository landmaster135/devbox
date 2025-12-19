# Arithmetic Calculator

算術計算を行うCLIツールです。基本的な四則演算から高度な数学関数まで、幅広い計算機能を提供します。

## 機能

### 基本機能
- **基本計算**: 加算、減算、乗算、除算の四則演算
- **配列計算**: 複数の数値の合計計算
- **ファイル評価**: ファイルの行数が指定された閾値を超えているかの評価
- **API料金抽出**: ファイルまたはテキストから「API料金が○円掛かった」形式の文字列を抽出し合計を計算

### 高度な数学機能
- **べき乗計算**: 指定された底と指数でべき乗を計算
- **平方根計算**: 指定された数値の平方根を計算
- **階乗計算**: 指定された整数の階乗を計算
- **三角関数**: sin, cos, tan関数（度数・ラジアン対応）
- **数式評価**: 複雑な数式を安全に評価
- **数学定数**: π, e, τなどの数学定数を提供

### その他の特徴
- **柔軟な入力方式**: フラグ指定と位置引数の両方に対応
- **短縮オプション**: 全てのオプションに短縮形を提供
- **JSON出力**: ファイル評価結果を構造化されたJSON形式で出力
- **安全性**: 数式評価時の危険なコード実行を防止

## インストール

```bash
# プロジェクトルートから
go build -o bin/arithmetic-calculator ./cmd/cli/arithmetic-calculator
```

## 使用方法

### 基本的な使用方法

```bash
go run ./cmd/cli/arithmetic-calculator -operation add -x 10 -y 5
```

出力例:
```
10.00 + 5.00 = 15.00
```

### オプション

| オプション | 短縮形 | 説明 | 必須 | 例 |
|-----------|--------|------|------|-----|
| `-operation` | `-o` | 実行する操作 | * | `-o add` |
| `-x` | | 第一オペランド (基本計算用) | | `-x 10` |
| `-y` | | 第二オペランド (基本計算用) | | `-y 5` |
| `-numbers` | `-nums` | カンマ区切りの数値リスト (sum操作用) | | `-nums 1,2,3,4,5` |
| `-file` | `-f` | 評価するファイルのパス | | `-f /path/to/file.txt` |
| `-threshold` | `-t` | 行数の閾値 (evaluate_line_count操作用) | | `-t 100` |
| `-text-input` | `-ti` | テキスト入力 (parse-api-cost操作用) | | `-ti "API料金が100円掛かった。"` |
| `-base` | `-b` | べき乗の底 (power操作用) | | `-base 2` |
| `-exponent` | `-exp` | べき乗の指数 (power操作用) | | `-exp 8` |
| `-number` | | 単一の数値 (square_root操作用) | | `-number 64` |
| `-n` | | 階乗の数値 (factorial操作用) | | `-n 5` |
| `-function` | `-func` | 三角関数の種類 (sin, cos, tan) | | `-func sin` |
| `-angle` | `-a` | 角度 (trigonometry操作用) | | `-angle 90` |
| `-unit` | `-u` | 角度の単位 (radians, degrees) | | `-unit degrees` |
| `-expression` | `-expr` | 評価する数式 (calculate操作用) | | `-expr "2+3*4"` |
| `-help` | `-h` | ヘルプを表示 | | `-h` |

### 利用可能な操作

| 操作 | 説明 | 必要なパラメータ |
|------|------|------------------|
| `add` | 加算 | `-x`, `-y` |
| `subtract` | 減算 | `-x`, `-y` |
| `multiply` | 乗算 | `-x`, `-y` |
| `divide` | 除算 | `-x`, `-y` |
| `sum` | 配列の合計 | `-numbers` |
| `power` | べき乗 | `-base`, `-exponent` |
| `square_root` | 平方根 | `-number` |
| `factorial` | 階乗 | `-n` |
| `trigonometry` | 三角関数 | `-function`, `-angle`, `-unit` |
| `calculate` | 数式評価 | `-expression` |
| `get_constants` | 数学定数一覧 | なし |
| `evaluate_line_count` | ファイル行数評価 | `-file`, `-threshold` |
| `parse-api-cost` | API料金抽出 | `-file` または `-text-input` |

## 使用例

### 基本計算

```bash
# 加算
go run ./cmd/cli/arithmetic-calculator -operation add -x 10 -y 5
go run ./cmd/cli/arithmetic-calculator -o add -x 10 -y 5
go run ./cmd/cli/arithmetic-calculator -o add 10 5  # 位置引数

# 減算
go run ./cmd/cli/arithmetic-calculator -operation subtract -x 10 -y 3

# 乗算
go run ./cmd/cli/arithmetic-calculator -operation multiply -x 4 -y 7

# 除算
go run ./cmd/cli/arithmetic-calculator -operation divide -x 20 -y 4
```

### 高度な数学計算

```bash
# べき乗計算
go run ./cmd/cli/arithmetic-calculator -operation power -base 2 -exponent 8
go run ./cmd/cli/arithmetic-calculator -o power -b 2 -exp 8

# 平方根計算
go run ./cmd/cli/arithmetic-calculator -operation square_root -number 64
go run ./cmd/cli/arithmetic-calculator -o square_root -number 64

# 階乗計算
go run ./cmd/cli/arithmetic-calculator -operation factorial -n 5
go run ./cmd/cli/arithmetic-calculator -o factorial -n 5

# 三角関数計算
go run ./cmd/cli/arithmetic-calculator -operation trigonometry -function sin -angle 90 -unit degrees
go run ./cmd/cli/arithmetic-calculator -o trigonometry -func cos -a 0 -u radians
go run ./cmd/cli/arithmetic-calculator -o trigonometry -func tan -a 45 -u degrees

# 数式評価
go run ./cmd/cli/arithmetic-calculator -operation calculate -expression "2+3*4"
go run ./cmd/cli/arithmetic-calculator -o calculate -expr "sqrt(16)+2**3"
go run ./cmd/cli/arithmetic-calculator -o calculate -expr "sin(pi/2)"

# 数学定数の表示
go run ./cmd/cli/arithmetic-calculator -operation get_constants
go run ./cmd/cli/arithmetic-calculator -o get_constants
```

### 配列計算

```bash
# 複数の数値の合計
go run ./cmd/cli/arithmetic-calculator -operation sum -numbers 1,2,3,4,5
go run ./cmd/cli/arithmetic-calculator -o sum -nums 1,2,3,4,5
```

### ファイル行数評価

```bash
# ファイルの行数が閾値を超えているかチェック
go run ./cmd/cli/arithmetic-calculator -operation evaluate_line_count -file /path/to/file.txt -threshold 100
go run ./cmd/cli/arithmetic-calculator -o evaluate_line_count -f /path/to/file.txt -t 100
go run ./cmd/cli/arithmetic-calculator -o evaluate_line_count -f /path/to/file.txt 100  # 位置引数
```

### API料金抽出

```bash
# ファイルからAPI料金を抽出
go run ./cmd/cli/arithmetic-calculator -operation parse-api-cost -file /path/to/api_log.txt
go run ./cmd/cli/arithmetic-calculator -o parse-api-cost -f /path/to/api_log.md

# テキストからAPI料金を抽出
go run ./cmd/cli/arithmetic-calculator -operation parse-api-cost -text-input "API料金が100円掛かった。別のAPI料金が200円掛かった。"
go run ./cmd/cli/arithmetic-calculator -o parse-api-cost -ti "API料金が150円掛かった。"
```

## 出力フォーマット

### 基本計算の出力

```
10.00 + 5.00 = 15.00
20.00 - 8.00 = 12.00
4.00 * 7.00 = 28.00
20.00 / 4.00 = 5.00
```

### 高度な数学計算の出力

```
2.00^8.00 = 256.00
√64.00 = 8.00
5! = 120
sin(90.00 degrees) = 1.000000
2+3*4 = 14.00
```

### 数学定数の出力

```
利用可能な数学定数:
pi = 3.141593
e = 2.718282
tau = 6.283185
```

### 配列計算の出力

```
sum([1 2 3 4 5]) = 15.00
```

### ファイル評価の出力

```json
{
  "is_greater": true,
  "line_count": 150,
  "threshold": 100,
  "file_path": "/path/to/file.txt",
  "description": "ファイル '/path/to/file.txt' の行数は 150 行で、閾値 100 より大きいです。"
}
```

### API料金抽出の出力

```
抽出されたAPI料金の合計: 300円
```

## エラーハンドリング

### 無効な操作タイプ

```bash
go run ./cmd/cli/arithmetic-calculator -operation invalid
```

```
エラー: 無効な操作タイプです: invalid
```

### 必須パラメータの不足

```bash
go run ./cmd/cli/arithmetic-calculator -o trigonometry -angle 90 -unit degrees
```

```
エラー: 三角関数の種類が指定されていません
```

### 無効な数値

```bash
go run ./cmd/cli/arithmetic-calculator -o add -x invalid -y 5
```

```
エラー: 無効なx値です: invalid
```

### 数学的エラー

```bash
# ゼロ除算
go run ./cmd/cli/arithmetic-calculator -o divide -x 10 -y 0
```

```
エラー: division by zero is not allowed
```

```bash
# 負数の平方根
go run ./cmd/cli/arithmetic-calculator -o square_root -number -4
```

```
エラー: 負数の平方根は計算できません
```

```bash
# 負数の階乗
go run ./cmd/cli/arithmetic-calculator -o factorial -n -5
```

```
エラー: 階乗は負数では定義されていません
```

### 危険な数式

```bash
go run ./cmd/cli/arithmetic-calculator -o calculate -expr "import os"
```

```
エラー: 危険なパターンが検出されました: import
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

### サービス構成

- **CalculatorService**: 基本的な四則演算
- **AdvancedMathService**: 高度な数学演算（べき乗、平方根、階乗）
- **TrigonometryService**: 三角関数計算
- **MathConstantsService**: 数学定数の提供
- **ExpressionEvaluatorService**: 安全な数式評価
- **FileEvaluatorService**: ファイル評価
- **ApiCostExtractorService**: API料金抽出

### 使用技術

- **Go**: プログラミング言語
- **標準ライブラリ**: `flag`, `fmt`, `strconv`, `strings`, `math`など
- **テストフレームワーク**: Go標準のtestingパッケージ

### 安全性機能

- **数式評価の安全性**: 危険なパターン（import, exec, evalなど）の検出と拒否
- **入力検証**: 数値の範囲チェックと型検証
- **エラーハンドリング**: 適切なエラーメッセージと例外処理

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

### 新機能の追加方法

1. **新しいサービスの作成**: `usecases/services.go`に新しいサービス構造体を追加
2. **設定の拡張**: `config/config.go`に必要なパラメータを追加
3. **メイン関数の更新**: `main.go`に新しいハンドラー関数を追加
4. **テストの作成**: 新機能に対するテストケースを作成

### 関連ツール

このCLIツールと同じサービスを使用するMCPツールも利用可能です：
- `devbox/cmd/mcp/arithmetic_calculator/mcp.go`

MCPツールを使用することで、Model Context Protocol経由で同じ計算機能にアクセスできます。

## ライセンス

このプロジェクトのライセンスについては、プロジェクトルートのLICENSEファイルを参照してください。
