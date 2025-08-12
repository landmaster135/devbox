# OCR Executor

TesseractOCRを使用して画像ファイルからテキストを抽出するCLIツールです。

## 機能

- 単一画像ファイルまたはディレクトリ内の画像ファイルからOCR処理を実行
- 複数言語での同時OCR処理（カンマ区切りで指定）
- 再帰的なディレクトリ検索
- テキスト形式またはJSON形式での出力
- ファイル出力機能（オプション）

## サポートしている画像形式

- JPEG (.jpg, .jpeg)
- PNG (.png)
- GIF (.gif)
- BMP (.bmp)
- WebP (.webp)
- TIFF (.tiff, .tif)

## 使用方法

```bash
# ヘルプを表示
./ocr-executor --help
```

## 使用例

```bash
# 単一画像ファイルをOCR処理
./ocr-executor -path=/path/to/image.jpg

# ディレクトリ内の画像ファイルをOCR処理
./ocr-executor -path=/path/to/directory

# 再帰的にディレクトリを検索してOCR処理
./ocr-executor -path=/path/to/directory -recursive=true

# 日本語と英語でOCR処理（デフォルト）
./ocr-executor -path=/path/to/image.jpg -language=jpn,eng

# 英語のみでOCR処理
./ocr-executor -path=/path/to/image.jpg -language=eng

# 複数言語でOCR処理
./ocr-executor -path=/path/to/image.jpg -language=jpn,eng,fra

# テキスト形式で出力（デフォルト）
./ocr-executor -path=/path/to/image.jpg -output-format=text

# JSON形式で出力
./ocr-executor -path=/path/to/image.jpg -output-format=json

# 結果をファイルに保存（テキスト形式）
./ocr-executor -path=/path/to/directory -output-dir=/path/to/output

# 結果をファイルに保存（JSON形式）
./ocr-executor -path=/path/to/directory -output-format=json -output-dir=/path/to/output
```

## コマンドラインオプション

| オプション | 型 | デフォルト | 説明 |
|-----------|---|-----------|------|
| `-path` | string | - | 対象パス（必須） |
| `-recursive` | bool | false | 再帰検索 |
| `-language` | string | jpn,eng | OCR言語（カンマ区切り） |
| `-output-format` | string | text | 出力形式（text/json） |
| `-output-dir` | string | - | 出力ディレクトリ（オプション） |

## 出力例

### テキスト形式

```
=== OCR Results ===
Total Images: 2

[1] /path/to/image1.jpg
Text: Hello World
これはテストです

[2] /path/to/image2.png
Text: Sample Text
サンプルテキスト
```

### JSON形式

```json
{
  "results": [
    {
      "file_path": "/path/to/image1.jpg",
      "text": "Hello World\nこれはテストです"
    },
    {
      "file_path": "/path/to/image2.png",
      "text": "Sample Text\nサンプルテキスト"
    }
  ],
  "total": 2
}
```

## 前提条件

- Tesseract OCRがシステムにインストールされている必要があります
- 使用する言語パックがインストールされている必要があります

### Ubuntu/Debianでのインストール

```bash
# Tesseract OCRのインストール
sudo apt-get install tesseract-ocr

# 日本語言語パックのインストール
sudo apt-get install tesseract-ocr-jpn

# その他の言語パック（例：フランス語）
sudo apt-get install tesseract-ocr-fra
```

## ビルド方法

```bash
go build -o ocr-executor ./cmd/cli/ocr-executor
```

## テスト実行

```bash
go test ./internal/ocr_executor/...
```

## エラーハンドリング

- 存在しないパスが指定された場合はエラーメッセージを表示
- サポートされていない画像形式は無視
- OCR処理に失敗した画像はエラー情報と共に結果に含まれる
- 無効な言語コードが指定された場合はエラーメッセージを表示

## 制限事項

- 言語コードは3文字以下である必要があります
- 出力形式は'text'または'json'のみサポート
- 大量の画像ファイルを処理する場合は時間がかかる場合があります
