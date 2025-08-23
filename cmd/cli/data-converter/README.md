# Data Converter

任意の拡張子のデータソースを別の任意の拡張子のデータソースに変換するためのCLIツールです。

## 概要

このツールは、JSON、CSV、TSV、HTML形式のデータを相互に変換できます。sample1111111.jsにあったJavaScript処理をGoに移植して実装されています。

## 対応形式

### 入力形式
- `json`: JSON形式の2次元配列
- `csv`: CSV（カンマ区切り）形式
- `tsv`: TSV（タブ区切り）形式
- `html`: HTMLテーブル形式

### 出力形式
- `html`: HTMLテーブル形式
- `csv`: CSV（カンマ区切り）形式
- `tsv`: TSV（タブ区切り）形式
- `json`: JSON形式の2次元配列

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
│   └── csv_converter.go               # CSV/TSV変換
└── usecases/services.go               # ビジネスロジック
```

## 対応する変換パターン

| 入力形式 → 出力形式 | 対応状況 |
|-------------------|---------|
| json → html       | ✅      |
| json → csv        | ✅      |
| json → tsv        | ✅      |
| csv → html        | ✅      |
| csv → json        | ✅      |
| csv → tsv         | ✅      |
| tsv → html        | ✅      |
| tsv → json        | ✅      |
| tsv → csv         | ✅      |
| html → json       | ✅      |
| html → csv        | ✅      |
| html → tsv        | ✅      |

## ヘルプ

```bash
./data-converter -help
```

## 実装の背景

このツールは、`sample1111111.js`にあったJavaScript関数群をGoに移植して作成されました：

- `getTableByValues()` → `HTMLConverter.ConvertToHTML()`
- `getCsvByValues()` → `CSVConverter.ConvertToCSV()`
- TSV対応は区切り文字をタブに変更して実装

## 今後の拡張予定

- XML形式のサポート
- YAML形式のサポート
- より高度なHTMLテーブルオプション（CSSクラス指定など）
- バッチ処理機能
- 設定ファイルによるデフォルト値設定
