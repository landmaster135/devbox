# YouTube動画ダウンロード実装ガイド

## 概要

このガイドでは、`~/package_references/youtube` ライブラリを使用したYouTube動画ダウンロードの実装方法について詳しく説明します。

## 基本的なアーキテクチャ

### 主要コンポーネント

```
package_references/youtube/
├── client.go              # YouTubeとの通信を担当するクライアント
├── video.go               # 動画情報の解析とデータ構造
├── downloader/
│   └── downloader.go      # 実際のダウンロード処理
├── decipher.go            # 暗号化されたURLの復号化
├── format_list.go         # 利用可能な動画フォーマットの管理
├── playlist.go            # プレイリスト処理
└── cmd/youtubedr/         # CLIツール実装
```

### 主要な構造体

```go
// クライアント
type Client struct {
  HTTPClient   *http.Client
  MaxRoutines  int
  ChunkSize    int64
  playerCache  playerCache
  client       *clientInfo
}

// 動画情報
type Video struct {
  ID              string
  Title           string
  Description     string
  Author          string
  Duration        time.Duration
  Formats         FormatList
  Thumbnails      Thumbnails
}

// フォーマット情報
type Format struct {
  ItagNo        int
  URL           string
  MimeType      string
  Quality       string
  Bitrate       int
  ContentLength int64
}
```

## 動画ダウンロードの流れ

### ステップ1: クライアントの初期化

```go
client := youtube.Client{
    HTTPClient:  http.DefaultClient,
    MaxRoutines: 10,           // 並列ダウンロード数
    ChunkSize:   10 * 1024 * 1024, // 10MB チャンク
}
```

### ステップ2: 動画情報の取得

```go
// URLまたは動画IDから動画情報を取得
video, err := client.GetVideo("https://www.youtube.com/watch?v=VIDEO_ID")
if err != nil {
    return fmt.Errorf("動画情報の取得に失敗: %w", err)
}

fmt.Printf("タイトル: %s\n", video.Title)
fmt.Printf("作者: %s\n", video.Author)
fmt.Printf("再生時間: %s\n", video.Duration)
```

### ステップ3: フォーマットの選択

```go
// 音声付きフォーマットのみを取得
formats := video.Formats.WithAudioChannels()

// 品質でフィルタリング
hdFormats := video.Formats.Quality("hd720")

// MIMEタイプでフィルタリング
mp4Formats := video.Formats.Type("video/mp4")

// 最高品質を選択
if len(formats) > 0 {
  selectedFormat := &formats[0]
  fmt.Printf("選択されたフォーマット: %s (%s)\n", 
    selectedFormat.Quality, selectedFormat.MimeType)
}
```

### ステップ4: ダウンロード実行

```go
// ストリームを取得
stream, size, err := client.GetStream(video, &formats[0])
if err != nil {
  return fmt.Errorf("ストリーム取得に失敗: %w", err)
}
defer stream.Close()

// ファイルに保存
file, err := os.Create("video.mp4")
if err != nil {
  return fmt.Errorf("ファイル作成に失敗: %w", err)
}
defer file.Close()

// ダウンロード実行
_, err = io.Copy(file, stream)
if err != nil {
  return fmt.Errorf("ダウンロードに失敗: %w", err)
}
```

## 技術的な実装詳細

### YouTube Innertube APIの活用

このライブラリは、YouTubeの内部API（Innertube API）を使用して動画情報を取得します。

```go
// 複数のクライアント情報を使い分けて制限を回避
var (
  WebClient = clientInfo{
    name:      "WEB",
    version:   "2.20220801.00.00",
    key:       "AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8",
    userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
  }
  
  AndroidClient = clientInfo{
    name:           "ANDROID",
    version:        "18.11.34",
    key:            "AIzaSyA8eiZmM1FaDVjRy-df2KTyQ_vz_yYM39w",
    userAgent:      "com.google.android.youtube/18.11.34 (Linux; U; Android 11)",
    androidVersion: 30,
  }
  
  IOSClient = clientInfo{
    name:        "IOS",
    version:     "19.45.4",
    key:         "AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8",
    userAgent:   "com.google.ios.youtube/19.45.4 (iPhone16,2; U; CPU iOS 18_1_0)",
    deviceModel: "iPhone16,2",
  }
)
```

### チャンク分割ダウンロード

大きなファイルを効率的にダウンロードするため、ファイルを複数のチャンクに分割して並列ダウンロードします。

```go
func (c *Client) downloadChunked(ctx context.Context, req *http.Request, w *io.PipeWriter, format *Format) {
  chunks := getChunks(format.ContentLength, c.getChunkSize())
  maxRoutines := c.getMaxRoutines(len(chunks))
  
  currentChunk := atomic.Uint32{}
  for i := 0; i < maxRoutines; i++ {
    go func() {
      for {
        chunkIndex := int(currentChunk.Add(1)) - 1
        if chunkIndex >= len(chunks) {
          return
        }
        
        chunk := &chunks[chunkIndex]
        err := c.downloadChunk(req.Clone(ctx), chunk)
        if err != nil {
          // エラーハンドリング
          return
        }
      }
    }()
  }
}
```

### 制限回避の仕組み

YouTubeの様々な制限を回避するための仕組み：

1. **年齢制限動画**: EmbeddedClientを使用
2. **地域制限**: 異なるクライアント情報で再試行
3. **プライベート動画**: 適切なエラーハンドリング

```go
func (c *Client) videoFromID(ctx context.Context, id string) (*Video, error) {
  body, err := c.videoDataByInnertube(ctx, id)
  if err != nil {
    return nil, err
  }
  
  v := Video{ID: id}
  
  if err = v.parseVideoInfo(body); err == nil {
    return &v, nil
  }
  
  // 埋め込み無効の場合
  if errors.Is(err, ErrNotPlayableInEmbed) {
    html, err := c.httpGetBodyBytes(ctx, "https://www.youtube.com/watch?v="+id)
    if err != nil {
      return nil, err
    }
    return &v, v.parseVideoPage(html)
  }
  
  // ログイン必須の場合（年齢制限）
  if errors.Is(err, ErrLoginRequired) {
    c.client = &EmbeddedClient
    bodyEmbed, errEmbed := c.videoDataByInnertube(ctx, id)
    if errEmbed == nil {
      errEmbed = v.parseVideoInfo(bodyEmbed)
    }
    return &v, errEmbed
  }
  
  return &v, err
}
```

## 高品質動画のダウンロード（Composite）

1080p以上の高品質動画は、音声と映像が分離されているため、別々にダウンロードしてFFmpegで結合する必要があります。

```go
func (dl *Downloader) DownloadComposite(ctx context.Context, outputFile string, v *youtube.Video, quality string, mimetype, language string) error {
  // 映像と音声のフォーマットを取得
  videoFormat, audioFormat, err := getVideoAudioFormats(v, quality, mimetype, language)
  if err != nil {
    return err
  }
  
  // 一時ファイルを作成
  videoFile, err := os.CreateTemp(outputDir, "youtube_*.m4v")
  if err != nil {
      return err
  }
  defer os.Remove(videoFile.Name())
  
  audioFile, err := os.CreateTemp(outputDir, "youtube_*.m4a")
  if err != nil {
    return err
  }
  defer os.Remove(audioFile.Name())
  
  // 映像をダウンロード
  err = dl.videoDLWorker(ctx, videoFile, v, videoFormat)
  if err != nil {
    return err
  }
  
  // 音声をダウンロード
  err = dl.videoDLWorker(ctx, audioFile, v, audioFormat)
  if err != nil {
    return err
  }
  
  // FFmpegで結合
  ffmpegCmd := exec.Command("ffmpeg", "-y",
    "-i", videoFile.Name(),
    "-i", audioFile.Name(),
    "-c", "copy",           // 再エンコードなし
    "-shortest",            // 短い方に合わせる
    destFile,
    "-loglevel", "warning",
  )
  
  return ffmpegCmd.Run()
}
```

## エラーハンドリング

様々なエラー状況に対する適切な処理：

```go
func (v *Video) isVideoDownloadable(prData playerResponseData, isVideoPage bool) error {
  switch prData.PlayabilityStatus.Status {
  case "OK":
    return nil
  case "LOGIN_REQUIRED":
    if strings.HasPrefix(prData.PlayabilityStatus.Reason, "This video is private") {
      return ErrVideoPrivate
    }
    return ErrLoginRequired
  case "UNPLAYABLE":
    return ErrVideoUnplayable
  case "ERROR":
    return ErrVideoNotFound
  }
  
  if !isVideoPage && !prData.PlayabilityStatus.PlayableInEmbed {
    return ErrNotPlayableInEmbed
  }
  
  return &ErrPlayabiltyStatus{
    Status: prData.PlayabilityStatus.Status,
    Reason: prData.PlayabilityStatus.Reason,
  }
}
```

## 実装のベストプラクティス

### 1. パフォーマンス最適化

```go
// 並列ダウンロード設定
client := youtube.Client{
  MaxRoutines: 10,                    // 並列数
  ChunkSize:   10 * 1024 * 1024,     // 10MBチャンク
}

// HTTPクライアントの再利用
client.HTTPClient = &http.Client{
  Timeout: 30 * time.Second,
  Transport: &http.Transport{
    MaxIdleConns:        100,
    MaxIdleConnsPerHost: 10,
  },
}
```

### 2. セキュリティ対策

```go
// 適切なUser-Agentの設定
req.Header.Set("User-Agent", c.client.userAgent)
req.Header.Set("Origin", "https://youtube.com")
req.Header.Set("Sec-Fetch-Mode", "navigate")

// Cookieによる同意処理
req.AddCookie(&http.Cookie{
  Name:   "CONSENT",
  Value:  "YES+cb.20210328-17-p0.en+FX+" + c.consentID,
  Path:   "/",
  Domain: ".youtube.com",
})
```

### 3. プログレス表示

```go
// プログレスバーの実装
progress := mpb.New(mpb.WithWidth(64))
bar := progress.AddBar(
  int64(contentLength),
  mpb.PrependDecorators(
    decor.CountersKibiByte("% .2f / % .2f"),
    decor.Percentage(decor.WCSyncSpace),
  ),
  mpb.AppendDecorators(
    decor.EwmaETA(decor.ET_STYLE_GO, 90),
    decor.EwmaSpeed(decor.UnitKiB, "% .2f", 60),
  ),
)

reader := bar.ProxyReader(stream)
_, err = io.Copy(file, reader)
```

## 使用例

### 基本的な使用例

```go
package main

import (
  "fmt"
  "io"
  "os"
  
  "github.com/kkdai/youtube/v2"
)

func main() {
  videoID := "BaW_jenozKc"
  client := youtube.Client{}
  
  // 動画情報を取得
  video, err := client.GetVideo(videoID)
  if err != nil {
    panic(err)
  }
  
  // 音声付きフォーマットを選択
  formats := video.Formats.WithAudioChannels()
  
  // ストリームを取得
  stream, _, err := client.GetStream(video, &formats[0])
  if err != nil {
    panic(err)
  }
  defer stream.Close()
  
  // ファイルに保存
  file, err := os.Create("video.mp4")
  if err != nil {
    panic(err)
  }
  defer file.Close()
  
  _, err = io.Copy(file, stream)
  if err != nil {
    panic(err)
  }
  
  fmt.Println("ダウンロード完了!")
}
```

### プレイリストのダウンロード

```go
func downloadPlaylist() {
  playlistID := "PLQZgI7en5XEgM0L1_ZcKmEzxW1sCOVZwP"
  client := youtube.Client{}
  
  // プレイリスト情報を取得
  playlist, err := client.GetPlaylist(playlistID)
  if err != nil {
    panic(err)
  }
  
  fmt.Printf("プレイリスト: %s by %s\n", playlist.Title, playlist.Author)
  
  // 各動画をダウンロード
  for i, entry := range playlist.Videos {
    video, err := client.VideoFromPlaylistEntry(entry)
    if err != nil {
      fmt.Printf("動画 %d の取得に失敗: %v\n", i+1, err)
      continue
    }
    
    // ダウンロード処理...
    fmt.Printf("(%d) %s をダウンロード中...\n", i+1, video.Title)
  }
}
```

## まとめ

このライブラリは、YouTubeの内部APIを巧妙に活用し、様々な制限を回避しながら安定した動画ダウンロード機能を提供します。主な特徴：

- **高性能**: 並列チャンクダウンロードによる高速化
- **柔軟性**: 複数のクライアント情報による制限回避
- **安定性**: 適切なエラーハンドリングと再試行機能
- **拡張性**: インターフェースベースの設計

適切に実装することで、効率的で安定したYouTube動画ダウンローダーを構築できます。
