# Base64 Extractor

任意のパス（ファイルもしくはディレクトリ）にある画像ファイルのbase64形式のバイト列を抽出するCLIツールです。

## 機能

- 単一の画像ファイルのbase64変換
- ディレクトリ内の画像ファイルの一括base64変換
- 再帰的なディレクトリ検索（オプション）
- テキスト形式またはJSON形式での出力

## サポートしている画像形式

- JPEG (.jpg, .jpeg)
- PNG (.png)
- GIF (.gif)
- BMP (.bmp)
- WebP (.webp)

## インストール

```bash
cd ~/devbox/cmd/cli/base64-extractor
go build -o base64-extractor main.go
```

## 使用方法

### 基本的な使用方法

```bash
./base64-extractor -path=/path/to/image.jpg
```

### オプション

| オプション | 説明 | デフォルト値 | 必須 |
|-----------|------|-------------|------|
| `-path` | 対象パス（ファイルまたはディレクトリ） | なし | ✓ |
| `-recursive` | 再帰検索を行うかどうか | true | |
| `-output-format` | 出力形式（text/json） | text | |

## 使用例
```bash
# 単一ファイルの変換
./base64-extractor -path=/path/to/image.jpg

# ディレクトリ内の画像を一括変換
./base64-extractor -path=/path/to/images/

# 非再帰検索（サブディレクトリを検索しない）
./base64-extractor -path=/path/to/images/ -recursive=false

# JSON形式で出力
./base64-extractor -path=/path/to/images/ -output-format=json
```

## 出力形式

### テキスト形式（デフォルト）

```
=== Base64 Extraction Results ===
Total Images: 2

[1] /path/to/image1.jpg
Base64: /9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAAYEBQYFBAYGBQYHBwYIChAKCgkJChQODwwQFxQYGBcUFhYaHSUfGhsjHBYWICwgIyYnKSopGR8tMC0oMCUoKSj/2wBDAQcHBwoIChMKChMoGhYaKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCj/wAARCAABAAEDASIAAhEBAxEB/8QAFQABAQAAAAAAAAAAAAAAAAAAAAv/xAAUEAEAAAAAAAAAAAAAAAAAAAAA/8QAFQEBAQAAAAAAAAAAAAAAAAAAAAX/xAAUEQEAAAAAAAAAAAAAAAAAAAAA/9oADAMBAAIRAxEAPwCdABmX/9k=

[2] /path/to/image2.png
Base64: iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChAGA4849ewAAAABJRU5ErkJggg==
```

### JSON形式

```json
{
  "images": [
    {
      "file_path": "/path/to/image1.jpg",
      "base64": "/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAAYEBQYFBAYGBQYHBwYIChAKCgkJChQODwwQFxQYGBcUFhYaHSUfGhsjHBYWICwgIyYnKSopGR8tMC0oMCUoKSj/2wBDAQcHBwoIChMKChMoGhYaKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCj/wAARCAABAAEDASIAAhEBAxEB/8QAFQABAQAAAAAAAAAAAAAAAAAAAAv/xAAUEAEAAAAAAAAAAAAAAAAAAAAA/8QAFQEBAQAAAAAAAAAAAAAAAAAAAAX/xAAUEQEAAAAAAAAAAAAAAAAAAAAA/9oADAMBAAIRAxEAPwCdABmX/9k="
    },
    {
      "file_path": "/path/to/image2.png",
      "base64": "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChAGA4849ewAAAABJRU5ErkJggg=="
    }
  ],
  "total": 2
}
```

## エラーハンドリング

- 指定されたパスが存在しない場合、エラーメッセージが表示されます
- 画像ファイルが見つからない場合、「画像ファイルが見つかりませんでした。」と表示されます
- ファイル読み込みエラーが発生した場合、個別のエラーメッセージが表示されます

## 開発

### テストの実行

```bash
cd $HOME/devbox
go test ./internal/base64_extractor/...
```

### カバレッジの確認

```bash
cd $HOME/devbox
go test -coverprofile=coverage.out ./internal/base64_extractor/...
go tool cover -html=coverage.out -o coverage.html
```

## アーキテクチャ

このツールは以下の層で構成されています：

- **CLI層** (`cmd/cli/base64-extractor/main.go`): コマンドライン引数の解析と結果の出力
- **設定層** (`internal/base64_extractor/config/`): フラグ解析と設定管理
- **ユースケース層** (`internal/base64_extractor/usecases/`): ビジネスロジックとサービス実装

### 主要コンポーネント

- `Base64ExtractorService`: 画像ファイルの検索とbase64変換を行うメインサービス
- `Config`: CLI設定を管理する構造体
- `ExtractResult`: 変換結果を保持し、テキスト/JSON形式で出力する構造体

## ライセンス

このプロジェクトは MIT ライセンスの下で公開されています。
