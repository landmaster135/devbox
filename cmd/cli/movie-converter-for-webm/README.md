# Movie Converter for WEBM

WEBMとMP4の相互変換を行うCLIツールです。

## 機能

- **単一ファイル変換**: 個別のファイルをWEBMからMP4、またはMP4からWEBMに変換
- **バッチ変換**: ディレクトリ内の複数ファイルを一括変換
- **再帰処理**: サブディレクトリも含めた再帰的な変換
- **並列処理**: 最大4並列での高速変換
- **高度なWEBM設定**: CRF/CBRモード、VP9コーデック、Opus/Vorbisオーディオ

## 前提条件

- Go 1.21以上
- FFmpeg（システムにインストール済み）

## インストール

```bash
cd devbox/cmd/cli/movie-converter-for-webm
go build -o movie-converter-for-webm .
```

## 使用方法

### 基本的な使用方法

```bash
# ヘルプを表示
./movie-converter-for-webm -help

# MP4からWEBM（デフォルト設定）
./movie-converter-for-webm -input video.mp4

# WEBMからMP4
./movie-converter-for-webm -input video.webm

# 出力ファイル名を指定
./movie-converter-for-webm -input video.mp4 -output custom.webm
```

### WEBM変換オプション

```bash
# CRFモード（品質重視）
./movie-converter-for-webm -input video.mp4 -crf 25 -audio-codec opus

# CBRモード（ビットレート固定）
./movie-converter-for-webm -input video.mp4 -conversion-mode cbr -video-bitrate 2M -video-quality 80

# Vorbisオーディオコーデック使用
./movie-converter-for-webm -input video.mp4 -audio-codec vorbis -audio-bitrate 192k
```

### バッチ変換

```bash
# MP4ファイルを一括でWEBMに変換
./movie-converter-for-webm -input-dir ./videos -input-ext mp4 -output-dir ./webm -output-ext webm

# WEBMファイルを一括でMP4に変換
./movie-converter-for-webm -input-dir ./webm -input-ext webm -output-dir ./videos -output-ext mp4

# 再帰的にサブディレクトリも処理
./movie-converter-for-webm -input-dir ./media -input-ext mp4 -output-dir ./converted -output-ext webm -recursive

# カスタム設定でバッチ変換
./movie-converter-for-webm -input-dir ./videos -input-ext mp4 -output-dir ./webm -output-ext webm -crf 28 -audio-codec vorbis
```

## コマンドラインオプション

### 基本オプション

| オプション | 説明 | デフォルト |
|-----------|------|-----------|
| `-input` | 入力ファイルのパス（単一ファイル処理時） | - |
| `-output` | 出力ファイルのパス（省略時は自動生成） | 自動生成 |
| `-help` | ヘルプを表示 | - |

### バッチ処理オプション

| オプション | 説明 | デフォルト |
|-----------|------|-----------|
| `-input-dir` | 入力ディレクトリのパス | - |
| `-input-ext` | 入力ファイルの拡張子（例: mp4） | - |
| `-output-dir` | 出力ディレクトリのパス | - |
| `-output-ext` | 出力ファイルの拡張子（例: webm） | - |
| `-recursive` | サブディレクトリも再帰的に処理 | false |

### WEBM変換オプション

| オプション | 説明 | デフォルト |
|-----------|------|-----------|
| `-crf` | Constant Rate Factor（0-63、低いほど高品質） | 30 |
| `-video-bitrate` | ビデオビットレート（例: 1M, 1500k） | 自動検出 |
| `-audio-bitrate` | オーディオビットレート | 128k |
| `-audio-codec` | オーディオコーデック（opus/vorbis） | opus |
| `-conversion-mode` | 変換モード（crf/cbr） | crf |
| `-video-quality` | ビデオ品質（CBRモード時、0-100） | 75 |

## 変換モードについて

### CRF (Constant Rate Factor)
- **特徴**: 品質重視、ファイルサイズは可変
- **用途**: 高品質な動画を作成したい場合
- **CRF値**: 0-63（低いほど高品質、推奨: 23-30）

### CBR (Constant Bit Rate)
- **特徴**: ビットレート固定、ファイルサイズ予測可能
- **用途**: ファイルサイズを制御したい場合
- **設定**: `-video-bitrate`でビットレートを指定

## サポートされているファイル形式

### 入力形式
- MP4 (.mp4)
- MKV (.mkv)
- AVI (.avi)
- MOV (.mov)
- FLV (.flv)
- WEBM (.webm)

### 出力形式
- MP4 (.mp4)
- WEBM (.webm)

## 使用例

WEBM変換
```bash
# 高品質設定（CRF 20）
./movie-converter-for-webm -input video.mp4 -crf 20 -audio-bitrate 192k

# 超高品質設定（CRF 15、Vorbis）
./movie-converter-for-webm -input video.mp4 -crf 15 -audio-codec vorbis -audio-bitrate 256k

# 1MbpsのWEBM
./movie-converter-for-webm -input video.mp4 -conversion-mode cbr -video-bitrate 1M

# 500kbpsの小サイズWEBM
./movie-converter-for-webm -input video.mp4 -conversion-mode cbr -video-bitrate 500k -audio-bitrate 96k

# プロジェクト全体を変換
./movie-converter-for-webm -input-dir ./project -input-ext mp4 -output-dir ./webm_output -output-ext webm -recursive -crf 25

# 複数の動画形式を一括変換
for ext in mp4 mkv avi; do
  ./movie-converter-for-webm -input-dir ./videos -input-ext $ext -output-dir ./webm -output-ext webm -crf 28
done
```

## トラブルシューティング

### FFmpegが見つからない場合
```bash
# Ubuntu/Debian
sudo apt update && sudo apt install ffmpeg

# macOS (Homebrew)
brew install ffmpeg

# Windows (Chocolatey)
choco install ffmpeg
```

### メモリ不足の場合
- 並列処理数を減らす（ソースコード内のセマフォ値を調整）
- より小さなファイルサイズで変換する

### 変換速度が遅い場合
- CRFモードではなくCBRモードを使用
- より高いCRF値を使用（品質は下がる）
- ハードウェアエンコーディングを検討

## 開発者向け情報

### テスト実行
```bash
cd devbox/internal/movie_converter_for_webm
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### アーキテクチャ
- **CLI層**: `main.go` - コマンドライン引数の解析と結果表示
- **サービス層**: `usecases/services.go` - 変換ロジックとビジネスルール
- **統合サービス**: 単一ファイルとバッチ処理の統合インターフェース

### 拡張方法
1. 新しいコーデックサポート: `convertMP4ToWEBM()`メソッドを拡張
2. 新しい入力形式: `GetSupportedExtensions()`を更新
3. 新しい変換オプション: `ConversionConfig`構造体に追加

## ライセンス

このプロジェクトのライセンスについては、プロジェクトルートのLICENSEファイルを参照してください。
