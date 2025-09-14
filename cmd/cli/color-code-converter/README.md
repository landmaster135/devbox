# カラーコード変換CLIツール

任意のカラーコード形式から別のカラーコード形式に変換するためのCLIツールです。

## サポートされているカラーコード形式

- **HEX** - 16進数形式 (例: `#FF0000`, `#ff0000`, `#F00`)
- **RGB** - RGB形式 (例: `rgb(255,0,0)`)
- **HSL** - HSL形式 (例: `hsl(0,100%,50%)`)
- **HSV** - HSV形式 (例: `hsv(0,100%,100%)`)

## 使用方法

### フラグ形式

```bash
# 基本的な使用方法
go run ./cmd/cli/color-code-converter -src-format hex -dest-format rgb -value "#FF0000"

# 短縮形フラグを使用
go run ./cmd/cli/color-code-converter -s hex -d rgb -v "#FF0000"
```

### 位置引数形式

```bash
go run ./cmd/cli/color-code-converter hex rgb "#FF0000"
```

### ヘルプの表示

```bash
go run ./cmd/cli/color-code-converter -help
go run ./cmd/cli/color-code-converter -h
```

## 使用例

### HEXからRGBに変換

```bash
$ go run ./cmd/cli/color-code-converter -s hex -d rgb -v "#FF0000"
rgb(255,0,0)
```

### RGBからHSLに変換

```bash
$ go run ./cmd/cli/color-code-converter -s rgb -d hsl -v "rgb(255,0,0)"
hsl(0,100%,50%)
```

### HEXからHSVに変換

```bash
$ go run ./cmd/cli/color-code-converter hex hsv "#00FF00"
hsv(120,100%,100%)
```

### 3桁HEXコードの変換

```bash
$ go run ./cmd/cli/color-code-converter -s hex -d rgb -v "#F00"
rgb(255,0,0)
```

### 大文字小文字を区別しない変換

```bash
$ go run ./cmd/cli/color-code-converter -s hex -d rgb -v "#ff0000"
rgb(255,0,0)
```

## オプション

| オプション | 短縮形 | 説明 | 必須 |
|-----------|--------|------|------|
| `-src-format` | `-s` | 変換元のカラーコード形式 (hex, rgb, hsl, hsv) | ✓ |
| `-dest-format` | `-d` | 変換先のカラーコード形式 (hex, rgb, hsl, hsv) | ✓ |
| `-value` | `-v` | 変換するカラーコード値 | ✓ |
| `-help` | `-h` | ヘルプを表示 | |

## エラーハンドリング

ツールは以下の場合にエラーを返します：

- サポートされていない形式を指定した場合
- 無効なカラーコード値を指定した場合
- 必須パラメータが不足している場合

```bash
$ ./color-code-converter -s hex -d rgb -v "invalid"
エラー: カラーコードの解析に失敗しました: 無効なHEX形式です: INVALID
```

## ビルド方法

```bash
# プロジェクトルートから
go build -o color-code-converter ./cmd/cli/color-code-converter
```

## テスト実行

```bash
# 全てのテストを実行
go test ./internal/color_code_converter/...

# カバレッジ付きでテスト実行
go test -cover ./internal/color_code_converter/...
```

## 技術仕様

### アーキテクチャ

- **Config層**: コマンドライン引数の解析と設定管理
- **Domain層**: カラーコードの内部表現と変換ロジック
- **Usecases層**: ビジネスロジックとサービス層
- **CLI層**: メイン処理とエラーハンドリング

### カラー変換アルゴリズム

- HEX ↔ RGB: 直接的な16進数変換
- RGB ↔ HSL: 標準的なHSL変換アルゴリズム
- RGB ↔ HSV: 標準的なHSV変換アルゴリズム
- 全ての変換はRGBを中間表現として使用

### 入力値の検証

- HEX: 3桁または6桁の16進数（#付きまたは無し）
- RGB: `rgb(r,g,b)` 形式（0-255の整数値）
- HSL: `hsl(h,s%,l%)` 形式（H: 0-360, S/L: 0-100%）
- HSV: `hsv(h,s%,v%)` 形式（H: 0-360, S/V: 0-100%）

## ライセンス

このプロジェクトはMITライセンスの下で公開されています。
