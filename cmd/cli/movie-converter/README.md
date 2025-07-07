# Movie Converter

GIFとMP4を相互に変換するCLIツールです。PowerShellスクリプト（Z5-5_convert_mp4_to_gif.ps1、Z5-10_convert_gif_to_mp4.ps1）と同等の機能を提供します。

## 機能

- MP4/MKV → GIF変換
- GIF → MP4変換
- 高品質なGIF生成（パレット生成とパレット使用）
- カスタマイズ可能なFPS、幅、速度、ループ設定
- 自動的な出力ファイル名生成

## インストール

```bash
go build -o movie-converter main.go
```

## 使用例

### 基本的な使用例

```bash
# MP4からGIF（デフォルト設定）
./movie-converter -input video.mp4

# GIFからMP4（デフォルト設定）
./movie-converter -input animation.gif

# MP4からGIF
./movie-converter -input video.mp4 -output animation.gif

# ヘルプを表示
./movie-converter -help
```

### 詳細なオプション

```bash
# MP4からGIF（カスタム設定）
./movie-converter -input video.mp4 -output custom.gif -fps 30 -width 320 -speed 1.5 -loop -1

# GIFからMP4（カスタムFPS）
./movie-converter -input animation.gif -output video.mp4 -fps 24
```

## オプション

| オプション | 説明 | デフォルト値 |
|-----------|------|-------------|
| `-input` | 入力ファイルのパス（必須） | - |
| `-output` | 出力ファイルのパス | 自動生成 |
| `-fps` | 出力のFPS | MP4→GIF: 60, GIF→MP4: 15 |
| `-width` | GIF出力時の幅（0=デフォルト品質） | 0 |
| `-speed` | GIF出力時の速度倍率 | 2.0 |
| `-loop` | GIF出力時のループ設定（0=無限ループ, -1=ループなし） | 0 |
| `-use-itsscale` | GIF出力時にitsscaleを使用するか | true |

## サポートされているファイル形式

### 入力
- `.mp4` - MP4動画ファイル
- `.mkv` - MKV動画ファイル
- `.gif` - GIFアニメーションファイル

### 出力
- `.gif` - GIFアニメーションファイル
- `.mp4` - MP4動画ファイル

## デフォルト値について

### MP4→GIF変換のデフォルト値
PowerShellスクリプト（Z5-5_convert_mp4_to_gif.bat）に基づく：
- FPS: 60
- 幅: 0（デフォルト品質）
- 速度: 2.0（2倍速）
- itsScale使用: true
- ループ: 0（無限ループ）

### GIF→MP4変換のデフォルト値
PowerShellスクリプト（Z5-10_convert_gif_to_mp4.ps1）に基づく：
- FPS: 15

## 例

### 1. 基本的なMP4からGIF変換
```bash
./movie-converter -input sample.mp4
# 出力: sample.gif（FPS=60, 2倍速, 無限ループ）
```

### 2. カスタム設定でのMP4からGIF変換
```bash
./movie-converter -input sample.mp4 -fps 15 -width 480 -speed 1 -loop -1
# 出力: sample.gif（FPS=15, 幅480px, 等倍速, ループなし）
```

### 3. GIFからMP4変換
```bash
./movie-converter -input animation.gif -fps 30
# 出力: animation.mp4（FPS=30）
```

### 4. 出力ファイル名を指定
```bash
./movie-converter -input video.mp4 -output converted_video.gif
```

## 注意事項

- ファイル名にスペースが含まれている場合、警告が表示されます
- ffmpegがシステムにインストールされている必要があります
- 入力ファイルが存在しない場合、エラーが表示されます
- サポートされていない拡張子の場合、エラーが表示されます

## エラーハンドリング

ツールは以下の場合にエラーを表示します：
- 入力ファイルが指定されていない
- 入力ファイルが存在しない
- ファイル名に拡張子が含まれていない
- サポートされていない変換（例：txt → gif）
- ffmpegの実行エラー

## 開発者向け情報

### アーキテクチャ
- `main.go`: CLIインターフェース（薄いレイヤー）
- `internal/movie_converter/usecases/services.go`: ビジネスロジック
- `internal/movie_converter/usecases/services_test.go`: テストコード

### テスト実行
```bash
cd internal/movie_converter/usecases
go test -v
```

### 依存関係
- `github.com/u2takey/ffmpeg-go`: ffmpegのGolangバインディング

## ライセンス

このプロジェクトは元のPowerShellスクリプトと同等の機能を提供することを目的としています。
