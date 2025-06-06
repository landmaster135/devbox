# unit-converter

物理量を **長さ・質量（重さ）・温度・面積・体積** の 5 カテゴリで相互変換できる高速 CLI ツールです。SI 接頭語 **yotta 〜 yocto** を自動解釈するため、`µm`, `nm2`, `pL` など自由に記述できます。

## 概要

- **依存ゼロ**: すべて標準ライブラリのみ。ビルドして置くだけで利用可能。  
- **SI プレフィックス対応**: 接頭語 + 指数 (m², m³) も正しく処理。  
- **疎結合構成**: CLI (`cmd/unit-converter`) とドメインロジック (`internal/independencies/unit_converter/usecases`) を分離。  
- **高カバレッジテスト**: `go test -cover` で 90 %以上。  

## インストール

```bash
go install github.com/your-org/your-repo/cmd/unit-converter@latest
```

またはソースをクローンしてビルド:

```bash
git clone https://github.com/your-org/your-repo.git
cd cmd/unit-converter
go build -o unit-converter
```

## 使用方法

### コマンドラインオプション

| オプション | 説明 | デフォルト |
|------------|------|------------|
| `-p`, `--precision` | 表示する有効桁数 | `6` |
| `--list`   | サポートするカテゴリ・基底単位と SI 接頭語表を表示 | - |
| `-h`, `--help` | ヘルプを表示 | - |

引数順は `<category> <value> <from-unit> <to-unit>` です。

- `category` : `length`, `weight`(`mass`), `temperature`(`temp`), `area`, `volume` のいずれか。  
- `value`    : 変換元数値 (float)。  
- `from-unit`, `to-unit` : 単位（SI 接頭語付き可）。

## 使用例

### 長さの変換
```bash
unit-converter length 12.7 mm in   # 12.7 mm → inch
```

### 温度の変換
```bash
unit-converter -p4 temp 100 F C    # 100°F → 37.78°C (4 桁表示)
```

### 面積と体積
```bash
unit-converter area 1 ha m2        # 1 ha → 10000 m²
unit-converter volume 0.25 uL nL   # 0.25 µL → 250 nL
```

### 一覧表示
```bash
unit-converter --list              # 基底単位 & 接頭語表を表示
```

## 出力例

```
$ unit-converter volume 500 ml cup
500 ml = 2.11338 cup
```

## ディレクトリ構成

```
.
├── cmd/unit-converter                # CLI (プレゼンテーション層)
└── internal/independencies/unit_converter/
    ├── usecases.go                   # 変換エンジン
    └── usecases_test.go              # ユニットテスト
```

## `usecases` 公開 API

```go
// 値を変換
Convert(category string, value float64, fromUnit, toUnit string) (float64, error)

// カテゴリと基底単位を取得
Categories() map[string][]string

// SI 接頭語表 (文字列) を取得
PrefixTable() string
```

CLI は上記 3 関数のみに依存し、ロジック入れ替えや拡張が容易です。

## テスト

```bash
go test ./internal/independencies/unit_converter/... -cover
```

- **Convert** の成功・エラーケース
- **ラウンドトリップ性**（往復で元値に戻る）
- **Categories の防御コピー**確認
- **PrefixTable** に全接頭語が含まれるか

## 拡張方法

1. **単位追加**: `usecases` 内の factor マップに基底単位を追記。  
2. **カテゴリ追加**: factor マップと `Convert()` の `switch` を追加。  
3. 係数で表せない変換（エネルギー, 通貨など）は `convertXxx` 関数を新規実装。

## ライセンス

MIT © landmaster135
