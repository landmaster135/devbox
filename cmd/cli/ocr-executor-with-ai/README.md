# OCR Executor with AI

AI（Gemini API）を使用して画像からテキストを抽出するCLIツールです。

## 概要

このツールは、Google Gemini APIを使用して画像に対してOCR（光学文字認識）を実行します。単一の画像ファイルまたはディレクトリ内の複数の画像を処理できます。

## 機能

- **画像OCR**: Gemini APIを使用した高精度なテキスト抽出
- **複数画像対応**: ディレクトリ内の画像を一括処理
- **再帰検索**: サブディレクトリも含めた画像検索
- **カスタマイズ可能**: プロンプト、システム指示、生成パラメータの調整
- **多様な画像形式**: jpg, jpeg, png, gif, bmp, webp をサポート

## インストール

```bash
# プロジェクトルートから
cd /home/nov/devbox
go build -o bin/ocr-executor-with-ai ./cmd/cli/ocr-executor-with-ai
```

## 環境設定

### Gemini Developer API使用時
```bash
export GOOGLE_API_KEY="your-api-key"
```

### Vertex AI使用時
```bash
export GOOGLE_GENAI_USE_VERTEXAI=true
export GOOGLE_CLOUD_PROJECT="your-project-id"
export GOOGLE_CLOUD_LOCATION="us-central1"
```

## 使用方法

### 基本的な使用方法

```bash
# 単一画像ファイル
./bin/ocr-executor-with-ai -path /path/to/image.webp

# ディレクトリ内の画像（非再帰）
./bin/ocr-executor-with-ai -path /path/to/directory

# ディレクトリ内の画像（再帰）
./bin/ocr-executor-with-ai -path /path/to/directory -recursive
```

### 詳細設定

```bash
# カスタムプロンプトとモデル指定
./bin/ocr-executor-with-ai \
  -path /path/to/screenshots \
  -recursive \
  -model gemini-2.0-flash \
  -prompt "この画像からテキストを抽出して" \
  -system-instruction "OCRして。" \
  -temperature 0.8 \
  -max-tokens 4096
```

### 短縮形オプション

```bash
# 短縮形を使用した例
./bin/ocr-executor-with-ai -p /path/to/image.webp -r -m gemini-1.5-pro-002 -pr "テキストを抽出" -t 0.5 -mt 2048
```

## オプション

| オプション | 短縮形 | 説明 | デフォルト値 |
|-----------|--------|------|-------------|
| `-path` | `-p` | 画像ファイルまたはディレクトリのパス（必須） | - |
| `-recursive` | `-r` | ディレクトリを再帰的に検索 | false |
| `-model` | `-m` | 使用するGeminiモデル | gemini-2.5-flash-lite |
| `-prompt` | `-pr` | OCR用プロンプト | "OCRして。補足や説明は不要です。" |
| `-system-instruction` | `-si` | システム指示 | "OCRして。" |
| `-temperature` | `-t` | 生成パラメータ（0.0-2.0） | 1.0 |
| `-max-tokens` | `-mt` | 最大トークン数 | 8192 |
| `-help` | `-h` | ヘルプを表示 | - |

## サポートされる画像形式

- JPEG (.jpg, .jpeg)
- PNG (.png)
- GIF (.gif)
- BMP (.bmp)
- WebP (.webp)

## 出力例

```
=== AI OCR実行結果 ===
処理総数: 2件 (成功: 2件, 失敗: 0件)

[1] /path/to/image1.png
OCR結果:
Hello World
This is a test image.

[2] /path/to/image2.webp
OCR結果:
Sample Text
Another line of text.
```

## エラーハンドリング

- パスが指定されていない場合
- 指定されたパスが存在しない場合
- サポートされていない画像形式の場合
- API呼び出しエラーの場合

すべてのエラーは適切な日本語メッセージで表示されます。

## 技術仕様

### アーキテクチャ

```
internal/ocr_executor_with_ai/
├── config/          # CLI設定とフラグ解析
├── usecases/        # ビジネスロジック
└── models/          # データ型定義
```

### 依存関係

- `google.golang.org/genai` - Gemini API クライアント
- 既存の `base64_services.go` - 画像のBase64変換

### レート制限対策

複数画像処理時は、API呼び出し間隔を2秒空けてレート制限を回避します。

## 開発者向け情報

### テスト実行

```bash
# ヘルプ表示テスト
go run ./cmd/cli/ocr-executor-with-ai -help

# エラーハンドリングテスト
go run ./cmd/cli/ocr-executor-with-ai
go run ./cmd/cli/ocr-executor-with-ai -path "/nonexistent/path"
```

### 拡張方法

新しい機能を追加する場合は、以下のパターンに従ってください：

1. `config/config.go` - 新しいオプションを追加
2. `usecases/ocr_executor_service.go` - ビジネスロジックを実装
3. `models/types.go` - 必要に応じて新しい型を定義

## ライセンス

このプロジェクトのライセンスに従います。

## 参考

このツールは、PythonのGoogle Colabで動作していたAI OCRコードをGoで再実装したものです。
