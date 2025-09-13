# Data Converter

任意の拡張子のデータソースを別の任意の拡張子のデータソースに変換するためのCLIツールです。

## 概要

このツールは、JSON、CSV、TSV、HTML、Markdownリスト形式のデータを相互に変換できます。sample1111111.jsにあったJavaScript処理をGoに移植して実装されています。

## 対応形式

### 入力形式
- `json`: JSON形式の2次元配列
- `csv`: CSV（カンマ区切り）形式
- `tsv`: TSV（タブ区切り）形式
- `html`: HTMLテーブル形式
- `list`: Markdownの箇条書きリスト形式（`- 項目`）
- `ordered-list`: Markdownの順序付きリスト形式（`1. 項目`）
- `table`: Markdownテーブル形式（`| 列1 | 列2 |`）

### 出力形式
- `json`: JSON形式の2次元配列
- `csv`: CSV（カンマ区切り）形式
- `tsv`: TSV（タブ区切り）形式
- `html`: HTMLテーブル形式
- `list`: Markdownの箇条書きリスト形式（`- 項目`）
- `ordered-list`: Markdownの順序付きリスト形式（`1. 項目`）
- `table`: Markdownテーブル形式（`| 列1 | 列2 |`）

## インストール

```bash
cd devbox
go build -o data-converter cmd/cli/data-converter/main.go
```

## 使用方法

### 基本構文

```bash
go run ./cmd/cli/data-converter -input-format=<入力形式> -output-format=<出力形式> [入力オプション]
```

### 入力オプション

以下のいずれかを指定してください（排他的）：
- `-input=<データ>`: 直接データを入力
- `-input-file-path=<ファイルパス>`: ファイルから入力

## 使用例

### 1. JSON文字列をHTMLテーブルに変換

```bash
go run ./cmd/cli/data-converter -input-format=json -output-format=html -input='[["名前","年齢"],["田中","25"],["佐藤","30"]]'
```

出力:
```html
<table>
<thead>
<tr><th>名前</th><th>年齢</th></tr>
</thead>
<tbody>
<tr><td>田中</td><td>25</td></tr>
<tr><td>佐藤</td><td>30</td></tr>
</tbody>
</table>
```

### 2. CSVファイルをTSVに変換

```bash
go run ./cmd/cli/data-converter -input-format=csv -output-format=tsv -input-file-path=data.csv
```

### 3. TSV文字列をJSONに変換

```bash
go run ./cmd/cli/data-converter -input-format=tsv -output-format=json -input='名前	年齢
田中	25
佐藤	30'
```

出力:
```json
[["名前","年齢"],["田中","25"],["佐藤","30"]]
```

### 4. JSONファイルをCSVに変換

```bash
go run ./cmd/cli/data-converter -input-format=json -output-format=csv -input-file-path=data.json
```

### 5. HTMLテーブルをJSONに変換

```bash
go run ./cmd/cli/data-converter -input-format=html -output-format=json -input='<table>
<tr><th>名前</th><th>年齢</th></tr>
<tr><td>田中</td><td>25</td></tr>
<tr><td>佐藤</td><td>30</td></tr>
</table>'
```

出力:
```json
[["名前","年齢"],["田中","25"],["佐藤","30"]]
```

### 6. CSVをJSONに変換

```bash
go run ./cmd/cli/data-converter -input-format=csv -output-format=json -input='"名前","年齢"
"田中","25"
"佐藤","30"'
```

出力:
```json
[["名前","年齢"],["田中","25"],["佐藤","30"]]
```

### 7. 箇条書きリストをHTMLテーブルに変換

```bash
go run ./cmd/cli/data-converter -input-format=list -output-format=html -input='- 項目1
- 項目2
- 項目3'
```

出力:
```html
<table>
<thead>
<tr><th>項目</th></tr>
</thead>
<tbody>
<tr><td>項目1</td></tr>
<tr><td>項目2</td></tr>
<tr><td>項目3</td></tr>
</tbody>
</table>
```

### 8. 順序付きリストをCSVに変換

```bash
go run ./cmd/cli/data-converter -input-format=ordered-list -output-format=csv -input='1. 項目1
2. 項目2
3. 項目3'
```

出力:
```csv
"番号","項目"
"1","項目1"
"2","項目2"
"3","項目3"
```

### 9. HTMLテーブルを箇条書きリストに変換

```bash
go run ./cmd/cli/data-converter -input-format=html -output-format=list -input='<table>
<tr><th>項目</th></tr>
<tr><td>項目1</td></tr>
<tr><td>項目2</td></tr>
</table>'
```

出力:
```
- 項目1
- 項目2
```

### 10. CSVを順序付きリストに変換

```bash
go run ./cmd/cli/data-converter -input-format=csv -output-format=ordered-list -input='"項目"
"項目1"
"項目2"
"項目3"'
```

出力:
```
1. 項目1
2. 項目2
3. 項目3
```

### 11. JSONをMarkdownテーブルに変換

```bash
go run ./cmd/cli/data-converter -input-format=json -output-format=table -input='[["名前","年齢","職業"],["田中","25","エンジニア"],["佐藤","30","デザイナー"]]'
```

出力:
```markdown
| 名前 | 年齢 | 職業          |
|--------|--------|-----------------|
| 田中 | 25     | エンジニア |
| 佐藤 | 30     | デザイナー |
```

### 12. MarkdownテーブルをJSONに変換

```bash
go run ./cmd/cli/data-converter -input-format=table -output-format=json -input='| 名前 | 年齢 | 職業 |
|------|------|------|
| 田中 | 25   | エンジニア |
| 佐藤 | 30   | デザイナー |'
```

出力:
```json
[{"名前":"田中","年齢":25,"職業":"エンジニア"},{"名前":"佐藤","年齢":30,"職業":"デザイナー"}]
```

### 13. Markdownテーブルを箇条書きリストに変換

```bash
go run ./cmd/cli/data-converter -input-format=table -output-format=list -input='| 名前 | 年齢 | 職業 |
|------|------|------|
| 田中 | 25   | エンジニア |
| 佐藤 | 30   | デザイナー |'
```

出力:
```
- 名前: 田中, 年齢: 25, 職業: エンジニア
- 名前: 佐藤, 年齢: 30, 職業: デザイナー
```

### 14. CSVをMarkdownテーブルに変換

```bash
go run ./cmd/cli/data-converter -input-format=csv -output-format=table -input='"名前","年齢","職業"
"田中","25","エンジニア"
"佐藤","30","デザイナー"'
```

出力:
```markdown
| 名前 | 年齢 | 職業          |
|--------|--------|-----------------|
| 田中 | 25     | エンジニア |
| 佐藤 | 30     | デザイナー |
```

## 機能詳細

### HTML変換の特徴

- 最初の行を自動的にヘッダー（`<th>`）として扱います
- 空のセルには「💩」が表示されます（ヘッダー行は除く）
- 適切な`<thead>`と`<tbody>`タグが生成されます

### HTML解析の特徴

- `<table>`タグ内の`<tr>`、`<td>`、`<th>`要素を自動解析
- HTMLタグの除去とHTMLエンティティのデコードに対応
- `<thead>`と`<tbody>`の区別なく、すべての行を順次処理
- ネストしたHTMLタグも適切に処理

### CSV/TSV変換の特徴

- 区切り文字、改行、ダブルクォートを含むセルは自動的にダブルクォートで囲まれます
- ダブルクォートのエスケープ（`""`）に対応しています
- RFC 4180準拠のCSV形式をサポートします

### JSON変換の特徴

- 2次元配列形式での入出力をサポートします
- 文字列以外の値も適切に処理されます

### Markdownリスト変換の特徴

#### 箇条書きリスト（list）
- `- 項目` または `* 項目` 形式をサポート
- テーブルから変換する場合、1列なら単純なリスト、複数列なら「項目名: 値」形式

#### 順序付きリスト（ordered-list）
- `1. 項目` 形式をサポート
- 番号は連続している必要があります
- テーブルから変換する場合、1列なら単純なリスト、複数列なら「項目名: 値」形式

### Markdownテーブル変換の特徴

#### Markdownテーブル（table）
- `| 列1 | 列2 |` 形式をサポート
- セパレーター行（`|---|---|`）が必要です
- 各列の幅は自動調整されます
- 空のセルも適切に処理されます
- パイプ文字（|）で区切られた形式を解析・生成します

## エラーハンドリング

- 必須パラメータの不足
- 未対応の入力/出力形式
- ファイル読み込みエラー
- データ解析エラー
- 排他的パラメータの同時指定

## アーキテクチャ

```
cmd/cli/data-converter/main.go          # CLIエントリーポイント
internal/data_converter/
├── config/config.go                    # フラグ解析とコンフィグ
├── parsers/data_parser.go             # 入力データ解析
├── converters/                        # 変換処理
│   ├── html_converter.go              # HTML変換
│   ├── csv_converter.go               # CSV/TSV変換
│   └── markdown_converter.go          # Markdownリスト変換
└── usecases/services.go               # ビジネスロジック
```

## 対応する変換パターン

| 入力形式 → 出力形式 | 対応状況 |
|-------------------|---------|
| json → html       | ✅      |
| json → csv        | ✅      |
| json → tsv        | ✅      |
| json → list       | ✅      |
| json → ordered-list | ✅    |
| csv → html        | ✅      |
| csv → json        | ✅      |
| csv → tsv         | ✅      |
| csv → list        | ✅      |
| csv → ordered-list | ✅     |
| tsv → html        | ✅      |
| tsv → json        | ✅      |
| tsv → csv         | ✅      |
| tsv → list        | ✅      |
| tsv → ordered-list | ✅     |
| html → json       | ✅      |
| html → csv        | ✅      |
| html → tsv        | ✅      |
| html → list       | ✅      |
| html → ordered-list | ✅    |
| list → html       | ✅      |
| list → csv        | ✅      |
| list → tsv        | ✅      |
| list → json       | ✅      |
| list → ordered-list | ✅    |
| ordered-list → html | ✅    |
| ordered-list → csv | ✅     |
| ordered-list → tsv | ✅     |
| ordered-list → json | ✅    |
| ordered-list → list | ✅    |
| table → html       | ✅      |
| table → csv        | ✅      |
| table → tsv        | ✅      |
| table → json       | ✅      |
| table → list       | ✅      |
| table → ordered-list | ✅    |
| json → table       | ✅      |
| csv → table        | ✅      |
| tsv → table        | ✅      |
| html → table       | ✅      |
| list → table       | ✅      |
| ordered-list → table | ✅    |

## ヘルプ

```bash
./data-converter -help
```
