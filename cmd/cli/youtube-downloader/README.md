# YouTube動画ダウンローダー

YouTube動画やプレイリストをダウンロードするためのCLIツールです。

## 機能

- YouTube動画の単体ダウンロード
- プレイリスト全体のダウンロード
- 音声のみのダウンロード
- 品質選択（720p、1080p等）
- **音声付き高品質動画の自動結合**（FFmpeg使用）
- 並列ダウンロードによる高速化
- プログレス表示
- エラーハンドリング

## インストール

```bash
# プロジェクトルートから
go build -o youtube-downloader ./cmd/cli/youtube-downloader/
```

## オプション

| オプション | 説明 | デフォルト値 |
|-----------|------|-------------|
| `-url` | YouTube動画またはプレイリストのURL（必須） | - |
| `-output` | ダウンロード先ディレクトリ | `./downloads` |
| `-quality` | 動画品質（best, worst, 720p, 1080p等） | `best` |
| `-format` | 動画形式（mp4, webm等） | `mp4` |
| `-audio-only` | 音声のみダウンロード | `false` |
| `-playlist` | プレイリスト全体をダウンロード | `false` |
| `-max-routines` | 並列ダウンロード数 | `10` |
| `-chunk-size` | チャンクサイズ（バイト） | `10485760` (10MB) |
| `-help` | ヘルプを表示 | `false` |

## 使用例

### 1. 基本的なダウンロード

```bash
go run ./cmd/cli/youtube-downloader -url "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
```

### 2. 高品質動画のダウンロード

```bash
go run ./cmd/cli/youtube-downloader -url "https://www.youtube.com/watch?v=dQw4w9WgXcQ" -quality "1080p" -output "./videos"
```

### 3. 音楽のダウンロード

```bash
go run ./cmd/cli/youtube-downloader -url "https://www.youtube.com/watch?v=dQw4w9WgXcQ" -audio-only -output "./music"
```

### 4. プレイリストのダウンロード

```bash
go run ./cmd/cli/youtube-downloader -url "https://www.youtube.com/playlist?list=PLQZgI7en5XEgM0L1_ZcKmEzxW1sCOVZwP" -playlist -output "./playlist"
```

### 5. 高速ダウンロード設定

```bash
go run ./cmd/cli/youtube-downloader -url "https://www.youtube.com/watch?v=dQw4w9WgXcQ" -max-routines 20 -chunk-size 20971520
```

## 出力形式

### 単一動画の場合

```
ダウンロード完了:
  タイトル: Never Gonna Give You Up
  作者: Rick Astley
  再生時間: 3m32s
  品質: 720p
  ファイル: /path/to/downloads/Never Gonna Give You Up [dQw4w9WgXcQ].mp4
  サイズ: 25.67 MB
```

### プレイリストの場合

```
プレイリストダウンロード開始:
  タイトル: My Favorite Songs
  作者: User Name
  動画数: 10

(1/10) Song 1 をダウンロード中...
  完了: 001. Song 1 [VIDEO_ID1].mp4
(2/10) Song 2 をダウンロード中...
  完了: 002. Song 2 [VIDEO_ID2].mp4
...

プレイリストダウンロード完了:
  成功: 10件
  失敗: 0件
  出力先: /path/to/downloads
```

## エラーハンドリング

ツールは以下のエラー状況を適切に処理します：

- **無効なURL**: URLの形式が正しくない場合
- **動画が見つからない**: 指定された動画が存在しない場合
- **プライベート動画**: 動画がプライベート設定の場合
- **年齢制限**: 年齢制限のある動画の場合
- **ネットワークエラー**: インターネット接続の問題
- **ファイルシステムエラー**: ディスク容量不足や権限の問題

## ファイル名の規則

### 単一動画
```
{タイトル} [{動画ID}].{拡張子}
例: Never Gonna Give You Up [dQw4w9WgXcQ].mp4
```

### プレイリスト
```
{連番}. {タイトル} [{動画ID}].{拡張子}
例: 001. Never Gonna Give You Up [dQw4w9WgXcQ].mp4
```

## 技術仕様

- **言語**: Go 1.23+
- **YouTubeライブラリ**: github.com/kkdai/youtube/v2
- **FFmpegライブラリ**: github.com/u2takey/ffmpeg-go
- **並列ダウンロード**: チャンク分割による高速化
- **対応フォーマット**: MP4, WebM, M4A
- **品質**: 144p～4K（利用可能な場合）
- **音声結合**: 高品質動画の映像+音声自動結合

## 制限事項

- YouTube Premium限定コンテンツはダウンロードできません
- 地域制限のあるコンテンツは制限地域からはダウンロードできません
- 著作権で保護されたコンテンツのダウンロードは法的な問題を引き起こす可能性があります

## ライセンス

このツールは教育目的で作成されています。YouTubeの利用規約を遵守してご使用ください。

## トラブルシューティング

### よくある問題

1. **"動画が見つかりません"エラー**
  - URLが正しいか確認してください
  - 動画が削除されていないか確認してください

2. **"ネットワークエラー"**
  - インターネット接続を確認してください
  - ファイアウォールの設定を確認してください

3. **"ファイルシステムエラー"**
  - ディスク容量を確認してください
  - 出力ディレクトリの書き込み権限を確認してください

4. **ダウンロードが遅い**
  - `-max-routines`を増やしてみてください（推奨: 5-20）
  - `-chunk-size`を調整してみてください

### デバッグ

詳細なエラー情報が必要な場合は、以下のコマンドでビルドしてください：

```bash
go build -ldflags "-X main.debug=true" -o youtube-downloader ./cmd/cli/youtube-downloader/
```

## 開発者向け情報

### アーキテクチャ

```
cmd/cli/youtube-downloader/
├── main.go                    # CLIエントリーポイント
└── README.md                  # このファイル

internal/youtube_downloader/
├── config/                    # 設定とフラグ処理
│   ├── config.go
│   └── config_test.go
├── domain/                    # ドメインモデル
│   └── models.go
└── usecases/                  # ビジネスロジック
    ├── services.go
    └── services_test.go
```

### テスト実行

```bash
# 全テスト実行
go test ./internal/youtube_downloader/... -v

# 設定テストのみ
go test ./internal/youtube_downloader/config/... -v

# カバレッジ付きテスト
go test -cover ./internal/youtube_downloader/...
```

### ビルド

```bash
# 開発用ビルド
go build -o youtube-downloader ./cmd/cli/youtube-downloader/

# リリース用ビルド
go build -ldflags "-s -w" -o youtube-downloader ./cmd/cli/youtube-downloader/
