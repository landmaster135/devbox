# image-converter-by-libwebp

`image-converter-by-libwebp` は、OS にインストールされた libwebp の `cwebp` を使って複数画像を WebP へ変換する CLI ツールです。

## 前提

`cwebp` が `PATH` から実行できる必要があります。見つからない場合はエラー終了します。

## フラグ一覧

| フラグ | 必須 | デフォルト値 | 説明 |
|---|---|---|---|
| `--src-dir` | 任意 | `.` | 入力画像を探索するディレクトリ |
| `--out-dir` | 任意 | `./999_converted_images` | 変換後画像の出力ディレクトリ |
| `--archive-dir` | 任意 | 空文字 | 処理済み元ファイルの退避先ディレクトリ。空なら退避しない |
| `--move` | 任意 | `false` | 退避時にコピーではなく移動する |
| `--ext` | 任意 | `webp` | 出力形式。対応形式は `webp` |
| `--q` | 任意 | `99` | WebP 品質。`1` から `100` |
| `--m` | 任意 | `4` | `cwebp -m` に渡す圧縮メソッド。`0` から `6` |
| `--workers` | 任意 | CPU 数 | 並列ワーカー数 |
| `--recursive`, `-R` | 任意 | `false` | サブディレクトリを再帰的に走査する |
| `--lossless` | 任意 | `false` | `cwebp -lossless` を有効にする |
| `--help`, `-h` | 任意 | `false` | ヘルプを表示する |

## 使用方法

```bash
go run ./cmd/cli/image-converter-by-libwebp \
  --src-dir ./photos \
  --out-dir ./webp \
  --q 99 \
  --m 6
```

## 使用例

再帰的に変換する:

```bash
go run ./cmd/cli/image-converter-by-libwebp \
  --src-dir ./photos \
  --out-dir ./webp \
  --recursive
```

変換済みの元画像を移動する:

```bash
go run ./cmd/cli/image-converter-by-libwebp \
  --src-dir ./photos \
  --out-dir ./webp \
  --archive-dir ./archive \
  --move
```

## cwebp オプション

変換時は常に以下のオプションを付与します。

```text
-preset photo -metadata icc -sharp_yuv -progress -short
```

`--m` は `cwebp -m` として渡します。デフォルトは `4` です。

`--lossless` 指定時は、上記に加えて `-lossless` を付与します。

## 出力例

成功時:

```text
画像変換が完了しました
  成功: 12 ファイル
  失敗: 0 ファイル
  出力先: /path/to/webp
```

エラー時:

```text
エラー: libwebp パッケージが見つかりません: cwebp をインストールしてください
```
