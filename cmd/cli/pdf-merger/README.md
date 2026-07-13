# PDF Merger CLI Tool

このツールは、画像からPDFを作成したり、PDFファイルから画像を抽出したりするためのコマンドラインツールです。

## 機能

1. **画像からPDFを作成**: ディレクトリ内の画像ファイルを収集してPDFファイルを生成
2. **既存PDFに画像を追加**: 既存のPDFファイルに画像ページを追加
3. **PDFから画像を抽出**: PDFファイルの任意のページ範囲を画像として抽出

## 使用例

### 1. 画像からPDFを作成

```bash
# 現在のディレクトリの画像からPDFを作成
go run ./cmd/cli/pdf-merger -operation=merge-into-new -output-dir ./dist

# 特定のディレクトリの画像からPDFを作成
go run ./cmd/cli/pdf-merger -operation=merge-into-new -src-dir /path/to/images -output-dir ./dist

# 出力ディレクトリ名をもとにPDFを作成
go run ./cmd/cli/pdf-merger -operation=merge-into-new -src-dir /path/to/images -output-dir ./merged
```

### 2. 既存PDFに画像を追加

```bash
# 既存のPDFに画像を追加
go run ./cmd/cli/pdf-merger -operation=add-into-exist -src-dir /path/to/images -add existing.pdf -output-dir ./merged
```

### 3. PDFから画像を抽出

```bash
# PDFの全ページを画像として抽出
go run ./cmd/cli/pdf-merger -operation=extract-images -src-file input.pdf -output-dir ./images

# 特定のページ範囲を抽出（例：3ページから10ページまで）
go run ./cmd/cli/pdf-merger -operation=extract-images -src-file input.pdf -output-dir ./images -start 3 -end 10

# JPEGで特定の開始ページから最後まで抽出
go run ./cmd/cli/pdf-merger -operation=extract-images -src-file input.pdf -output-dir ./images -start 5 -format jpg
```

## オプション

### 基本オプション

| オプション | 必須/任意 | デフォルト | 説明 |
|---|---|---|---|
| `-operation string` | 必須 | なし | 実行する処理。`merge-into-new`, `add-into-exist`, `extract-images` のいずれか |
| `-src-dir string` | 任意 | `.` | 画像を検索するフォルダー |
| `-output-dir string` | 必須 | なし | PDFの出力ディレクトリ。出力PDFは`<output-dir>/<basename(output-dir)>.pdf`として作成 |
| `-add string` | `add-into-exist`時必須 | なし | 画像を追加する既存PDFファイルパス |
| `-recursive bool` | 任意 | `false` | サブディレクトリまで再帰的に画像を検索する |

### PDF画像抽出オプション

| オプション | 必須/任意 | デフォルト | 説明 |
|---|---|---|---|
| `-operation string` | 必須 | なし | `extract-images` を指定 |
| `-src-file string` | `extract-images`時必須 | なし | 画像を抽出するPDFファイルパス |
| `-output-dir string` | 抽出時必須 | なし | 画像の出力ディレクトリ |
| `-format string` | 任意 | `jpg` | 出力画像形式。サポート形式: `jpg`, `jpeg`, `png`, `tiff`, `webp` |
| `-start int` | 任意 | `0` | 抽出開始ページ。1から開始し、0は全ページ |
| `-end int` | 任意 | `0` | 抽出終了ページ。0は最終ページまで |

## 出力例

### 成功時

```text
検出した画像: 1 枚
PDF を生成しました。完了です。
出力 PDF   : ./dist/dist.pdf
```

### エラー時

```text
-operation は必須です
```

## 注意事項

- **PDF作成・既存PDF追加・画像抽出のいずれでも、`-operation` と `-output-dir` オプションは必須です**
- `-operation=merge-into-new` では `-add` と `-src-file` は指定できません
- `-operation=add-into-exist` では `-add` が必須で、`-src-file` は指定できません
- `-operation=extract-images` では `-src-file` が必須で、`-add` は指定できません
- PDF作成・既存PDF追加時の出力PDFは、`-output-dir ./dist`の場合は`./dist/dist.pdf`として作成されます
- PDF作成・既存PDF追加時に`-output-dir .`を指定した場合は、カレントディレクトリ直下に`<カレントディレクトリ名>.pdf`として作成されます
- ページ番号は1から始まります
- 開始ページが終了ページより大きい場合はエラーになります
- サポートされていない画像形式を指定した場合はエラーメッセージが表示されます
- 出力ディレクトリが存在しない場合は自動的に作成されます
- **ファイル名について**: 抽出された画像は`document_0001.jpg`、`document_0002.jpg`のような4桁連番形式で命名されます
- 画像作成時は、デフォルトでは指定ディレクトリ直下のみを検索します。サブディレクトリも対象にする場合は`-recursive`を指定してください
- 画像作成時の検索元ディレクトリは`-src-dir`で指定します。旧`-dir`オプションは使用できません
- 現在のバージョンでは、画像作成時はJPG形式（`.jpg`拡張子）の画像ファイルのみがサポートされています
- 画像ファイルはアルファベット順にソートされてPDFに配置されます
- 出力PDFファイルは、既存のファイルがある場合は上書きされます
