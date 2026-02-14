# Data Converter

入力ファイルをいったん **key-value 型リスト（`[]map[string]string`）** に正規化してから、別形式ファイルへ変換する CLI ツールです。

## 対応形式

- `json`
- `yaml` (`.yaml` / `.yml`)
- `csv`
- `html` (`table` 要素)

## インストール

```bash
go build -o data-converter ./cmd/cli/data-converter
```

## 使用方法

```bash
go run ./cmd/cli/data-converter \
  -input-file-path ./input.json \
  -output-file-path ./output.csv
```

形式は `-input-format` / `-output-format` で明示できます。省略時は拡張子から推定します。

```bash
go run ./cmd/cli/data-converter \
  -input-file-path ./input.data \
  -output-file-path ./output.data \
  -input-format json \
  -output-format html
```

## フラグ一覧

| フラグ | 必須 | デフォルト値 | 説明 |
| --- | --- | --- | --- |
| `-input-file-path` | 必須 | なし | 入力ファイルパス |
| `-output-file-path` | 必須 | なし | 出力ファイルパス |
| `-input-format` | 任意 | 入力ファイル拡張子から推定 | 入力形式（`json|yaml|csv|html`） |
| `-output-format` | 任意 | 出力ファイル拡張子から推定 | 出力形式（`json|yaml|csv|html`） |
| `-help` | 任意 | `false` | ヘルプを表示 |

## 使用例

### JSON から CSV へ変換

```bash
go run ./cmd/cli/data-converter \
  -input-file-path ./users.json \
  -output-file-path ./users.csv
```

### CSV から YAML へ変換

```bash
go run ./cmd/cli/data-converter \
  -input-file-path ./users.csv \
  -output-file-path ./users.yaml
```

### HTML から JSON へ変換

```bash
go run ./cmd/cli/data-converter \
  -input-file-path ./users.html \
  -output-file-path ./users.json
```

### YAML から HTML へ変換

```bash
go run ./cmd/cli/data-converter \
  -input-file-path ./users.yaml \
  -output-file-path ./users.html
```

## 出力例

### 成功時

```text
変換完了: ./users.json (json) -> ./users.csv (csv)
```

### エラー時

```text
エラー: 未対応の入力形式です: txt
```
