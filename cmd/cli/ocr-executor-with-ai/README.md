# OCR Executor with AI

AI（Gemini API / Vertex AI / Ollama）を使用して画像からテキストを抽出するCLIツールです。

## 概要

このツールは、Google Gemini APIやVertex AIに加えて、ローカルのOllamaモデルを利用して画像に対してOCR（光学文字認識）を実行します。単一の画像ファイルまたはディレクトリ内の複数の画像を処理できます。

## 機能

- **画像OCR**: Gemini APIを使用した高精度なテキスト抽出
- **Markdownテーブル生成**: 表形式の画像をMarkdownテーブル形式で出力
- **複数画像対応**: ディレクトリ内の画像を一括処理
- **再帰検索**: サブディレクトリも含めた画像検索
- **カスタマイズ可能**: プロンプト、システム指示、生成パラメータの調整
- **多様な画像形式**: jpg, jpeg, png, gif, bmp, webp をサポート

## インストール

```bash
# プロジェクトルートから
go build -o bin/ocr-executor-with-ai ./cmd/cli/ocr-executor-with-ai
```

## 使用例

### 基本的な使用方法

```bash
# Gemini API使用（単一画像ファイル）
go run ./cmd/cli/ocr-executor-with-ai -path /path/to/image.webp -ai-type gemini -api-key "your-api-key"

# Vertex AI使用（単一画像ファイル）
go run ./cmd/cli/ocr-executor-with-ai -path /path/to/image.webp -ai-type vertex -project "your-project-id"

# Ollama使用（単一画像ファイル）
go run ./cmd/cli/ocr-executor-with-ai -path /path/to/image.webp -ai-type ollama

# ディレクトリ内の画像（再帰）
go run ./cmd/cli/ocr-executor-with-ai -path /path/to/directory -recursive -ai-type gemini -api-key "your-api-key"
```

### 詳細設定

```bash
# Gemini APIでカスタムプロンプトとモデル指定
go run ./cmd/cli/ocr-executor-with-ai \
  -path /path/to/screenshots \
  -recursive \
  -ai-type gemini \
  -api-key "your-api-key" \
  -model gemini-2.0-flash \
  -prompt "この画像からテキストを抽出して" \
  -system-instruction "OCRして。" \
  -temperature 0.8 \
  -max-tokens 4096

# Vertex AIで詳細設定
go run ./cmd/cli/ocr-executor-with-ai \
  -path /path/to/screenshots \
  -recursive \
  -ai-type vertex \
  -project "your-project-id" \
  -location "us-central1" \
  -model gemini-1.5-pro-002 \
  -temperature 0.5 \
  -max-tokens 2048

# Ollamaで詳細設定
go run ./cmd/cli/ocr-executor-with-ai \
  -path /path/to/screenshots \
  -recursive \
  -ai-type ollama \
  -model qwen2.5vl \
  -prompt "表をOCRして" \
  -temperature 0.8 \
  -max-tokens 2048
```

### 短縮形オプション

```bash
# Gemini API使用（短縮形）
go run ./cmd/cli/ocr-executor-with-ai -p /path/to/image.webp -at gemini -ak "your-api-key" -m gemini-1.5-pro-002 -pr "テキストを抽出" -t 0.5 -mt 2048

# Vertex AI使用（短縮形）
go run ./cmd/cli/ocr-executor-with-ai -p /path/to/image.webp -at vertex -pj "your-project-id" -loc "us-central1"

# Ollama使用（短縮形）
go run ./cmd/cli/ocr-executor-with-ai -p /path/to/image.webp -at ollama -m qwen2.5vl
```

### Markdownテーブル生成

```bash
# Markdownテーブル形式でOCRを実行
go run ./cmd/cli/ocr-executor-with-ai -path /path/to/table-image.webp -generates-markdown-table -ai-type gemini -api-key "your-api-key"

# 短縮形
go run ./cmd/cli/ocr-executor-with-ai -p /path/to/table-image.webp -gmt -at gemini -ak "your-api-key"

# ディレクトリ内の複数画像をMarkdownテーブル形式で処理
go run ./cmd/cli/ocr-executor-with-ai -path /path/to/table-images -recursive -generates-markdown-table -ai-type gemini -api-key "your-api-key"

# OllamaでMarkdownテーブル形式を生成
go run ./cmd/cli/ocr-executor-with-ai -path /path/to/table-image.webp -generates-markdown-table -ai-type ollama -m qwen2.5vl
```

## オプション

### 共通オプション

| オプション | 短縮形 | 説明 | デフォルト値 |
|-----------|--------|------|-------------|
| `-path` | `-p` | 画像ファイルまたはディレクトリのパス | `.` (カレントディレクトリ) |
| `-recursive` | `-r` | ディレクトリを再帰的に検索 | false |
| `-ai-type` | `-at` | 利用するAIタイプ（`gemini` / `vertex` / `ollama`） | gemini |
| `-prompt` | `-pr` | OCR用プロンプト（`-generates-markdown-table`と併用不可） | "OCRして。補足や説明は不要です。" |
| `-system-instruction` | `-si` | システム指示 | "OCRして。" |
| `-generates-markdown-table` | `-gmt` | Markdownテーブル形式でOCRを実行 | false |
| `-temperature` | `-t` | 生成パラメータ（0.0-2.0） | 1.0 |
| `-max-tokens` | `-mt` | 最大トークン数 | 8192 |
| `-help` | `-h` | ヘルプを表示 | - |

### Gemini（`-ai-type gemini`）専用オプション

| オプション | 短縮形 | 説明 | デフォルト値 |
|-----------|--------|------|-------------|
| `-model` | `-m` | 利用するGeminiモデル名 | `gemini-2.5-flash-lite` |
| `-api-key` | `-ak` | Gemini API キー（必須） | - |

### Vertex AI（`-ai-type vertex`）専用オプション

| オプション | 短縮形 | 説明 | デフォルト値 |
|-----------|--------|------|-------------|
| `-model` | `-m` | 利用するVertex AIモデル名 | `gemini-1.5-pro-002` |
| `-project` | `-pj` | Google Cloud プロジェクトID（必須） | - |
| `-location` | `-loc` | Google Cloud ロケーション | us-central1 |

### Ollama（`-ai-type ollama`）専用オプション

| オプション | 短縮形 | 説明 | デフォルト値 |
|-----------|--------|------|-------------|
| `-model` | `-m` | 利用するOllamaモデル名 | `qwen2.5vl` |

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
