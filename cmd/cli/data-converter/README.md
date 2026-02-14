# Data Converter

入力ファイルをいったん **key-value 型リスト（`[]map[string]string`）** に正規化してから、別形式ファイルへ変換する CLI ツールです。

## 対応形式

- `json`
- `yaml` (`.yaml` / `.yml`)
- `csv`
- `tsv`
- `html` (`table` 要素)
- `md-ordered-list`
- `md-unordered-list`
- `md-table`

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
| `-input-format` | 任意 | 入力ファイル拡張子から推定 | 入力形式（`json|yaml|csv|tsv|html|md-ordered-list|md-unordered-list|md-table`） |
| `-output-format` | 任意 | 出力ファイル拡張子から推定 | 出力形式（`json|yaml|csv|tsv|html|md-ordered-list|md-unordered-list|md-table`） |
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

### TSV から JSON へ変換

```bash
go run ./cmd/cli/data-converter \
  -input-file-path ./users.tsv \
  -output-file-path ./users.json
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

### Markdown ordered-list から JSON へ変換

```bash
go run ./cmd/cli/data-converter \
  -input-file-path ./users.md \
  -output-file-path ./users.json \
  -input-format md-ordered-list
```

### JSON から Markdown unordered-list へ変換

```bash
go run ./cmd/cli/data-converter \
  -input-file-path ./users.json \
  -output-file-path ./users.md \
  -output-format md-unordered-list
```

### CSV から Markdown table へ変換

```bash
go run ./cmd/cli/data-converter \
  -input-file-path ./users.csv \
  -output-file-path ./users.md \
  -output-format md-table
```

## 注意事項

- `md-ordered-list` / `md-unordered-list` へ出力する場合、単一キーのデータ（または `item` キーが含まれるデータ）を指定してください。
- Markdown系形式（`md-ordered-list` / `md-unordered-list` / `md-table`）は `.md` 拡張子だけでは自動判定できないため、必要に応じて `-input-format` / `-output-format` を明示してください。

## 出力例

### 成功時

```text
変換完了: ./users.json (json) -> ./users.csv (csv)
```

### エラー時

```text
エラー: 未対応の入力形式です: txt
```
