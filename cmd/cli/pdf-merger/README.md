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
go run ./cmd/cli/pdf-merger

# 特定のディレクトリの画像からPDFを作成
go run ./cmd/cli/pdf-merger -dir /path/to/images

# 出力ファイル名を指定
go run ./cmd/cli/pdf-merger -dir /path/to/images -out output.pdf
```

### 2. 既存PDFに画像を追加

```bash
# 既存のPDFに画像を追加
go run ./cmd/cli/pdf-merger -dir /path/to/images -add existing.pdf -out merged.pdf
```

### 3. PDFから画像を抽出

```bash
# PDFの全ページを画像として抽出
go run ./cmd/cli/pdf-merger -extract input.pdf -output-dir ./images

# 特定のページ範囲を抽出（例：3ページから10ページまで）
go run ./cmd/cli/pdf-merger -extract input.pdf -output-dir ./images -start 3 -end 10

# JPEGで特定の開始ページから最後まで抽出
go run ./cmd/cli/pdf-merger -extract input.pdf -output-dir ./images -start 5 -format jpg
```

## オプション

### 基本オプション
- `-dir string`: 画像を検索するフォルダー（デフォルト: "."）
- `-out string`: 出力PDFファイル名（未指定なら <dir名>.pdf）
- `-add string`: 既存のPDFファイルパス（指定時は既存PDFに画像を追加）
- `-recursive bool`: サブディレクトリまで再帰的に画像を検索する（デフォルト: false）

### PDF画像抽出オプション
- `-extract string`: **PDFファイルから画像を抽出する（PDFファイルパス）**
- `-output-dir string`: **画像の出力ディレクトリ（抽出時必須）**
- `-format string`: **出力画像形式（デフォルト: "jpg"）**
  - サポート形式: jpg, jpeg, png, tiff, webp
- `-start int`: **抽出開始ページ（1から開始、0は全ページ）**
- `-end int`: **抽出終了ページ（0は最終ページまで）**

## 注意事項

- **画像抽出機能を使用する場合、`-output-dir`オプションは必須です**
- ページ番号は1から始まります
- 開始ページが終了ページより大きい場合はエラーになります
- サポートされていない画像形式を指定した場合はエラーメッセージが表示されます
- 出力ディレクトリが存在しない場合は自動的に作成されます
- **ファイル名について**: 抽出された画像は`document_0001.jpg`、`document_0002.jpg`のような4桁連番形式で命名されます
- 画像作成時は、デフォルトでは指定ディレクトリ直下のみを検索します。サブディレクトリも対象にする場合は`-recursive`を指定してください
- 現在のバージョンでは、画像作成時はJPG形式（`.jpg`拡張子）の画像ファイルのみがサポートされています
- 画像ファイルはアルファベット順にソートされてPDFに配置されます
- 出力PDFファイルは、既存のファイルがある場合は上書きされます
