# Movie Converter

GIFとMP4を相互に変換するCLIツールです。PowerShellスクリプト（Z5-5_convert_mp4_to_gif.ps1、Z5-10_convert_gif_to_mp4.ps1）と同等の機能を提供します。

## 機能

- MP4/MKV → GIF変換
- GIF → MP4変換
- 高品質なGIF生成（パレット生成とパレット使用）
- カスタマイズ可能なFPS、幅、速度、ループ設定
- 自動的な出力ファイル名生成
- **バッチ処理対応** - 複数ファイルの一括変換
- 並列処理による高速変換
- 再帰的なディレクトリ処理

## インストール

```bash
go build -o movie-converter main.go
```

## 使用方法

### 単一ファイル処理

```bash
# MP4からGIF（デフォルト設定）
./movie-converter -input video.mp4

# GIFからMP4（デフォルト設定）
./movie-converter -input animation.gif

# MP4からGIF（カスタム設定）
./movie-converter -input video.mp4 -output custom.gif -fps 30 -width 320 -speed 1.5 -loop -1

# GIFからMP4（カスタムFPS）
./movie-converter -input animation.gif -output video.mp4 -fps 24

# ヘルプを表示
./movie-converter -help
```

### バッチ処理

```bash
# MP4ファイルを一括でGIFに変換
./movie-converter -input-dir ./videos -input-ext .mp4 -output-dir ./gifs -output-ext .gif

# GIFファイルを一括でMP4に変換
./movie-converter -input-dir ./animations -input-ext .gif -output-dir ./videos -output-ext .mp4

# 再帰的にサブディレクトリも処理
./movie-converter -input-dir ./media -input-ext .mp4 -output-dir ./converted -output-ext .gif -recursive

# バッチ処理でカスタム設定
./movie-converter -input-dir ./videos -input-ext .mp4 -output-dir ./gifs -output-ext .gif -fps 15 -width 480 -speed 1
```

## オプション

### 単一ファイル処理用オプション

| オプション | 説明 | デフォルト値 |
|-----------|------|-------------|
| `-input` | 入力ファイルのパス（必須） | - |
| `-output` | 出力ファイルのパス | 自動生成 |

### バッチ処理用オプション

| オプション | 説明 | デフォルト値 |
|-----------|------|-------------|
| `-input-dir` | 入力ディレクトリのパス（必須） | - |
| `-input-ext` | 入力ファイルの拡張子（例: .mp4） | - |
| `-output-dir` | 出力ディレクトリのパス（必須） | - |
| `-output-ext` | 出力ファイルの拡張子（例: .gif） | - |
| `-recursive` | サブディレクトリも再帰的に処理するか | false |

### 共通オプション

| オプション | 説明 | デフォルト値 |
|-----------|------|-------------|
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

## 使用例

```bash
# 1. 基本的なMP4からGIF変換
./movie-converter -input sample.mp4
# 出力: sample.gif（FPS=60, 2倍速, 無限ループ）

# 2. カスタム設定でのMP4からGIF変換
./movie-converter -input sample.mp4 -fps 15 -width 480 -speed 1 -loop -1
# 出力: sample.gif（FPS=15, 幅480px, 等倍速, ループなし）

# 3. GIFからMP4変換
./movie-converter -input animation.gif -fps 30
# 出力: animation.mp4（FPS=30）

# 4. 出力ファイル名を指定
./movie-converter -input video.mp4 -output converted_video.gif

# 5. 基本的なバッチ変換
./movie-converter -input-dir ./videos -input-ext .mp4 -output-dir ./gifs -output-ext .gif
# ./videosディレクトリ内の全ての.mp4ファイルを./gifsディレクトリに.gifとして変換

# 6. 再帰的なバッチ変換
./movie-converter -input-dir ./media -input-ext .mp4 -output-dir ./converted -output-ext .gif -recursive
# ./mediaディレクトリとそのサブディレクトリ内の全ての.mp4ファイルを変換

# 7. カスタム設定でのバッチ変換
./movie-converter -input-dir ./videos -input-ext .mp4 -output-dir ./gifs -output-ext .gif -fps 15 -width 320 -speed 1
# カスタム設定（FPS=15, 幅=320px, 等倍速）でバッチ変換

# 8. GIFからMP4のバッチ変換
./movie-converter -input-dir ./animations -input-ext .gif -output-dir ./videos -output-ext .mp4 -fps 24
# 全てのGIFファイルをFPS=24でMP4に変換
```

## バッチ処理の特徴

- **並列処理**: 最大4つのファイルを同時に変換（高速化）
- **進捗表示**: 変換の進捗とファイル数を表示
- **エラー継続**: 一部のファイルが失敗しても他のファイルの変換を継続
- **詳細レポート**: 成功・失敗の詳細な結果を表示
- **ディレクトリ構造保持**: 入力ディレクトリの構造を出力ディレクトリに保持

## 注意事項

- ファイル名にスペースが含まれている場合、警告が表示されます
- ffmpegがシステムにインストールされている必要があります
- 入力ファイルが存在しない場合、エラーが表示されます
- サポートされていない拡張子の場合、エラーが表示されます
- バッチ処理では出力ディレクトリが自動的に作成されます

## エラーハンドリング

ツールは以下の場合にエラーを表示します：

### 単一ファイル処理
- 入力ファイルが指定されていない
- 入力ファイルが存在しない
- ファイル名に拡張子が含まれていない
- サポートされていない変換（例：txt → gif）
- ffmpegの実行エラー

### バッチ処理
- 入力ディレクトリが指定されていない
- 入力・出力拡張子が指定されていない
- サポートされていない拡張子
- 入力ディレクトリが存在しない
- 出力ディレクトリの作成に失敗

## 開発者向け情報

### アーキテクチャ
- `main.go`: CLIインターフェース（薄いレイヤー）
- `internal/movie_converter/usecases/services.go`: ビジネスロジック
  - 単一ファイル処理: `MovieConverterService`
  - バッチ処理: `BatchMovieConverterService`
- `internal/movie_converter/usecases/services_test.go`: テストコード

### テスト実行
```bash
cd internal/movie_converter/usecases
go test -v
```

### 依存関係
- `github.com/u2takey/ffmpeg-go`: ffmpegのGolangバインディング

### 設計原則
- **SOLID原則**: 単一責任、開放閉鎖、依存性逆転の原則を適用
- **TDD**: テスト駆動開発によるRed-Green-Refactorサイクル
- **並行処理**: Goroutineとチャネルによる安全な並列処理
- **エラーハンドリング**: 包括的なエラー処理と詳細なログ出力

## ライセンス

このプロジェクトは元のPowerShellスクリプトと同等の機能を提供することを目的としています。
