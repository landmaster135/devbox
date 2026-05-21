# Movie Extractor

動画からフレーム画像を抽出する CLI ツールです。

## ツール概要
- `operation: extract-frames` を実行します。
- `ffmpeg` を利用して、指定動画から連番 JPEG を出力します。

## フラグ一覧

| フラグ | 必須/任意 | デフォルト値 | 説明 |
|---|---|---|---|
| `-operation` | 任意 | `extract-frames` | 実行する操作。現在は `extract-frames` のみ対応 |
| `-src-file` | 必須 | なし | 入力動画ファイルのパス |
| `-fps` | 任意 | `2` | 1秒あたりに抽出するフレーム数 |
| `-quality` | 任意 | `2` | JPEG品質 (`1-31`、小さいほど高品質) |
| `-start-position` | 任意 | 空 | 抽出開始位置。秒数または `HH:MM:SS[.ms]` |
| `-out-dir` | 必須 | なし | 抽出画像の出力先ディレクトリ |
| `-help`, `-h` | 任意 | `false` | ヘルプ表示 |

## 使用方法

```bash
./movie-extractor -operation extract-frames -src-file ./sample.mp4 -out-dir ./frames
```

## 使用例

```bash
# 1fps で先頭から抽出
./movie-extractor -operation extract-frames -src-file ./sample.mp4 -out-dir ./frames

# 5fps・quality=3 で抽出
./movie-extractor -operation extract-frames -src-file ./sample.mp4 -fps 5 -quality 3 -out-dir ./frames

# 00:00:10.5 から抽出
./movie-extractor -operation extract-frames -src-file ./sample.mp4 -start-position 00:00:10.5 -out-dir ./frames
```

## 出力例

成功時:

```text
フレーム抽出が完了しました。
出力ディレクトリ: /path/to/frames
出力ファイルパターン: frame_%06d.jpg
```

エラー時（入力ファイル不存在）:

```text
エラー: 入力動画ファイルが存在しません: ./not-found.mp4
```

## 注意事項
- 実行環境に `ffmpeg` がインストールされ、`PATH` に含まれている必要があります。
- 出力ファイルは `frame_%06d.jpg` 形式で保存されます。
